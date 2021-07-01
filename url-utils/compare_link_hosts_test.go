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
		{
			name: "success with same links",
			args: args{
				linkOne: "http://example.com/",
				linkTwo: "http://example.com/test",
			},
			wantResult: Same,
			wantErr:    assert.NoError,
		},
		{
			name: "success with different links",
			args: args{
				linkOne: "http://example1.com/",
				linkTwo: "http://example2.com/test",
			},
			wantResult: Different,
			wantErr:    assert.NoError,
		},
		{
			name: "error with link one",
			args: args{
				linkOne: ":",
				linkTwo: "http://example.com/test",
			},
			wantResult: 0,
			wantErr:    assert.Error,
		},
		{
			name: "error with link two",
			args: args{
				linkOne: "http://example.com/",
				linkTwo: ":",
			},
			wantResult: 0,
			wantErr:    assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			gotResult, gotErr := CompareLinkHosts(data.args.linkOne, data.args.linkTwo)

			assert.Equal(test, data.wantResult, gotResult)
			data.wantErr(test, gotErr)
		})
	}
}
