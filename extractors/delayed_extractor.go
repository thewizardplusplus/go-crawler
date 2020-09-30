package extractors

import (
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
