package registers

import (
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

// HTTPClient ...
type HTTPClient interface {
	Do(request *http.Request) (*http.Response, error)
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
