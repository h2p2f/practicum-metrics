// Code generated by mockery v2.30.1. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// Pinger is an autogenerated mock type for the Pinger type
type Pinger struct {
	mock.Mock
}

// Ping provides a mock function with given fields:
func (_m *Pinger) Ping() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewPinger creates a new instance of Pinger. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewPinger(t interface {
	mock.TestingT
	Cleanup(func())
}) *Pinger {
	mock := &Pinger{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
