# go-crawler

[![GoDoc](https://godoc.org/github.com/thewizardplusplus/go-crawler?status.svg)](https://godoc.org/github.com/thewizardplusplus/go-crawler)
[![Go Report Card](https://goreportcard.com/badge/github.com/thewizardplusplus/go-crawler)](https://goreportcard.com/report/github.com/thewizardplusplus/go-crawler)
[![Build Status](https://travis-ci.org/thewizardplusplus/go-crawler.svg?branch=master)](https://travis-ci.org/thewizardplusplus/go-crawler)
[![codecov](https://codecov.io/gh/thewizardplusplus/go-crawler/branch/master/graph/badge.svg)](https://codecov.io/gh/thewizardplusplus/go-crawler)

The library that implements crawling of all relative links for specified ones.

## Features

- crawling of all relative links for specified ones:
  - names of tags and attributes of links may be configured;
  - supporting of an outer transformer for the extracted links (optional):
    - data passed to the transformer:
      - extracted links;
      - service data of the HTTP response;
      - content of the HTTP response as bytes;
    - transformers:
      - leading and trailing spaces trimming in the extracted links;
      - resolving of relative links:
        - by the base tag:
          - tag and attribute names may be configured (`<base href="..." />` by default);
          - tag selection:
            - first occurrence;
            - last occurrence;
        - by the header list:
          - the headers are listed in the descending order of the priority;
          - `Content-Base` and `Content-Location` by default;
        - by the request URI;
    - supporting of grouping of transformers:
      - the transformers are processed sequentially, so one transformer can influence another one;
  - supporting of leading and trailing spaces trimming in extracted links (optional):
    - as the transformer for the extracted links (see above);
    - as the wrapper for a link extractor;
  - repeated extracting of relative links on error (optional):
    - only the specified repeat count;
    - supporting of a delay between repeats;
  - delayed extracting of relative links (optional):
    - reducing of a delay time by the time elapsed since the last request;
    - using of individual delays for each thread;
  - extracting links from a `sitemap.xml` file (optional):
    - in-memory caching of the loaded `sitemap.xml` files;
    - ignoring of the error on loading of the `sitemap.xml` file:
      - logging of the received error;
      - returning of the empty Sitemap instead;
    - supporting of few `sitemap.xml` files for a single link:
      - processing of each `sitemap.xml` file is done in a separate goroutine;
      - supporting of an outer generator for the `sitemap.xml` links:
        - generators:
          - hierarchical generator:
            - returns the suitable `sitemap.xml` file for each part of the URL path;
            - supporting of sanitizing of the base link before generating of the `sitemap.xml` links;
            - supporting of the restriction of the maximal depth;
          - generator based on the `robots.txt` file;
        - supporting of grouping of generators:
          - result of group generating is merged results of each generator in the group;
          - processing of each generator is done in a separate goroutine;
    - supporting of a Sitemap index file:
      - supporting of a delay before loading of each `sitemap.xml` file listed in the index;
    - supporting of a gzip compression of a `sitemap.xml` file;
  - supporting of grouping of link extractors:
    - result of group extracting is merged results of each link extractor in the group;
    - processing of each link extractor is done in a separate goroutine;
- calling of an outer handler for each extracted link:
  - handling of the extracted links directly during the crawling, i.e., immediately after they have been extracted;
  - data passed to the handler:
    - extracted link;
    - source link for the extracted link;
  - handling only of those extracted links that have been filtered by a link filter (see below; optional);
  - handling of the extracted links concurrently, i.e., in the goroutine pool (optional);
  - supporting of grouping of handlers:
    - processing of each handler is done in a separate goroutine;
- filtering of the extracted links by an outer link filter:
  - by relativity of the extracted link (optional):
    - supporting of result inverting;
  - by uniqueness of the extracted link (optional):
    - supporting of sanitizing of the link before checking of uniqueness;
  - by a `robots.txt` file (optional):
    - customized user agent;
    - in-memory caching of the loaded `robots.txt` files;
  - supporting of grouping of link filters:
    - the link filters are processed sequentially, so one link filter can influence another one;
    - result of group filtering is successful only when all link filters are successful;
    - the empty group of link filters is always failed;
- parallelization possibilities:
  - crawling of relative links concurrently, i.e., in the goroutine pool;
  - simulation of an unbounded channel of links to avoid a deadlock;
  - waiting of completion of processing of all extracted links;
  - supporting of stopping of all operations via the context.

## Installation

Prepare the directory:

```
$ mkdir --parents "$(go env GOPATH)/src/github.com/thewizardplusplus/"
$ cd "$(go env GOPATH)/src/github.com/thewizardplusplus/"
```

Clone this repository:

```
$ git clone https://github.com/thewizardplusplus/go-crawler.git
$ cd go-crawler
```

Install dependencies with the [dep](https://golang.github.io/dep/) tool:

```
$ dep ensure -vendor-only
```

## Examples

`crawler.Crawl()` with all the features:

```go
package main

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
			fmt.Fprintf( // nolint: errcheck
				writer,
				`
					User-agent: go-crawler
					Disallow: /2

					Sitemap: %s
				`,
				sitemapLink,
			)

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

func main() {
	server := RunServer()
	defer server.Close()

	logger := stdlog.New(os.Stderr, "", stdlog.LstdFlags|stdlog.Lmicroseconds)
	// wrap the standard logger via the github.com/go-log/log package
	wrappedLogger := print.New(logger)

	robotsTXTRegister := registers.NewRobotsTXTRegister(http.DefaultClient)
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
			LinkExtractor: extractors.NewDelayingExtractor(
				time.Second,
				time.Sleep,
				extractors.ExtractorGroup{
					Name: "main extractors",
					LinkExtractors: []models.LinkExtractor{
						extractors.RepeatingExtractor{
							LinkExtractor: extractors.DefaultExtractor{
								HTTPClient: http.DefaultClient,
								Filters: htmlselector.OptimizeFilters(htmlselector.FilterGroup{
									"a": {"href"},
								}),
								LinkTransformer: transformers.TransformerGroup{
									transformers.TrimmingTransformer{
										TrimLink: urlutils.TrimLink,
									},
									transformers.ResolvingTransformer{
										BaseTagSelection: transformers.SelectFirstBaseTag,
										BaseTagFilters:   transformers.DefaultBaseTagFilters,
										BaseHeaderNames:  urlutils.DefaultBaseHeaderNames,
										Logger:           wrappedLogger,
									},
								},
							},
							RepeatCount:  5,
							RepeatDelay:  time.Second,
							Logger:       wrappedLogger,
							SleepHandler: time.Sleep,
						},
						extractors.RepeatingExtractor{
							LinkExtractor: extractors.TrimmingExtractor{
								TrimLink: urlutils.TrimLink,
								LinkExtractor: extractors.SitemapExtractor{
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
													RobotsTXTRegister: robotsTXTRegister,
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
							RepeatCount:  5,
							RepeatDelay:  time.Second,
							Logger:       wrappedLogger,
							SleepHandler: time.Sleep,
						},
					},
					Logger: wrappedLogger,
				},
			),
			LinkChecker: checkers.CheckerGroup{
				checkers.HostChecker{
					Logger: wrappedLogger,
				},
				checkers.DuplicateChecker{
					LinkRegister: registers.NewLinkRegister(urlutils.SanitizeLink),
					Logger:       wrappedLogger,
				},
				checkers.RobotsTXTChecker{
					UserAgent:         "go-crawler",
					RobotsTXTRegister: robotsTXTRegister,
					Logger:            wrappedLogger,
				},
			},
			LinkHandler: handlers.CheckedHandler{
				LinkChecker: checkers.DuplicateChecker{
					// don't use here the link register from the duplicate checker above
					LinkRegister: registers.NewLinkRegister(urlutils.SanitizeLink),
					Logger:       wrappedLogger,
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
			},
			Logger: wrappedLogger,
		},
	)

	// Unordered output:
	// [inner] received link "http://example.com/1" from page "http://example.com"
	// [inner] received link "http://example.com/1/1" from page "http://example.com/1"
	// [inner] received link "http://example.com/1/2" from page "http://example.com/1"
	// [inner] received link "http://example.com/2" from page "http://example.com"
	// [inner] received link "http://example.com/hidden/1" from page "http://example.com"
	// [inner] received link "http://example.com/hidden/1/test" from page "http://example.com/hidden/1"
	// [inner] received link "http://example.com/hidden/2" from page "http://example.com"
	// [inner] received link "http://example.com/hidden/3" from page "http://example.com"
	// [inner] received link "http://example.com/hidden/4" from page "http://example.com"
	// [inner] received link "http://example.com/hidden/5" from page "http://example.com/hidden/1/test"
	// [inner] received link "http://example.com/hidden/6" from page "http://example.com/hidden/1/test"
	// [outer] received link "https://golang.org/" from page "http://example.com"
}
```

---

`crawler.Crawl()`:

```go
package main

import (
	"context"
	"fmt"
	"html/template"
	stdlog "log"
	"net/http"
	"net/http/httptest"
	"os"
	"runtime"
	"strings"

	"github.com/go-log/log/print"
	crawler "github.com/thewizardplusplus/go-crawler"
	"github.com/thewizardplusplus/go-crawler/checkers"
	"github.com/thewizardplusplus/go-crawler/extractors"
	"github.com/thewizardplusplus/go-crawler/extractors/transformers"
	"github.com/thewizardplusplus/go-crawler/models"
	urlutils "github.com/thewizardplusplus/go-crawler/url-utils"
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
		"received link %q from page %q\n",
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
		var links []string
		switch request.URL.Path {
		case "/":
			links = []string{"/1", "/2", "/2", "https://golang.org/"}
		case "/1":
			links = []string{"/1/1", "/1/2"}
		case "/2":
			links = []string{"/2/1", "/2/2"}
		}

		template, _ := template.New("").Parse( // nolint: errcheck
			`<ul>
				{{ range $link := . }}
					<li><a href="{{ $link }}">{{ $link }}</a></li>
				{{ end }}
			</ul>`,
		)
		template.Execute(writer, links) // nolint: errcheck
	}))
}

func main() {
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
```

`crawler.Crawl()` without duplicates on extracting:

```go
package main

import (
	"context"
	"fmt"
	"html/template"
	stdlog "log"
	"net/http"
	"net/http/httptest"
	"os"
	"runtime"
	"strings"

	"github.com/go-log/log/print"
	crawler "github.com/thewizardplusplus/go-crawler"
	"github.com/thewizardplusplus/go-crawler/checkers"
	"github.com/thewizardplusplus/go-crawler/extractors"
	"github.com/thewizardplusplus/go-crawler/extractors/transformers"
	"github.com/thewizardplusplus/go-crawler/models"
	"github.com/thewizardplusplus/go-crawler/registers"
	urlutils "github.com/thewizardplusplus/go-crawler/url-utils"
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
		"received link %q from page %q\n",
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
		var links []string
		switch request.URL.Path {
		case "/":
			links = []string{"/1", "/2", "/2", "https://golang.org/"}
		case "/1":
			links = []string{"/1/1", "/1/2"}
		case "/2":
			links = []string{"/2/1", "/2/2"}
		}

		template, _ := template.New("").Parse( // nolint: errcheck
			`<ul>
				{{ range $link := . }}
					<li><a href="{{ $link }}">{{ $link }}</a></li>
				{{ end }}
			</ul>`,
		)
		template.Execute(writer, links) // nolint: errcheck
	}))
}

func main() {
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
```

`crawler.Crawl()` without duplicates on handling:

```go
package main

import (
	"context"
	"fmt"
	"html/template"
	stdlog "log"
	"net/http"
	"net/http/httptest"
	"os"
	"runtime"
	"strings"

	"github.com/go-log/log/print"
	crawler "github.com/thewizardplusplus/go-crawler"
	"github.com/thewizardplusplus/go-crawler/checkers"
	"github.com/thewizardplusplus/go-crawler/extractors"
	"github.com/thewizardplusplus/go-crawler/extractors/transformers"
	"github.com/thewizardplusplus/go-crawler/handlers"
	"github.com/thewizardplusplus/go-crawler/models"
	"github.com/thewizardplusplus/go-crawler/registers"
	urlutils "github.com/thewizardplusplus/go-crawler/url-utils"
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
		"received link %q from page %q\n",
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
		var links []string
		switch request.URL.Path {
		case "/":
			links = []string{"/1", "/2", "/2", "https://golang.org/"}
		case "/1":
			links = []string{"/1/1", "/1/2"}
		case "/2":
			links = []string{"/2/1", "/2/2"}
		}

		template, _ := template.New("").Parse( // nolint: errcheck
			`<ul>
				{{ range $link := . }}
					<li><a href="{{ $link }}">{{ $link }}</a></li>
				{{ end }}
			</ul>`,
		)
		template.Execute(writer, links) // nolint: errcheck
	}))
}

func main() {
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
```

`crawler.Crawl()` with processing a `robots.txt` file:

```go
package main

import (
	"context"
	"fmt"
	"html/template"
	stdlog "log"
	"net/http"
	"net/http/httptest"
	"os"
	"runtime"
	"strings"

	"github.com/go-log/log/print"
	crawler "github.com/thewizardplusplus/go-crawler"
	"github.com/thewizardplusplus/go-crawler/checkers"
	"github.com/thewizardplusplus/go-crawler/extractors"
	"github.com/thewizardplusplus/go-crawler/extractors/transformers"
	"github.com/thewizardplusplus/go-crawler/models"
	"github.com/thewizardplusplus/go-crawler/registers"
	urlutils "github.com/thewizardplusplus/go-crawler/url-utils"
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
		"received link %q from page %q\n",
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
			fmt.Fprint( // nolint: errcheck
				writer,
				`
					User-agent: go-crawler
					Disallow: /2
				`,
			)

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

		template, _ := template.New("").Parse( // nolint: errcheck
			`<ul>
				{{ range $link := . }}
					<li><a href="{{ $link }}">{{ $link }}</a></li>
				{{ end }}
			</ul>`,
		)
		template.Execute(writer, links) // nolint: errcheck
	}))
}

func main() {
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
```

`crawler.Crawl()` with processing a `sitemap.xml` file:

```go
package main

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
	ServerURL string
}

func (handler LinkHandler) HandleLink(
	ctx context.Context,
	link models.SourcedLink,
) {
	fmt.Printf(
		"received link %q from page %q\n",
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
			fmt.Fprintf( // nolint: errcheck
				writer,
				`
					User-agent: go-crawler
					Disallow: /2

					Sitemap: %s
				`,
				sitemapLink,
			)

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

func main() {
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
```

`crawler.Crawl()` with few handlers:

```go
package main

import (
	"context"
	"fmt"
	"html/template"
	stdlog "log"
	"net/http"
	"net/http/httptest"
	"os"
	"runtime"
	"strings"

	"github.com/go-log/log/print"
	crawler "github.com/thewizardplusplus/go-crawler"
	"github.com/thewizardplusplus/go-crawler/checkers"
	"github.com/thewizardplusplus/go-crawler/extractors"
	"github.com/thewizardplusplus/go-crawler/extractors/transformers"
	"github.com/thewizardplusplus/go-crawler/handlers"
	"github.com/thewizardplusplus/go-crawler/models"
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
	fmt.Printf(
		"[%s] received link %q from page %q\n",
		handler.Name,
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
		var links []string
		switch request.URL.Path {
		case "/":
			links = []string{"/1", "/2", "/2", "https://golang.org/"}
		case "/1":
			links = []string{"/1/1", "/1/2"}
		case "/2":
			links = []string{"/2/1", "/2/2"}
		}

		template, _ := template.New("").Parse( // nolint: errcheck
			`<ul>
				{{ range $link := . }}
					<li><a href="{{ $link }}">{{ $link }}</a></li>
				{{ end }}
			</ul>`,
		)
		template.Execute(writer, links) // nolint: errcheck
	}))
}

func main() {
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
```

## Bibliography

- [Official Website of the Robots Exclusion Standard (`robots.txt`)](http://www.robotstxt.org/)
- [Official Website of the Sitemaps Protocol (`sitemap.xml`)](https://www.sitemaps.org/):
  - [Using Sitemap Index Files](https://www.sitemaps.org/protocol.html#index)
  - [Sitemap File Location](https://www.sitemaps.org/protocol.html#location)
- [Is There an HTTP Header to Say What Base URL to Use for Relative Links?](https://stackoverflow.com/a/48409040)
  - [HTML 4.0 Specification, 5.1.3 URLs in HTML](https://www.w3.org/TR/WD-html40-970917/htmlweb.html#h-5.1.3):
    1. [`<base>`: The Document Base URL element](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/base)
    2. HTTP headers:
       1. RFC 2068 Hypertext Transfer Protocol &mdash; HTTP/1.1:
          1. [14.11 `Content-Base`](https://datatracker.ietf.org/doc/html/rfc2068#section-14.11)
          2. [14.15 `Content-Location`](https://datatracker.ietf.org/doc/html/rfc2068#section-14.15)
       2. RFC 2616 Hypertext Transfer Protocol &mdash; HTTP/1.1:
          - [14.14 `Content-Location`](https://datatracker.ietf.org/doc/html/rfc2616#section-14.14)
       3. RFC 7231 Hypertext Transfer Protocol (HTTP/1.1): Semantics and Content:
          - [Appendix B Changes from RFC 2616](https://datatracker.ietf.org/doc/html/rfc7231#appendix-B)
  - [HTML Living Standard, 2.4.1 URLs: Terminology](https://html.spec.whatwg.org/multipage/urls-and-fetching.html#terminology-2)

## License

The MIT License (MIT)

Copyright &copy; 2020-2021 thewizardplusplus
