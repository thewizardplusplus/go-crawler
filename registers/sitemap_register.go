package registers

import (
	"context"
	"sync"
	"time"

	"github.com/go-log/log"
	"github.com/pkg/errors"
	"github.com/thewizardplusplus/go-crawler/models"
	"github.com/yterajima/go-sitemap"
)

// SitemapRegister ...
type SitemapRegister struct {
	linkGenerator models.LinkExtractor
	logger        log.Logger

	sitemapRegister BasicRegister
}

// NewSitemapRegister ...
func NewSitemapRegister(
	loadingInterval time.Duration,
	linkGenerator models.LinkExtractor,
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

		sitemapRegister: NewBasicRegister(),
	}
}

// RegisterSitemap ...
func (register SitemapRegister) RegisterSitemap(
	ctx context.Context,
	threadID int,
	link string,
) (
	sitemap.Sitemap,
	error,
) {
	sitemapLinks, err := register.linkGenerator.ExtractLinks(ctx, threadID, link)
	if err != nil {
		return sitemap.Sitemap{}, errors.Wrap(err, "unable to generate Sitemap links")
	}

	var waiter sync.WaitGroup
	waiter.Add(len(sitemapLinks))

	sitemapDataGroup := make([]sitemap.Sitemap, len(sitemapLinks))
	for index, sitemapLink := range sitemapLinks {
		go func(index int, sitemapLink string) {
			defer waiter.Done()

			sitemapData := register.loadSitemapData(ctx, sitemapLink)
			sitemapDataGroup[index] = sitemapData
		}(index, sitemapLink)
	}

	waiter.Wait()

	var totalSitemapData sitemap.Sitemap
	for _, sitemapData := range sitemapDataGroup {
		totalSitemapData.URL = append(totalSitemapData.URL, sitemapData.URL...)
	}

	return totalSitemapData, nil
}

func (register SitemapRegister) loadSitemapData(
	ctx context.Context,
	sitemapLink string,
) sitemap.Sitemap {
	sitemapData, _ := register.sitemapRegister.RegisterValue( // nolint: gosec
		ctx,
		sitemapLink,
		func(ctx context.Context, sitemapLink interface{}) (interface{}, error) {
			sitemapData, err := sitemap.Get(sitemapLink.(string), ctx)
			if err != nil {
				register.logger.Logf("unable to load Sitemap link %q: %s", sitemapLink, err)
			}

			return sitemapData, nil
		},
	)

	return sitemapData.(sitemap.Sitemap)
}
