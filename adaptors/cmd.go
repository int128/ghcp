package adaptors

import (
	"context"

	"github.com/int128/ghcp/adaptors/interfaces"
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

	if err := c.Push.Do(ctx, usecases.PushIn{
		RepositoryOwner: o.RepositoryOwner,
		RepositoryName:  o.RepositoryName,
		Paths:           o.Paths,
	}); err != nil {
		return errors.Wrapf(err, "error while commit and push")
	}
	return nil
}
