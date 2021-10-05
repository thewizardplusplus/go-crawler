package transformers

import (
	"bytes"
	"net/http"

	"github.com/pkg/errors"
	urlutils "github.com/thewizardplusplus/go-crawler/url-utils"
	htmlselector "github.com/thewizardplusplus/go-html-selector"
)

// ResolvingTransformer ...
type ResolvingTransformer struct {
	BaseHeaderNames []string
}

// TransformLinks ...
func (transformer ResolvingTransformer) TransformLinks(
	links []string,
	response *http.Response,
	responseContent []byte,
) ([]string, error) {
	baseTag := selectBaseTag(responseContent)
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
			return nil, errors.Wrapf(err, "unable to resolve link %q", link)
		}

		resolvedLinks = append(resolvedLinks, resolvedLink)
	}

	return resolvedLinks, nil
}

func selectBaseTag(data []byte) string {
	var builder BaseTagBuilder
	htmlselector.SelectTags( // nolint: errcheck, gosec
		bytes.NewReader(data),
		BaseTagFilters,
		&builder,
		htmlselector.SkipEmptyTags(),
		htmlselector.SkipEmptyAttributes(),
	)

	baseLink, _ := builder.BaseLink()
	return string(baseLink)
}
