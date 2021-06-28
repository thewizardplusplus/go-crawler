package extractors

import (
	"context"
	"time"

	"github.com/go-log/log"
	"github.com/thewizardplusplus/go-crawler/models"
)

// SleepHandler ...
type SleepHandler func(duration time.Duration)

// RepeatingExtractor ...
type RepeatingExtractor struct {
	LinkExtractor models.LinkExtractor
	RepeatCount   int
	RepeatDelay   time.Duration
	Logger        log.Logger
	SleepHandler  SleepHandler
}

// ExtractLinks ...
func (extractor RepeatingExtractor) ExtractLinks(
	ctx context.Context,
	threadID int,
	link string,
) ([]string, error) {
	var links []string
	for repeat := 0; repeat < extractor.RepeatCount; repeat++ {
		var err error
		links, err = extractor.LinkExtractor.ExtractLinks(ctx, threadID, link)
		if err == nil {
			break
		}
		if repeat == extractor.RepeatCount-1 {
			return nil, err
		}

		const logMessage = "unable to extract links for link %q (repeat #%d): %s"
		extractor.Logger.Logf(logMessage, link, repeat, err)

		extractor.SleepHandler(extractor.RepeatDelay)
	}

	return links, nil
}
