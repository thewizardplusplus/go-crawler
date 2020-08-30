package checkers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCheckerGroup_CheckLink(test *testing.T) {
	type args struct {
		parentLink string
		link       string
	}

	for _, data := range []struct {
		name     string
		checkers CheckerGroup
		args     args
		want     assert.BoolAssertionFunc
	}{
		// TODO: Add test cases.
	} {
		test.Run(data.name, func(test *testing.T) {
			got := data.checkers.CheckLink(data.args.parentLink, data.args.link)

			for _, checker := range data.checkers {
				mock.AssertExpectationsForObjects(test, checker)
			}
			data.want(test, got)
		})
	}
}
