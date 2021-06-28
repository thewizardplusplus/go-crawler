package checkers

import (
	"context"
	"net/url"

	"github.com/go-log/log"
	"github.com/thewizardplusplus/go-crawler/models"
)

// HostChecker ...
type HostChecker struct {
	Logger log.Logger
}

// CheckLink ...
func (checker HostChecker) CheckLink(
	ctx context.Context,
	link models.SourcedLink,
) bool {
	const logPrefix = "host checking"

	parsedSourceLink, err := url.Parse(link.SourceLink)
	if err != nil {
		const logMessage = "%s: unable to parse parent link %q: %s"
		checker.Logger.Logf(logMessage, logPrefix, link.SourceLink, err)

		return false
	}

	parsedLink, err := url.Parse(link.Link)
	if err != nil {
		const logMessage = "%s: unable to parse link %q: %s"
		checker.Logger.Logf(logMessage, logPrefix, link.Link, err)

		return false
	}

	return parsedLink.Host == parsedSourceLink.Host
}
