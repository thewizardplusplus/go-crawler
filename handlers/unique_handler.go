package handlers

import (
	crawler "github.com/thewizardplusplus/go-crawler"
	"github.com/thewizardplusplus/go-crawler/registers"
)

// UniqueHandler ...
type UniqueHandler struct {
	LinkRegister registers.LinkRegister
	LinkHandler  crawler.LinkHandler
}

// HandleLink ...
func (handler UniqueHandler) HandleLink(link crawler.SourcedLink) {
	wasRegistered := handler.LinkRegister.RegisterLink(link.Link)
	if !wasRegistered {
		return
	}

	handler.LinkHandler.HandleLink(link)
}
