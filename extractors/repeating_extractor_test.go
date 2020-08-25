package extractors

import (
	"context"
	"testing"
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
		// TODO: Add test cases.
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
