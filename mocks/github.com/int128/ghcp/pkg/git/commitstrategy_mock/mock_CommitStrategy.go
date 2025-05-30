// Code generated by mockery v2.53.3. DO NOT EDIT.

package commitstrategy_mock

import (
	git "github.com/int128/ghcp/pkg/git"
	mock "github.com/stretchr/testify/mock"
)

// MockCommitStrategy is an autogenerated mock type for the CommitStrategy type
type MockCommitStrategy struct {
	mock.Mock
}

type MockCommitStrategy_Expecter struct {
	mock *mock.Mock
}

func (_m *MockCommitStrategy) EXPECT() *MockCommitStrategy_Expecter {
	return &MockCommitStrategy_Expecter{mock: &_m.Mock}
}

// IsFastForward provides a mock function with no fields
func (_m *MockCommitStrategy) IsFastForward() bool {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for IsFastForward")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// MockCommitStrategy_IsFastForward_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'IsFastForward'
type MockCommitStrategy_IsFastForward_Call struct {
	*mock.Call
}

// IsFastForward is a helper method to define mock.On call
func (_e *MockCommitStrategy_Expecter) IsFastForward() *MockCommitStrategy_IsFastForward_Call {
	return &MockCommitStrategy_IsFastForward_Call{Call: _e.mock.On("IsFastForward")}
}

func (_c *MockCommitStrategy_IsFastForward_Call) Run(run func()) *MockCommitStrategy_IsFastForward_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockCommitStrategy_IsFastForward_Call) Return(_a0 bool) *MockCommitStrategy_IsFastForward_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockCommitStrategy_IsFastForward_Call) RunAndReturn(run func() bool) *MockCommitStrategy_IsFastForward_Call {
	_c.Call.Return(run)
	return _c
}

// IsRebase provides a mock function with no fields
func (_m *MockCommitStrategy) IsRebase() bool {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for IsRebase")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// MockCommitStrategy_IsRebase_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'IsRebase'
type MockCommitStrategy_IsRebase_Call struct {
	*mock.Call
}

// IsRebase is a helper method to define mock.On call
func (_e *MockCommitStrategy_Expecter) IsRebase() *MockCommitStrategy_IsRebase_Call {
	return &MockCommitStrategy_IsRebase_Call{Call: _e.mock.On("IsRebase")}
}

func (_c *MockCommitStrategy_IsRebase_Call) Run(run func()) *MockCommitStrategy_IsRebase_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockCommitStrategy_IsRebase_Call) Return(_a0 bool) *MockCommitStrategy_IsRebase_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockCommitStrategy_IsRebase_Call) RunAndReturn(run func() bool) *MockCommitStrategy_IsRebase_Call {
	_c.Call.Return(run)
	return _c
}

// NoParent provides a mock function with no fields
func (_m *MockCommitStrategy) NoParent() bool {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for NoParent")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// MockCommitStrategy_NoParent_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'NoParent'
type MockCommitStrategy_NoParent_Call struct {
	*mock.Call
}

// NoParent is a helper method to define mock.On call
func (_e *MockCommitStrategy_Expecter) NoParent() *MockCommitStrategy_NoParent_Call {
	return &MockCommitStrategy_NoParent_Call{Call: _e.mock.On("NoParent")}
}

func (_c *MockCommitStrategy_NoParent_Call) Run(run func()) *MockCommitStrategy_NoParent_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockCommitStrategy_NoParent_Call) Return(_a0 bool) *MockCommitStrategy_NoParent_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockCommitStrategy_NoParent_Call) RunAndReturn(run func() bool) *MockCommitStrategy_NoParent_Call {
	_c.Call.Return(run)
	return _c
}

// RebaseUpstream provides a mock function with no fields
func (_m *MockCommitStrategy) RebaseUpstream() git.RefName {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for RebaseUpstream")
	}

	var r0 git.RefName
	if rf, ok := ret.Get(0).(func() git.RefName); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(git.RefName)
	}

	return r0
}

// MockCommitStrategy_RebaseUpstream_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RebaseUpstream'
type MockCommitStrategy_RebaseUpstream_Call struct {
	*mock.Call
}

// RebaseUpstream is a helper method to define mock.On call
func (_e *MockCommitStrategy_Expecter) RebaseUpstream() *MockCommitStrategy_RebaseUpstream_Call {
	return &MockCommitStrategy_RebaseUpstream_Call{Call: _e.mock.On("RebaseUpstream")}
}

func (_c *MockCommitStrategy_RebaseUpstream_Call) Run(run func()) *MockCommitStrategy_RebaseUpstream_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockCommitStrategy_RebaseUpstream_Call) Return(_a0 git.RefName) *MockCommitStrategy_RebaseUpstream_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockCommitStrategy_RebaseUpstream_Call) RunAndReturn(run func() git.RefName) *MockCommitStrategy_RebaseUpstream_Call {
	_c.Call.Return(run)
	return _c
}

// String provides a mock function with no fields
func (_m *MockCommitStrategy) String() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for String")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// MockCommitStrategy_String_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'String'
type MockCommitStrategy_String_Call struct {
	*mock.Call
}

// String is a helper method to define mock.On call
func (_e *MockCommitStrategy_Expecter) String() *MockCommitStrategy_String_Call {
	return &MockCommitStrategy_String_Call{Call: _e.mock.On("String")}
}

func (_c *MockCommitStrategy_String_Call) Run(run func()) *MockCommitStrategy_String_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockCommitStrategy_String_Call) Return(_a0 string) *MockCommitStrategy_String_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockCommitStrategy_String_Call) RunAndReturn(run func() string) *MockCommitStrategy_String_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockCommitStrategy creates a new instance of MockCommitStrategy. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockCommitStrategy(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockCommitStrategy {
	mock := &MockCommitStrategy{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
