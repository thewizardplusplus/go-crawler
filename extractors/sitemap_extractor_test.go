package extractors

import (
	"context"
	"testing"
	"time"

	"github.com/go-log/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/go-crawler/registers"
)

func TestSitemapExtractor_ExtractLinks(test *testing.T) {
	type fields struct {
		loadingInterval time.Duration
		linkGenerator   registers.LinkGenerator
		logger          log.Logger
		sleeper         Sleeper
		linkLoader      LinkLoader
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
		// TODO: Add test cases.
	} {
		test.Run(data.name, func(test *testing.T) {
			register := registers.NewSitemapRegister(
				data.fields.loadingInterval,
				data.fields.linkGenerator,
				data.fields.logger,
				data.fields.sleeper.Sleep,
				data.fields.linkLoader.LoadLink,
			)
			extractor := SitemapExtractor{
				SitemapRegister: register,
				Logger:          data.fields.logger,
			}
			gotLinks, gotErr := extractor.ExtractLinks(
				data.args.ctx,
				data.args.threadID,
				data.args.link,
			)

			mock.AssertExpectationsForObjects(
				test,
				data.fields.linkGenerator,
				data.fields.logger,
				data.fields.sleeper,
				data.fields.linkLoader,
			)
			assert.Equal(test, data.wantLinks, gotLinks)
			data.wantErr(test, gotErr)
		})
	}
}
