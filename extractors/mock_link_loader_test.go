// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package extractors

import mock "github.com/stretchr/testify/mock"

// MockLinkLoader is an autogenerated mock type for the LinkLoader type
type MockLinkLoader struct {
	mock.Mock
}

// LoadLink provides a mock function with given fields: link, options
func (_m *MockLinkLoader) LoadLink(link string, options interface{}) ([]byte, error) {
	ret := _m.Called(link, options)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(string, interface{}) []byte); ok {
		r0 = rf(link, options)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, interface{}) error); ok {
		r1 = rf(link, options)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
