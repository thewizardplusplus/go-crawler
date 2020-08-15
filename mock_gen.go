package crawler

import (
	"github.com/go-log/log"
	"github.com/thewizardplusplus/go-crawler/waiter"
)

//go:generate mockery -name=Logger -inpkg -case=underscore -testonly

// Logger ...
//
// It's used only for mock generating.
type Logger interface {
	log.Logger
}

//go:generate mockery -name=Waiter -inpkg -case=underscore -testonly

// Waiter ...
//
// It's used only for mock generating.
type Waiter interface {
	waiter.Waiter
}
