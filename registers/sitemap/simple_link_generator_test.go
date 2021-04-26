package sitemap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleLinkGenerator_GenerateLinks(test *testing.T) {
	type args struct {
		baseLink string
	}

	for _, data := range []struct {
		name             string
		args             args
		wantSitemapLinks []string
		wantErr          assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	} {
		test.Run(data.name, func(test *testing.T) {
			var generator SimpleLinkGenerator
			gotSitemapLinks, gotErr := generator.GenerateLinks(data.args.baseLink)

			assert.Equal(test, data.wantSitemapLinks, gotSitemapLinks)
			data.wantErr(test, gotErr)
		})
	}
}
