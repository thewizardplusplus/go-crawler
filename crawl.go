package crawler

import (
	"context"
	"sync"

	"github.com/go-log/log"
	"github.com/thewizardplusplus/go-crawler/models"
	syncutils "github.com/thewizardplusplus/go-sync-utils"
)

// CrawlDependencies ...
type CrawlDependencies struct {
	LinkExtractor models.LinkExtractor
	LinkChecker   models.LinkChecker
	LinkHandler   models.LinkHandler
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
	for _, link := range links {
		// use unbounded sending to avoid a deadlock
		syncutils.UnboundedSend(linkChannel, link)
	}

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
	// it should be called after the waiter.Wait() call
	close(linkChannel)
}
