package extractors

import (
	"context"

	"github.com/thewizardplusplus/go-crawler/extractors/transformers"
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
		return nil, err
	}

	trimmingTransformer := transformers.TrimmingTransformer{
		TrimLink: extractor.TrimLink,
	}
	// this method never returns an error
	trimmedLinks, _ := trimmingTransformer.TransformLinks(links, nil, nil)
	return trimmedLinks, nil
}
