package handlers

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/go-log/log"
	"github.com/stretchr/testify/mock"
	crawler "github.com/thewizardplusplus/go-crawler"
	"github.com/thewizardplusplus/go-crawler/registers"
)

func TestRobotsTXTHandler_HandleLink(test *testing.T) {
	type fields struct {
		UserAgent         string
		RobotsTXTRegister registers.RobotsTXTRegister
		LinkHandler       crawler.LinkHandler
		Logger            log.Logger
	}
	type args struct {
		ctx  context.Context
		link crawler.SourcedLink
	}

	for _, data := range []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "success with an allowed link",
			fields: fields{
				UserAgent: "go-crawler",
				RobotsTXTRegister: func() registers.RobotsTXTRegister {
					request, _ :=
						http.NewRequest(http.MethodGet, "http://example.com/robots.txt", nil)
					request = request.WithContext(context.Background())

					response := &http.Response{
						StatusCode: http.StatusOK,
						Body: ioutil.NopCloser(strings.NewReader(`
							User-agent: *
							Disallow: /

							User-agent: go-crawler
							Disallow: /
							Allow: /$
							Allow: /sitemap.xml$
							Allow: /post/
							Allow: /storage/app/media/
						`)),
					}

					httpClient := new(MockHTTPClient)
					httpClient.On("Do", request).Return(response, nil)

					register := registers.NewRobotsTXTRegister(httpClient)
					return register
				}(),
				LinkHandler: func() LinkHandler {
					handler := new(MockLinkHandler)
					handler.
						On("HandleLink", context.Background(), crawler.SourcedLink{
							SourceLink: "http://example.com/",
							Link:       "http://example.com/post/23",
						}).
						Return()

					return handler
				}(),
				Logger: new(MockLogger),
			},
			args: args{
				ctx: context.Background(),
				link: crawler.SourcedLink{
					SourceLink: "http://example.com/",
					Link:       "http://example.com/post/23",
				},
			},
		},
		{
			name: "success with a link disallowed by a user agent",
			fields: fields{
				UserAgent: "go-crawler",
				RobotsTXTRegister: func() registers.RobotsTXTRegister {
					request, _ :=
						http.NewRequest(http.MethodGet, "http://example.com/robots.txt", nil)
					request = request.WithContext(context.Background())

					response := &http.Response{
						StatusCode: http.StatusOK,
						Body: ioutil.NopCloser(strings.NewReader(`
							User-agent: *
							Disallow: /

							User-agent: test
							Disallow: /
							Allow: /$
							Allow: /sitemap.xml$
							Allow: /post/
							Allow: /storage/app/media/
						`)),
					}

					httpClient := new(MockHTTPClient)
					httpClient.On("Do", request).Return(response, nil)

					register := registers.NewRobotsTXTRegister(httpClient)
					return register
				}(),
				LinkHandler: new(MockLinkHandler),
				Logger:      new(MockLogger),
			},
			args: args{
				ctx: context.Background(),
				link: crawler.SourcedLink{
					SourceLink: "http://example.com/",
					Link:       "http://example.com/post/23",
				},
			},
		},
		{
			name: "success with a link disallowed by a path",
			fields: fields{
				UserAgent: "go-crawler",
				RobotsTXTRegister: func() registers.RobotsTXTRegister {
					request, _ :=
						http.NewRequest(http.MethodGet, "http://example.com/robots.txt", nil)
					request = request.WithContext(context.Background())

					response := &http.Response{
						StatusCode: http.StatusOK,
						Body: ioutil.NopCloser(strings.NewReader(`
							User-agent: *
							Disallow: /

							User-agent: go-crawler
							Disallow: /
							Allow: /$
							Allow: /sitemap.xml$
							Allow: /post/
							Allow: /storage/app/media/
						`)),
					}

					httpClient := new(MockHTTPClient)
					httpClient.On("Do", request).Return(response, nil)

					register := registers.NewRobotsTXTRegister(httpClient)
					return register
				}(),
				LinkHandler: new(MockLinkHandler),
				Logger:      new(MockLogger),
			},
			args: args{
				ctx: context.Background(),
				link: crawler.SourcedLink{
					SourceLink: "http://example.com/",
					Link:       "http://example.com/test",
				},
			},
		},
		{
			name: "error with link parsing",
			fields: fields{
				UserAgent: "go-crawler",
				RobotsTXTRegister: func() registers.RobotsTXTRegister {
					httpClient := new(MockHTTPClient)
					register := registers.NewRobotsTXTRegister(httpClient)
					return register
				}(),
				LinkHandler: new(MockLinkHandler),
				Logger: func() Logger {
					err := errors.New("missing protocol scheme")
					urlErr := &url.Error{Op: "parse", URL: ":", Err: err}

					logger := new(MockLogger)
					logger.On("Logf", "unable to parse the link: %s", urlErr).Return()

					return logger
				}(),
			},
			args: args{
				ctx: context.Background(),
				link: crawler.SourcedLink{
					SourceLink: "http://example.com/",
					Link:       ":",
				},
			},
		},
		{
			name: "error with link registering",
			fields: fields{
				UserAgent: "go-crawler",
				RobotsTXTRegister: func() registers.RobotsTXTRegister {
					request, _ :=
						http.NewRequest(http.MethodGet, "http://example.com/robots.txt", nil)
					request = request.WithContext(context.Background())

					httpClient := new(MockHTTPClient)
					httpClient.On("Do", request).Return(nil, iotest.ErrTimeout)

					register := registers.NewRobotsTXTRegister(httpClient)
					return register
				}(),
				LinkHandler: new(MockLinkHandler),
				Logger: func() Logger {
					errMatcher := mock.MatchedBy(func(err error) bool {
						wantErrMessage := "unable to load the robots.txt data: " +
							"unable to send the request: " +
							"timeout"
						return err.Error() == wantErrMessage
					})

					logger := new(MockLogger)
					logger.
						On("Logf", "unable to register the robots.txt link: %s", errMatcher).
						Return()

					return logger
				}(),
			},
			args: args{
				ctx: context.Background(),
				link: crawler.SourcedLink{
					SourceLink: "http://example.com/",
					Link:       "http://example.com/post/23",
				},
			},
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			handler := RobotsTXTHandler{
				UserAgent:         data.fields.UserAgent,
				RobotsTXTRegister: data.fields.RobotsTXTRegister,
				LinkHandler:       data.fields.LinkHandler,
				Logger:            data.fields.Logger,
			}
			handler.HandleLink(data.args.ctx, data.args.link)

			mock.AssertExpectationsForObjects(
				test,
				data.fields.LinkHandler,
				data.fields.Logger,
			)
		})
	}
}
