package transformers

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTransformerGroup_TransformLinks(test *testing.T) {
	type args struct {
		links           []string
		response        *http.Response
		responseContent []byte
	}

	for _, data := range []struct {
		name         string
		transformers TransformerGroup
		args         args
		wantLinks    []string
		wantErr      assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	} {
		test.Run(data.name, func(test *testing.T) {
			gotLinks, gotErr := data.transformers.TransformLinks(
				data.args.links,
				data.args.response,
				data.args.responseContent,
			)

			for _, transformer := range data.transformers {
				mock.AssertExpectationsForObjects(test, transformer)
			}
			assert.Equal(test, data.wantLinks, gotLinks)
			data.wantErr(test, gotErr)
		})
	}
}
