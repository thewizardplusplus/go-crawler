package registers

import (
	"context"
	"encoding/xml"
	"reflect"
	"sync"
	"testing"
	"testing/iotest"
	"time"

	"github.com/go-log/log"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/go-crawler/models"
	"github.com/yterajima/go-sitemap"
)

func TestNewSitemapRegister(test *testing.T) {
	type args struct {
		loadingInterval time.Duration
		linkGenerator   models.LinkExtractor
		logger          log.Logger
		linkLoader      LinkLoader
	}

	for _, data := range []struct {
		name                   string
		args                   args
		wantLinkGenerator      models.LinkExtractor
		wantLogger             log.Logger
		wantRegisteredSitemaps *sync.Map
	}{
		{
			name: "with a link loader",
			args: args{
				loadingInterval: 5 * time.Second,
				linkGenerator:   new(MockLinkExtractor),
				logger:          new(MockLogger),
				linkLoader:      new(MockLinkLoader),
			},
			wantLinkGenerator:      new(MockLinkExtractor),
			wantLogger:             new(MockLogger),
			wantRegisteredSitemaps: new(sync.Map),
		},
		{
			name: "without a link loader",
			args: args{
				loadingInterval: 5 * time.Second,
				linkGenerator:   new(MockLinkExtractor),
				logger:          new(MockLogger),
				linkLoader:      nil,
			},
			wantLinkGenerator:      new(MockLinkExtractor),
			wantLogger:             new(MockLogger),
			wantRegisteredSitemaps: new(sync.Map),
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			var linkLoader func(link string, options interface{}) ([]byte, error)
			if data.args.linkLoader != nil {
				linkLoader = data.args.linkLoader.LoadLink
			}
			register := NewSitemapRegister(
				data.args.loadingInterval,
				data.args.linkGenerator,
				data.args.logger,
				linkLoader,
			)

			mock.AssertExpectationsForObjects(
				test,
				data.args.linkGenerator,
				data.args.logger,
			)
			if data.args.linkLoader != nil {
				mock.AssertExpectationsForObjects(test, data.args.linkLoader)
			}
			assert.Equal(test, data.wantLinkGenerator, register.linkGenerator)
			assert.Equal(test, data.wantLogger, register.logger)
			assert.Equal(test, data.wantRegisteredSitemaps, register.registeredSitemaps)
		})
	}
}

func TestSitemapRegister_RegisterSitemap(test *testing.T) {
	type fields struct {
		linkGenerator      models.LinkExtractor
		linkLoader         LinkLoader
		registeredSitemaps *sync.Map
	}
	type args struct {
		ctx  context.Context
		link string
	}

	for _, data := range []struct {
		name            string
		fields          fields
		args            args
		wantSitemapData sitemap.Sitemap
		wantErr         assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			fields: fields{
				linkGenerator: func() models.LinkExtractor {
					sitemapLinks := []string{
						"http://example.com/sitemap_1.xml",
						"http://example.com/sitemap_2.xml",
					}

					linkGenerator := new(MockLinkExtractor)
					linkGenerator.
						On("ExtractLinks", context.Background(), -1, "http://example.com/test").
						Return(sitemapLinks, nil)

					return linkGenerator
				}(),
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
				registeredSitemaps: new(sync.Map),
			},
			args: args{
				ctx:  context.Background(),
				link: "http://example.com/test",
			},
			wantSitemapData: sitemap.Sitemap{
				URL: []sitemap.URL{
					{Loc: "http://example.com/1"},
					{Loc: "http://example.com/2"},
					{Loc: "http://example.com/3"},
					{Loc: "http://example.com/4"},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "error",
			fields: fields{
				linkGenerator: func() models.LinkExtractor {
					linkGenerator := new(MockLinkExtractor)
					linkGenerator.
						On("ExtractLinks", context.Background(), -1, "http://example.com/test").
						Return(nil, iotest.ErrTimeout)

					return linkGenerator
				}(),
				linkLoader:         new(MockLinkLoader),
				registeredSitemaps: new(sync.Map),
			},
			args: args{
				ctx:  context.Background(),
				link: "http://example.com/test",
			},
			wantSitemapData: sitemap.Sitemap{},
			wantErr:         assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			sitemap.SetFetch(data.fields.linkLoader.LoadLink)

			register := SitemapRegister{
				linkGenerator:      data.fields.linkGenerator,
				registeredSitemaps: data.fields.registeredSitemaps,
			}
			gotSitemapData, gotErr :=
				register.RegisterSitemap(data.args.ctx, data.args.link)

			mock.AssertExpectationsForObjects(
				test,
				data.fields.linkGenerator,
				data.fields.linkLoader,
			)
			assert.Equal(test, data.wantSitemapData.XMLName, gotSitemapData.XMLName)
			assert.ElementsMatch(test, data.wantSitemapData.URL, gotSitemapData.URL)
			data.wantErr(test, gotErr)
		})
	}
}

