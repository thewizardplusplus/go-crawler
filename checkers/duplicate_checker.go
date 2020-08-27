package checkers

// DuplicateChecker ...
type DuplicateChecker struct {
	checkedLinks map[string]struct{}
}

// CheckLink ...
func (checker DuplicateChecker) CheckLink(parentLink string, link string) bool {
	_, isDuplicate := checker.checkedLinks[link]
	checker.checkedLinks[link] = struct{}{}

	return isDuplicate
}
