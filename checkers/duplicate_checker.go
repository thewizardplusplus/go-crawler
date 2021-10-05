package checkers

import (
	"context"

	"github.com/go-log/log"
	"github.com/thewizardplusplus/go-crawler/models"
	"github.com/thewizardplusplus/go-crawler/registers"
)

// DuplicateChecker ...
type DuplicateChecker struct {
	LinkRegister registers.LinkRegister
	Logger       log.Logger
}

// CheckLink ...
func (checker DuplicateChecker) CheckLink(
	ctx context.Context,
	link models.SourcedLink,
) bool {
	const logPrefix = "duplicate checking"

	wasRegistered, err := checker.LinkRegister.RegisterLink(link.Link)
	if err != nil {
		const logMessage = "%s: unable to register link %q: %s"
		checker.Logger.Logf(logMessage, logPrefix, link.Link, err)

		return false
	}

	return wasRegistered
}
