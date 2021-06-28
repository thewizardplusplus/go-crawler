package extractors

import (
	"context"
	"testing"
	"testing/iotest"
	"time"

	"github.com/go-log/log"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/go-crawler/models"
	"github.com/thewizardplusplus/go-crawler/registers"
)

func TestSitemapExtractor_ExtractLinks(test *testing.T) {
	type fields struct {
		loadingInterval time.Duration
		linkGenerator   models.LinkExtractor
		logger          log.Logger
		linkLoader      LinkLoader
	}
	type args struct {
		ctx      context.Context
		threadID int
		link     string
	}

	for _, data := range []struct {
		name      string
		fields    fields
		args      args
		wantLinks []string
		wantErr   assert.ErrorAssertionFunc
	}{
		{
			name: "success without sitemap.xml links",
			fields: fields{
				loadingInterval: 5 * time.Second,
				linkGenerator: func() models.LinkExtractor {
					linkGenerator := new(MockLinkExtractor)
					linkGenerator.
						On("ExtractLinks", context.Background(), 23, "http://example.com/").
						Return(nil, nil)

					return linkGenerator
				}(),
				logger:     new(MockLogger),
				linkLoader: new(MockLinkLoader),
			},
			args: args{
				ctx:      context.Background(),
				threadID: 23,
				link:     "http://example.com/",
			},
			wantLinks: nil,
			wantErr:   assert.NoError,
		},
		{
			name: "success without links",
			fields: fields{
				loadingInterval: 5 * time.Second,
				linkGenerator: func() models.LinkExtractor {
					sitemapLinks := []string{
						"http://example.com/sitemap_1.xml",
						"http://example.com/sitemap_2.xml",
					}

					linkGenerator := new(MockLinkExtractor)
					linkGenerator.
						On("ExtractLinks", context.Background(), 23, "http://example.com/").
						Return(sitemapLinks, nil)

					return linkGenerator
				}(),
				logger: new(MockLogger),
				linkLoader: func() LinkLoader {
					const responseOne = `
						<?xml version="1.0" encoding="UTF-8" ?>
						<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
						</urlset>
					`
					const responseTwo = `
						<?xml version="1.0" encoding="UTF-8" ?>
						<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
						</urlset>
					`

					linkLoader := new(MockLinkLoader)
					linkLoader.
						On("LoadLink", "http://example.com/sitemap_1.xml", context.Background()).
						Return([]byte(responseOne), nil)
					linkLoader.
						On("LoadLink", "http://example.com/sitemap_2.xml", context.Background()).
						Return([]byte(responseTwo), nil)

					return linkLoader
				}(),
			},
			args: args{
				ctx:      context.Background(),
				threadID: 23,
				link:     "http://example.com/",
			},
			wantLinks: nil,
			wantErr:   assert.NoError,
		},
		{
			name: "success with links",
			fields: fields{
				loadingInterval: 5 * time.Second,
				linkGenerator: func() models.LinkExtractor {
					sitemapLinks := []string{
						"http://example.com/sitemap_1.xml",
						"http://example.com/sitemap_2.xml",
					}

					linkGenerator := new(MockLinkExtractor)
					linkGenerator.
						On("ExtractLinks", context.Background(), 23, "http://example.com/").
						Return(sitemapLinks, nil)

					return linkGenerator
				}(),
				logger: new(MockLogger),
				linkLoader: func() LinkLoader {
					const responseOne = `
						<?xml version="1.0" encoding="UTF-8" ?>
						<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
							<url>
								<loc>http://example.com/1</loc>
							</url>
							<url>
								<loc>http://example.com/2</loc>
							</url>
						</urlset>
					`
					const responseTwo = `
						<?xml version="1.0" encoding="UTF-8" ?>
						<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
							<url>
								<loc>http://example.com/3</loc>
							</url>
							<url>
								<loc>http://example.com/4</loc>
							</url>
						</urlset>
					`

					linkLoader := new(MockLinkLoader)
					linkLoader.
						On("LoadLink", "http://example.com/sitemap_1.xml", context.Background()).
						Return([]byte(responseOne), nil)
					linkLoader.
						On("LoadLink", "http://example.com/sitemap_2.xml", context.Background()).
						Return([]byte(responseTwo), nil)

					return linkLoader
				}(),
			},
			args: args{
				ctx:      context.Background(),
				threadID: 23,
				link:     "http://example.com/",
			},
			wantLinks: []string{
				"http://example.com/1",
				"http://example.com/2",
				"http://example.com/3",
				"http://example.com/4",
			},
			wantErr: assert.NoError,
		},
		{
			name: "error with generation",
			fields: fields{
				loadingInterval: 5 * time.Second,
				linkGenerator: func() models.LinkExtractor {
					linkGenerator := new(MockLinkExtractor)
					linkGenerator.
						On("ExtractLinks", context.Background(), 23, "http://example.com/").
						Return(nil, iotest.ErrTimeout)

					return linkGenerator
				}(),
				logger: func() Logger {
					wantErr :=
						errors.Wrap(iotest.ErrTimeout, "unable to generate Sitemap links")

					logger := new(MockLogger)
					logger.On(
						"Logf",
						"unable to register the sitemap.xml link for link %q: %s",
						"http://example.com/",
						mock.MatchedBy(func(gotErr error) bool {
							return gotErr.Error() == wantErr.Error()
						}),
					).Return()

					return logger
				}(),
				linkLoader: new(MockLinkLoader),
			},
			args: args{
				ctx:      context.Background(),
				threadID: 23,
				link:     "http://example.com/",
			},
			wantLinks: nil,
			wantErr:   assert.NoError,
		},
		{
			name: "error with loading",
			fields: fields{
				loadingInterval: 5 * time.Second,
				linkGenerator: func() models.LinkExtractor {
					sitemapLinks := []string{
						"http://example.com/sitemap_1.xml",
						"http://example.com/sitemap_2.xml",
					}

					linkGenerator := new(MockLinkExtractor)
					linkGenerator.
						On("ExtractLinks", context.Background(), 23, "http://example.com/").
						Return(sitemapLinks, nil)

					return linkGenerator
				}(),
				logger: func() Logger {
					wantErr := iotest.ErrTimeout

					logger := new(MockLogger)
					logger.On(
						"Logf",
						"unable to load Sitemap link %q: %s",
						"http://example.com/sitemap_1.xml",
						mock.MatchedBy(func(gotErr error) bool {
							return gotErr.Error() == wantErr.Error()
						}),
					).Return()

					return logger
				}(),
				linkLoader: func() LinkLoader {
					const responseTwo = `
						<?xml version="1.0" encoding="UTF-8" ?>
						<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
							<url>
								<loc>http://example.com/3</loc>
							</url>
							<url>
								<loc>http://example.com/4</loc>
							</url>
						</urlset>
					`

					linkLoader := new(MockLinkLoader)
					linkLoader.
						On("LoadLink", "http://example.com/sitemap_1.xml", context.Background()).
						Return(nil, iotest.ErrTimeout)
					linkLoader.
						On("LoadLink", "http://example.com/sitemap_2.xml", context.Background()).
						Return([]byte(responseTwo), nil)

					return linkLoader
				}(),
			},
			args: args{
				ctx:      context.Background(),
				threadID: 23,
				link:     "http://example.com/",
			},
			wantLinks: []string{
				"http://example.com/3",
				"http://example.com/4",
			},
			wantErr: assert.NoError,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			register := registers.NewSitemapRegister(
				data.fields.loadingInterval,
				data.fields.linkGenerator,
				data.fields.logger,
				data.fields.linkLoader.LoadLink,
			)
			extractor := SitemapExtractor{
				SitemapRegister: register,
				Logger:          data.fields.logger,
			}
			gotLinks, gotErr := extractor.ExtractLinks(
				data.args.ctx,
				data.args.threadID,
				data.args.link,
			)

			mock.AssertExpectationsForObjects(
				test,
				data.fields.linkGenerator,
				data.fields.logger,
				data.fields.linkLoader,
			)
			assert.ElementsMatch(test, data.wantLinks, gotLinks)
			data.wantErr(test, gotErr)
		})
	}
}
