package crawler_test

import (
	"compress/gzip"
	"context"
	"fmt"
	"html/template"
	"io"
	stdlog "log"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/go-log/log/print"
	crawler "github.com/thewizardplusplus/go-crawler"
	"github.com/thewizardplusplus/go-crawler/checkers"
	"github.com/thewizardplusplus/go-crawler/extractors"
	"github.com/thewizardplusplus/go-crawler/extractors/transformers"
	"github.com/thewizardplusplus/go-crawler/handlers"
	"github.com/thewizardplusplus/go-crawler/models"
	"github.com/thewizardplusplus/go-crawler/registers"
	"github.com/thewizardplusplus/go-crawler/registers/sitemap"
	urlutils "github.com/thewizardplusplus/go-crawler/url-utils"
	htmlselector "github.com/thewizardplusplus/go-html-selector"
)

type LinkHandler struct {
	Name      string
	ServerURL string
}

func (handler LinkHandler) HandleLink(
	ctx context.Context,
	link models.SourcedLink,
) {
	var prefix string
	if handler.Name != "" {
		prefix = fmt.Sprintf("[%s] ", handler.Name)
	}

	fmt.Printf(
		"%sreceived link %q from page %q\n",
		prefix,
		handler.replaceServerURL(link.Link),
		handler.replaceServerURL(link.SourceLink),
	)
}

// replace the test server URL for reproducibility of the example
func (handler LinkHandler) replaceServerURL(link string) string {
	return strings.Replace(link, handler.ServerURL, "http://example.com", -1)
}

// nolint: gocyclo
func RunServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(
		writer http.ResponseWriter,
		request *http.Request,
	) {
		if request.URL.Path == "/robots.txt" {
			sitemapLink :=
				completeLinkWithHost("/sitemap_from_robots_txt.xml", request.Host)
			fmt.Fprintf(writer, `
				User-agent: go-crawler
				Disallow: /2

				Sitemap: %s
			`, sitemapLink)

			return
		}

		var links []string
		switch request.URL.Path {
		case "/sitemap.xml":
			links = []string{"/1", "/2", "/hidden/1", "/hidden/2"}
		case "/sitemap_from_robots_txt.xml":
			links = []string{"/hidden/3", "/hidden/4"}
		case "/hidden/1/sitemap.xml":
			links = []string{"/hidden/5", "/hidden/6"}
		case "/1/sitemap.xml", "/2/sitemap.xml", "/hidden/sitemap.xml":
			// render the empty Sitemap to escape the error logging
			// for reproducibility of the example
			links = []string{}
		}
		for index := range links {
			links[index] = completeLinkWithHost(links[index], request.Host)
		}

		if links != nil {
			writer.Header().Set("Content-Encoding", "gzip")

			compressingWriter := gzip.NewWriter(writer)
			defer compressingWriter.Close() // nolint: errcheck

			renderTemplate(compressingWriter, links, `
				<?xml version="1.0" encoding="UTF-8" ?>
				<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
					{{ range $link := . }}
						<url>
							<loc>{{ $link }}</loc>
						</url>
					{{ end }}
				</urlset>
			`)

			return
		}

		switch request.URL.Path {
		case "/":
			links = []string{"/1", "/2", "/2", "https://golang.org/"}
		case "/1":
			links = []string{"/1/1", "/1/2"}
		case "/2":
			links = []string{"/2/1", "/2/2"}
		case "/hidden/1":
			links = []string{"/hidden/1/test"}
		}

		renderTemplate(writer, links, `
			<ul>
				{{ range $link := . }}
					<li>
						<a href="{{ $link }}">{{ $link }}</a>
					</li>
				{{ end }}
			</ul>
		`)
	}))
}

func completeLinkWithHost(link string, host string) string {
	return "http://" + path.Join(host, link)
}

func renderTemplate(writer io.Writer, data interface{}, text string) {
	template, _ := template.New("").Parse(text) // nolint: errcheck
	template.Execute(writer, data)              // nolint: errcheck
}

