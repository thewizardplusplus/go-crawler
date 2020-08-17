package extractors

import (
	"context"
	"net/http"

	htmlselector "github.com/thewizardplusplus/go-html-selector"
	"github.com/thewizardplusplus/go-html-selector/builders"
)

//go:generate mockery -name=HTTPClient -inpkg -case=underscore -testonly

// HTTPClient ...
type HTTPClient interface {
	Do(request *http.Request) (*http.Response, error)
}

// DefaultExtractor ...
type DefaultExtractor struct {
	HTTPClient HTTPClient
	Filters    htmlselector.OptimizedFilterGroup
}

// ExtractLinks ...
func (extractor DefaultExtractor) ExtractLinks(ctx context.Context, link string) ([]string, error) {
	request, err := http.NewRequest(http.MethodGet, link, nil)
	if err != nil {
		return nil, err
	}
	request = request.WithContext(ctx)

	response, err := extractor.HTTPClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var builder builders.FlattenBuilder
	if err := htmlselector.SelectTags(
		response.Body,
		extractor.Filters,
		&builder,
		htmlselector.SkipEmptyTags(),
		htmlselector.SkipEmptyAttributes(),
	); err != nil {
		return nil, err
	}

	var links []string
	for _, attributeValue := range builder.AttributeValues() {
		link := string(attributeValue)
		links = append(links, link)
	}

	return links, nil
}
