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
		{
			name: "was not found",
			fields: fields{
				baseLink:     []byte("base link"),
				isFirstFound: false,
			},
			wantBaseLink: nil,
			wantIsFound:  assert.False,
		},
		{
			name: "was found",
			fields: fields{
				baseLink:     []byte("base link"),
				isFirstFound: true,
			},
			wantBaseLink: []byte("base link"),
			wantIsFound:  assert.True,
		},
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

func TestBaseTagBuilder_AddTag(test *testing.T) {
	var builder BaseTagBuilder
	builder.AddTag([]byte("tag"))

	assert.Equal(test, BaseTagBuilder{}, builder)
}
