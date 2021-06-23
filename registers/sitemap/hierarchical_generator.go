package sitemap

import (
	"context"
	"net/url"
	"strings"

	"github.com/pkg/errors"
	urlutils "github.com/thewizardplusplus/go-crawler/url-utils"
)

// HierarchicalGenerator ...
type HierarchicalGenerator struct {
	SanitizeLink urlutils.LinkSanitizing
}

// ExtractLinks ...
func (generator HierarchicalGenerator) ExtractLinks(
	ctx context.Context,
	threadID int,
	baseLink string,
) (
	[]string,
	error,
) {
	if generator.SanitizeLink == urlutils.SanitizeLink {
		var err error
		baseLink, err = urlutils.ApplyLinkSanitizing(baseLink)
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
