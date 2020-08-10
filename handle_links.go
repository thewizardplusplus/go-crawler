package crawler

import (
	"context"
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
	CheckLink(link string) bool
}

//go:generate mockery -name=LinkHandler -inpkg -case=underscore -testonly

// LinkHandler ...
type LinkHandler interface {
	HandleLink(link string)
}

//go:generate mockery -name=ErrorHandler -inpkg -case=underscore -testonly

// ErrorHandler ...
type ErrorHandler interface {
	HandleError(err error)
}

// Dependencies ...
type Dependencies struct {
	Waiter        Waiter
	LinkExtractor LinkExtractor
	LinkChecker   LinkChecker
	LinkHandler   LinkHandler
	ErrorHandler  ErrorHandler
}

// HandleLinks ...
func HandleLinks(
	ctx context.Context,
	links chan string,
	dependencies Dependencies,
) {
	for link := range links {
		HandleLink(ctx, links, link, dependencies)
	}
}

// HandleLink ...
func HandleLink(
	ctx context.Context,
	links chan string,
	link string,
	dependencies Dependencies,
) {
	defer dependencies.Waiter.Done()

	dependencies.LinkHandler.HandleLink(link)

	extractedLinks, err := dependencies.LinkExtractor.ExtractLinks(ctx, link)
	if err != nil {
		dependencies.ErrorHandler.HandleError(err)
		return
	}

	for _, link := range extractedLinks {
		if !dependencies.LinkChecker.CheckLink(link) {
			continue
		}

		dependencies.Waiter.Add(1)
		links <- link
	}
}
