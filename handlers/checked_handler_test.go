package handlers

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	crawler "github.com/thewizardplusplus/go-crawler"
)

func TestCheckedHandler_HandleLink(test *testing.T) {
	type fields struct {
		LinkChecker crawler.LinkChecker
		LinkHandler crawler.LinkHandler
	}
	type args struct {
		ctx  context.Context
		link crawler.SourcedLink
	}

	for _, data := range []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	} {
		test.Run(data.name, func(test *testing.T) {
			handler := CheckedHandler{
				LinkChecker: data.fields.LinkChecker,
				LinkHandler: data.fields.LinkHandler,
			}
			handler.HandleLink(data.args.ctx, data.args.link)

			mock.AssertExpectationsForObjects(
				test,
				data.fields.LinkChecker,
				data.fields.LinkHandler,
			)
		})
	}
}
