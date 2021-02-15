package handlers

import (
	"context"

	crawler "github.com/thewizardplusplus/go-crawler"
)

// ConcurrentHandler ...
type ConcurrentHandler struct {
	linkHandler crawler.LinkHandler

	links chan crawler.SourcedLink
}

// NewConcurrentHandler ...
func NewConcurrentHandler(
	bufferSize int,
	linkHandler crawler.LinkHandler,
) ConcurrentHandler {
	return ConcurrentHandler{
		linkHandler: linkHandler,

		links: make(chan crawler.SourcedLink, bufferSize),
	}
}

// HandleLink ...
func (handler ConcurrentHandler) HandleLink(
	ctx context.Context,
	link crawler.SourcedLink,
) {
	handler.links <- link
}
