package registers

import (
	"context"
	"sync"
	"time"

	"github.com/go-log/log"
	"github.com/pkg/errors"
	"github.com/yterajima/go-sitemap"
)

//go:generate mockery --name=LinkGenerator --inpackage --case=underscore --testonly

// LinkGenerator ...
type LinkGenerator interface {
	GenerateLinks(baseLink string) ([]string, error)
}

// SleepHandler ...
type SleepHandler func(duration time.Duration)

// SitemapRegister ...
type SitemapRegister struct {
	loadingInterval time.Duration
	linkGenerator   LinkGenerator
	logger          log.Logger
	sleeper         SleepHandler

	registeredSitemaps *sync.Map
}

// NewSitemapRegister ...
func NewSitemapRegister(
	loadingInterval time.Duration,
	linkGenerator LinkGenerator,
	logger log.Logger,
	linkLoader func(link string, options interface{}) ([]byte, error),
) SitemapRegister {
	sitemap.SetInterval(loadingInterval)
	if linkLoader != nil {
		sitemap.SetFetch(linkLoader)
	}

	return SitemapRegister{
		linkGenerator: linkGenerator,
		logger:        logger,

		registeredSitemaps: new(sync.Map),
	}
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
