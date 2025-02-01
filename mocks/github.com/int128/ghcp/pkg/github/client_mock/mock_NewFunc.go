// Code generated by mockery v2.52.1. DO NOT EDIT.

package client_mock

import (
	client "github.com/int128/ghcp/pkg/github/client"
	mock "github.com/stretchr/testify/mock"
)

// MockNewFunc is an autogenerated mock type for the NewFunc type
type MockNewFunc struct {
	mock.Mock
}

type MockNewFunc_Expecter struct {
	mock *mock.Mock
}

func (_m *MockNewFunc) EXPECT() *MockNewFunc_Expecter {
	return &MockNewFunc_Expecter{mock: &_m.Mock}
}

// Execute provides a mock function with given fields: _a0
func (_m *MockNewFunc) Execute(_a0 client.Option) (client.Interface, error) {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for Execute")
	}

	var r0 client.Interface
	var r1 error
	if rf, ok := ret.Get(0).(func(client.Option) (client.Interface, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(client.Option) client.Interface); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(client.Interface)
		}
	}

	if rf, ok := ret.Get(1).(func(client.Option) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockNewFunc_Execute_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Execute'
type MockNewFunc_Execute_Call struct {
	*mock.Call
}

// Execute is a helper method to define mock.On call
//   - _a0 client.Option
func (_e *MockNewFunc_Expecter) Execute(_a0 interface{}) *MockNewFunc_Execute_Call {
	return &MockNewFunc_Execute_Call{Call: _e.mock.On("Execute", _a0)}
}

func (_c *MockNewFunc_Execute_Call) Run(run func(_a0 client.Option)) *MockNewFunc_Execute_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(client.Option))
	})
	return _c
}

func (_c *MockNewFunc_Execute_Call) Return(_a0 client.Interface, _a1 error) *MockNewFunc_Execute_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockNewFunc_Execute_Call) RunAndReturn(run func(client.Option) (client.Interface, error)) *MockNewFunc_Execute_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockNewFunc creates a new instance of MockNewFunc. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockNewFunc(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockNewFunc {
	mock := &MockNewFunc{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
