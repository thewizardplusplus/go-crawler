package registers

import (
	"net/url"

	"github.com/pkg/errors"
)

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
