package transformers

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	urlutils "github.com/thewizardplusplus/go-crawler/url-utils"
)

func TestTrimmingTransformer_TransformLinks(test *testing.T) {
	type fields struct {
		TrimLink urlutils.LinkTrimming
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
				TrimLink: urlutils.TrimLink,
			},
			args: args{
				links:           nil,
				response:        nil,
				responseContent: nil,
			},
			wantLinks: nil,
			wantErr:   assert.NoError,
		},
		{
			name: "success with links and without spaces",
			fields: fields{
				TrimLink: urlutils.TrimLink,
			},
			args: args{
				links:           []string{"http://example.com/1", "http://example.com/2"},
				response:        nil,
				responseContent: nil,
			},
			wantLinks: []string{"http://example.com/1", "http://example.com/2"},
			wantErr:   assert.NoError,
		},
		{
			name: "success with links and spaces (DoNotTrimLink)",
			fields: fields{
				TrimLink: urlutils.DoNotTrimLink,
			},
			args: args{
				links: []string{
					"  http://example.com/1  ",
					"  http://example.com/2  ",
				},
				response:        nil,
				responseContent: nil,
			},
			wantLinks: []string{"  http://example.com/1  ", "  http://example.com/2  "},
			wantErr:   assert.NoError,
		},
		{
			name: "success with links and spaces (TrimLink)",
			fields: fields{
				TrimLink: urlutils.TrimLink,
			},
			args: args{
				links: []string{
					"  http://example.com/1  ",
					"  http://example.com/2  ",
				},
				response:        nil,
				responseContent: nil,
			},
			wantLinks: []string{"http://example.com/1", "http://example.com/2"},
			wantErr:   assert.NoError,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			transformer := TrimmingTransformer{
				TrimLink: data.fields.TrimLink,
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
