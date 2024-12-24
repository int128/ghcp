// Code generated by mockery v2.50.1. DO NOT EDIT.

package cmd_mock

import (
	cmd "github.com/int128/ghcp/pkg/cmd"
	client "github.com/int128/ghcp/pkg/github/client"

	logger "github.com/int128/ghcp/pkg/logger"

	mock "github.com/stretchr/testify/mock"
)

// MockNewInternalRunnerFunc is an autogenerated mock type for the NewInternalRunnerFunc type
type MockNewInternalRunnerFunc struct {
	mock.Mock
}

type MockNewInternalRunnerFunc_Expecter struct {
	mock *mock.Mock
}

func (_m *MockNewInternalRunnerFunc) EXPECT() *MockNewInternalRunnerFunc_Expecter {
	return &MockNewInternalRunnerFunc_Expecter{mock: &_m.Mock}
}

// Execute provides a mock function with given fields: _a0, _a1
func (_m *MockNewInternalRunnerFunc) Execute(_a0 logger.Interface, _a1 client.Interface) *cmd.InternalRunner {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for Execute")
	}

	var r0 *cmd.InternalRunner
	if rf, ok := ret.Get(0).(func(logger.Interface, client.Interface) *cmd.InternalRunner); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*cmd.InternalRunner)
		}
	}

	return r0
}

// MockNewInternalRunnerFunc_Execute_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Execute'
type MockNewInternalRunnerFunc_Execute_Call struct {
	*mock.Call
}

// Execute is a helper method to define mock.On call
//   - _a0 logger.Interface
//   - _a1 client.Interface
func (_e *MockNewInternalRunnerFunc_Expecter) Execute(_a0 interface{}, _a1 interface{}) *MockNewInternalRunnerFunc_Execute_Call {
	return &MockNewInternalRunnerFunc_Execute_Call{Call: _e.mock.On("Execute", _a0, _a1)}
}

func (_c *MockNewInternalRunnerFunc_Execute_Call) Run(run func(_a0 logger.Interface, _a1 client.Interface)) *MockNewInternalRunnerFunc_Execute_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(logger.Interface), args[1].(client.Interface))
	})
	return _c
}

func (_c *MockNewInternalRunnerFunc_Execute_Call) Return(_a0 *cmd.InternalRunner) *MockNewInternalRunnerFunc_Execute_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockNewInternalRunnerFunc_Execute_Call) RunAndReturn(run func(logger.Interface, client.Interface) *cmd.InternalRunner) *MockNewInternalRunnerFunc_Execute_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockNewInternalRunnerFunc creates a new instance of MockNewInternalRunnerFunc. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockNewInternalRunnerFunc(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockNewInternalRunnerFunc {
	mock := &MockNewInternalRunnerFunc{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
