package handlers

import (
	"context"

	"github.com/thewizardplusplus/go-crawler/models"
	syncutils "github.com/thewizardplusplus/go-sync-utils"
)

type linkHandlerWrapper struct {
	linkHandler models.LinkHandler
}

func (wrapper linkHandlerWrapper) Handle(
	ctx context.Context,
	data interface{},
) {
	wrapper.linkHandler.HandleLink(ctx, data.(models.SourcedLink))
}

// ConcurrentHandler ...
type ConcurrentHandler struct {
	// do not use embedding to hide the Handle() method
	innerConcurrentHandler syncutils.ConcurrentHandler
}

// NewConcurrentHandler ...
func NewConcurrentHandler(
	bufferSize int,
	linkHandler models.LinkHandler,
) ConcurrentHandler {
	return ConcurrentHandler{
		innerConcurrentHandler: syncutils.NewConcurrentHandler(
			bufferSize,
			linkHandlerWrapper{
				linkHandler: linkHandler,
			},
		),
	}
}

// HandleLink ...
func (handler ConcurrentHandler) HandleLink(
	ctx context.Context,
	link models.SourcedLink,
) {
	handler.innerConcurrentHandler.Handle(link)
}

// Start ...
func (handler ConcurrentHandler) Start(ctx context.Context) {
	handler.innerConcurrentHandler.Start(ctx)
}

// StartConcurrently ...
func (handler ConcurrentHandler) StartConcurrently(
	ctx context.Context,
	concurrencyFactor int,
) {
	handler.innerConcurrentHandler.StartConcurrently(ctx, concurrencyFactor)
}

// Stop ...
func (handler ConcurrentHandler) Stop() {
	handler.innerConcurrentHandler.Stop()
}
