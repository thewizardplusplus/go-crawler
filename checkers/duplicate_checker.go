package checkers

import (
	mapset "github.com/deckarep/golang-set"
	"github.com/go-log/log"
	"github.com/thewizardplusplus/go-crawler/sanitizing"
)

// DuplicateChecker ...
type DuplicateChecker struct {
	sanitizeLink sanitizing.LinkSanitizing
	logger       log.Logger

	checkedLinks mapset.Set
}

// NewDuplicateChecker ...
func NewDuplicateChecker(
	sanitizeLink sanitizing.LinkSanitizing,
	logger log.Logger,
) DuplicateChecker {
	return DuplicateChecker{
		sanitizeLink: sanitizeLink,
		logger:       logger,

		checkedLinks: mapset.NewSet(),
	}
}

// CheckLink ...
func (checker DuplicateChecker) CheckLink(
	sourceLink string,
	link string,
) bool {
	if checker.sanitizeLink == sanitizing.SanitizeLink {
		var err error
		link, err = sanitizing.ApplyLinkSanitizing(link)
		if err != nil {
			checker.logger.Logf("unable to sanitize the link: %s", err)
			return false
		}
	}

	wasAdded := checker.checkedLinks.Add(link)
	return wasAdded
}
