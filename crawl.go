package crawler

import (
	"context"
	"sync"

	"github.com/go-log/log"
	"github.com/thewizardplusplus/go-crawler/handlers"
	"github.com/thewizardplusplus/go-crawler/models"
	syncutils "github.com/thewizardplusplus/go-sync-utils"
)

// ConcurrencyConfig ...
type ConcurrencyConfig struct {
	ConcurrencyFactor int
	BufferSize        int
}

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
	concurrencyConfig ConcurrencyConfig,
	links []string,
	dependencies CrawlDependencies,
) {
	linkChannel := make(chan string, concurrencyConfig.BufferSize)
	for _, link := range links {
		// use unbounded sending to avoid a deadlock
		syncutils.UnboundedSend(linkChannel, link)
	}

	var waiter sync.WaitGroup
	waiter.Add(len(links))

	HandleLinksConcurrently(
		ctx,
		concurrencyConfig.ConcurrencyFactor,
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

// CrawlByConcurrentHandler ...
func CrawlByConcurrentHandler(
	ctx context.Context,
	concurrencyFactor int,
	bufferSize int,
	handlerConcurrencyFactor int,
	handlerBufferSize int,
	links []string,
	dependencies CrawlDependencies,
) {
	concurrentHandler :=
		handlers.NewConcurrentHandler(handlerBufferSize, dependencies.LinkHandler)
	go concurrentHandler.RunConcurrently(ctx, handlerConcurrencyFactor)
	defer concurrentHandler.Stop()

	concurrencyConfig := ConcurrencyConfig{
		ConcurrencyFactor: concurrencyFactor,
		BufferSize:        bufferSize,
	}
	Crawl(ctx, concurrencyConfig, links, CrawlDependencies{
		LinkExtractor: dependencies.LinkExtractor,
		LinkChecker:   dependencies.LinkChecker,
		LinkHandler:   concurrentHandler,
		Logger:        dependencies.Logger,
	})
}
