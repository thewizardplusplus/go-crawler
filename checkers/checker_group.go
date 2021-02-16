package checkers

import (
	"context"

	crawler "github.com/thewizardplusplus/go-crawler"
	"github.com/thewizardplusplus/go-crawler/models"
)

// CheckerGroup ...
type CheckerGroup []crawler.LinkChecker

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
