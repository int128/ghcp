// Code generated by mockery v2.48.0. DO NOT EDIT.

package client_mock

import (
	context "context"

	github "github.com/google/go-github/v66/github"
	githubv4 "github.com/shurcooL/githubv4"

	mock "github.com/stretchr/testify/mock"

	os "os"
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

// CreateBlob provides a mock function with given fields: ctx, owner, repo, blob
func (_m *MockInterface) CreateBlob(ctx context.Context, owner string, repo string, blob *github.Blob) (*github.Blob, *github.Response, error) {
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

// MockInterface_CreateBlob_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateBlob'
type MockInterface_CreateBlob_Call struct {
	*mock.Call
}

// CreateBlob is a helper method to define mock.On call
//   - ctx context.Context
//   - owner string
//   - repo string
//   - blob *github.Blob
func (_e *MockInterface_Expecter) CreateBlob(ctx interface{}, owner interface{}, repo interface{}, blob interface{}) *MockInterface_CreateBlob_Call {
	return &MockInterface_CreateBlob_Call{Call: _e.mock.On("CreateBlob", ctx, owner, repo, blob)}
}

func (_c *MockInterface_CreateBlob_Call) Run(run func(ctx context.Context, owner string, repo string, blob *github.Blob)) *MockInterface_CreateBlob_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string), args[3].(*github.Blob))
	})
	return _c
}

func (_c *MockInterface_CreateBlob_Call) Return(_a0 *github.Blob, _a1 *github.Response, _a2 error) *MockInterface_CreateBlob_Call {
	_c.Call.Return(_a0, _a1, _a2)
	return _c
}

func (_c *MockInterface_CreateBlob_Call) RunAndReturn(run func(context.Context, string, string, *github.Blob) (*github.Blob, *github.Response, error)) *MockInterface_CreateBlob_Call {
	_c.Call.Return(run)
	return _c
}

// CreateCommit provides a mock function with given fields: ctx, owner, repo, commit, opts
func (_m *MockInterface) CreateCommit(ctx context.Context, owner string, repo string, commit *github.Commit, opts *github.CreateCommitOptions) (*github.Commit, *github.Response, error) {
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

// MockInterface_CreateCommit_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateCommit'
type MockInterface_CreateCommit_Call struct {
	*mock.Call
}

// CreateCommit is a helper method to define mock.On call
//   - ctx context.Context
//   - owner string
//   - repo string
//   - commit *github.Commit
//   - opts *github.CreateCommitOptions
func (_e *MockInterface_Expecter) CreateCommit(ctx interface{}, owner interface{}, repo interface{}, commit interface{}, opts interface{}) *MockInterface_CreateCommit_Call {
	return &MockInterface_CreateCommit_Call{Call: _e.mock.On("CreateCommit", ctx, owner, repo, commit, opts)}
}

func (_c *MockInterface_CreateCommit_Call) Run(run func(ctx context.Context, owner string, repo string, commit *github.Commit, opts *github.CreateCommitOptions)) *MockInterface_CreateCommit_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string), args[3].(*github.Commit), args[4].(*github.CreateCommitOptions))
	})
	return _c
}

func (_c *MockInterface_CreateCommit_Call) Return(_a0 *github.Commit, _a1 *github.Response, _a2 error) *MockInterface_CreateCommit_Call {
	_c.Call.Return(_a0, _a1, _a2)
	return _c
}

func (_c *MockInterface_CreateCommit_Call) RunAndReturn(run func(context.Context, string, string, *github.Commit, *github.CreateCommitOptions) (*github.Commit, *github.Response, error)) *MockInterface_CreateCommit_Call {
	_c.Call.Return(run)
	return _c
}

// CreateFork provides a mock function with given fields: ctx, owner, repo, opt
func (_m *MockInterface) CreateFork(ctx context.Context, owner string, repo string, opt *github.RepositoryCreateForkOptions) (*github.Repository, *github.Response, error) {
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

// MockInterface_CreateFork_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateFork'
type MockInterface_CreateFork_Call struct {
	*mock.Call
}

// CreateFork is a helper method to define mock.On call
//   - ctx context.Context
//   - owner string
//   - repo string
//   - opt *github.RepositoryCreateForkOptions
func (_e *MockInterface_Expecter) CreateFork(ctx interface{}, owner interface{}, repo interface{}, opt interface{}) *MockInterface_CreateFork_Call {
	return &MockInterface_CreateFork_Call{Call: _e.mock.On("CreateFork", ctx, owner, repo, opt)}
}

func (_c *MockInterface_CreateFork_Call) Run(run func(ctx context.Context, owner string, repo string, opt *github.RepositoryCreateForkOptions)) *MockInterface_CreateFork_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string), args[3].(*github.RepositoryCreateForkOptions))
	})
	return _c
}

func (_c *MockInterface_CreateFork_Call) Return(_a0 *github.Repository, _a1 *github.Response, _a2 error) *MockInterface_CreateFork_Call {
	_c.Call.Return(_a0, _a1, _a2)
	return _c
}

