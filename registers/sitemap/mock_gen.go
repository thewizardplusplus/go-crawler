package sitemap

import (
	"github.com/thewizardplusplus/go-crawler/registers"
)

//go:generate mockery --name=LinkGenerator --inpackage --case=underscore --testonly

// LinkGenerator ...
//
// It's used only for mock generating.
//
type LinkGenerator interface {
	registers.LinkGenerator
}
