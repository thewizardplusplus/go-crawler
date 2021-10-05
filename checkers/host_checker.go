package checkers

import (
	"context"

	"github.com/go-log/log"
	"github.com/thewizardplusplus/go-crawler/models"
	urlutils "github.com/thewizardplusplus/go-crawler/url-utils"
)

// HostChecker ...
type HostChecker struct {
	ComparisonResult urlutils.ComparisonResult
	Logger           log.Logger
}

// CheckLink ...
func (checker HostChecker) CheckLink(
	ctx context.Context,
	link models.SourcedLink,
) bool {
	const logPrefix = "host checking"

	result, err := urlutils.CompareLinkHosts(link.SourceLink, link.Link)
	if err != nil {
		const logMessage = "%s: unable to compare the hosts of links %q and %q: %s"
		checker.Logger.Logf(logMessage, logPrefix, link.SourceLink, link.Link, err)

		return false
	}

	return result == checker.ComparisonResult
}
