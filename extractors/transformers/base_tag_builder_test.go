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
		attributeValues [][]byte
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
				attributeValues: nil,
			},
			wantBaseLink: nil,
			wantIsFound:  assert.False,
		},
		{
			name: "was found",
			fields: fields{
				attributeValues: [][]byte{
					[]byte("base link 1"),
					[]byte("base link 2"),
				},
			},
			wantBaseLink: []byte("base link 2"),
			wantIsFound:  assert.True,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			var builder BaseTagBuilder
			for _, attributeValue := range data.fields.attributeValues {
				builder.AddAttribute(nil, attributeValue)
			}

			gotBaseLink, gotIsFound := builder.BaseLink()

			assert.Equal(test, data.wantBaseLink, gotBaseLink)
			data.wantIsFound(test, gotIsFound)
		})
	}
}

func TestBaseTagBuilder_IsSelectionTerminated(test *testing.T) {
	type fields struct {
		attributeValues  [][]byte
		baseTagSelection BaseTagSelection
	}

	for _, data := range []struct {
		name   string
		fields fields
		wantOk assert.BoolAssertionFunc
	}{
		{
			name: "selection is not terminated (with SelectFirstBaseTag)",
			fields: fields{
				attributeValues:  nil,
				baseTagSelection: SelectFirstBaseTag,
			},
			wantOk: assert.False,
		},
		{
			name: "selection is not terminated (with SelectLastBaseTag)",
			fields: fields{
				attributeValues:  nil,
				baseTagSelection: SelectLastBaseTag,
			},
			wantOk: assert.False,
		},
		{
			name: "selection is not terminated " +
				"(with SelectLastBaseTag and attribute values)",
			fields: fields{
				attributeValues:  [][]byte{[]byte("base link")},
				baseTagSelection: SelectLastBaseTag,
			},
			wantOk: assert.False,
		},
		{
			name: "selection is terminated",
			fields: fields{
				attributeValues:  [][]byte{[]byte("base link")},
				baseTagSelection: SelectFirstBaseTag,
			},
			wantOk: assert.True,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			builder := BaseTagBuilder{
				baseTagSelection: data.fields.baseTagSelection,
			}
			for _, attributeValue := range data.fields.attributeValues {
				builder.AddAttribute(nil, attributeValue)
			}

			gotOk := builder.IsSelectionTerminated()

			data.wantOk(test, gotOk)
		})
	}
}
