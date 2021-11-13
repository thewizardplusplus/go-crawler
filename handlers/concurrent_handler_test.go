package handlers

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/go-crawler/models"
)

func TestConcurrentHandler(test *testing.T) {
	type args struct {
		links []models.SourcedLink
	}

	for _, data := range []struct {
		name         string
		args         args
		startHandler func(ctx context.Context, handler ConcurrentHandler)
	}{
		{
			name: "with the Start() method",
			args: args{
				links: []models.SourcedLink{
					{
						SourceLink: "http://example.com/",
						Link:       "http://example.com/1",
					},
					{
						SourceLink: "http://example.com/",
						Link:       "http://example.com/2",
					},
				},
			},
			startHandler: func(ctx context.Context, handler ConcurrentHandler) {
				handler.Start(ctx)
			},
		},
		{
			name: "with the StartConcurrently() method",
			args: args{
				links: []models.SourcedLink{
					{
						SourceLink: "http://example.com/",
						Link:       "http://example.com/1",
					},
					{
						SourceLink: "http://example.com/",
						Link:       "http://example.com/2",
					},
				},
			},
			startHandler: func(ctx context.Context, handler ConcurrentHandler) {
				handler.StartConcurrently(ctx, 10)
			},
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			innerHandler := new(MockLinkHandler)
			for _, link := range data.args.links {
				innerHandler.On("HandleLink", context.Background(), link).Return().Times(1)
			}

			concurrentHandler := NewConcurrentHandler(1000, innerHandler)
			go data.startHandler(context.Background(), concurrentHandler)

			for _, link := range data.args.links {
				concurrentHandler.HandleLink(context.Background(), link)
			}
			concurrentHandler.Stop()

			mock.AssertExpectationsForObjects(test, innerHandler)
		})
	}
}