func TestSitemapRegister_loadSitemapData(test *testing.T) {
	type fields struct {
		linkLoader         LinkLoader
		logger             log.Logger
		registeredSitemaps *sync.Map
	}
	type args struct {
		ctx         context.Context
		sitemapLink string
	}

	for _, data := range []struct {
		name            string
		fields          fields
		args            args
		wantSitemapData sitemap.Sitemap
	}{
		{
			name: "success with an unregistered Sitemap link",
			fields: fields{
				linkLoader: func() LinkLoader {
					const response = `
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

					linkLoader := new(MockLinkLoader)
					linkLoader.
						On("LoadLink", "http://example.com/sitemap.xml", context.Background()).
						Return([]byte(response), nil)

					return linkLoader
				}(),
				logger:             new(MockLogger),
				registeredSitemaps: new(sync.Map),
			},
			args: args{
				ctx:         context.Background(),
				sitemapLink: "http://example.com/sitemap.xml",
			},
			wantSitemapData: sitemap.Sitemap{
				XMLName: xml.Name{
					Space: "http://www.sitemaps.org/schemas/sitemap/0.9",
					Local: "urlset",
				},
				URL: []sitemap.URL{
					{Loc: "http://example.com/1"},
					{Loc: "http://example.com/2"},
				},
			},
		},
		{
			name: "success with a registered Sitemap link",
			fields: fields{
				linkLoader: new(MockLinkLoader),
				logger:     new(MockLogger),
				registeredSitemaps: func() *sync.Map {
					sitemapData := sitemap.Sitemap{
						XMLName: xml.Name{
							Space: "http://www.sitemaps.org/schemas/sitemap/0.9",
							Local: "urlset",
						},
						URL: []sitemap.URL{
							{Loc: "http://example.com/1"},
							{Loc: "http://example.com/2"},
						},
					}

					registeredSitemaps := new(sync.Map)
					registeredSitemaps.Store("http://example.com/sitemap.xml", sitemapData)

					return registeredSitemaps
				}(),
			},
			args: args{
				ctx:         context.Background(),
				sitemapLink: "http://example.com/sitemap.xml",
			},
			wantSitemapData: sitemap.Sitemap{
				XMLName: xml.Name{
					Space: "http://www.sitemaps.org/schemas/sitemap/0.9",
					Local: "urlset",
				},
				URL: []sitemap.URL{
					{Loc: "http://example.com/1"},
					{Loc: "http://example.com/2"},
				},
			},
		},
		{
			name: "error",
			fields: fields{
				linkLoader: func() LinkLoader {
					linkLoader := new(MockLinkLoader)
					linkLoader.
						On("LoadLink", "http://example.com/sitemap.xml", context.Background()).
						Return(nil, iotest.ErrTimeout)

					return linkLoader
				}(),
				logger: func() Logger {
					logger := new(MockLogger)
					logger.
						On(
							"Logf",
							"unable to load the Sitemap link %q: %s",
							"http://example.com/sitemap.xml",
							mock.MatchedBy(func(err error) bool {
								unwrappedErr := errors.Cause(err)
								return reflect.DeepEqual(unwrappedErr, iotest.ErrTimeout)
							}),
						).
						Return()

					return logger
				}(),
				registeredSitemaps: new(sync.Map),
			},
			args: args{
				ctx:         context.Background(),
				sitemapLink: "http://example.com/sitemap.xml",
			},
			wantSitemapData: sitemap.Sitemap{},
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			sitemap.SetFetch(data.fields.linkLoader.LoadLink)

			register := SitemapRegister{
				logger:             data.fields.logger,
				registeredSitemaps: data.fields.registeredSitemaps,
			}
			gotSitemapData :=
				register.loadSitemapData(data.args.ctx, data.args.sitemapLink)

			mock.AssertExpectationsForObjects(test, data.fields.linkLoader)
			assert.Equal(test, data.wantSitemapData, gotSitemapData)
		})
	}
}
