package registers

import (
	"context"
	"net/http"
	"net/url"
	"sync"

	"github.com/pkg/errors"
	"github.com/temoto/robotstxt"
)

//go:generate mockery -name=HTTPClient -inpkg -case=underscore -testonly

// HTTPClient ...
type HTTPClient interface {
	Do(request *http.Request) (*http.Response, error)
}

// RobotsTXTRegister ...
type RobotsTXTRegister struct {
	httpClient HTTPClient

	registeredRobotsTXT *sync.Map
}

// NewRobotsTXTRegister ...
func NewRobotsTXTRegister(httpClient HTTPClient) RobotsTXTRegister {
	return RobotsTXTRegister{
		httpClient: httpClient,

		registeredRobotsTXT: new(sync.Map),
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
	robotsTXTLink, err := makeRobotsTXTLink(link)
	if err != nil {
		return nil, errors.Wrap(err, "unable to make the robots.txt link")
	}

	robotsTXTData, ok := register.registeredRobotsTXT.Load(robotsTXTLink)
	if !ok {
		var err error
		robotsTXTData, err =
			loadRobotsTXTData(ctx, register.httpClient, robotsTXTLink)
		if err != nil {
			return nil, errors.Wrap(err, "unable to load the robots.txt data")
		}

		register.registeredRobotsTXT.Store(robotsTXTLink, robotsTXTData)
	}

	return robotsTXTData.(*robotstxt.RobotsData), nil
}

func makeRobotsTXTLink(regularLink string) (robotsTXTLink string, err error) {
	parsedRegularLink, err := url.Parse(regularLink)
	if err != nil {
		return "", errors.Wrap(err, "unable to parse the regular link")
	}

	parsedRobotsTXTLink := &url.URL{
		Scheme: parsedRegularLink.Scheme,
		User:   parsedRegularLink.User,
		Host:   parsedRegularLink.Host,
		Path:   "/robots.txt",
	}
	return parsedRobotsTXTLink.String(), nil
}

func loadRobotsTXTData(
	ctx context.Context,
	httpClient HTTPClient,
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

	response, err := httpClient.Do(request)
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
