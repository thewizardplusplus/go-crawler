package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	crawler "github.com/thewizardplusplus/go-crawler"
	"github.com/thewizardplusplus/go-crawler/register"
	"github.com/thewizardplusplus/go-crawler/sanitizing"
)

func TestNewUniqueHandler(test *testing.T) {
	linkRegister := register.NewLinkRegister(sanitizing.DoNotSanitizeLink, nil)
	linkHandler := new(MockLinkHandler)
	got := NewUniqueHandler(linkRegister, linkHandler)

	mock.AssertExpectationsForObjects(test, linkHandler)
	require.NotNil(test, got)
	assert.Equal(test, linkRegister, got.linkRegister)
	assert.Equal(test, linkHandler, got.linkHandler)
}

func TestUniqueHandler_HandleLink(test *testing.T) {
	type fields struct {
		linkRegister register.LinkRegister
		linkHandler  crawler.LinkHandler
	}
	type args struct {
		sourceLink string
		link       string
	}

	for _, data := range []struct {
		name             string
		fields           fields
		args             args
		wantLinkRegister register.LinkRegister
	}{
		{
			name: "without a duplicate",
			fields: fields{
				linkRegister: register.NewLinkRegister(sanitizing.DoNotSanitizeLink, nil),
				linkHandler: func() LinkHandler {
					handler := new(MockLinkHandler)
					handler.
						On("HandleLink", "http://example.com/", "http://example.com/test").
						Return()

					return handler
				}(),
			},
			args: args{
				sourceLink: "http://example.com/",
				link:       "http://example.com/test",
			},
			wantLinkRegister: func() register.LinkRegister {
				linkRegister := register.NewLinkRegister(sanitizing.DoNotSanitizeLink, nil)
				linkRegister.RegisterLink("http://example.com/test")

				return linkRegister
			}(),
		},
		{
			name: "with a duplicate",
			fields: fields{
				linkRegister: func() register.LinkRegister {
					linkRegister := register.NewLinkRegister(sanitizing.DoNotSanitizeLink, nil)
					linkRegister.RegisterLink("http://example.com/test")

					return linkRegister
				}(),
				linkHandler: new(MockLinkHandler),
			},
			args: args{
				sourceLink: "http://example.com/",
				link:       "http://example.com/test",
			},
			wantLinkRegister: func() register.LinkRegister {
				linkRegister := register.NewLinkRegister(sanitizing.DoNotSanitizeLink, nil)
				linkRegister.RegisterLink("http://example.com/test")

				return linkRegister
			}(),
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			handler := UniqueHandler{
				linkRegister: data.fields.linkRegister,
				linkHandler:  data.fields.linkHandler,
			}
			handler.HandleLink(data.args.sourceLink, data.args.link)

			mock.AssertExpectationsForObjects(test, data.fields.linkHandler)
			assert.Equal(test, data.wantLinkRegister, handler.linkRegister)
		})
	}
}
