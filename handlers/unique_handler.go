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
