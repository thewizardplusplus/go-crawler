package handlers

import (
	"context"

	crawler "github.com/thewizardplusplus/go-crawler"
)

// CheckedHandler ...
type CheckedHandler struct {
	LinkChecker crawler.LinkChecker
	LinkHandler crawler.LinkHandler
}

// HandleLink ...
func (handler CheckedHandler) HandleLink(
	ctx context.Context,
	link crawler.SourcedLink,
) {
	if !handler.LinkChecker.CheckLink(ctx, link) {
		return
	}

	handler.LinkHandler.HandleLink(ctx, link)
}
