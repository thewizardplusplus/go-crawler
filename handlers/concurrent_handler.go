package handlers

import (
	"context"
	"sync"

	"github.com/thewizardplusplus/go-crawler/models"
)

// ConcurrentHandler ...
type ConcurrentHandler struct {
	linkHandler models.LinkHandler

	links chan models.SourcedLink
}

// NewConcurrentHandler ...
func NewConcurrentHandler(
	bufferSize int,
	linkHandler models.LinkHandler,
) ConcurrentHandler {
	return ConcurrentHandler{
		linkHandler: linkHandler,

		links: make(chan models.SourcedLink, bufferSize),
	}
}

// HandleLink ...
func (handler ConcurrentHandler) HandleLink(
	ctx context.Context,
	link models.SourcedLink,
) {
	handler.links <- link
}

// Run ...
func (handler ConcurrentHandler) Run(ctx context.Context) {
	for link := range handler.links {
		handler.linkHandler.HandleLink(ctx, link)
	}
}

// RunConcurrently ...
func (handler ConcurrentHandler) RunConcurrently(
	ctx context.Context,
	concurrencyFactor int,
) {
	var waiter sync.WaitGroup
	waiter.Add(concurrencyFactor)

	for threadID := 0; threadID < concurrencyFactor; threadID++ {
		go func() {
			defer waiter.Done()

			handler.Run(ctx)
		}()
	}

	waiter.Wait()
}

// Stop ...
func (handler ConcurrentHandler) Stop() {
	close(handler.links)
}
