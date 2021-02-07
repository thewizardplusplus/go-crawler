package checkers

import (
	"github.com/go-log/log"
	crawler "github.com/thewizardplusplus/go-crawler"
	httputils "github.com/thewizardplusplus/go-http-utils"
)

//go:generate mockery -name=LinkChecker -inpkg -case=underscore -testonly

// LinkChecker ...
//
// It's used only for mock generating.
//
type LinkChecker interface {
	crawler.LinkChecker
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
	httputils.HTTPClient
}
