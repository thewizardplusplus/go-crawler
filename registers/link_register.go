package registers

import (
	mapset "github.com/deckarep/golang-set"
	"github.com/pkg/errors"
	urlutils "github.com/thewizardplusplus/go-crawler/url-utils"
)

// LinkRegister ...
type LinkRegister struct {
	sanitizeLink urlutils.LinkSanitizing

	registeredLinks mapset.Set
}

// NewLinkRegister ...
func NewLinkRegister(sanitizeLink urlutils.LinkSanitizing) LinkRegister {
	return LinkRegister{
		sanitizeLink: sanitizeLink,

		registeredLinks: mapset.NewSet(),
	}
}

// RegisterLink ...
func (register LinkRegister) RegisterLink(link string) (
	wasRegistered bool,
	err error,
) {
	if register.sanitizeLink == urlutils.SanitizeLink {
		sanitizedLink, err := urlutils.ApplyLinkSanitizing(link)
		if err != nil {
			return false, errors.Wrapf(err, "unable to sanitize the link")
		}

		link = sanitizedLink
	}

	wasRegistered = register.registeredLinks.Add(link)
	return wasRegistered, nil
}
