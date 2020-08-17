package checkers

import (
	"net/url"

	"github.com/go-log/log"
)

// HostChecker ...
type HostChecker struct {
	Logger log.Logger
}

// CheckLink ...
func (checker HostChecker) CheckLink(parentLink string, link string) bool {
	parsedParentLink, err := url.Parse(parentLink)
	if err != nil {
		checker.Logger.Logf("unable to parse the parent link: %s", err)
		return false
	}

	parsedLink, err := url.Parse(link)
	if err != nil {
		checker.Logger.Logf("unable to parse the link: %s", err)
		return false
	}

	return parsedLink.Host == parsedParentLink.Host
}
