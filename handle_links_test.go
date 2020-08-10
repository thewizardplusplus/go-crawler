package crawler

import (
	"context"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandleLink(test *testing.T) {
	type args struct {
		ctx          context.Context
		link         string
		dependencies Dependencies
	}

	for _, data := range []struct {
		name      string
		args      args
		wantLinks []string
	}{
		{
			name: "success with all correct links",
			args: args{
				ctx:  context.Background(),
				link: "http://example.com/",
				dependencies: Dependencies{
					Waiter: func() Waiter {
						waiter := new(MockWaiter)
						waiter.On("Add", 1).Return().Times(2)
						waiter.On("Done").Return().Times(1)

						return waiter
					}(),
					LinkExtractor: func() LinkExtractor {
						extractor := new(MockLinkExtractor)
						extractor.
							On("ExtractLinks", context.Background(), "http://example.com/").
							Return([]string{"http://example.com/1", "http://example.com/2"}, nil)

						return extractor
					}(),
					LinkChecker: func() LinkChecker {
						checker := new(MockLinkChecker)
						checker.
							On("CheckLink", "http://example.com/", "http://example.com/1").
							Return(true)
						checker.
							On("CheckLink", "http://example.com/", "http://example.com/2").
							Return(true)

						return checker
					}(),
					LinkHandler: func() LinkHandler {
						linkHandler := new(MockLinkHandler)
						linkHandler.On("HandleLink", "http://example.com/").Return()

						return linkHandler
					}(),
					ErrorHandler: new(MockErrorHandler),
				},
			},
			wantLinks: []string{"http://example.com/1", "http://example.com/2"},
		},
		{
			name: "success with some correct links",
			args: args{
				ctx:  context.Background(),
				link: "http://example.com/",
				dependencies: Dependencies{
					Waiter: func() Waiter {
						waiter := new(MockWaiter)
						waiter.On("Add", 1).Return().Times(1)
						waiter.On("Done").Return().Times(1)

						return waiter
					}(),
					LinkExtractor: func() LinkExtractor {
						extractor := new(MockLinkExtractor)
						extractor.
							On("ExtractLinks", context.Background(), "http://example.com/").
							Return([]string{"http://example.com/1", "http://example.com/2"}, nil)

						return extractor
					}(),
					LinkChecker: func() LinkChecker {
						checker := new(MockLinkChecker)
						checker.
							On("CheckLink", "http://example.com/", "http://example.com/1").
							Return(false)
						checker.
							On("CheckLink", "http://example.com/", "http://example.com/2").
							Return(true)

						return checker
					}(),
					LinkHandler: func() LinkHandler {
						linkHandler := new(MockLinkHandler)
						linkHandler.On("HandleLink", "http://example.com/").Return()

						return linkHandler
					}(),
					ErrorHandler: new(MockErrorHandler),
				},
			},
			wantLinks: []string{"http://example.com/2"},
		},
		{
			name: "error",
			args: args{
				ctx:  context.Background(),
				link: "http://example.com/",
				dependencies: Dependencies{
					Waiter: func() Waiter {
						waiter := new(MockWaiter)
						waiter.On("Done").Return().Times(1)

						return waiter
					}(),
					LinkExtractor: func() LinkExtractor {
						extractor := new(MockLinkExtractor)
						extractor.
							On("ExtractLinks", context.Background(), "http://example.com/").
							Return(nil, iotest.ErrTimeout)

						return extractor
					}(),
					LinkChecker: new(MockLinkChecker),
					LinkHandler: func() LinkHandler {
						linkHandler := new(MockLinkHandler)
						linkHandler.On("HandleLink", "http://example.com/").Return()

						return linkHandler
					}(),
					ErrorHandler: func() ErrorHandler {
						errorHandler := new(MockErrorHandler)
						errorHandler.On("HandleError", iotest.ErrTimeout).Return()

						return errorHandler
					}(),
				},
			},
			wantLinks: nil,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			gotLinks := HandleLink(data.args.ctx, data.args.link, data.args.dependencies)

			mock.AssertExpectationsForObjects(
				test,
				data.args.dependencies.Waiter,
				data.args.dependencies.LinkExtractor,
				data.args.dependencies.LinkChecker,
				data.args.dependencies.LinkHandler,
				data.args.dependencies.ErrorHandler,
			)
			assert.Equal(test, data.wantLinks, gotLinks)
		})
	}
}
