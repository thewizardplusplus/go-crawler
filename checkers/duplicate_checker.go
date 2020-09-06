package checkers

import (
	"net/url"
	"path"
	"sync"

	mapset "github.com/deckarep/golang-set"
	"github.com/go-log/log"
	"github.com/thewizardplusplus/go-crawler/sanitizing"
)

// DuplicateChecker ...
type DuplicateChecker struct {
	sanitizeLink sanitizing.LinkSanitizing
	logger       log.Logger

	locker       sync.RWMutex
	checkedLinks mapset.Set
}

// NewDuplicateChecker ...
func NewDuplicateChecker(
	sanitizeLink sanitizing.LinkSanitizing,
	logger log.Logger,
) *DuplicateChecker {
	return &DuplicateChecker{
		sanitizeLink: sanitizeLink,
		logger:       logger,

		checkedLinks: mapset.NewThreadUnsafeSet(),
	}
}

// CheckLink ...
func (checker *DuplicateChecker) CheckLink(
	sourceLink string,
	link string,
) bool {
	if checker.sanitizeLink == sanitizing.SanitizeLink {
		var err error
		link, err = sanitizeLink(link)
		if err != nil {
			checker.logger.Logf("unable to parse the link: %s", err)
			return false
		}
	}

	checker.locker.Lock()
	defer checker.locker.Unlock()

	isDuplicate := checker.checkedLinks.Contains(link)
	if !isDuplicate {
		checker.checkedLinks.Add(link)
	}

	return !isDuplicate
}

func sanitizeLink(link string) (string, error) {
	parsedLink, err := url.Parse(link)
	if err != nil {
		return "", err
	}

	parsedLink.Path = path.Clean(parsedLink.Path)
	return parsedLink.String(), nil
}
