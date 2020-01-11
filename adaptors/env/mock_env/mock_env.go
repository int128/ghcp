// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/int128/ghcp/adaptors/env (interfaces: Interface)

// Package mock_env is a generated GoMock package.
package mock_env

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockInterface is a mock of Interface interface
type MockInterface struct {
	ctrl     *gomock.Controller
	recorder *MockInterfaceMockRecorder
}

// MockInterfaceMockRecorder is the mock recorder for MockInterface
type MockInterfaceMockRecorder struct {
	mock *MockInterface
}

// NewMockInterface creates a new mock instance
func NewMockInterface(ctrl *gomock.Controller) *MockInterface {
	mock := &MockInterface{ctrl: ctrl}
	mock.recorder = &MockInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockInterface) EXPECT() *MockInterfaceMockRecorder {
	return m.recorder
}

// Chdir mocks base method
func (m *MockInterface) Chdir(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Chdir", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Chdir indicates an expected call of Chdir
func (mr *MockInterfaceMockRecorder) Chdir(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Chdir", reflect.TypeOf((*MockInterface)(nil).Chdir), arg0)
}

// Getenv mocks base method
func (m *MockInterface) Getenv(arg0 string) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Getenv", arg0)
	ret0, _ := ret[0].(string)
	return ret0
}

// Getenv indicates an expected call of Getenv
func (mr *MockInterfaceMockRecorder) Getenv(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Getenv", reflect.TypeOf((*MockInterface)(nil).Getenv), arg0)
}
