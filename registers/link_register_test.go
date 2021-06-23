package registers

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
	urlutils "github.com/thewizardplusplus/go-crawler/url-utils"
)

func TestNewLinkRegister(test *testing.T) {
	logger := new(MockLogger)
	got := NewLinkRegister(urlutils.SanitizeLink, logger)

	mock.AssertExpectationsForObjects(test, logger)
	assert.Equal(test, urlutils.SanitizeLink, got.sanitizeLink)
	assert.Equal(test, logger, got.logger)
	assert.Equal(test, mapset.NewSet(), got.registeredLinks)
}

func TestLinkRegister_RegisterLink(test *testing.T) {
	type fields struct {
		sanitizeLink    urlutils.LinkSanitizing
		logger          log.Logger
		registeredLinks mapset.Set
	}
	type args struct {
		link string
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
				sanitizeLink: urlutils.DoNotSanitizeLink,
				logger:       new(MockLogger),
				registeredLinks: mapset.NewSet(
					"http://example.com/1",
					"http://example.com/2",
				),
			},
			args: args{
				link: "http://example.com/3",
			},
			want: assert.True,
		},
		{
			name: "success with a duplicate and without link sanitizing",
			fields: fields{
				sanitizeLink: urlutils.DoNotSanitizeLink,
				logger:       new(MockLogger),
				registeredLinks: mapset.NewSet(
					"http://example.com/1",
					"http://example.com/2",
				),
			},
			args: args{
				link: "http://example.com/2",
			},
			want: assert.False,
		},
		{
			name: "success with a duplicate and with link sanitizing",
			fields: fields{
				sanitizeLink: urlutils.SanitizeLink,
				logger:       new(MockLogger),
				registeredLinks: mapset.NewSet(
					"http://example.com/1",
					"http://example.com/2",
				),
			},
			args: args{
				link: "http://example.com/test/../2",
			},
			want: assert.False,
		},
		{
			name: "error",
			fields: fields{
				sanitizeLink: urlutils.SanitizeLink,
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
				registeredLinks: mapset.NewSet(
					"http://example.com/1",
					"http://example.com/2",
				),
			},
			args: args{
				link: ":",
			},
			want: assert.False,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			register := LinkRegister{
				sanitizeLink:    data.fields.sanitizeLink,
				logger:          data.fields.logger,
				registeredLinks: data.fields.registeredLinks,
			}
			got := register.RegisterLink(data.args.link)

			mock.AssertExpectationsForObjects(test, data.fields.logger)
			data.want(test, got)
		})
	}
}