func (_c *MockInterface_CreateFork_Call) RunAndReturn(run func(context.Context, string, string, *github.RepositoryCreateForkOptions) (*github.Repository, *github.Response, error)) *MockInterface_CreateFork_Call {
	_c.Call.Return(run)
	return _c
}

// CreateRelease provides a mock function with given fields: ctx, owner, repo, release
func (_m *MockInterface) CreateRelease(ctx context.Context, owner string, repo string, release *github.RepositoryRelease) (*github.RepositoryRelease, *github.Response, error) {
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

// MockInterface_CreateRelease_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateRelease'
type MockInterface_CreateRelease_Call struct {
	*mock.Call
}

// CreateRelease is a helper method to define mock.On call
//   - ctx context.Context
//   - owner string
//   - repo string
//   - release *github.RepositoryRelease
func (_e *MockInterface_Expecter) CreateRelease(ctx interface{}, owner interface{}, repo interface{}, release interface{}) *MockInterface_CreateRelease_Call {
	return &MockInterface_CreateRelease_Call{Call: _e.mock.On("CreateRelease", ctx, owner, repo, release)}
}

func (_c *MockInterface_CreateRelease_Call) Run(run func(ctx context.Context, owner string, repo string, release *github.RepositoryRelease)) *MockInterface_CreateRelease_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string), args[3].(*github.RepositoryRelease))
	})
	return _c
}

func (_c *MockInterface_CreateRelease_Call) Return(_a0 *github.RepositoryRelease, _a1 *github.Response, _a2 error) *MockInterface_CreateRelease_Call {
	_c.Call.Return(_a0, _a1, _a2)
	return _c
}

func (_c *MockInterface_CreateRelease_Call) RunAndReturn(run func(context.Context, string, string, *github.RepositoryRelease) (*github.RepositoryRelease, *github.Response, error)) *MockInterface_CreateRelease_Call {
	_c.Call.Return(run)
	return _c
}

// CreateTree provides a mock function with given fields: ctx, owner, repo, baseTree, entries
func (_m *MockInterface) CreateTree(ctx context.Context, owner string, repo string, baseTree string, entries []*github.TreeEntry) (*github.Tree, *github.Response, error) {
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

// MockInterface_CreateTree_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateTree'
type MockInterface_CreateTree_Call struct {
	*mock.Call
}

// CreateTree is a helper method to define mock.On call
//   - ctx context.Context
//   - owner string
//   - repo string
//   - baseTree string
//   - entries []*github.TreeEntry
func (_e *MockInterface_Expecter) CreateTree(ctx interface{}, owner interface{}, repo interface{}, baseTree interface{}, entries interface{}) *MockInterface_CreateTree_Call {
	return &MockInterface_CreateTree_Call{Call: _e.mock.On("CreateTree", ctx, owner, repo, baseTree, entries)}
}

func (_c *MockInterface_CreateTree_Call) Run(run func(ctx context.Context, owner string, repo string, baseTree string, entries []*github.TreeEntry)) *MockInterface_CreateTree_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string), args[3].(string), args[4].([]*github.TreeEntry))
	})
	return _c
}

func (_c *MockInterface_CreateTree_Call) Return(_a0 *github.Tree, _a1 *github.Response, _a2 error) *MockInterface_CreateTree_Call {
	_c.Call.Return(_a0, _a1, _a2)
	return _c
}

func (_c *MockInterface_CreateTree_Call) RunAndReturn(run func(context.Context, string, string, string, []*github.TreeEntry) (*github.Tree, *github.Response, error)) *MockInterface_CreateTree_Call {
	_c.Call.Return(run)
	return _c
}

// GetReleaseByTag provides a mock function with given fields: ctx, owner, repo, tag
func (_m *MockInterface) GetReleaseByTag(ctx context.Context, owner string, repo string, tag string) (*github.RepositoryRelease, *github.Response, error) {
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

// MockInterface_GetReleaseByTag_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetReleaseByTag'
type MockInterface_GetReleaseByTag_Call struct {
	*mock.Call
}

// GetReleaseByTag is a helper method to define mock.On call
//   - ctx context.Context
//   - owner string
//   - repo string
//   - tag string
func (_e *MockInterface_Expecter) GetReleaseByTag(ctx interface{}, owner interface{}, repo interface{}, tag interface{}) *MockInterface_GetReleaseByTag_Call {
	return &MockInterface_GetReleaseByTag_Call{Call: _e.mock.On("GetReleaseByTag", ctx, owner, repo, tag)}
}

func (_c *MockInterface_GetReleaseByTag_Call) Run(run func(ctx context.Context, owner string, repo string, tag string)) *MockInterface_GetReleaseByTag_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string), args[3].(string))
	})
	return _c
}

func (_c *MockInterface_GetReleaseByTag_Call) Return(_a0 *github.RepositoryRelease, _a1 *github.Response, _a2 error) *MockInterface_GetReleaseByTag_Call {
	_c.Call.Return(_a0, _a1, _a2)
	return _c
}

func (_c *MockInterface_GetReleaseByTag_Call) RunAndReturn(run func(context.Context, string, string, string) (*github.RepositoryRelease, *github.Response, error)) *MockInterface_GetReleaseByTag_Call {
	_c.Call.Return(run)
	return _c
}

