package extractors

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thewizardplusplus/go-crawler/models"
)

// ExtractorGroup ...
type ExtractorGroup []models.LinkExtractor

// ExtractLinks ...
func (extractors ExtractorGroup) ExtractLinks(
	ctx context.Context,
	threadID int,
	link string,
) ([]string, error) {
	var totalLinks []string
	for _, extractor := range extractors {
		links, err := extractor.ExtractLinks(ctx, threadID, link)
		if err != nil {
			return nil, errors.Wrap(err, "unable to extract links")
		}

		totalLinks = append(totalLinks, links...)
	}

	return totalLinks, nil
}
