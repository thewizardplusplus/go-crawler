package extractors

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
		// TODO: Add test cases.
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
