package sanitizing

import (
	"net/url"
	"path"
)

// LinkSanitizing ...
type LinkSanitizing int

// ...
const (
	DoNotSanitizeLink LinkSanitizing = iota
	SanitizeLink
)

// ApplyLinkSanitizing ...
func ApplyLinkSanitizing(link string) (string, error) {
	parsedLink, err := url.Parse(link)
	if err != nil {
		return "", err
	}

	parsedLink.Path = path.Clean(parsedLink.Path)
	return parsedLink.String(), nil
}
