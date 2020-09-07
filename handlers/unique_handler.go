package handlers

import (
	"sync"

	mapset "github.com/deckarep/golang-set"
	"github.com/go-log/log"
	crawler "github.com/thewizardplusplus/go-crawler"
	"github.com/thewizardplusplus/go-crawler/sanitizing"
)

// UniqueHandler ...
type UniqueHandler struct {
	sanitizeLink sanitizing.LinkSanitizing
	linkHandler  crawler.LinkHandler
	logger       log.Logger

	locker       sync.RWMutex
	handledLinks mapset.Set
}

// NewUniqueHandler ...
func NewUniqueHandler(
	sanitizeLink sanitizing.LinkSanitizing,
	linkHandler crawler.LinkHandler,
	logger log.Logger,
) *UniqueHandler {
	return &UniqueHandler{
		sanitizeLink: sanitizeLink,
		linkHandler:  linkHandler,
		logger:       logger,

		handledLinks: mapset.NewThreadUnsafeSet(),
	}
}

// HandleLink ...
func (handler *UniqueHandler) HandleLink(sourceLink string, link string) {
	if handler.sanitizeLink == sanitizing.SanitizeLink {
		var err error
		link, err = sanitizing.ApplyLinkSanitizing(link)
		if err != nil {
			handler.logger.Logf("unable to sanitize the link: %s", err)
			return
		}
	}

	handler.locker.Lock()
	defer handler.locker.Unlock()

	// checking should be occurred under lock
	if handler.handledLinks.Contains(link) {
		return
	}
	// use the defer operator to remember a link
	// even if the inner handler throws a panic
	defer handler.handledLinks.Add(link)

	handler.linkHandler.HandleLink(sourceLink, link)
}
