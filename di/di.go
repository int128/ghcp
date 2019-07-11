// Package di provides dependency injection.
package di

import (
	"github.com/int128/ghcp/adaptors/cmd"
	"github.com/int128/ghcp/adaptors/env"
	"github.com/int128/ghcp/adaptors/github"
	"github.com/int128/ghcp/adaptors/logger"
	"github.com/int128/ghcp/infrastructure"
	"github.com/int128/ghcp/usecases"
	"github.com/pkg/errors"
	"go.uber.org/dig"
)

var dependencies = []interface{}{
	usecases.NewUpdateBranch,
	usecases.NewCreateBranch,
	usecases.NewCommit,

	cmd.NewCmd,
	env.NewFileSystem,
	env.NewEnv,
	logger.NewLogger,
	github.NewGitHub,

	infrastructure.NewGitHubClient,
}

func Invoke(runner interface{}) error {
	c := dig.New()
	for _, d := range dependencies {
		if err := c.Provide(d); err != nil {
			return errors.Wrapf(err, "error while providing predefined dependency %T", d)
		}
	}
	if err := c.Invoke(runner); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
