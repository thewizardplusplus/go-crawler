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

// GenerateLinks ...
func (generator RobotsTXTGenerator) GenerateLinks(
	ctx context.Context,
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
