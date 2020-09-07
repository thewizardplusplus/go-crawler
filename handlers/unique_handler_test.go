package handlers

import (
	"testing"

	mapset "github.com/deckarep/golang-set"
	"github.com/go-log/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	crawler "github.com/thewizardplusplus/go-crawler"
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

func TestUniqueHandler_HandleLink(test *testing.T) {
	type fields struct {
		sanitizeLink sanitizing.LinkSanitizing
		linkHandler  crawler.LinkHandler
		logger       log.Logger

		handledLinks mapset.Set
	}
	type args struct {
		sourceLink string
		link       string
	}

	for _, data := range []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	} {
		test.Run(data.name, func(test *testing.T) {
			handler := UniqueHandler{
				sanitizeLink: data.fields.sanitizeLink,
				linkHandler:  data.fields.linkHandler,
				logger:       data.fields.logger,

				handledLinks: data.fields.handledLinks,
			}
			handler.HandleLink(data.args.sourceLink, data.args.link)

			mock.AssertExpectationsForObjects(
				test,
				data.fields.linkHandler,
				data.fields.logger,
			)
		})
	}
}
