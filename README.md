# go-crawler

[![GoDoc](https://godoc.org/github.com/thewizardplusplus/go-crawler?status.svg)](https://godoc.org/github.com/thewizardplusplus/go-crawler)
[![Go Report Card](https://goreportcard.com/badge/github.com/thewizardplusplus/go-crawler)](https://goreportcard.com/report/github.com/thewizardplusplus/go-crawler)
[![Build Status](https://travis-ci.org/thewizardplusplus/go-crawler.svg?branch=master)](https://travis-ci.org/thewizardplusplus/go-crawler)
[![codecov](https://codecov.io/gh/thewizardplusplus/go-crawler/branch/master/graph/badge.svg)](https://codecov.io/gh/thewizardplusplus/go-crawler)

The library that implements crawling of all relative links for specified ones.

## Features

- crawling of all relative links for specified ones:
  - resolving of relative links:
    - by the `base` tag;
    - by the `Content-Base` and `Content-Location` headers;
    - by the request URI;
  - supporting of leading and trailing spaces trimming in extracted links (optional);
  - repeated extracting of relative links on error (optional):
    - only specified repeat count;
    - supporting of delay between repeats;
  - delayed extracting of relative links (optional):
    - reducing of a delay time by the time elapsed since the last request;
    - using of individual delays for each thread;
  - extracting links from a `sitemap.xml` file (optional):
    - ignoring of the error on loading of the `sitemap.xml` file:
      - logging of the received error;
      - returning of an empty Sitemap instead;
    - supporting of few `sitemap.xml` files for a single link:
      - processing of each `sitemap.xml` file is done in a separate goroutine;
      - supporting of an outer generator for `sitemap.xml` links:
        - generators:
          - simple generator (it returns the `sitemap.xml` file in the site root);
          - hierarchical generator (it returns the suitable `sitemap.xml` file for each part of the URL path);
          - generator based on the `robots.txt` file;
        - supporting of grouping of generators:
          - result of group generating is merged results of each generator in the group;
          - generating concurrently:
            - processing of each generator is done in a separate goroutine;
    - supporting of a Sitemap index file:
      - supporting of a delay before loading of each `sitemap.xml` file listed in the index;
    - supporting of a gzip compression of a `sitemap.xml` file;
  - supporting of grouping of link extractors:
    - result of group extracting is merged results of each extractor in the group;
    - extracting links concurrently:
      - processing of each link extractor is done in a separate goroutine;
- calling of an outer handler for an each found link:
  - it's called directly during crawling;
  - handling of links immediately after they have been extracted;
  - passing of the source link in the outer handler;
  - handling links filtered by a custom link filter (optional);
  - handling links concurrently (optional);
  - supporting of grouping of outer handlers:
    - processing of each outer handler is done in a separate goroutine;
- custom filtering of considered links:
  - by relativity of a link (optional):
    - supporting of result inverting;
  - by uniqueness of an extracted link (optional):
    - supporting of sanitizing of a link before checking of uniqueness (optional);
  - by a `robots.txt` file (optional):
    - customized user agent;
  - supporting of grouping of link filters:
    - result of group filtering is successful only when all filters are successful;
- parallelization possibilities:
  - crawling of relative links in parallel;
  - supporting of background working:
    - automatic completion after processing all filtered links;
  - simulate an unbounded channel of links to avoid a deadlock.

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
	"time"

	"github.com/go-log/log/print"
	crawler "github.com/thewizardplusplus/go-crawler"
	"github.com/thewizardplusplus/go-crawler/checkers"
	"github.com/thewizardplusplus/go-crawler/extractors"
	"github.com/thewizardplusplus/go-crawler/extractors/transformers"
	"github.com/thewizardplusplus/go-crawler/models"
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
	"time"

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
	"time"

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
	"time"

	"github.com/go-log/log/print"
	crawler "github.com/thewizardplusplus/go-crawler"
	"github.com/thewizardplusplus/go-crawler/checkers"
	"github.com/thewizardplusplus/go-crawler/extractors"
	"github.com/thewizardplusplus/go-crawler/extractors/transformers"
	"github.com/thewizardplusplus/go-crawler/models"
	"github.com/thewizardplusplus/go-crawler/registers"
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
			fmt.Fprint(writer, `
				User-agent: go-crawler
				Disallow: /2
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
	"time"

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

## License

The MIT License (MIT)

Copyright &copy; 2020-2021 thewizardplusplus
