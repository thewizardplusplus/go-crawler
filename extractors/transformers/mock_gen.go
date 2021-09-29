package transformers

import (
	"github.com/thewizardplusplus/go-crawler/extractors"
)

//go:generate mockery --name=LinkTransformer --inpackage --case=underscore --testonly

// LinkTransformer ...
//
// It's used only for mock generating.
//
type LinkTransformer interface {
	extractors.LinkTransformer
}
