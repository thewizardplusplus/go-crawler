package registers

import (
	"context"
	"sync"

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

	registeredSitemaps *sync.Map
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
