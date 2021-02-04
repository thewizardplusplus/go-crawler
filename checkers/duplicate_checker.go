package checkers

import (
	"context"

	crawler "github.com/thewizardplusplus/go-crawler"
	"github.com/thewizardplusplus/go-crawler/registers"
)

// DuplicateChecker ...
type DuplicateChecker struct {
	LinkRegister registers.LinkRegister
}

// CheckLink ...
func (checker DuplicateChecker) CheckLink(
	ctx context.Context,
	link crawler.SourcedLink,
) bool {
	return checker.LinkRegister.RegisterLink(link.Link)
}
