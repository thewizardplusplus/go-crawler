// Code generated by mockery v1.0.0. DO NOT EDIT.

package handlers

import mock "github.com/stretchr/testify/mock"

// MockLinkHandler is an autogenerated mock type for the LinkHandler type
type MockLinkHandler struct {
	mock.Mock
}

// HandleLink provides a mock function with given fields: sourceLink, link
func (_m *MockLinkHandler) HandleLink(sourceLink string, link string) {
	_m.Called(sourceLink, link)
}
