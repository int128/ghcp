package adaptors

import (
	"context"

	"github.com/int128/ghcp/git"
)

//go:generate mockgen -package mock_adaptors -destination ../mock_adaptors/mock_adaptors.go github.com/int128/ghcp/adaptors/interfaces Cmd,FileSystem,Env,LoggerConfig,GitHub

type Cmd interface {
	Run(ctx context.Context, args []string) int
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

type Env interface {
	Get(key string) string
}

type Logger interface {
	Errorf(format string, v ...interface{})
	Warnf(format string, v ...interface{})
	Infof(format string, v ...interface{})
	Debugf(format string, v ...interface{})
}

type LoggerConfig interface {
	SetDebug(debug bool)
}

type GitHub interface {
	QueryRepository(ctx context.Context, in QueryRepositoryIn) (*QueryRepositoryOut, error)
	CreateBranch(ctx context.Context, branch git.NewBranch) error
	UpdateBranch(ctx context.Context, branch git.NewBranch, force bool) error
	CreateCommit(ctx context.Context, commit git.NewCommit) (git.CommitSHA, error)
	QueryCommit(ctx context.Context, in QueryCommitIn) (*QueryCommitOut, error)
	CreateTree(ctx context.Context, tree git.NewTree) (git.TreeSHA, error)
	CreateBlob(ctx context.Context, blob git.NewBlob) (git.BlobSHA, error)
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

type QueryCommitIn struct {
	Repository git.RepositoryID
	CommitSHA  git.CommitSHA
}

type QueryCommitOut struct {
	ChangedFiles int
}
