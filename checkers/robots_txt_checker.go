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
	parsedLink, err := url.Parse(link.Link)
	if err != nil {
		checker.Logger.Logf("unable to parse the link: %s", err)
		return false
	}

	robotsTXTData, err :=
		checker.RobotsTXTRegister.RegisterRobotsTXT(ctx, link.Link)
	if err != nil {
		checker.Logger.Logf("unable to register the robots.txt link: %s", err)
		return false
	}

	group := robotsTXTData.FindGroup(checker.UserAgent)
	return group.Test(parsedLink.Path)
}
