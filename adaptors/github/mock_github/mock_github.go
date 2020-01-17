// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/int128/ghcp/adaptors/github (interfaces: Interface)

// Package mock_github is a generated GoMock package.
package mock_github

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	github "github.com/int128/ghcp/adaptors/github"
	git "github.com/int128/ghcp/domain/git"
	reflect "reflect"
)

// MockInterface is a mock of Interface interface
type MockInterface struct {
	ctrl     *gomock.Controller
	recorder *MockInterfaceMockRecorder
}

// MockInterfaceMockRecorder is the mock recorder for MockInterface
type MockInterfaceMockRecorder struct {
	mock *MockInterface
}

// NewMockInterface creates a new mock instance
func NewMockInterface(ctrl *gomock.Controller) *MockInterface {
	mock := &MockInterface{ctrl: ctrl}
	mock.recorder = &MockInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockInterface) EXPECT() *MockInterfaceMockRecorder {
	return m.recorder
}

// CreateBlob mocks base method
func (m *MockInterface) CreateBlob(arg0 context.Context, arg1 git.NewBlob) (git.BlobSHA, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateBlob", arg0, arg1)
	ret0, _ := ret[0].(git.BlobSHA)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateBlob indicates an expected call of CreateBlob
func (mr *MockInterfaceMockRecorder) CreateBlob(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateBlob", reflect.TypeOf((*MockInterface)(nil).CreateBlob), arg0, arg1)
}

// CreateBranch mocks base method
func (m *MockInterface) CreateBranch(arg0 context.Context, arg1 git.NewBranch) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateBranch", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateBranch indicates an expected call of CreateBranch
func (mr *MockInterfaceMockRecorder) CreateBranch(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateBranch", reflect.TypeOf((*MockInterface)(nil).CreateBranch), arg0, arg1)
}

// CreateCommit mocks base method
func (m *MockInterface) CreateCommit(arg0 context.Context, arg1 git.NewCommit) (git.CommitSHA, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateCommit", arg0, arg1)
	ret0, _ := ret[0].(git.CommitSHA)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateCommit indicates an expected call of CreateCommit
func (mr *MockInterfaceMockRecorder) CreateCommit(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCommit", reflect.TypeOf((*MockInterface)(nil).CreateCommit), arg0, arg1)
}

// CreateFork mocks base method
func (m *MockInterface) CreateFork(arg0 context.Context, arg1 git.RepositoryID) (*git.RepositoryID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateFork", arg0, arg1)
	ret0, _ := ret[0].(*git.RepositoryID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateFork indicates an expected call of CreateFork
func (mr *MockInterfaceMockRecorder) CreateFork(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateFork", reflect.TypeOf((*MockInterface)(nil).CreateFork), arg0, arg1)
}

// CreateRelease mocks base method
func (m *MockInterface) CreateRelease(arg0 context.Context, arg1 git.Release) (*git.Release, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateRelease", arg0, arg1)
	ret0, _ := ret[0].(*git.Release)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateRelease indicates an expected call of CreateRelease
func (mr *MockInterfaceMockRecorder) CreateRelease(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateRelease", reflect.TypeOf((*MockInterface)(nil).CreateRelease), arg0, arg1)
}

// CreateReleaseAsset mocks base method
func (m *MockInterface) CreateReleaseAsset(arg0 context.Context, arg1 git.ReleaseAsset) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateReleaseAsset", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateReleaseAsset indicates an expected call of CreateReleaseAsset
func (mr *MockInterfaceMockRecorder) CreateReleaseAsset(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateReleaseAsset", reflect.TypeOf((*MockInterface)(nil).CreateReleaseAsset), arg0, arg1)
}

// CreateTree mocks base method
func (m *MockInterface) CreateTree(arg0 context.Context, arg1 git.NewTree) (git.TreeSHA, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTree", arg0, arg1)
	ret0, _ := ret[0].(git.TreeSHA)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateTree indicates an expected call of CreateTree
func (mr *MockInterfaceMockRecorder) CreateTree(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTree", reflect.TypeOf((*MockInterface)(nil).CreateTree), arg0, arg1)
}

// GetReleaseByTagOrNil mocks base method
func (m *MockInterface) GetReleaseByTagOrNil(arg0 context.Context, arg1 git.RepositoryID, arg2 git.TagName) (*git.Release, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetReleaseByTagOrNil", arg0, arg1, arg2)
	ret0, _ := ret[0].(*git.Release)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetReleaseByTagOrNil indicates an expected call of GetReleaseByTagOrNil
func (mr *MockInterfaceMockRecorder) GetReleaseByTagOrNil(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetReleaseByTagOrNil", reflect.TypeOf((*MockInterface)(nil).GetReleaseByTagOrNil), arg0, arg1, arg2)
}

// QueryCommit mocks base method
func (m *MockInterface) QueryCommit(arg0 context.Context, arg1 github.QueryCommitIn) (*github.QueryCommitOut, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryCommit", arg0, arg1)
	ret0, _ := ret[0].(*github.QueryCommitOut)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryCommit indicates an expected call of QueryCommit
func (mr *MockInterfaceMockRecorder) QueryCommit(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryCommit", reflect.TypeOf((*MockInterface)(nil).QueryCommit), arg0, arg1)
}

// QueryForCommitToBranch mocks base method
func (m *MockInterface) QueryForCommitToBranch(arg0 context.Context, arg1 github.QueryForCommitToBranchIn) (*github.QueryForCommitToBranchOut, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryForCommitToBranch", arg0, arg1)
	ret0, _ := ret[0].(*github.QueryForCommitToBranchOut)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryForCommitToBranch indicates an expected call of QueryForCommitToBranch
func (mr *MockInterfaceMockRecorder) QueryForCommitToBranch(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryForCommitToBranch", reflect.TypeOf((*MockInterface)(nil).QueryForCommitToBranch), arg0, arg1)
}

// UpdateBranch mocks base method
func (m *MockInterface) UpdateBranch(arg0 context.Context, arg1 git.NewBranch, arg2 bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateBranch", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateBranch indicates an expected call of UpdateBranch
func (mr *MockInterfaceMockRecorder) UpdateBranch(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateBranch", reflect.TypeOf((*MockInterface)(nil).UpdateBranch), arg0, arg1, arg2)
}