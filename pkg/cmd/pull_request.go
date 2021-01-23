package cmd

import (
	"context"
	"fmt"

	"github.com/int128/ghcp/pkg/git"
	"github.com/int128/ghcp/usecases/pullrequest"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"golang.org/x/xerrors"
)

const pullRequestCmdExample = ` To create a pull request from the feature branch to the default branch:
    ghcp pull-request -u OWNER -r REPO -b feature --title TITLE --body BODY

  To create a pull request from the feature branch to the develop branch:
    ghcp pull-request -u OWNER -r REPO -b feature --base develop --title TITLE --body BODY

  To create a pull request from the feature branch of the OWNER/REPO repository to the default branch of the UPSTREAM/REPO repository:
    ghcp pull-request -u OWNER -r REPO -b feature --base-owner UPSTREAM --base-repo REPO --title TITLE --body BODY

  To create a pull request from the feature branch of the OWNER/REPO repository to the default branch of the UPSTREAM/REPO repository:
    ghcp pull-request -u OWNER -r REPO -b feature --base-owner UPSTREAM --base-repo REPO --base feature --title TITLE --body BODY
`

func (r *Runner) newPullRequestCmd(ctx context.Context, gOpts *globalOptions) *cobra.Command {
	var o pullRequestOptions
	c := &cobra.Command{
		Use:     fmt.Sprintf("%s [flags] FILES...", pullRequestCmdName),
		Short:   "Create a pull request",
		Long:    `This creates a pull request. Do nothing if it already exists.`,
		Example: pullRequestCmdExample,
		Args: func(*cobra.Command, []string) error {
			if o.HeadRepositoryOwner == "" ||
				o.HeadRepositoryName == "" ||
				o.HeadBranchName == "" ||
				o.Title == "" {
				return xerrors.New("mandatory options: -u, -r, -b, --title")
			}
			if o.BaseRepositoryOwner == "" {
				o.BaseRepositoryOwner = o.HeadRepositoryOwner
			}
			if o.BaseRepositoryName == "" {
				o.BaseRepositoryName = o.HeadRepositoryName
			}
			return nil
		},
		RunE: func(_ *cobra.Command, args []string) error {
			ir, err := r.newInternalRunner(gOpts)
			if err != nil {
				return xerrors.Errorf("error while bootstrap of the dependencies: %w", err)
			}
			in := pullrequest.Input{
				BaseRepository: git.RepositoryID{
					Owner: o.BaseRepositoryOwner,
					Name:  o.BaseRepositoryName,
				},
				BaseBranchName: git.BranchName(o.BaseBranchName),
				HeadRepository: git.RepositoryID{
					Owner: o.HeadRepositoryOwner,
					Name:  o.HeadRepositoryName,
				},
				HeadBranchName: git.BranchName(o.HeadBranchName),
				Title:          o.Title,
				Body:           o.Body,
			}
			if err := ir.PullRequestUseCase.Do(ctx, in); err != nil {
				ir.Logger.Debugf("Stacktrace:\n%+v", err)
				return xerrors.Errorf("could not create a pull request: %s", err)
			}
			return nil
		},
	}
	o.register(c.Flags())
	return c
}

type pullRequestOptions struct {
	BaseRepositoryOwner string
	BaseRepositoryName  string
	BaseBranchName      string
	HeadRepositoryOwner string
	HeadRepositoryName  string
	HeadBranchName      string
	Title               string
	Body                string
}

func (o *pullRequestOptions) register(f *pflag.FlagSet) {
	f.StringVarP(&o.HeadRepositoryOwner, "head-owner", "u", "", "Head repository owner (mandatory)")
	f.StringVarP(&o.HeadRepositoryName, "head-repo", "r", "", "Head repository name (mandatory)")
	f.StringVarP(&o.HeadBranchName, "head", "b", "", "Head branch name (mandatory)")
	f.StringVar(&o.BaseRepositoryOwner, "base-owner", "", "Base repository owner (default: head)")
	f.StringVar(&o.BaseRepositoryName, "base-repo", "", "Base repository name (default: head)")
	f.StringVar(&o.BaseBranchName, "base", "", "Base branch name (default: default branch of base repository)")
	f.StringVar(&o.Title, "title", "", "Title of a pull request (mandatory)")
	f.StringVar(&o.Body, "body", "", "Body of a pull request")
}
