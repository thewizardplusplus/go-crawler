package crawler

import (
	"github.com/go-log/log"
	syncutils "github.com/thewizardplusplus/go-sync-utils"
)

//go:generate mockery --name=Waiter --inpackage --case=underscore --testonly

// Waiter ...
//
// It's used only for mock generating.
//
type Waiter interface {
	syncutils.WaitGroup
}

//go:generate mockery --name=Logger --inpackage --case=underscore --testonly

// Logger ...
//
// It's used only for mock generating.
//
type Logger interface {
	log.Logger
}
