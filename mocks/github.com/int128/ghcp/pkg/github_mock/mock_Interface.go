// Code generated by mockery v2.49.1. DO NOT EDIT.

package github_mock

import (
	context "context"

	git "github.com/int128/ghcp/pkg/git"
	github "github.com/int128/ghcp/pkg/github"

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

// CreateBlob provides a mock function with given fields: ctx, blob
func (_m *MockInterface) CreateBlob(ctx context.Context, blob git.NewBlob) (git.BlobSHA, error) {
	ret := _m.Called(ctx, blob)

	if len(ret) == 0 {
		panic("no return value specified for CreateBlob")
	}

	var r0 git.BlobSHA
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, git.NewBlob) (git.BlobSHA, error)); ok {
		return rf(ctx, blob)
	}
	if rf, ok := ret.Get(0).(func(context.Context, git.NewBlob) git.BlobSHA); ok {
		r0 = rf(ctx, blob)
	} else {
		r0 = ret.Get(0).(git.BlobSHA)
	}

	if rf, ok := ret.Get(1).(func(context.Context, git.NewBlob) error); ok {
		r1 = rf(ctx, blob)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockInterface_CreateBlob_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateBlob'
type MockInterface_CreateBlob_Call struct {
	*mock.Call
}

// CreateBlob is a helper method to define mock.On call
//   - ctx context.Context
//   - blob git.NewBlob
func (_e *MockInterface_Expecter) CreateBlob(ctx interface{}, blob interface{}) *MockInterface_CreateBlob_Call {
	return &MockInterface_CreateBlob_Call{Call: _e.mock.On("CreateBlob", ctx, blob)}
}

func (_c *MockInterface_CreateBlob_Call) Run(run func(ctx context.Context, blob git.NewBlob)) *MockInterface_CreateBlob_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(git.NewBlob))
	})
	return _c
}

func (_c *MockInterface_CreateBlob_Call) Return(_a0 git.BlobSHA, _a1 error) *MockInterface_CreateBlob_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockInterface_CreateBlob_Call) RunAndReturn(run func(context.Context, git.NewBlob) (git.BlobSHA, error)) *MockInterface_CreateBlob_Call {
	_c.Call.Return(run)
	return _c
}

// CreateBranch provides a mock function with given fields: ctx, in
func (_m *MockInterface) CreateBranch(ctx context.Context, in github.CreateBranchInput) error {
	ret := _m.Called(ctx, in)

	if len(ret) == 0 {
		panic("no return value specified for CreateBranch")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, github.CreateBranchInput) error); ok {
		r0 = rf(ctx, in)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockInterface_CreateBranch_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateBranch'
type MockInterface_CreateBranch_Call struct {
	*mock.Call
}

// CreateBranch is a helper method to define mock.On call
//   - ctx context.Context
//   - in github.CreateBranchInput
func (_e *MockInterface_Expecter) CreateBranch(ctx interface{}, in interface{}) *MockInterface_CreateBranch_Call {
	return &MockInterface_CreateBranch_Call{Call: _e.mock.On("CreateBranch", ctx, in)}
}

func (_c *MockInterface_CreateBranch_Call) Run(run func(ctx context.Context, in github.CreateBranchInput)) *MockInterface_CreateBranch_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(github.CreateBranchInput))
	})
	return _c
}

func (_c *MockInterface_CreateBranch_Call) Return(_a0 error) *MockInterface_CreateBranch_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockInterface_CreateBranch_Call) RunAndReturn(run func(context.Context, github.CreateBranchInput) error) *MockInterface_CreateBranch_Call {
	_c.Call.Return(run)
	return _c
}

