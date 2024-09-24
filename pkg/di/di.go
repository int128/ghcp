//go:build wireinject
// +build wireinject

// Package di provides dependency injection.
package di

import (
	"github.com/google/wire"
	"github.com/int128/ghcp/pkg/cmd"
	"github.com/int128/ghcp/pkg/env"
	"github.com/int128/ghcp/pkg/fs"
	"github.com/int128/ghcp/pkg/github"
	"github.com/int128/ghcp/pkg/github/client"
	"github.com/int128/ghcp/pkg/logger"
	"github.com/int128/ghcp/pkg/usecases/commit"
	"github.com/int128/ghcp/pkg/usecases/forkcommit"
	"github.com/int128/ghcp/pkg/usecases/gitobject"
	"github.com/int128/ghcp/pkg/usecases/pullrequest"
	"github.com/int128/ghcp/pkg/usecases/release"
)

func NewCmd() cmd.Interface {
	wire.Build(
		cmd.Set,
		logger.Set,
		client.Set,
		env.Set,

		wire.Value(cmd.NewInternalRunnerFunc(NewCmdInternalRunner)),
	)
	return nil
}

func NewCmdInternalRunner(logger.Interface, client.Interface) *cmd.InternalRunner {
	wire.Build(
		cmd.Set,
		fs.Set,
		github.Set,

		gitobject.Set,
		commit.Set,
		forkcommit.Set,
		pullrequest.Set,
		release.Set,
	)
	return nil
}
