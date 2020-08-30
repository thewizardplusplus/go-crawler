package checkers

import (
	"github.com/go-log/log"
	crawler "github.com/thewizardplusplus/go-crawler"
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