// Mutate provides a mock function with given fields: ctx, m, input, variables
func (_m *MockInterface) Mutate(ctx context.Context, m interface{}, input githubv4.Input, variables map[string]interface{}) error {
	ret := _m.Called(ctx, m, input, variables)

	if len(ret) == 0 {
		panic("no return value specified for Mutate")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, interface{}, githubv4.Input, map[string]interface{}) error); ok {
		r0 = rf(ctx, m, input, variables)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockInterface_Mutate_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Mutate'
type MockInterface_Mutate_Call struct {
	*mock.Call
}

// Mutate is a helper method to define mock.On call
//   - ctx context.Context
//   - m interface{}
//   - input githubv4.Input
//   - variables map[string]interface{}
func (_e *MockInterface_Expecter) Mutate(ctx interface{}, m interface{}, input interface{}, variables interface{}) *MockInterface_Mutate_Call {
	return &MockInterface_Mutate_Call{Call: _e.mock.On("Mutate", ctx, m, input, variables)}
}

func (_c *MockInterface_Mutate_Call) Run(run func(ctx context.Context, m interface{}, input githubv4.Input, variables map[string]interface{})) *MockInterface_Mutate_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(interface{}), args[2].(githubv4.Input), args[3].(map[string]interface{}))
	})
	return _c
}

func (_c *MockInterface_Mutate_Call) Return(_a0 error) *MockInterface_Mutate_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockInterface_Mutate_Call) RunAndReturn(run func(context.Context, interface{}, githubv4.Input, map[string]interface{}) error) *MockInterface_Mutate_Call {
	_c.Call.Return(run)
	return _c
}

// Query provides a mock function with given fields: ctx, q, variables
func (_m *MockInterface) Query(ctx context.Context, q interface{}, variables map[string]interface{}) error {
	ret := _m.Called(ctx, q, variables)

	if len(ret) == 0 {
		panic("no return value specified for Query")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, interface{}, map[string]interface{}) error); ok {
		r0 = rf(ctx, q, variables)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockInterface_Query_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Query'
type MockInterface_Query_Call struct {
	*mock.Call
}

// Query is a helper method to define mock.On call
//   - ctx context.Context
//   - q interface{}
//   - variables map[string]interface{}
func (_e *MockInterface_Expecter) Query(ctx interface{}, q interface{}, variables interface{}) *MockInterface_Query_Call {
	return &MockInterface_Query_Call{Call: _e.mock.On("Query", ctx, q, variables)}
}

func (_c *MockInterface_Query_Call) Run(run func(ctx context.Context, q interface{}, variables map[string]interface{})) *MockInterface_Query_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(interface{}), args[2].(map[string]interface{}))
	})
	return _c
}

func (_c *MockInterface_Query_Call) Return(_a0 error) *MockInterface_Query_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockInterface_Query_Call) RunAndReturn(run func(context.Context, interface{}, map[string]interface{}) error) *MockInterface_Query_Call {
	_c.Call.Return(run)
	return _c
}

// UploadReleaseAsset provides a mock function with given fields: ctx, owner, repo, id, opt, file
func (_m *MockInterface) UploadReleaseAsset(ctx context.Context, owner string, repo string, id int64, opt *github.UploadOptions, file *os.File) (*github.ReleaseAsset, *github.Response, error) {
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

// MockInterface_UploadReleaseAsset_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UploadReleaseAsset'
type MockInterface_UploadReleaseAsset_Call struct {
	*mock.Call
}

// UploadReleaseAsset is a helper method to define mock.On call
//   - ctx context.Context
//   - owner string
//   - repo string
//   - id int64
//   - opt *github.UploadOptions
//   - file *os.File
func (_e *MockInterface_Expecter) UploadReleaseAsset(ctx interface{}, owner interface{}, repo interface{}, id interface{}, opt interface{}, file interface{}) *MockInterface_UploadReleaseAsset_Call {
	return &MockInterface_UploadReleaseAsset_Call{Call: _e.mock.On("UploadReleaseAsset", ctx, owner, repo, id, opt, file)}
}

func (_c *MockInterface_UploadReleaseAsset_Call) Run(run func(ctx context.Context, owner string, repo string, id int64, opt *github.UploadOptions, file *os.File)) *MockInterface_UploadReleaseAsset_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string), args[3].(int64), args[4].(*github.UploadOptions), args[5].(*os.File))
	})
	return _c
}

func (_c *MockInterface_UploadReleaseAsset_Call) Return(_a0 *github.ReleaseAsset, _a1 *github.Response, _a2 error) *MockInterface_UploadReleaseAsset_Call {
	_c.Call.Return(_a0, _a1, _a2)
	return _c
}

func (_c *MockInterface_UploadReleaseAsset_Call) RunAndReturn(run func(context.Context, string, string, int64, *github.UploadOptions, *os.File) (*github.ReleaseAsset, *github.Response, error)) *MockInterface_UploadReleaseAsset_Call {
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
