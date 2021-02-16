package crawler

import (
	"context"
	"sync"

	"github.com/go-log/log"
	"github.com/thewizardplusplus/go-crawler/models"
	syncutils "github.com/thewizardplusplus/go-sync-utils"
)

//go:generate mockery --name=LinkExtractor --inpackage --case=underscore --testonly

// LinkExtractor ...
type LinkExtractor interface {
	ExtractLinks(ctx context.Context, threadID int, link string) ([]string, error)
}

//go:generate mockery --name=LinkChecker --inpackage --case=underscore --testonly

// LinkChecker ...
type LinkChecker interface {
	CheckLink(ctx context.Context, link models.SourcedLink) bool
}

//go:generate mockery --name=LinkHandler --inpackage --case=underscore --testonly

// LinkHandler ...
type LinkHandler interface {
	HandleLink(ctx context.Context, link models.SourcedLink)
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
