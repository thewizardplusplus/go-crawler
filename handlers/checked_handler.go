package handlers

import (
	crawler "github.com/thewizardplusplus/go-crawler"
)

// CheckedHandler ...
type CheckedHandler struct {
	LinkChecker crawler.LinkChecker
	LinkHandler crawler.LinkHandler
}
