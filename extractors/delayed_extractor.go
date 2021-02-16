package extractors

import (
	"context"
	"sync"
	"time"

	"github.com/thewizardplusplus/go-crawler/models"
)

// DelayingExtractor ...
type DelayingExtractor struct {
	timestamps    sync.Map // map[threadID]time.Time
	minimalDelay  time.Duration
	sleeper       SleepHandler
	linkExtractor models.LinkExtractor
}

// NewDelayingExtractor ...
func NewDelayingExtractor(
	minimalDelay time.Duration,
	sleeper SleepHandler,
	linkExtractor models.LinkExtractor,
) *DelayingExtractor {
	return &DelayingExtractor{
		minimalDelay:  minimalDelay,
		sleeper:       sleeper,
		linkExtractor: linkExtractor,
	}
}

// ExtractLinks ...
func (extractor *DelayingExtractor) ExtractLinks(
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
