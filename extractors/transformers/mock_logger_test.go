// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package transformers

import mock "github.com/stretchr/testify/mock"

// MockLogger is an autogenerated mock type for the Logger type
type MockLogger struct {
	mock.Mock
}

// Log provides a mock function with given fields: v
func (_m *MockLogger) Log(v ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, v...)
	_m.Called(_ca...)
}

// Logf provides a mock function with given fields: format, v
func (_m *MockLogger) Logf(format string, v ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, format)
	_ca = append(_ca, v...)
	_m.Called(_ca...)
}
