package transformers

import (
	htmlselector "github.com/thewizardplusplus/go-html-selector"
	"github.com/thewizardplusplus/go-html-selector/builders"
)

// DefaultBaseTagFilters ...
var DefaultBaseTagFilters = htmlselector.OptimizeFilters(
	htmlselector.FilterGroup{
		"base": {"href"},
	},
)

// BaseTagSelection ...
type BaseTagSelection int

// ...
const (
	SelectFirstBaseTag BaseTagSelection = iota
	SelectLastBaseTag
)

// BaseTagBuilder ...
type BaseTagBuilder struct {
	builders.FlattenBuilder

	baseTagSelection BaseTagSelection
}

// NewBaseTagBuilder ...
func NewBaseTagBuilder(baseTagSelection BaseTagSelection) BaseTagBuilder {
	return BaseTagBuilder{
		baseTagSelection: baseTagSelection,
	}
}

// BaseLink ...
func (builder BaseTagBuilder) BaseLink() (baseLink []byte, isFound bool) {
	if len(builder.AttributeValues()) == 0 {
		return nil, false
	}

	baseLink = builder.AttributeValues()[len(builder.AttributeValues())-1]
	return baseLink, true
}

// IsSelectionTerminated ...
func (builder BaseTagBuilder) IsSelectionTerminated() bool {
	return builder.baseTagSelection == SelectFirstBaseTag &&
		len(builder.AttributeValues()) > 0
}
