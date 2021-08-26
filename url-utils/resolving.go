package urlutils

import (
	"net/url"

	"github.com/pkg/errors"
)

// LinkResolver ...
type LinkResolver struct {
	BaseLink *url.URL
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
