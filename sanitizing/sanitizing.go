package sanitizing

import (
	"net/url"
	"path"

	"github.com/pkg/errors"
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
		return "", errors.Wrap(err, "unable to parse the link")
	}

	parsedLink.Path = path.Clean(parsedLink.Path)
	return parsedLink.String(), nil
}
