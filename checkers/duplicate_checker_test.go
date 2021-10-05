package checkers

import (
	"context"
	"testing"

	"github.com/go-log/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/go-crawler/models"
	"github.com/thewizardplusplus/go-crawler/registers"
	urlutils "github.com/thewizardplusplus/go-crawler/url-utils"
)

func TestDuplicateChecker_CheckLink(test *testing.T) {
	type fields struct {
		LinkRegister registers.LinkRegister
		Logger       log.Logger
	}
	type args struct {
		ctx  context.Context
		link models.SourcedLink
	}

	for _, data := range []struct {
		name             string
		fields           fields
		args             args
		wantLinkRegister registers.LinkRegister
		wantOk           assert.BoolAssertionFunc
	}{
		{
			name: "without a duplicate",
			fields: fields{
				LinkRegister: registers.NewLinkRegister(urlutils.DoNotSanitizeLink),
				Logger:       new(MockLogger),
			},
			args: args{
				ctx: context.Background(),
				link: models.SourcedLink{
					SourceLink: "http://example.com/",
					Link:       "http://example.com/test",
				},
			},
			wantLinkRegister: func() registers.LinkRegister {
				linkRegister := registers.NewLinkRegister(urlutils.DoNotSanitizeLink)
				linkRegister.RegisterLink("http://example.com/test") // nolint: errcheck

				return linkRegister
			}(),
			wantOk: assert.True,
		},
		{
			name: "with a duplicate",
			fields: fields{
				LinkRegister: func() registers.LinkRegister {
					linkRegister := registers.NewLinkRegister(urlutils.DoNotSanitizeLink)
					linkRegister.RegisterLink("http://example.com/test") // nolint: errcheck

					return linkRegister
				}(),
				Logger: new(MockLogger),
			},
			args: args{
				ctx: context.Background(),
				link: models.SourcedLink{
					SourceLink: "http://example.com/",
					Link:       "http://example.com/test",
				},
			},
			wantLinkRegister: func() registers.LinkRegister {
				linkRegister := registers.NewLinkRegister(urlutils.DoNotSanitizeLink)
				linkRegister.RegisterLink("http://example.com/test") // nolint: errcheck

				return linkRegister
			}(),
			wantOk: assert.False,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			checker := DuplicateChecker{
				LinkRegister: data.fields.LinkRegister,
				Logger:       data.fields.Logger,
			}
			got := checker.CheckLink(data.args.ctx, data.args.link)

			mock.AssertExpectationsForObjects(test, data.fields.Logger)
			assert.Equal(test, data.wantLinkRegister, checker.LinkRegister)
			data.wantOk(test, got)
		})
	}
}
