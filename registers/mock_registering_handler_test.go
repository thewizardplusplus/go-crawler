// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package registers

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockRegisteringHandler is an autogenerated mock type for the RegisteringHandler type
type MockRegisteringHandler struct {
	mock.Mock
}

// HandleRegistering provides a mock function with given fields: ctx, key
func (_m *MockRegisteringHandler) HandleRegistering(ctx context.Context, key interface{}) (interface{}, error) {
	ret := _m.Called(ctx, key)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(context.Context, interface{}) interface{}); ok {
		r0 = rf(ctx, key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, interface{}) error); ok {
		r1 = rf(ctx, key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
