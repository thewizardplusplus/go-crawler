package urlutils

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateBaseLinks(test *testing.T) {
	type args struct {
		response        *http.Response
		baseTagValue    string
		baseHeaderNames []string
	}

	for _, data := range []struct {
		name string
		args args
		want []string
	}{
		{
			name: "with the full data",
			args: args{
				response: &http.Response{
					Header: http.Header{
						"Content-Base":     {"e/f/"},
						"Content-Location": {"c/d/"},
					},
					Request: httptest.NewRequest(
						http.MethodGet,
						"http://example.com/a/b/",
						nil,
					),
				},
				baseTagValue:    "g/h/",
				baseHeaderNames: []string{"Content-Base", "Content-Location"},
			},
			want: []string{"g/h/", "e/f/", "c/d/", "http://example.com/a/b/"},
		},
		{
			name: "without the base tag",
			args: args{
				response: &http.Response{
					Header: http.Header{
						"Content-Base":     {"e/f/"},
						"Content-Location": {"c/d/"},
					},
					Request: httptest.NewRequest(
						http.MethodGet,
						"http://example.com/a/b/",
						nil,
					),
				},
				baseTagValue:    "",
				baseHeaderNames: []string{"Content-Base", "Content-Location"},
			},
			want: []string{"e/f/", "c/d/", "http://example.com/a/b/"},
		},
		{
			name: "without the several header values",
			args: args{
				response: &http.Response{
					Header: http.Header{"Content-Location": {"c/d/"}},
					Request: httptest.NewRequest(
						http.MethodGet,
						"http://example.com/a/b/",
						nil,
					),
				},
				baseTagValue:    "g/h/",
				baseHeaderNames: []string{"Content-Base", "Content-Location"},
			},
			want: []string{"g/h/", "c/d/", "http://example.com/a/b/"},
		},
		{
			name: "without the header names",
			args: args{
				response: &http.Response{
					Header: http.Header{
						"Content-Base":     {"e/f/"},
						"Content-Location": {"c/d/"},
					},
					Request: httptest.NewRequest(
						http.MethodGet,
						"http://example.com/a/b/",
						nil,
					),
				},
				baseTagValue:    "g/h/",
				baseHeaderNames: nil,
			},
			want: []string{"g/h/", "http://example.com/a/b/"},
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			got := GenerateBaseLinks(
				data.args.response,
				data.args.baseTagValue,
				data.args.baseHeaderNames,
			)

			assert.Equal(test, data.want, got)
		})
	}
}

func TestNewLinkResolver(test *testing.T) {
	type args struct {
		baseLinks []string
	}

	for _, data := range []struct {
		name             string
		args             args
		wantLinkResolver LinkResolver
		wantErr          assert.ErrorAssertionFunc
	}{
		{
			name: "success with the single absolute link",
			args: args{
				baseLinks: []string{"e/f/", "c/d/", "http://example.com/a/b/"},
			},
			wantLinkResolver: LinkResolver{
				BaseLink: &url.URL{
					Scheme: "http",
					Host:   "example.com",
					Path:   "/a/b/c/d/e/f/",
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "success with the several absolute links",
			args: args{
				baseLinks: []string{
					"e/f/", "c/d/", "http://example.com/a/b/",
					"3/4/", "http://example.com/1/2/",
				},
			},
			wantLinkResolver: LinkResolver{
				BaseLink: &url.URL{
					Scheme: "http",
					Host:   "example.com",
					Path:   "/a/b/c/d/e/f/",
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "error",
			args: args{
				baseLinks: []string{"e/f/", ":", "http://example.com/a/b/"},
			},
			wantLinkResolver: LinkResolver{
				BaseLink: nil,
			},
			wantErr: assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			gotLinkResolver, gotErr := NewLinkResolver(data.args.baseLinks)

			assert.Equal(test, data.wantLinkResolver, gotLinkResolver)
			data.wantErr(test, gotErr)
		})
	}
}

func TestLinkResolver_ResolveLink(test *testing.T) {
	type fields struct {
		BaseLink *url.URL
	}
	type args struct {
		link string
	}

	for _, data := range []struct {
		name     string
		fields   fields
		args     args
		wantLink string
		wantErr  assert.ErrorAssertionFunc
	}{
		{
			name: "success with a relative link (relative to the root)",
			fields: fields{
				BaseLink: func() *url.URL {
					baseLink, err := url.Parse("http://example.com/one/")
					require.NoError(test, err)

					return baseLink
				}(),
			},
			args: args{
				link: "/two",
			},
			wantLink: "http://example.com/two",
			wantErr:  assert.NoError,
		},
		{
			name: "success with a relative link (relative to the current directory)",
			fields: fields{
				BaseLink: func() *url.URL {
					baseLink, err := url.Parse("http://example.com/one/")
					require.NoError(test, err)

					return baseLink
				}(),
			},
			args: args{
				link: "two",
			},
			wantLink: "http://example.com/one/two",
			wantErr:  assert.NoError,
		},
		{
			name: "success with an absolute link",
			fields: fields{
				BaseLink: func() *url.URL {
					baseLink, err := url.Parse("http://example-1.com/one/")
					require.NoError(test, err)

					return baseLink
				}(),
			},
			args: args{
				link: "http://example-2.com/two",
			},
			wantLink: "http://example-2.com/two",
			wantErr:  assert.NoError,
		},
		{
			name: "error",
			fields: fields{
				BaseLink: func() *url.URL {
					baseLink, err := url.Parse("http://example.com/one/")
					require.NoError(test, err)

					return baseLink
				}(),
			},
			args: args{
				link: ":",
			},
			wantLink: "",
			wantErr:  assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			resolver := LinkResolver{
				BaseLink: data.fields.BaseLink,
			}
			gotLink, gotErr := resolver.ResolveLink(data.args.link)

			assert.Equal(test, data.wantLink, gotLink)
			data.wantErr(test, gotErr)
		})
	}
}