// CreateCommit provides a mock function with given fields: ctx, commit
func (_m *MockInterface) CreateCommit(ctx context.Context, commit git.NewCommit) (git.CommitSHA, error) {
	ret := _m.Called(ctx, commit)

	if len(ret) == 0 {
		panic("no return value specified for CreateCommit")
	}

	var r0 git.CommitSHA
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, git.NewCommit) (git.CommitSHA, error)); ok {
		return rf(ctx, commit)
	}
	if rf, ok := ret.Get(0).(func(context.Context, git.NewCommit) git.CommitSHA); ok {
		r0 = rf(ctx, commit)
	} else {
		r0 = ret.Get(0).(git.CommitSHA)
	}

	if rf, ok := ret.Get(1).(func(context.Context, git.NewCommit) error); ok {
		r1 = rf(ctx, commit)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockInterface_CreateCommit_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateCommit'
type MockInterface_CreateCommit_Call struct {
	*mock.Call
}

// CreateCommit is a helper method to define mock.On call
//   - ctx context.Context
//   - commit git.NewCommit
func (_e *MockInterface_Expecter) CreateCommit(ctx interface{}, commit interface{}) *MockInterface_CreateCommit_Call {
	return &MockInterface_CreateCommit_Call{Call: _e.mock.On("CreateCommit", ctx, commit)}
}

func (_c *MockInterface_CreateCommit_Call) Run(run func(ctx context.Context, commit git.NewCommit)) *MockInterface_CreateCommit_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(git.NewCommit))
	})
	return _c
}

func (_c *MockInterface_CreateCommit_Call) Return(_a0 git.CommitSHA, _a1 error) *MockInterface_CreateCommit_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockInterface_CreateCommit_Call) RunAndReturn(run func(context.Context, git.NewCommit) (git.CommitSHA, error)) *MockInterface_CreateCommit_Call {
	_c.Call.Return(run)
	return _c
}

// CreateFork provides a mock function with given fields: ctx, id
func (_m *MockInterface) CreateFork(ctx context.Context, id git.RepositoryID) (*git.RepositoryID, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for CreateFork")
	}

	var r0 *git.RepositoryID
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, git.RepositoryID) (*git.RepositoryID, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, git.RepositoryID) *git.RepositoryID); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*git.RepositoryID)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, git.RepositoryID) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockInterface_CreateFork_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateFork'
type MockInterface_CreateFork_Call struct {
	*mock.Call
}

// CreateFork is a helper method to define mock.On call
//   - ctx context.Context
//   - id git.RepositoryID
func (_e *MockInterface_Expecter) CreateFork(ctx interface{}, id interface{}) *MockInterface_CreateFork_Call {
	return &MockInterface_CreateFork_Call{Call: _e.mock.On("CreateFork", ctx, id)}
}

func (_c *MockInterface_CreateFork_Call) Run(run func(ctx context.Context, id git.RepositoryID)) *MockInterface_CreateFork_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(git.RepositoryID))
	})
	return _c
}

func (_c *MockInterface_CreateFork_Call) Return(_a0 *git.RepositoryID, _a1 error) *MockInterface_CreateFork_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockInterface_CreateFork_Call) RunAndReturn(run func(context.Context, git.RepositoryID) (*git.RepositoryID, error)) *MockInterface_CreateFork_Call {
	_c.Call.Return(run)
	return _c
}

// CreatePullRequest provides a mock function with given fields: ctx, in
func (_m *MockInterface) CreatePullRequest(ctx context.Context, in github.CreatePullRequestInput) (*github.CreatePullRequestOutput, error) {
	ret := _m.Called(ctx, in)

	if len(ret) == 0 {
		panic("no return value specified for CreatePullRequest")
	}

	var r0 *github.CreatePullRequestOutput
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, github.CreatePullRequestInput) (*github.CreatePullRequestOutput, error)); ok {
		return rf(ctx, in)
	}
	if rf, ok := ret.Get(0).(func(context.Context, github.CreatePullRequestInput) *github.CreatePullRequestOutput); ok {
		r0 = rf(ctx, in)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*github.CreatePullRequestOutput)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, github.CreatePullRequestInput) error); ok {
		r1 = rf(ctx, in)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockInterface_CreatePullRequest_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreatePullRequest'
type MockInterface_CreatePullRequest_Call struct {
	*mock.Call
}

// CreatePullRequest is a helper method to define mock.On call
//   - ctx context.Context
//   - in github.CreatePullRequestInput
func (_e *MockInterface_Expecter) CreatePullRequest(ctx interface{}, in interface{}) *MockInterface_CreatePullRequest_Call {
	return &MockInterface_CreatePullRequest_Call{Call: _e.mock.On("CreatePullRequest", ctx, in)}
}

func (_c *MockInterface_CreatePullRequest_Call) Run(run func(ctx context.Context, in github.CreatePullRequestInput)) *MockInterface_CreatePullRequest_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(github.CreatePullRequestInput))
	})
	return _c
}

