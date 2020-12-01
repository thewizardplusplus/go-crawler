package registers

import (
	mapset "github.com/deckarep/golang-set"
	"github.com/go-log/log"
	"github.com/thewizardplusplus/go-crawler/sanitizing"
)

// LinkRegister ...
type LinkRegister struct {
	sanitizeLink sanitizing.LinkSanitizing
	logger       log.Logger

	registeredLinks mapset.Set
}

// NewLinkRegister ...
func NewLinkRegister(
	sanitizeLink sanitizing.LinkSanitizing,
	logger log.Logger,
) LinkRegister {
	return LinkRegister{
		sanitizeLink: sanitizeLink,
		logger:       logger,

		registeredLinks: mapset.NewSet(),
	}
}

// RegisterLink ...
func (register LinkRegister) RegisterLink(link string) bool {
	if register.sanitizeLink == sanitizing.SanitizeLink {
		var err error
		link, err = sanitizing.ApplyLinkSanitizing(link)
		if err != nil {
			register.logger.Logf("unable to sanitize the link: %s", err)
			return false
		}
	}

	return register.registeredLinks.Add(link)
}
