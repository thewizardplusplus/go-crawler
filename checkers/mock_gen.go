package checkers

import (
	"github.com/go-log/log"
	"github.com/thewizardplusplus/go-crawler/models"
	httputils "github.com/thewizardplusplus/go-http-utils"
)

//go:generate mockery --name=LinkChecker --inpackage --case=underscore --testonly

// LinkChecker ...
//
// It's used only for mock generating.
//
type LinkChecker interface {
	models.LinkChecker
}

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
