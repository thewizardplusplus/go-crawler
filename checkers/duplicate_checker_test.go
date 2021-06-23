package checkers

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thewizardplusplus/go-crawler/models"
	"github.com/thewizardplusplus/go-crawler/registers"
	urlutils "github.com/thewizardplusplus/go-crawler/url-utils"
)

func TestDuplicateChecker_CheckLink(test *testing.T) {
	type fields struct {
		LinkRegister registers.LinkRegister
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
				LinkRegister: registers.NewLinkRegister(urlutils.DoNotSanitizeLink, nil),
			},
			args: args{
				ctx: context.Background(),
				link: models.SourcedLink{
					SourceLink: "http://example.com/",
					Link:       "http://example.com/test",
				},
			},
			wantLinkRegister: func() registers.LinkRegister {
				linkRegister := registers.NewLinkRegister(urlutils.DoNotSanitizeLink, nil)
				linkRegister.RegisterLink("http://example.com/test")

				return linkRegister
			}(),
			wantOk: assert.True,
		},
		{
			name: "with a duplicate",
			fields: fields{
				LinkRegister: func() registers.LinkRegister {
					linkRegister :=
						registers.NewLinkRegister(urlutils.DoNotSanitizeLink, nil)
					linkRegister.RegisterLink("http://example.com/test")

					return linkRegister
				}(),
			},
			args: args{
				ctx: context.Background(),
				link: models.SourcedLink{
					SourceLink: "http://example.com/",
					Link:       "http://example.com/test",
				},
			},
			wantLinkRegister: func() registers.LinkRegister {
				linkRegister := registers.NewLinkRegister(urlutils.DoNotSanitizeLink, nil)
				linkRegister.RegisterLink("http://example.com/test")

				return linkRegister
			}(),
			wantOk: assert.False,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			checker := DuplicateChecker{
				LinkRegister: data.fields.LinkRegister,
			}
			got := checker.CheckLink(data.args.ctx, data.args.link)

			assert.Equal(test, data.wantLinkRegister, checker.LinkRegister)
			data.wantOk(test, got)
		})
	}
}
