package sitemap

import (
	"testing"

	"github.com/stretchr/testify/assert"
	httputils "github.com/thewizardplusplus/go-http-utils"
)

func TestLoader_LoadLink(test *testing.T) {
	type fields struct {
		HTTPClient httputils.HTTPClient
	}
	type args struct {
		link    string
		options interface{}
	}

	for _, data := range []struct {
		name             string
		fields           fields
		args             args
		wantResponseData []byte
		wantErr          assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	} {
		test.Run(data.name, func(test *testing.T) {
			loader := Loader{
				HTTPClient: data.fields.HTTPClient,
			}
			gotResponseData, gotErr := loader.LoadLink(data.args.link, data.args.options)

			assert.Equal(test, data.wantResponseData, gotResponseData)
			data.wantErr(test, gotErr)
		})
	}
}
