package extractors

import (
	"context"
	"sync"
	"testing"
	"testing/iotest"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/thewizardplusplus/go-crawler/models"
)

func TestNewDelayingExtractor(test *testing.T) {
	sleeper := new(MockSleeper)
	linkExtractor := new(MockLinkExtractor)
	got := NewDelayingExtractor(100*time.Millisecond, sleeper.Sleep, linkExtractor)

	mock.AssertExpectationsForObjects(test, sleeper, linkExtractor)
	require.NotNil(test, got)
	assert.Equal(test, 100*time.Millisecond, got.minimalDelay)
	// don't use the reflect.Value.Pointer() method for this check; see details:
	// * https://golang.org/pkg/reflect/#Value.Pointer
	// * https://stackoverflow.com/a/9644797
	assert.NotNil(test, got.sleeper)
	assert.Equal(test, linkExtractor, got.linkExtractor)
}

func TestDelayingExtractor_ExtractLinks(test *testing.T) {
	type fields struct {
		minimalDelay  time.Duration
		sleeper       Sleeper
		linkExtractor models.LinkExtractor
	}
	type args struct {
		ctx      context.Context
		threadID int
		link     string
	}

	for _, data := range []struct {
		name                 string
		fields               fields
		initializeTimestamps func(timestamps *sync.Map)
		args                 args
		wantLinks            []string
		wantErr              assert.ErrorAssertionFunc
	}{
		{
			name: "success with an unknown thread ID",
			fields: fields{
				minimalDelay: 100 * time.Millisecond,
				sleeper:      new(MockSleeper),
				linkExtractor: func() models.LinkExtractor {
					extractor := new(MockLinkExtractor)
					extractor.
						On("ExtractLinks", context.Background(), 23, "http://example.com/").
						Return([]string{"http://example.com/1", "http://example.com/2"}, nil)

					return extractor
				}(),
			},
			initializeTimestamps: func(timestamps *sync.Map) {},
			args: args{
				ctx:      context.Background(),
				threadID: 23,
				link:     "http://example.com/",
			},
			wantLinks: []string{"http://example.com/1", "http://example.com/2"},
			wantErr:   assert.NoError,
		},
		{
			name: "success with a known thread ID",
			fields: fields{
				minimalDelay: 100 * time.Millisecond,
				sleeper: func() Sleeper {
					durationChecker := mock.MatchedBy(func(duration time.Duration) bool {
						return duration <= 100*time.Millisecond
					})

					sleeper := new(MockSleeper)
					sleeper.On("Sleep", durationChecker).Return()

					return sleeper
				}(),
				linkExtractor: func() models.LinkExtractor {
					extractor := new(MockLinkExtractor)
					extractor.
						On("ExtractLinks", context.Background(), 23, "http://example.com/").
						Return([]string{"http://example.com/1", "http://example.com/2"}, nil)

					return extractor
				}(),
			},
			initializeTimestamps: func(timestamps *sync.Map) {
				timestamps.Store(23, time.Now())
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
				minimalDelay: 100 * time.Millisecond,
				sleeper:      new(MockSleeper),
				linkExtractor: func() models.LinkExtractor {
					extractor := new(MockLinkExtractor)
					extractor.
						On("ExtractLinks", context.Background(), 23, "http://example.com/").
						Return(nil, iotest.ErrTimeout)

					return extractor
				}(),
			},
			initializeTimestamps: func(timestamps *sync.Map) {},
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
			extractor := &DelayingExtractor{
				minimalDelay:  data.fields.minimalDelay,
				sleeper:       data.fields.sleeper.Sleep,
				linkExtractor: data.fields.linkExtractor,
			}
			data.initializeTimestamps(&extractor.timestamps)

			gotLinks, gotErr := extractor.ExtractLinks(
				data.args.ctx,
				data.args.threadID,
				data.args.link,
			)

			mock.AssertExpectationsForObjects(
				test,
				data.fields.sleeper,
				data.fields.linkExtractor,
			)
			assert.Equal(test, data.wantLinks, gotLinks)
			data.wantErr(test, gotErr)
		})
	}
}
