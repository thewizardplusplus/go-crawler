package crawler

import (
	"context"
	"sync"

	"github.com/go-log/log"
	syncutils "github.com/thewizardplusplus/go-sync-utils"
)

// SourcedLink ...
type SourcedLink struct {
	SourceLink string
	Link       string
}

//go:generate mockery -name=LinkExtractor -inpkg -case=underscore -testonly

// LinkExtractor ...
type LinkExtractor interface {
	ExtractLinks(ctx context.Context, threadID int, link string) ([]string, error)
}

//go:generate mockery -name=LinkChecker -inpkg -case=underscore -testonly

// LinkChecker ...
type LinkChecker interface {
	CheckLink(link SourcedLink) bool
}

//go:generate mockery -name=LinkHandler -inpkg -case=underscore -testonly

// LinkHandler ...
type LinkHandler interface {
	HandleLink(link SourcedLink)
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
}
