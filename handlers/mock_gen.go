package handlers

import (
	"github.com/thewizardplusplus/go-crawler/models"
)

//go:generate mockery --name=LinkChecker --inpackage --case=underscore --testonly

// LinkChecker ...
//
// It's used only for mock generating.
//
type LinkChecker interface {
	models.LinkChecker
}

//go:generate mockery --name=LinkHandler --inpackage --case=underscore --testonly

// LinkHandler ...
//
// It's used only for mock generating.
//
type LinkHandler interface {
	models.LinkHandler
}

//go:generate mockery --name=ContextCancellerInterface --inpackage --case=underscore --testonly

// ContextCancellerInterface ...
//
// It is used only for mock generating.
//
type ContextCancellerInterface interface {
	CancelContext()
}
