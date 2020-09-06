package checkers

import (
	"errors"
	"net/url"
	"testing"

	mapset "github.com/deckarep/golang-set"
	"github.com/go-log/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/thewizardplusplus/go-crawler/sanitizing"
)

func TestNewDuplicateChecker(test *testing.T) {
	logger := new(MockLogger)
	got := NewDuplicateChecker(sanitizing.SanitizeLink, logger)

	mock.AssertExpectationsForObjects(test, logger)
	require.NotNil(test, got)
	assert.Equal(test, sanitizing.SanitizeLink, got.sanitizeLink)
	assert.Equal(test, logger, got.logger)
	assert.NotNil(test, got.checkedLinks)
}

func TestDuplicateChecker_CheckLink(test *testing.T) {
	type fields struct {
		sanitizeLink sanitizing.LinkSanitizing
		logger       log.Logger

		checkedLinks mapset.Set
	}
	type args struct {
		sourceLink string
		link       string
	}

	for _, data := range []struct {
		name   string
		fields fields
		args   args
		want   assert.BoolAssertionFunc
	}{
		{
			name: "success without a duplicate",
			fields: fields{
				sanitizeLink: sanitizing.DoNotSanitizeLink,
				logger:       new(MockLogger),

				checkedLinks: mapset.NewSet("http://example.com/1", "http://example.com/2"),
			},
			args: args{
				sourceLink: "http://example.com/",
				link:       "http://example.com/3",
			},
			want: assert.True,
		},
		{
			name: "success with a duplicate and without link sanitizing",
			fields: fields{
				sanitizeLink: sanitizing.DoNotSanitizeLink,
				logger:       new(MockLogger),

				checkedLinks: mapset.NewSet("http://example.com/1", "http://example.com/2"),
			},
			args: args{
				sourceLink: "http://example.com/",
				link:       "http://example.com/2",
			},
			want: assert.False,
		},
		{
			name: "success with a duplicate and with link sanitizing",
			fields: fields{
				sanitizeLink: sanitizing.SanitizeLink,
				logger:       new(MockLogger),

				checkedLinks: mapset.NewSet("http://example.com/1", "http://example.com/2"),
			},
			args: args{
				sourceLink: "http://example.com/",
				link:       "http://example.com/test/../2",
			},
			want: assert.False,
		},
		{
			name: "error",
			fields: fields{
				sanitizeLink: sanitizing.SanitizeLink,
				logger: func() Logger {
					err := errors.New("missing protocol scheme")
					urlErr := &url.Error{Op: "parse", URL: ":", Err: err}

					logger := new(MockLogger)
					logger.On("Logf", "unable to parse the link: %s", urlErr).Return()

					return logger
				}(),

				checkedLinks: mapset.NewSet("http://example.com/1", "http://example.com/2"),
			},
			args: args{
				sourceLink: "http://example.com/",
				link:       ":",
			},
			want: assert.False,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			checker := DuplicateChecker{
				sanitizeLink: data.fields.sanitizeLink,
				logger:       data.fields.logger,

				checkedLinks: data.fields.checkedLinks,
			}
			got := checker.CheckLink(data.args.sourceLink, data.args.link)

			mock.AssertExpectationsForObjects(test, data.fields.logger)
			data.want(test, got)
		})
	}
}
