package handlers

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	crawler "github.com/thewizardplusplus/go-crawler"
	"github.com/thewizardplusplus/go-crawler/registers"
	"github.com/thewizardplusplus/go-crawler/sanitizing"
)

func TestUniqueHandler_HandleLink(test *testing.T) {
	type fields struct {
		LinkRegister registers.LinkRegister
		LinkHandler  crawler.LinkHandler
	}
	type args struct {
		ctx  context.Context
		link crawler.SourcedLink
	}

	for _, data := range []struct {
		name             string
		fields           fields
		args             args
		wantLinkRegister registers.LinkRegister
	}{
		{
			name: "without a duplicate",
			fields: fields{
				LinkRegister: registers.NewLinkRegister(sanitizing.DoNotSanitizeLink, nil),
				LinkHandler: func() LinkHandler {
					handler := new(MockLinkHandler)
					handler.
						On("HandleLink", context.Background(), crawler.SourcedLink{
							SourceLink: "http://example.com/",
							Link:       "http://example.com/test",
						}).
						Return()

					return handler
				}(),
			},
			args: args{
				ctx: context.Background(),
				link: crawler.SourcedLink{
					SourceLink: "http://example.com/",
					Link:       "http://example.com/test",
				},
			},
			wantLinkRegister: func() registers.LinkRegister {
				linkRegister := registers.NewLinkRegister(sanitizing.DoNotSanitizeLink, nil)
				linkRegister.RegisterLink("http://example.com/test")

				return linkRegister
			}(),
		},
		{
			name: "with a duplicate",
			fields: fields{
				LinkRegister: func() registers.LinkRegister {
					linkRegister :=
						registers.NewLinkRegister(sanitizing.DoNotSanitizeLink, nil)
					linkRegister.RegisterLink("http://example.com/test")

					return linkRegister
				}(),
				LinkHandler: new(MockLinkHandler),
			},
			args: args{
				ctx: context.Background(),
				link: crawler.SourcedLink{
					SourceLink: "http://example.com/",
					Link:       "http://example.com/test",
				},
			},
			wantLinkRegister: func() registers.LinkRegister {
				linkRegister := registers.NewLinkRegister(sanitizing.DoNotSanitizeLink, nil)
				linkRegister.RegisterLink("http://example.com/test")

				return linkRegister
			}(),
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			handler := UniqueHandler{
				LinkRegister: data.fields.LinkRegister,
				LinkHandler:  data.fields.LinkHandler,
			}
			handler.HandleLink(data.args.ctx, data.args.link)

			mock.AssertExpectationsForObjects(test, data.fields.LinkHandler)
			assert.Equal(test, data.wantLinkRegister, handler.LinkRegister)
		})
	}
}
