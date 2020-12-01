package checkers

import (
	crawler "github.com/thewizardplusplus/go-crawler"
	"github.com/thewizardplusplus/go-crawler/registers"
)

// DuplicateChecker ...
type DuplicateChecker struct {
	LinkRegister registers.LinkRegister
}

// CheckLink ...
func (checker DuplicateChecker) CheckLink(link crawler.SourcedLink) bool {
	return checker.LinkRegister.RegisterLink(link.Link)
}
