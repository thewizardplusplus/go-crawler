package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/thewizardplusplus/go-crawler/sanitizing"
)

func TestNewUniqueHandler(test *testing.T) {
	linkHandler := new(MockLinkHandler)
	logger := new(MockLogger)
	got := NewUniqueHandler(sanitizing.SanitizeLink, linkHandler, logger)

	mock.AssertExpectationsForObjects(test, linkHandler, logger)
	require.NotNil(test, got)
	assert.Equal(test, sanitizing.SanitizeLink, got.sanitizeLink)
	assert.Equal(test, linkHandler, got.linkHandler)
	assert.Equal(test, logger, got.logger)
	assert.NotNil(test, got.handledLinks)
}
