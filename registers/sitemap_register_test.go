package registers

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/yterajima/go-sitemap"
)

func TestSitemapRegister_loadSitemapData(test *testing.T) {
	type fields struct {
		linkLoader         LinkLoader
		registeredSitemaps *sync.Map
	}
	type args struct {
		ctx         context.Context
		sitemapLink string
	}

	for _, data := range []struct {
		name            string
		fields          fields
		args            args
		wantSitemapData sitemap.Sitemap
		wantErr         assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	} {
		test.Run(data.name, func(test *testing.T) {
			sitemap.SetFetch(data.fields.linkLoader.LoadLink)

			register := SitemapRegister{
				registeredSitemaps: data.fields.registeredSitemaps,
			}
			gotSitemapData, gotErr :=
				register.loadSitemapData(data.args.ctx, data.args.sitemapLink)

			mock.AssertExpectationsForObjects(test, data.fields.linkLoader)
			assert.Equal(test, data.wantSitemapData, gotSitemapData)
			data.wantErr(test, gotErr)
		})
	}
}
