package extractors

import (
	"github.com/go-log/log"
	"github.com/thewizardplusplus/go-crawler/registers"
)

// SitemapExtractor ...
type SitemapExtractor struct {
	SitemapRegister registers.SitemapRegister
	Logger          log.Logger
}
