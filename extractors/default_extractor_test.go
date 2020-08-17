package extractors

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	htmlselector "github.com/thewizardplusplus/go-html-selector"
)

func TestDefaultExtractor_ExtractLinks(test *testing.T) {
	type fields struct {
		HTTPClient HTTPClient
		Filters    htmlselector.OptimizedFilterGroup
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
	}{} {
		test.Run(data.name, func(test *testing.T) {
			extractor := DefaultExtractor{
				HTTPClient: data.fields.HTTPClient,
				Filters:    data.fields.Filters,
			}
			gotLinks, gotErr := extractor.ExtractLinks(data.args.ctx, data.args.link)

			mock.AssertExpectationsForObjects(test, data.fields.HTTPClient)
			assert.Equal(test, data.wantLinks, gotLinks)
			data.wantErr(test, gotErr)
		})
	}
}
