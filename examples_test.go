package crawler_test

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
	"github.com/thewizardplusplus/go-crawler/sanitizing"
	htmlselector "github.com/thewizardplusplus/go-html-selector"
)

type LinkHandler struct {
	ServerURL string
}

func (handler LinkHandler) HandleLink(link crawler.SourcedLink) {
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

func ExampleCrawl() {
	server := RunServer()
	defer server.Close()

	logger := stdlog.New(os.Stderr, "", stdlog.LstdFlags|stdlog.Lmicroseconds)
	// wrap the standard logger via the github.com/go-log/log package
	wrappedLogger := print.New(logger)

	crawler.Crawl(
		context.Background(),
		runtime.NumCPU(),
		1000,
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
			LinkHandler: handlers.UniqueHandler{
				// don't use here the link register from the duplicate checker above
				LinkRegister: registers.NewLinkRegister(
					sanitizing.SanitizeLink,
					wrappedLogger,
				),
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
				LinkHandler: handlers.UniqueHandler{
					// don't use here the link register from the duplicate checker above
					LinkRegister: registers.NewLinkRegister(
						sanitizing.SanitizeLink,
						wrappedLogger,
					),
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
				LinkHandler: handlers.UniqueHandler{
					// don't use here the link register from the duplicate checker above
					LinkRegister: registers.NewLinkRegister(
						sanitizing.SanitizeLink,
						wrappedLogger,
					),
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
