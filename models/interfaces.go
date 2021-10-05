package models

import (
	"context"
	"net/http"
)

// LinkExtractor ...
type LinkExtractor interface {
	ExtractLinks(ctx context.Context, threadID int, link string) ([]string, error)
}

// LinkTransformer ...
type LinkTransformer interface {
	TransformLinks(
		links []string,
		response *http.Response,
		responseContent []byte,
	) ([]string, error)
}

// LinkChecker ...
type LinkChecker interface {
	CheckLink(ctx context.Context, link SourcedLink) bool
}

// LinkHandler ...
type LinkHandler interface {
	HandleLink(ctx context.Context, link SourcedLink)
}
