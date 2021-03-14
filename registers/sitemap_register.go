package registers

import (
	"context"
	"sync"

	"github.com/go-log/log"
	"github.com/pkg/errors"
	"github.com/yterajima/go-sitemap"
)

// LinkGenerator ...
type LinkGenerator interface {
	GenerateLinks(baseLink string) ([]string, error)
}

// SitemapRegister ...
type SitemapRegister struct {
	linkGenerator LinkGenerator
	logger        log.Logger

	registeredSitemaps *sync.Map
}

// RegisterSitemap ...
func (register SitemapRegister) RegisterSitemap(
	ctx context.Context,
	link string,
) (
	sitemap.Sitemap,
	error,
) {
	sitemapLinks, err := register.linkGenerator.GenerateLinks(link)
	if err != nil {
		return sitemap.Sitemap{}, errors.Wrap(err, "unable to generate Sitemap links")
	}

	var totalSitemapData sitemap.Sitemap
	for _, sitemapLink := range sitemapLinks {
		sitemapData, err := register.loadSitemapData(ctx, sitemapLink)
		if err != nil {
			register.logger.
				Logf("unable to process the Sitemap link %q: %s", sitemapLink, err)

			continue
		}

		totalSitemapData.URL = append(totalSitemapData.URL, sitemapData.URL...)
	}

	return totalSitemapData, nil
}

func (register SitemapRegister) loadSitemapData(
	ctx context.Context,
	sitemapLink string,
) (
	sitemap.Sitemap,
	error,
) {
	sitemapData, ok := register.registeredSitemaps.Load(sitemapLink)
	if !ok {
		var err error
		sitemapData, err = sitemap.Get(sitemapLink, ctx)
		if err != nil {
			return sitemap.Sitemap{}, errors.Wrap(err, "unable to load the Sitemap data")
		}

		register.registeredSitemaps.Store(sitemapLink, sitemapData)
	}

	return sitemapData.(sitemap.Sitemap), nil
}
