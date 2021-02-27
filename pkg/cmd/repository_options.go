package cmd

import (
	"fmt"
	"strings"

	"github.com/int128/ghcp/pkg/git"
	"github.com/spf13/pflag"
)

type repositoryOptions struct {
	RepositoryOwner string
	RepositoryName  string
}

func (o *repositoryOptions) register(f *pflag.FlagSet) {
	f.StringVarP(&o.RepositoryName, "repo", "r", "", "Repository name, either -r OWNER/REPO or -u OWNER -r REPO (mandatory)")
	f.StringVarP(&o.RepositoryOwner, "owner", "u", "", "Repository owner")
}

func (o repositoryOptions) repositoryID() (git.RepositoryID, error) {
	if o.RepositoryName == "" {
		return git.RepositoryID{}, fmt.Errorf("you need to set the repository name")
	}
	if o.RepositoryOwner != "" {
		return git.RepositoryID{Owner: o.RepositoryOwner, Name: o.RepositoryName}, nil
	}
	s := strings.SplitN(o.RepositoryName, "/", 2)
	if len(s) != 2 {
		return git.RepositoryID{}, fmt.Errorf("you need to set OWNER/REPO when owner flag is omitted")
	}
	return git.RepositoryID{Owner: s[0], Name: s[1]}, nil
}
