package registers

import (
	"net/url"

	"github.com/pkg/errors"
)

func makeRobotsTXTLink(regularLink string) (robotsTXTLink *url.URL, err error) {
	parsedLink, err := url.Parse(regularLink)
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse the regular link")
	}

	robotsTXTLink = &url.URL{
		Scheme: parsedLink.Scheme,
		User:   parsedLink.User,
		Host:   parsedLink.Host,
		Path:   "/robots.txt",
	}
	return robotsTXTLink, nil
}
