package register

import (
	"testing"

	mapset "github.com/deckarep/golang-set"
	"github.com/go-log/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/go-crawler/sanitizing"
)

func TestNewLinkRegister(test *testing.T) {
	logger := new(MockLogger)
	got := NewLinkRegister(sanitizing.SanitizeLink, logger)

	mock.AssertExpectationsForObjects(test, logger)
	assert.Equal(test, sanitizing.SanitizeLink, got.sanitizeLink)
	assert.Equal(test, logger, got.logger)
	assert.NotNil(test, got.registeredLinks)
}

func TestLinkRegister_RegisterLink(test *testing.T) {
	type fields struct {
		sanitizeLink sanitizing.LinkSanitizing
		logger       log.Logger

		registeredLinks mapset.Set
	}
	type args struct {
		link string
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
			register := LinkRegister{
				sanitizeLink: data.fields.sanitizeLink,
				logger:       data.fields.logger,

				registeredLinks: data.fields.registeredLinks,
			}
			got := register.RegisterLink(data.args.link)

			mock.AssertExpectationsForObjects(test, data.fields.logger)
			data.want(test, got)
		})
	}
}
