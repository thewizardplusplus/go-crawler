package urlutils

import (
	"strings"
	"unicode"
)

// LinkTrimming ...
type LinkTrimming int

// ...
const (
	DoNotTrimLink LinkTrimming = iota
	TrimLinkLeft
	TrimLinkRight
	TrimLink
)

// ApplyLinkTrimming ...
func ApplyLinkTrimming(link string, trimming LinkTrimming) string {
	switch trimming {
	case TrimLinkLeft:
		link = strings.TrimLeftFunc(link, unicode.IsSpace)
	case TrimLinkRight:
		link = strings.TrimRightFunc(link, unicode.IsSpace)
	case TrimLink:
		link = strings.TrimSpace(link)
	}

	return link
}
