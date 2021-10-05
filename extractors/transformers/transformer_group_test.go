package transformers

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/go-crawler/models"
)

func TestTransformerGroup_TransformLinks(test *testing.T) {
	type args struct {
		links           []string
		response        *http.Response
		responseContent []byte
	}

	for _, data := range []struct {
		name         string
		transformers TransformerGroup
		args         args
		wantLinks    []string
		wantErr      assert.ErrorAssertionFunc
	}{
		{
			name:         "success without transformers",
			transformers: nil,
			args: args{
				links: []string{"http://example.com/1", "http://example.com/2"},
				response: &http.Response{
					Body: ioutil.NopCloser(strings.NewReader(`
						<ul>
							<li><a href="http://example.com/1">1</a></li>
							<li><a href="http://example.com/2">2</a></li>
						</ul>
					`)),
					Request: httptest.NewRequest(http.MethodGet, "http://example.com/", nil),
				},
				responseContent: []byte(`
					<ul>
						<li><a href="http://example.com/1">1</a></li>
						<li><a href="http://example.com/2">2</a></li>
					</ul>
				`),
			},
			wantLinks: []string{"http://example.com/1", "http://example.com/2"},
			wantErr:   assert.NoError,
		},
		{
			name: "success with transformers",
			transformers: TransformerGroup{
				func() models.LinkTransformer {
					links := []string{"http://example.com/1", "http://example.com/2"}
					transformedLinks := []string{
						"http://example.com/1/transformed/1",
						"http://example.com/2/transformed/1",
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

					transformer := new(MockLinkTransformer)
					transformer.
						On("TransformLinks", links, response, []byte(responseContent)).
						Return(transformedLinks, nil)

					return transformer
				}(),
				func() models.LinkTransformer {
					links := []string{
						"http://example.com/1/transformed/1",
						"http://example.com/2/transformed/1",
					}
					transformedLinks := []string{
						"http://example.com/1/transformed/2",
						"http://example.com/2/transformed/2",
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

					transformer := new(MockLinkTransformer)
					transformer.
						On("TransformLinks", links, response, []byte(responseContent)).
						Return(transformedLinks, nil)

					return transformer
				}(),
			},
			args: args{
				links: []string{"http://example.com/1", "http://example.com/2"},
				response: &http.Response{
					Body: ioutil.NopCloser(strings.NewReader(`
						<ul>
							<li><a href="http://example.com/1">1</a></li>
							<li><a href="http://example.com/2">2</a></li>
						</ul>
					`)),
					Request: httptest.NewRequest(http.MethodGet, "http://example.com/", nil),
				},
				responseContent: func() []byte {
					return []byte(`
						<ul>
							<li><a href="http://example.com/1">1</a></li>
							<li><a href="http://example.com/2">2</a></li>
						</ul>
					`)
				}(),
			},
			wantLinks: []string{
				"http://example.com/1/transformed/2",
				"http://example.com/2/transformed/2",
			},
			wantErr: assert.NoError,
		},
		{
			name: "error",
			transformers: TransformerGroup{
				func() models.LinkTransformer {
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

					transformer := new(MockLinkTransformer)
					transformer.
						On("TransformLinks", links, response, []byte(responseContent)).
						Return(nil, iotest.ErrTimeout)

					return transformer
				}(),
				new(MockLinkTransformer),
			},
			args: args{
				links: []string{"http://example.com/1", "http://example.com/2"},
				response: &http.Response{
					Body: ioutil.NopCloser(strings.NewReader(`
						<ul>
							<li><a href="http://example.com/1">1</a></li>
							<li><a href="http://example.com/2">2</a></li>
						</ul>
					`)),
					Request: httptest.NewRequest(http.MethodGet, "http://example.com/", nil),
				},
				responseContent: func() []byte {
					return []byte(`
						<ul>
							<li><a href="http://example.com/1">1</a></li>
							<li><a href="http://example.com/2">2</a></li>
						</ul>
					`)
				}(),
			},
			wantLinks: nil,
			wantErr:   assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			gotLinks, gotErr := data.transformers.TransformLinks(
				data.args.links,
				data.args.response,
				data.args.responseContent,
			)

			for _, transformer := range data.transformers {
				mock.AssertExpectationsForObjects(test, transformer)
			}
			assert.Equal(test, data.wantLinks, gotLinks)
			data.wantErr(test, gotErr)
		})
	}
}
