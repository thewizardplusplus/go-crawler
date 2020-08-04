package crawler

import (
	"context"
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
