// Package di provides dependency injection.
package di

import (
	"github.com/google/go-github/v24/github"
	"github.com/int128/ghcp/adaptors"
	"github.com/int128/ghcp/usecases"
	"github.com/pkg/errors"
	"github.com/shurcooL/githubv4"
	"go.uber.org/dig"
)

//go:generate mockgen -destination mock_di/mock_di.go github.com/int128/ghcp/di Container

var dependencies = []interface{}{
	adaptors.NewCmd,
	adaptors.NewFileSystem,
	adaptors.NewLogger,
	adaptors.NewGitHub,

	usecases.NewPush,
}

// ExtraDependencies are given in runtime.
type ExtraDependencies struct {
	dig.Out
	GitHubV3 *github.Client
	GitHubV4 *githubv4.Client
}

// Container provides dependency injection.
type Container interface {
	Run(ExtraDependencies, interface{}) error
}

// New returns a new container with predefined dependencies.
func New() Container {
	return &container{}
}

type container struct{}

func (*container) Run(extraDependencies ExtraDependencies, runner interface{}) error {
	c := dig.New()
	for _, d := range dependencies {
		if err := c.Provide(d); err != nil {
			return errors.Wrapf(err, "error while providing predefined dependency %T", d)
		}
	}
	if err := c.Provide(func() ExtraDependencies { return extraDependencies }); err != nil {
		return errors.Wrapf(err, "error while providing extra dependencies")
	}
	if err := c.Invoke(runner); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
