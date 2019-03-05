package infrastructure

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/int128/ghcp/adaptors/interfaces"
	"github.com/int128/ghcp/infrastructure/interfaces"
	"go.uber.org/dig"
)

const usage = `Usage: %s [options] [file or directory...]

  ghcp commits and pushes files to a repository.
  It depends on GitHub API and works without git commands.

Options:
`

const envGitHubToken = "GITHUB_TOKEN"

func NewCmd(in cmdIn) infrastructure.Cmd {
	return &Cmd{in, os.Getenv}
}

type cmdIn struct {
	dig.In
	Cmd                adaptors.Cmd
	GitHubClientConfig infrastructure.GitHubClientConfig
}

type Cmd struct {
	cmdIn
	Getenv func(string) string
}

// Run parses the arguments and bootstraps the application.
func (cmd *Cmd) Run(ctx context.Context, args []string) int {
	f := flag.NewFlagSet(args[0], flag.ContinueOnError)
	f.Usage = func() {
		fmt.Fprintf(f.Output(), usage, args[0])
		f.PrintDefaults()
	}
	var o options
	f.StringVar(&o.RepositoryOwner, "u", "", "GitHub repository owner (mandatory)")
	f.StringVar(&o.RepositoryName, "r", "", "GitHub repository name (mandatory)")
	f.StringVar(&o.CommitMessage, "m", "", "Commit message (mandatory)")
	f.StringVar(&o.GitHubToken, "token", "", fmt.Sprintf("GitHub API token [$%s]", envGitHubToken))

	if err := f.Parse(args[1:]); err != nil {
		return 1
	}
	o.Paths = f.Args()
	if o.GitHubToken == "" {
		o.GitHubToken = cmd.Getenv(envGitHubToken)
	}
	if o.GitHubToken == "" {
		log.Printf("Error: provide GitHub API token by $%s or -token", envGitHubToken)
		return 1
	}
	cmd.GitHubClientConfig.SetToken(o.GitHubToken)

	if err := cmd.Cmd.Run(ctx, o.CmdOptions); err != nil {
		log.Printf("Error: %s", err)
		return 1
	}
	return 0
}

type options struct {
	adaptors.CmdOptions

	GitHubToken string
}
