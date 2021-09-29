package sitemap

import (
	"context"

	"github.com/pkg/errors"
	urlutils "github.com/thewizardplusplus/go-crawler/url-utils"
)

// HierarchicalGenerator ...
type HierarchicalGenerator struct {
	SanitizeLink urlutils.LinkSanitizing
	MaximalDepth int
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
	hierarchicalLinks, err := urlutils.GenerateHierarchicalLinks(
		baseLink,
		"sitemap.xml",
		urlutils.SanitizeBaseLink(generator.SanitizeLink),
		urlutils.WithMaximalHierarchyDepth(generator.MaximalDepth),
	)
	if err != nil {
		return nil, errors.Wrap(err, "unable to generate the hierarchical links")
	}

	return hierarchicalLinks, nil
}
