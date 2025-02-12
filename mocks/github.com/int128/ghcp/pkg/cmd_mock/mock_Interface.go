// Code generated by mockery v2.52.2. DO NOT EDIT.

package cmd_mock

import mock "github.com/stretchr/testify/mock"

// MockInterface is an autogenerated mock type for the Interface type
type MockInterface struct {
	mock.Mock
}

type MockInterface_Expecter struct {
	mock *mock.Mock
}

func (_m *MockInterface) EXPECT() *MockInterface_Expecter {
	return &MockInterface_Expecter{mock: &_m.Mock}
}

// Run provides a mock function with given fields: args, version
func (_m *MockInterface) Run(args []string, version string) int {
	ret := _m.Called(args, version)

	if len(ret) == 0 {
		panic("no return value specified for Run")
	}

	var r0 int
	if rf, ok := ret.Get(0).(func([]string, string) int); ok {
		r0 = rf(args, version)
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// MockInterface_Run_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Run'
type MockInterface_Run_Call struct {
	*mock.Call
}

// Run is a helper method to define mock.On call
//   - args []string
//   - version string
func (_e *MockInterface_Expecter) Run(args interface{}, version interface{}) *MockInterface_Run_Call {
	return &MockInterface_Run_Call{Call: _e.mock.On("Run", args, version)}
}

func (_c *MockInterface_Run_Call) Run(run func(args []string, version string)) *MockInterface_Run_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].([]string), args[1].(string))
	})
	return _c
}

func (_c *MockInterface_Run_Call) Return(_a0 int) *MockInterface_Run_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockInterface_Run_Call) RunAndReturn(run func([]string, string) int) *MockInterface_Run_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockInterface creates a new instance of MockInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockInterface {
	mock := &MockInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
