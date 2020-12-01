package registers

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_makeRobotsTXTLink(test *testing.T) {
	type args struct {
		regularLink string
	}

	for _, data := range []struct {
		name              string
		args              args
		wantRobotsTXTLink *url.URL
		wantErr           assert.ErrorAssertionFunc
	}{
		{
			name: "success with a path only",
			args: args{
				regularLink: "http://example.com/test",
			},
			wantRobotsTXTLink: &url.URL{
				Scheme: "http",
				Host:   "example.com",
				Path:   "/robots.txt",
			},
			wantErr: assert.NoError,
		},
		{
			name: "success with an HTTPS scheme",
			args: args{
				regularLink: "https://example.com/test",
			},
			wantRobotsTXTLink: &url.URL{
				Scheme: "https",
				Host:   "example.com",
				Path:   "/robots.txt",
			},
			wantErr: assert.NoError,
		},
		{
			name: "success with an user",
			args: args{
				regularLink: "http://username:password@example.com/test",
			},
			wantRobotsTXTLink: &url.URL{
				Scheme: "http",
				User:   url.UserPassword("username", "password"),
				Host:   "example.com",
				Path:   "/robots.txt",
			},
			wantErr: assert.NoError,
		},
		{
			name: "success with a query",
			args: args{
				regularLink: "http://example.com/test?key=value",
			},
			wantRobotsTXTLink: &url.URL{
				Scheme: "http",
				Host:   "example.com",
				Path:   "/robots.txt",
			},
			wantErr: assert.NoError,
		},
		{
			name: "success with a fragment",
			args: args{
				regularLink: "http://example.com/test#fragment",
			},
			wantRobotsTXTLink: &url.URL{
				Scheme: "http",
				Host:   "example.com",
				Path:   "/robots.txt",
			},
			wantErr: assert.NoError,
		},
		{
			name: "error",
			args: args{
				regularLink: ":",
			},
			wantRobotsTXTLink: nil,
			wantErr:           assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			gotRobotsTXTLink, gotErr := makeRobotsTXTLink(data.args.regularLink)

			assert.Equal(test, data.wantRobotsTXTLink, gotRobotsTXTLink)
			data.wantErr(test, gotErr)
		})
	}
}
