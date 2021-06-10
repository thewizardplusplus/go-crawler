package sitemap

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thewizardplusplus/go-crawler/registers"
)

func TestRobotsTXTGenerator_GenerateLinks(test *testing.T) {
	type fields struct {
		RobotsTXTRegister registers.RobotsTXTRegister
	}
	type args struct {
		ctx      context.Context
		baseLink string
	}

	for _, data := range []struct {
		name             string
		fields           fields
		args             args
		wantSitemapLinks []string
		wantErr          assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	} {
		test.Run(data.name, func(test *testing.T) {
			generator := RobotsTXTGenerator{
				RobotsTXTRegister: data.fields.RobotsTXTRegister,
			}
			gotSitemapLinks, gotErr :=
				generator.GenerateLinks(data.args.ctx, data.args.baseLink)

			assert.Equal(test, data.wantSitemapLinks, gotSitemapLinks)
			data.wantErr(test, gotErr)
		})
	}
}
