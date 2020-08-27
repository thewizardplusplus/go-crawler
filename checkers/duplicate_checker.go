package checkers

import (
	"net/url"
	"path"

	mapset "github.com/deckarep/golang-set"
	"github.com/go-log/log"
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
		parsedLink, err := url.Parse(link)
		if err != nil {
			checker.Logger.Logf("unable to parse the link: %s", err)
			return false
		}
		parsedLink.Path = path.Clean(parsedLink.Path)

		link = parsedLink.String()
	}

	isDuplicate := checker.checkedLinks.Contains(link)
	checker.checkedLinks.Add(link)

	return isDuplicate
}
