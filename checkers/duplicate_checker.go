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
	sanitizeLink bool
	logger       log.Logger

	checkedLinks mapset.Set
}

// NewDuplicateChecker ...
func NewDuplicateChecker(
	sanitizeLink bool,
	logger log.Logger,
) DuplicateChecker {
	return DuplicateChecker{
		sanitizeLink: sanitizeLink,
		logger:       logger,

		checkedLinks: mapset.NewSet(),
	}
}

// CheckLink ...
func (checker DuplicateChecker) CheckLink(parentLink string, link string) bool {
	if checker.sanitizeLink {
		var err error
		link, err = sanitizeLink(link)
		if err != nil {
			checker.logger.Logf("unable to sanitize the link: %s", err)
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
