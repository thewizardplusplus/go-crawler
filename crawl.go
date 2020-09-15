package crawler

import (
	"context"
	"sync"
)

// Crawl ...
func Crawl(
	ctx context.Context,
	concurrencyFactor int,
	links []string,
	dependencies Dependencies,
) {
	linkChannel := make(chan string, len(links))
	for _, link := range links {
		linkChannel <- link
	}

	var waiter sync.WaitGroup
	waiter.Add(len(links))

	HandleLinksConcurrently(ctx, concurrencyFactor, linkChannel, Dependencies{
		Waiter:        &waiter,
		LinkExtractor: dependencies.LinkExtractor,
		LinkChecker:   dependencies.LinkChecker,
		LinkHandler:   dependencies.LinkHandler,
		Logger:        dependencies.Logger,
	})

	waiter.Wait()
}
