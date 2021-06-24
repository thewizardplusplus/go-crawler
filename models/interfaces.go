package models

import (
	"context"
)

// LinkExtractor ...
type LinkExtractor interface {
	ExtractLinks(ctx context.Context, threadID int, link string) ([]string, error)
}

// LinkChecker ...
type LinkChecker interface {
	CheckLink(ctx context.Context, link SourcedLink) bool
}

// LinkHandler ...
type LinkHandler interface {
	HandleLink(ctx context.Context, link SourcedLink)
}
