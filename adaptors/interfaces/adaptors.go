package adaptors

import (
	"context"

	"github.com/int128/ghcp/git"
)

//go:generate mockgen -package mock_adaptors -destination ../mock_adaptors/mock_adaptors.go github.com/int128/ghcp/adaptors/interfaces Cmd,FileSystem,GitHub

type Cmd interface {
	Run(ctx context.Context, o CmdOptions) error
}

type CmdOptions struct {
	RepositoryOwner string
	RepositoryName  string
	CommitMessage   string
	Paths           []string
}

type FileSystem interface {
	FindFiles(paths []string) ([]string, error)
	ReadAsBase64EncodedContent(filename string) (string, error)
}

type Logger interface {
	Infof(format string, v ...interface{})
	Debugf(format string, v ...interface{})
}

type GitHub interface {
	QueryRepository(ctx context.Context, in QueryRepositoryIn) (*QueryRepositoryOut, error)
	CreateBranch(ctx context.Context, branch NewBranch) error
	UpdateBranch(ctx context.Context, branch NewBranch, force bool) error
	CreateCommit(ctx context.Context, commit NewCommit) (git.CommitSHA, error)
	CreateTree(ctx context.Context, tree NewTree) (git.TreeSHA, error)
	CreateBlob(ctx context.Context, blob NewBlob) (git.BlobSHA, error)
}

type QueryRepositoryIn struct {
	Repository git.RepositoryID
}

type QueryRepositoryOut struct {
	CurrentUserName        string
	Repository             git.RepositoryID
	DefaultBranchName      git.BranchName
	DefaultBranchCommitSHA git.CommitSHA
	DefaultBranchTreeSHA   git.TreeSHA
}

type NewBranch struct {
	Repository git.RepositoryID
	BranchName git.BranchName
	CommitSHA  git.CommitSHA
}

type NewCommit struct {
	Repository      git.RepositoryID
	Message         git.CommitMessage
	ParentCommitSHA git.CommitSHA
	TreeSHA         git.TreeSHA
}

type NewTree struct {
	Repository  git.RepositoryID
	BaseTreeSHA git.TreeSHA
	Files       []git.File
}

type NewBlob struct {
	Repository git.RepositoryID
	Content    string // base64 encoded content
}
