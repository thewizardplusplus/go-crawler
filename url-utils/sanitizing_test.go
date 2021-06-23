package urlutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApplyLinkSanitizing(test *testing.T) {
	type args struct {
		link string
	}

	for _, data := range []struct {
		name     string
		args     args
		wantLink string
		wantErr  assert.ErrorAssertionFunc
	}{
		{
			name: "success without sanitizing",
			args: args{
				link: "http://example.com/test",
			},
			wantLink: "http://example.com/test",
			wantErr:  assert.NoError,
		},
		{
			name: "success with sanitizing",
			args: args{
				link: "http://example.com/one/../two",
			},
			wantLink: "http://example.com/two",
			wantErr:  assert.NoError,
		},
		{
			name: "error",
			args: args{
				link: ":",
			},
			wantLink: "",
			wantErr:  assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			gotLink, gotErr := ApplyLinkSanitizing(data.args.link)

			assert.Equal(test, data.wantLink, gotLink)
			data.wantErr(test, gotErr)
		})
	}
}
