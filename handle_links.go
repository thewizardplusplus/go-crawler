package crawler

import (
	"context"

	"github.com/go-log/log"
)

//go:generate mockery -name=Waiter -inpkg -case=underscore -testonly

// Waiter ...
type Waiter interface {
	Add(delta int)
	Done()
	Wait()
}

//go:generate mockery -name=LinkExtractor -inpkg -case=underscore -testonly

// LinkExtractor ...
type LinkExtractor interface {
	ExtractLinks(ctx context.Context, link string) ([]string, error)
}

//go:generate mockery -name=LinkChecker -inpkg -case=underscore -testonly

// LinkChecker ...
type LinkChecker interface {
	CheckLink(parentLink string, link string) bool
}

//go:generate mockery -name=LinkHandler -inpkg -case=underscore -testonly

// LinkHandler ...
type LinkHandler interface {
	HandleLink(link string)
}

// Dependencies ...
type Dependencies struct {
	Waiter        Waiter
	LinkExtractor LinkExtractor
	LinkChecker   LinkChecker
	LinkHandler   LinkHandler
	Logger        log.Logger
}

// HandleLinks ...
func HandleLinks(
	ctx context.Context,
	links chan string,
	dependencies Dependencies,
) {
	for link := range links {
		extractedLinks := HandleLink(ctx, link, dependencies)
		for _, extractedLink := range extractedLinks {
			links <- extractedLink
		}
	}
}

// HandleLink ...
func HandleLink(
	ctx context.Context,
	link string,
	dependencies Dependencies,
) []string {
	defer dependencies.Waiter.Done()

	dependencies.LinkHandler.HandleLink(link)

	extractedLinks, err := dependencies.LinkExtractor.ExtractLinks(ctx, link)
	if err != nil {
		dependencies.Logger.Logf("unable to extract links: %s", err)
		return nil
	}

	var checkedExtractedLinks []string
	for _, extractedLink := range extractedLinks {
		if !dependencies.LinkChecker.CheckLink(link, extractedLink) {
			continue
		}

		checkedExtractedLinks = append(checkedExtractedLinks, extractedLink)
		// it should be called before the dependencies.Waiter.Done() call
		dependencies.Waiter.Add(1)
	}

	return checkedExtractedLinks
}
