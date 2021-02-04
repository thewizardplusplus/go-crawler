package handlers

import (
	crawler "github.com/thewizardplusplus/go-crawler"
)

//go:generate mockery -name=LinkChecker -inpkg -case=underscore -testonly

// LinkChecker ...
//
// It's used only for mock generating.
//
type LinkChecker interface {
	crawler.LinkChecker
}

//go:generate mockery -name=LinkHandler -inpkg -case=underscore -testonly

// LinkHandler ...
//
// It's used only for mock generating.
//
type LinkHandler interface {
	crawler.LinkHandler
}
