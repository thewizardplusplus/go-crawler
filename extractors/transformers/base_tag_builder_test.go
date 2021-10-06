package transformers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBaseTagBuilder(test *testing.T) {
	got := NewBaseTagBuilder(SelectLastBaseTag)

	want := BaseTagBuilder{
		baseTagSelection: SelectLastBaseTag,
	}
	assert.Equal(test, want, got)
}

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
		baseTagSelection BaseTagSelection
		isFirstFound     bool
	}

	for _, data := range []struct {
		name   string
		fields fields
		wantOk assert.BoolAssertionFunc
	}{
		{
			name: "selection is not terminated (SelectFirstBaseTag)",
			fields: fields{
				baseTagSelection: SelectFirstBaseTag,
				isFirstFound:     false,
			},
			wantOk: assert.False,
		},
		{
			name: "selection is not terminated (SelectLastBaseTag)",
			fields: fields{
				baseTagSelection: SelectLastBaseTag,
				isFirstFound:     false,
			},
			wantOk: assert.False,
		},
		{
			name: "selection is not terminated (SelectLastBaseTag and isFirstFound)",
			fields: fields{
				baseTagSelection: SelectLastBaseTag,
				isFirstFound:     true,
			},
			wantOk: assert.False,
		},
		{
			name: "selection is terminated",
			fields: fields{
				baseTagSelection: SelectFirstBaseTag,
				isFirstFound:     true,
			},
			wantOk: assert.True,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			builder := BaseTagBuilder{
				baseTagSelection: data.fields.baseTagSelection,
				isFirstFound:     data.fields.isFirstFound,
			}
			gotOk := builder.IsSelectionTerminated()

			data.wantOk(test, gotOk)
		})
	}
}
