package urlutils

import (
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

// DefaultBaseHeaderNames ...
var DefaultBaseHeaderNames = []string{"Content-Base", "Content-Location"} // nolint: gochecknoglobals, lll

// LinkResolver ...
type LinkResolver struct {
	BaseLink *url.URL
}

// GenerateBaseLinks ...
func GenerateBaseLinks(
	response *http.Response,
	baseTagValue string,
	baseHeaderNames []string,
) []string {
	var baseLinks []string
	if baseTagValue != "" {
		baseLinks = append(baseLinks, baseTagValue)
	}

	for _, baseHeaderName := range baseHeaderNames {
		baseHeader := response.Header.Get(baseHeaderName)
		if baseHeader != "" {
			baseLinks = append(baseLinks, baseHeader)
		}
	}

	requestURI := response.Request.URL.String()
	return append(baseLinks, requestURI)
}

// NewLinkResolver ...
func NewLinkResolver(baseLinks []string) (LinkResolver, error) {
	var parsedBaseLinks []*url.URL
	for _, baseLink := range baseLinks {
		parsedBaseLink, err := url.Parse(baseLink)
		if err != nil {
			return LinkResolver{},
				errors.Wrapf(err, "unable to parse base link %q", baseLink)
		}

		parsedBaseLinks = append(parsedBaseLinks, parsedBaseLink)
		if parsedBaseLink.IsAbs() {
			break
		}
	}

	var totalBaseLink *url.URL
	for index := len(parsedBaseLinks) - 1; index >= 0; index-- {
		parsedBaseLink := parsedBaseLinks[index]
		if totalBaseLink != nil {
			totalBaseLink = totalBaseLink.ResolveReference(parsedBaseLink)
		} else {
			totalBaseLink = parsedBaseLink
		}
	}

	resolver := LinkResolver{BaseLink: totalBaseLink}
	return resolver, nil
}

// ResolveLink ...
func (resolver LinkResolver) ResolveLink(link string) (string, error) {
	parsedLink, err := url.Parse(link)
	if err != nil {
		return "", errors.Wrap(err, "unable to parse the link")
	}

	resolvedLink := resolver.BaseLink.ResolveReference(parsedLink)
	return resolvedLink.String(), nil
}
