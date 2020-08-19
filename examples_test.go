package crawler_test

import (
	"context"
	"fmt"
	stdlog "log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"

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

func ExampleHandleLinks() {
	server := httptest.NewServer(http.HandlerFunc(func(
		writer http.ResponseWriter,
		request *http.Request,
	) {
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

	go crawler.HandleLinks(context.Background(), links, crawler.Dependencies{
		Waiter: &waiter,
		LinkExtractor: extractors.DefaultExtractor{
			HTTPClient: http.DefaultClient,
			Filters: htmlselector.OptimizeFilters(htmlselector.FilterGroup{
				"a": {"href"},
			}),
		},
		LinkChecker: checkers.HostChecker{
			Logger: wrappedLogger,
		},
		LinkHandler: LinkHandler{
			ServerURL: server.URL,
		},
		Logger: wrappedLogger,
	})

	waiter.Wait()

	// Output:
	// have got the link: http://example.com
	// have got the link: http://example.com/1
	// have got the link: http://example.com/2
}
