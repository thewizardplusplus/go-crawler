package transformers

import (
	"github.com/thewizardplusplus/go-crawler/models"
)

//go:generate mockery --name=LinkTransformer --inpackage --case=underscore --testonly

// LinkTransformer ...
//
// It's used only for mock generating.
//
type LinkTransformer interface {
	models.LinkTransformer
}
