package handlers

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	crawler "github.com/thewizardplusplus/go-crawler"
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
	link := crawler.SourcedLink{
		SourceLink: "http://example.com/",
		Link:       "http://example.com/test",
	}

	links := make(chan crawler.SourcedLink, 1)
	handler := ConcurrentHandler{links: links}
	handler.HandleLink(context.Background(), link)

	gotLink := <-handler.links
	assert.Equal(test, link, gotLink)
}

func TestConcurrentHandler_Run(test *testing.T) {
	links := []crawler.SourcedLink{
		{
			SourceLink: "http://example.com/",
			Link:       "http://example.com/1",
		},
		{
			SourceLink: "http://example.com/",
			Link:       "http://example.com/2",
		},
	}

	innerHandler := new(MockLinkHandler)
	for _, link := range links {
		innerHandler.On("HandleLink", context.Background(), link).Return()
	}

	linkChannel := make(chan crawler.SourcedLink, len(links))
	for _, link := range links {
		linkChannel <- link
	}
	close(linkChannel)

	handler := ConcurrentHandler{linkHandler: innerHandler, links: linkChannel}
	handler.Run(context.Background())

	mock.AssertExpectationsForObjects(test, innerHandler)
}

func TestConcurrentHandler_RunConcurrently(test *testing.T) {
	links := []crawler.SourcedLink{
		{
			SourceLink: "http://example.com/",
			Link:       "http://example.com/1",
		},
		{
			SourceLink: "http://example.com/",
			Link:       "http://example.com/2",
		},
	}

	innerHandler := new(MockLinkHandler)
	for _, link := range links {
		innerHandler.On("HandleLink", context.Background(), link).Return()
	}

	linkChannel := make(chan crawler.SourcedLink, len(links))
	for _, link := range links {
		linkChannel <- link
	}
	close(linkChannel)

	handler := ConcurrentHandler{linkHandler: innerHandler, links: linkChannel}
	handler.RunConcurrently(context.Background(), 10)

	mock.AssertExpectationsForObjects(test, innerHandler)
}
