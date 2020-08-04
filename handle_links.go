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
