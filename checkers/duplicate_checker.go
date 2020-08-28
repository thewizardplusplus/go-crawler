package checkers

import (
	"net/url"
	"path"
	"sync"

	mapset "github.com/deckarep/golang-set"
	"github.com/go-log/log"
	"github.com/pkg/errors"
)

// LinkSanitizing ...
type LinkSanitizing int

// ...
const (
	DoNotSanitizeLink LinkSanitizing = iota
	SanitizeLink
)

// DuplicateChecker ...
type DuplicateChecker struct {
	sanitizeLink LinkSanitizing
	logger       log.Logger

	locker       sync.RWMutex
	checkedLinks mapset.Set
}

// NewDuplicateChecker ...
func NewDuplicateChecker(
	sanitizeLink LinkSanitizing,
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
	parentLink string,
	link string,
) bool {
	if checker.sanitizeLink == SanitizeLink {
		var err error
		link, err = sanitizeLink(link)
		if err != nil {
			checker.logger.Logf("unable to sanitize the link: %s", err)
			return false
		}
	}

	checker.locker.Lock()
	defer checker.locker.Unlock()

	isDuplicate := checker.checkedLinks.Contains(link)
	if !isDuplicate {
		checker.checkedLinks.Add(link)
	}

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
