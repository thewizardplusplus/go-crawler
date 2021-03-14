package registers

import (
	"context"
	"encoding/xml"
	"sync"
	"testing"
	"testing/iotest"

	"github.com/go-log/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/yterajima/go-sitemap"
)

func TestSitemapRegister_RegisterSitemap(test *testing.T) {
	type fields struct {
		linkGenerator      LinkGenerator
		linkLoader         LinkLoader
		logger             log.Logger
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
		// TODO: Add test cases.
	} {
		test.Run(data.name, func(test *testing.T) {
			sitemap.SetFetch(data.fields.linkLoader.LoadLink)

			register := SitemapRegister{
				linkGenerator:      data.fields.linkGenerator,
				logger:             data.fields.logger,
				registeredSitemaps: data.fields.registeredSitemaps,
			}
			gotSitemapData, gotErr :=
				register.RegisterSitemap(data.args.ctx, data.args.link)

			mock.AssertExpectationsForObjects(
				test,
				data.fields.linkGenerator,
				data.fields.linkLoader,
				data.fields.logger,
			)
			assert.Equal(test, data.wantSitemapData, gotSitemapData)
			data.wantErr(test, gotErr)
		})
	}
}

func TestSitemapRegister_loadSitemapData(test *testing.T) {
	type fields struct {
		linkLoader         LinkLoader
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
		wantErr         assert.ErrorAssertionFunc
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
			wantErr: assert.NoError,
		},
		{
			name: "success with a registered Sitemap link",
			fields: fields{
				linkLoader: new(MockLinkLoader),
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
			wantErr: assert.NoError,
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
				registeredSitemaps: new(sync.Map),
			},
			args: args{
				ctx:         context.Background(),
				sitemapLink: "http://example.com/sitemap.xml",
			},
			wantSitemapData: sitemap.Sitemap{},
			wantErr:         assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			sitemap.SetFetch(data.fields.linkLoader.LoadLink)

			register := SitemapRegister{
				registeredSitemaps: data.fields.registeredSitemaps,
			}
			gotSitemapData, gotErr :=
				register.loadSitemapData(data.args.ctx, data.args.sitemapLink)

			mock.AssertExpectationsForObjects(test, data.fields.linkLoader)
			assert.Equal(test, data.wantSitemapData, gotSitemapData)
			data.wantErr(test, gotErr)
		})
	}
}
