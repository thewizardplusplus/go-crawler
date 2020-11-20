package checkers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	crawler "github.com/thewizardplusplus/go-crawler"
	"github.com/thewizardplusplus/go-crawler/register"
	"github.com/thewizardplusplus/go-crawler/sanitizing"
)

func TestDuplicateChecker_CheckLink(test *testing.T) {
	type fields struct {
		LinkRegister register.LinkRegister
	}
	type args struct {
		link crawler.SourcedLink
	}

	for _, data := range []struct {
		name             string
		fields           fields
		args             args
		wantLinkRegister register.LinkRegister
		wantOk           assert.BoolAssertionFunc
	}{
		{
			name: "without a duplicate",
			fields: fields{
				LinkRegister: register.NewLinkRegister(sanitizing.DoNotSanitizeLink, nil),
			},
			args: args{
				link: crawler.SourcedLink{
					SourceLink: "http://example.com/",
					Link:       "http://example.com/test",
				},
			},
			wantLinkRegister: func() register.LinkRegister {
				linkRegister := register.NewLinkRegister(sanitizing.DoNotSanitizeLink, nil)
				linkRegister.RegisterLink("http://example.com/test")

				return linkRegister
			}(),
			wantOk: assert.True,
		},
		{
			name: "with a duplicate",
			fields: fields{
				LinkRegister: func() register.LinkRegister {
					linkRegister := register.NewLinkRegister(sanitizing.DoNotSanitizeLink, nil)
					linkRegister.RegisterLink("http://example.com/test")

					return linkRegister
				}(),
			},
			args: args{
				link: crawler.SourcedLink{
					SourceLink: "http://example.com/",
					Link:       "http://example.com/test",
				},
			},
			wantLinkRegister: func() register.LinkRegister {
				linkRegister := register.NewLinkRegister(sanitizing.DoNotSanitizeLink, nil)
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
			got := checker.CheckLink(data.args.link)

			assert.Equal(test, data.wantLinkRegister, checker.LinkRegister)
			data.wantOk(test, got)
		})
	}
}
