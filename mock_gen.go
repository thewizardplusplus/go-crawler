package crawler

import (
	"github.com/go-log/log"
	"github.com/thewizardplusplus/go-crawler/models"
	syncutils "github.com/thewizardplusplus/go-sync-utils"
)

//go:generate mockery --name=LinkExtractor --inpackage --case=underscore --testonly

// LinkExtractor ...
//
// It's used only for mock generating.
//
type LinkExtractor interface {
	models.LinkExtractor
}

//go:generate mockery --name=LinkChecker --inpackage --case=underscore --testonly

// LinkChecker ...
//
// It's used only for mock generating.
//
type LinkChecker interface {
	models.LinkChecker
}

//go:generate mockery --name=LinkHandler --inpackage --case=underscore --testonly

// LinkHandler ...
//
// It's used only for mock generating.
//
type LinkHandler interface {
	models.LinkHandler
}

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
