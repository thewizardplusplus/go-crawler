package transformers

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

func TestBaseTagBuilder_AddAttribute(test *testing.T) {
	var builder BaseTagBuilder
	builder.AddAttribute([]byte("tag"), []byte("attribute"))

	wantBuilder := BaseTagBuilder{
		baseLink:     []byte("attribute"),
		isFirstFound: true,
	}
	assert.Equal(test, wantBuilder, builder)
}

func TestBaseTagBuilder_IsSelectionTerminated(test *testing.T) {
	type fields struct {
		isFirstFound bool
	}

	for _, data := range []struct {
		name   string
		fields fields
		wantOk assert.BoolAssertionFunc
	}{
		{
			name: "selection is terminated",
			fields: fields{
				isFirstFound: false,
			},
			wantOk: assert.False,
		},
		{
			name: "selection is not terminated",
			fields: fields{
				isFirstFound: true,
			},
			wantOk: assert.True,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			builder := BaseTagBuilder{
				isFirstFound: data.fields.isFirstFound,
			}
			gotOk := builder.IsSelectionTerminated()

			data.wantOk(test, gotOk)
		})
	}
}
