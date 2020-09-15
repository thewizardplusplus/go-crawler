package crawler

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
)

func TestCrawl(test *testing.T) {
	type args struct {
		ctx               context.Context
		concurrencyFactor int
		bufferSize        int
		links             []string
		dependencies      Dependencies
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
				bufferSize:        1000,
				links:             []string{"http://example.com/"},
				dependencies: Dependencies{
					Waiter: nil,
					LinkExtractor: func() LinkExtractor {
						extractor := new(MockLinkExtractor)
						extractor.
							On("ExtractLinks", context.Background(), "http://example.com/").
							Return([]string{"http://example.com/1", "http://example.com/2"}, nil)
						extractor.
							On("ExtractLinks", context.Background(), "http://example.com/1").
							Return(nil, nil)
						extractor.
							On("ExtractLinks", context.Background(), "http://example.com/2").
							Return(nil, nil)

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
						handler := new(MockLinkHandler)
						handler.
							On("HandleLink", "http://example.com/", "http://example.com/1").
							Return()
						handler.
							On("HandleLink", "http://example.com/", "http://example.com/2").
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
				data.args.concurrencyFactor,
				data.args.bufferSize,
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
