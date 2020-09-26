package crawler

import (
	"github.com/go-log/log"
	syncutils "github.com/thewizardplusplus/go-sync-utils"
)

//go:generate mockery -name=Waiter -inpkg -case=underscore -testonly

// Waiter ...
//
// It's used only for mock generating.
//
type Waiter interface {
	syncutils.WaitGroup
}

//go:generate mockery -name=Logger -inpkg -case=underscore -testonly

// Logger ...
//
// It's used only for mock generating.
//
type Logger interface {
	log.Logger
}
