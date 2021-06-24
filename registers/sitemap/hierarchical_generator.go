package sitemap

import (
	"context"

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
	return urlutils.GenerateHierarchicalLinks(
		baseLink,
		"sitemap.xml",
		urlutils.SanitizeBaseLink(generator.SanitizeLink),
		urlutils.WithMaximalHierarchyDepth(generator.MaximalDepth),
	)
}
