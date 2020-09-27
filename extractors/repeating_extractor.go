package extractors

import (
	"context"
	"time"

	"github.com/go-log/log"
	crawler "github.com/thewizardplusplus/go-crawler"
)

// RepeatingExtractor ...
type RepeatingExtractor struct {
	LinkExtractor crawler.LinkExtractor
	RepeatCount   int
	RepeatDelay   time.Duration
	Logger        log.Logger
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

		extractor.Logger.Logf("unable to extract links (repeat #%d): %s", repeat, err)
		time.Sleep(extractor.RepeatDelay)
	}

	return links, nil
}