func ExampleCrawl() {
	server := RunServer()
	defer server.Close()

	logger := stdlog.New(os.Stderr, "", stdlog.LstdFlags|stdlog.Lmicroseconds)
	// wrap the standard logger via the github.com/go-log/log package
	wrappedLogger := print.New(logger)

	crawler.Crawl(
		context.Background(),
		crawler.ConcurrencyConfig{
			ConcurrencyFactor: runtime.NumCPU(),
			BufferSize:        1000,
		},
		[]string{server.URL},
		crawler.CrawlDependencies{
			LinkExtractor: extractors.DefaultExtractor{
				HTTPClient: http.DefaultClient,
				Filters: htmlselector.OptimizeFilters(htmlselector.FilterGroup{
					"a": {"href"},
				}),
				LinkTransformer: transformers.ResolvingTransformer{
					BaseTagSelection: transformers.SelectFirstBaseTag,
					BaseTagFilters:   transformers.DefaultBaseTagFilters,
					BaseHeaderNames:  urlutils.DefaultBaseHeaderNames,
					Logger:           wrappedLogger,
				},
			},
			LinkChecker: checkers.HostChecker{
				ComparisonResult: urlutils.Same,
				Logger:           wrappedLogger,
			},
			LinkHandler: LinkHandler{
				ServerURL: server.URL,
			},
			Logger: wrappedLogger,
		},
	)

	// Unordered output:
	// received link "http://example.com/1" from page "http://example.com"
	// received link "http://example.com/1/1" from page "http://example.com/1"
	// received link "http://example.com/1/2" from page "http://example.com/1"
	// received link "http://example.com/2" from page "http://example.com"
	// received link "http://example.com/2" from page "http://example.com"
	// received link "http://example.com/2/1" from page "http://example.com/2"
	// received link "http://example.com/2/1" from page "http://example.com/2"
	// received link "http://example.com/2/2" from page "http://example.com/2"
	// received link "http://example.com/2/2" from page "http://example.com/2"
	// received link "https://golang.org/" from page "http://example.com"
}

func ExampleCrawl_withoutDuplicatesOnExtracting() {
	server := RunServer()
	defer server.Close()

	logger := stdlog.New(os.Stderr, "", stdlog.LstdFlags|stdlog.Lmicroseconds)
	// wrap the standard logger via the github.com/go-log/log package
	wrappedLogger := print.New(logger)

	crawler.Crawl(
		context.Background(),
		crawler.ConcurrencyConfig{
			ConcurrencyFactor: runtime.NumCPU(),
			BufferSize:        1000,
		},
		[]string{server.URL},
		crawler.CrawlDependencies{
			LinkExtractor: extractors.DefaultExtractor{
				HTTPClient: http.DefaultClient,
				Filters: htmlselector.OptimizeFilters(htmlselector.FilterGroup{
					"a": {"href"},
				}),
				LinkTransformer: transformers.ResolvingTransformer{
					BaseTagSelection: transformers.SelectFirstBaseTag,
					BaseTagFilters:   transformers.DefaultBaseTagFilters,
					BaseHeaderNames:  urlutils.DefaultBaseHeaderNames,
					Logger:           wrappedLogger,
				},
			},
			LinkChecker: checkers.CheckerGroup{
				checkers.HostChecker{
					ComparisonResult: urlutils.Same,
					Logger:           wrappedLogger,
				},
				checkers.DuplicateChecker{
					LinkRegister: registers.NewLinkRegister(urlutils.SanitizeLink),
					Logger:       wrappedLogger,
				},
			},
			LinkHandler: LinkHandler{
				ServerURL: server.URL,
			},
			Logger: wrappedLogger,
		},
	)

	// Unordered output:
	// received link "http://example.com/1" from page "http://example.com"
	// received link "http://example.com/1/1" from page "http://example.com/1"
	// received link "http://example.com/1/2" from page "http://example.com/1"
	// received link "http://example.com/2" from page "http://example.com"
	// received link "http://example.com/2" from page "http://example.com"
	// received link "http://example.com/2/1" from page "http://example.com/2"
	// received link "http://example.com/2/2" from page "http://example.com/2"
	// received link "https://golang.org/" from page "http://example.com"
}

func ExampleCrawl_withoutDuplicatesOnHandling() {
	server := RunServer()
	defer server.Close()

	logger := stdlog.New(os.Stderr, "", stdlog.LstdFlags|stdlog.Lmicroseconds)
	// wrap the standard logger via the github.com/go-log/log package
	wrappedLogger := print.New(logger)

	crawler.Crawl(
		context.Background(),
		crawler.ConcurrencyConfig{
			ConcurrencyFactor: runtime.NumCPU(),
			BufferSize:        1000,
		},
		[]string{server.URL},
		crawler.CrawlDependencies{
			LinkExtractor: extractors.DefaultExtractor{
				HTTPClient: http.DefaultClient,
				Filters: htmlselector.OptimizeFilters(htmlselector.FilterGroup{
					"a": {"href"},
				}),
				LinkTransformer: transformers.ResolvingTransformer{
					BaseTagSelection: transformers.SelectFirstBaseTag,
					BaseTagFilters:   transformers.DefaultBaseTagFilters,
					BaseHeaderNames:  urlutils.DefaultBaseHeaderNames,
					Logger:           wrappedLogger,
				},
			},
			LinkChecker: checkers.CheckerGroup{
				checkers.HostChecker{
					ComparisonResult: urlutils.Same,
					Logger:           wrappedLogger,
				},
				checkers.DuplicateChecker{
					LinkRegister: registers.NewLinkRegister(urlutils.SanitizeLink),
					Logger:       wrappedLogger,
				},
			},
			LinkHandler: handlers.CheckedHandler{
				LinkChecker: checkers.DuplicateChecker{
					// don't use here the link register from the duplicate checker above
					LinkRegister: registers.NewLinkRegister(urlutils.SanitizeLink),
					Logger:       wrappedLogger,
				},
				LinkHandler: LinkHandler{
					ServerURL: server.URL,
				},
			},
			Logger: wrappedLogger,
		},
	)

	// Unordered output:
	// received link "http://example.com/1" from page "http://example.com"
	// received link "http://example.com/1/1" from page "http://example.com/1"
	// received link "http://example.com/1/2" from page "http://example.com/1"
	// received link "http://example.com/2" from page "http://example.com"
	// received link "http://example.com/2/1" from page "http://example.com/2"
	// received link "http://example.com/2/2" from page "http://example.com/2"
	// received link "https://golang.org/" from page "http://example.com"
}

