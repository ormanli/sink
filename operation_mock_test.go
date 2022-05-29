// Code generated by MockGen. DO NOT EDIT.
// Source: sink_test.go

// Package sink_test is a generated GoMock package.
package sink_test

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockOperation is a mock of Operation interface.
type MockOperation[I, O any] struct {
	ctrl     *gomock.Controller
	recorder *MockOperationMockRecorder[I, O]
}

// MockOperationMockRecorder is the mock recorder for MockOperation.
type MockOperationMockRecorder[I, O any] struct {
	mock *MockOperation[I, O]
}

// NewMockOperation creates a new mock instance.
func NewMockOperation[I, O any](ctrl *gomock.Controller) *MockOperation[I, O] {
	mock := &MockOperation[I, O]{ctrl: ctrl}
	mock.recorder = &MockOperationMockRecorder[I, O]{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOperation[I, O]) EXPECT() *MockOperationMockRecorder[I, O] {
	return m.recorder
}

// Op mocks base method.
func (m *MockOperation[I, O]) Op(arg0 I) (O, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Op", arg0)
	ret0, _ := ret[0].(O)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Op indicates an expected call of Op.
func (mr *MockOperationMockRecorder[I, O]) Op(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Op", reflect.TypeOf((*MockOperation[I, O])(nil).Op), arg0)
}
