package checkers

import (
	"context"
	"errors"
	"net/url"
	"testing"

	"github.com/go-log/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/go-crawler/models"
)

func TestHostChecker_CheckLink(test *testing.T) {
	type fields struct {
		Logger log.Logger
	}
	type args struct {
		ctx  context.Context
		link models.SourcedLink
	}

	for _, data := range []struct {
		name   string
		fields fields
		args   args
		want   assert.BoolAssertionFunc
	}{
		{
			name: "success with different hosts",
			fields: fields{
				Logger: new(MockLogger),
			},
			args: args{
				ctx: context.Background(),
				link: models.SourcedLink{
					SourceLink: "http://example1.com/",
					Link:       "http://example2.com/test",
				},
			},
			want: assert.False,
		},
		{
			name: "success with same hosts",
			fields: fields{
				Logger: new(MockLogger),
			},
			args: args{
				ctx: context.Background(),
				link: models.SourcedLink{
					SourceLink: "http://example.com/",
					Link:       "http://example.com/test",
				},
			},
			want: assert.True,
		},
		{
			name: "error with the parent link",
			fields: fields{
				Logger: func() Logger {
					err := errors.New("missing protocol scheme")
					urlErr := &url.Error{Op: "parse", URL: ":", Err: err}

					logger := new(MockLogger)
					logger.On("Logf", "unable to parse the parent link: %s", urlErr).Return()

					return logger
				}(),
			},
			args: args{
				ctx: context.Background(),
				link: models.SourcedLink{
					SourceLink: ":",
					Link:       "http://example.com/test",
				},
			},
			want: assert.False,
		},
		{
			name: "error with the link",
			fields: fields{
				Logger: func() Logger {
					err := errors.New("missing protocol scheme")
					urlErr := &url.Error{Op: "parse", URL: ":", Err: err}

					logger := new(MockLogger)
					logger.On("Logf", "unable to parse the link: %s", urlErr).Return()

					return logger
				}(),
			},
			args: args{
				ctx: context.Background(),
				link: models.SourcedLink{
					SourceLink: "http://example.com/",
					Link:       ":",
				},
			},
			want: assert.False,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			checker := HostChecker{
				Logger: data.fields.Logger,
			}
			got := checker.CheckLink(data.args.ctx, data.args.link)

			mock.AssertExpectationsForObjects(test, data.fields.Logger)
			data.want(test, got)
		})
	}
}
