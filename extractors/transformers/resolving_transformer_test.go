package transformers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	urlutils "github.com/thewizardplusplus/go-crawler/url-utils"
)

func TestResolvingTransformer_TransformLinks(test *testing.T) {
	type fields struct {
		BaseHeaderNames []string
	}
	type args struct {
		links           []string
		response        *http.Response
		responseContent []byte
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
				BaseHeaderNames: urlutils.DefaultBaseHeaderNames,
			},
			args: args{
				links: nil,
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
				responseContent: []byte(`<base href="g/h/" />`),
			},
			wantLinks: nil,
			wantErr:   assert.NoError,
		},
		{
			name: "success with links",
			fields: fields{
				BaseHeaderNames: urlutils.DefaultBaseHeaderNames,
			},
			args: args{
				links: []string{"one", "two"},
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
				responseContent: []byte(`<base href="g/h/" />`),
			},
			wantLinks: []string{
				"http://example.com/a/b/c/d/e/f/g/h/one",
				"http://example.com/a/b/c/d/e/f/g/h/two",
			},
			wantErr: assert.NoError,
		},
		{
			name: "error with constructing of the link resolver",
			fields: fields{
				BaseHeaderNames: urlutils.DefaultBaseHeaderNames,
			},
			args: args{
				links: []string{"one", "two"},
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
				responseContent: []byte(`<base href=":" />`),
			},
			wantLinks: nil,
			wantErr:   assert.Error,
		},
		{
			name: "error with resolving of the link",
			fields: fields{
				BaseHeaderNames: urlutils.DefaultBaseHeaderNames,
			},
			args: args{
				links: []string{":", "two"},
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
				responseContent: []byte(`<base href="g/h/" />`),
			},
			wantLinks: nil,
			wantErr:   assert.Error,
		},
	} {
		test.Run(data.name, func(t *testing.T) {
			transformer := ResolvingTransformer{
				BaseHeaderNames: data.fields.BaseHeaderNames,
			}
			gotLinks, gotErr := transformer.TransformLinks(
				data.args.links,
				data.args.response,
				data.args.responseContent,
			)

			assert.Equal(test, data.wantLinks, gotLinks)
			data.wantErr(test, gotErr)
		})
	}
}

func Test_selectBaseTag(test *testing.T) {
	type args struct {
		data []byte
	}

	for _, data := range []struct {
		name string
		args args
		want string
	}{
		{
			name: "without the base tag",
			args: args{
				data: []byte(`
					<ul>
						<li><a href="http://example.com/1">1</a></li>
						<li><a href="http://example.com/2">2</a></li>
					</ul>
				`),
			},
			want: "",
		},
		{
			name: "with the base tag without the href attribute",
			args: args{
				data: []byte(`
					<base target="_blank" />

					<ul>
						<li><a href="http://example.com/1">1</a></li>
						<li><a href="http://example.com/2">2</a></li>
					</ul>
				`),
			},
			want: "",
		},
		{
			name: "with the base tag with the href attribute",
			args: args{
				data: []byte(`
					<base href="http://example.com/" />

					<ul>
						<li><a href="1">1</a></li>
						<li><a href="2">2</a></li>
					</ul>
				`),
			},
			want: "http://example.com/",
		},
		{
			name: "with the several base tags with the href attribute",
			args: args{
				data: []byte(`
					<base href="http://example.com/1/" />
					<base href="http://example.com/2/" />

					<ul>
						<li><a href="3">3</a></li>
						<li><a href="4">4</a></li>
					</ul>
				`),
			},
			want: "http://example.com/1/",
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			got := selectBaseTag(data.args.data)

			assert.Equal(test, data.want, got)
		})
	}
}
