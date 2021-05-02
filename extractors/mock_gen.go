package extractors

import (
	"time"

	"github.com/go-log/log"
	"github.com/thewizardplusplus/go-crawler/models"
	"github.com/thewizardplusplus/go-crawler/registers"
	httputils "github.com/thewizardplusplus/go-http-utils"
)

//go:generate mockery --name=HTTPClient --inpackage --case=underscore --testonly

// HTTPClient ...
//
// It's used only for mock generating.
//
type HTTPClient interface {
	httputils.HTTPClient
}

//go:generate mockery --name=LinkExtractor --inpackage --case=underscore --testonly

// LinkExtractor ...
//
// It's used only for mock generating.
//
type LinkExtractor interface {
	models.LinkExtractor
}

//go:generate mockery --name=Logger --inpackage --case=underscore --testonly

// Logger ...
//
// It's used only for mock generating.
//
type Logger interface {
	log.Logger
}

//go:generate mockery --name=LinkGenerator --inpackage --case=underscore --testonly

// LinkGenerator ...
//
// It's used only for mock generating.
//
type LinkGenerator interface {
	registers.LinkGenerator
}

//go:generate mockery --name=Sleeper --inpackage --case=underscore --testonly

// Sleeper ...
//
// It's used only for mock generating.
//
type Sleeper interface {
	Sleep(duration time.Duration)
}