func (_c *MockInterface_CreatePullRequest_Call) Return(_a0 *github.CreatePullRequestOutput, _a1 error) *MockInterface_CreatePullRequest_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockInterface_CreatePullRequest_Call) RunAndReturn(run func(context.Context, github.CreatePullRequestInput) (*github.CreatePullRequestOutput, error)) *MockInterface_CreatePullRequest_Call {
	_c.Call.Return(run)
	return _c
}

// CreateRelease provides a mock function with given fields: ctx, r
func (_m *MockInterface) CreateRelease(ctx context.Context, r git.Release) (*git.Release, error) {
	ret := _m.Called(ctx, r)

	if len(ret) == 0 {
		panic("no return value specified for CreateRelease")
	}

	var r0 *git.Release
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, git.Release) (*git.Release, error)); ok {
		return rf(ctx, r)
	}
	if rf, ok := ret.Get(0).(func(context.Context, git.Release) *git.Release); ok {
		r0 = rf(ctx, r)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*git.Release)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, git.Release) error); ok {
		r1 = rf(ctx, r)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockInterface_CreateRelease_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateRelease'
type MockInterface_CreateRelease_Call struct {
	*mock.Call
}

// CreateRelease is a helper method to define mock.On call
//   - ctx context.Context
//   - r git.Release
func (_e *MockInterface_Expecter) CreateRelease(ctx interface{}, r interface{}) *MockInterface_CreateRelease_Call {
	return &MockInterface_CreateRelease_Call{Call: _e.mock.On("CreateRelease", ctx, r)}
}

func (_c *MockInterface_CreateRelease_Call) Run(run func(ctx context.Context, r git.Release)) *MockInterface_CreateRelease_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(git.Release))
	})
	return _c
}

func (_c *MockInterface_CreateRelease_Call) Return(_a0 *git.Release, _a1 error) *MockInterface_CreateRelease_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockInterface_CreateRelease_Call) RunAndReturn(run func(context.Context, git.Release) (*git.Release, error)) *MockInterface_CreateRelease_Call {
	_c.Call.Return(run)
	return _c
}

// CreateReleaseAsset provides a mock function with given fields: ctx, a
func (_m *MockInterface) CreateReleaseAsset(ctx context.Context, a git.ReleaseAsset) error {
	ret := _m.Called(ctx, a)

	if len(ret) == 0 {
		panic("no return value specified for CreateReleaseAsset")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, git.ReleaseAsset) error); ok {
		r0 = rf(ctx, a)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockInterface_CreateReleaseAsset_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateReleaseAsset'
type MockInterface_CreateReleaseAsset_Call struct {
	*mock.Call
}

// CreateReleaseAsset is a helper method to define mock.On call
//   - ctx context.Context
//   - a git.ReleaseAsset
func (_e *MockInterface_Expecter) CreateReleaseAsset(ctx interface{}, a interface{}) *MockInterface_CreateReleaseAsset_Call {
	return &MockInterface_CreateReleaseAsset_Call{Call: _e.mock.On("CreateReleaseAsset", ctx, a)}
}

func (_c *MockInterface_CreateReleaseAsset_Call) Run(run func(ctx context.Context, a git.ReleaseAsset)) *MockInterface_CreateReleaseAsset_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(git.ReleaseAsset))
	})
	return _c
}

