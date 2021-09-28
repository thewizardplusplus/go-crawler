package transformers

import (
	urlutils "github.com/thewizardplusplus/go-crawler/url-utils"
)

// TrimmingTransformer ...
type TrimmingTransformer struct {
	TrimLink urlutils.LinkTrimming
}
