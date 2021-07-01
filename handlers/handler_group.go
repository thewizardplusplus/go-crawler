package handlers

import (
	"context"
	"sync"

	"github.com/thewizardplusplus/go-crawler/models"
)

// HandlerGroup ...
type HandlerGroup []models.LinkHandler

// HandleLink ...
func (handlers HandlerGroup) HandleLink(
	ctx context.Context,
	link models.SourcedLink,
) {
	var waiter sync.WaitGroup
	waiter.Add(len(handlers))

	for _, handler := range handlers {
		go func(handler models.LinkHandler) {
			defer waiter.Done()

			handler.HandleLink(ctx, link)
		}(handler)
	}

	waiter.Wait()
}
