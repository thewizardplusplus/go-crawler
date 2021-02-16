package crawler

import (
	"context"
	"sync"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/go-crawler/models"
	syncutils "github.com/thewizardplusplus/go-sync-utils"
)

func TestHandleLinksConcurrently(test *testing.T) {
	type args struct {
		ctx               context.Context
		concurrencyFactor int
		links             chan string
		dependencies      HandleLinkDependencies
	}

	for _, data := range []struct {
		name string
		args args
	}{
		{
			name: "success",
			args: args{
				ctx:               context.Background(),
				concurrencyFactor: 10,
				links: func() chan string {
					links := make(chan string, 1)
					links <- "http://example.com/"

					return links
				}(),
				dependencies: HandleLinkDependencies{
					CrawlDependencies: CrawlDependencies{
						LinkExtractor: func() LinkExtractor {
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
						LinkChecker: func() LinkChecker {
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
						LinkHandler: func() LinkHandler {
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
					Waiter: func() syncutils.WaitGroup {
						waiter := new(MockWaiter)
						waiter.On("Add", 1).Return().Times(3)
						waiter.On("Done").Return().Times(3)
						waiter.On("Wait").Return().Times(1)

						return waiter
					}(),
				},
			},
		},
	} {
		test.Run(data.name, func(t *testing.T) {
			waiter := data.args.dependencies.Waiter
			synchronousWaiter := syncutils.MultiWaitGroup{waiter, new(sync.WaitGroup)}
			synchronousWaiter.Add(len(data.args.links))

			data.args.dependencies.Waiter = synchronousWaiter

			HandleLinksConcurrently(
				data.args.ctx,
				data.args.concurrencyFactor,
				data.args.links,
				data.args.dependencies,
			)
			synchronousWaiter.Wait()

			mock.AssertExpectationsForObjects(
				test,
				waiter,
				data.args.dependencies.LinkExtractor,
				data.args.dependencies.LinkChecker,
				data.args.dependencies.LinkHandler,
				data.args.dependencies.Logger,
			)
		})
	}
}

func TestHandleLinks(test *testing.T) {
	type args struct {
		ctx          context.Context
		threadID     int
		links        chan string
		dependencies HandleLinkDependencies
	}

	for _, data := range []struct {
		name string
		args args
	}{
		{
			name: "success",
			args: args{
				ctx:      context.Background(),
				threadID: 23,
				links: func() chan string {
					links := make(chan string, 1)
					links <- "http://example.com/"

					return links
				}(),
				dependencies: HandleLinkDependencies{
					CrawlDependencies: CrawlDependencies{
						LinkExtractor: func() LinkExtractor {
							extractor := new(MockLinkExtractor)
							extractor.
								On("ExtractLinks", context.Background(), 23, "http://example.com/").
								Return([]string{"http://example.com/1", "http://example.com/2"}, nil)
							extractor.
								On("ExtractLinks", context.Background(), 23, "http://example.com/1").
								Return(nil, nil)
							extractor.
								On("ExtractLinks", context.Background(), 23, "http://example.com/2").
								Return(nil, nil)

							return extractor
						}(),
						LinkChecker: func() LinkChecker {
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
						LinkHandler: func() LinkHandler {
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
					Waiter: func() syncutils.WaitGroup {
						waiter := new(MockWaiter)
						waiter.On("Add", 1).Return().Times(3)
						waiter.On("Done").Return().Times(3)
						waiter.On("Wait").Return().Times(1)

						return waiter
					}(),
				},
			},
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			waiter := data.args.dependencies.Waiter
			synchronousWaiter := syncutils.MultiWaitGroup{waiter, new(sync.WaitGroup)}
			synchronousWaiter.Add(len(data.args.links))

			data.args.dependencies.Waiter = synchronousWaiter

			go HandleLinks(
				data.args.ctx,
				data.args.threadID,
				data.args.links,
				data.args.dependencies,
			)

			synchronousWaiter.Wait()

			mock.AssertExpectationsForObjects(
				test,
				waiter,
				data.args.dependencies.LinkExtractor,
				data.args.dependencies.LinkChecker,
				data.args.dependencies.LinkHandler,
				data.args.dependencies.Logger,
			)
		})
	}
}

func TestHandleLink(test *testing.T) {
	type args struct {
		ctx          context.Context
		threadID     int
		link         string
		dependencies HandleLinkDependencies
	}

	for _, data := range []struct {
		name      string
		args      args
		wantLinks []string
	}{
		{
			name: "success with all correct links",
			args: args{
				ctx:      context.Background(),
				threadID: 23,
				link:     "http://example.com/",
				dependencies: HandleLinkDependencies{
					CrawlDependencies: CrawlDependencies{
						LinkExtractor: func() LinkExtractor {
							extractor := new(MockLinkExtractor)
							extractor.
								On("ExtractLinks", context.Background(), 23, "http://example.com/").
								Return([]string{"http://example.com/1", "http://example.com/2"}, nil)

							return extractor
						}(),
						LinkChecker: func() LinkChecker {
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
						LinkHandler: func() LinkHandler {
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
					Waiter: func() syncutils.WaitGroup {
						waiter := new(MockWaiter)
						waiter.On("Add", 1).Return().Times(2)
						waiter.On("Done").Return().Times(1)

						return waiter
					}(),
				},
			},
			wantLinks: []string{"http://example.com/1", "http://example.com/2"},
		},
		{
			name: "success with some correct links",
			args: args{
				ctx:      context.Background(),
				threadID: 23,
				link:     "http://example.com/",
				dependencies: HandleLinkDependencies{
					CrawlDependencies: CrawlDependencies{
						LinkExtractor: func() LinkExtractor {
							extractor := new(MockLinkExtractor)
							extractor.
								On("ExtractLinks", context.Background(), 23, "http://example.com/").
								Return([]string{"http://example.com/1", "http://example.com/2"}, nil)

							return extractor
						}(),
						LinkChecker: func() LinkChecker {
							checker := new(MockLinkChecker)
							checker.
								On("CheckLink", context.Background(), models.SourcedLink{
									SourceLink: "http://example.com/",
									Link:       "http://example.com/1",
								}).
								Return(false)
							checker.
								On("CheckLink", context.Background(), models.SourcedLink{
									SourceLink: "http://example.com/",
									Link:       "http://example.com/2",
								}).
								Return(true)

							return checker
						}(),
						LinkHandler: func() LinkHandler {
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
					Waiter: func() syncutils.WaitGroup {
						waiter := new(MockWaiter)
						waiter.On("Add", 1).Return().Times(1)
						waiter.On("Done").Return().Times(1)

						return waiter
					}(),
				},
			},
			wantLinks: []string{"http://example.com/2"},
		},
		{
			name: "error",
			args: args{
				ctx:      context.Background(),
				threadID: 23,
				link:     "http://example.com/",
				dependencies: HandleLinkDependencies{
					CrawlDependencies: CrawlDependencies{
						LinkExtractor: func() LinkExtractor {
							extractor := new(MockLinkExtractor)
							extractor.
								On("ExtractLinks", context.Background(), 23, "http://example.com/").
								Return(nil, iotest.ErrTimeout)

							return extractor
						}(),
						LinkChecker: new(MockLinkChecker),
						LinkHandler: new(MockLinkHandler),
						Logger: func() Logger {
							logger := new(MockLogger)
							logger.
								On("Logf", "unable to extract links: %s", iotest.ErrTimeout).
								Return()

							return logger
						}(),
					},
					Waiter: func() syncutils.WaitGroup {
						waiter := new(MockWaiter)
						waiter.On("Done").Return().Times(1)

						return waiter
					}(),
				},
			},
			wantLinks: nil,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			gotLinks := HandleLink(
				data.args.ctx,
				data.args.threadID,
				data.args.link,
				data.args.dependencies,
			)

			mock.AssertExpectationsForObjects(
				test,
				data.args.dependencies.Waiter,
				data.args.dependencies.LinkExtractor,
				data.args.dependencies.LinkChecker,
				data.args.dependencies.LinkHandler,
				data.args.dependencies.Logger,
			)
			assert.Equal(test, data.wantLinks, gotLinks)
		})
	}
}
