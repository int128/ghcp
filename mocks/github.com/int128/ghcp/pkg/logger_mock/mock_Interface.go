// Code generated by mockery v2.46.0. DO NOT EDIT.

package logger_mock

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

// Debugf provides a mock function with given fields: format, v
func (_m *MockInterface) Debugf(format string, v ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, format)
	_ca = append(_ca, v...)
	_m.Called(_ca...)
}

// MockInterface_Debugf_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Debugf'
type MockInterface_Debugf_Call struct {
	*mock.Call
}

// Debugf is a helper method to define mock.On call
//   - format string
//   - v ...interface{}
func (_e *MockInterface_Expecter) Debugf(format interface{}, v ...interface{}) *MockInterface_Debugf_Call {
	return &MockInterface_Debugf_Call{Call: _e.mock.On("Debugf",
		append([]interface{}{format}, v...)...)}
}

func (_c *MockInterface_Debugf_Call) Run(run func(format string, v ...interface{})) *MockInterface_Debugf_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]interface{}, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(interface{})
			}
		}
		run(args[0].(string), variadicArgs...)
	})
	return _c
}

func (_c *MockInterface_Debugf_Call) Return() *MockInterface_Debugf_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockInterface_Debugf_Call) RunAndReturn(run func(string, ...interface{})) *MockInterface_Debugf_Call {
	_c.Call.Return(run)
	return _c
}

// Errorf provides a mock function with given fields: format, v
func (_m *MockInterface) Errorf(format string, v ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, format)
	_ca = append(_ca, v...)
	_m.Called(_ca...)
}

// MockInterface_Errorf_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Errorf'
type MockInterface_Errorf_Call struct {
	*mock.Call
}

// Errorf is a helper method to define mock.On call
//   - format string
//   - v ...interface{}
func (_e *MockInterface_Expecter) Errorf(format interface{}, v ...interface{}) *MockInterface_Errorf_Call {
	return &MockInterface_Errorf_Call{Call: _e.mock.On("Errorf",
		append([]interface{}{format}, v...)...)}
}

func (_c *MockInterface_Errorf_Call) Run(run func(format string, v ...interface{})) *MockInterface_Errorf_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]interface{}, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(interface{})
			}
		}
		run(args[0].(string), variadicArgs...)
	})
	return _c
}

func (_c *MockInterface_Errorf_Call) Return() *MockInterface_Errorf_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockInterface_Errorf_Call) RunAndReturn(run func(string, ...interface{})) *MockInterface_Errorf_Call {
	_c.Call.Return(run)
	return _c
}

// Infof provides a mock function with given fields: format, v
func (_m *MockInterface) Infof(format string, v ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, format)
	_ca = append(_ca, v...)
	_m.Called(_ca...)
}

// MockInterface_Infof_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Infof'
type MockInterface_Infof_Call struct {
	*mock.Call
}

// Infof is a helper method to define mock.On call
//   - format string
//   - v ...interface{}
func (_e *MockInterface_Expecter) Infof(format interface{}, v ...interface{}) *MockInterface_Infof_Call {
	return &MockInterface_Infof_Call{Call: _e.mock.On("Infof",
		append([]interface{}{format}, v...)...)}
}

func (_c *MockInterface_Infof_Call) Run(run func(format string, v ...interface{})) *MockInterface_Infof_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]interface{}, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(interface{})
			}
		}
		run(args[0].(string), variadicArgs...)
	})
	return _c
}

func (_c *MockInterface_Infof_Call) Return() *MockInterface_Infof_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockInterface_Infof_Call) RunAndReturn(run func(string, ...interface{})) *MockInterface_Infof_Call {
	_c.Call.Return(run)
	return _c
}

// Warnf provides a mock function with given fields: format, v
func (_m *MockInterface) Warnf(format string, v ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, format)
	_ca = append(_ca, v...)
	_m.Called(_ca...)
}

// MockInterface_Warnf_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Warnf'
type MockInterface_Warnf_Call struct {
	*mock.Call
}

// Warnf is a helper method to define mock.On call
//   - format string
//   - v ...interface{}
func (_e *MockInterface_Expecter) Warnf(format interface{}, v ...interface{}) *MockInterface_Warnf_Call {
	return &MockInterface_Warnf_Call{Call: _e.mock.On("Warnf",
		append([]interface{}{format}, v...)...)}
}

func (_c *MockInterface_Warnf_Call) Run(run func(format string, v ...interface{})) *MockInterface_Warnf_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]interface{}, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(interface{})
			}
		}
		run(args[0].(string), variadicArgs...)
	})
	return _c
}

func (_c *MockInterface_Warnf_Call) Return() *MockInterface_Warnf_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockInterface_Warnf_Call) RunAndReturn(run func(string, ...interface{})) *MockInterface_Warnf_Call {
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
