// Package adaptors bridges between external interfaces and use cases.
package adaptors

import (
	"context"

	"github.com/int128/ghcp/git"
)

//go:generate mockgen -destination mock_adaptors/mock_adaptors.go github.com/int128/ghcp/adaptors Cmd,FileSystem,Env,LoggerConfig,GitHub

type Cmd interface {
	Run(ctx context.Context, args []string, version string) int
}

type CmdOptions struct {
	RepositoryOwner string
	RepositoryName  string
	CommitMessage   string
	Paths           []string
}

type FileSystem interface {
	FindFiles(paths []string, filter FindFilesFilter) ([]File, error)
	ReadAsBase64EncodedContent(filename string) (string, error)
}

// FindFilesFilter is an interface to filter directories and files.
type FindFilesFilter interface {
	SkipDir(path string) bool     // If true, it skips entering the directory
	ExcludeFile(path string) bool // If true, it excludes the file from the result
}

type File struct {
	Path       string
	Executable bool
}

type Env interface {
	Getenv(key string) string
	Chdir(dir string) error
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
	CreateFork(ctx context.Context, id git.RepositoryID) (*git.RepositoryID, error)
	QueryForCommitToBranch(ctx context.Context, in QueryForCommitToBranchIn) (*QueryForCommitToBranchOut, error)
	CreateBranch(ctx context.Context, branch git.NewBranch) error
	UpdateBranch(ctx context.Context, branch git.NewBranch, force bool) error
	CreateCommit(ctx context.Context, commit git.NewCommit) (git.CommitSHA, error)
	QueryCommit(ctx context.Context, in QueryCommitIn) (*QueryCommitOut, error)
	CreateTree(ctx context.Context, tree git.NewTree) (git.TreeSHA, error)
	CreateBlob(ctx context.Context, blob git.NewBlob) (git.BlobSHA, error)
}

type QueryForCommitToBranchIn struct {
	ParentRepository git.RepositoryID
	ParentRef        git.RefName // optional
	TargetRepository git.RepositoryID
	TargetBranchName git.BranchName // optional
}

type QueryForCommitToBranchOut struct {
	CurrentUserName              string
	ParentDefaultBranchCommitSHA git.CommitSHA
	ParentDefaultBranchTreeSHA   git.TreeSHA
	ParentRefCommitSHA           git.CommitSHA // empty if the parent ref does not exist
	ParentRefTreeSHA             git.TreeSHA   // empty if the parent ref does not exist
	TargetRepository             git.RepositoryID
	TargetDefaultBranchName      git.BranchName
	TargetBranchCommitSHA        git.CommitSHA // empty if the branch does not exist
	TargetBranchTreeSHA          git.TreeSHA   // empty if the branch does not exist
}

type QueryCommitIn struct {
	Repository git.RepositoryID
	CommitSHA  git.CommitSHA
}

type QueryCommitOut struct {
	ChangedFiles int
}
