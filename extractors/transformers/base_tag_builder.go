package transformers

import (
	htmlselector "github.com/thewizardplusplus/go-html-selector"
	byteutils "github.com/thewizardplusplus/go-html-selector/byte-utils"
)

// BaseTagFilters ...
var BaseTagFilters = htmlselector.OptimizeFilters(htmlselector.FilterGroup{
	"base": {"href"},
})

// BaseTagSelection ...
type BaseTagSelection int

// ...
const (
	SelectFirstBaseTag BaseTagSelection = iota
	SelectLastBaseTag
)

// BaseTagBuilder ...
type BaseTagBuilder struct {
	baseTagSelection BaseTagSelection
	baseLink         []byte
	isFirstFound     bool
}

// BaseLink ...
func (builder BaseTagBuilder) BaseLink() (baseLink []byte, isFound bool) {
	if !builder.isFirstFound {
		return nil, false
	}

	return builder.baseLink, true
}

// AddTag ...
func (builder BaseTagBuilder) AddTag(name []byte) {}

// AddAttribute ...
func (builder *BaseTagBuilder) AddAttribute(name []byte, value []byte) {
	builder.baseLink = byteutils.Copy(value)
	builder.isFirstFound = true
}

// IsSelectionTerminated ...
func (builder BaseTagBuilder) IsSelectionTerminated() bool {
	return builder.baseTagSelection == SelectFirstBaseTag && builder.isFirstFound
}
