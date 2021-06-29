package urlutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApplyLinkTrimming(test *testing.T) {
	type args struct {
		link     string
		trimming LinkTrimming
	}

	for _, data := range []struct {
		name       string
		args       args
		wantedLink string
	}{
		// TODO: Add test cases.
	} {
		test.Run(data.name, func(test *testing.T) {
			receivedLink := ApplyLinkTrimming(data.args.link, data.args.trimming)

			assert.Equal(test, data.wantedLink, receivedLink)
		})
	}
}
