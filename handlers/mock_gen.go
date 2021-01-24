package handlers

import (
	"net/http"

	"github.com/go-log/log"
	crawler "github.com/thewizardplusplus/go-crawler"
)

//go:generate mockery -name=LinkHandler -inpkg -case=underscore -testonly

// LinkHandler ...
//
// It's used only for mock generating.
//
type LinkHandler interface {
	crawler.LinkHandler
}

//go:generate mockery -name=Logger -inpkg -case=underscore -testonly

// Logger ...
//
// It's used only for mock generating.
//
type Logger interface {
	log.Logger
}

//go:generate mockery -name=HTTPClient -inpkg -case=underscore -testonly

// HTTPClient ...
//
// It's used only for mock generating.
//
type HTTPClient interface {
	Do(request *http.Request) (*http.Response, error)
}
