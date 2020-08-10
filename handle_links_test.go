package crawler

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandleLink(test *testing.T) {
	type args struct {
		ctx          context.Context
		link         string
		dependencies Dependencies
	}

	for _, data := range []struct {
		name      string
		args      args
		wantLinks []string
	}{
		// TODO: Add test cases.
	} {
		test.Run(data.name, func(test *testing.T) {
			gotLinks := HandleLink(data.args.ctx, data.args.link, data.args.dependencies)

			mock.AssertExpectationsForObjects(
				test,
				data.args.dependencies.Waiter,
				data.args.dependencies.LinkExtractor,
				data.args.dependencies.LinkChecker,
				data.args.dependencies.LinkHandler,
				data.args.dependencies.ErrorHandler,
			)
			assert.Equal(test, data.wantLinks, gotLinks)
		})
	}
}
