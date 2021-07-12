package extractors

import (
	"context"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
	"github.com/thewizardplusplus/go-crawler/models"
	urlutils "github.com/thewizardplusplus/go-crawler/url-utils"
)

func TestTrimmingExtractor_ExtractLinks(test *testing.T) {
	type fields struct {
		TrimLink      urlutils.LinkTrimming
		LinkExtractor models.LinkExtractor
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
			name: "success without links",
			fields: fields{
				TrimLink: urlutils.TrimLink,
				LinkExtractor: func() models.LinkExtractor {
					extractor := new(MockLinkExtractor)
					extractor.
						On("ExtractLinks", context.Background(), 23, "http://example.com/").
						Return(nil, nil)

					return extractor
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
		{
			name: "success with links and without spaces",
			fields: fields{
				TrimLink: urlutils.TrimLink,
				LinkExtractor: func() models.LinkExtractor {
					extractor := new(MockLinkExtractor)
					extractor.
						On("ExtractLinks", context.Background(), 23, "http://example.com/").
						Return([]string{"http://example.com/1", "http://example.com/2"}, nil)

					return extractor
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
			name: "success with links and spaces (DoNotTrimLink)",
			fields: fields{
				TrimLink: urlutils.DoNotTrimLink,
				LinkExtractor: func() models.LinkExtractor {
					extractor := new(MockLinkExtractor)
					extractor.
						On("ExtractLinks", context.Background(), 23, "http://example.com/").
						Return(
							[]string{"  http://example.com/1  ", "  http://example.com/2  "},
							nil,
						)

					return extractor
				}(),
			},
			args: args{
				ctx:      context.Background(),
				threadID: 23,
				link:     "http://example.com/",
			},
			wantLinks: []string{"  http://example.com/1  ", "  http://example.com/2  "},
			wantErr:   assert.NoError,
		},
		{
			name: "success with links and spaces (TrimLink)",
			fields: fields{
				TrimLink: urlutils.TrimLink,
				LinkExtractor: func() models.LinkExtractor {
					extractor := new(MockLinkExtractor)
					extractor.
						On("ExtractLinks", context.Background(), 23, "http://example.com/").
						Return(
							[]string{"  http://example.com/1  ", "  http://example.com/2  "},
							nil,
						)

					return extractor
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
				TrimLink: urlutils.TrimLink,
				LinkExtractor: func() models.LinkExtractor {
					extractor := new(MockLinkExtractor)
					extractor.
						On("ExtractLinks", context.Background(), 23, "http://example.com/").
						Return(nil, iotest.ErrTimeout)

					return extractor
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
			extractor := TrimmingExtractor{
				TrimLink:      data.fields.TrimLink,
				LinkExtractor: data.fields.LinkExtractor,
			}
			gotLinks, gotErr := extractor.ExtractLinks(
				data.args.ctx,
				data.args.threadID,
				data.args.link,
			)

			assert.Equal(test, data.wantLinks, gotLinks)
			data.wantErr(test, gotErr)
		})
	}
}
