package sanitizing

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
		// TODO: Add test cases.
	} {
		test.Run(data.name, func(test *testing.T) {
			gotLink, gotErr := ApplyLinkSanitizing(data.args.link)

			assert.Equal(test, data.wantLink, gotLink)
			data.wantErr(test, gotErr)
		})
	}
}
