package handlers

import (
	"context"

	crawler "github.com/thewizardplusplus/go-crawler"
	"github.com/thewizardplusplus/go-crawler/models"
)

// CheckedHandler ...
type CheckedHandler struct {
	LinkChecker models.LinkChecker
	LinkHandler crawler.LinkHandler
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
