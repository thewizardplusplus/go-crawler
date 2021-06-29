package extractors

import (
	"context"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	urlutils "github.com/thewizardplusplus/go-crawler/url-utils"
	htmlselector "github.com/thewizardplusplus/go-html-selector"
	httputils "github.com/thewizardplusplus/go-http-utils"
)

func TestDefaultExtractor_ExtractLinks(test *testing.T) {
	type fields struct {
		TrimLink   urlutils.LinkTrimming
		HTTPClient httputils.HTTPClient
		Filters    htmlselector.OptimizedFilterGroup
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
			name: "success without links",
			fields: fields{
				TrimLink: urlutils.DoNotTrimLink,
				HTTPClient: func() httputils.HTTPClient {
					request, _ := http.NewRequest(http.MethodGet, "http://example.com/", nil)
					request = request.WithContext(context.Background())

					response := &http.Response{
						Body: ioutil.NopCloser(strings.NewReader("")),
					}

					httpClient := new(MockHTTPClient)
					httpClient.On("Do", request).Return(response, nil)

					return httpClient
				}(),
				Filters: htmlselector.OptimizeFilters(htmlselector.FilterGroup{
					"a": {"href"},
				}),
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
			name: "success with links (without link trimming)",
			fields: fields{
				TrimLink: urlutils.DoNotTrimLink,
				HTTPClient: func() httputils.HTTPClient {
					request, _ := http.NewRequest(http.MethodGet, "http://example.com/", nil)
					request = request.WithContext(context.Background())

					response := &http.Response{
						Body: ioutil.NopCloser(strings.NewReader(`
							<ul>
								<li><a href="  http://example.com/1  ">1</a></li>
								<li><a href="  http://example.com/2  ">2</a></li>
							</ul>
						`)),
					}

					httpClient := new(MockHTTPClient)
					httpClient.On("Do", request).Return(response, nil)

					return httpClient
				}(),
				Filters: htmlselector.OptimizeFilters(htmlselector.FilterGroup{
					"a": {"href"},
				}),
			},
			args: args{
				ctx:      context.Background(),
				threadID: 23,
				link:     "http://example.com/",
			},
			wantLinks: []string{"  http://example.com/1  ", "  http://example.com/2  "},
			wantErr:   assert.NoError,
		},
		{
			name: "success with links (with link trimming)",
			fields: fields{
				TrimLink: urlutils.TrimLink,
				HTTPClient: func() httputils.HTTPClient {
					request, _ := http.NewRequest(http.MethodGet, "http://example.com/", nil)
					request = request.WithContext(context.Background())

					response := &http.Response{
						Body: ioutil.NopCloser(strings.NewReader(`
							<ul>
								<li><a href="  http://example.com/1  ">1</a></li>
								<li><a href="  http://example.com/2  ">2</a></li>
							</ul>
						`)),
					}

					httpClient := new(MockHTTPClient)
					httpClient.On("Do", request).Return(response, nil)

					return httpClient
				}(),
				Filters: htmlselector.OptimizeFilters(htmlselector.FilterGroup{
					"a": {"href"},
				}),
			},
			args: args{
				ctx:      context.Background(),
				threadID: 23,
				link:     "http://example.com/",
			},
			wantLinks: []string{"http://example.com/1", "http://example.com/2"},
			wantErr:   assert.NoError,
		},
		{
			name: "error with request creating",
			fields: fields{
				TrimLink:   urlutils.DoNotTrimLink,
				HTTPClient: new(MockHTTPClient),
				Filters: htmlselector.OptimizeFilters(htmlselector.FilterGroup{
					"a": {"href"},
				}),
			},
			args: args{
				ctx:      context.Background(),
				threadID: 23,
				link:     ":",
			},
			wantLinks: nil,
			wantErr:   assert.Error,
		},
		{
			name: "error with request sending",
			fields: fields{
				TrimLink: urlutils.DoNotTrimLink,
				HTTPClient: func() httputils.HTTPClient {
					request, _ := http.NewRequest(http.MethodGet, "http://example.com/", nil)
					request = request.WithContext(context.Background())

					httpClient := new(MockHTTPClient)
					httpClient.On("Do", request).Return(nil, iotest.ErrTimeout)

					return httpClient
				}(),
				Filters: htmlselector.OptimizeFilters(htmlselector.FilterGroup{
					"a": {"href"},
				}),
			},
			args: args{
				ctx:      context.Background(),
				threadID: 23,
				link:     "http://example.com/",
			},
			wantLinks: nil,
			wantErr:   assert.Error,
		},
		{
			name: "error with tags selecting",
			fields: fields{
				TrimLink: urlutils.DoNotTrimLink,
				HTTPClient: func() httputils.HTTPClient {
					request, _ := http.NewRequest(http.MethodGet, "http://example.com/", nil)
					request = request.WithContext(context.Background())

					response := &http.Response{
						Body: ioutil.NopCloser(iotest.TimeoutReader(strings.NewReader(`
							<ul>
								<li><a href="  http://example.com/1  ">1</a></li>
								<li><a href="  http://example.com/2  ">2</a></li>
							</ul>
						`))),
					}

					httpClient := new(MockHTTPClient)
					httpClient.On("Do", request).Return(response, nil)

					return httpClient
				}(),
				Filters: htmlselector.OptimizeFilters(htmlselector.FilterGroup{
					"a": {"href"},
				}),
			},
			args: args{
				ctx:      context.Background(),
				threadID: 23,
				link:     "http://example.com/",
			},
			wantLinks: nil,
			wantErr:   assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			extractor := DefaultExtractor{
				TrimLink:   data.fields.TrimLink,
				HTTPClient: data.fields.HTTPClient,
				Filters:    data.fields.Filters,
			}
			gotLinks, gotErr := extractor.ExtractLinks(
				data.args.ctx,
				data.args.threadID,
				data.args.link,
			)

			mock.AssertExpectationsForObjects(test, data.fields.HTTPClient)
			assert.Equal(test, data.wantLinks, gotLinks)
			data.wantErr(test, gotErr)
		})
	}
}
