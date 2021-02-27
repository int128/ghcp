package cmd

import (
	"context"
	"errors"
	"fmt"

	"github.com/int128/ghcp/pkg/git"
	"github.com/int128/ghcp/pkg/usecases/pullrequest"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const pullRequestCmdExample = ` To create a pull request from the feature branch to the default branch:
    ghcp pull-request -r OWNER/REPO -b feature --title TITLE --body BODY

  To create a pull request from the feature branch to the develop branch:
    ghcp pull-request -r OWNER/REPO -b feature --base develop --title TITLE --body BODY

  To create a pull request from the feature branch of the OWNER/REPO repository to the default branch of the upstream repository:
    ghcp pull-request -r OWNER/REPO -b feature --base-repo UPSTREAM/REPO --title TITLE --body BODY

  To create a pull request from the feature branch of the OWNER/REPO repository to the default branch of the upstream repository:
    ghcp pull-request -r OWNER/REPO -b feature --base-repo UPSTREAM/REPO --base feature --title TITLE --body BODY
`

func (r *Runner) newPullRequestCmd(ctx context.Context, gOpts *globalOptions) *cobra.Command {
	var o pullRequestOptions
	c := &cobra.Command{
		Use:     fmt.Sprintf("%s [flags] FILES...", pullRequestCmdName),
		Short:   "Create a pull request",
		Long:    `This creates a pull request. Do nothing if it already exists.`,
		Example: pullRequestCmdExample,
		RunE: func(_ *cobra.Command, args []string) error {
			if err := o.validate(); err != nil {
				return fmt.Errorf("invalid flag: %w", err)
			}
			headRepository, err := o.Head.repositoryID()
			if err != nil {
				return fmt.Errorf("invalid flag: %w", err)
			}

			baseRepository := headRepository
			if o.Base.RepositoryName != "" {
				baseRepository, err = o.Base.repositoryID()
				if err != nil {
					return fmt.Errorf("invalid flag: %w", err)
				}
			}

			ir, err := r.newInternalRunner(gOpts)
			if err != nil {
				return fmt.Errorf("error while bootstrap of the dependencies: %w", err)
			}
			in := pullrequest.Input{
				BaseRepository: baseRepository,
				BaseBranchName: git.BranchName(o.BaseBranchName),
				HeadRepository: headRepository,
				HeadBranchName: git.BranchName(o.HeadBranchName),
				Title:          o.Title,
				Body:           o.Body,
				Reviewer:       o.Reviewer,
				Draft:          o.Draft,
			}
			if err := ir.PullRequestUseCase.Do(ctx, in); err != nil {
				ir.Logger.Debugf("Stacktrace:\n%+v", err)
				return fmt.Errorf("could not create a pull request: %s", err)
			}
			return nil
		},
	}
	o.register(c.Flags())
	return c
}

type pullRequestOptions struct {
	Base repositoryOptions
	Head repositoryOptions

	BaseBranchName string
	HeadBranchName string
	Title          string
	Body           string
	Reviewer       string
	Draft          bool
}

func (o pullRequestOptions) validate() error {
	if o.HeadBranchName == "" || o.Title == "" {
		return errors.New("you need to set -b and --title")
	}
	return nil
}

func (o *pullRequestOptions) register(f *pflag.FlagSet) {
	f.StringVarP(&o.Head.RepositoryName, "head-repo", "r", "", "Head repository name, either -r OWNER/REPO or -u OWNER -r REPO (mandatory)")
	f.StringVarP(&o.Head.RepositoryOwner, "head-owner", "u", "", "Head repository owner")
	f.StringVarP(&o.HeadBranchName, "head", "b", "", "Head branch name (mandatory)")
	f.StringVar(&o.Base.RepositoryName, "base-repo", "", "Base repository name, either --base-repo OWNER/REPO or --base-owner OWNER --base-repo REPO (default: head)")
	f.StringVar(&o.Base.RepositoryOwner, "base-owner", "", "Base repository owner (default: head)")
	f.StringVar(&o.BaseBranchName, "base", "", "Base branch name (default: default branch of base repository)")
	f.StringVar(&o.Title, "title", "", "Title of a pull request (mandatory)")
	f.StringVar(&o.Body, "body", "", "Body of a pull request")
	f.StringVar(&o.Reviewer, "reviewer", "", "If set, request a review")
	f.BoolVar(&o.Draft, "draft", false, "If set, mark as a draft")
}
