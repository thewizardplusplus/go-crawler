package transformers

import (
	"bytes"
	"net/http"

	"github.com/go-log/log"
	"github.com/pkg/errors"
	urlutils "github.com/thewizardplusplus/go-crawler/url-utils"
	htmlselector "github.com/thewizardplusplus/go-html-selector"
)

// ResolvingTransformer ...
type ResolvingTransformer struct {
	BaseTagFilters  htmlselector.OptimizedFilterGroup
	BaseHeaderNames []string
	Logger          log.Logger
}

// TransformLinks ...
func (transformer ResolvingTransformer) TransformLinks(
	links []string,
	response *http.Response,
	responseContent []byte,
) ([]string, error) {
	baseTag := transformer.selectBaseTag(responseContent)
	baseLinks :=
		urlutils.GenerateBaseLinks(response, baseTag, transformer.BaseHeaderNames)
	linkResolver, err := urlutils.NewLinkResolver(baseLinks)
	if err != nil {
		return nil, errors.Wrap(err, "unable to construct the link resolver")
	}

	var resolvedLinks []string
	for _, link := range links {
		resolvedLink, err := linkResolver.ResolveLink(link)
		if err != nil {
			transformer.Logger.Logf("unable to resolve link %q: %s", link, err)
			continue
		}

		resolvedLinks = append(resolvedLinks, resolvedLink)
	}

	return resolvedLinks, nil
}

func (transformer ResolvingTransformer) selectBaseTag(data []byte) string {
	var builder BaseTagBuilder
	htmlselector.SelectTags( // nolint: errcheck, gosec
		bytes.NewReader(data),
		transformer.BaseTagFilters,
		&builder,
		htmlselector.SkipEmptyTags(),
		htmlselector.SkipEmptyAttributes(),
	)

	baseLink, _ := builder.BaseLink()
	return string(baseLink)
}
