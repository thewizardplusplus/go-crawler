package crawler

import (
	"context"
)

// LinkExtractor ...
type LinkExtractor interface {
	ExtractLinks(ctx context.Context, link string) ([]string, error)
}
