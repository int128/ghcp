package cmd

import (
	"github.com/spf13/pflag"
	"golang.org/x/xerrors"

	"github.com/int128/ghcp/pkg/git"
)

type commitAttributeOptions struct {
	CommitMessage  string
	AuthorName     string
	AuthorEmail    string
	CommitterName  string
	CommitterEmail string
}

func (o *commitAttributeOptions) register(f *pflag.FlagSet) {
	f.StringVarP(&o.CommitMessage, "message", "m", "", "Commit message (mandatory)")
	f.StringVarP(&o.AuthorName, "author-name", "", "", "Author name (default: login name)")
	f.StringVarP(&o.AuthorEmail, "author-email", "", "", "Author email (default: login email)")
	f.StringVarP(&o.CommitterName, "committer-name", "", "", "Committer name (default: login name)")
	f.StringVarP(&o.CommitterEmail, "committer-email", "", "", "Committer email (default: login email)")
}

func (o *commitAttributeOptions) validate() error {
	if (o.AuthorName == "" && o.AuthorEmail != "") || (o.AuthorName != "" && o.AuthorEmail == "") {
		return xerrors.Errorf("you need to set both --author-name and --author-email")
	}
	if (o.CommitterName == "" && o.CommitterEmail != "") || (o.CommitterName != "" && o.CommitterEmail == "") {
		return xerrors.Errorf("you need to set both --committer-name and --committer-email")
	}
	return nil
}

func (o *commitAttributeOptions) committer() *git.CommitAuthor {
	if o.CommitterName != "" && o.CommitterEmail != "" {
		return &git.CommitAuthor{
			Name:  o.CommitterName,
			Email: o.CommitterEmail,
		}
	}
	return nil
}

func (o *commitAttributeOptions) author() *git.CommitAuthor {
	if o.AuthorName != "" && o.AuthorEmail != "" {
		return &git.CommitAuthor{
			Name:  o.AuthorName,
			Email: o.AuthorEmail,
		}
	}
	return nil
}
