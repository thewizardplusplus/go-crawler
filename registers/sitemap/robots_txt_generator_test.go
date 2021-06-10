package sitemap

import (
	"context"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
	"github.com/thewizardplusplus/go-crawler/registers"
)

func TestRobotsTXTGenerator_GenerateLinks(test *testing.T) {
	type fields struct {
		RobotsTXTRegister registers.RobotsTXTRegister
	}
	type args struct {
		ctx      context.Context
		baseLink string
	}

	for _, data := range []struct {
		name             string
		fields           fields
		args             args
		wantSitemapLinks []string
		wantErr          assert.ErrorAssertionFunc
	}{
		{
			name: "success without Sitemap links",
			fields: fields{
				RobotsTXTRegister: func() registers.RobotsTXTRegister {
					request, _ :=
						http.NewRequest(http.MethodGet, "http://example.com/robots.txt", nil)
					request = request.WithContext(context.Background())

					response := &http.Response{
						StatusCode: http.StatusOK,
						Body: ioutil.NopCloser(strings.NewReader(`
							User-agent: *
							Disallow: /

							User-agent: go-crawler
							Disallow: /
							Allow: /$
							Allow: /sitemap.xml$
							Allow: /post/
							Allow: /storage/app/media/
						`)),
					}

					httpClient := new(MockHTTPClient)
					httpClient.On("Do", request).Return(response, nil)

					register := registers.NewRobotsTXTRegister(httpClient)
					return register
				}(),
			},
			args: args{
				ctx:      context.Background(),
				baseLink: "http://example.com/test",
			},
			wantSitemapLinks: nil,
			wantErr:          assert.NoError,
		},
		{
			name: "success with Sitemap links",
			fields: fields{
				RobotsTXTRegister: func() registers.RobotsTXTRegister {
					request, _ :=
						http.NewRequest(http.MethodGet, "http://example.com/robots.txt", nil)
					request = request.WithContext(context.Background())

					response := &http.Response{
						StatusCode: http.StatusOK,
						Body: ioutil.NopCloser(strings.NewReader(`
							User-agent: *
							Disallow: /

							User-agent: go-crawler
							Disallow: /
							Allow: /$
							Allow: /sitemap.xml$
							Allow: /post/
							Allow: /storage/app/media/

							Sitemap: http://example.com/sitemap_1.xml
							Sitemap: http://example.com/sitemap_2.xml
						`)),
					}

					httpClient := new(MockHTTPClient)
					httpClient.On("Do", request).Return(response, nil)

					register := registers.NewRobotsTXTRegister(httpClient)
					return register
				}(),
			},
			args: args{
				ctx:      context.Background(),
				baseLink: "http://example.com/test",
			},
			wantSitemapLinks: []string{
				"http://example.com/sitemap_1.xml",
				"http://example.com/sitemap_2.xml",
			},
			wantErr: assert.NoError,
		},
		{
			name: "error",
			fields: fields{
				RobotsTXTRegister: func() registers.RobotsTXTRegister {
					request, _ :=
						http.NewRequest(http.MethodGet, "http://example.com/robots.txt", nil)
					request = request.WithContext(context.Background())

					httpClient := new(MockHTTPClient)
					httpClient.On("Do", request).Return(nil, iotest.ErrTimeout)

					register := registers.NewRobotsTXTRegister(httpClient)
					return register
				}(),
			},
			args: args{
				ctx:      context.Background(),
				baseLink: "http://example.com/test",
			},
			wantSitemapLinks: nil,
			wantErr:          assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			generator := RobotsTXTGenerator{
				RobotsTXTRegister: data.fields.RobotsTXTRegister,
			}
			gotSitemapLinks, gotErr :=
				generator.GenerateLinks(data.args.ctx, data.args.baseLink)

			assert.Equal(test, data.wantSitemapLinks, gotSitemapLinks)
			data.wantErr(test, gotErr)
		})
	}
}