func (_c *MockInterface_CreateReleaseAsset_Call) Return(_a0 error) *MockInterface_CreateReleaseAsset_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockInterface_CreateReleaseAsset_Call) RunAndReturn(run func(context.Context, git.ReleaseAsset) error) *MockInterface_CreateReleaseAsset_Call {
	_c.Call.Return(run)
	return _c
}

// CreateTree provides a mock function with given fields: ctx, tree
func (_m *MockInterface) CreateTree(ctx context.Context, tree git.NewTree) (git.TreeSHA, error) {
	ret := _m.Called(ctx, tree)

	if len(ret) == 0 {
		panic("no return value specified for CreateTree")
	}

	var r0 git.TreeSHA
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, git.NewTree) (git.TreeSHA, error)); ok {
		return rf(ctx, tree)
	}
	if rf, ok := ret.Get(0).(func(context.Context, git.NewTree) git.TreeSHA); ok {
		r0 = rf(ctx, tree)
	} else {
		r0 = ret.Get(0).(git.TreeSHA)
	}

	if rf, ok := ret.Get(1).(func(context.Context, git.NewTree) error); ok {
		r1 = rf(ctx, tree)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockInterface_CreateTree_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateTree'
type MockInterface_CreateTree_Call struct {
	*mock.Call
}

// CreateTree is a helper method to define mock.On call
//   - ctx context.Context
//   - tree git.NewTree
func (_e *MockInterface_Expecter) CreateTree(ctx interface{}, tree interface{}) *MockInterface_CreateTree_Call {
	return &MockInterface_CreateTree_Call{Call: _e.mock.On("CreateTree", ctx, tree)}
}

func (_c *MockInterface_CreateTree_Call) Run(run func(ctx context.Context, tree git.NewTree)) *MockInterface_CreateTree_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(git.NewTree))
	})
	return _c
}

func (_c *MockInterface_CreateTree_Call) Return(_a0 git.TreeSHA, _a1 error) *MockInterface_CreateTree_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockInterface_CreateTree_Call) RunAndReturn(run func(context.Context, git.NewTree) (git.TreeSHA, error)) *MockInterface_CreateTree_Call {
	_c.Call.Return(run)
	return _c
}

// GetReleaseByTagOrNil provides a mock function with given fields: ctx, repo, tag
func (_m *MockInterface) GetReleaseByTagOrNil(ctx context.Context, repo git.RepositoryID, tag git.TagName) (*git.Release, error) {
	ret := _m.Called(ctx, repo, tag)

	if len(ret) == 0 {
		panic("no return value specified for GetReleaseByTagOrNil")
	}

	var r0 *git.Release
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, git.RepositoryID, git.TagName) (*git.Release, error)); ok {
		return rf(ctx, repo, tag)
	}
	if rf, ok := ret.Get(0).(func(context.Context, git.RepositoryID, git.TagName) *git.Release); ok {
		r0 = rf(ctx, repo, tag)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*git.Release)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, git.RepositoryID, git.TagName) error); ok {
		r1 = rf(ctx, repo, tag)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockInterface_GetReleaseByTagOrNil_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetReleaseByTagOrNil'
type MockInterface_GetReleaseByTagOrNil_Call struct {
	*mock.Call
}

// GetReleaseByTagOrNil is a helper method to define mock.On call
//   - ctx context.Context
//   - repo git.RepositoryID
//   - tag git.TagName
func (_e *MockInterface_Expecter) GetReleaseByTagOrNil(ctx interface{}, repo interface{}, tag interface{}) *MockInterface_GetReleaseByTagOrNil_Call {
	return &MockInterface_GetReleaseByTagOrNil_Call{Call: _e.mock.On("GetReleaseByTagOrNil", ctx, repo, tag)}
}

func (_c *MockInterface_GetReleaseByTagOrNil_Call) Run(run func(ctx context.Context, repo git.RepositoryID, tag git.TagName)) *MockInterface_GetReleaseByTagOrNil_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(git.RepositoryID), args[2].(git.TagName))
	})
	return _c
}

