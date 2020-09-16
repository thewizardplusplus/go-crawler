package crawler

import (
	"context"
	"sync"
)

// Crawl ...
func Crawl(
	ctx context.Context,
	concurrencyFactor int,
	bufferSize int,
	links []string,
	dependencies HandleLinkDependencies,
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
