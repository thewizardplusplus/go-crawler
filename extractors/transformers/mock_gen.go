package transformers

import (
	"github.com/go-log/log"
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

//go:generate mockery --name=Logger --inpackage --case=underscore --testonly

// Logger ...
//
// It's used only for mock generating.
//
type Logger interface {
	log.Logger
}
