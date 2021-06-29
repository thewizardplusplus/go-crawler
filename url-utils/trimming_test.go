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
		{
			name: "DoNotTrimLink",
			args: args{
				link:     "  http://example.com/  ",
				trimming: DoNotTrimLink,
			},
			wantedLink: "  http://example.com/  ",
		},
		{
			name: "TrimLinkLeft",
			args: args{
				link:     "  http://example.com/  ",
				trimming: TrimLinkLeft,
			},
			wantedLink: "http://example.com/  ",
		},
		{
			name: "TrimLinkRight",
			args: args{
				link:     "  http://example.com/  ",
				trimming: TrimLinkRight,
			},
			wantedLink: "  http://example.com/",
		},
		{
			name: "TrimLink",
			args: args{
				link:     "  http://example.com/  ",
				trimming: TrimLink,
			},
			wantedLink: "http://example.com/",
		},
		{
			name: "without spaces",
			args: args{
				link:     "http://example.com/",
				trimming: TrimLink,
			},
			wantedLink: "http://example.com/",
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			receivedLink := ApplyLinkTrimming(data.args.link, data.args.trimming)

			assert.Equal(test, data.wantedLink, receivedLink)
		})
	}
}
