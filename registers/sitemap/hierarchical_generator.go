package sitemap

import (
	"context"
	"net/url"
	"strings"

	"github.com/pkg/errors"
	"github.com/thewizardplusplus/go-crawler/sanitizing"
)

// HierarchicalGenerator ...
type HierarchicalGenerator struct {
	SanitizeLink sanitizing.LinkSanitizing
}

// GenerateLinks ...
func (generator HierarchicalGenerator) GenerateLinks(
	ctx context.Context,
	baseLink string,
) (
	[]string,
	error,
) {
	if generator.SanitizeLink == sanitizing.SanitizeLink {
		var err error
		baseLink, err = sanitizing.ApplyLinkSanitizing(baseLink)
		if err != nil {
			return nil, errors.Wrap(err, "unable to sanitize the base link")
		}
	}

	parsedBaseLink, err := url.Parse(baseLink)
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse the base link")
	}

	var sitemapLinks []string
	var pathPrefix string
	pathParts := strings.Split(parsedBaseLink.Path, "/")
	for _, pathPart := range pathParts[:len(pathParts)-1] {
		if pathPart != "" {
			pathPrefix += "/" + pathPart
		}

		parsedSitemapLink := &url.URL{
			Scheme: parsedBaseLink.Scheme,
			User:   parsedBaseLink.User,
			Host:   parsedBaseLink.Host,
			Path:   pathPrefix + "/sitemap.xml",
		}
		sitemapLinks = append(sitemapLinks, parsedSitemapLink.String())
	}

	return sitemapLinks, nil
}
