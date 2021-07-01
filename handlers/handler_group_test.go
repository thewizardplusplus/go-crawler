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
		{
			name:     "empty",
			handlers: nil,
			args: args{
				ctx: context.Background(),
				link: models.SourcedLink{
					SourceLink: "http://example.com/",
					Link:       "http://example.com/test",
				},
			},
		},
		{
			name: "non-empty",
			handlers: HandlerGroup{
				func() models.LinkHandler {
					handler := new(MockLinkHandler)
					handler.
						On("HandleLink", context.Background(), models.SourcedLink{
							SourceLink: "http://example.com/",
							Link:       "http://example.com/test",
						}).
						Return()

					return handler
				}(),
				func() models.LinkHandler {
					handler := new(MockLinkHandler)
					handler.
						On("HandleLink", context.Background(), models.SourcedLink{
							SourceLink: "http://example.com/",
							Link:       "http://example.com/test",
						}).
						Return()

					return handler
				}(),
			},
			args: args{
				ctx: context.Background(),
				link: models.SourcedLink{
					SourceLink: "http://example.com/",
					Link:       "http://example.com/test",
				},
			},
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			data.handlers.HandleLink(data.args.ctx, data.args.link)

			for _, handler := range data.handlers {
				mock.AssertExpectationsForObjects(test, handler)
			}
		})
	}
}
