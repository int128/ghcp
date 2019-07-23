package cmd

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/google/wire"
	"github.com/int128/ghcp/adaptors"
	"github.com/int128/ghcp/git"
	"github.com/int128/ghcp/infrastructure"
	"github.com/int128/ghcp/usecases"
	"github.com/spf13/cobra"
	"golang.org/x/xerrors"
)

var Set = wire.NewSet(
	wire.Struct(new(Cmd), "*"),
	wire.Bind(new(adaptors.Cmd), new(*Cmd)),
)

const (
	envGitHubToken = "GITHUB_TOKEN"
	envGitHubAPI   = "GITHUB_API"

	exitCodeOK    = 0
	exitCodeError = 1

	commitCmdName = "commit"
)

// Cmd interacts with command line interface.
type Cmd struct {
	Commit           usecases.Commit
	Env              adaptors.Env
	Logger           adaptors.Logger
	LoggerConfig     adaptors.LoggerConfig
	GitHubClientInit infrastructure.GitHubClientInit
}

// Run parses the arguments and executes the use-case.
func (c *Cmd) Run(ctx context.Context, args []string) int {
	rootCmd := newRootCmd(filepath.Base(args[0]), c)
	commitCmd := newCommitCmd(ctx, c)
	rootCmd.AddCommand(commitCmd)

	rootCmd.SetArgs(args[1:])
	if err := rootCmd.Execute(); err != nil {
		return exitCodeError
	}
	return exitCodeOK
}

type globalOptions struct {
	Chdir       string
	GitHubToken string
	GitHubAPI   string // optional
	Debug       bool
}

func newRootCmd(name string, cmd *Cmd) *cobra.Command {
	var o globalOptions
	c := &cobra.Command{
		Use:          name,
		Short:        "A command to commit files to a GitHub repository",
		SilenceUsage: true,
		PersistentPreRunE: func(*cobra.Command, []string) error {
			if o.Debug {
				cmd.LoggerConfig.SetDebug(true)
				cmd.Logger.Debugf("Debug enabled")
			}
			if o.Chdir != "" {
				if err := cmd.Env.Chdir(o.Chdir); err != nil {
					return xerrors.Errorf("could not change to directory %s: %w", o.Chdir, err)
				}
				cmd.Logger.Infof("Changed to directory %s", o.Chdir)
			}
			if o.GitHubToken == "" {
				o.GitHubToken = cmd.Env.Getenv(envGitHubToken)
				if o.GitHubToken != "" {
					cmd.Logger.Debugf("Using token from environment variable $%s", envGitHubToken)
				}
			}
			if o.GitHubToken == "" {
				return xerrors.Errorf("no GitHub API token. Set environment variable %s or --token option", envGitHubToken)
			}
			if o.GitHubAPI == "" {
				o.GitHubAPI = cmd.Env.Getenv(envGitHubAPI)
				if o.GitHubAPI != "" {
					cmd.Logger.Debugf("Using GitHub Enterprise URL from environment variable $%s", envGitHubAPI)
				}
			}
			if err := cmd.GitHubClientInit.Init(infrastructure.GitHubClientInitOptions{
				Token: o.GitHubToken,
				URLv3: o.GitHubAPI,
			}); err != nil {
				return xerrors.Errorf("could not connect to GitHub API: %w", err)
			}
			return nil
		},
	}
	f := c.PersistentFlags()
	f.StringVarP(&o.Chdir, "directory", "C", "", "Change to directory before operation")
	f.StringVar(&o.GitHubToken, "token", "", fmt.Sprintf("GitHub API token [$%s]", envGitHubToken))
	f.StringVar(&o.GitHubAPI, "api", "", fmt.Sprintf("GitHub API v3 URL (v4 will be inferred) [$%s]", envGitHubAPI))
	f.BoolVar(&o.Debug, "debug", false, "Show debug logs")
	return c
}

type commitOptions struct {
	RepositoryOwner string
	RepositoryName  string
	CommitMessage   string
	BranchName      string
	ParentRef       string
	NoParent        bool
	NoFileMode      bool
	DryRun          bool
}

const commitCmdExample = `  To commit files to the default branch:
    ghcp commit -u OWNER -r REPO -m MESSAGE FILES...

  To commit files to the branch:
    ghcp commit -u OWNER -r REPO -b BRANCH -m MESSAGE FILES...

  If the branch does not exist, ghcp creates a branch from the default branch.
  It the branch exists, ghcp updates the branch by fast-forward.

  To commit files to a new branch from the parent branch:
    ghcp commit -u OWNER -r REPO -b BRANCH --parent PARENT -m MESSAGE FILES...

  If the branch exists, ghcp cannot update the branch by fast-forward and will fail.

  To commit files to a new branch without any parent:
    ghcp commit -u OWNER -r REPO -b BRANCH --no-parent -m MESSAGE FILES...

  If the branch exists, ghcp cannot update the branch by fast-forward and will fail.`

func newCommitCmd(ctx context.Context, cmd *Cmd) *cobra.Command {
	var o commitOptions
	c := &cobra.Command{
		Use:     fmt.Sprintf("%s [flags] FILES...", commitCmdName),
		Short:   "Commit files to the branch",
		Long:    `This commits the files to the branch. This will create a branch if it does not exist.`,
		Example: commitCmdExample,
		Args: func(*cobra.Command, []string) error {
			if o.ParentRef != "" && o.NoParent {
				return xerrors.Errorf("do not set both --parent and --no-parent")
			}
			return nil
		},
		RunE: func(_ *cobra.Command, args []string) error {
			in := usecases.CommitIn{
				Repository: git.RepositoryID{
					Owner: o.RepositoryOwner,
					Name:  o.RepositoryName,
				},
				BranchName: git.BranchName(o.BranchName),
				ParentOfBranch: usecases.ParentOfBranch{
					FastForward: o.ParentRef == "" && !o.NoParent,
					NoParent:    o.NoParent,
					FromRef:     git.RefName(o.ParentRef),
				},
				CommitMessage: git.CommitMessage(o.CommitMessage),
				Paths:         args,
				NoFileMode:    o.NoFileMode,
				DryRun:        o.DryRun,
			}
			if err := cmd.Commit.Do(ctx, in); err != nil {
				cmd.Logger.Debugf("Stacktrace:\n%+v", err)
				return xerrors.Errorf("could not commit the files: %s", err)
			}
			return nil
		},
	}
	f := c.Flags()
	f.StringVarP(&o.RepositoryOwner, "owner", "u", "", "GitHub repository owner (mandatory)")
	f.StringVarP(&o.RepositoryName, "repo", "r", "", "GitHub repository name (mandatory)")
	f.StringVarP(&o.CommitMessage, "message", "m", "", "Commit message (mandatory)")
	f.StringVarP(&o.BranchName, "branch", "b", "", "Name of the branch to create or update (default: the default branch of repository)")
	f.StringVar(&o.ParentRef, "parent", "", "Create a commit from the parent branch/tag (default: fast-forward)")
	f.BoolVar(&o.NoParent, "no-parent", false, "Create a commit without a parent")
	f.BoolVar(&o.NoFileMode, "no-file-mode", false, "Ignore executable bit of file and treat as 0644")
	f.BoolVar(&o.DryRun, "dry-run", false, "Upload files but do not update the branch actually")
	return c
}
