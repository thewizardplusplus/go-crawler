package handlers

import (
	crawler "github.com/thewizardplusplus/go-crawler"
	"github.com/thewizardplusplus/go-crawler/register"
)

// UniqueHandler ...
type UniqueHandler struct {
	LinkRegister register.LinkRegister
	LinkHandler  crawler.LinkHandler
}

// HandleLink ...
func (handler UniqueHandler) HandleLink(sourceLink string, link string) {
	wasRegistered := handler.LinkRegister.RegisterLink(link)
	if !wasRegistered {
		return
	}

	handler.LinkHandler.HandleLink(sourceLink, link)
}
