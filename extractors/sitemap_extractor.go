package extractors

import (
	"context"

	"github.com/go-log/log"
	"github.com/thewizardplusplus/go-crawler/registers"
)

// SitemapExtractor ...
type SitemapExtractor struct {
	SitemapRegister registers.SitemapRegister
	Logger          log.Logger
}

// ExtractLinks ...
func (extractor SitemapExtractor) ExtractLinks(
	ctx context.Context,
	threadID int,
	link string,
) ([]string, error) {
	sitemapData, err := extractor.SitemapRegister.RegisterSitemap(ctx, link)
	if err != nil {
		extractor.Logger.Logf("unable to register the sitemap.xml link: %s", err)
		return nil, nil
	}

	var links []string
	for _, url := range sitemapData.URL {
		links = append(links, url.Loc)
	}

	return links, nil
}
