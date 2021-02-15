package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
