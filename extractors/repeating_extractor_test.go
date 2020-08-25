package extractors

import (
	"context"
	"testing"
	"testing/iotest"
	"time"

	"github.com/go-log/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	crawler "github.com/thewizardplusplus/go-crawler"
)

func TestRepeatingExtractor_ExtractLinks(test *testing.T) {
	type fields struct {
		LinkExtractor crawler.LinkExtractor
		RepeatCount   int
		RepeatDelay   time.Duration
		Logger        log.Logger
	}
	type args struct {
		ctx  context.Context
		link string
	}

	for _, data := range []struct {
		name      string
		fields    fields
		args      args
		wantLinks []string
		wantErr   assert.ErrorAssertionFunc
	}{
		{
			name: "success on the first repeat",
			fields: fields{
				LinkExtractor: func() LinkExtractor {
					extractor := new(MockLinkExtractor)
					extractor.
						On("ExtractLinks", context.Background(), "http://example.com/").
						Return([]string{"http://example.com/1", "http://example.com/2"}, nil)

					return extractor
				}(),
				RepeatCount: 5,
				RepeatDelay: 100 * time.Millisecond,
				Logger:      new(MockLogger),
			},
			args: args{
				ctx:  context.Background(),
				link: "http://example.com/",
			},
			wantLinks: []string{"http://example.com/1", "http://example.com/2"},
			wantErr:   assert.NoError,
		},
		{
			name: "success on the last repeat",
			fields: fields{
				LinkExtractor: func() LinkExtractor {
					var repeat int

					extractor := new(MockLinkExtractor)
					extractor.
						On("ExtractLinks", context.Background(), "http://example.com/").
						Return(
							func(context.Context, string) []string {
								if repeat < 4 {
									return nil
								}

								return []string{"http://example.com/1", "http://example.com/2"}
							},
							func(context.Context, string) error {
								defer func() { repeat++ }()

								if repeat < 4 {
									return iotest.ErrTimeout
								}

								return nil
							},
						)

					return extractor
				}(),
				RepeatCount: 5,
				RepeatDelay: 100 * time.Millisecond,
				Logger: func() Logger {
					logger := new(MockLogger)
					for repeat := 0; repeat < 4; repeat++ {
						logger.
							On(
								"Logf",
								"unable to extract links (repeat #%d): %s",
								repeat,
								iotest.ErrTimeout,
							).
							Return()
					}

					return logger
				}(),
			},
			args: args{
				ctx:  context.Background(),
				link: "http://example.com/",
			},
			wantLinks: []string{"http://example.com/1", "http://example.com/2"},
			wantErr:   assert.NoError,
		},
		{
			name: "error",
			fields: fields{
				LinkExtractor: func() LinkExtractor {
					extractor := new(MockLinkExtractor)
					extractor.
						On("ExtractLinks", context.Background(), "http://example.com/").
						Return(nil, iotest.ErrTimeout)

					return extractor
				}(),
				RepeatCount: 5,
				RepeatDelay: 100 * time.Millisecond,
				Logger: func() Logger {
					logger := new(MockLogger)
					for repeat := 0; repeat < 4; repeat++ {
						logger.
							On(
								"Logf",
								"unable to extract links (repeat #%d): %s",
								repeat,
								iotest.ErrTimeout,
							).
							Return()
					}

					return logger
				}(),
			},
			args: args{
				ctx:  context.Background(),
				link: "http://example.com/",
			},
			wantLinks: nil,
			wantErr:   assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			extractor := RepeatingExtractor{
				LinkExtractor: data.fields.LinkExtractor,
				RepeatCount:   data.fields.RepeatCount,
				RepeatDelay:   data.fields.RepeatDelay,
				Logger:        data.fields.Logger,
			}
			gotLinks, gotErr := extractor.ExtractLinks(data.args.ctx, data.args.link)

			mock.AssertExpectationsForObjects(
				test,
				data.fields.LinkExtractor,
				data.fields.Logger,
			)
			assert.Equal(test, data.wantLinks, gotLinks)
			data.wantErr(test, gotErr)
		})
	}
}
