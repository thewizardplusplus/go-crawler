// Code generated by mockery v1.0.0. DO NOT EDIT.

package handlers

import context "context"
import crawler "github.com/thewizardplusplus/go-crawler"
import mock "github.com/stretchr/testify/mock"

// MockLinkChecker is an autogenerated mock type for the LinkChecker type
type MockLinkChecker struct {
	mock.Mock
}

// CheckLink provides a mock function with given fields: ctx, link
func (_m *MockLinkChecker) CheckLink(ctx context.Context, link crawler.SourcedLink) bool {
	ret := _m.Called(ctx, link)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, crawler.SourcedLink) bool); ok {
		r0 = rf(ctx, link)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}
