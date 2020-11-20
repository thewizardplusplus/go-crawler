package checkers

import (
	"net/url"

	"github.com/go-log/log"
	crawler "github.com/thewizardplusplus/go-crawler"
)

// HostChecker ...
type HostChecker struct {
	Logger log.Logger
}

// CheckLink ...
func (checker HostChecker) CheckLink(link crawler.SourcedLink) bool {
	parsedSourceLink, err := url.Parse(link.SourceLink)
	if err != nil {
		checker.Logger.Logf("unable to parse the parent link: %s", err)
		return false
	}

	parsedLink, err := url.Parse(link.Link)
	if err != nil {
		checker.Logger.Logf("unable to parse the link: %s", err)
		return false
	}

	return parsedLink.Host == parsedSourceLink.Host
}