func (_c *MockInterface_GetReleaseByTagOrNil_Call) Return(_a0 *git.Release, _a1 error) *MockInterface_GetReleaseByTagOrNil_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockInterface_GetReleaseByTagOrNil_Call) RunAndReturn(run func(context.Context, git.RepositoryID, git.TagName) (*git.Release, error)) *MockInterface_GetReleaseByTagOrNil_Call {
	_c.Call.Return(run)
	return _c
}

// QueryCommit provides a mock function with given fields: ctx, in
func (_m *MockInterface) QueryCommit(ctx context.Context, in github.QueryCommitInput) (*github.QueryCommitOutput, error) {
	ret := _m.Called(ctx, in)

	if len(ret) == 0 {
		panic("no return value specified for QueryCommit")
	}

	var r0 *github.QueryCommitOutput
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, github.QueryCommitInput) (*github.QueryCommitOutput, error)); ok {
		return rf(ctx, in)
	}
	if rf, ok := ret.Get(0).(func(context.Context, github.QueryCommitInput) *github.QueryCommitOutput); ok {
		r0 = rf(ctx, in)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*github.QueryCommitOutput)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, github.QueryCommitInput) error); ok {
		r1 = rf(ctx, in)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockInterface_QueryCommit_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'QueryCommit'
type MockInterface_QueryCommit_Call struct {
	*mock.Call
}

// QueryCommit is a helper method to define mock.On call
//   - ctx context.Context
//   - in github.QueryCommitInput
func (_e *MockInterface_Expecter) QueryCommit(ctx interface{}, in interface{}) *MockInterface_QueryCommit_Call {
	return &MockInterface_QueryCommit_Call{Call: _e.mock.On("QueryCommit", ctx, in)}
}

func (_c *MockInterface_QueryCommit_Call) Run(run func(ctx context.Context, in github.QueryCommitInput)) *MockInterface_QueryCommit_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(github.QueryCommitInput))
	})
	return _c
}

func (_c *MockInterface_QueryCommit_Call) Return(_a0 *github.QueryCommitOutput, _a1 error) *MockInterface_QueryCommit_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockInterface_QueryCommit_Call) RunAndReturn(run func(context.Context, github.QueryCommitInput) (*github.QueryCommitOutput, error)) *MockInterface_QueryCommit_Call {
	_c.Call.Return(run)
	return _c
}

// QueryDefaultBranch provides a mock function with given fields: ctx, in
func (_m *MockInterface) QueryDefaultBranch(ctx context.Context, in github.QueryDefaultBranchInput) (*github.QueryDefaultBranchOutput, error) {
	ret := _m.Called(ctx, in)

	if len(ret) == 0 {
		panic("no return value specified for QueryDefaultBranch")
	}

	var r0 *github.QueryDefaultBranchOutput
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, github.QueryDefaultBranchInput) (*github.QueryDefaultBranchOutput, error)); ok {
		return rf(ctx, in)
	}
	if rf, ok := ret.Get(0).(func(context.Context, github.QueryDefaultBranchInput) *github.QueryDefaultBranchOutput); ok {
		r0 = rf(ctx, in)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*github.QueryDefaultBranchOutput)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, github.QueryDefaultBranchInput) error); ok {
		r1 = rf(ctx, in)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockInterface_QueryDefaultBranch_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'QueryDefaultBranch'
type MockInterface_QueryDefaultBranch_Call struct {
	*mock.Call
}

// QueryDefaultBranch is a helper method to define mock.On call
//   - ctx context.Context
//   - in github.QueryDefaultBranchInput
func (_e *MockInterface_Expecter) QueryDefaultBranch(ctx interface{}, in interface{}) *MockInterface_QueryDefaultBranch_Call {
	return &MockInterface_QueryDefaultBranch_Call{Call: _e.mock.On("QueryDefaultBranch", ctx, in)}
}

func (_c *MockInterface_QueryDefaultBranch_Call) Run(run func(ctx context.Context, in github.QueryDefaultBranchInput)) *MockInterface_QueryDefaultBranch_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(github.QueryDefaultBranchInput))
	})
	return _c
}

