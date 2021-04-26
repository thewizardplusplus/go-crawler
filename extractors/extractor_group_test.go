package extractors

import (
	"context"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/go-crawler/models"
)

func TestExtractorGroup_ExtractLinks(test *testing.T) {
	type args struct {
		ctx      context.Context
		threadID int
		link     string
	}

	for _, data := range []struct {
		name       string
		extractors ExtractorGroup
		args       args
		wantLinks  []string
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name:       "empty",
			extractors: nil,
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
			extractors: ExtractorGroup{
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
			name: "with a failed extracting",
			extractors: ExtractorGroup{
				func() models.LinkExtractor {
					extractor := new(MockLinkExtractor)
					extractor.
						On("ExtractLinks", context.Background(), 23, "http://example.com/").
						Return(nil, iotest.ErrTimeout)

					return extractor
				}(),
				new(MockLinkExtractor),
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
			gotLinks, gotErr := data.extractors.ExtractLinks(
				data.args.ctx,
				data.args.threadID,
				data.args.link,
			)

			for _, extractor := range data.extractors {
				mock.AssertExpectationsForObjects(test, extractor)
			}
			assert.Equal(test, data.wantLinks, gotLinks)
			data.wantErr(test, gotErr)
		})
	}
}
