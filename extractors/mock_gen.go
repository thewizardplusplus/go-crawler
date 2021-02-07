package extractors

import (
	"time"

	"github.com/go-log/log"
	crawler "github.com/thewizardplusplus/go-crawler"
	httputils "github.com/thewizardplusplus/go-http-utils"
)

//go:generate mockery -name=HTTPClient -inpkg -case=underscore -testonly

// HTTPClient ...
//
// It's used only for mock generating.
//
type HTTPClient interface {
	httputils.HTTPClient
}

//go:generate mockery -name=LinkExtractor -inpkg -case=underscore -testonly

// LinkExtractor ...
//
// It's used only for mock generating.
//
type LinkExtractor interface {
	crawler.LinkExtractor
}

//go:generate mockery -name=Logger -inpkg -case=underscore -testonly

// Logger ...
//
// It's used only for mock generating.
//
type Logger interface {
	log.Logger
}

//go:generate mockery -name=Sleeper -inpkg -case=underscore -testonly

// Sleeper ...
//
// It's used only for mock generating.
//
type Sleeper interface {
	Sleep(duration time.Duration)
}
