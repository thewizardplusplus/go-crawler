package crawler

import (
	"context"
	"sync"

	"github.com/go-log/log"
)

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
			Waiter:        &waiter,
			LinkExtractor: dependencies.LinkExtractor,
			LinkChecker:   dependencies.LinkChecker,
			LinkHandler:   dependencies.LinkHandler,
			Logger:        dependencies.Logger,
		},
	)

	waiter.Wait()
}
