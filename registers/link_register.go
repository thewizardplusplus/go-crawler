package registers

import (
	mapset "github.com/deckarep/golang-set"
	"github.com/go-log/log"
	urlutils "github.com/thewizardplusplus/go-crawler/url-utils"
)

// LinkRegister ...
type LinkRegister struct {
	sanitizeLink urlutils.LinkSanitizing
	logger       log.Logger

	registeredLinks mapset.Set
}

// NewLinkRegister ...
func NewLinkRegister(
	sanitizeLink urlutils.LinkSanitizing,
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
	if register.sanitizeLink == urlutils.SanitizeLink {
		sanitizedLink, err := urlutils.ApplyLinkSanitizing(link)
		if err != nil {
			register.logger.Logf("unable to sanitize link %q: %s", link, err)
			return false
		}

		link = sanitizedLink
	}

	return register.registeredLinks.Add(link)
}
