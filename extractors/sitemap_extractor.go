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
	sitemapData, err :=
		extractor.SitemapRegister.RegisterSitemap(ctx, threadID, link)
	if err != nil {
		const logMessage = "unable to register the sitemap.xml link for link %q: %s"
		extractor.Logger.Logf(logMessage, link, err)

		return nil, nil
	}

	var links []string
	for _, url := range sitemapData.URL {
		links = append(links, url.Loc)
	}

	return links, nil
}
