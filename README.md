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

func (handler LinkHandler) HandleLink(link string) {
	// replace the test server URL for reproducibility of the example
	link = strings.Replace(link, handler.ServerURL, "http://example.com", -1)

	fmt.Printf("have got the link: %s\n", link)
}

func main() {
	server := httptest.NewServer(http.HandlerFunc(func(
		writer http.ResponseWriter,
		request *http.Request,
	) {
		if request.URL.Path != "/common" {
			fmt.Fprintf( // nolint: errcheck
				writer,
				`<p><a href="http://%[1]s/common">common</a></p>`,
				request.Host,
			)
		}
		if request.URL.Path != "/" {
			return
		}

		fmt.Fprintf( // nolint: errcheck
			writer,
			`<ul>
				<li><a href="http://%[1]s/1">1</a></li>
				<li><a href="http://%[1]s/2">2</a></li>
			</ul>`,
			request.Host,
		)
	}))
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
	// have got the link: http://example.com
	// have got the link: http://example.com/1
	// have got the link: http://example.com/2
	// have got the link: http://example.com/common
	// have got the link: http://example.com/common
	// have got the link: http://example.com/common
}
```

## License

The MIT License (MIT)

Copyright &copy; 2020 thewizardplusplus
