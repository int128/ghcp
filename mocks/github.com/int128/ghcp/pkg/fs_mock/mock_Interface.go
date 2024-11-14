// Code generated by mockery v2.47.0. DO NOT EDIT.

package fs_mock

import (
	fs "github.com/int128/ghcp/pkg/fs"
	mock "github.com/stretchr/testify/mock"
)

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

// FindFiles provides a mock function with given fields: paths, filter
func (_m *MockInterface) FindFiles(paths []string, filter fs.FindFilesFilter) ([]fs.File, error) {
	ret := _m.Called(paths, filter)

	if len(ret) == 0 {
		panic("no return value specified for FindFiles")
	}

	var r0 []fs.File
	var r1 error
	if rf, ok := ret.Get(0).(func([]string, fs.FindFilesFilter) ([]fs.File, error)); ok {
		return rf(paths, filter)
	}
	if rf, ok := ret.Get(0).(func([]string, fs.FindFilesFilter) []fs.File); ok {
		r0 = rf(paths, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]fs.File)
		}
	}

	if rf, ok := ret.Get(1).(func([]string, fs.FindFilesFilter) error); ok {
		r1 = rf(paths, filter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockInterface_FindFiles_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindFiles'
type MockInterface_FindFiles_Call struct {
	*mock.Call
}

// FindFiles is a helper method to define mock.On call
//   - paths []string
//   - filter fs.FindFilesFilter
func (_e *MockInterface_Expecter) FindFiles(paths interface{}, filter interface{}) *MockInterface_FindFiles_Call {
	return &MockInterface_FindFiles_Call{Call: _e.mock.On("FindFiles", paths, filter)}
}

func (_c *MockInterface_FindFiles_Call) Run(run func(paths []string, filter fs.FindFilesFilter)) *MockInterface_FindFiles_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].([]string), args[1].(fs.FindFilesFilter))
	})
	return _c
}

func (_c *MockInterface_FindFiles_Call) Return(_a0 []fs.File, _a1 error) *MockInterface_FindFiles_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockInterface_FindFiles_Call) RunAndReturn(run func([]string, fs.FindFilesFilter) ([]fs.File, error)) *MockInterface_FindFiles_Call {
	_c.Call.Return(run)
	return _c
}

// ReadAsBase64EncodedContent provides a mock function with given fields: filename
func (_m *MockInterface) ReadAsBase64EncodedContent(filename string) (string, error) {
	ret := _m.Called(filename)

	if len(ret) == 0 {
		panic("no return value specified for ReadAsBase64EncodedContent")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (string, error)); ok {
		return rf(filename)
	}
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(filename)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(filename)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockInterface_ReadAsBase64EncodedContent_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ReadAsBase64EncodedContent'
type MockInterface_ReadAsBase64EncodedContent_Call struct {
	*mock.Call
}

// ReadAsBase64EncodedContent is a helper method to define mock.On call
//   - filename string
func (_e *MockInterface_Expecter) ReadAsBase64EncodedContent(filename interface{}) *MockInterface_ReadAsBase64EncodedContent_Call {
	return &MockInterface_ReadAsBase64EncodedContent_Call{Call: _e.mock.On("ReadAsBase64EncodedContent", filename)}
}

func (_c *MockInterface_ReadAsBase64EncodedContent_Call) Run(run func(filename string)) *MockInterface_ReadAsBase64EncodedContent_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockInterface_ReadAsBase64EncodedContent_Call) Return(_a0 string, _a1 error) *MockInterface_ReadAsBase64EncodedContent_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockInterface_ReadAsBase64EncodedContent_Call) RunAndReturn(run func(string) (string, error)) *MockInterface_ReadAsBase64EncodedContent_Call {
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
