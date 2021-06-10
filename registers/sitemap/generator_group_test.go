package sitemap

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGeneratorGroup_GenerateLinks(test *testing.T) {
	type args struct {
		ctx      context.Context
		baseLink string
	}

	for _, data := range []struct {
		name             string
		generators       GeneratorGroup
		args             args
		wantSitemapLinks []string
		wantErr          assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	} {
		test.Run(data.name, func(test *testing.T) {
			gotSitemapLinks, gotErr :=
				data.generators.GenerateLinks(data.args.ctx, data.args.baseLink)

			assert.Equal(test, data.wantSitemapLinks, gotSitemapLinks)
			data.wantErr(test, gotErr)
		})
	}
}
