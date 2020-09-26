package crawler

import (
	"context"
	"sync"

	"github.com/go-log/log"
)

//go:generate mockery -name=LinkExtractor -inpkg -case=underscore -testonly

// LinkExtractor ...
type LinkExtractor interface {
	ExtractLinks(ctx context.Context, link string) ([]string, error)
}

//go:generate mockery -name=LinkChecker -inpkg -case=underscore -testonly

// LinkChecker ...
type LinkChecker interface {
	CheckLink(sourceLink string, link string) bool
}

//go:generate mockery -name=LinkHandler -inpkg -case=underscore -testonly

// LinkHandler ...
type LinkHandler interface {
	HandleLink(sourceLink string, link string)
}

// CrawlDependencies ...
type CrawlDependencies struct {
	LinkExtractor LinkExtractor
	LinkChecker   LinkChecker
	LinkHandler   LinkHandler
	Logger        log.Logger
}

// Crawl ...
func Crawl(
	ctx context.Context,
	concurrencyFactor int,
	bufferSize int,
	links []string,
	dependencies CrawlDependencies,
) {
	linkChannel := make(chan string, bufferSize)
	go func() {
		for _, link := range links {
			linkChannel <- link
		}
	}()

	var waiter sync.WaitGroup
	waiter.Add(len(links))

	HandleLinksConcurrently(
		ctx,
		concurrencyFactor,
		linkChannel,
		HandleLinkDependencies{
			CrawlDependencies: dependencies,
			Waiter:            &waiter,
		},
	)

	waiter.Wait()
}