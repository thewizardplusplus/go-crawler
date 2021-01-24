package checkers

import (
	"testing"

	"github.com/go-log/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	crawler "github.com/thewizardplusplus/go-crawler"
	"github.com/thewizardplusplus/go-crawler/registers"
)

func TestRobotsTXTChecker_CheckLink(test *testing.T) {
	type fields struct {
		UserAgent         string
		RobotsTXTRegister registers.RobotsTXTRegister
		Logger            log.Logger
	}
	type args struct {
		link crawler.SourcedLink
	}

	for _, data := range []struct {
		name   string
		fields fields
		args   args
		wantOk assert.BoolAssertionFunc
	}{
		// TODO: Add test cases.
	} {
		test.Run(data.name, func(test *testing.T) {
			checker := RobotsTXTChecker{
				UserAgent:         data.fields.UserAgent,
				RobotsTXTRegister: data.fields.RobotsTXTRegister,
				Logger:            data.fields.Logger,
			}
			got := checker.CheckLink(data.args.link)

			mock.AssertExpectationsForObjects(test, data.fields.Logger)
			data.wantOk(test, got)
		})
	}
}
