package sitemap

import (
	"github.com/thewizardplusplus/go-crawler/registers"
	httputils "github.com/thewizardplusplus/go-http-utils"
)

//go:generate mockery --name=LinkGenerator --inpackage --case=underscore --testonly

// LinkGenerator ...
//
// It's used only for mock generating.
//
type LinkGenerator interface {
	registers.LinkGenerator
}

//go:generate mockery --name=HTTPClient --inpackage --case=underscore --testonly

// HTTPClient ...
//
// It's used only for mock generating.
//
type HTTPClient interface {
	httputils.HTTPClient
}
