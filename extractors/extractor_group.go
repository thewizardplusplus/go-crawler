package extractors

import (
	"context"
	"fmt"
	"sync"

	"github.com/go-log/log"
	"github.com/thewizardplusplus/go-crawler/models"
)

// ExtractorGroup ...
type ExtractorGroup struct {
	Name           string
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

	var logPrefix string
	if extractors.Name != "" {
		logPrefix = fmt.Sprintf("%s: ", extractors.Name)
	}

	linkGroups := make([][]string, len(extractors.LinkExtractors))
	for index, extractor := range extractors.LinkExtractors {
		go func(index int, extractor models.LinkExtractor) {
			defer waiter.Done()

			links, err := extractor.ExtractLinks(ctx, threadID, link)
			if err != nil {
				const logMessage = "%sunable to extract links for link %q " +
					"via extractor #%d: %s"
				extractors.Logger.Logf(logMessage, logPrefix, link, index, err)

				return
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
