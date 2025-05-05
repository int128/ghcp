// Package cmd parses command line args and runs the corresponding use-case.
package cmd

import (
	"context"
	"fmt"
	"log"
	"log/slog"

	"github.com/google/wire"
	"github.com/int128/ghcp/pkg/env"
	"github.com/int128/ghcp/pkg/github/client"
	"github.com/int128/ghcp/pkg/usecases/commit"
	"github.com/int128/ghcp/pkg/usecases/forkcommit"
	"github.com/int128/ghcp/pkg/usecases/pullrequest"
	"github.com/int128/ghcp/pkg/usecases/release"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	envGitHubToken = "GITHUB_TOKEN"
	envGitHubAPI   = "GITHUB_API"

	exitCodeOK    = 0
	exitCodeError = 1

	commitCmdName      = "commit"
	emptyCommitCmdName = "empty-commit"
	forkCommitCmdName  = "fork-commit"
	pullRequestCmdName = "pull-request"
	releaseCmdName     = "release"
)

var Set = wire.NewSet(
	wire.Bind(new(Interface), new(*Runner)),
	wire.Struct(new(Runner), "*"),
	wire.Struct(new(InternalRunner), "*"),
)

type Interface interface {
	Run(args []string, version string) int
}

// Runner is the entry point for the command line application.
// It bootstraps the InternalRunner and runs the specified use-case.
type Runner struct {
	Env               env.Interface
	NewGitHub         client.NewFunc
	NewInternalRunner NewInternalRunnerFunc
}

// Run parses the command line args and runs the corresponding use-case.
func (r *Runner) Run(args []string, version string) int {
	ctx := context.Background()

	var o globalOptions
	rootCmd := r.newRootCmd(&o)
	commitCmd := r.newCommitCmd(ctx, &o)
	rootCmd.AddCommand(commitCmd)
	emptyCommitCmd := r.newEmptyCommitCmd(ctx, &o)
	rootCmd.AddCommand(emptyCommitCmd)
	forkCommitCmd := r.newForkCommitCmd(ctx, &o)
	rootCmd.AddCommand(forkCommitCmd)
	pullRequestCmd := r.newPullRequestCmd(ctx, &o)
	rootCmd.AddCommand(pullRequestCmd)
	releaseCmd := r.newReleaseCmd(ctx, &o)
	rootCmd.AddCommand(releaseCmd)

	rootCmd.Version = version
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

func (o *globalOptions) register(f *pflag.FlagSet) {
	f.StringVarP(&o.Chdir, "directory", "C", "", "Change to directory before operation")
	f.StringVar(&o.GitHubToken, "token", "", fmt.Sprintf("GitHub API token [$%s]", envGitHubToken))
	f.StringVar(&o.GitHubAPI, "api", "", fmt.Sprintf("GitHub API v3 URL (v4 will be inferred) [$%s]", envGitHubAPI))
	f.BoolVar(&o.Debug, "debug", false, "Show debug logs")
}

func (r *Runner) newRootCmd(o *globalOptions) *cobra.Command {
	c := &cobra.Command{
		Use:          "ghcp",
		Short:        "A command to commit files to a GitHub repository",
		SilenceUsage: true,
	}
	o.register(c.PersistentFlags())
	return c
}

type NewInternalRunnerFunc func(client.Interface) *InternalRunner

// InternalRunner has the set of use-cases.
type InternalRunner struct {
	CommitUseCase      commit.Interface
	ForkCommitUseCase  forkcommit.Interface
	PullRequestUseCase pullrequest.Interface
	ReleaseUseCase     release.Interface
}

func (r *Runner) newInternalRunner(o *globalOptions) (*InternalRunner, error) {
	log.SetFlags(log.Lmicroseconds)
	if o.Debug {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}
	if o.Chdir != "" {
		if err := r.Env.Chdir(o.Chdir); err != nil {
			return nil, fmt.Errorf("could not change to directory %s: %w", o.Chdir, err)
		}
		slog.Info("Changed to directory", "directory", o.Chdir)
	}
	if o.GitHubToken == "" {
		o.GitHubToken = r.Env.Getenv(envGitHubToken)
		if o.GitHubToken != "" {
			slog.Debug("Using token from environment variable", "variable", envGitHubToken)
		}
	}
	if o.GitHubToken == "" {
		return nil, fmt.Errorf("no GitHub API token. Set environment variable %s or --token option", envGitHubToken)
	}
	if o.GitHubAPI == "" {
		o.GitHubAPI = r.Env.Getenv(envGitHubAPI)
		if o.GitHubAPI != "" {
			slog.Debug("Using GitHub Enterprise URL from environment variable", "variable", envGitHubAPI)
		}
	}
	gh, err := r.NewGitHub(client.Option{
		Token: o.GitHubToken,
		URLv3: o.GitHubAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("could not connect to GitHub API: %w", err)
	}
	return r.NewInternalRunner(gh), nil
}
