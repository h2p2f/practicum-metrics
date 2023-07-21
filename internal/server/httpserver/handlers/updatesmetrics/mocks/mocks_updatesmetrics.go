// Code generated by mockery v2.30.1. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// Updater is an autogenerated mock type for the Updater type
type Updater struct {
	mock.Mock
}

// SetCounter provides a mock function with given fields: name, value
func (_m *Updater) SetCounter(name string, value int64) {
	_m.Called(name, value)
}

// SetGauge provides a mock function with given fields: name, value
func (_m *Updater) SetGauge(name string, value float64) {
	_m.Called(name, value)
}

// NewUpdater creates a new instance of Updater. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUpdater(t interface {
	mock.TestingT
	Cleanup(func())
}) *Updater {
	mock := &Updater{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}