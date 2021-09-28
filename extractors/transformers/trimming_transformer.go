package transformers

import (
	"net/http"

	urlutils "github.com/thewizardplusplus/go-crawler/url-utils"
)

// TrimmingTransformer ...
type TrimmingTransformer struct {
	TrimLink urlutils.LinkTrimming
}

// TransformLinks ...
func (transformer TrimmingTransformer) TransformLinks(
	links []string,
	response *http.Response,
	responseContent []byte,
) ([]string, error) {
	var trimmedLinks []string
	for _, link := range links {
		trimmedLink := urlutils.ApplyLinkTrimming(link, transformer.TrimLink)
		trimmedLinks = append(trimmedLinks, trimmedLink)
	}

	return trimmedLinks, nil
}
