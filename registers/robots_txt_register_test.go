package registers

import (
	"context"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/temoto/robotstxt"
	httputils "github.com/thewizardplusplus/go-http-utils"
)

func TestNewRobotsTXTRegister(test *testing.T) {
	httpClient := new(MockHTTPClient)
	got := NewRobotsTXTRegister(httpClient)

	mock.AssertExpectationsForObjects(test, httpClient)
	assert.Equal(test, httpClient, got.httpClient)
	assert.Equal(test, new(sync.Map), got.registeredRobotsTXT)
}

func TestRobotsTXTRegister_RegisterRobotsTXT(test *testing.T) {
	type fields struct {
		httpClient          httputils.HTTPClient
		registeredRobotsTXT *sync.Map
	}
	type args struct {
		ctx  context.Context
		link string
	}

	for _, data := range []struct {
		name              string
		fields            fields
		args              args
		wantRobotsTXTData *robotstxt.RobotsData
		wantErr           assert.ErrorAssertionFunc
	}{
		{
			name: "success with an unregistered robots.txt link",
			fields: fields{
				httpClient: func() httputils.HTTPClient {
					request, _ :=
						http.NewRequest(http.MethodGet, "http://example.com/robots.txt", nil)
					request = request.WithContext(context.Background())

					response := &http.Response{
						StatusCode: http.StatusOK,
						Body: ioutil.NopCloser(strings.NewReader(`
							User-agent: *
							Disallow: /
							Allow: /$
							Allow: /sitemap.xml$
							Allow: /post/
							Allow: /storage/app/media/
						`)),
					}

					httpClient := new(MockHTTPClient)
					httpClient.On("Do", request).Return(response, nil)

					return httpClient
				}(),
				registeredRobotsTXT: new(sync.Map),
			},
			args: args{
				ctx:  context.Background(),
				link: "http://example.com/test",
			},
			wantRobotsTXTData: func() *robotstxt.RobotsData {
				robotsTXTData, err := robotstxt.FromString(`
					User-agent: *
					Disallow: /
					Allow: /$
					Allow: /sitemap.xml$
					Allow: /post/
					Allow: /storage/app/media/
				`)
				require.NoError(test, err)

				return robotsTXTData
			}(),
			wantErr: assert.NoError,
		},
		{
			name: "success with a registered robots.txt link",
			fields: fields{
				httpClient: new(MockHTTPClient),
				registeredRobotsTXT: func() *sync.Map {
					robotsTXTData, err := robotstxt.FromString(`
						User-agent: *
						Disallow: /
						Allow: /$
						Allow: /sitemap.xml$
						Allow: /post/
						Allow: /storage/app/media/
					`)
					require.NoError(test, err)

					registeredRobotsTXT := new(sync.Map)
					registeredRobotsTXT.Store("http://example.com/robots.txt", robotsTXTData)

					return registeredRobotsTXT
				}(),
			},
			args: args{
				ctx:  context.Background(),
				link: "http://example.com/test",
			},
			wantRobotsTXTData: func() *robotstxt.RobotsData {
				robotsTXTData, err := robotstxt.FromString(`
					User-agent: *
					Disallow: /
					Allow: /$
					Allow: /sitemap.xml$
					Allow: /post/
					Allow: /storage/app/media/
				`)
				require.NoError(test, err)

				return robotsTXTData
			}(),
			wantErr: assert.NoError,
		},
		{
			name: "error with making of a robots.txt link",
			fields: fields{
				httpClient:          new(MockHTTPClient),
				registeredRobotsTXT: new(sync.Map),
			},
			args: args{
				ctx:  context.Background(),
				link: ":",
			},
			wantRobotsTXTData: nil,
			wantErr:           assert.Error,
		},
		{
			name: "error with loading of a robots.txt data",
			fields: fields{
				httpClient: func() httputils.HTTPClient {
					request, _ :=
						http.NewRequest(http.MethodGet, "http://example.com/robots.txt", nil)
					request = request.WithContext(context.Background())

					httpClient := new(MockHTTPClient)
					httpClient.On("Do", request).Return(nil, iotest.ErrTimeout)

					return httpClient
				}(),
				registeredRobotsTXT: new(sync.Map),
			},
			args: args{
				ctx:  context.Background(),
				link: "http://example.com/test",
			},
			wantRobotsTXTData: nil,
			wantErr:           assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			register := RobotsTXTRegister{
				httpClient:          data.fields.httpClient,
				registeredRobotsTXT: data.fields.registeredRobotsTXT,
			}
			gotRobotsTXTData, gotErr :=
				register.RegisterRobotsTXT(data.args.ctx, data.args.link)

			mock.AssertExpectationsForObjects(test, data.fields.httpClient)
			assert.Equal(test, data.wantRobotsTXTData, gotRobotsTXTData)
			data.wantErr(test, gotErr)
		})
	}
}

