package mocks

import (
	"context"
	"github.com/golang/mock/gomock"
	"reflect"
)

// MockDataBaser is a mock of DataBaser interface.
type MockDataBaser struct {
	ctrl     *gomock.Controller
	recorder *MockDataBaserMockRecorder
}

// MockDataBaserMockRecorder is the mock recorder for MockDataBaser.
type MockDataBaserMockRecorder struct {
	mock *MockDataBaser
}

// NewMockDataBaser creates a new mock instance.
func NewMockDataBaser(ctrl *gomock.Controller) *MockDataBaser {
	mock := &MockDataBaser{ctrl: ctrl}
	mock.recorder = &MockDataBaserMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDataBaser) EXPECT() *MockDataBaserMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockDataBaser) Create(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockDataBaserMockRecorder) Create(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockDataBaser)(nil).Create), arg0)
}

// Read mocks base method.
func (m *MockDataBaser) Read(arg0 context.Context) ([][]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Read", arg0)
	ret0, _ := ret[0].([][]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Read indicates an expected call of Read.
func (mr *MockDataBaserMockRecorder) Read(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Read", reflect.TypeOf((*MockDataBaser)(nil).Read), arg0)
}

// Write mocks base method.
func (m *MockDataBaser) Write(arg0 context.Context, arg1 [][]byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Write", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Write indicates an expected call of Write.
func (mr *MockDataBaserMockRecorder) Write(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Write", reflect.TypeOf((*MockDataBaser)(nil).Write), arg0, arg1)
}
