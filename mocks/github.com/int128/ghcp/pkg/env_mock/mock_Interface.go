// Code generated by mockery v2.51.0. DO NOT EDIT.

package env_mock

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

// Chdir provides a mock function with given fields: dir
func (_m *MockInterface) Chdir(dir string) error {
	ret := _m.Called(dir)

	if len(ret) == 0 {
		panic("no return value specified for Chdir")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(dir)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockInterface_Chdir_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Chdir'
type MockInterface_Chdir_Call struct {
	*mock.Call
}

// Chdir is a helper method to define mock.On call
//   - dir string
func (_e *MockInterface_Expecter) Chdir(dir interface{}) *MockInterface_Chdir_Call {
	return &MockInterface_Chdir_Call{Call: _e.mock.On("Chdir", dir)}
}

func (_c *MockInterface_Chdir_Call) Run(run func(dir string)) *MockInterface_Chdir_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockInterface_Chdir_Call) Return(_a0 error) *MockInterface_Chdir_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockInterface_Chdir_Call) RunAndReturn(run func(string) error) *MockInterface_Chdir_Call {
	_c.Call.Return(run)
	return _c
}

// Getenv provides a mock function with given fields: key
func (_m *MockInterface) Getenv(key string) string {
	ret := _m.Called(key)

	if len(ret) == 0 {
		panic("no return value specified for Getenv")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(key)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// MockInterface_Getenv_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Getenv'
type MockInterface_Getenv_Call struct {
	*mock.Call
}

// Getenv is a helper method to define mock.On call
//   - key string
func (_e *MockInterface_Expecter) Getenv(key interface{}) *MockInterface_Getenv_Call {
	return &MockInterface_Getenv_Call{Call: _e.mock.On("Getenv", key)}
}

func (_c *MockInterface_Getenv_Call) Run(run func(key string)) *MockInterface_Getenv_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockInterface_Getenv_Call) Return(_a0 string) *MockInterface_Getenv_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockInterface_Getenv_Call) RunAndReturn(run func(string) string) *MockInterface_Getenv_Call {
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
