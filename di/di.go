// Package di provides dependency injection.
package di

import (
	"github.com/int128/ghcp/adaptors"
	"github.com/int128/ghcp/usecases"
	"github.com/pkg/errors"
	"go.uber.org/dig"
)

var dependencies = []interface{}{
	adaptors.NewCmd,
	adaptors.NewGitHub,

	usecases.NewPush,
}

// New returns a new container with given dependencies.
func New(runtimeDependencies ...interface{}) (*dig.Container, error) {
	c := dig.New()
	for _, d := range dependencies {
		if err := c.Provide(d); err != nil {
			return nil, errors.Wrapf(err, "error while providing dependency %T", d)
		}
	}
	for _, d := range runtimeDependencies {
		if err := c.Provide(d); err != nil {
			return nil, errors.Wrapf(err, "error while providing dependency %T", d)
		}
	}
	return c, nil
}
