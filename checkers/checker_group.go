package checkers

import (
	"context"

	crawler "github.com/thewizardplusplus/go-crawler"
)

// CheckerGroup ...
type CheckerGroup []crawler.LinkChecker

// CheckLink ...
func (checkers CheckerGroup) CheckLink(
	ctx context.Context,
	link crawler.SourcedLink,
) bool {
	for _, checker := range checkers {
		if !checker.CheckLink(ctx, link) {
			return false
		}
	}

	// to prohibit using an empty group as a filter that passes everything
	return len(checkers) != 0
}
