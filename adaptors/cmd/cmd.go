package cmd

import (
	"context"
	"fmt"

	"github.com/int128/ghcp/adaptors"
	"github.com/int128/ghcp/git"
	"github.com/int128/ghcp/infrastructure"
	"github.com/int128/ghcp/usecases/interfaces"
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
	exitCodeUseCaseError      = 20
)

func NewCmd(i Cmd) adaptors.Cmd {
	return &i
}

// Cmd interacts with command line interface.
type Cmd struct {
	dig.In
	UpdateBranch     usecases.UpdateBranch
	CreateBranch     usecases.CreateBranch
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
	var o cmdOptions
	f.StringVarP(&o.RepositoryOwner, "owner", "u", "", "GitHub repository owner (mandatory)")
	f.StringVarP(&o.RepositoryName, "repo", "r", "", "GitHub repository name (mandatory)")
	f.StringVarP(&o.CommitMessage, "message", "m", "", "Commit message (mandatory)")
	f.StringVarP(&o.UpdateBranch, "branch", "b", "", "Update the branch (default: default branch of repository)")
	f.StringVarP(&o.NewBranch, "new-branch", "B", "", "Create a branch")
	f.StringVar(&o.ParentRef, "parent", "", "Create a commit from the parent branch or tag (default: default branch of repository)")
	f.BoolVar(&o.NoParent, "no-parent", false, "Create a commit without a parent")
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

	if o.UpdateBranch != "" && o.NewBranch != "" {
		c.Logger.Errorf("Do not set both --branch and --new-branch")
		return exitCodePreconditionError
	}
	if o.NewBranch != "" {
		return c.runCreateBranch(ctx, o)
	}
	return c.runUpdateBranch(ctx, o)
}

func (c *Cmd) runCreateBranch(ctx context.Context, o cmdOptions) int {
	if o.ParentRef != "" && o.NoParent {
		c.Logger.Errorf("Do not set both --parent and --no-parent")
		return exitCodePreconditionError
	}

	in := usecases.CreateBranchIn{
		Repository: git.RepositoryID{
			Owner: o.RepositoryOwner,
			Name:  o.RepositoryName,
		},
		NewBranchName: git.BranchName(o.NewBranch),
		ParentOfNewBranch: usecases.ParentOfNewBranch{
			NoParent:          o.NoParent,
			FromDefaultBranch: o.ParentRef == "" && !o.NoParent,
			FromRef:           git.RefName(o.ParentRef),
		},
		CommitMessage: git.CommitMessage(o.CommitMessage),
		Paths:         o.Paths,
		NoFileMode:    o.NoFileMode,
		DryRun:        o.DryRun,
	}
	if err := c.CreateBranch.Do(ctx, in); err != nil {
		c.Logger.Errorf("Could not copy files: %s", err)
		c.Logger.Debugf("Stacktrace:\n%+v", err)
		return exitCodeUseCaseError
	}
	return exitCodeOK
}

func (c *Cmd) runUpdateBranch(ctx context.Context, o cmdOptions) int {
	if o.ParentRef != "" {
		c.Logger.Errorf("Do not set --parent on updating the branch")
		return exitCodePreconditionError
	}

	in := usecases.UpdateBranchIn{
		Repository: git.RepositoryID{
			Owner: o.RepositoryOwner,
			Name:  o.RepositoryName,
		},
		BranchName:    git.BranchName(o.UpdateBranch),
		CommitMessage: git.CommitMessage(o.CommitMessage),
		Paths:         o.Paths,
		NoFileMode:    o.NoFileMode,
		DryRun:        o.DryRun,
	}
	if err := c.UpdateBranch.Do(ctx, in); err != nil {
		c.Logger.Errorf("Could not copy files: %s", err)
		c.Logger.Debugf("Stacktrace:\n%+v", err)
		return exitCodeUseCaseError
	}
	return exitCodeOK
}

type cmdOptions struct {
	envOptions
	useCaseOptions
}

type envOptions struct {
	Chdir       string
	GitHubToken string
	GitHubAPI   string // optional
	Debug       bool
}

type useCaseOptions struct {
	RepositoryOwner string
	RepositoryName  string
	CommitMessage   string
	UpdateBranch    string
	NewBranch       string
	ParentRef       string
	NoParent        bool
	Paths           []string
	NoFileMode      bool
	DryRun          bool
}
