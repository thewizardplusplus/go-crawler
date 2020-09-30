package extractors

import (
	"context"
	"sync"
	"time"

	crawler "github.com/thewizardplusplus/go-crawler"
)

// DelayedExtractor ...
type DelayedExtractor struct {
	timestamps    sync.Map // map[threadID]time.Time
	minimalDelay  time.Duration
	sleeper       SleepHandler
	linkExtractor crawler.LinkExtractor
}

// NewDelayedExtractor ...
func NewDelayedExtractor(
	minimalDelay time.Duration,
	sleeper SleepHandler,
	linkExtractor crawler.LinkExtractor,
) *DelayedExtractor {
	return &DelayedExtractor{
		minimalDelay:  minimalDelay,
		sleeper:       sleeper,
		linkExtractor: linkExtractor,
	}
}

// ExtractLinks ...
func (extractor *DelayedExtractor) ExtractLinks(
	ctx context.Context,
	threadID int,
	link string,
) ([]string, error) {
	// function wrapping is necessary to correctly compute an extraction time
	defer func() { extractor.timestamps.Store(threadID, time.Now()) }()

	if lastExtractionTime, ok := extractor.timestamps.Load(threadID); ok {
		expiredTime := time.Since(lastExtractionTime.(time.Time))
		extractor.sleeper(extractor.minimalDelay - expiredTime)
	}

	return extractor.linkExtractor.ExtractLinks(ctx, threadID, link)
}
