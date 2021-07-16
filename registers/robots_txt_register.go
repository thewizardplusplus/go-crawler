package registers

import (
	"context"
	"net/http"

	"github.com/pkg/errors"
	"github.com/temoto/robotstxt"
	urlutils "github.com/thewizardplusplus/go-crawler/url-utils"
	httputils "github.com/thewizardplusplus/go-http-utils"
)

// RobotsTXTRegister ...
type RobotsTXTRegister struct {
	httpClient httputils.HTTPClient

	robotsTXTRegister BasicRegister
}

// NewRobotsTXTRegister ...
func NewRobotsTXTRegister(httpClient httputils.HTTPClient) RobotsTXTRegister {
	return RobotsTXTRegister{
		httpClient: httpClient,

		robotsTXTRegister: NewBasicRegister(),
	}
}

// RegisterRobotsTXT ...
func (register RobotsTXTRegister) RegisterRobotsTXT(
	ctx context.Context,
	link string,
) (
	*robotstxt.RobotsData,
	error,
) {
	robotsTXTLinks, err := urlutils.GenerateHierarchicalLinks(
		link,
		"robots.txt",
		urlutils.WithMaximalHierarchyDepth(0),
	)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create the robots.txt link")
	}

	// if successful, the result will always be one link
	robotsTXTLink := robotsTXTLinks[0]
	robotsTXTData, err := register.robotsTXTRegister.RegisterValue(
		ctx,
		robotsTXTLink,
		func(ctx context.Context, robotsTXTLink interface{}) (interface{}, error) {
			return register.loadRobotsTXTData(ctx, robotsTXTLink.(string))
		},
	)
	if err != nil {
		return nil, errors.Wrap(err, "unable to load the robots.txt data")
	}

	return robotsTXTData.(*robotstxt.RobotsData), nil
}

func (register RobotsTXTRegister) loadRobotsTXTData(
	ctx context.Context,
	robotsTXTLink string,
) (
	*robotstxt.RobotsData,
	error,
) {
	request, err := http.NewRequest(http.MethodGet, robotsTXTLink, nil)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create the request")
	}
	request = request.WithContext(ctx)

	response, err := register.httpClient.Do(request)
	if err != nil {
		return nil, errors.Wrap(err, "unable to send the request")
	}
	defer response.Body.Close() // nolint: errcheck

	robotsTXTData, err := robotstxt.FromResponse(response)
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse the response")
	}

	return robotsTXTData, nil
}
