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
		Sleeper       SleeperInterface
	}
	type args struct {
		ctx      context.Context
		threadID int
		link     string
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
						On("ExtractLinks", context.Background(), 23, "http://example.com/").
						Return([]string{"http://example.com/1", "http://example.com/2"}, nil).
						Times(1)

					return extractor
				}(),
				RepeatCount: 5,
				RepeatDelay: 100 * time.Millisecond,
				Logger:      new(MockLogger),
				Sleeper:     new(MockSleeperInterface),
			},
			args: args{
				ctx:      context.Background(),
				threadID: 23,
				link:     "http://example.com/",
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
						On("ExtractLinks", context.Background(), 23, "http://example.com/").
						Return(
							func(context.Context, int, string) []string {
								if repeat < 4 {
									return nil
								}

								return []string{"http://example.com/1", "http://example.com/2"}
							},
							func(context.Context, int, string) error {
								defer func() { repeat++ }()

								if repeat < 4 {
									return iotest.ErrTimeout
								}

								return nil
							},
						).
						Times(5)

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
							Return().
							Times(1)
					}

					return logger
				}(),
				Sleeper: func() SleeperInterface {
					sleeper := new(MockSleeperInterface)
					sleeper.On("Sleep", 100*time.Millisecond).Return().Times(4)

					return sleeper
				}(),
			},
			args: args{
				ctx:      context.Background(),
				threadID: 23,
				link:     "http://example.com/",
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
						On("ExtractLinks", context.Background(), 23, "http://example.com/").
						Return(nil, iotest.ErrTimeout).
						Times(5)

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
							Return().
							Times(1)
					}

					return logger
				}(),
				Sleeper: func() SleeperInterface {
					sleeper := new(MockSleeperInterface)
					sleeper.On("Sleep", 100*time.Millisecond).Return().Times(4)

					return sleeper
				}(),
			},
			args: args{
				ctx:      context.Background(),
				threadID: 23,
				link:     "http://example.com/",
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
				Sleeper:       data.fields.Sleeper.Sleep,
			}
			gotLinks, gotErr := extractor.ExtractLinks(
				data.args.ctx,
				data.args.threadID,
				data.args.link,
			)

			mock.AssertExpectationsForObjects(
				test,
				data.fields.LinkExtractor,
				data.fields.Logger,
				data.fields.Sleeper,
			)
			assert.Equal(test, data.wantLinks, gotLinks)
			data.wantErr(test, gotErr)
		})
	}
}