func (_c *MockInterface_QueryDefaultBranch_Call) Return(_a0 *github.QueryDefaultBranchOutput, _a1 error) *MockInterface_QueryDefaultBranch_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockInterface_QueryDefaultBranch_Call) RunAndReturn(run func(context.Context, github.QueryDefaultBranchInput) (*github.QueryDefaultBranchOutput, error)) *MockInterface_QueryDefaultBranch_Call {
	_c.Call.Return(run)
	return _c
}

// QueryForCommit provides a mock function with given fields: ctx, in
func (_m *MockInterface) QueryForCommit(ctx context.Context, in github.QueryForCommitInput) (*github.QueryForCommitOutput, error) {
	ret := _m.Called(ctx, in)

	if len(ret) == 0 {
		panic("no return value specified for QueryForCommit")
	}

	var r0 *github.QueryForCommitOutput
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, github.QueryForCommitInput) (*github.QueryForCommitOutput, error)); ok {
		return rf(ctx, in)
	}
	if rf, ok := ret.Get(0).(func(context.Context, github.QueryForCommitInput) *github.QueryForCommitOutput); ok {
		r0 = rf(ctx, in)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*github.QueryForCommitOutput)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, github.QueryForCommitInput) error); ok {
		r1 = rf(ctx, in)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockInterface_QueryForCommit_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'QueryForCommit'
type MockInterface_QueryForCommit_Call struct {
	*mock.Call
}

// QueryForCommit is a helper method to define mock.On call
//   - ctx context.Context
//   - in github.QueryForCommitInput
func (_e *MockInterface_Expecter) QueryForCommit(ctx interface{}, in interface{}) *MockInterface_QueryForCommit_Call {
	return &MockInterface_QueryForCommit_Call{Call: _e.mock.On("QueryForCommit", ctx, in)}
}

func (_c *MockInterface_QueryForCommit_Call) Run(run func(ctx context.Context, in github.QueryForCommitInput)) *MockInterface_QueryForCommit_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(github.QueryForCommitInput))
	})
	return _c
}

func (_c *MockInterface_QueryForCommit_Call) Return(_a0 *github.QueryForCommitOutput, _a1 error) *MockInterface_QueryForCommit_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockInterface_QueryForCommit_Call) RunAndReturn(run func(context.Context, github.QueryForCommitInput) (*github.QueryForCommitOutput, error)) *MockInterface_QueryForCommit_Call {
	_c.Call.Return(run)
	return _c
}

// QueryForPullRequest provides a mock function with given fields: ctx, in
func (_m *MockInterface) QueryForPullRequest(ctx context.Context, in github.QueryForPullRequestInput) (*github.QueryForPullRequestOutput, error) {
	ret := _m.Called(ctx, in)

	if len(ret) == 0 {
		panic("no return value specified for QueryForPullRequest")
	}

	var r0 *github.QueryForPullRequestOutput
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, github.QueryForPullRequestInput) (*github.QueryForPullRequestOutput, error)); ok {
		return rf(ctx, in)
	}
	if rf, ok := ret.Get(0).(func(context.Context, github.QueryForPullRequestInput) *github.QueryForPullRequestOutput); ok {
		r0 = rf(ctx, in)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*github.QueryForPullRequestOutput)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, github.QueryForPullRequestInput) error); ok {
		r1 = rf(ctx, in)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockInterface_QueryForPullRequest_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'QueryForPullRequest'
type MockInterface_QueryForPullRequest_Call struct {
	*mock.Call
}

// QueryForPullRequest is a helper method to define mock.On call
//   - ctx context.Context
//   - in github.QueryForPullRequestInput
func (_e *MockInterface_Expecter) QueryForPullRequest(ctx interface{}, in interface{}) *MockInterface_QueryForPullRequest_Call {
	return &MockInterface_QueryForPullRequest_Call{Call: _e.mock.On("QueryForPullRequest", ctx, in)}
}

func (_c *MockInterface_QueryForPullRequest_Call) Run(run func(ctx context.Context, in github.QueryForPullRequestInput)) *MockInterface_QueryForPullRequest_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(github.QueryForPullRequestInput))
	})
	return _c
}

