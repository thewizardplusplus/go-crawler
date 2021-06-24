package urlutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateHierarchicalLinks(test *testing.T) {
	type args struct {
		baseLink   string
		linkSuffix string
		options    []HierarchicalLinkOption
	}

	for _, data := range []struct {
		name      string
		args      args
		wantLinks []string
		wantErr   assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	} {
		test.Run(data.name, func(test *testing.T) {
			gotLinks, gotErr := GenerateHierarchicalLinks(
				data.args.baseLink,
				data.args.linkSuffix,
				data.args.options...,
			)

			assert.Equal(test, data.wantLinks, gotLinks)
			data.wantErr(test, gotErr)
		})
	}
}
