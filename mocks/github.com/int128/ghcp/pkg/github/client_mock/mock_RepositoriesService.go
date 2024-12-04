// Code generated by mockery v2.50.0. DO NOT EDIT.

package client_mock

import (
	context "context"

	github "github.com/google/go-github/v66/github"
	mock "github.com/stretchr/testify/mock"

	os "os"
)

// MockRepositoriesService is an autogenerated mock type for the RepositoriesService type
type MockRepositoriesService struct {
	mock.Mock
}

type MockRepositoriesService_Expecter struct {
	mock *mock.Mock
}

func (_m *MockRepositoriesService) EXPECT() *MockRepositoriesService_Expecter {
	return &MockRepositoriesService_Expecter{mock: &_m.Mock}
}

// CreateFork provides a mock function with given fields: ctx, owner, repo, opt
func (_m *MockRepositoriesService) CreateFork(ctx context.Context, owner string, repo string, opt *github.RepositoryCreateForkOptions) (*github.Repository, *github.Response, error) {
	ret := _m.Called(ctx, owner, repo, opt)

	if len(ret) == 0 {
		panic("no return value specified for CreateFork")
	}

	var r0 *github.Repository
	var r1 *github.Response
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, *github.RepositoryCreateForkOptions) (*github.Repository, *github.Response, error)); ok {
		return rf(ctx, owner, repo, opt)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string, *github.RepositoryCreateForkOptions) *github.Repository); ok {
		r0 = rf(ctx, owner, repo, opt)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*github.Repository)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string, *github.RepositoryCreateForkOptions) *github.Response); ok {
		r1 = rf(ctx, owner, repo, opt)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*github.Response)
		}
	}

	if rf, ok := ret.Get(2).(func(context.Context, string, string, *github.RepositoryCreateForkOptions) error); ok {
		r2 = rf(ctx, owner, repo, opt)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// MockRepositoriesService_CreateFork_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateFork'
type MockRepositoriesService_CreateFork_Call struct {
	*mock.Call
}

// CreateFork is a helper method to define mock.On call
//   - ctx context.Context
//   - owner string
//   - repo string
//   - opt *github.RepositoryCreateForkOptions
func (_e *MockRepositoriesService_Expecter) CreateFork(ctx interface{}, owner interface{}, repo interface{}, opt interface{}) *MockRepositoriesService_CreateFork_Call {
	return &MockRepositoriesService_CreateFork_Call{Call: _e.mock.On("CreateFork", ctx, owner, repo, opt)}
}

func (_c *MockRepositoriesService_CreateFork_Call) Run(run func(ctx context.Context, owner string, repo string, opt *github.RepositoryCreateForkOptions)) *MockRepositoriesService_CreateFork_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string), args[3].(*github.RepositoryCreateForkOptions))
	})
	return _c
}

func (_c *MockRepositoriesService_CreateFork_Call) Return(_a0 *github.Repository, _a1 *github.Response, _a2 error) *MockRepositoriesService_CreateFork_Call {
	_c.Call.Return(_a0, _a1, _a2)
	return _c
}

func (_c *MockRepositoriesService_CreateFork_Call) RunAndReturn(run func(context.Context, string, string, *github.RepositoryCreateForkOptions) (*github.Repository, *github.Response, error)) *MockRepositoriesService_CreateFork_Call {
	_c.Call.Return(run)
	return _c
}

