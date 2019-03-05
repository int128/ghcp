// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/int128/ghcp/adaptors/interfaces (interfaces: Cmd,FileSystem,GitHub)

// Package mock_adaptors is a generated GoMock package.
package mock_adaptors

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	interfaces "github.com/int128/ghcp/adaptors/interfaces"
	git "github.com/int128/ghcp/git"
	reflect "reflect"
)

// MockCmd is a mock of Cmd interface
type MockCmd struct {
	ctrl     *gomock.Controller
	recorder *MockCmdMockRecorder
}

// MockCmdMockRecorder is the mock recorder for MockCmd
type MockCmdMockRecorder struct {
	mock *MockCmd
}

// NewMockCmd creates a new mock instance
func NewMockCmd(ctrl *gomock.Controller) *MockCmd {
	mock := &MockCmd{ctrl: ctrl}
	mock.recorder = &MockCmdMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCmd) EXPECT() *MockCmdMockRecorder {
	return m.recorder
}

// Run mocks base method
func (m *MockCmd) Run(arg0 context.Context, arg1 interfaces.CmdOptions) error {
	ret := m.ctrl.Call(m, "Run", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Run indicates an expected call of Run
func (mr *MockCmdMockRecorder) Run(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Run", reflect.TypeOf((*MockCmd)(nil).Run), arg0, arg1)
}

// MockFileSystem is a mock of FileSystem interface
type MockFileSystem struct {
	ctrl     *gomock.Controller
	recorder *MockFileSystemMockRecorder
}

// MockFileSystemMockRecorder is the mock recorder for MockFileSystem
type MockFileSystemMockRecorder struct {
	mock *MockFileSystem
}

// NewMockFileSystem creates a new mock instance
func NewMockFileSystem(ctrl *gomock.Controller) *MockFileSystem {
	mock := &MockFileSystem{ctrl: ctrl}
	mock.recorder = &MockFileSystemMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockFileSystem) EXPECT() *MockFileSystemMockRecorder {
	return m.recorder
}

// FindFiles mocks base method
func (m *MockFileSystem) FindFiles(arg0 []string) ([]string, error) {
	ret := m.ctrl.Call(m, "FindFiles", arg0)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindFiles indicates an expected call of FindFiles
func (mr *MockFileSystemMockRecorder) FindFiles(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindFiles", reflect.TypeOf((*MockFileSystem)(nil).FindFiles), arg0)
}

// ReadAsBase64EncodedContent mocks base method
func (m *MockFileSystem) ReadAsBase64EncodedContent(arg0 string) (string, error) {
	ret := m.ctrl.Call(m, "ReadAsBase64EncodedContent", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadAsBase64EncodedContent indicates an expected call of ReadAsBase64EncodedContent
func (mr *MockFileSystemMockRecorder) ReadAsBase64EncodedContent(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadAsBase64EncodedContent", reflect.TypeOf((*MockFileSystem)(nil).ReadAsBase64EncodedContent), arg0)
}

// MockGitHub is a mock of GitHub interface
type MockGitHub struct {
	ctrl     *gomock.Controller
	recorder *MockGitHubMockRecorder
}

// MockGitHubMockRecorder is the mock recorder for MockGitHub
type MockGitHubMockRecorder struct {
	mock *MockGitHub
}

// NewMockGitHub creates a new mock instance
func NewMockGitHub(ctrl *gomock.Controller) *MockGitHub {
	mock := &MockGitHub{ctrl: ctrl}
	mock.recorder = &MockGitHubMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockGitHub) EXPECT() *MockGitHubMockRecorder {
	return m.recorder
}

// CreateBlob mocks base method
func (m *MockGitHub) CreateBlob(arg0 context.Context, arg1 git.NewBlob) (git.BlobSHA, error) {
	ret := m.ctrl.Call(m, "CreateBlob", arg0, arg1)
	ret0, _ := ret[0].(git.BlobSHA)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateBlob indicates an expected call of CreateBlob
func (mr *MockGitHubMockRecorder) CreateBlob(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateBlob", reflect.TypeOf((*MockGitHub)(nil).CreateBlob), arg0, arg1)
}

// CreateBranch mocks base method
func (m *MockGitHub) CreateBranch(arg0 context.Context, arg1 git.NewBranch) error {
	ret := m.ctrl.Call(m, "CreateBranch", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateBranch indicates an expected call of CreateBranch
func (mr *MockGitHubMockRecorder) CreateBranch(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateBranch", reflect.TypeOf((*MockGitHub)(nil).CreateBranch), arg0, arg1)
}

// CreateCommit mocks base method
func (m *MockGitHub) CreateCommit(arg0 context.Context, arg1 git.NewCommit) (git.CommitSHA, error) {
	ret := m.ctrl.Call(m, "CreateCommit", arg0, arg1)
	ret0, _ := ret[0].(git.CommitSHA)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateCommit indicates an expected call of CreateCommit
func (mr *MockGitHubMockRecorder) CreateCommit(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCommit", reflect.TypeOf((*MockGitHub)(nil).CreateCommit), arg0, arg1)
}

// CreateTree mocks base method
func (m *MockGitHub) CreateTree(arg0 context.Context, arg1 git.NewTree) (git.TreeSHA, error) {
	ret := m.ctrl.Call(m, "CreateTree", arg0, arg1)
	ret0, _ := ret[0].(git.TreeSHA)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateTree indicates an expected call of CreateTree
func (mr *MockGitHubMockRecorder) CreateTree(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTree", reflect.TypeOf((*MockGitHub)(nil).CreateTree), arg0, arg1)
}

// QueryCommit mocks base method
func (m *MockGitHub) QueryCommit(arg0 context.Context, arg1 interfaces.QueryCommitIn) (*interfaces.QueryCommitOut, error) {
	ret := m.ctrl.Call(m, "QueryCommit", arg0, arg1)
	ret0, _ := ret[0].(*interfaces.QueryCommitOut)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryCommit indicates an expected call of QueryCommit
func (mr *MockGitHubMockRecorder) QueryCommit(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryCommit", reflect.TypeOf((*MockGitHub)(nil).QueryCommit), arg0, arg1)
}

// QueryRepository mocks base method
func (m *MockGitHub) QueryRepository(arg0 context.Context, arg1 interfaces.QueryRepositoryIn) (*interfaces.QueryRepositoryOut, error) {
	ret := m.ctrl.Call(m, "QueryRepository", arg0, arg1)
	ret0, _ := ret[0].(*interfaces.QueryRepositoryOut)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryRepository indicates an expected call of QueryRepository
func (mr *MockGitHubMockRecorder) QueryRepository(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryRepository", reflect.TypeOf((*MockGitHub)(nil).QueryRepository), arg0, arg1)
}

// UpdateBranch mocks base method
func (m *MockGitHub) UpdateBranch(arg0 context.Context, arg1 git.NewBranch, arg2 bool) error {
	ret := m.ctrl.Call(m, "UpdateBranch", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateBranch indicates an expected call of UpdateBranch
func (mr *MockGitHubMockRecorder) UpdateBranch(arg0, arg1, arg2 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateBranch", reflect.TypeOf((*MockGitHub)(nil).UpdateBranch), arg0, arg1, arg2)
}