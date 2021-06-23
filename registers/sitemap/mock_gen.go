package sitemap

import (
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
