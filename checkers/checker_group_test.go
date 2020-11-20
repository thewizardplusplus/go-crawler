package checkers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	crawler "github.com/thewizardplusplus/go-crawler"
)

func TestCheckerGroup_CheckLink(test *testing.T) {
	type args struct {
		link crawler.SourcedLink
	}

	for _, data := range []struct {
		name     string
		checkers CheckerGroup
		args     args
		want     assert.BoolAssertionFunc
	}{
		{
			name:     "empty",
			checkers: nil,
			args: args{
				link: crawler.SourcedLink{
					SourceLink: "http://example.com/",
					Link:       "http://example.com/test",
				},
			},
			want: assert.False,
		},
		{
			name: "without failed checkings",
			checkers: CheckerGroup{
				func() LinkChecker {
					checker := new(MockLinkChecker)
					checker.
						On("CheckLink", crawler.SourcedLink{
							SourceLink: "http://example.com/",
							Link:       "http://example.com/test",
						}).
						Return(true)

					return checker
				}(),
				func() LinkChecker {
					checker := new(MockLinkChecker)
					checker.
						On("CheckLink", crawler.SourcedLink{
							SourceLink: "http://example.com/",
							Link:       "http://example.com/test",
						}).
						Return(true)

					return checker
				}(),
			},
			args: args{
				link: crawler.SourcedLink{
					SourceLink: "http://example.com/",
					Link:       "http://example.com/test",
				},
			},
			want: assert.True,
		},
		{
			name: "with a failed checking",
			checkers: CheckerGroup{
				func() LinkChecker {
					checker := new(MockLinkChecker)
					checker.
						On("CheckLink", crawler.SourcedLink{
							SourceLink: "http://example.com/",
							Link:       "http://example.com/test",
						}).
						Return(false)

					return checker
				}(),
				new(MockLinkChecker),
			},
			args: args{
				link: crawler.SourcedLink{
					SourceLink: "http://example.com/",
					Link:       "http://example.com/test",
				},
			},
			want: assert.False,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			got := data.checkers.CheckLink(data.args.link)

			for _, checker := range data.checkers {
				mock.AssertExpectationsForObjects(test, checker)
			}
			data.want(test, got)
		})
	}
}
