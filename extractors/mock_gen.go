package extractors

import (
	"time"

	"github.com/go-log/log"
	crawler "github.com/thewizardplusplus/go-crawler"
)

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

//go:generate mockery -name=SleeperInterface -inpkg -case=underscore -testonly

// SleeperInterface ...
//
// It's used only for mock generating.
//
type SleeperInterface interface {
	Sleep(duration time.Duration)
}