// CreateRelease provides a mock function with given fields: ctx, owner, repo, release
func (_m *MockRepositoriesService) CreateRelease(ctx context.Context, owner string, repo string, release *github.RepositoryRelease) (*github.RepositoryRelease, *github.Response, error) {
	ret := _m.Called(ctx, owner, repo, release)

	if len(ret) == 0 {
		panic("no return value specified for CreateRelease")
	}

	var r0 *github.RepositoryRelease
	var r1 *github.Response
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, *github.RepositoryRelease) (*github.RepositoryRelease, *github.Response, error)); ok {
		return rf(ctx, owner, repo, release)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string, *github.RepositoryRelease) *github.RepositoryRelease); ok {
		r0 = rf(ctx, owner, repo, release)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*github.RepositoryRelease)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string, *github.RepositoryRelease) *github.Response); ok {
		r1 = rf(ctx, owner, repo, release)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*github.Response)
		}
	}

	if rf, ok := ret.Get(2).(func(context.Context, string, string, *github.RepositoryRelease) error); ok {
		r2 = rf(ctx, owner, repo, release)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// MockRepositoriesService_CreateRelease_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateRelease'
type MockRepositoriesService_CreateRelease_Call struct {
	*mock.Call
}

// CreateRelease is a helper method to define mock.On call
//   - ctx context.Context
//   - owner string
//   - repo string
//   - release *github.RepositoryRelease
func (_e *MockRepositoriesService_Expecter) CreateRelease(ctx interface{}, owner interface{}, repo interface{}, release interface{}) *MockRepositoriesService_CreateRelease_Call {
	return &MockRepositoriesService_CreateRelease_Call{Call: _e.mock.On("CreateRelease", ctx, owner, repo, release)}
}

func (_c *MockRepositoriesService_CreateRelease_Call) Run(run func(ctx context.Context, owner string, repo string, release *github.RepositoryRelease)) *MockRepositoriesService_CreateRelease_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string), args[3].(*github.RepositoryRelease))
	})
	return _c
}

func (_c *MockRepositoriesService_CreateRelease_Call) Return(_a0 *github.RepositoryRelease, _a1 *github.Response, _a2 error) *MockRepositoriesService_CreateRelease_Call {
	_c.Call.Return(_a0, _a1, _a2)
	return _c
}

func (_c *MockRepositoriesService_CreateRelease_Call) RunAndReturn(run func(context.Context, string, string, *github.RepositoryRelease) (*github.RepositoryRelease, *github.Response, error)) *MockRepositoriesService_CreateRelease_Call {
	_c.Call.Return(run)
	return _c
}

// GetReleaseByTag provides a mock function with given fields: ctx, owner, repo, tag
func (_m *MockRepositoriesService) GetReleaseByTag(ctx context.Context, owner string, repo string, tag string) (*github.RepositoryRelease, *github.Response, error) {
	ret := _m.Called(ctx, owner, repo, tag)

	if len(ret) == 0 {
		panic("no return value specified for GetReleaseByTag")
	}

	var r0 *github.RepositoryRelease
	var r1 *github.Response
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string) (*github.RepositoryRelease, *github.Response, error)); ok {
		return rf(ctx, owner, repo, tag)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string) *github.RepositoryRelease); ok {
		r0 = rf(ctx, owner, repo, tag)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*github.RepositoryRelease)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string, string) *github.Response); ok {
		r1 = rf(ctx, owner, repo, tag)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*github.Response)
		}
	}

	if rf, ok := ret.Get(2).(func(context.Context, string, string, string) error); ok {
		r2 = rf(ctx, owner, repo, tag)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// MockRepositoriesService_GetReleaseByTag_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetReleaseByTag'
type MockRepositoriesService_GetReleaseByTag_Call struct {
	*mock.Call
}

// GetReleaseByTag is a helper method to define mock.On call
//   - ctx context.Context
//   - owner string
//   - repo string
//   - tag string
func (_e *MockRepositoriesService_Expecter) GetReleaseByTag(ctx interface{}, owner interface{}, repo interface{}, tag interface{}) *MockRepositoriesService_GetReleaseByTag_Call {
	return &MockRepositoriesService_GetReleaseByTag_Call{Call: _e.mock.On("GetReleaseByTag", ctx, owner, repo, tag)}
}

func (_c *MockRepositoriesService_GetReleaseByTag_Call) Run(run func(ctx context.Context, owner string, repo string, tag string)) *MockRepositoriesService_GetReleaseByTag_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string), args[3].(string))
	})
	return _c
}

