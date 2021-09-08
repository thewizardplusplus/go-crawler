package extractors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBaseTagBuilder_BaseLink(test *testing.T) {
	type fields struct {
		baseLink     []byte
		isFirstFound bool
	}

	for _, data := range []struct {
		name         string
		fields       fields
		wantBaseLink []byte
		wantIsFound  assert.BoolAssertionFunc
	}{
		// TODO: Add test cases.
	} {
		test.Run(data.name, func(test *testing.T) {
			builder := BaseTagBuilder{
				baseLink:     data.fields.baseLink,
				isFirstFound: data.fields.isFirstFound,
			}
			gotBaseLink, gotIsFound := builder.BaseLink()

			assert.Equal(test, data.wantBaseLink, gotBaseLink)
			data.wantIsFound(test, gotIsFound)
		})
	}
}
