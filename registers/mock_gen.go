package registers

import (
	"context"

	"github.com/go-log/log"
	"github.com/thewizardplusplus/go-crawler/models"
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

//go:generate mockery --name=LinkExtractor --inpackage --case=underscore --testonly

// LinkExtractor ...
//
// It's used only for mock generating.
//
type LinkExtractor interface {
	models.LinkExtractor
}

//go:generate mockery --name=LinkLoader --inpackage --case=underscore --testonly

// LinkLoader ...
//
// It's used only for mock generating.
//
type LinkLoader interface {
	LoadLink(link string, options interface{}) ([]byte, error)
}

//go:generate mockery --name=RegisteringHandler --inpackage --case=underscore --testonly

// RegisteringHandler ...
//
// It's used only for mock generating.
//
type RegisteringHandler interface {
	HandleRegistering(ctx context.Context, key interface{}) (
		value interface{},
		err error,
	)
}
