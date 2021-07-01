package urlutils

import (
	"net/url"

	"github.com/pkg/errors"
)

// ComparisonResult ...
type ComparisonResult int

// ...
const (
	Same ComparisonResult = iota
	Different
)

// CompareLinkHosts ...
func CompareLinkHosts(linkOne string, linkTwo string) (
	ComparisonResult,
	error,
) {
	parsedLinkOne, err := url.Parse(linkOne)
	if err != nil {
		return 0, errors.Wrapf(err, "unable to parse link %q", linkOne)
	}

	parsedLinkTwo, err := url.Parse(linkTwo)
	if err != nil {
		return 0, errors.Wrapf(err, "unable to parse link %q", linkTwo)
	}

	var result ComparisonResult
	if parsedLinkOne.Host == parsedLinkTwo.Host {
		result = Same
	} else {
		result = Different
	}

	return result, nil
}