func ExampleCrawl_withRobotsTXT() {
	server := RunServer()
	defer server.Close()

	logger := stdlog.New(os.Stderr, "", stdlog.LstdFlags|stdlog.Lmicroseconds)
	// wrap the standard logger via the github.com/go-log/log package
	wrappedLogger := print.New(logger)

	crawler.Crawl(
		context.Background(),
		crawler.ConcurrencyConfig{
			ConcurrencyFactor: runtime.NumCPU(),
			BufferSize:        1000,
		},
		[]string{server.URL},
		crawler.CrawlDependencies{
			LinkExtractor: extractors.DefaultExtractor{
				HTTPClient: http.DefaultClient,
				Filters: htmlselector.OptimizeFilters(htmlselector.FilterGroup{
					"a": {"href"},
				}),
				LinkTransformer: transformers.ResolvingTransformer{
					BaseTagSelection: transformers.SelectFirstBaseTag,
					BaseTagFilters:   transformers.DefaultBaseTagFilters,
					BaseHeaderNames:  urlutils.DefaultBaseHeaderNames,
					Logger:           wrappedLogger,
				},
			},
			LinkChecker: checkers.CheckerGroup{
				checkers.HostChecker{
					ComparisonResult: urlutils.Same,
					Logger:           wrappedLogger,
				},
				checkers.RobotsTXTChecker{
					UserAgent:         "go-crawler",
					RobotsTXTRegister: registers.NewRobotsTXTRegister(http.DefaultClient),
					Logger:            wrappedLogger,
				},
			},
			LinkHandler: LinkHandler{
				ServerURL: server.URL,
			},
			Logger: wrappedLogger,
		},
	)

	// Unordered output:
	// received link "http://example.com/1" from page "http://example.com"
	// received link "http://example.com/1/1" from page "http://example.com/1"
	// received link "http://example.com/1/2" from page "http://example.com/1"
	// received link "http://example.com/2" from page "http://example.com"
	// received link "http://example.com/2" from page "http://example.com"
	// received link "https://golang.org/" from page "http://example.com"
}

