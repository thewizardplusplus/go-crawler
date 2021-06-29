package extractors

import (
	"context"
	"net/http"

	"github.com/pkg/errors"
	urlutils "github.com/thewizardplusplus/go-crawler/url-utils"
	htmlselector "github.com/thewizardplusplus/go-html-selector"
	"github.com/thewizardplusplus/go-html-selector/builders"
	httputils "github.com/thewizardplusplus/go-http-utils"
)

// DefaultExtractor ...
type DefaultExtractor struct {
	TrimLink   urlutils.LinkTrimming
	HTTPClient httputils.HTTPClient
	Filters    htmlselector.OptimizedFilterGroup
}

// ExtractLinks ...
func (extractor DefaultExtractor) ExtractLinks(
	ctx context.Context,
	threadID int,
	link string,
) ([]string, error) {
	request, err := http.NewRequest(http.MethodGet, link, nil)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create the request")
	}
	request = request.WithContext(ctx)

	response, err := extractor.HTTPClient.Do(request)
	if err != nil {
		return nil, errors.Wrap(err, "unable to send the request")
	}
	defer response.Body.Close() // nolint: errcheck

	var builder builders.FlattenBuilder
	if err := htmlselector.SelectTags(
		response.Body,
		extractor.Filters,
		&builder,
		htmlselector.SkipEmptyTags(),
		htmlselector.SkipEmptyAttributes(),
	); err != nil {
		return nil, errors.Wrap(err, "unable to select tags")
	}

	var links []string
	for _, attributeValue := range builder.AttributeValues() {
		link := string(attributeValue)
		link = urlutils.ApplyLinkTrimming(link, extractor.TrimLink)

		links = append(links, link)
	}

	return links, nil
}
