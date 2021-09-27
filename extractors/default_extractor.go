package extractors

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
	urlutils "github.com/thewizardplusplus/go-crawler/url-utils"
	htmlselector "github.com/thewizardplusplus/go-html-selector"
	"github.com/thewizardplusplus/go-html-selector/builders"
	httputils "github.com/thewizardplusplus/go-http-utils"
)

// DefaultExtractor ...
type DefaultExtractor struct {
	HTTPClient      httputils.HTTPClient
	Filters         htmlselector.OptimizedFilterGroup
	BaseHeaderNames []string
}

// ExtractLinks ...
func (extractor DefaultExtractor) ExtractLinks(
	ctx context.Context,
	threadID int,
	link string,
) ([]string, error) {
	data, response, err := extractor.loadData(ctx, link)
	if err != nil {
		return nil, errors.Wrap(err, "unable to load the data")
	}

	links := extractor.selectLinks(data)
	resolvedLinks, err := extractor.resolveLinks(links, data, response)
	if err != nil {
		return nil, errors.Wrap(err, "unable to resolve the links")
	}

	return resolvedLinks, nil
}

func (extractor DefaultExtractor) loadData(
	ctx context.Context,
	link string,
) ([]byte, *http.Response, error) {
	request, err := http.NewRequest(http.MethodGet, link, nil)
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to create the request")
	}
	request = request.WithContext(ctx)

	response, err := extractor.HTTPClient.Do(request)
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to send the request")
	}
	defer response.Body.Close() // nolint: errcheck

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to read the response")
	}

	return data, response, nil
}

func (extractor DefaultExtractor) selectLinks(data []byte) []string {
	var builder builders.FlattenBuilder
	htmlselector.SelectTags( // nolint: errcheck, gosec
		bytes.NewReader(data),
		extractor.Filters,
		&builder,
		htmlselector.SkipEmptyTags(),
		htmlselector.SkipEmptyAttributes(),
	)

	var links []string
	for _, attributeValue := range builder.AttributeValues() {
		link := string(attributeValue)
		links = append(links, link)
	}

	return links
}

func (extractor DefaultExtractor) resolveLinks(
	links []string,
	data []byte,
	response *http.Response,
) ([]string, error) {
	linkResolver, err := urlutils.NewLinkResolver(urlutils.GenerateBaseLinks(
		response,
		selectBaseTag(data),
		extractor.BaseHeaderNames,
	))
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
