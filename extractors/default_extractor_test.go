package extractors

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	htmlselector "github.com/thewizardplusplus/go-html-selector"
	httputils "github.com/thewizardplusplus/go-http-utils"
)

func TestDefaultExtractor_ExtractLinks(test *testing.T) {
	type fields struct {
		HTTPClient      httputils.HTTPClient
		Filters         htmlselector.OptimizedFilterGroup
		LinkTransformer LinkTransformer
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
			name: "success without transformation of the links",
			fields: fields{
				HTTPClient: func() httputils.HTTPClient {
					request, _ := http.NewRequest(http.MethodGet, "http://example.com/", nil)
					request = request.WithContext(context.Background())

					response := &http.Response{
						Body: ioutil.NopCloser(strings.NewReader(`
							<ul>
								<li><a href="http://example.com/1">1</a></li>
								<li><a href="http://example.com/2">2</a></li>
							</ul>
						`)),
						Request: httptest.NewRequest(http.MethodGet, "http://example.com/", nil),
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
			name: "success with transformation of the links",
			fields: fields{
				HTTPClient: func() httputils.HTTPClient {
					request, _ := http.NewRequest(http.MethodGet, "http://example.com/", nil)
					request = request.WithContext(context.Background())

					responseContent := `
						<ul>
							<li><a href="http://example.com/1">1</a></li>
							<li><a href="http://example.com/2">2</a></li>
						</ul>
					`
					response := &http.Response{
						Body:    ioutil.NopCloser(strings.NewReader(responseContent)),
						Request: httptest.NewRequest(http.MethodGet, "http://example.com/", nil),
					}

					httpClient := new(MockHTTPClient)
					httpClient.On("Do", request).Return(response, nil)

					return httpClient
				}(),
				Filters: htmlselector.OptimizeFilters(htmlselector.FilterGroup{
					"a": {"href"},
				}),
				LinkTransformer: func() LinkTransformer {
					links := []string{"http://example.com/1", "http://example.com/2"}
					transformedLinks := []string{
						"http://example.com/transformed/1",
						"http://example.com/transformed/2",
					}

					responseContent := `
						<ul>
							<li><a href="http://example.com/1">1</a></li>
							<li><a href="http://example.com/2">2</a></li>
						</ul>
					`
					response := &http.Response{
						Body:    ioutil.NopCloser(strings.NewReader(responseContent)),
						Request: httptest.NewRequest(http.MethodGet, "http://example.com/", nil),
					}
					ioutil.ReadAll(response.Body) // nolint: errcheck

					linkTransformer := new(MockLinkTransformer)
					linkTransformer.
						On("TransformLinks", links, response, []byte(responseContent)).
						Return(transformedLinks, nil)

					return linkTransformer
				}(),
			},
			args: args{
				ctx:      context.Background(),
				threadID: 23,
				link:     "http://example.com/",
			},
			wantLinks: []string{
				"http://example.com/transformed/1",
				"http://example.com/transformed/2",
			},
			wantErr: assert.NoError,
		},
		{
			name: "error with loading of the data",
			fields: fields{
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
			name: "error with transformation of the links",
			fields: fields{
				HTTPClient: func() httputils.HTTPClient {
					request, _ := http.NewRequest(http.MethodGet, "http://example.com/", nil)
					request = request.WithContext(context.Background())

					responseContent := `
						<ul>
							<li><a href="http://example.com/1">1</a></li>
							<li><a href="http://example.com/2">2</a></li>
						</ul>
					`
					response := &http.Response{
						Body:    ioutil.NopCloser(strings.NewReader(responseContent)),
						Request: httptest.NewRequest(http.MethodGet, "http://example.com/", nil),
					}

					httpClient := new(MockHTTPClient)
					httpClient.On("Do", request).Return(response, nil)

					return httpClient
				}(),
				Filters: htmlselector.OptimizeFilters(htmlselector.FilterGroup{
					"a": {"href"},
				}),
				LinkTransformer: func() LinkTransformer {
					links := []string{"http://example.com/1", "http://example.com/2"}

					responseContent := `
						<ul>
							<li><a href="http://example.com/1">1</a></li>
							<li><a href="http://example.com/2">2</a></li>
						</ul>
					`
					response := &http.Response{
						Body:    ioutil.NopCloser(strings.NewReader(responseContent)),
						Request: httptest.NewRequest(http.MethodGet, "http://example.com/", nil),
					}
					ioutil.ReadAll(response.Body) // nolint: errcheck

					linkTransformer := new(MockLinkTransformer)
					linkTransformer.
						On("TransformLinks", links, response, []byte(responseContent)).
						Return(nil, iotest.ErrTimeout)

					return linkTransformer
				}(),
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
				HTTPClient:      data.fields.HTTPClient,
				Filters:         data.fields.Filters,
				LinkTransformer: data.fields.LinkTransformer,
			}
			gotLinks, gotErr := extractor.ExtractLinks(
				data.args.ctx,
				data.args.threadID,
				data.args.link,
			)

			mock.AssertExpectationsForObjects(test, data.fields.HTTPClient)
			if data.fields.LinkTransformer != nil {
				mock.AssertExpectationsForObjects(test, data.fields.LinkTransformer)
			}
			assert.Equal(test, data.wantLinks, gotLinks)
			data.wantErr(test, gotErr)
		})
	}
}

func TestDefaultExtractor_loadData(test *testing.T) {
	type fields struct {
		HTTPClient httputils.HTTPClient
	}
	type args struct {
		ctx  context.Context
		link string
	}

	for _, data := range []struct {
		name         string
		fields       fields
		args         args
		wantData     []byte
		wantResponse *http.Response
		wantErr      assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			fields: fields{
				HTTPClient: func() httputils.HTTPClient {
					request, _ := http.NewRequest(http.MethodGet, "http://example.com/", nil)
					request = request.WithContext(context.Background())

					response := &http.Response{
						Body: ioutil.NopCloser(strings.NewReader("data")),
					}

					httpClient := new(MockHTTPClient)
					httpClient.On("Do", request).Return(response, nil)

					return httpClient
				}(),
			},
			args: args{
				ctx:  context.Background(),
				link: "http://example.com/",
			},
			wantData: []byte("data"),
			wantResponse: func() *http.Response {
				data := strings.NewReader("data")
				data.Seek(0, io.SeekEnd) // nolint: errcheck

				return &http.Response{
					Body: ioutil.NopCloser(data),
				}
			}(),
			wantErr: assert.NoError,
		},
		{
			name: "error with request creating",
			fields: fields{
				HTTPClient: new(MockHTTPClient),
			},
			args: args{
				ctx:  context.Background(),
				link: ":",
			},
			wantData:     nil,
			wantResponse: nil,
			wantErr:      assert.Error,
		},
		{
			name: "error with request sending",
			fields: fields{
				HTTPClient: func() httputils.HTTPClient {
					request, _ := http.NewRequest(http.MethodGet, "http://example.com/", nil)
					request = request.WithContext(context.Background())

					httpClient := new(MockHTTPClient)
					httpClient.On("Do", request).Return(nil, iotest.ErrTimeout)

					return httpClient
				}(),
			},
			args: args{
				ctx:  context.Background(),
				link: "http://example.com/",
			},
			wantData:     nil,
			wantResponse: nil,
			wantErr:      assert.Error,
		},
		{
			name: "error with request reading",
			fields: fields{
				HTTPClient: func() httputils.HTTPClient {
					request, _ := http.NewRequest(http.MethodGet, "http://example.com/", nil)
					request = request.WithContext(context.Background())

					response := &http.Response{
						Body: ioutil.NopCloser(iotest.TimeoutReader(strings.NewReader("data"))),
					}

					httpClient := new(MockHTTPClient)
					httpClient.On("Do", request).Return(response, nil)

					return httpClient
				}(),
			},
			args: args{
				ctx:  context.Background(),
				link: "http://example.com/",
			},
			wantData:     nil,
			wantResponse: nil,
			wantErr:      assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			extractor := DefaultExtractor{
				HTTPClient: data.fields.HTTPClient,
			}
			gotData, gotResponse, gotErr :=
				extractor.loadData(data.args.ctx, data.args.link)

			mock.AssertExpectationsForObjects(test, data.fields.HTTPClient)
			assert.Equal(test, data.wantData, gotData)
			assert.Equal(test, data.wantResponse, gotResponse)
			data.wantErr(test, gotErr)
		})
	}
}

func TestDefaultExtractor_selectLinks(test *testing.T) {
	type fields struct {
		Filters htmlselector.OptimizedFilterGroup
	}
	type args struct {
		data []byte
	}

	for _, data := range []struct {
		name      string
		fields    fields
		args      args
		wantLinks []string
	}{
		{
			name: "without links",
			fields: fields{
				Filters: htmlselector.OptimizeFilters(htmlselector.FilterGroup{
					"a": {"href"},
				}),
			},
			args: args{
				data: []byte(""),
			},
			wantLinks: nil,
		},
		{
			name: "with links",
			fields: fields{
				Filters: htmlselector.OptimizeFilters(htmlselector.FilterGroup{
					"a": {"href"},
				}),
			},
			args: args{
				data: []byte(`
					<ul>
						<li><a href="http://example.com/1">1</a></li>
						<li><a href="http://example.com/2">2</a></li>
					</ul>
				`),
			},
			wantLinks: []string{"http://example.com/1", "http://example.com/2"},
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			extractor := DefaultExtractor{
				Filters: data.fields.Filters,
			}
			gotLinks := extractor.selectLinks(data.args.data)

			assert.Equal(test, data.wantLinks, gotLinks)
		})
	}
}
