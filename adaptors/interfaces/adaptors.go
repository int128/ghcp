package adaptors

import "context"

//go:generate mockgen -package mock_adaptors -destination ../mock_adaptors/mock_adaptors.go github.com/int128/ghcp/adaptors/interfaces GitHub

type Cmd interface {
	Run(ctx context.Context, o CmdOptions) error
}

type CmdOptions struct {
	RepositoryOwner string
	RepositoryName  string
	Paths           []string
}

type GitHub interface {
	GetRepository(ctx context.Context, in GetRepositoryIn) (*GetRepositoryOut, error)
}

type GetRepositoryIn struct {
	Owner string
	Name  string
}

type GetRepositoryOut struct {
	CurrentUserName   string
	DefaultBranchName string
}
