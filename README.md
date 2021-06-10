# go-crawler

[![GoDoc](https://godoc.org/github.com/thewizardplusplus/go-crawler?status.svg)](https://godoc.org/github.com/thewizardplusplus/go-crawler)
[![Go Report Card](https://goreportcard.com/badge/github.com/thewizardplusplus/go-crawler)](https://goreportcard.com/report/github.com/thewizardplusplus/go-crawler)
[![Build Status](https://travis-ci.org/thewizardplusplus/go-crawler.svg?branch=master)](https://travis-ci.org/thewizardplusplus/go-crawler)
[![codecov](https://codecov.io/gh/thewizardplusplus/go-crawler/branch/master/graph/badge.svg)](https://codecov.io/gh/thewizardplusplus/go-crawler)

The library that implements crawling of all relative links for specified ones.

## Features

- crawling of all relative links for specified ones:
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
    - supporting of a Sitemap index file:
      - supporting of a delay before loading of each `sitemap.xml` file listed in the index;
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
- custom filtering of considered links:
  - by relativity of a link (optional);
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
	"sync"
	"time"

	"github.com/go-log/log/print"
	crawler "github.com/thewizardplusplus/go-crawler"
	"github.com/thewizardplusplus/go-crawler/checkers"
	"github.com/thewizardplusplus/go-crawler/extractors"
	"github.com/thewizardplusplus/go-crawler/handlers"
	"github.com/thewizardplusplus/go-crawler/register"
	"github.com/thewizardplusplus/go-crawler/sanitizing"
	htmlselector "github.com/thewizardplusplus/go-html-selector"
	"github.com/thewizardplusplus/go-crawler/models"
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
		var links []string
		switch request.URL.Path {
		case "/":
			links = []string{"/1", "/2", "/2", "https://golang.org/"}
		case "/1":
			links = []string{"/1/1", "/1/2"}
		case "/2":
			links = []string{"/2/1", "/2/2"}
		}
		for index := range links {
			if strings.HasPrefix(links[index], "/") {
				links[index] = "http://" + request.Host + links[index]
			}
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
```

`crawler.CrawlByConcurrentHandler()`:

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
	"github.com/thewizardplusplus/go-crawler/handlers"
	"github.com/thewizardplusplus/go-crawler/registers"
	"github.com/thewizardplusplus/go-crawler/sanitizing"
	htmlselector "github.com/thewizardplusplus/go-html-selector"
	"github.com/thewizardplusplus/go-crawler/models"
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

		var links []string
		switch request.URL.Path {
		case "/":
			links = []string{"/1", "/2", "/2", "https://golang.org/"}
		case "/1":
			links = []string{"/1/1", "/1/2"}
		case "/2":
			links = []string{"/2/1", "/2/2"}
		}
		for index := range links {
			if strings.HasPrefix(links[index], "/") {
				links[index] = "http://" + request.Host + links[index]
			}
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
```

`crawler.HandleLinksConcurrently()`:

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
	"sync"
	"time"

	"github.com/go-log/log/print"
	crawler "github.com/thewizardplusplus/go-crawler"
	"github.com/thewizardplusplus/go-crawler/checkers"
	"github.com/thewizardplusplus/go-crawler/extractors"
	htmlselector "github.com/thewizardplusplus/go-html-selector"
	"github.com/thewizardplusplus/go-crawler/models"
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
		var links []string
		switch request.URL.Path {
		case "/":
			links = []string{"/1", "/2", "/2", "https://golang.org/"}
		case "/1":
			links = []string{"/1/1", "/1/2"}
		case "/2":
			links = []string{"/2/1", "/2/2"}
		}
		for index := range links {
			if strings.HasPrefix(links[index], "/") {
				links[index] = "http://" + request.Host + links[index]
			}
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
```

`crawler.HandleLinksConcurrently()` without duplicates on extracting:

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
	"sync"
	"time"

	"github.com/go-log/log/print"
	crawler "github.com/thewizardplusplus/go-crawler"
	"github.com/thewizardplusplus/go-crawler/checkers"
	"github.com/thewizardplusplus/go-crawler/extractors"
	"github.com/thewizardplusplus/go-crawler/register"
	"github.com/thewizardplusplus/go-crawler/sanitizing"
	htmlselector "github.com/thewizardplusplus/go-html-selector"
	"github.com/thewizardplusplus/go-crawler/models"
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
		var links []string
		switch request.URL.Path {
		case "/":
			links = []string{"/1", "/2", "/2", "https://golang.org/"}
		case "/1":
			links = []string{"/1/1", "/1/2"}
		case "/2":
			links = []string{"/2/1", "/2/2"}
		}
		for index := range links {
			if strings.HasPrefix(links[index], "/") {
				links[index] = "http://" + request.Host + links[index]
			}
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
```

`crawler.HandleLinksConcurrently()` without duplicates on handling:

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
	"sync"
	"time"

	"github.com/go-log/log/print"
	crawler "github.com/thewizardplusplus/go-crawler"
	"github.com/thewizardplusplus/go-crawler/checkers"
	"github.com/thewizardplusplus/go-crawler/extractors"
	"github.com/thewizardplusplus/go-crawler/handlers"
	"github.com/thewizardplusplus/go-crawler/register"
	"github.com/thewizardplusplus/go-crawler/sanitizing"
	htmlselector "github.com/thewizardplusplus/go-html-selector"
	"github.com/thewizardplusplus/go-crawler/models"
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
		var links []string
		switch request.URL.Path {
		case "/":
			links = []string{"/1", "/2", "/2", "https://golang.org/"}
		case "/1":
			links = []string{"/1/1", "/1/2"}
		case "/2":
			links = []string{"/2/1", "/2/2"}
		}
		for index := range links {
			if strings.HasPrefix(links[index], "/") {
				links[index] = "http://" + request.Host + links[index]
			}
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
```

`crawler.HandleLinksConcurrently()` with delaying extracting:

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
	"sync"
	"time"

	"github.com/go-log/log/print"
	crawler "github.com/thewizardplusplus/go-crawler"
	"github.com/thewizardplusplus/go-crawler/checkers"
	"github.com/thewizardplusplus/go-crawler/extractors"
	"github.com/thewizardplusplus/go-crawler/handlers"
	"github.com/thewizardplusplus/go-crawler/register"
	"github.com/thewizardplusplus/go-crawler/sanitizing"
	htmlselector "github.com/thewizardplusplus/go-html-selector"
	"github.com/thewizardplusplus/go-crawler/models"
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
		var links []string
		switch request.URL.Path {
		case "/":
			links = []string{"/1", "/2", "/2", "https://golang.org/"}
		case "/1":
			links = []string{"/1/1", "/1/2"}
		case "/2":
			links = []string{"/2/1", "/2/2"}
		}
		for index := range links {
			if strings.HasPrefix(links[index], "/") {
				links[index] = "http://" + request.Host + links[index]
			}
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
```

`crawler.HandleLinksConcurrently()` with processing a `robots.txt` file on extracting:

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
	"sync"
	"time"

	"github.com/go-log/log/print"
	crawler "github.com/thewizardplusplus/go-crawler"
	"github.com/thewizardplusplus/go-crawler/checkers"
	"github.com/thewizardplusplus/go-crawler/extractors"
	"github.com/thewizardplusplus/go-crawler/registers"
	htmlselector "github.com/thewizardplusplus/go-html-selector"
	"github.com/thewizardplusplus/go-crawler/models"
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
		for index := range links {
			if strings.HasPrefix(links[index], "/") {
				links[index] = "http://" + request.Host + links[index]
			}
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
```

`crawler.HandleLinksConcurrently()` with processing a `robots.txt` file on handling:

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
	"sync"
	"time"

	"github.com/go-log/log/print"
	crawler "github.com/thewizardplusplus/go-crawler"
	"github.com/thewizardplusplus/go-crawler/checkers"
	"github.com/thewizardplusplus/go-crawler/extractors"
	"github.com/thewizardplusplus/go-crawler/handlers"
	"github.com/thewizardplusplus/go-crawler/registers"
	htmlselector "github.com/thewizardplusplus/go-html-selector"
	"github.com/thewizardplusplus/go-crawler/models"
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
		for index := range links {
			if strings.HasPrefix(links[index], "/") {
				links[index] = "http://" + request.Host + links[index]
			}
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
```

`crawler.HandleLinksConcurrently()` with processing a `sitemap.xml` file:

```go
package main

import (
	"context"
	"fmt"
	"io"
	stdlog "log"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"runtime"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/go-log/log/print"
	"github.com/thewizardplusplus/go-crawler"
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
			links := []string{"/sitemap_from_robots_txt.xml"}
			completeLinksWithHost(links, request.Host)

			// nolint: errcheck
			fmt.Fprintf(
				writer,
				`
					User-agent: go-crawler
					Disallow: /2

					Sitemap: %s
				`,
				links[0],
			)

			return
		}

		if request.URL.Path == "/sitemap.xml" {
			links := []string{"/1", "/2", "/hidden/1", "/hidden/2"}
			completeLinksWithHost(links, request.Host)
			renderSitemap(writer, links) // nolint: errcheck

			return
		}
		if request.URL.Path == "/sitemap_from_robots_txt.xml" {
			links := []string{"/hidden/3", "/hidden/4"}
			completeLinksWithHost(links, request.Host)
			renderSitemap(writer, links) // nolint: errcheck

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

func renderTemplate(writer io.Writer, data interface{}, text string) error {
	template, err := template.New("").Parse(text)
	if err != nil {
		return err
	}

	return template.Execute(writer, data)
}

// nolint: unparam
func renderSitemap(writer io.Writer, links []string) error {
	return renderTemplate(writer, links, `
		<?xml version="1.0" encoding="UTF-8" ?>
		<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
			{{ range $link := . }}
				<url>
					<loc>{{ $link }}</loc>
				</url>
			{{ end }}
		</urlset>
	`)
}

func main() {
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
								sitemap.GeneratorGroup{
									sitemap.SimpleLinkGenerator{},
									sitemap.RobotsTXTGenerator{
										RobotsTXTRegister: registers.NewRobotsTXTRegister(http.DefaultClient),
									},
								},
								wrappedLogger,
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
	// have got the link "http://example.com/hidden/3" from the page "http://example.com"
	// have got the link "http://example.com/hidden/4" from the page "http://example.com"
	// have got the link "https://golang.org/" from the page "http://example.com"
}
```

## License

The MIT License (MIT)

Copyright &copy; 2020-2021 thewizardplusplus
