// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package di

import (
	"github.com/int128/ghcp/adaptors"
	"github.com/int128/ghcp/adaptors/cmd"
	"github.com/int128/ghcp/adaptors/env"
	github2 "github.com/int128/ghcp/adaptors/github"
	"github.com/int128/ghcp/adaptors/logger"
	"github.com/int128/ghcp/infrastructure/github"
	"github.com/int128/ghcp/usecases/btc"
	"github.com/int128/ghcp/usecases/commit"
	"github.com/int128/ghcp/usecases/fork"
)

// Injectors from di.go:

func NewCmd() adaptors.Cmd {
	fileSystem := &env.FileSystem{}
	loggerLogger := &logger.Logger{}
	client := &github.Client{}
	gitHub := &github2.GitHub{
		Client: client,
		Logger: loggerLogger,
	}
	createBlobTreeCommit := &btc.CreateBlobTreeCommit{
		FileSystem: fileSystem,
		Logger:     loggerLogger,
		GitHub:     gitHub,
	}
	commitCommit := &commit.Commit{
		CreateBlobTreeCommit: createBlobTreeCommit,
		FileSystem:           fileSystem,
		Logger:               loggerLogger,
		GitHub:               gitHub,
	}
	commitToFork := &fork.CommitToFork{
		Commit: commitCommit,
		Logger: loggerLogger,
		GitHub: gitHub,
	}
	envEnv := &env.Env{}
	cmdCmd := &cmd.Cmd{
		Commit:           commitCommit,
		CommitToFork:     commitToFork,
		Env:              envEnv,
		Logger:           loggerLogger,
		LoggerConfig:     loggerLogger,
		GitHubClientInit: client,
	}
	return cmdCmd
}
