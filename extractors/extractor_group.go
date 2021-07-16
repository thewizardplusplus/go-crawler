package extractors

import (
	"context"
	"sync"

	"github.com/go-log/log"
	"github.com/thewizardplusplus/go-crawler/models"
)

// ExtractorGroup ...
type ExtractorGroup struct {
	LinkExtractors []models.LinkExtractor
	Logger         log.Logger
}

// ExtractLinks ...
func (extractors ExtractorGroup) ExtractLinks(
	ctx context.Context,
	threadID int,
	link string,
) ([]string, error) {
	var waiter sync.WaitGroup
	waiter.Add(len(extractors.LinkExtractors))

	linkGroups := make([][]string, len(extractors.LinkExtractors))
	for index, extractor := range extractors.LinkExtractors {
		go func(index int, extractor models.LinkExtractor) {
			defer waiter.Done()

			links, err := extractor.ExtractLinks(ctx, threadID, link)
			if err != nil {
				const logMessage = "unable to extract links for link %q " +
					"via extractor #%d: %s"
				extractors.Logger.Logf(logMessage, link, index, err)
			}

			linkGroups[index] = links
		}(index, extractor)
	}

	waiter.Wait()

	var totalLinks []string
	for _, linkGroup := range linkGroups {
		totalLinks = append(totalLinks, linkGroup...)
	}

	return totalLinks, nil
}
