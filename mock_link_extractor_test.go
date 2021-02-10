// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package crawler

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockLinkExtractor is an autogenerated mock type for the LinkExtractor type
type MockLinkExtractor struct {
	mock.Mock
}

// ExtractLinks provides a mock function with given fields: ctx, threadID, link
func (_m *MockLinkExtractor) ExtractLinks(ctx context.Context, threadID int, link string) ([]string, error) {
	ret := _m.Called(ctx, threadID, link)

	var r0 []string
	if rf, ok := ret.Get(0).(func(context.Context, int, string) []string); ok {
		r0 = rf(ctx, threadID, link)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int, string) error); ok {
		r1 = rf(ctx, threadID, link)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
