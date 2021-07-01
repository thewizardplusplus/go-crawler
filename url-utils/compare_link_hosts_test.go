package urlutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompareLinkHosts(test *testing.T) {
	type args struct {
		linkOne string
		linkTwo string
	}

	for _, data := range []struct {
		name       string
		args       args
		wantResult ComparisonResult
		wantErr    assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	} {
		test.Run(data.name, func(test *testing.T) {
			gotResult, gotErr := CompareLinkHosts(data.args.linkOne, data.args.linkTwo)

			assert.Equal(test, data.wantResult, gotResult)
			data.wantErr(test, gotErr)
		})
	}
}
