package register

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/go-crawler/sanitizing"
)

func TestNewLinkRegister(test *testing.T) {
	logger := new(MockLogger)
	got := NewLinkRegister(sanitizing.SanitizeLink, logger)

	mock.AssertExpectationsForObjects(test, logger)
	assert.Equal(test, sanitizing.SanitizeLink, got.sanitizeLink)
	assert.Equal(test, logger, got.logger)
	assert.NotNil(test, got.registeredLinks)
}
