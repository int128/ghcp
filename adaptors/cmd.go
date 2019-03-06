package adaptors

import (
	"context"
	"flag"
	"fmt"

	"github.com/int128/ghcp/adaptors/interfaces"
	"github.com/int128/ghcp/git"
	"github.com/int128/ghcp/infrastructure/interfaces"
	"github.com/int128/ghcp/usecases/interfaces"
	"github.com/pkg/errors"
	"go.uber.org/dig"
)

const usage = `Usage: %s [options] [file or directory...]

  ghcp commits and pushes files to a repository.
  It depends on GitHub API and works without git commands.

Options:
`

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
	f := flag.NewFlagSet(args[0], flag.ContinueOnError)
	f.Usage = func() {
		fmt.Fprintf(f.Output(), usage, args[0])
		f.PrintDefaults()
	}
	var o struct {
		pushOptions
		GitHubToken string
		Debug       bool
	}
	f.StringVar(&o.RepositoryOwner, "u", "", "GitHub repository owner (mandatory)")
	f.StringVar(&o.RepositoryName, "r", "", "GitHub repository name (mandatory)")
	f.StringVar(&o.CommitMessage, "m", "", "Commit message (mandatory)")
	f.StringVar(&o.GitHubToken, "token", "", fmt.Sprintf("GitHub API token [$%s]", envGitHubToken))
	f.BoolVar(&o.DryRun, "dry-run", false, "Upload files but do not update the branch actually")
	f.BoolVar(&o.Debug, "debug", false, "Show debug logs")

	if err := f.Parse(args[1:]); err != nil {
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
		c.Logger.Infof("Error: provide GitHub API token by $%s or -token", envGitHubToken)
		return 1
	}
	c.GitHubClientInit.Init(infrastructure.GitHubClientInitOptions{
		Token: o.GitHubToken,
	})

	if err := c.push(ctx, o.pushOptions); err != nil {
		c.Logger.Infof("Error: %s", err)
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
