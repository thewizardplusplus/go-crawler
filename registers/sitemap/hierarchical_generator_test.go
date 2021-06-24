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
		MaximalDepth int
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
			name: "success",
			fields: fields{
				SanitizeLink: urlutils.DoNotSanitizeLink,
				MaximalDepth: -1,
			},
			args: args{
				ctx:      context.Background(),
				threadID: 23,
				baseLink: "http://example.com/one/two/three/test",
			},
			wantSitemapLinks: []string{
				"http://example.com/sitemap.xml",
				"http://example.com/one/sitemap.xml",
				"http://example.com/one/two/sitemap.xml",
				"http://example.com/one/two/three/sitemap.xml",
			},
			wantErr: assert.NoError,
		},
		{
			name: "success with a depth limit",
			fields: fields{
				SanitizeLink: urlutils.DoNotSanitizeLink,
				MaximalDepth: 2,
			},
			args: args{
				ctx:      context.Background(),
				threadID: 23,
				baseLink: "http://example.com/one/two/three/test",
			},
			wantSitemapLinks: []string{
				"http://example.com/sitemap.xml",
				"http://example.com/one/sitemap.xml",
				"http://example.com/one/two/sitemap.xml",
			},
			wantErr: assert.NoError,
		},
		{
			name: "success with sanitizing",
			fields: fields{
				SanitizeLink: urlutils.SanitizeLink,
				MaximalDepth: -1,
			},
			args: args{
				ctx:      context.Background(),
				threadID: 23,
				baseLink: "http://example.com/one/two//three/../test",
			},
			wantSitemapLinks: []string{
				"http://example.com/sitemap.xml",
				"http://example.com/one/sitemap.xml",
				"http://example.com/one/two/sitemap.xml",
			},
			wantErr: assert.NoError,
		},
		{
			name: "error",
			fields: fields{
				SanitizeLink: urlutils.DoNotSanitizeLink,
				MaximalDepth: -1,
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
				MaximalDepth: data.fields.MaximalDepth,
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
