package checkers

import (
	mapset "github.com/deckarep/golang-set"
)

// DuplicateChecker ...
type DuplicateChecker struct {
	checkedLinks mapset.Set
}

// CheckLink ...
func (checker DuplicateChecker) CheckLink(parentLink string, link string) bool {
	isDuplicate := checker.checkedLinks.Contains(link)
	checker.checkedLinks.Add(link)

	return isDuplicate
}
