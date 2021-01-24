package handlers

import (
	"testing"

	"github.com/go-log/log"
	"github.com/stretchr/testify/mock"
	crawler "github.com/thewizardplusplus/go-crawler"
	"github.com/thewizardplusplus/go-crawler/registers"
)

func TestRobotsTXTHandler_HandleLink(test *testing.T) {
	type fields struct {
		UserAgent         string
		RobotsTXTRegister registers.RobotsTXTRegister
		LinkHandler       crawler.LinkHandler
		Logger            log.Logger
	}
	type args struct {
		link crawler.SourcedLink
	}

	for _, data := range []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	} {
		test.Run(data.name, func(test *testing.T) {
			handler := RobotsTXTHandler{
				UserAgent:         data.fields.UserAgent,
				RobotsTXTRegister: data.fields.RobotsTXTRegister,
				LinkHandler:       data.fields.LinkHandler,
				Logger:            data.fields.Logger,
			}
			handler.HandleLink(data.args.link)

			mock.AssertExpectationsForObjects(
				test,
				data.fields.LinkHandler,
				data.fields.Logger,
			)
		})
	}
}
