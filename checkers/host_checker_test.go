package checkers

import (
	"context"
	"testing"

	"github.com/go-log/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/go-crawler/models"
	urlutils "github.com/thewizardplusplus/go-crawler/url-utils"
)

func TestHostChecker_CheckLink(test *testing.T) {
	type fields struct {
		ComparisonResult urlutils.ComparisonResult
		Logger           log.Logger
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
			name: "success with different hosts (false)",
			fields: fields{
				ComparisonResult: urlutils.Same,
				Logger:           new(MockLogger),
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
			name: "success with different hosts (true)",
			fields: fields{
				ComparisonResult: urlutils.Different,
				Logger:           new(MockLogger),
			},
			args: args{
				ctx: context.Background(),
				link: models.SourcedLink{
					SourceLink: "http://example1.com/",
					Link:       "http://example2.com/test",
				},
			},
			want: assert.True,
		},
		{
			name: "success with same hosts (false)",
			fields: fields{
				ComparisonResult: urlutils.Different,
				Logger:           new(MockLogger),
			},
			args: args{
				ctx: context.Background(),
				link: models.SourcedLink{
					SourceLink: "http://example.com/",
					Link:       "http://example.com/test",
				},
			},
			want: assert.False,
		},
		{
			name: "success with same hosts (true)",
			fields: fields{
				ComparisonResult: urlutils.Same,
				Logger:           new(MockLogger),
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
				ComparisonResult: urlutils.Same,
				Logger: func() Logger {
					logger := new(MockLogger)
					logger.
						On(
							"Logf",
							"%s: unable to compare link hosts: %s",
							"host checking",
							mock.MatchedBy(func(err error) bool {
								wantErrMessage := `unable to parse link ":": ` +
									`parse :: ` +
									"missing protocol scheme"
								return err.Error() == wantErrMessage
							}),
						).
						Return()

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
				ComparisonResult: urlutils.Same,
				Logger: func() Logger {
					logger := new(MockLogger)
					logger.
						On(
							"Logf",
							"%s: unable to compare link hosts: %s",
							"host checking",
							mock.MatchedBy(func(err error) bool {
								wantErrMessage := `unable to parse link ":": ` +
									`parse :: ` +
									"missing protocol scheme"
								return err.Error() == wantErrMessage
							}),
						).
						Return()

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
				ComparisonResult: data.fields.ComparisonResult,
				Logger:           data.fields.Logger,
			}
			got := checker.CheckLink(data.args.ctx, data.args.link)

			mock.AssertExpectationsForObjects(test, data.fields.Logger)
			data.want(test, got)
		})
	}
}
