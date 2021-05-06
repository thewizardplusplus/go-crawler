package crawler_test

import (
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
	"sync"
	"time"

	"github.com/go-log/log/print"
	crawler "github.com/thewizardplusplus/go-crawler"
	"github.com/thewizardplusplus/go-crawler/checkers"
	"github.com/thewizardplusplus/go-crawler/extractors"
	"github.com/thewizardplusplus/go-crawler/handlers"
	"github.com/thewizardplusplus/go-crawler/models"
	"github.com/thewizardplusplus/go-crawler/registers"
	"github.com/thewizardplusplus/go-crawler/registers/sitemap"
	"github.com/thewizardplusplus/go-crawler/sanitizing"
	htmlselector "github.com/thewizardplusplus/go-html-selector"
)

type LinkHandler struct {
	ServerURL string
}

func (handler LinkHandler) HandleLink(
	ctx context.Context,
	link models.SourcedLink,
) {
	fmt.Printf(
		"have got the link %q from the page %q\n",
		handler.replaceServerURL(link.Link),
		handler.replaceServerURL(link.SourceLink),
	)
}

// replace the test server URL for reproducibility of the example
func (handler LinkHandler) replaceServerURL(link string) string {
	return strings.Replace(link, handler.ServerURL, "http://example.com", -1)
}

func RunServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(
		writer http.ResponseWriter,
		request *http.Request,
	) {
		if request.URL.Path == "/robots.txt" {
			// nolint: errcheck
			fmt.Fprint(writer, `
				User-agent: go-crawler
				Disallow: /2
			`)

			return
		}

		if request.URL.Path == "/sitemap.xml" {
			links := []string{"/1", "/2", "/hidden/1", "/hidden/2"}
			completeLinksWithHost(links, request.Host)

			// nolint: errcheck
			renderTemplate(writer, links, `
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

		var links []string
		switch request.URL.Path {
		case "/":
			links = []string{"/1", "/2", "/2", "https://golang.org/"}
		case "/1":
			links = []string{"/1/1", "/1/2"}
		case "/2":
			links = []string{"/2/1", "/2/2"}
		}
		completeLinksWithHost(links, request.Host)

		// nolint: errcheck
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

func completeLinksWithHost(links []string, host string) {
	for index := range links {
		if strings.HasPrefix(links[index], "/") {
			links[index] = "http://" + path.Join(host, links[index])
		}
	}
}

// nolint: unparam
func renderTemplate(writer io.Writer, data interface{}, text string) error {
	template, err := template.New("").Parse(text)
	if err != nil {
		return err
	}

	return template.Execute(writer, data)
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
			LinkExtractor: extractors.RepeatingExtractor{
				LinkExtractor: extractors.NewDelayingExtractor(
					time.Second,
					time.Sleep,
					extractors.DefaultExtractor{
						HTTPClient: http.DefaultClient,
						Filters: htmlselector.OptimizeFilters(htmlselector.FilterGroup{
							"a": {"href"},
						}),
					},
				),
				RepeatCount:  5,
				RepeatDelay:  0,
				Logger:       wrappedLogger,
				SleepHandler: time.Sleep,
			},
			LinkChecker: checkers.CheckerGroup{
				checkers.HostChecker{
					Logger: wrappedLogger,
				},
				checkers.DuplicateChecker{
					LinkRegister: registers.NewLinkRegister(
						sanitizing.SanitizeLink,
						wrappedLogger,
					),
				},
			},
			LinkHandler: handlers.CheckedHandler{
				LinkChecker: checkers.DuplicateChecker{
					// don't use here the link register from the duplicate checker above
					LinkRegister: registers.NewLinkRegister(
						sanitizing.SanitizeLink,
						wrappedLogger,
					),
				},
				LinkHandler: LinkHandler{
					ServerURL: server.URL,
				},
			},
			Logger: wrappedLogger,
		},
	)

	// Unordered output:
	// have got the link "http://example.com/1" from the page "http://example.com"
	// have got the link "http://example.com/1/1" from the page "http://example.com/1"
	// have got the link "http://example.com/1/2" from the page "http://example.com/1"
	// have got the link "http://example.com/2" from the page "http://example.com"
	// have got the link "http://example.com/2/1" from the page "http://example.com/2"
	// have got the link "http://example.com/2/2" from the page "http://example.com/2"
	// have got the link "https://golang.org/" from the page "http://example.com"
}

func ExampleCrawl_withConcurrentHandling() {
	server := RunServer()
	defer server.Close()

	logger := stdlog.New(os.Stderr, "", stdlog.LstdFlags|stdlog.Lmicroseconds)
	// wrap the standard logger via the github.com/go-log/log package
	wrappedLogger := print.New(logger)

	crawler.CrawlByConcurrentHandler(
		context.Background(),
		crawler.ConcurrencyConfig{
			ConcurrencyFactor: runtime.NumCPU(),
			BufferSize:        1000,
		},
		crawler.ConcurrencyConfig{
			ConcurrencyFactor: runtime.NumCPU(),
			BufferSize:        1000,
		},
		[]string{server.URL},
		crawler.CrawlDependencies{
			LinkExtractor: extractors.RepeatingExtractor{
				LinkExtractor: extractors.NewDelayingExtractor(
					time.Second,
					time.Sleep,
					extractors.DefaultExtractor{
						HTTPClient: http.DefaultClient,
						Filters: htmlselector.OptimizeFilters(htmlselector.FilterGroup{
							"a": {"href"},
						}),
					},
				),
				RepeatCount:  5,
				RepeatDelay:  0,
				Logger:       wrappedLogger,
				SleepHandler: time.Sleep,
			},
			LinkChecker: checkers.CheckerGroup{
				checkers.HostChecker{
					Logger: wrappedLogger,
				},
				checkers.DuplicateChecker{
					LinkRegister: registers.NewLinkRegister(
						sanitizing.SanitizeLink,
						wrappedLogger,
					),
				},
			},
			LinkHandler: handlers.CheckedHandler{
				LinkChecker: checkers.DuplicateChecker{
					// don't use here the link register from the duplicate checker above
					LinkRegister: registers.NewLinkRegister(
						sanitizing.SanitizeLink,
						wrappedLogger,
					),
				},
				LinkHandler: LinkHandler{
					ServerURL: server.URL,
				},
			},
			Logger: wrappedLogger,
		},
	)

	// Unordered output:
	// have got the link "http://example.com/1" from the page "http://example.com"
	// have got the link "http://example.com/1/1" from the page "http://example.com/1"
	// have got the link "http://example.com/1/2" from the page "http://example.com/1"
	// have got the link "http://example.com/2" from the page "http://example.com"
	// have got the link "http://example.com/2/1" from the page "http://example.com/2"
	// have got the link "http://example.com/2/2" from the page "http://example.com/2"
	// have got the link "https://golang.org/" from the page "http://example.com"
}

func ExampleHandleLinksConcurrently() {
	server := RunServer()
	defer server.Close()

	links := make(chan string, 1000)
	links <- server.URL

	var waiter sync.WaitGroup
	waiter.Add(1)

	logger := stdlog.New(os.Stderr, "", stdlog.LstdFlags|stdlog.Lmicroseconds)
	// wrap the standard logger via the github.com/go-log/log package
	wrappedLogger := print.New(logger)

	crawler.HandleLinksConcurrently(
		context.Background(),
		runtime.NumCPU(),
		links,
		crawler.HandleLinkDependencies{
			CrawlDependencies: crawler.CrawlDependencies{
				LinkExtractor: extractors.RepeatingExtractor{
					LinkExtractor: extractors.DefaultExtractor{
						HTTPClient: http.DefaultClient,
						Filters: htmlselector.OptimizeFilters(htmlselector.FilterGroup{
							"a": {"href"},
						}),
					},
					RepeatCount:  5,
					RepeatDelay:  time.Second,
					Logger:       wrappedLogger,
					SleepHandler: time.Sleep,
				},
				LinkChecker: checkers.HostChecker{
					Logger: wrappedLogger,
				},
				LinkHandler: LinkHandler{
					ServerURL: server.URL,
				},
				Logger: wrappedLogger,
			},
			Waiter: &waiter,
		},
	)

	waiter.Wait()

	// Unordered output:
	// have got the link "http://example.com/1" from the page "http://example.com"
	// have got the link "http://example.com/1/1" from the page "http://example.com/1"
	// have got the link "http://example.com/1/2" from the page "http://example.com/1"
	// have got the link "http://example.com/2" from the page "http://example.com"
	// have got the link "http://example.com/2" from the page "http://example.com"
	// have got the link "http://example.com/2/1" from the page "http://example.com/2"
	// have got the link "http://example.com/2/1" from the page "http://example.com/2"
	// have got the link "http://example.com/2/2" from the page "http://example.com/2"
	// have got the link "http://example.com/2/2" from the page "http://example.com/2"
	// have got the link "https://golang.org/" from the page "http://example.com"
}

func ExampleHandleLinksConcurrently_withoutDuplicatesOnExtracting() {
	server := RunServer()
	defer server.Close()

	links := make(chan string, 1000)
	links <- server.URL

	var waiter sync.WaitGroup
	waiter.Add(1)

	logger := stdlog.New(os.Stderr, "", stdlog.LstdFlags|stdlog.Lmicroseconds)
	// wrap the standard logger via the github.com/go-log/log package
	wrappedLogger := print.New(logger)

	crawler.HandleLinksConcurrently(
		context.Background(),
		runtime.NumCPU(),
		links,
		crawler.HandleLinkDependencies{
			CrawlDependencies: crawler.CrawlDependencies{
				LinkExtractor: extractors.RepeatingExtractor{
					LinkExtractor: extractors.DefaultExtractor{
						HTTPClient: http.DefaultClient,
						Filters: htmlselector.OptimizeFilters(htmlselector.FilterGroup{
							"a": {"href"},
						}),
					},
					RepeatCount:  5,
					RepeatDelay:  time.Second,
					Logger:       wrappedLogger,
					SleepHandler: time.Sleep,
				},
				LinkChecker: checkers.CheckerGroup{
					checkers.HostChecker{
						Logger: wrappedLogger,
					},
					checkers.DuplicateChecker{
						LinkRegister: registers.NewLinkRegister(
							sanitizing.SanitizeLink,
							wrappedLogger,
						),
					},
				},
				LinkHandler: LinkHandler{
					ServerURL: server.URL,
				},
				Logger: wrappedLogger,
			},
			Waiter: &waiter,
		},
	)

	waiter.Wait()

	// Unordered output:
	// have got the link "http://example.com/1" from the page "http://example.com"
	// have got the link "http://example.com/1/1" from the page "http://example.com/1"
	// have got the link "http://example.com/1/2" from the page "http://example.com/1"
	// have got the link "http://example.com/2" from the page "http://example.com"
	// have got the link "http://example.com/2" from the page "http://example.com"
	// have got the link "http://example.com/2/1" from the page "http://example.com/2"
	// have got the link "http://example.com/2/2" from the page "http://example.com/2"
	// have got the link "https://golang.org/" from the page "http://example.com"
}

func ExampleHandleLinksConcurrently_withoutDuplicatesOnHandling() {
	server := RunServer()
	defer server.Close()

	links := make(chan string, 1000)
	links <- server.URL

	var waiter sync.WaitGroup
	waiter.Add(1)

	logger := stdlog.New(os.Stderr, "", stdlog.LstdFlags|stdlog.Lmicroseconds)
	// wrap the standard logger via the github.com/go-log/log package
	wrappedLogger := print.New(logger)

	crawler.HandleLinksConcurrently(
		context.Background(),
		runtime.NumCPU(),
		links,
		crawler.HandleLinkDependencies{
			CrawlDependencies: crawler.CrawlDependencies{
				LinkExtractor: extractors.RepeatingExtractor{
					LinkExtractor: extractors.DefaultExtractor{
						HTTPClient: http.DefaultClient,
						Filters: htmlselector.OptimizeFilters(htmlselector.FilterGroup{
							"a": {"href"},
						}),
					},
					RepeatCount:  5,
					RepeatDelay:  time.Second,
					Logger:       wrappedLogger,
					SleepHandler: time.Sleep,
				},
				LinkChecker: checkers.CheckerGroup{
					checkers.HostChecker{
						Logger: wrappedLogger,
					},
					checkers.DuplicateChecker{
						LinkRegister: registers.NewLinkRegister(
							sanitizing.SanitizeLink,
							wrappedLogger,
						),
					},
				},
				LinkHandler: handlers.CheckedHandler{
					LinkChecker: checkers.DuplicateChecker{
						// don't use here the link register from the duplicate checker above
						LinkRegister: registers.NewLinkRegister(
							sanitizing.SanitizeLink,
							wrappedLogger,
						),
					},
					LinkHandler: LinkHandler{
						ServerURL: server.URL,
					},
				},
				Logger: wrappedLogger,
			},
			Waiter: &waiter,
		},
	)

	waiter.Wait()

	// Unordered output:
	// have got the link "http://example.com/1" from the page "http://example.com"
	// have got the link "http://example.com/1/1" from the page "http://example.com/1"
	// have got the link "http://example.com/1/2" from the page "http://example.com/1"
	// have got the link "http://example.com/2" from the page "http://example.com"
	// have got the link "http://example.com/2/1" from the page "http://example.com/2"
	// have got the link "http://example.com/2/2" from the page "http://example.com/2"
	// have got the link "https://golang.org/" from the page "http://example.com"
}

func ExampleHandleLinksConcurrently_withDelayingExtracting() {
	server := RunServer()
	defer server.Close()

	links := make(chan string, 1000)
	links <- server.URL

	var waiter sync.WaitGroup
	waiter.Add(1)

	logger := stdlog.New(os.Stderr, "", stdlog.LstdFlags|stdlog.Lmicroseconds)
	// wrap the standard logger via the github.com/go-log/log package
	wrappedLogger := print.New(logger)

	crawler.HandleLinksConcurrently(
		context.Background(),
		runtime.NumCPU(),
		links,
		crawler.HandleLinkDependencies{
			CrawlDependencies: crawler.CrawlDependencies{
				LinkExtractor: extractors.RepeatingExtractor{
					LinkExtractor: extractors.NewDelayingExtractor(
						time.Second,
						time.Sleep,
						extractors.DefaultExtractor{
							HTTPClient: http.DefaultClient,
							Filters: htmlselector.OptimizeFilters(htmlselector.FilterGroup{
								"a": {"href"},
							}),
						},
					),
					RepeatCount:  5,
					RepeatDelay:  0,
					Logger:       wrappedLogger,
					SleepHandler: time.Sleep,
				},
				LinkChecker: checkers.CheckerGroup{
					checkers.HostChecker{
						Logger: wrappedLogger,
					},
					checkers.DuplicateChecker{
						LinkRegister: registers.NewLinkRegister(
							sanitizing.SanitizeLink,
							wrappedLogger,
						),
					},
				},
				LinkHandler: handlers.CheckedHandler{
					LinkChecker: checkers.DuplicateChecker{
						// don't use here the link register from the duplicate checker above
						LinkRegister: registers.NewLinkRegister(
							sanitizing.SanitizeLink,
							wrappedLogger,
						),
					},
					LinkHandler: LinkHandler{
						ServerURL: server.URL,
					},
				},
				Logger: wrappedLogger,
			},
			Waiter: &waiter,
		},
	)

	waiter.Wait()

	// Unordered output:
	// have got the link "http://example.com/1" from the page "http://example.com"
	// have got the link "http://example.com/1/1" from the page "http://example.com/1"
	// have got the link "http://example.com/1/2" from the page "http://example.com/1"
	// have got the link "http://example.com/2" from the page "http://example.com"
	// have got the link "http://example.com/2/1" from the page "http://example.com/2"
	// have got the link "http://example.com/2/2" from the page "http://example.com/2"
	// have got the link "https://golang.org/" from the page "http://example.com"
}

func ExampleHandleLinksConcurrently_withRobotsTXTOnExtracting() {
	server := RunServer()
	defer server.Close()

	links := make(chan string, 1000)
	links <- server.URL

	var waiter sync.WaitGroup
	waiter.Add(1)

	logger := stdlog.New(os.Stderr, "", stdlog.LstdFlags|stdlog.Lmicroseconds)
	// wrap the standard logger via the github.com/go-log/log package
	wrappedLogger := print.New(logger)

	crawler.HandleLinksConcurrently(
		context.Background(),
		runtime.NumCPU(),
		links,
		crawler.HandleLinkDependencies{
			CrawlDependencies: crawler.CrawlDependencies{
				LinkExtractor: extractors.RepeatingExtractor{
					LinkExtractor: extractors.DefaultExtractor{
						HTTPClient: http.DefaultClient,
						Filters: htmlselector.OptimizeFilters(htmlselector.FilterGroup{
							"a": {"href"},
						}),
					},
					RepeatCount:  5,
					RepeatDelay:  time.Second,
					Logger:       wrappedLogger,
					SleepHandler: time.Sleep,
				},
				LinkChecker: checkers.CheckerGroup{
					checkers.HostChecker{
						Logger: wrappedLogger,
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
			Waiter: &waiter,
		},
	)

	waiter.Wait()

	// Unordered output:
	// have got the link "http://example.com/1" from the page "http://example.com"
	// have got the link "http://example.com/1/1" from the page "http://example.com/1"
	// have got the link "http://example.com/1/2" from the page "http://example.com/1"
	// have got the link "http://example.com/2" from the page "http://example.com"
	// have got the link "http://example.com/2" from the page "http://example.com"
	// have got the link "https://golang.org/" from the page "http://example.com"
}

func ExampleHandleLinksConcurrently_withRobotsTXTOnHandling() {
	server := RunServer()
	defer server.Close()

	links := make(chan string, 1000)
	links <- server.URL

	var waiter sync.WaitGroup
	waiter.Add(1)

	logger := stdlog.New(os.Stderr, "", stdlog.LstdFlags|stdlog.Lmicroseconds)
	// wrap the standard logger via the github.com/go-log/log package
	wrappedLogger := print.New(logger)

	crawler.HandleLinksConcurrently(
		context.Background(),
		runtime.NumCPU(),
		links,
		crawler.HandleLinkDependencies{
			CrawlDependencies: crawler.CrawlDependencies{
				LinkExtractor: extractors.RepeatingExtractor{
					LinkExtractor: extractors.DefaultExtractor{
						HTTPClient: http.DefaultClient,
						Filters: htmlselector.OptimizeFilters(htmlselector.FilterGroup{
							"a": {"href"},
						}),
					},
					RepeatCount:  5,
					RepeatDelay:  time.Second,
					Logger:       wrappedLogger,
					SleepHandler: time.Sleep,
				},
				LinkChecker: checkers.HostChecker{
					Logger: wrappedLogger,
				},
				LinkHandler: handlers.CheckedHandler{
					LinkChecker: checkers.RobotsTXTChecker{
						UserAgent:         "go-crawler",
						RobotsTXTRegister: registers.NewRobotsTXTRegister(http.DefaultClient),
						Logger:            wrappedLogger,
					},
					LinkHandler: LinkHandler{
						ServerURL: server.URL,
					},
				},
				Logger: wrappedLogger,
			},
			Waiter: &waiter,
		},
	)

	waiter.Wait()

	// Unordered output:
	// have got the link "http://example.com/1" from the page "http://example.com"
	// have got the link "http://example.com/1/1" from the page "http://example.com/1"
	// have got the link "http://example.com/1/2" from the page "http://example.com/1"
	// have got the link "https://golang.org/" from the page "http://example.com"
}

func ExampleHandleLinksConcurrently_withSitemap() {
	server := RunServer()
	defer server.Close()

	links := make(chan string, 1000)
	links <- server.URL

	var waiter sync.WaitGroup
	waiter.Add(1)

	logger := stdlog.New(os.Stderr, "", stdlog.LstdFlags|stdlog.Lmicroseconds)
	// wrap the standard logger via the github.com/go-log/log package
	wrappedLogger := print.New(logger)

	crawler.HandleLinksConcurrently(
		context.Background(),
		runtime.NumCPU(),
		links,
		crawler.HandleLinkDependencies{
			CrawlDependencies: crawler.CrawlDependencies{
				LinkExtractor: extractors.RepeatingExtractor{
					LinkExtractor: extractors.ExtractorGroup{
						extractors.DefaultExtractor{
							HTTPClient: http.DefaultClient,
							Filters: htmlselector.OptimizeFilters(htmlselector.FilterGroup{
								"a": {"href"},
							}),
						},
						extractors.SitemapExtractor{
							SitemapRegister: registers.NewSitemapRegister(
								time.Second,
								sitemap.SimpleLinkGenerator{},
								wrappedLogger,
								time.Sleep,
								nil,
							),
							Logger: wrappedLogger,
						},
					},
					RepeatCount:  5,
					RepeatDelay:  time.Second,
					Logger:       wrappedLogger,
					SleepHandler: time.Sleep,
				},
				LinkChecker: checkers.CheckerGroup{
					checkers.HostChecker{
						Logger: wrappedLogger,
					},
					checkers.DuplicateChecker{
						LinkRegister: registers.NewLinkRegister(
							sanitizing.SanitizeLink,
							wrappedLogger,
						),
					},
				},
				LinkHandler: handlers.CheckedHandler{
					LinkChecker: checkers.DuplicateChecker{
						// don't use here the link register from the duplicate checker above
						LinkRegister: registers.NewLinkRegister(
							sanitizing.SanitizeLink,
							wrappedLogger,
						),
					},
					LinkHandler: LinkHandler{
						ServerURL: server.URL,
					},
				},
				Logger: wrappedLogger,
			},
			Waiter: &waiter,
		},
	)

	waiter.Wait()

	// Unordered output:
	// have got the link "http://example.com/1" from the page "http://example.com"
	// have got the link "http://example.com/1/1" from the page "http://example.com/1"
	// have got the link "http://example.com/1/2" from the page "http://example.com/1"
	// have got the link "http://example.com/2" from the page "http://example.com"
	// have got the link "http://example.com/2/1" from the page "http://example.com/2"
	// have got the link "http://example.com/2/2" from the page "http://example.com/2"
	// have got the link "http://example.com/hidden/1" from the page "http://example.com"
	// have got the link "http://example.com/hidden/2" from the page "http://example.com"
	// have got the link "https://golang.org/" from the page "http://example.com"
}
