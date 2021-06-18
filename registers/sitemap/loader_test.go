package sitemap

import (
	"bytes"
	"compress/gzip"
	"context"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	httputils "github.com/thewizardplusplus/go-http-utils"
)

func TestLoader_LoadLink(test *testing.T) {
	type fields struct {
		HTTPClient httputils.HTTPClient
	}
	type args struct {
		link    string
		options interface{}
	}

	for _, data := range []struct {
		name             string
		fields           fields
		args             args
		wantResponseData []byte
		wantErr          assert.ErrorAssertionFunc
	}{
		{
			name: "success without compression",
			fields: fields{
				HTTPClient: func() httputils.HTTPClient {
					request, _ :=
						http.NewRequest(http.MethodGet, "http://example.com/sitemap.xml", nil)
					request = request.WithContext(context.Background())

					response := &http.Response{
						StatusCode: http.StatusOK,
						Body: ioutil.NopCloser(strings.NewReader(
							`<?xml version="1.0" encoding="UTF-8" ?>` +
								`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">` +
								"<url><loc>http://example.com/1</loc></url>" +
								"<url><loc>http://example.com/2</loc></url>" +
								"</urlset>",
						)),
					}

					httpClient := new(MockHTTPClient)
					httpClient.On("Do", request).Return(response, nil)

					return httpClient
				}(),
			},
			args: args{
				link:    "http://example.com/sitemap.xml",
				options: context.Background(),
			},
			wantResponseData: []byte(
				`<?xml version="1.0" encoding="UTF-8" ?>` +
					`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">` +
					"<url><loc>http://example.com/1</loc></url>" +
					"<url><loc>http://example.com/2</loc></url>" +
					"</urlset>",
			),
			wantErr: assert.NoError,
		},
		{
			name: "success with compression",
			fields: fields{
				HTTPClient: func() httputils.HTTPClient {
					request, _ :=
						http.NewRequest(http.MethodGet, "http://example.com/sitemap.xml", nil)
					request = request.WithContext(context.Background())

					var buffer bytes.Buffer
					compressingWriter := gzip.NewWriter(&buffer)
					_, err := compressingWriter.Write([]byte(
						`<?xml version="1.0" encoding="UTF-8" ?>` +
							`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">` +
							"<url><loc>http://example.com/1</loc></url>" +
							"<url><loc>http://example.com/2</loc></url>" +
							"</urlset>",
					))
					require.NoError(test, err)
					err = compressingWriter.Close()
					require.NoError(test, err)

					response := &http.Response{
						StatusCode: http.StatusOK,
						Header:     http.Header{"Content-Encoding": {"gzip"}},
						Body:       ioutil.NopCloser(&buffer),
					}

					httpClient := new(MockHTTPClient)
					httpClient.On("Do", request).Return(response, nil)

					return httpClient
				}(),
			},
			args: args{
				link:    "http://example.com/sitemap.xml",
				options: context.Background(),
			},
			wantResponseData: []byte(
				`<?xml version="1.0" encoding="UTF-8" ?>` +
					`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">` +
					"<url><loc>http://example.com/1</loc></url>" +
					"<url><loc>http://example.com/2</loc></url>" +
					"</urlset>",
			),
			wantErr: assert.NoError,
		},
		{
			name: "error with the request creating",
			fields: fields{
				HTTPClient: new(MockHTTPClient),
			},
			args: args{
				link:    ":",
				options: context.Background(),
			},
			wantResponseData: nil,
			wantErr:          assert.Error,
		},
		{
			name: "error with the request sending",
			fields: fields{
				HTTPClient: func() httputils.HTTPClient {
					request, _ :=
						http.NewRequest(http.MethodGet, "http://example.com/sitemap.xml", nil)
					request = request.WithContext(context.Background())

					httpClient := new(MockHTTPClient)
					httpClient.On("Do", request).Return(nil, iotest.ErrTimeout)

					return httpClient
				}(),
			},
			args: args{
				link:    "http://example.com/sitemap.xml",
				options: context.Background(),
			},
			wantResponseData: nil,
			wantErr:          assert.Error,
		},
		{
			name: "error with the gzip reader creating",
			fields: fields{
				HTTPClient: func() httputils.HTTPClient {
					request, _ :=
						http.NewRequest(http.MethodGet, "http://example.com/sitemap.xml", nil)
					request = request.WithContext(context.Background())

					response := &http.Response{
						StatusCode: http.StatusOK,
						Header:     http.Header{"Content-Encoding": {"gzip"}},
						Body: ioutil.NopCloser(strings.NewReader(
							`<?xml version="1.0" encoding="UTF-8" ?>` +
								`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">` +
								"<url><loc>http://example.com/1</loc></url>" +
								"<url><loc>http://example.com/2</loc></url>" +
								"</urlset>",
						)),
					}

					httpClient := new(MockHTTPClient)
					httpClient.On("Do", request).Return(response, nil)

					return httpClient
				}(),
			},
			args: args{
				link:    "http://example.com/sitemap.xml",
				options: context.Background(),
			},
			wantResponseData: nil,
			wantErr:          assert.Error,
		},
		{
			name: "error with the response reading",
			fields: fields{
				HTTPClient: func() httputils.HTTPClient {
					request, _ :=
						http.NewRequest(http.MethodGet, "http://example.com/sitemap.xml", nil)
					request = request.WithContext(context.Background())

					response := &http.Response{
						StatusCode: http.StatusOK,
						Body: ioutil.NopCloser(iotest.TimeoutReader(strings.NewReader(
							`<?xml version="1.0" encoding="UTF-8" ?>` +
								`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">` +
								"<url><loc>http://example.com/1</loc></url>" +
								"<url><loc>http://example.com/2</loc></url>" +
								"</urlset>",
						))),
					}

					httpClient := new(MockHTTPClient)
					httpClient.On("Do", request).Return(response, nil)

					return httpClient
				}(),
			},
			args: args{
				link:    "http://example.com/sitemap.xml",
				options: context.Background(),
			},
			wantResponseData: nil,
			wantErr:          assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			loader := Loader{
				HTTPClient: data.fields.HTTPClient,
			}
			gotResponseData, gotErr := loader.LoadLink(data.args.link, data.args.options)

			assert.Equal(test, data.wantResponseData, gotResponseData)
			data.wantErr(test, gotErr)
		})
	}
}
