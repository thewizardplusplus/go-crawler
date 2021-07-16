package extractors

import (
	"context"
	"testing"
	"testing/iotest"

	"github.com/go-log/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/go-crawler/models"
)

func TestExtractorGroup_ExtractLinks(test *testing.T) {
	type fields struct {
		LinkExtractors []models.LinkExtractor
		Logger         log.Logger
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
			name: "empty",
			fields: fields{
				LinkExtractors: nil,
				Logger:         new(MockLogger),
			},
			args: args{
				ctx:      context.Background(),
				threadID: 23,
				link:     "http://example.com/",
			},
			wantLinks: nil,
			wantErr:   assert.NoError,
		},
		{
			name: "without failed extractings",
			fields: fields{
				LinkExtractors: []models.LinkExtractor{
					func() models.LinkExtractor {
						extractor := new(MockLinkExtractor)
						extractor.
							On("ExtractLinks", context.Background(), 23, "http://example.com/").
							Return([]string{"http://example.com/1", "http://example.com/2"}, nil)

						return extractor
					}(),
					func() models.LinkExtractor {
						extractor := new(MockLinkExtractor)
						extractor.
							On("ExtractLinks", context.Background(), 23, "http://example.com/").
							Return([]string{"http://example.com/3", "http://example.com/4"}, nil)

						return extractor
					}(),
				},
				Logger: new(MockLogger),
			},
			args: args{
				ctx:      context.Background(),
				threadID: 23,
				link:     "http://example.com/",
			},
			wantLinks: []string{
				"http://example.com/1",
				"http://example.com/2",
				"http://example.com/3",
				"http://example.com/4",
			},
			wantErr: assert.NoError,
		},
		{
			name: "with some failed extractings",
			fields: fields{
				LinkExtractors: []models.LinkExtractor{
					func() models.LinkExtractor {
						extractor := new(MockLinkExtractor)
						extractor.
							On("ExtractLinks", context.Background(), 23, "http://example.com/").
							Return(nil, iotest.ErrTimeout)

						return extractor
					}(),
					func() models.LinkExtractor {
						extractor := new(MockLinkExtractor)
						extractor.
							On("ExtractLinks", context.Background(), 23, "http://example.com/").
							Return([]string{"http://example.com/3", "http://example.com/4"}, nil)

						return extractor
					}(),
				},
				Logger: func() Logger {
					logger := new(MockLogger)
					logger.
						On(
							"Logf",
							"unable to extract links for link %q via extractor #%d: %s",
							"http://example.com/",
							0,
							iotest.ErrTimeout,
						).
						Return()

					return logger
				}(),
			},
			args: args{
				ctx:      context.Background(),
				threadID: 23,
				link:     "http://example.com/",
			},
			wantLinks: []string{"http://example.com/3", "http://example.com/4"},
			wantErr:   assert.NoError,
		},
		{
			name: "with all failed extractings",
			fields: fields{
				LinkExtractors: []models.LinkExtractor{
					func() models.LinkExtractor {
						extractor := new(MockLinkExtractor)
						extractor.
							On("ExtractLinks", context.Background(), 23, "http://example.com/").
							Return(nil, iotest.ErrTimeout)

						return extractor
					}(),
					func() models.LinkExtractor {
						extractor := new(MockLinkExtractor)
						extractor.
							On("ExtractLinks", context.Background(), 23, "http://example.com/").
							Return(nil, iotest.ErrTimeout)

						return extractor
					}(),
				},
				Logger: func() Logger {
					logger := new(MockLogger)
					for _, index := range []int{0, 1} {
						logger.
							On(
								"Logf",
								"unable to extract links for link %q via extractor #%d: %s",
								"http://example.com/",
								index,
								iotest.ErrTimeout,
							).
							Return()
					}

					return logger
				}(),
			},
			args: args{
				ctx:      context.Background(),
				threadID: 23,
				link:     "http://example.com/",
			},
			wantLinks: nil,
			wantErr:   assert.NoError,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			extractors := ExtractorGroup{
				LinkExtractors: data.fields.LinkExtractors,
				Logger:         data.fields.Logger,
			}
			gotLinks, gotErr := extractors.ExtractLinks(
				data.args.ctx,
				data.args.threadID,
				data.args.link,
			)

			for _, extractor := range data.fields.LinkExtractors {
				mock.AssertExpectationsForObjects(test, extractor)
			}
			mock.AssertExpectationsForObjects(test, data.fields.Logger)
			assert.Equal(test, data.wantLinks, gotLinks)
			data.wantErr(test, gotErr)
		})
	}
}