func (_c *MockRepositoriesService_GetReleaseByTag_Call) Return(_a0 *github.RepositoryRelease, _a1 *github.Response, _a2 error) *MockRepositoriesService_GetReleaseByTag_Call {
	_c.Call.Return(_a0, _a1, _a2)
	return _c
}

func (_c *MockRepositoriesService_GetReleaseByTag_Call) RunAndReturn(run func(context.Context, string, string, string) (*github.RepositoryRelease, *github.Response, error)) *MockRepositoriesService_GetReleaseByTag_Call {
	_c.Call.Return(run)
	return _c
}

// UploadReleaseAsset provides a mock function with given fields: ctx, owner, repo, id, opt, file
func (_m *MockRepositoriesService) UploadReleaseAsset(ctx context.Context, owner string, repo string, id int64, opt *github.UploadOptions, file *os.File) (*github.ReleaseAsset, *github.Response, error) {
	ret := _m.Called(ctx, owner, repo, id, opt, file)

	if len(ret) == 0 {
		panic("no return value specified for UploadReleaseAsset")
	}

	var r0 *github.ReleaseAsset
	var r1 *github.Response
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, int64, *github.UploadOptions, *os.File) (*github.ReleaseAsset, *github.Response, error)); ok {
		return rf(ctx, owner, repo, id, opt, file)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string, int64, *github.UploadOptions, *os.File) *github.ReleaseAsset); ok {
		r0 = rf(ctx, owner, repo, id, opt, file)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*github.ReleaseAsset)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string, int64, *github.UploadOptions, *os.File) *github.Response); ok {
		r1 = rf(ctx, owner, repo, id, opt, file)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*github.Response)
		}
	}

	if rf, ok := ret.Get(2).(func(context.Context, string, string, int64, *github.UploadOptions, *os.File) error); ok {
		r2 = rf(ctx, owner, repo, id, opt, file)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// MockRepositoriesService_UploadReleaseAsset_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UploadReleaseAsset'
type MockRepositoriesService_UploadReleaseAsset_Call struct {
	*mock.Call
}

// UploadReleaseAsset is a helper method to define mock.On call
//   - ctx context.Context
//   - owner string
//   - repo string
//   - id int64
//   - opt *github.UploadOptions
//   - file *os.File
func (_e *MockRepositoriesService_Expecter) UploadReleaseAsset(ctx interface{}, owner interface{}, repo interface{}, id interface{}, opt interface{}, file interface{}) *MockRepositoriesService_UploadReleaseAsset_Call {
	return &MockRepositoriesService_UploadReleaseAsset_Call{Call: _e.mock.On("UploadReleaseAsset", ctx, owner, repo, id, opt, file)}
}

func (_c *MockRepositoriesService_UploadReleaseAsset_Call) Run(run func(ctx context.Context, owner string, repo string, id int64, opt *github.UploadOptions, file *os.File)) *MockRepositoriesService_UploadReleaseAsset_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string), args[3].(int64), args[4].(*github.UploadOptions), args[5].(*os.File))
	})
	return _c
}

func (_c *MockRepositoriesService_UploadReleaseAsset_Call) Return(_a0 *github.ReleaseAsset, _a1 *github.Response, _a2 error) *MockRepositoriesService_UploadReleaseAsset_Call {
	_c.Call.Return(_a0, _a1, _a2)
	return _c
}

func (_c *MockRepositoriesService_UploadReleaseAsset_Call) RunAndReturn(run func(context.Context, string, string, int64, *github.UploadOptions, *os.File) (*github.ReleaseAsset, *github.Response, error)) *MockRepositoriesService_UploadReleaseAsset_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockRepositoriesService creates a new instance of MockRepositoriesService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockRepositoriesService(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockRepositoriesService {
	mock := &MockRepositoriesService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
