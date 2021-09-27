// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package extractors

import (
	http "net/http"

	mock "github.com/stretchr/testify/mock"
)

// MockLinkTransformer is an autogenerated mock type for the LinkTransformer type
type MockLinkTransformer struct {
	mock.Mock
}

// TransformLinks provides a mock function with given fields: links, response, responseContent
func (_m *MockLinkTransformer) TransformLinks(links []string, response *http.Response, responseContent []byte) ([]string, error) {
	ret := _m.Called(links, response, responseContent)

	var r0 []string
	if rf, ok := ret.Get(0).(func([]string, *http.Response, []byte) []string); ok {
		r0 = rf(links, response, responseContent)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]string, *http.Response, []byte) error); ok {
		r1 = rf(links, response, responseContent)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
