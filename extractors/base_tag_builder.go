package extractors

import (
	htmlselector "github.com/thewizardplusplus/go-html-selector"
)

// BaseTagFilters ...
var BaseTagFilters = htmlselector.OptimizeFilters(htmlselector.FilterGroup{
	"base": {"href"},
})

// BaseTagBuilder ...
type BaseTagBuilder struct {
	baseLink     []byte
	isFirstFound bool
}
