package sitemap

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleGenerator_ExtractLinks(test *testing.T) {
	type args struct {
		ctx      context.Context
		threadID int
		baseLink string
	}

	for _, data := range []struct {
		name             string
		args             args
		wantSitemapLinks []string
		wantErr          assert.ErrorAssertionFunc
	}{
		{
			name: "success with a path only",
			args: args{
				ctx:      context.Background(),
				threadID: 23,
				baseLink: "http://example.com/test",
			},
			wantSitemapLinks: []string{"http://example.com/sitemap.xml"},
			wantErr:          assert.NoError,
		},
		{
			name: "success with an HTTPS scheme",
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
			args: args{
				ctx:      context.Background(),
				threadID: 23,
				baseLink: "http://example.com/test#fragment",
			},
			wantSitemapLinks: []string{"http://example.com/sitemap.xml"},
			wantErr:          assert.NoError,
		},
		{
			name: "error",
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
			var generator SimpleGenerator
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
