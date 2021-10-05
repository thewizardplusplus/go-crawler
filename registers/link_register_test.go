package registers

import (
	"testing"

	mapset "github.com/deckarep/golang-set"
	"github.com/stretchr/testify/assert"
	urlutils "github.com/thewizardplusplus/go-crawler/url-utils"
)

func TestNewLinkRegister(test *testing.T) {
	got := NewLinkRegister(urlutils.SanitizeLink)

	assert.Equal(test, urlutils.SanitizeLink, got.sanitizeLink)
	assert.Equal(test, mapset.NewSet(), got.registeredLinks)
}

func TestLinkRegister_RegisterLink(test *testing.T) {
	type fields struct {
		sanitizeLink    urlutils.LinkSanitizing
		registeredLinks mapset.Set
	}
	type args struct {
		link string
	}

	for _, data := range []struct {
		name              string
		fields            fields
		args              args
		wantWasRegistered assert.BoolAssertionFunc
		wantErr           assert.ErrorAssertionFunc
	}{
		{
			name: "success without a duplicate",
			fields: fields{
				sanitizeLink: urlutils.DoNotSanitizeLink,
				registeredLinks: mapset.NewSet(
					"http://example.com/1",
					"http://example.com/2",
				),
			},
			args: args{
				link: "http://example.com/3",
			},
			wantWasRegistered: assert.True,
			wantErr:           assert.NoError,
		},
		{
			name: "success with a duplicate and without link sanitizing",
			fields: fields{
				sanitizeLink: urlutils.DoNotSanitizeLink,
				registeredLinks: mapset.NewSet(
					"http://example.com/1",
					"http://example.com/2",
				),
			},
			args: args{
				link: "http://example.com/2",
			},
			wantWasRegistered: assert.False,
			wantErr:           assert.NoError,
		},
		{
			name: "success with a duplicate and with link sanitizing",
			fields: fields{
				sanitizeLink: urlutils.SanitizeLink,
				registeredLinks: mapset.NewSet(
					"http://example.com/1",
					"http://example.com/2",
				),
			},
			args: args{
				link: "http://example.com/test/../2",
			},
			wantWasRegistered: assert.False,
			wantErr:           assert.NoError,
		},
		{
			name: "error",
			fields: fields{
				sanitizeLink: urlutils.SanitizeLink,
				registeredLinks: mapset.NewSet(
					"http://example.com/1",
					"http://example.com/2",
				),
			},
			args: args{
				link: ":",
			},
			wantWasRegistered: assert.False,
			wantErr:           assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			register := LinkRegister{
				sanitizeLink:    data.fields.sanitizeLink,
				registeredLinks: data.fields.registeredLinks,
			}
			gotWasRegistered, gotErr := register.RegisterLink(data.args.link)

			data.wantWasRegistered(test, gotWasRegistered)
			data.wantErr(test, gotErr)
		})
	}
}
