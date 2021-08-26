package urlutils

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLinkResolver_ResolveLink(test *testing.T) {
	type fields struct {
		BaseLink *url.URL
	}
	type args struct {
		link string
	}

	for _, data := range []struct {
		name     string
		fields   fields
		args     args
		wantLink string
		wantErr  assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	} {
		test.Run(data.name, func(test *testing.T) {
			resolver := LinkResolver{
				BaseLink: data.fields.BaseLink,
			}
			gotLink, gotErr := resolver.ResolveLink(data.args.link)

			assert.Equal(test, data.wantLink, gotLink)
			data.wantErr(test, gotErr)
		})
	}
}
