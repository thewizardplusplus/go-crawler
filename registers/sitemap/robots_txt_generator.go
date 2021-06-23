package sitemap

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thewizardplusplus/go-crawler/registers"
)

// RobotsTXTGenerator ...
type RobotsTXTGenerator struct {
	RobotsTXTRegister registers.RobotsTXTRegister
}

// ExtractLinks ...
func (generator RobotsTXTGenerator) ExtractLinks(
	ctx context.Context,
	threadID int,
	baseLink string,
) (
	[]string,
	error,
) {
	robotsTXTData, err :=
		generator.RobotsTXTRegister.RegisterRobotsTXT(ctx, baseLink)
	if err != nil {
		return nil, errors.Wrap(err, "unable to register the robots.txt link")
	}

	return robotsTXTData.Sitemaps, nil
}
