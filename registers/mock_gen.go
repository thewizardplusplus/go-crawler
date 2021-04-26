package registers

import (
	"time"

	"github.com/go-log/log"
	httputils "github.com/thewizardplusplus/go-http-utils"
)

//go:generate mockery --name=Logger --inpackage --case=underscore --testonly

// Logger ...
//
// It's used only for mock generating.
//
type Logger interface {
	log.Logger
}

//go:generate mockery --name=HTTPClient --inpackage --case=underscore --testonly

// HTTPClient ...
//
// It's used only for mock generating.
//
type HTTPClient interface {
	httputils.HTTPClient
}

//go:generate mockery --name=LinkLoader --inpackage --case=underscore --testonly

// LinkLoader ...
//
// It's used only for mock generating.
//
type LinkLoader interface {
	LoadLink(link string, options interface{}) ([]byte, error)
}

//go:generate mockery --name=Sleeper --inpackage --case=underscore --testonly

// Sleeper ...
//
// It's used only for mock generating.
//
type Sleeper interface {
	Sleep(duration time.Duration)
}