func ExampleCrawl_withSitemap() {
	server := RunServer()
	defer server.Close()

	logger := stdlog.New(os.Stderr, "", stdlog.LstdFlags|stdlog.Lmicroseconds)
	// wrap the standard logger via the github.com/go-log/log package
	wrappedLogger := print.New(logger)

	crawler.Crawl(
		context.Background(),
		crawler.ConcurrencyConfig{
			ConcurrencyFactor: runtime.NumCPU(),
			BufferSize:        1000,
		},
		[]string{server.URL},
		crawler.CrawlDependencies{
			LinkExtractor: extractors.ExtractorGroup{
				Name: "main extractors",
				LinkExtractors: []models.LinkExtractor{
					extractors.DefaultExtractor{
						HTTPClient: http.DefaultClient,
						Filters: htmlselector.OptimizeFilters(htmlselector.FilterGroup{
							"a": {"href"},
						}),
						LinkTransformer: transformers.ResolvingTransformer{
							BaseTagSelection: transformers.SelectFirstBaseTag,
							BaseTagFilters:   transformers.DefaultBaseTagFilters,
							BaseHeaderNames:  urlutils.DefaultBaseHeaderNames,
							Logger:           wrappedLogger,
						},
					},
					extractors.SitemapExtractor{
						SitemapRegister: registers.NewSitemapRegister(
							time.Second,
							extractors.ExtractorGroup{
								Name: "extractors of Sitemap links",
								LinkExtractors: []models.LinkExtractor{
									sitemap.HierarchicalGenerator{
										SanitizeLink: urlutils.SanitizeLink,
										MaximalDepth: -1,
									},
									sitemap.RobotsTXTGenerator{
										RobotsTXTRegister: registers.NewRobotsTXTRegister(
											http.DefaultClient,
										),
									},
								},
								Logger: wrappedLogger,
							},
							wrappedLogger,
							sitemap.Loader{HTTPClient: http.DefaultClient}.LoadLink,
						),
						Logger: wrappedLogger,
					},
				},
				Logger: wrappedLogger,
			},
			LinkChecker: checkers.CheckerGroup{
				checkers.HostChecker{
					ComparisonResult: urlutils.Same,
					Logger:           wrappedLogger,
				},
				checkers.DuplicateChecker{
					LinkRegister: registers.NewLinkRegister(urlutils.SanitizeLink),
					Logger:       wrappedLogger,
				},
			},
			LinkHandler: handlers.CheckedHandler{
				LinkChecker: checkers.DuplicateChecker{
					// don't use here the link register from the duplicate checker above
					LinkRegister: registers.NewLinkRegister(urlutils.SanitizeLink),
					Logger:       wrappedLogger,
				},
				LinkHandler: LinkHandler{
					ServerURL: server.URL,
				},
			},
			Logger: wrappedLogger,
		},
	)

	// Unordered output:
	// received link "http://example.com/1" from page "http://example.com"
	// received link "http://example.com/1/1" from page "http://example.com/1"
	// received link "http://example.com/1/2" from page "http://example.com/1"
	// received link "http://example.com/2" from page "http://example.com"
	// received link "http://example.com/2/1" from page "http://example.com/2"
	// received link "http://example.com/2/2" from page "http://example.com/2"
	// received link "http://example.com/hidden/1" from page "http://example.com"
	// received link "http://example.com/hidden/1/test" from page "http://example.com/hidden/1"
	// received link "http://example.com/hidden/2" from page "http://example.com"
	// received link "http://example.com/hidden/3" from page "http://example.com"
	// received link "http://example.com/hidden/4" from page "http://example.com"
	// received link "http://example.com/hidden/5" from page "http://example.com/hidden/1/test"
	// received link "http://example.com/hidden/6" from page "http://example.com/hidden/1/test"
	// received link "https://golang.org/" from page "http://example.com"
}

func ExampleCrawl_withFewHandlers() {
	server := RunServer()
	defer server.Close()

	logger := stdlog.New(os.Stderr, "", stdlog.LstdFlags|stdlog.Lmicroseconds)
	// wrap the standard logger via the github.com/go-log/log package
	wrappedLogger := print.New(logger)

	crawler.Crawl(
		context.Background(),
		crawler.ConcurrencyConfig{
			ConcurrencyFactor: runtime.NumCPU(),
			BufferSize:        1000,
		},
		[]string{server.URL},
		crawler.CrawlDependencies{
			LinkExtractor: extractors.DefaultExtractor{
				HTTPClient: http.DefaultClient,
				Filters: htmlselector.OptimizeFilters(htmlselector.FilterGroup{
					"a": {"href"},
				}),
				LinkTransformer: transformers.ResolvingTransformer{
					BaseTagSelection: transformers.SelectFirstBaseTag,
					BaseTagFilters:   transformers.DefaultBaseTagFilters,
					BaseHeaderNames:  urlutils.DefaultBaseHeaderNames,
					Logger:           wrappedLogger,
				},
			},
			LinkChecker: checkers.HostChecker{
				ComparisonResult: urlutils.Same,
				Logger:           wrappedLogger,
			},
			LinkHandler: handlers.HandlerGroup{
				handlers.CheckedHandler{
					LinkChecker: checkers.HostChecker{
						ComparisonResult: urlutils.Same,
						Logger:           wrappedLogger,
					},
					LinkHandler: LinkHandler{
						Name:      "inner",
						ServerURL: server.URL,
					},
				},
				handlers.CheckedHandler{
					LinkChecker: checkers.HostChecker{
						ComparisonResult: urlutils.Different,
						Logger:           wrappedLogger,
					},
					LinkHandler: LinkHandler{
						Name:      "outer",
						ServerURL: server.URL,
					},
				},
			},
			Logger: wrappedLogger,
		},
	)

	// Unordered output:
	// [inner] received link "http://example.com/1" from page "http://example.com"
	// [inner] received link "http://example.com/1/1" from page "http://example.com/1"
	// [inner] received link "http://example.com/1/2" from page "http://example.com/1"
	// [inner] received link "http://example.com/2" from page "http://example.com"
	// [inner] received link "http://example.com/2" from page "http://example.com"
	// [inner] received link "http://example.com/2/1" from page "http://example.com/2"
	// [inner] received link "http://example.com/2/1" from page "http://example.com/2"
	// [inner] received link "http://example.com/2/2" from page "http://example.com/2"
	// [inner] received link "http://example.com/2/2" from page "http://example.com/2"
	// [outer] received link "https://golang.org/" from page "http://example.com"
}
