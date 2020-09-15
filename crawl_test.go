package crawler

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
)

func TestCrawl(test *testing.T) {
	type args struct {
		ctx               context.Context
		concurrencyFactor int
		links             []string
		dependencies      Dependencies
	}

	for _, data := range []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	} {
		test.Run(data.name, func(t *testing.T) {
			Crawl(
				data.args.ctx,
				data.args.concurrencyFactor,
				data.args.links,
				data.args.dependencies,
			)

			mock.AssertExpectationsForObjects(
				test,
				data.args.dependencies.LinkExtractor,
				data.args.dependencies.LinkChecker,
				data.args.dependencies.LinkHandler,
				data.args.dependencies.Logger,
			)
		})
	}
}
