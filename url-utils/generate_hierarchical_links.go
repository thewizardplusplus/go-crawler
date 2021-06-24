package urlutils

import (
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

// GenerateHierarchicalLinks ...
func GenerateHierarchicalLinks(
	baseLink string,
	linkSuffix string,
	options ...HierarchicalLinkOption,
) ([]string, error) {
	// default config
	config := HierarchicalLinkConfig{
		sanitizeBaseLink:      DoNotSanitizeLink,
		maximalHierarchyDepth: -1,
	}
	for _, option := range options {
		option(&config)
	}

	if config.sanitizeBaseLink == SanitizeLink {
		var err error
		baseLink, err = ApplyLinkSanitizing(baseLink)
		if err != nil {
			return nil, errors.Wrap(err, "unable to sanitize the base link")
		}
	}

	parsedBaseLink, err := url.Parse(baseLink)
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse the base link")
	}
	if parsedBaseLink.Path != "" && !strings.HasPrefix(parsedBaseLink.Path, "/") {
		return nil, errors.New("link path is not absolute")
	}
	if parsedBaseLink.Path == "" {
		parsedBaseLink.Path = "/"
	}

	pathParts := strings.Split(parsedBaseLink.Path, "/")
	// increase due to the additional empty part at the beginning
	// (since the link path is absolute)
	fixedMaximalDepth := config.maximalHierarchyDepth + 1
	if config.maximalHierarchyDepth == -1 || fixedMaximalDepth > len(pathParts)-1 {
		fixedMaximalDepth = len(pathParts) - 1
	}
	pathParts = pathParts[:fixedMaximalDepth]

	var hierarchicalLinks []string
	var pathPrefix string
	for index, pathPart := range pathParts {
		// first part is always empty, so do not append it to avoid redundant slashes
		if index != 0 {
			pathPrefix += "/" + pathPart
		}

		parsedHierarchicalLink := &url.URL{
			Scheme: parsedBaseLink.Scheme,
			User:   parsedBaseLink.User,
			Host:   parsedBaseLink.Host,
			Path:   pathPrefix + "/" + linkSuffix,
		}
		hierarchicalLink := parsedHierarchicalLink.String()
		hierarchicalLinks = append(hierarchicalLinks, hierarchicalLink)
	}

	return hierarchicalLinks, nil
}
