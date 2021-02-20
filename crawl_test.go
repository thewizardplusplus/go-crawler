package crawler

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/go-crawler/models"
)

func TestCrawl(test *testing.T) {
	type args struct {
		ctx               context.Context
		concurrencyConfig ConcurrencyConfig
		links             []string
		dependencies      CrawlDependencies
	}

	for _, data := range []struct {
		name string
		args args
	}{
		{
			name: "success with fewer links than the buffer size",
			args: args{
				ctx: context.Background(),
				concurrencyConfig: ConcurrencyConfig{
					ConcurrencyFactor: 10,
					BufferSize:        1000,
				},
				links: []string{"http://example.com/"},
				dependencies: CrawlDependencies{
					LinkExtractor: func() models.LinkExtractor {
						threadIDChecker := mock.MatchedBy(func(threadID int) bool {
							return threadID >= 0 && threadID < 10
						})

						extractor := new(MockLinkExtractor)
						extractor.
							On(
								"ExtractLinks",
								context.Background(),
								threadIDChecker,
								"http://example.com/",
							).
							Return([]string{"http://example.com/1", "http://example.com/2"}, nil)
						extractor.
							On(
								"ExtractLinks",
								context.Background(),
								threadIDChecker,
								"http://example.com/1",
							).
							Return(nil, nil)
						extractor.
							On(
								"ExtractLinks",
								context.Background(),
								threadIDChecker,
								"http://example.com/2",
							).
							Return(nil, nil)

						return extractor
					}(),
					LinkChecker: func() models.LinkChecker {
						checker := new(MockLinkChecker)
						checker.
							On("CheckLink", context.Background(), models.SourcedLink{
								SourceLink: "http://example.com/",
								Link:       "http://example.com/1",
							}).
							Return(true)
						checker.
							On("CheckLink", context.Background(), models.SourcedLink{
								SourceLink: "http://example.com/",
								Link:       "http://example.com/2",
							}).
							Return(true)

						return checker
					}(),
					LinkHandler: func() models.LinkHandler {
						handler := new(MockLinkHandler)
						handler.
							On("HandleLink", context.Background(), models.SourcedLink{
								SourceLink: "http://example.com/",
								Link:       "http://example.com/1",
							}).
							Return()
						handler.
							On("HandleLink", context.Background(), models.SourcedLink{
								SourceLink: "http://example.com/",
								Link:       "http://example.com/2",
							}).
							Return()

						return handler
					}(),
					Logger: new(MockLogger),
				},
			},
		},
		{
			name: "success without a buffer",
			args: args{
				ctx: context.Background(),
				concurrencyConfig: ConcurrencyConfig{
					ConcurrencyFactor: 10,
					BufferSize:        0,
				},
				links: []string{"http://example.com/"},
				dependencies: CrawlDependencies{
					LinkExtractor: func() models.LinkExtractor {
						threadIDChecker := mock.MatchedBy(func(threadID int) bool {
							return threadID >= 0 && threadID < 10
						})

						extractor := new(MockLinkExtractor)
						extractor.
							On(
								"ExtractLinks",
								context.Background(),
								threadIDChecker,
								"http://example.com/",
							).
							Return([]string{"http://example.com/1", "http://example.com/2"}, nil)
						extractor.
							On(
								"ExtractLinks",
								context.Background(),
								threadIDChecker,
								"http://example.com/1",
							).
							Return(nil, nil)
						extractor.
							On(
								"ExtractLinks",
								context.Background(),
								threadIDChecker,
								"http://example.com/2",
							).
							Return(nil, nil)

						return extractor
					}(),
					LinkChecker: func() models.LinkChecker {
						checker := new(MockLinkChecker)
						checker.
							On("CheckLink", context.Background(), models.SourcedLink{
								SourceLink: "http://example.com/",
								Link:       "http://example.com/1",
							}).
							Return(true)
						checker.
							On("CheckLink", context.Background(), models.SourcedLink{
								SourceLink: "http://example.com/",
								Link:       "http://example.com/2",
							}).
							Return(true)

						return checker
					}(),
					LinkHandler: func() models.LinkHandler {
						handler := new(MockLinkHandler)
						handler.
							On("HandleLink", context.Background(), models.SourcedLink{
								SourceLink: "http://example.com/",
								Link:       "http://example.com/1",
							}).
							Return()
						handler.
							On("HandleLink", context.Background(), models.SourcedLink{
								SourceLink: "http://example.com/",
								Link:       "http://example.com/2",
							}).
							Return()

						return handler
					}(),
					Logger: new(MockLogger),
				},
			},
		},
	} {
		test.Run(data.name, func(t *testing.T) {
			Crawl(
				data.args.ctx,
				data.args.concurrencyConfig,
				data.args.links,
				data.args.dependencies,
			)

			mock.AssertExpectationsForObjects(
				test,
				data.args.dependencies.LinkExtractor,
				data.args.dependencies.LinkChecker,
				data.args.dependencies.LinkHandler,
				data.args.dependencies.Logger,
			)
		})
	}
}

