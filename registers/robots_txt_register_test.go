package registers

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/temoto/robotstxt"
)

func Test_makeRobotsTXTLink(test *testing.T) {
	type args struct {
		regularLink string
	}

	for _, data := range []struct {
		name              string
		args              args
		wantRobotsTXTLink string
		wantErr           assert.ErrorAssertionFunc
	}{
		{
			name: "success with a path only",
			args: args{
				regularLink: "http://example.com/test",
			},
			wantRobotsTXTLink: "http://example.com/robots.txt",
			wantErr:           assert.NoError,
		},
		{
			name: "success with an HTTPS scheme",
			args: args{
				regularLink: "https://example.com/test",
			},
			wantRobotsTXTLink: "https://example.com/robots.txt",
			wantErr:           assert.NoError,
		},
		{
			name: "success with an user",
			args: args{
				regularLink: "http://username:password@example.com/test",
			},
			wantRobotsTXTLink: "http://username:password@example.com/robots.txt",
			wantErr:           assert.NoError,
		},
		{
			name: "success with a query",
			args: args{
				regularLink: "http://example.com/test?key=value",
			},
			wantRobotsTXTLink: "http://example.com/robots.txt",
			wantErr:           assert.NoError,
		},
		{
			name: "success with a fragment",
			args: args{
				regularLink: "http://example.com/test#fragment",
			},
			wantRobotsTXTLink: "http://example.com/robots.txt",
			wantErr:           assert.NoError,
		},
		{
			name: "error",
			args: args{
				regularLink: ":",
			},
			wantRobotsTXTLink: "",
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

func Test_loadRobotsTXTData(test *testing.T) {
	type args struct {
		ctx           context.Context
		httpClient    HTTPClient
		robotsTXTLink string
	}

	for _, data := range []struct {
		name              string
		args              args
		wantRobotsTXTData *robotstxt.RobotsData
		wantErr           assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	} {
		test.Run(data.name, func(test *testing.T) {
			gotRobotsTXTData, gotErr := loadRobotsTXTData(
				data.args.ctx,
				data.args.httpClient,
				data.args.robotsTXTLink,
			)

			mock.AssertExpectationsForObjects(test, data.args.httpClient)
			assert.Equal(test, data.wantRobotsTXTData, gotRobotsTXTData)
			data.wantErr(test, gotErr)
		})
	}
}
