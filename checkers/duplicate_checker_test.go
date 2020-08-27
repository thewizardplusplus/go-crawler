package checkers

import (
	"testing"

	mapset "github.com/deckarep/golang-set"
	"github.com/go-log/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDuplicateChecker_CheckLink(test *testing.T) {
	type fields struct {
		SanitizeLink bool
		Logger       log.Logger

		checkedLinks mapset.Set
	}
	type args struct {
		parentLink string
		link       string
	}

	for _, data := range []struct {
		name   string
		fields fields
		args   args
		want   assert.BoolAssertionFunc
	}{
		// TODO: Add test cases.
	} {
		test.Run(data.name, func(test *testing.T) {
			checker := DuplicateChecker{
				SanitizeLink: data.fields.SanitizeLink,
				Logger:       data.fields.Logger,

				checkedLinks: data.fields.checkedLinks,
			}
			got := checker.CheckLink(data.args.parentLink, data.args.link)

			mock.AssertExpectationsForObjects(test, data.fields.Logger)
			data.want(test, got)
		})
	}
}
