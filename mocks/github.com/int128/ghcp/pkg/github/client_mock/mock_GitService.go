// Code generated by mockery v2.49.2. DO NOT EDIT.

package client_mock

import (
	context "context"

	github "github.com/google/go-github/v66/github"
	mock "github.com/stretchr/testify/mock"
)

// MockGitService is an autogenerated mock type for the GitService type
type MockGitService struct {
	mock.Mock
}

type MockGitService_Expecter struct {
	mock *mock.Mock
}

func (_m *MockGitService) EXPECT() *MockGitService_Expecter {
	return &MockGitService_Expecter{mock: &_m.Mock}
}

// CreateBlob provides a mock function with given fields: ctx, owner, repo, blob
func (_m *MockGitService) CreateBlob(ctx context.Context, owner string, repo string, blob *github.Blob) (*github.Blob, *github.Response, error) {
	ret := _m.Called(ctx, owner, repo, blob)

	if len(ret) == 0 {
		panic("no return value specified for CreateBlob")
	}

	var r0 *github.Blob
	var r1 *github.Response
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, *github.Blob) (*github.Blob, *github.Response, error)); ok {
		return rf(ctx, owner, repo, blob)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string, *github.Blob) *github.Blob); ok {
		r0 = rf(ctx, owner, repo, blob)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*github.Blob)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string, *github.Blob) *github.Response); ok {
		r1 = rf(ctx, owner, repo, blob)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*github.Response)
		}
	}

	if rf, ok := ret.Get(2).(func(context.Context, string, string, *github.Blob) error); ok {
		r2 = rf(ctx, owner, repo, blob)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// MockGitService_CreateBlob_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateBlob'
type MockGitService_CreateBlob_Call struct {
	*mock.Call
}

// CreateBlob is a helper method to define mock.On call
//   - ctx context.Context
//   - owner string
//   - repo string
//   - blob *github.Blob
func (_e *MockGitService_Expecter) CreateBlob(ctx interface{}, owner interface{}, repo interface{}, blob interface{}) *MockGitService_CreateBlob_Call {
	return &MockGitService_CreateBlob_Call{Call: _e.mock.On("CreateBlob", ctx, owner, repo, blob)}
}

func (_c *MockGitService_CreateBlob_Call) Run(run func(ctx context.Context, owner string, repo string, blob *github.Blob)) *MockGitService_CreateBlob_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string), args[3].(*github.Blob))
	})
	return _c
}

func (_c *MockGitService_CreateBlob_Call) Return(_a0 *github.Blob, _a1 *github.Response, _a2 error) *MockGitService_CreateBlob_Call {
	_c.Call.Return(_a0, _a1, _a2)
	return _c
}

func (_c *MockGitService_CreateBlob_Call) RunAndReturn(run func(context.Context, string, string, *github.Blob) (*github.Blob, *github.Response, error)) *MockGitService_CreateBlob_Call {
	_c.Call.Return(run)
	return _c
}

// CreateCommit provides a mock function with given fields: ctx, owner, repo, commit, opts
func (_m *MockGitService) CreateCommit(ctx context.Context, owner string, repo string, commit *github.Commit, opts *github.CreateCommitOptions) (*github.Commit, *github.Response, error) {
	ret := _m.Called(ctx, owner, repo, commit, opts)

	if len(ret) == 0 {
		panic("no return value specified for CreateCommit")
	}

	var r0 *github.Commit
	var r1 *github.Response
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, *github.Commit, *github.CreateCommitOptions) (*github.Commit, *github.Response, error)); ok {
		return rf(ctx, owner, repo, commit, opts)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string, *github.Commit, *github.CreateCommitOptions) *github.Commit); ok {
		r0 = rf(ctx, owner, repo, commit, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*github.Commit)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string, *github.Commit, *github.CreateCommitOptions) *github.Response); ok {
		r1 = rf(ctx, owner, repo, commit, opts)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*github.Response)
		}
	}

	if rf, ok := ret.Get(2).(func(context.Context, string, string, *github.Commit, *github.CreateCommitOptions) error); ok {
		r2 = rf(ctx, owner, repo, commit, opts)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// MockGitService_CreateCommit_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateCommit'
type MockGitService_CreateCommit_Call struct {
	*mock.Call
}

// CreateCommit is a helper method to define mock.On call
//   - ctx context.Context
//   - owner string
//   - repo string
//   - commit *github.Commit
//   - opts *github.CreateCommitOptions
func (_e *MockGitService_Expecter) CreateCommit(ctx interface{}, owner interface{}, repo interface{}, commit interface{}, opts interface{}) *MockGitService_CreateCommit_Call {
	return &MockGitService_CreateCommit_Call{Call: _e.mock.On("CreateCommit", ctx, owner, repo, commit, opts)}
}

func (_c *MockGitService_CreateCommit_Call) Run(run func(ctx context.Context, owner string, repo string, commit *github.Commit, opts *github.CreateCommitOptions)) *MockGitService_CreateCommit_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string), args[3].(*github.Commit), args[4].(*github.CreateCommitOptions))
	})
	return _c
}

