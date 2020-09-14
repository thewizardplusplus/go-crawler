package checkers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thewizardplusplus/go-crawler/register"
	"github.com/thewizardplusplus/go-crawler/sanitizing"
)

func TestDuplicateChecker_CheckLink(test *testing.T) {
	type fields struct {
		LinkRegister register.LinkRegister
	}
	type args struct {
		sourceLink string
		link       string
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
				sourceLink: "http://example.com/",
				link:       "http://example.com/test",
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
				sourceLink: "http://example.com/",
				link:       "http://example.com/test",
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
			got := checker.CheckLink(data.args.sourceLink, data.args.link)

			assert.Equal(test, data.wantLinkRegister, checker.LinkRegister)
			data.wantOk(test, got)
		})
	}
}