func TestCrawlByConcurrentHandler(test *testing.T) {
	type args struct {
		ctx                      context.Context
		concurrencyConfig        ConcurrencyConfig
		handlerConcurrencyConfig ConcurrencyConfig
		links                    []string
		dependencies             CrawlDependencies
	}

	for _, data := range []struct {
		name string
		args args
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				concurrencyConfig: ConcurrencyConfig{
					ConcurrencyFactor: 10,
					BufferSize:        1000,
				},
				handlerConcurrencyConfig: ConcurrencyConfig{
					ConcurrencyFactor: 10,
					BufferSize:        1000,
				},
				links: []string{"http://example.com/"},
				dependencies: CrawlDependencies{
					LinkExtractor: func() models.LinkExtractor {
						threadIDChecker := mock.MatchedBy(func(threadID int) bool {
							return threadID >= 0 && threadID < 10
						})

						extractor := new(MockLinkExtractor)
						extractor.
							On(
								"ExtractLinks",
								context.Background(),
								threadIDChecker,
								"http://example.com/",
							).
							Return([]string{"http://example.com/1", "http://example.com/2"}, nil)
						extractor.
							On(
								"ExtractLinks",
								context.Background(),
								threadIDChecker,
								"http://example.com/1",
							).
							Return(nil, nil)
						extractor.
							On(
								"ExtractLinks",
								context.Background(),
								threadIDChecker,
								"http://example.com/2",
							).
							Return(nil, nil)

						return extractor
					}(),
					LinkChecker: func() models.LinkChecker {
						checker := new(MockLinkChecker)
						checker.
							On("CheckLink", context.Background(), models.SourcedLink{
								SourceLink: "http://example.com/",
								Link:       "http://example.com/1",
							}).
							Return(true)
						checker.
							On("CheckLink", context.Background(), models.SourcedLink{
								SourceLink: "http://example.com/",
								Link:       "http://example.com/2",
							}).
							Return(true)

						return checker
					}(),
					LinkHandler: func() models.LinkHandler {
						handler := new(MockLinkHandler)
						handler.
							On("HandleLink", context.Background(), models.SourcedLink{
								SourceLink: "http://example.com/",
								Link:       "http://example.com/1",
							}).
							Return()
						handler.
							On("HandleLink", context.Background(), models.SourcedLink{
								SourceLink: "http://example.com/",
								Link:       "http://example.com/2",
							}).
							Return()

						return handler
					}(),
					Logger: new(MockLogger),
				},
			},
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			CrawlByConcurrentHandler(
				data.args.ctx,
				data.args.concurrencyConfig,
				data.args.handlerConcurrencyConfig,
				data.args.links,
				data.args.dependencies,
			)

			mock.AssertExpectationsForObjects(
				test,
				data.args.dependencies.LinkExtractor,
				data.args.dependencies.LinkChecker,
				data.args.dependencies.LinkHandler,
				data.args.dependencies.Logger,
			)
		})
	}
}
