// Code generated by mockery v1.0.0. DO NOT EDIT.

package crawler

import mock "github.com/stretchr/testify/mock"

// MockLinkChecker is an autogenerated mock type for the LinkChecker type
type MockLinkChecker struct {
	mock.Mock
}

// CheckLink provides a mock function with given fields: parentLink, link
func (_m *MockLinkChecker) CheckLink(parentLink string, link string) bool {
	ret := _m.Called(parentLink, link)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string, string) bool); ok {
		r0 = rf(parentLink, link)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}
