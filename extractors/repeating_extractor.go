package extractors

import (
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
