package adaptors

import (
	"context"

	"github.com/int128/ghcp/adaptors/interfaces"
	"github.com/int128/ghcp/git"
	"github.com/int128/ghcp/usecases/interfaces"
	"github.com/pkg/errors"
	"go.uber.org/dig"
)

func NewCmd(i Cmd) adaptors.Cmd {
	return &i
}

// Cmd represents a controller for command line interface.
type Cmd struct {
	dig.In
	Push usecases.Push
}

func (c *Cmd) Run(ctx context.Context, o adaptors.CmdOptions) error {
	if o.RepositoryOwner == "" {
		return errors.New("provide GitHub repository owner")
	}
	if o.RepositoryName == "" {
		return errors.New("provide GitHub repository name")
	}
	if o.CommitMessage == "" {
		return errors.New("provide commit message")
	}
	if len(o.Paths) == 0 {
		return errors.New("nothing to commit; provide one or more paths")
	}

	if err := c.Push.Do(ctx, usecases.PushIn{
		Repository: git.RepositoryID{
			Owner: o.RepositoryOwner,
			Name:  o.RepositoryName,
		},
		CommitMessage: git.CommitMessage(o.CommitMessage),
		Paths:         o.Paths,
	}); err != nil {
		return errors.Wrapf(err, "error while commit and push")
	}
	return nil
}