func (_c *MockInterface_QueryForPullRequest_Call) Return(_a0 *github.QueryForPullRequestOutput, _a1 error) *MockInterface_QueryForPullRequest_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockInterface_QueryForPullRequest_Call) RunAndReturn(run func(context.Context, github.QueryForPullRequestInput) (*github.QueryForPullRequestOutput, error)) *MockInterface_QueryForPullRequest_Call {
	_c.Call.Return(run)
	return _c
}

// RequestPullRequestReview provides a mock function with given fields: ctx, in
func (_m *MockInterface) RequestPullRequestReview(ctx context.Context, in github.RequestPullRequestReviewInput) error {
	ret := _m.Called(ctx, in)

	if len(ret) == 0 {
		panic("no return value specified for RequestPullRequestReview")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, github.RequestPullRequestReviewInput) error); ok {
		r0 = rf(ctx, in)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockInterface_RequestPullRequestReview_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RequestPullRequestReview'
type MockInterface_RequestPullRequestReview_Call struct {
	*mock.Call
}

// RequestPullRequestReview is a helper method to define mock.On call
//   - ctx context.Context
//   - in github.RequestPullRequestReviewInput
func (_e *MockInterface_Expecter) RequestPullRequestReview(ctx interface{}, in interface{}) *MockInterface_RequestPullRequestReview_Call {
	return &MockInterface_RequestPullRequestReview_Call{Call: _e.mock.On("RequestPullRequestReview", ctx, in)}
}

func (_c *MockInterface_RequestPullRequestReview_Call) Run(run func(ctx context.Context, in github.RequestPullRequestReviewInput)) *MockInterface_RequestPullRequestReview_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(github.RequestPullRequestReviewInput))
	})
	return _c
}

func (_c *MockInterface_RequestPullRequestReview_Call) Return(_a0 error) *MockInterface_RequestPullRequestReview_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockInterface_RequestPullRequestReview_Call) RunAndReturn(run func(context.Context, github.RequestPullRequestReviewInput) error) *MockInterface_RequestPullRequestReview_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateBranch provides a mock function with given fields: ctx, in
func (_m *MockInterface) UpdateBranch(ctx context.Context, in github.UpdateBranchInput) error {
	ret := _m.Called(ctx, in)

	if len(ret) == 0 {
		panic("no return value specified for UpdateBranch")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, github.UpdateBranchInput) error); ok {
		r0 = rf(ctx, in)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockInterface_UpdateBranch_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateBranch'
type MockInterface_UpdateBranch_Call struct {
	*mock.Call
}

// UpdateBranch is a helper method to define mock.On call
//   - ctx context.Context
//   - in github.UpdateBranchInput
func (_e *MockInterface_Expecter) UpdateBranch(ctx interface{}, in interface{}) *MockInterface_UpdateBranch_Call {
	return &MockInterface_UpdateBranch_Call{Call: _e.mock.On("UpdateBranch", ctx, in)}
}

func (_c *MockInterface_UpdateBranch_Call) Run(run func(ctx context.Context, in github.UpdateBranchInput)) *MockInterface_UpdateBranch_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(github.UpdateBranchInput))
	})
	return _c
}

func (_c *MockInterface_UpdateBranch_Call) Return(_a0 error) *MockInterface_UpdateBranch_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockInterface_UpdateBranch_Call) RunAndReturn(run func(context.Context, github.UpdateBranchInput) error) *MockInterface_UpdateBranch_Call {
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
