package sitemap

import (
	"context"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/go-crawler/registers"
)

func TestGeneratorGroup_GenerateLinks(test *testing.T) {
	type args struct {
		ctx      context.Context
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
				baseLink: "http://example.com/test",
			},
			wantSitemapLinks: nil,
			wantErr:          assert.NoError,
		},
		{
			name: "without failed generatings",
			generators: GeneratorGroup{
				func() registers.LinkGenerator {
					ctxMatcher := mock.MatchedBy(func(context.Context) bool { return true })

					generator := new(MockLinkGenerator)
					generator.
						On("GenerateLinks", ctxMatcher, "http://example.com/test").
						Return(
							[]string{
								"http://example.com/sitemap_1.xml",
								"http://example.com/sitemap_2.xml",
							},
							nil,
						)

					return generator
				}(),
				func() registers.LinkGenerator {
					ctxMatcher := mock.MatchedBy(func(context.Context) bool { return true })

					generator := new(MockLinkGenerator)
					generator.
						On("GenerateLinks", ctxMatcher, "http://example.com/test").
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
				func() registers.LinkGenerator {
					ctxMatcher := mock.MatchedBy(func(context.Context) bool { return true })

					generator := new(MockLinkGenerator)
					generator.
						On("GenerateLinks", ctxMatcher, "http://example.com/test").
						Return(nil, iotest.ErrTimeout)

					return generator
				}(),
				func() registers.LinkGenerator {
					ctxMatcher := mock.MatchedBy(func(context.Context) bool { return true })

					generator := new(MockLinkGenerator)
					generator.
						On("GenerateLinks", ctxMatcher, "http://example.com/test").
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
				baseLink: "http://example.com/test",
			},
			wantSitemapLinks: nil,
			wantErr:          assert.Error,
		},
		{
			name: "with all failed generatings",
			generators: GeneratorGroup{
				func() registers.LinkGenerator {
					ctxMatcher := mock.MatchedBy(func(context.Context) bool { return true })

					generator := new(MockLinkGenerator)
					generator.
						On("GenerateLinks", ctxMatcher, "http://example.com/test").
						Return(nil, iotest.ErrTimeout)

					return generator
				}(),
				func() registers.LinkGenerator {
					ctxMatcher := mock.MatchedBy(func(context.Context) bool { return true })

					generator := new(MockLinkGenerator)
					generator.
						On("GenerateLinks", ctxMatcher, "http://example.com/test").
						Return(nil, iotest.ErrTimeout)

					return generator
				}(),
			},
			args: args{
				ctx:      context.Background(),
				baseLink: "http://example.com/test",
			},
			wantSitemapLinks: nil,
			wantErr:          assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			gotSitemapLinks, gotErr :=
				data.generators.GenerateLinks(data.args.ctx, data.args.baseLink)

			assert.Equal(test, data.wantSitemapLinks, gotSitemapLinks)
			data.wantErr(test, gotErr)
		})
	}
}
