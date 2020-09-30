package extractors

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	crawler "github.com/thewizardplusplus/go-crawler"
)

func TestNewDelayedExtractor(test *testing.T) {
	sleeper := new(MockSleeper)
	linkExtractor := new(MockLinkExtractor)
	got := NewDelayedExtractor(100*time.Millisecond, sleeper.Sleep, linkExtractor)

	mock.AssertExpectationsForObjects(test, sleeper, linkExtractor)
	require.NotNil(test, got)
	assert.Equal(test, 100*time.Millisecond, got.minimalDelay)
	// don't use the reflect.Value.Pointer() method for this check; see details:
	// * https://golang.org/pkg/reflect/#Value.Pointer
	// * https://stackoverflow.com/a/9644797
	assert.NotNil(test, got.sleeper)
	assert.Equal(test, linkExtractor, got.linkExtractor)
}

func TestDelayedExtractor_ExtractLinks(test *testing.T) {
	type fields struct {
		minimalDelay  time.Duration
		sleeper       Sleeper
		linkExtractor crawler.LinkExtractor
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
		// TODO: Add test cases.
	} {
		test.Run(data.name, func(test *testing.T) {
			extractor := &DelayedExtractor{
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
