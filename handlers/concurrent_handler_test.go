package handlers

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/go-crawler/models"
)

func TestNewConcurrentHandler(test *testing.T) {
	innerHandler := new(MockLinkHandler)
	handler := NewConcurrentHandler(1000, innerHandler)

	mock.AssertExpectationsForObjects(test, innerHandler)
	assert.Equal(test, innerHandler, handler.linkHandler)
	assert.NotNil(test, handler.links)
	assert.Len(test, handler.links, 0)
	assert.Equal(test, 1000, cap(handler.links))
}

func TestConcurrentHandler_HandleLink(test *testing.T) {
	link := models.SourcedLink{
		SourceLink: "http://example.com/",
		Link:       "http://example.com/test",
	}

	links := make(chan models.SourcedLink, 1)
	handler := ConcurrentHandler{links: links}
	handler.HandleLink(context.Background(), link)

	gotLink := <-handler.links
	assert.Equal(test, link, gotLink)
}

func TestConcurrentHandler_running(test *testing.T) {
	for _, data := range []struct {
		name       string
		links      []models.SourcedLink
		runHandler func(ctx context.Context, handler ConcurrentHandler)
	}{
		{
			name: "with the Run() method",
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
			runHandler: func(ctx context.Context, handler ConcurrentHandler) {
				handler.Run(ctx)
			},
		},
		{
			name: "with the RunConcurrently() method",
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
			runHandler: func(ctx context.Context, handler ConcurrentHandler) {
				handler.RunConcurrently(ctx, 10)
			},
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			innerHandler := new(MockLinkHandler)
			for _, link := range data.links {
				innerHandler.On("HandleLink", context.Background(), link).Return()
			}

			linkChannel := make(chan models.SourcedLink, len(data.links))
			for _, link := range data.links {
				linkChannel <- link
			}
			close(linkChannel)

			handler := ConcurrentHandler{linkHandler: innerHandler, links: linkChannel}
			data.runHandler(context.Background(), handler)

			mock.AssertExpectationsForObjects(test, innerHandler)
		})
	}
}

func TestConcurrentHandler_Stop(test *testing.T) {
	links := make(chan models.SourcedLink)
	handler := ConcurrentHandler{links: links}
	handler.Stop()

	isNotClosed := true
	select {
	case _, isNotClosed = <-handler.links:
	default: // to prevent blocking
	}

	assert.False(test, isNotClosed)
}
