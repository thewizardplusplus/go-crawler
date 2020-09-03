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
- calling of an outer handler for an each found link:
  - it's called directly during crawling;
- custom filtering of considered links:
  - by relativity of a link (optional);
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
)

type LinkHandler struct {
	ServerURL string
}

func (handler LinkHandler) HandleLink(sourceLink string, link string) {
	fmt.Printf(
		"have got the link %q from the page %q\n",
		handler.replaceServerURL(link),
		handler.replaceServerURL(sourceLink),
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
		crawler.Dependencies{
			Waiter: &waiter,
			LinkExtractor: extractors.RepeatingExtractor{
				LinkExtractor: extractors.DefaultExtractor{
					HTTPClient: http.DefaultClient,
					Filters: htmlselector.OptimizeFilters(htmlselector.FilterGroup{
						"a": {"href"},
					}),
				},
				RepeatCount: 5,
				RepeatDelay: time.Second,
				Logger:      wrappedLogger,
			},
			LinkChecker: checkers.HostChecker{
				Logger: wrappedLogger,
			},
			LinkHandler: LinkHandler{
				ServerURL: server.URL,
			},
			Logger: wrappedLogger,
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

`crawler.HandleLinksConcurrently()` without duplicates:

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
)

type LinkHandler struct {
	ServerURL string
}

func (handler LinkHandler) HandleLink(sourceLink string, link string) {
	fmt.Printf(
		"have got the link %q from the page %q\n",
		handler.replaceServerURL(link),
		handler.replaceServerURL(sourceLink),
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
		crawler.Dependencies{
			Waiter: &waiter,
			LinkExtractor: extractors.RepeatingExtractor{
				LinkExtractor: extractors.DefaultExtractor{
					HTTPClient: http.DefaultClient,
					Filters: htmlselector.OptimizeFilters(htmlselector.FilterGroup{
						"a": {"href"},
					}),
				},
				RepeatCount: 5,
				RepeatDelay: time.Second,
				Logger:      wrappedLogger,
			},
			LinkChecker: checkers.CheckerGroup{
				checkers.HostChecker{
					Logger: wrappedLogger,
				},
				checkers.NewDuplicateChecker(checkers.SanitizeLink, wrappedLogger),
			},
			LinkHandler: LinkHandler{
				ServerURL: server.URL,
			},
			Logger: wrappedLogger,
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

## License

The MIT License (MIT)

Copyright &copy; 2020 thewizardplusplus
