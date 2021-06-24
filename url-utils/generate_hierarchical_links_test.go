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
		{
			name: "success without a trailing slash",
			args: args{
				baseLink:   "http://example.com/test",
				linkSuffix: "suffix",
				options:    nil,
			},
			wantLinks: []string{"http://example.com/suffix"},
			wantErr:   assert.NoError,
		},
		{
			name: "success with a trailing slash",
			args: args{
				baseLink:   "http://example.com/test/",
				linkSuffix: "suffix",
				options:    nil,
			},
			wantLinks: []string{
				"http://example.com/suffix",
				"http://example.com/test/suffix",
			},
			wantErr: assert.NoError,
		},
		{
			name: "success with a long path and without a depth limit",
			args: args{
				baseLink:   "http://example.com/one/two/three/test",
				linkSuffix: "suffix",
				options:    nil,
			},
			wantLinks: []string{
				"http://example.com/suffix",
				"http://example.com/one/suffix",
				"http://example.com/one/two/suffix",
				"http://example.com/one/two/three/suffix",
			},
			wantErr: assert.NoError,
		},
		{
			name: "success with a long path and normal depth limit",
			args: args{
				baseLink:   "http://example.com/one/two/three/test",
				linkSuffix: "suffix",
				options:    []HierarchicalLinkOption{WithMaximalHierarchyDepth(2)},
			},
			wantLinks: []string{
				"http://example.com/suffix",
				"http://example.com/one/suffix",
				"http://example.com/one/two/suffix",
			},
			wantErr: assert.NoError,
		},
		{
			name: "success with a long path and zero depth limit",
			args: args{
				baseLink:   "http://example.com/one/two/three/test",
				linkSuffix: "suffix",
				options:    []HierarchicalLinkOption{WithMaximalHierarchyDepth(0)},
			},
			wantLinks: []string{"http://example.com/suffix"},
			wantErr:   assert.NoError,
		},
		{
			name: "success with sanitizing",
			args: args{
				baseLink:   "http://example.com/one/two//three/../test",
				linkSuffix: "suffix",
				options:    []HierarchicalLinkOption{SanitizeBaseLink(SanitizeLink)},
			},
			wantLinks: []string{
				"http://example.com/suffix",
				"http://example.com/one/suffix",
				"http://example.com/one/two/suffix",
			},
			wantErr: assert.NoError,
		},
		{
			name: "success without sanitizing",
			args: args{
				baseLink:   "http://example.com/one/two//three/../test",
				linkSuffix: "suffix",
				options:    nil,
			},
			wantLinks: []string{
				"http://example.com/suffix",
				"http://example.com/one/suffix",
				"http://example.com/one/two/suffix",
				"http://example.com/one/two//suffix",
				"http://example.com/one/two//three/suffix",
				"http://example.com/one/two//three/../suffix",
			},
			wantErr: assert.NoError,
		},
		{
			name: "success with an HTTPS scheme",
			args: args{
				baseLink:   "https://example.com/test",
				linkSuffix: "suffix",
				options:    nil,
			},
			wantLinks: []string{"https://example.com/suffix"},
			wantErr:   assert.NoError,
		},
		{
			name: "success with a user",
			args: args{
				baseLink:   "http://username:password@example.com/test",
				linkSuffix: "suffix",
				options:    nil,
			},
			wantLinks: []string{"http://username:password@example.com/suffix"},
			wantErr:   assert.NoError,
		},
		{
			name: "success with a query",
			args: args{
				baseLink:   "http://example.com/test?key=value",
				linkSuffix: "suffix",
				options:    nil,
			},
			wantLinks: []string{"http://example.com/suffix"},
			wantErr:   assert.NoError,
		},
		{
			name: "success with a fragment",
			args: args{
				baseLink:   "http://example.com/test#fragment",
				linkSuffix: "suffix",
				options:    nil,
			},
			wantLinks: []string{"http://example.com/suffix"},
			wantErr:   assert.NoError,
		},
		{
			name: "error with sanitizing",
			args: args{
				baseLink:   ":",
				linkSuffix: "suffix",
				options:    []HierarchicalLinkOption{SanitizeBaseLink(SanitizeLink)},
			},
			wantLinks: nil,
			wantErr:   assert.Error,
		},
		{
			name: "error with parsing",
			args: args{
				baseLink:   ":",
				linkSuffix: "suffix",
				options:    nil,
			},
			wantLinks: nil,
			wantErr:   assert.Error,
		},
		{
			name: "error with a not absolute path",
			args: args{
				baseLink:   "one/two/three/test",
				linkSuffix: "suffix",
				options:    nil,
			},
			wantLinks: nil,
			wantErr:   assert.Error,
		},
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
