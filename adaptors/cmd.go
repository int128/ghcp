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

const (
	envGitHubToken = "GITHUB_TOKEN"
	envGitHubAPI   = "GITHUB_API"
)

const (
	exitCodeOK                = 0
	exitCodeGenericError      = 10
	exitCodePreconditionError = 11
	exitCodeCopyError         = 20
)

func NewCmd(i Cmd) adaptors.Cmd {
	return &i
}

// Cmd interacts with command line interface.
type Cmd struct {
	dig.In
	CopyUseCase      usecases.CopyUseCase
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
		copyOptions
		Chdir       string // optional
		GitHubToken string
		GitHubAPI   string // optional
		Debug       bool
	}
	f.StringVarP(&o.RepositoryOwner, "owner", "u", "", "GitHub repository owner (mandatory)")
	f.StringVarP(&o.RepositoryName, "repo", "r", "", "GitHub repository name (mandatory)")
	f.StringVarP(&o.CommitMessage, "message", "m", "", "Commit message (mandatory)")
	f.StringVarP(&o.Branch, "branch", "b", "", "Branch name (default: default branch of repository)")
	f.StringVarP(&o.Chdir, "directory", "C", "", "Change to directory before copy")
	f.StringVar(&o.GitHubToken, "token", "", fmt.Sprintf("GitHub API token [$%s]", envGitHubToken))
	f.StringVar(&o.GitHubAPI, "api", "", fmt.Sprintf("GitHub API v3 URL (v4 will be inferred) [$%s]", envGitHubAPI))
	f.BoolVar(&o.NoFileMode, "no-file-mode", false, "Ignore executable bit of file and treat as 0644")
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
		c.Logger.Debugf("Debug enabled")
	}
	if o.Chdir != "" {
		if err := c.Env.Chdir(o.Chdir); err != nil {
			c.Logger.Errorf("Could not change to directory %s: %s", o.Chdir, err)
			return exitCodePreconditionError
		}
		c.Logger.Infof("Changed to directory %s", o.Chdir)
	}
	if o.GitHubToken == "" {
		o.GitHubToken = c.Env.Getenv(envGitHubToken)
		if o.GitHubToken != "" {
			c.Logger.Debugf("Using token from environment variable $%s", envGitHubToken)
		}
	}
	if o.GitHubToken == "" {
		c.Logger.Errorf("No GitHub API token. Set environment variable %s or --token option", envGitHubToken)
		return exitCodePreconditionError
	}
	if o.GitHubAPI == "" {
		o.GitHubAPI = c.Env.Getenv(envGitHubAPI)
		if o.GitHubAPI != "" {
			c.Logger.Debugf("Using GitHub Enterprise URL from environment variable $%s", envGitHubAPI)
		}
	}
	if err := c.GitHubClientInit.Init(infrastructure.GitHubClientInitOptions{
		Token: o.GitHubToken,
		URLv3: o.GitHubAPI,
	}); err != nil {
		c.Logger.Errorf("Could not connect to GitHub API: %s", err)
		return exitCodePreconditionError
	}

	return c.copy(ctx, o.copyOptions)
}

func (c *Cmd) copy(ctx context.Context, o copyOptions) int {
	if err := o.validate(); err != nil {
		c.Logger.Errorf("Invalid arguments: %s", err)
		return exitCodePreconditionError
	}
	if err := c.CopyUseCase.Do(ctx, usecases.CopyUseCaseIn{
		Repository: git.RepositoryID{
			Owner: o.RepositoryOwner,
			Name:  o.RepositoryName,
		},
		CommitMessage: git.CommitMessage(o.CommitMessage),
		BranchName:    git.BranchName(o.Branch),
		Paths:         o.Paths,
		NoFileMode:    o.NoFileMode,
		DryRun:        o.DryRun,
	}); err != nil {
		c.Logger.Errorf("Could not copy files: %s", err)
		c.Logger.Debugf("Stacktrace:\n%+v", err)
		return exitCodeCopyError
	}
	return exitCodeOK
}

type copyOptions struct {
	RepositoryOwner string
	RepositoryName  string
	CommitMessage   string
	Branch          string // optional
	Paths           []string
	NoFileMode      bool
	DryRun          bool
}

func (o *copyOptions) validate() error {
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
