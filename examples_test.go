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
	"github.com/thewizardplusplus/go-crawler/register"
	"github.com/thewizardplusplus/go-crawler/sanitizing"
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
				checkers.DuplicateChecker{
					LinkRegister: register.NewLinkRegister(
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
				checkers.DuplicateChecker{
					LinkRegister: register.NewLinkRegister(
						sanitizing.SanitizeLink,
						wrappedLogger,
					),
				},
			},
			LinkHandler: handlers.NewUniqueHandler(
				// don't use here the link register from the duplicate checker above
				register.NewLinkRegister(sanitizing.SanitizeLink, wrappedLogger),
				LinkHandler{ServerURL: server.URL},
			),
			Logger: wrappedLogger,
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