func TestRobotsTXTRegister_loadRobotsTXTData(test *testing.T) {
	type fields struct {
		httpClient httputils.HTTPClient
	}
	type args struct {
		ctx           context.Context
		robotsTXTLink string
	}

	for _, data := range []struct {
		name              string
		fields            fields
		args              args
		wantRobotsTXTData *robotstxt.RobotsData
		wantErr           assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			fields: fields{
				httpClient: func() httputils.HTTPClient {
					request, _ :=
						http.NewRequest(http.MethodGet, "http://example.com/robots.txt", nil)
					request = request.WithContext(context.Background())

					response := &http.Response{
						StatusCode: http.StatusOK,
						Body: ioutil.NopCloser(strings.NewReader(`
							User-agent: *
							Disallow: /
							Allow: /$
							Allow: /sitemap.xml$
							Allow: /post/
							Allow: /storage/app/media/
						`)),
					}

					httpClient := new(MockHTTPClient)
					httpClient.On("Do", request).Return(response, nil)

					return httpClient
				}(),
			},
			args: args{
				ctx:           context.Background(),
				robotsTXTLink: "http://example.com/robots.txt",
			},
			wantRobotsTXTData: func() *robotstxt.RobotsData {
				robotsTXTData, err := robotstxt.FromString(`
					User-agent: *
					Disallow: /
					Allow: /$
					Allow: /sitemap.xml$
					Allow: /post/
					Allow: /storage/app/media/
				`)
				require.NoError(test, err)

				return robotsTXTData
			}(),
			wantErr: assert.NoError,
		},
		{
			name: "error with request creating",
			fields: fields{
				httpClient: new(MockHTTPClient),
			},
			args: args{
				ctx:           context.Background(),
				robotsTXTLink: ":",
			},
			wantRobotsTXTData: nil,
			wantErr:           assert.Error,
		},
		{
			name: "error with request sending",
			fields: fields{
				httpClient: func() httputils.HTTPClient {
					request, _ :=
						http.NewRequest(http.MethodGet, "http://example.com/robots.txt", nil)
					request = request.WithContext(context.Background())

					httpClient := new(MockHTTPClient)
					httpClient.On("Do", request).Return(nil, iotest.ErrTimeout)

					return httpClient
				}(),
			},
			args: args{
				ctx:           context.Background(),
				robotsTXTLink: "http://example.com/robots.txt",
			},
			wantRobotsTXTData: nil,
			wantErr:           assert.Error,
		},
		{
			name: "error with response parsing",
			fields: fields{
				httpClient: func() httputils.HTTPClient {
					request, _ :=
						http.NewRequest(http.MethodGet, "http://example.com/robots.txt", nil)
					request = request.WithContext(context.Background())

					response := &http.Response{
						StatusCode: http.StatusOK,
						Body: ioutil.NopCloser(iotest.TimeoutReader(strings.NewReader(`
							User-agent: *
							Disallow: /
							Allow: /$
							Allow: /sitemap.xml$
							Allow: /post/
							Allow: /storage/app/media/
						`))),
					}

					httpClient := new(MockHTTPClient)
					httpClient.On("Do", request).Return(response, nil)

					return httpClient
				}(),
			},
			args: args{
				ctx:           context.Background(),
				robotsTXTLink: "http://example.com/robots.txt",
			},
			wantRobotsTXTData: nil,
			wantErr:           assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			register := RobotsTXTRegister{
				httpClient: data.fields.httpClient,
			}
			gotRobotsTXTData, gotErr :=
				register.loadRobotsTXTData(data.args.ctx, data.args.robotsTXTLink)

			mock.AssertExpectationsForObjects(test, data.fields.httpClient)
			assert.Equal(test, data.wantRobotsTXTData, gotRobotsTXTData)
			data.wantErr(test, gotErr)
		})
	}
}
