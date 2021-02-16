package checkers

import (
	"context"

	"github.com/thewizardplusplus/go-crawler/models"
	"github.com/thewizardplusplus/go-crawler/registers"
)

// DuplicateChecker ...
type DuplicateChecker struct {
	LinkRegister registers.LinkRegister
}

// CheckLink ...
func (checker DuplicateChecker) CheckLink(
	ctx context.Context,
	link models.SourcedLink,
) bool {
	return checker.LinkRegister.RegisterLink(link.Link)
}
