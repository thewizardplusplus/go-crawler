package checkers

import (
	crawler "github.com/thewizardplusplus/go-crawler"
)

// CheckerGroup ...
type CheckerGroup []crawler.LinkChecker

// CheckLink ...
func (checkers CheckerGroup) CheckLink(sourceLink string, link string) bool {
	for _, checker := range checkers {
		if !checker.CheckLink(sourceLink, link) {
			return false
		}
	}

	// to prohibit using an empty group as a filter that passes everything
	return len(checkers) != 0
}
