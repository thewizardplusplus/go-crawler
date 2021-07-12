package extractors

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thewizardplusplus/go-crawler/models"
	urlutils "github.com/thewizardplusplus/go-crawler/url-utils"
)

// TrimmingExtractor ...
type TrimmingExtractor struct {
	TrimLink      urlutils.LinkTrimming
	LinkExtractor models.LinkExtractor
}

// ExtractLinks ...
func (extractor TrimmingExtractor) ExtractLinks(
	ctx context.Context,
	threadID int,
	link string,
) ([]string, error) {
	links, err := extractor.LinkExtractor.ExtractLinks(ctx, threadID, link)
	if err != nil {
		return nil, errors.Wrap(err, "unable to extract links")
	}

	var trimmedLinks []string
	for _, link := range links {
		trimmedLink := urlutils.ApplyLinkTrimming(link, extractor.TrimLink)
		trimmedLinks = append(trimmedLinks, trimmedLink)
	}

	return trimmedLinks, nil
}
