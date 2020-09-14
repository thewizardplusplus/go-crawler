package checkers

import (
	"github.com/thewizardplusplus/go-crawler/register"
)

// DuplicateChecker ...
type DuplicateChecker struct {
	LinkRegister register.LinkRegister
}

// CheckLink ...
func (checker DuplicateChecker) CheckLink(
	sourceLink string,
	link string,
) bool {
	return checker.LinkRegister.RegisterLink(link)
}
