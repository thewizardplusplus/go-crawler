package registers

import (
	"context"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/temoto/robotstxt"
)

func TestNewRobotsTXTRegister(test *testing.T) {
	httpClient := new(MockHTTPClient)
	got := NewRobotsTXTRegister(httpClient)

	mock.AssertExpectationsForObjects(test, httpClient)
	assert.Equal(test, httpClient, got.httpClient)
}

func Test_makeRobotsTXTLink(test *testing.T) {
	type args struct {
		regularLink string
	}

	for _, data := range []struct {
		name              string
		args              args
		wantRobotsTXTLink string
		wantErr           assert.ErrorAssertionFunc
	}{
		{
			name: "success with a path only",
			args: args{
				regularLink: "http://example.com/test",
			},
			wantRobotsTXTLink: "http://example.com/robots.txt",
			wantErr:           assert.NoError,
		},
		{
			name: "success with an HTTPS scheme",
			args: args{
				regularLink: "https://example.com/test",
			},
			wantRobotsTXTLink: "https://example.com/robots.txt",
			wantErr:           assert.NoError,
		},
		{
			name: "success with an user",
			args: args{
				regularLink: "http://username:password@example.com/test",
			},
			wantRobotsTXTLink: "http://username:password@example.com/robots.txt",
			wantErr:           assert.NoError,
		},
		{
			name: "success with a query",
			args: args{
				regularLink: "http://example.com/test?key=value",
			},
			wantRobotsTXTLink: "http://example.com/robots.txt",
			wantErr:           assert.NoError,
		},
		{
			name: "success with a fragment",
			args: args{
				regularLink: "http://example.com/test#fragment",
			},
			wantRobotsTXTLink: "http://example.com/robots.txt",
			wantErr:           assert.NoError,
		},
		{
			name: "error",
			args: args{
				regularLink: ":",
			},
			wantRobotsTXTLink: "",
			wantErr:           assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			gotRobotsTXTLink, gotErr := makeRobotsTXTLink(data.args.regularLink)

			assert.Equal(test, data.wantRobotsTXTLink, gotRobotsTXTLink)
			data.wantErr(test, gotErr)
		})
	}
}

func Test_loadRobotsTXTData(test *testing.T) {
	type args struct {
		ctx           context.Context
		httpClient    HTTPClient
		robotsTXTLink string
	}

	for _, data := range []struct {
		name              string
		args              args
		wantRobotsTXTData *robotstxt.RobotsData
		wantErr           assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				httpClient: func() HTTPClient {
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
			args: args{
				ctx:           context.Background(),
				httpClient:    new(MockHTTPClient),
				robotsTXTLink: ":",
			},
			wantRobotsTXTData: nil,
			wantErr:           assert.Error,
		},
		{
			name: "error with request sending",
			args: args{
				ctx: context.Background(),
				httpClient: func() HTTPClient {
					request, _ :=
						http.NewRequest(http.MethodGet, "http://example.com/robots.txt", nil)
					request = request.WithContext(context.Background())

					httpClient := new(MockHTTPClient)
					httpClient.On("Do", request).Return(nil, iotest.ErrTimeout)

					return httpClient
				}(),
				robotsTXTLink: "http://example.com/robots.txt",
			},
			wantRobotsTXTData: nil,
			wantErr:           assert.Error,
		},
		{
			name: "error with response parsing",
			args: args{
				ctx: context.Background(),
				httpClient: func() HTTPClient {
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
				robotsTXTLink: "http://example.com/robots.txt",
			},
			wantRobotsTXTData: nil,
			wantErr:           assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			gotRobotsTXTData, gotErr := loadRobotsTXTData(
				data.args.ctx,
				data.args.httpClient,
				data.args.robotsTXTLink,
			)

			mock.AssertExpectationsForObjects(test, data.args.httpClient)
			assert.Equal(test, data.wantRobotsTXTData, gotRobotsTXTData)
			data.wantErr(test, gotErr)
		})
	}
}
