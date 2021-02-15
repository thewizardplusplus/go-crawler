package handlers

import (
	crawler "github.com/thewizardplusplus/go-crawler"
)

// ConcurrentHandler ...
type ConcurrentHandler struct {
	linkHandler crawler.LinkHandler

	links chan crawler.SourcedLink
}
