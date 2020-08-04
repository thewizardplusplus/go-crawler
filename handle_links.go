package crawler

import (
	"context"
	"sync"
)

// LinkExtractor ...
type LinkExtractor interface {
	ExtractLinks(ctx context.Context, link string) ([]string, error)
}

// LinkChecker ...
type LinkChecker interface {
	CheckLink(link string) bool
}

// LinkHandler ...
type LinkHandler interface {
	HandleLink(link string)
}

// ErrorHandler ...
type ErrorHandler interface {
	HandleError(err error)
}

// Dependencies ...
type Dependencies struct {
	LinkExtractor LinkExtractor
	LinkChecker   LinkChecker
	LinkHandler   LinkHandler
	ErrorHandler  ErrorHandler
}

// HandleLinks ...
func HandleLinks(
	ctx context.Context,
	waiter *sync.WaitGroup,
	links chan string,
	dependencies Dependencies,
) {
	for link := range links {
		HandleLink(ctx, waiter, links, link, dependencies)
	}
}

// HandleLink ...
func HandleLink(
	ctx context.Context,
	waiter *sync.WaitGroup,
	links chan string,
	link string,
	dependencies Dependencies,
) {
	defer waiter.Done()

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

		waiter.Add(1)
		links <- link
	}
}
