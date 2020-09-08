package handlers

import (
	stderrors "errors"
	"net/url"
	"reflect"
	"testing"

	mapset "github.com/deckarep/golang-set"
	"github.com/go-log/log"
	"github.com/pkg/errors"
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
		{
			name: "success without a duplicate",
			fields: fields{
				sanitizeLink: sanitizing.DoNotSanitizeLink,
				linkHandler: func() crawler.LinkHandler {
					handler := new(MockLinkHandler)
					handler.
						On("HandleLink", "http://example.com/", "http://example.com/3").
						Return()

					return handler
				}(),
				logger: new(MockLogger),

				handledLinks: mapset.NewSet("http://example.com/1", "http://example.com/2"),
			},
			args: args{
				sourceLink: "http://example.com/",
				link:       "http://example.com/3",
			},
		},
		{
			name: "success with a duplicate and without link sanitizing",
			fields: fields{
				sanitizeLink: sanitizing.DoNotSanitizeLink,
				linkHandler:  new(MockLinkHandler),
				logger:       new(MockLogger),

				handledLinks: mapset.NewSet("http://example.com/1", "http://example.com/2"),
			},
			args: args{
				sourceLink: "http://example.com/",
				link:       "http://example.com/2",
			},
		},
		{
			name: "success with a duplicate and with link sanitizing",
			fields: fields{
				sanitizeLink: sanitizing.SanitizeLink,
				linkHandler:  new(MockLinkHandler),
				logger:       new(MockLogger),

				handledLinks: mapset.NewSet("http://example.com/1", "http://example.com/2"),
			},
			args: args{
				sourceLink: "http://example.com/",
				link:       "http://example.com/test/../2",
			},
		},
		{
			name: "error",
			fields: fields{
				sanitizeLink: sanitizing.SanitizeLink,
				linkHandler:  new(MockLinkHandler),
				logger: func() Logger {
					err := stderrors.New("missing protocol scheme")
					urlErr := &url.Error{Op: "parse", URL: ":", Err: err}

					logger := new(MockLogger)
					logger.
						On(
							"Logf",
							"unable to sanitize the link: %s",
							mock.MatchedBy(func(err error) bool {
								unwrappedErr := errors.Cause(err)
								return reflect.DeepEqual(unwrappedErr, urlErr)
							}),
						).
						Return()

					return logger
				}(),

				handledLinks: mapset.NewSet("http://example.com/1", "http://example.com/2"),
			},
			args: args{
				sourceLink: "http://example.com/",
				link:       ":",
			},
		},
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
