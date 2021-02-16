package handlers

import (
	"context"

	"github.com/thewizardplusplus/go-crawler/models"
)

// CheckedHandler ...
type CheckedHandler struct {
	LinkChecker models.LinkChecker
	LinkHandler models.LinkHandler
}

// HandleLink ...
func (handler CheckedHandler) HandleLink(
	ctx context.Context,
	link models.SourcedLink,
) {
	if !handler.LinkChecker.CheckLink(ctx, link) {
		return
	}

	handler.LinkHandler.HandleLink(ctx, link)
}
