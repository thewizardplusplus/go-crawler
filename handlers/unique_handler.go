package handlers

import (
	"sync"

	crawler "github.com/thewizardplusplus/go-crawler"
	"github.com/thewizardplusplus/go-crawler/register"
)

// UniqueHandler ...
type UniqueHandler struct {
	locker       sync.Mutex
	linkRegister register.LinkRegister
	linkHandler  crawler.LinkHandler
}

// NewUniqueHandler ...
func NewUniqueHandler(
	linkRegister register.LinkRegister,
	linkHandler crawler.LinkHandler,
) *UniqueHandler {
	return &UniqueHandler{
		linkRegister: linkRegister,
		linkHandler:  linkHandler,
	}
}

// HandleLink ...
func (handler *UniqueHandler) HandleLink(sourceLink string, link string) {
	handler.locker.Lock()
	defer handler.locker.Unlock()

	wasRegistered := handler.linkRegister.RegisterLink(link)
	if !wasRegistered {
		return
	}

	handler.linkHandler.HandleLink(sourceLink, link)
}
