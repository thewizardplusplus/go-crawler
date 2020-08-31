package checkers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCheckerGroup_CheckLink(test *testing.T) {
	type args struct {
		sourceLink string
		link       string
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
				sourceLink: "http://example.com/",
				link:       "http://example.com/test",
			},
			want: assert.False,
		},
		{
			name: "without failed checkings",
			checkers: CheckerGroup{
				func() LinkChecker {
					checker := new(MockLinkChecker)
					checker.
						On("CheckLink", "http://example.com/", "http://example.com/test").
						Return(true)

					return checker
				}(),
				func() LinkChecker {
					checker := new(MockLinkChecker)
					checker.
						On("CheckLink", "http://example.com/", "http://example.com/test").
						Return(true)

					return checker
				}(),
			},
			args: args{
				sourceLink: "http://example.com/",
				link:       "http://example.com/test",
			},
			want: assert.True,
		},
		{
			name: "with a failed checking",
			checkers: CheckerGroup{
				func() LinkChecker {
					checker := new(MockLinkChecker)
					checker.
						On("CheckLink", "http://example.com/", "http://example.com/test").
						Return(false)

					return checker
				}(),
				new(MockLinkChecker),
			},
			args: args{
				sourceLink: "http://example.com/",
				link:       "http://example.com/test",
			},
			want: assert.False,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			got := data.checkers.CheckLink(data.args.sourceLink, data.args.link)

			for _, checker := range data.checkers {
				mock.AssertExpectationsForObjects(test, checker)
			}
			data.want(test, got)
		})
	}
}
