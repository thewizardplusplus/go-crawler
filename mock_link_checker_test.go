// Code generated by mockery v1.0.0. DO NOT EDIT.

package crawler

import mock "github.com/stretchr/testify/mock"

// MockLinkChecker is an autogenerated mock type for the LinkChecker type
type MockLinkChecker struct {
	mock.Mock
}

// CheckLink provides a mock function with given fields: link
func (_m *MockLinkChecker) CheckLink(link SourcedLink) bool {
	ret := _m.Called(link)

	var r0 bool
	if rf, ok := ret.Get(0).(func(SourcedLink) bool); ok {
		r0 = rf(link)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}
