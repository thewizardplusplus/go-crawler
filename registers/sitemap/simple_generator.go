package sitemap

import (
	"context"
	"net/url"

	"github.com/pkg/errors"
)

// SimpleGenerator ...
type SimpleGenerator struct{}

// ExtractLinks ...
func (generator SimpleGenerator) ExtractLinks(
	ctx context.Context,
	threadID int,
	baseLink string,
) (
	[]string,
	error,
) {
	parsedBaseLink, err := url.Parse(baseLink)
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse the base link")
	}

	parsedSitemapLink := &url.URL{
		Scheme: parsedBaseLink.Scheme,
		User:   parsedBaseLink.User,
		Host:   parsedBaseLink.Host,
		Path:   "/sitemap.xml",
	}
	return []string{parsedSitemapLink.String()}, nil
}
