package checkers

import (
	"net/url"
	"path"

	mapset "github.com/deckarep/golang-set"
	"github.com/go-log/log"
	"github.com/pkg/errors"
)

// DuplicateChecker ...
type DuplicateChecker struct {
	SanitizeLink bool
	Logger       log.Logger

	checkedLinks mapset.Set
}

// CheckLink ...
func (checker DuplicateChecker) CheckLink(parentLink string, link string) bool {
	if checker.SanitizeLink {
		var err error
		link, err = sanitizeLink(link)
		if err != nil {
			checker.Logger.Logf("unable to sanitize the link: %s", err)
			return false
		}
	}

	isDuplicate := checker.checkedLinks.Contains(link)
	checker.checkedLinks.Add(link)

	return isDuplicate
}

func sanitizeLink(link string) (string, error) {
	parsedLink, err := url.Parse(link)
	if err != nil {
		return "", errors.Wrap(err, "unable to parse the link")
	}

	parsedLink.Path = path.Clean(parsedLink.Path)
	return parsedLink.String(), nil
}
