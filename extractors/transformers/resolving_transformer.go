package transformers

import (
	"bytes"

	htmlselector "github.com/thewizardplusplus/go-html-selector"
)

func selectBaseTag(data []byte) string {
	var builder BaseTagBuilder
	htmlselector.SelectTags( // nolint: errcheck, gosec
		bytes.NewReader(data),
		BaseTagFilters,
		&builder,
		htmlselector.SkipEmptyTags(),
		htmlselector.SkipEmptyAttributes(),
	)

	baseLink, _ := builder.BaseLink()
	return string(baseLink)
}
