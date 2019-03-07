package adaptors

import (
	"context"
	"fmt"

	"github.com/int128/ghcp/adaptors/interfaces"
	"github.com/int128/ghcp/git"
	"github.com/int128/ghcp/infrastructure/interfaces"
	"github.com/int128/ghcp/usecases/interfaces"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"go.uber.org/dig"
)

const usage = `Help:

Usage: %s [options] [file or directory...]

  ghcp copies files to a GitHub repository.
  It depends on GitHub API and works without git commands.

Options:
%s`

const envGitHubToken = "GITHUB_TOKEN"

func NewCmd(i Cmd) adaptors.Cmd {
	return &i
}

// Cmd interacts with command line interface.
type Cmd struct {
	dig.In
	Push             usecases.Push
	Env              adaptors.Env
	Logger           adaptors.Logger
	LoggerConfig     adaptors.LoggerConfig
	GitHubClientInit infrastructure.GitHubClientInit
}

func (c *Cmd) Run(ctx context.Context, args []string) int {
	f := pflag.NewFlagSet(args[0], pflag.ContinueOnError)
	f.Usage = func() {
		c.Logger.Infof(usage, args[0], f.FlagUsages())
	}
	var o struct {
		pushOptions
		GitHubToken string
		Debug       bool
	}
	f.StringVarP(&o.RepositoryOwner, "owner", "u", "", "GitHub repository owner (mandatory)")
	f.StringVarP(&o.RepositoryName, "repo", "r", "", "GitHub repository name (mandatory)")
	f.StringVarP(&o.CommitMessage, "message", "m", "", "Commit message (mandatory)")
	f.StringVar(&o.GitHubToken, "token", "", fmt.Sprintf("GitHub API token [$%s]", envGitHubToken))
	f.BoolVar(&o.DryRun, "dry-run", false, "Upload files but do not update the branch actually")
	f.BoolVar(&o.Debug, "debug", false, "Show debug logs")

	if err := f.Parse(args[1:]); err != nil {
		if err == pflag.ErrHelp {
			return 1
		}
		c.Logger.Errorf("Invalid arguments: %s", err)
		return 1
	}
	o.Paths = f.Args()
	if o.Debug {
		c.LoggerConfig.SetDebug(true)
	}
	if o.GitHubToken == "" {
		o.GitHubToken = c.Env.Get(envGitHubToken)
	}
	if o.GitHubToken == "" {
		c.Logger.Errorf("Error: provide GitHub API token by $%s or -token", envGitHubToken)
		return 1
	}
	c.GitHubClientInit.Init(infrastructure.GitHubClientInitOptions{
		Token: o.GitHubToken,
	})

	if err := c.push(ctx, o.pushOptions); err != nil {
		c.Logger.Errorf("Error: %s", err)
		return 1
	}
	return 0
}

type pushOptions struct {
	RepositoryOwner string
	RepositoryName  string
	CommitMessage   string
	Paths           []string
	DryRun          bool
}

func (c *Cmd) push(ctx context.Context, o pushOptions) error {
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
		DryRun:        o.DryRun,
	}); err != nil {
		return errors.Wrapf(err, "error while commit and push")
	}
	return nil
}
