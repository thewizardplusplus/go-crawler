package sitemap

import (
	"github.com/thewizardplusplus/go-crawler/models"
	httputils "github.com/thewizardplusplus/go-http-utils"
)

//go:generate mockery --name=LinkExtractor --inpackage --case=underscore --testonly

// LinkExtractor ...
//
// It's used only for mock generating.
//
type LinkExtractor interface {
	models.LinkExtractor
}

//go:generate mockery --name=HTTPClient --inpackage --case=underscore --testonly

// HTTPClient ...
//
// It's used only for mock generating.
//
type HTTPClient interface {
	httputils.HTTPClient
}
