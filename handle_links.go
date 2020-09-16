package crawler

import (
	"context"

	"github.com/thewizardplusplus/go-crawler/waiter"
)

// HandleLinkDependencies ...
type HandleLinkDependencies struct {
	CrawlDependencies

	Waiter waiter.Waiter
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
		go HandleLinks(ctx, links, dependencies)
	}
}

// HandleLinks ...
func HandleLinks(
	ctx context.Context,
	links chan string,
	dependencies HandleLinkDependencies,
) {
	for link := range links {
		extractedLinks := HandleLink(ctx, link, dependencies)
		for _, extractedLink := range extractedLinks {
			// simulate an unbounded channel to avoid a deadlock
			select {
			case links <- extractedLink:
			default:
				go func(link string) { links <- link }(extractedLink)
			}
		}
	}
}

// HandleLink ...
func HandleLink(
	ctx context.Context,
	link string,
	dependencies HandleLinkDependencies,
) []string {
	defer dependencies.Waiter.Done()

	extractedLinks, err := dependencies.LinkExtractor.ExtractLinks(ctx, link)
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
