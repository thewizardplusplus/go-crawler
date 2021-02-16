package models

import (
	"context"
)

// SourcedLink ...
type SourcedLink struct {
	SourceLink string
	Link       string
}

// LinkExtractor ...
type LinkExtractor interface {
	ExtractLinks(ctx context.Context, threadID int, link string) ([]string, error)
}
