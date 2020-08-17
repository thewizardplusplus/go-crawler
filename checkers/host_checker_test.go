package checkers

import (
	"testing"

	"github.com/go-log/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHostChecker_CheckLink(test *testing.T) {
	type fields struct {
		Logger log.Logger
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
			checker := HostChecker{
				Logger: data.fields.Logger,
			}
			got := checker.CheckLink(data.args.parentLink, data.args.link)

			mock.AssertExpectationsForObjects(test, data.fields.Logger)
			data.want(test, got)
		})
	}
}
