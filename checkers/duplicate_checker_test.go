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
		sanitizeLink bool
		logger       log.Logger

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
				sanitizeLink: data.fields.sanitizeLink,
				logger:       data.fields.logger,

				checkedLinks: data.fields.checkedLinks,
			}
			got := checker.CheckLink(data.args.parentLink, data.args.link)

			mock.AssertExpectationsForObjects(test, data.fields.logger)
			data.want(test, got)
		})
	}
}
