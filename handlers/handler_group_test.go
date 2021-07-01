package handlers

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/go-crawler/models"
)

func TestHandlerGroup_HandleLink(test *testing.T) {
	type args struct {
		ctx  context.Context
		link models.SourcedLink
	}

	for _, data := range []struct {
		name     string
		handlers HandlerGroup
		args     args
	}{
		// TODO: Add test cases.
	} {
		test.Run(data.name, func(test *testing.T) {
			data.handlers.HandleLink(data.args.ctx, data.args.link)

			for _, handler := range data.handlers {
				mock.AssertExpectationsForObjects(test, handler)
			}
		})
	}
}
