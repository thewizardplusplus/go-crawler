package extractors

import (
	"context"
	"testing"

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
		// TODO: Add test cases.
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