func (_c *MockGitService_CreateCommit_Call) Return(_a0 *github.Commit, _a1 *github.Response, _a2 error) *MockGitService_CreateCommit_Call {
	_c.Call.Return(_a0, _a1, _a2)
	return _c
}

func (_c *MockGitService_CreateCommit_Call) RunAndReturn(run func(context.Context, string, string, *github.Commit, *github.CreateCommitOptions) (*github.Commit, *github.Response, error)) *MockGitService_CreateCommit_Call {
	_c.Call.Return(run)
	return _c
}

// CreateTree provides a mock function with given fields: ctx, owner, repo, baseTree, entries
func (_m *MockGitService) CreateTree(ctx context.Context, owner string, repo string, baseTree string, entries []*github.TreeEntry) (*github.Tree, *github.Response, error) {
	ret := _m.Called(ctx, owner, repo, baseTree, entries)

	if len(ret) == 0 {
		panic("no return value specified for CreateTree")
	}

	var r0 *github.Tree
	var r1 *github.Response
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string, []*github.TreeEntry) (*github.Tree, *github.Response, error)); ok {
		return rf(ctx, owner, repo, baseTree, entries)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string, []*github.TreeEntry) *github.Tree); ok {
		r0 = rf(ctx, owner, repo, baseTree, entries)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*github.Tree)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string, string, []*github.TreeEntry) *github.Response); ok {
		r1 = rf(ctx, owner, repo, baseTree, entries)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*github.Response)
		}
	}

	if rf, ok := ret.Get(2).(func(context.Context, string, string, string, []*github.TreeEntry) error); ok {
		r2 = rf(ctx, owner, repo, baseTree, entries)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// MockGitService_CreateTree_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateTree'
type MockGitService_CreateTree_Call struct {
	*mock.Call
}

// CreateTree is a helper method to define mock.On call
//   - ctx context.Context
//   - owner string
//   - repo string
//   - baseTree string
//   - entries []*github.TreeEntry
func (_e *MockGitService_Expecter) CreateTree(ctx interface{}, owner interface{}, repo interface{}, baseTree interface{}, entries interface{}) *MockGitService_CreateTree_Call {
	return &MockGitService_CreateTree_Call{Call: _e.mock.On("CreateTree", ctx, owner, repo, baseTree, entries)}
}

func (_c *MockGitService_CreateTree_Call) Run(run func(ctx context.Context, owner string, repo string, baseTree string, entries []*github.TreeEntry)) *MockGitService_CreateTree_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string), args[3].(string), args[4].([]*github.TreeEntry))
	})
	return _c
}

func (_c *MockGitService_CreateTree_Call) Return(_a0 *github.Tree, _a1 *github.Response, _a2 error) *MockGitService_CreateTree_Call {
	_c.Call.Return(_a0, _a1, _a2)
	return _c
}

func (_c *MockGitService_CreateTree_Call) RunAndReturn(run func(context.Context, string, string, string, []*github.TreeEntry) (*github.Tree, *github.Response, error)) *MockGitService_CreateTree_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockGitService creates a new instance of MockGitService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockGitService(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockGitService {
	mock := &MockGitService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
