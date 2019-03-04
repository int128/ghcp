package infrastructure

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/int128/ghcp/adaptors/interfaces"
	"github.com/int128/ghcp/di"
	"github.com/pkg/errors"
)

const usage = `Usage: %s [options] [path...]

  ghcp performs commit and push to a GitHub repository.
  It depends on GitHub API and works without git commands.

Options:
`

const envGitHubToken = "GITHUB_TOKEN"

// Run parses the arguments and bootstraps the application.
func Run(ctx context.Context, args []string) int {
	f := flag.NewFlagSet(args[0], flag.ContinueOnError)
	f.Usage = func() {
		fmt.Fprintf(f.Output(), usage, args[0])
		f.PrintDefaults()
	}
	var o options
	f.StringVar(&o.RepositoryOwner, "u", "", "GitHub repository owner")
	f.StringVar(&o.RepositoryName, "r", "", "GitHub repository name")
	f.StringVar(&o.CommitMessage, "m", "", "Commit message")
	f.StringVar(&o.GitHubToken, "token", "", fmt.Sprintf("GitHub API token [$%s]", envGitHubToken))

	if err := f.Parse(args[1:]); err != nil {
		return 1
	}
	o.Paths = f.Args()
	if o.GitHubToken == "" {
		o.GitHubToken = os.Getenv(envGitHubToken)
	}
	if o.GitHubToken == "" {
		log.Printf("Error: provide GitHub API token by $%s or -token", envGitHubToken)
		return 1
	}

	c, err := di.New(GitHubClientFactory(o.GitHubToken))
	if err != nil {
		log.Printf("Error: %s", err)
		return 1
	}
	if err := c.Invoke(func(cmd adaptors.Cmd) error {
		if err := cmd.Run(ctx, o.CmdOptions); err != nil {
			return errors.WithStack(err)
		}
		return nil
	}); err != nil {
		log.Printf("Error: %s", err)
		return 1
	}
	return 0
}

type options struct {
	adaptors.CmdOptions

	GitHubToken string
}
