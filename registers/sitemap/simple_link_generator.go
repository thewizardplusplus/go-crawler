package sitemap

import (
	"net/url"

	"github.com/pkg/errors"
)

// SimpleLinkGenerator ...
type SimpleLinkGenerator struct{}

// GenerateLinks ...
func (generator SimpleLinkGenerator) GenerateLinks(baseLink string) (
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
