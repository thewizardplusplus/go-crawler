package crawler

import (
	"context"

	syncutils "github.com/thewizardplusplus/go-sync-utils"
)

// HandleLinkDependencies ...
type HandleLinkDependencies struct {
	CrawlDependencies

	Waiter syncutils.WaitGroup
}

// HandleLinksConcurrently ...
func HandleLinksConcurrently(
	ctx context.Context,
	concurrencyFactor int,
	links chan string,
	dependencies HandleLinkDependencies,
) {
	for i := 0; i < concurrencyFactor; i++ {
		// waiting for completion is done via dependencies.Waiter
		go HandleLinks(ctx, 0, links, dependencies)
	}
}

// HandleLinks ...
func HandleLinks(
	ctx context.Context,
	threadID int,
	links chan string,
	dependencies HandleLinkDependencies,
) {
	for link := range links {
		extractedLinks := HandleLink(ctx, threadID, link, dependencies)
		for _, extractedLink := range extractedLinks {
			// use unbounded sending to avoid a deadlock
			syncutils.UnboundedSend(links, extractedLink)
		}
	}
}

// HandleLink ...
func HandleLink(
	ctx context.Context,
	threadID int,
	link string,
	dependencies HandleLinkDependencies,
) []string {
	defer dependencies.Waiter.Done()

	extractedLinks, err :=
		dependencies.LinkExtractor.ExtractLinks(ctx, threadID, link)
	if err != nil {
		dependencies.Logger.Logf("unable to extract links: %s", err)
		return nil
	}

	var checkedExtractedLinks []string
	for _, extractedLink := range extractedLinks {
		dependencies.LinkHandler.HandleLink(link, extractedLink)

		if !dependencies.LinkChecker.CheckLink(link, extractedLink) {
			continue
		}

		checkedExtractedLinks = append(checkedExtractedLinks, extractedLink)
		// it should be called before the dependencies.Waiter.Done() call
		dependencies.Waiter.Add(1)
	}

	return checkedExtractedLinks
}
