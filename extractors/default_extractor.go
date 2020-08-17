package extractors

import (
	"net/http"

	htmlselector "github.com/thewizardplusplus/go-html-selector"
)

// HTTPClient ...
type HTTPClient interface {
	Do(request *http.Request) (*http.Response, error)
}

// DefaultExtractor ...
type DefaultExtractor struct {
	HTTPClient HTTPClient
	Filters    htmlselector.OptimizedFilterGroup
}
