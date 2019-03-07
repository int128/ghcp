package adaptors

import (
	"context"
	"fmt"
	"strings"

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

const (
	exitCodeOK                = 0
	exitCodeGenericError      = 10
	exitCodePreconditionError = 11
	exitCodePushError         = 20
)

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

// Run parses the arguments and executes the use case.
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
			return exitCodeGenericError
		}
		c.Logger.Errorf("Invalid arguments: %s", err)
		return exitCodeGenericError
	}
	o.Paths = f.Args()

	if o.Debug {
		c.LoggerConfig.SetDebug(true)
	}
	if o.GitHubToken == "" {
		o.GitHubToken = c.Env.Getenv(envGitHubToken)
	}
	if o.GitHubToken == "" {
		c.Logger.Errorf("No GitHub API token. Set $%s or -token", envGitHubToken)
		return exitCodePreconditionError
	}
	c.GitHubClientInit.Init(infrastructure.GitHubClientInitOptions{
		Token: o.GitHubToken,
	})

	return c.push(ctx, o.pushOptions)
}

func (c *Cmd) push(ctx context.Context, o pushOptions) int {
	if err := o.validate(); err != nil {
		c.Logger.Errorf("Invalid arguments: %s", err)
		return exitCodePreconditionError
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
		c.Logger.Errorf("Could not push files: %s", err)
		return exitCodePushError
	}
	return exitCodeOK
}

type pushOptions struct {
	RepositoryOwner string
	RepositoryName  string
	CommitMessage   string
	Paths           []string
	DryRun          bool
}

func (o *pushOptions) validate() error {
	var msg []string
	if o.RepositoryOwner == "" {
		msg = append(msg, "GitHub repository owner")
	}
	if o.RepositoryName == "" {
		msg = append(msg, "GitHub repository name")
	}
	if o.CommitMessage == "" {
		msg = append(msg, "commit message")
	}
	if len(o.Paths) == 0 {
		msg = append(msg, "one or more paths")
	}
	if len(msg) > 0 {
		return errors.Errorf("you need to set %s", strings.Join(msg, ", "))
	}
	return nil
}
