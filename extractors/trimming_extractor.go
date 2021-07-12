package extractors

import (
	"github.com/thewizardplusplus/go-crawler/models"
	urlutils "github.com/thewizardplusplus/go-crawler/url-utils"
)

// TrimmingExtractor ...
type TrimmingExtractor struct {
	TrimLink      urlutils.LinkTrimming
	LinkExtractor models.LinkExtractor
}
