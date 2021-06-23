package sitemap

import (
	"context"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/go-crawler/models"
)

func TestGeneratorGroup_ExtractLinks(test *testing.T) {
	type args struct {
		ctx      context.Context
		threadID int
		baseLink string
	}

	for _, data := range []struct {
		name             string
		generators       GeneratorGroup
		args             args
		wantSitemapLinks []string
		wantErr          assert.ErrorAssertionFunc
	}{
		{
			name:       "empty",
			generators: nil,
			args: args{
				ctx:      context.Background(),
				threadID: 23,
				baseLink: "http://example.com/test",
			},
			wantSitemapLinks: nil,
			wantErr:          assert.NoError,
		},
		{
			name: "without failed generatings",
			generators: GeneratorGroup{
				func() models.LinkExtractor {
					ctxMatcher := mock.MatchedBy(func(context.Context) bool { return true })

					generator := new(MockLinkExtractor)
					generator.
						On("ExtractLinks", ctxMatcher, 23, "http://example.com/test").
						Return(
							[]string{
								"http://example.com/sitemap_1.xml",
								"http://example.com/sitemap_2.xml",
							},
							nil,
						)

					return generator
				}(),
				func() models.LinkExtractor {
					ctxMatcher := mock.MatchedBy(func(context.Context) bool { return true })

					generator := new(MockLinkExtractor)
					generator.
						On("ExtractLinks", ctxMatcher, 23, "http://example.com/test").
						Return(
							[]string{
								"http://example.com/sitemap_3.xml",
								"http://example.com/sitemap_4.xml",
							},
							nil,
						)

					return generator
				}(),
			},
			args: args{
				ctx:      context.Background(),
				threadID: 23,
				baseLink: "http://example.com/test",
			},
			wantSitemapLinks: []string{
				"http://example.com/sitemap_1.xml",
				"http://example.com/sitemap_2.xml",
				"http://example.com/sitemap_3.xml",
				"http://example.com/sitemap_4.xml",
			},
			wantErr: assert.NoError,
		},
		{
			name: "with some failed generatings",
			generators: GeneratorGroup{
				func() models.LinkExtractor {
					ctxMatcher := mock.MatchedBy(func(context.Context) bool { return true })

					generator := new(MockLinkExtractor)
					generator.
						On("ExtractLinks", ctxMatcher, 23, "http://example.com/test").
						Return(nil, iotest.ErrTimeout)

					return generator
				}(),
				func() models.LinkExtractor {
					ctxMatcher := mock.MatchedBy(func(context.Context) bool { return true })

					generator := new(MockLinkExtractor)
					generator.
						On("ExtractLinks", ctxMatcher, 23, "http://example.com/test").
						Return(
							[]string{
								"http://example.com/sitemap_3.xml",
								"http://example.com/sitemap_4.xml",
							},
							nil,
						)

					return generator
				}(),
			},
			args: args{
				ctx:      context.Background(),
				threadID: 23,
				baseLink: "http://example.com/test",
			},
			wantSitemapLinks: nil,
			wantErr:          assert.Error,
		},
		{
			name: "with all failed generatings",
			generators: GeneratorGroup{
				func() models.LinkExtractor {
					ctxMatcher := mock.MatchedBy(func(context.Context) bool { return true })

					generator := new(MockLinkExtractor)
					generator.
						On("ExtractLinks", ctxMatcher, 23, "http://example.com/test").
						Return(nil, iotest.ErrTimeout)

					return generator
				}(),
				func() models.LinkExtractor {
					ctxMatcher := mock.MatchedBy(func(context.Context) bool { return true })

					generator := new(MockLinkExtractor)
					generator.
						On("ExtractLinks", ctxMatcher, 23, "http://example.com/test").
						Return(nil, iotest.ErrTimeout)

					return generator
				}(),
			},
			args: args{
				ctx:      context.Background(),
				threadID: 23,
				baseLink: "http://example.com/test",
			},
			wantSitemapLinks: nil,
			wantErr:          assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			gotSitemapLinks, gotErr := data.generators.ExtractLinks(
				data.args.ctx,
				data.args.threadID,
				data.args.baseLink,
			)

			assert.Equal(test, data.wantSitemapLinks, gotSitemapLinks)
			data.wantErr(test, gotErr)
		})
	}
}
