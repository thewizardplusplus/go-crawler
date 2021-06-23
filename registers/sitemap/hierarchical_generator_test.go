package sitemap

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	urlutils "github.com/thewizardplusplus/go-crawler/url-utils"
)

func TestHierarchicalGenerator_ExtractLinks(test *testing.T) {
	type fields struct {
		SanitizeLink urlutils.LinkSanitizing
	}
	type args struct {
		ctx      context.Context
		threadID int
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
			name: "success without a trailing slash",
			fields: fields{
				SanitizeLink: urlutils.DoNotSanitizeLink,
			},
			args: args{
				ctx:      context.Background(),
				threadID: 23,
				baseLink: "http://example.com/test",
			},
			wantSitemapLinks: []string{"http://example.com/sitemap.xml"},
			wantErr:          assert.NoError,
		},
		{
			name: "success with a trailing slash",
			fields: fields{
				SanitizeLink: urlutils.DoNotSanitizeLink,
			},
			args: args{
				ctx:      context.Background(),
				threadID: 23,
				baseLink: "http://example.com/test/",
			},
			wantSitemapLinks: []string{
				"http://example.com/sitemap.xml",
				"http://example.com/test/sitemap.xml",
			},
			wantErr: assert.NoError,
		},
		{
			name: "success with a long path",
			fields: fields{
				SanitizeLink: urlutils.DoNotSanitizeLink,
			},
			args: args{
				ctx:      context.Background(),
				threadID: 23,
				baseLink: "http://example.com/one/two/test",
			},
			wantSitemapLinks: []string{
				"http://example.com/sitemap.xml",
				"http://example.com/one/sitemap.xml",
				"http://example.com/one/two/sitemap.xml",
			},
			wantErr: assert.NoError,
		},
		{
			name: "success with an HTTPS scheme",
			fields: fields{
				SanitizeLink: urlutils.DoNotSanitizeLink,
			},
			args: args{
				ctx:      context.Background(),
				threadID: 23,
				baseLink: "https://example.com/test",
			},
			wantSitemapLinks: []string{"https://example.com/sitemap.xml"},
			wantErr:          assert.NoError,
		},
		{
			name: "success with an user",
			fields: fields{
				SanitizeLink: urlutils.DoNotSanitizeLink,
			},
			args: args{
				ctx:      context.Background(),
				threadID: 23,
				baseLink: "http://username:password@example.com/test",
			},
			wantSitemapLinks: []string{
				"http://username:password@example.com/sitemap.xml",
			},
			wantErr: assert.NoError,
		},
		{
			name: "success with a query",
			fields: fields{
				SanitizeLink: urlutils.DoNotSanitizeLink,
			},
			args: args{
				ctx:      context.Background(),
				threadID: 23,
				baseLink: "http://example.com/test?key=value",
			},
			wantSitemapLinks: []string{"http://example.com/sitemap.xml"},
			wantErr:          assert.NoError,
		},
		{
			name: "success with a fragment",
			fields: fields{
				SanitizeLink: urlutils.DoNotSanitizeLink,
			},
			args: args{
				ctx:      context.Background(),
				threadID: 23,
				baseLink: "http://example.com/test#fragment",
			},
			wantSitemapLinks: []string{"http://example.com/sitemap.xml"},
			wantErr:          assert.NoError,
		},
		{
			name: "success with sanitizing",
			fields: fields{
				SanitizeLink: urlutils.SanitizeLink,
			},
			args: args{
				ctx:      context.Background(),
				threadID: 23,
				baseLink: "http://example.com/one/two/../test",
			},
			wantSitemapLinks: []string{
				"http://example.com/sitemap.xml",
				"http://example.com/one/sitemap.xml",
			},
			wantErr: assert.NoError,
		},
		{
			name: "error with link sanitizing",
			fields: fields{
				SanitizeLink: urlutils.SanitizeLink,
			},
			args: args{
				ctx:      context.Background(),
				threadID: 23,
				baseLink: ":",
			},
			wantSitemapLinks: nil,
			wantErr:          assert.Error,
		},
		{
			name: "error with link parsing",
			fields: fields{
				SanitizeLink: urlutils.DoNotSanitizeLink,
			},
			args: args{
				ctx:      context.Background(),
				threadID: 23,
				baseLink: ":",
			},
			wantSitemapLinks: nil,
			wantErr:          assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			generator := HierarchicalGenerator{
				SanitizeLink: data.fields.SanitizeLink,
			}
			gotSitemapLinks, gotErr := generator.ExtractLinks(
				data.args.ctx,
				data.args.threadID,
				data.args.baseLink,
			)

			assert.Equal(test, data.wantSitemapLinks, gotSitemapLinks)
			data.wantErr(test, gotErr)
		})
	}
}
