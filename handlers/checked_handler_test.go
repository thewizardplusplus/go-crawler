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
		{
			name: "with a passed checking",
			fields: fields{
				LinkChecker: func() crawler.LinkChecker {
					checker := new(MockLinkChecker)
					checker.
						On("CheckLink", context.Background(), crawler.SourcedLink{
							SourceLink: "http://example.com/",
							Link:       "http://example.com/test",
						}).
						Return(true)

					return checker
				}(),
				LinkHandler: func() crawler.LinkHandler {
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
		},
		{
			name: "with a not passed checking",
			fields: fields{
				LinkChecker: func() crawler.LinkChecker {
					checker := new(MockLinkChecker)
					checker.
						On("CheckLink", context.Background(), crawler.SourcedLink{
							SourceLink: "http://example.com/",
							Link:       "http://example.com/test",
						}).
						Return(false)

					return checker
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
		},
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
