package checkers

import (
	"context"
	"net/url"

	"github.com/go-log/log"
	"github.com/thewizardplusplus/go-crawler/models"
	"github.com/thewizardplusplus/go-crawler/registers"
)

// RobotsTXTChecker ...
type RobotsTXTChecker struct {
	UserAgent         string
	RobotsTXTRegister registers.RobotsTXTRegister
	Logger            log.Logger
}

// CheckLink ...
func (checker RobotsTXTChecker) CheckLink(
	ctx context.Context,
	link models.SourcedLink,
) bool {
	const logPrefix = "robots.txt checking"

	parsedLink, err := url.Parse(link.Link)
	if err != nil {
		const logMessage = "%s: unable to parse link %q: %s"
		checker.Logger.Logf(logMessage, logPrefix, link.Link, err)

		return false
	}

	robotsTXTData, err :=
		checker.RobotsTXTRegister.RegisterRobotsTXT(ctx, link.Link)
	if err != nil {
		const logMessage = "%s: " +
			"unable to register the robots.txt link for link %q: " +
			"%s"
		checker.Logger.Logf(logMessage, logPrefix, link.Link, err)

		return false
	}

	group := robotsTXTData.FindGroup(checker.UserAgent)
	return group.Test(parsedLink.Path)
}
