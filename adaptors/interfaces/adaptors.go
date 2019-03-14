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
	FindFiles(paths []string) ([]File, error)
	ReadAsBase64EncodedContent(filename string) (string, error)
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
	QueryForUpdateBranch(ctx context.Context, in QueryForUpdateBranchIn) (*QueryForUpdateBranchOut, error)
	QueryForCreateBranch(ctx context.Context, in QueryForCreateBranchIn) (*QueryForCreateBranchOut, error)
	CreateBranch(ctx context.Context, branch git.NewBranch) error
	UpdateBranch(ctx context.Context, branch git.NewBranch, force bool) error
	CreateCommit(ctx context.Context, commit git.NewCommit) (git.CommitSHA, error)
	QueryCommit(ctx context.Context, in QueryCommitIn) (*QueryCommitOut, error)
	CreateTree(ctx context.Context, tree git.NewTree) (git.TreeSHA, error)
	CreateBlob(ctx context.Context, blob git.NewBlob) (git.BlobSHA, error)
}

type QueryForUpdateBranchIn struct {
	Repository git.RepositoryID
	BranchName git.BranchName // optional
}

type QueryForUpdateBranchOut struct {
	CurrentUserName        string
	Repository             git.RepositoryID
	DefaultBranchName      git.BranchName
	DefaultBranchCommitSHA git.CommitSHA
	DefaultBranchTreeSHA   git.TreeSHA
	BranchCommitSHA        git.CommitSHA // empty if the branch does not exist
	BranchTreeSHA          git.TreeSHA   // empty if the branch does not exist
}

type QueryForCreateBranchIn struct {
	Repository    git.RepositoryID
	ParentRef     git.RefName // optional
	NewBranchName git.BranchName
}

type QueryForCreateBranchOut struct {
	CurrentUserName        string
	Repository             git.RepositoryID
	DefaultBranchRefName   git.RefQualifiedName
	DefaultBranchCommitSHA git.CommitSHA
	DefaultBranchTreeSHA   git.TreeSHA
	ParentRefName          git.RefQualifiedName // empty if the parent ref does not exist
	ParentRefCommitSHA     git.CommitSHA        // empty if the parent ref does not exist
	ParentRefTreeSHA       git.TreeSHA          // empty if the parent ref does not exist
	NewBranchExists        bool
}

type QueryCommitIn struct {
	Repository git.RepositoryID
	CommitSHA  git.CommitSHA
}

type QueryCommitOut struct {
	ChangedFiles int
}
