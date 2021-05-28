package extractors

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thewizardplusplus/go-crawler/models"
	"golang.org/x/sync/errgroup"
)

// ExtractorGroup ...
type ExtractorGroup []models.LinkExtractor

// ExtractLinks ...
func (extractors ExtractorGroup) ExtractLinks(
	ctx context.Context,
	threadID int,
	link string,
) ([]string, error) {
	waiter, ctx := errgroup.WithContext(ctx)

	linkGroups := make([][]string, len(extractors))
	for index, extractor := range extractors {
		index, extractor := index, extractor

		waiter.Go(func() error {
			links, err := extractor.ExtractLinks(ctx, threadID, link)
			if err != nil {
				return errors.Wrapf(err, "error with extractor #%d", index)
			}

			linkGroups[index] = links
			return nil
		})
	}

	if err := waiter.Wait(); err != nil {
		return nil, errors.Wrap(err, "unable to extract links")
	}

	var totalLinks []string
	for _, linkGroup := range linkGroups {
		totalLinks = append(totalLinks, linkGroup...)
	}

	return totalLinks, nil
}
