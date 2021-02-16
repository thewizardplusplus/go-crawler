package checkers

import (
	"context"

	"github.com/thewizardplusplus/go-crawler/models"
)

// CheckerGroup ...
type CheckerGroup []models.LinkChecker

// CheckLink ...
func (checkers CheckerGroup) CheckLink(
	ctx context.Context,
	link models.SourcedLink,
) bool {
	for _, checker := range checkers {
		if !checker.CheckLink(ctx, link) {
			return false
		}
	}

	// to prohibit using an empty group as a filter that passes everything
	return len(checkers) != 0
}
