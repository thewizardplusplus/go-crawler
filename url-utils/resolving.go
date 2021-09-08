package urlutils

import (
	"net/url"

	"github.com/pkg/errors"
)

// LinkResolver ...
type LinkResolver struct {
	BaseLink *url.URL
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
