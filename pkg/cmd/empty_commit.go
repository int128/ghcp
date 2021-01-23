package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"golang.org/x/xerrors"

	"github.com/int128/ghcp/pkg/git"
	"github.com/int128/ghcp/pkg/git/commitstrategy"
	"github.com/int128/ghcp/pkg/usecases/commit"
)

const emptyCommitCmdExample = `  To create an empty commit to the default branch:
    ghcp empty-commit -u OWNER -r REPO -m MESSAGE

  To create an empty commit to the branch:
    ghcp empty-commit -u OWNER -r REPO -b BRANCH -m MESSAGE

  If the branch does not exist, ghcp creates a branch from the default branch.
  It the branch exists, ghcp updates the branch by fast-forward.

  To create an empty commit to a new branch from the parent branch:
    ghcp empty-commit -u OWNER -r REPO -b BRANCH --parent PARENT -m MESSAGE

  If the branch exists, it will fail.`

func (r *Runner) newEmptyCommitCmd(ctx context.Context, gOpts *globalOptions) *cobra.Command {
	var o emptyCommitOptions
	c := &cobra.Command{
		Use:     fmt.Sprintf("%s [flags]", emptyCommitCmdName),
		Short:   "Create an empty commit to the branch",
		Long:    `This creates an empty commit to the branch. This will create a branch if it does not exist.`,
		Example: emptyCommitCmdExample,
		Args: func(_ *cobra.Command, args []string) error {
			if err := o.validate(); err != nil {
				return xerrors.Errorf("invalid flag: %w", err)
			}
			if len(args) > 0 {
				return xerrors.New("do not set any argument")
			}
			return nil
		},
		RunE: func(*cobra.Command, []string) error {
			ir, err := r.newInternalRunner(gOpts)
			if err != nil {
				return xerrors.Errorf("error while bootstrap of the dependencies: %w", err)
			}
			in := commit.Input{
				TargetRepository: git.RepositoryID{
					Owner: o.RepositoryOwner,
					Name:  o.RepositoryName,
				},
				TargetBranchName: git.BranchName(o.BranchName),
				ParentRepository: git.RepositoryID{
					Owner: o.RepositoryOwner,
					Name:  o.RepositoryName,
				},
				CommitStrategy: o.commitStrategy(),
				CommitMessage:  git.CommitMessage(o.CommitMessage),
				Author:         o.author(),
				Committer:      o.committer(),
				DryRun:         o.DryRun,
			}
			if err := ir.CommitUseCase.Do(ctx, in); err != nil {
				ir.Logger.Debugf("Stacktrace:\n%+v", err)
				return xerrors.Errorf("could not create an empty commit: %s", err)
			}
			return nil
		},
	}
	o.register(c.Flags())
	return c
}

type emptyCommitOptions struct {
	commitAttributeOptions

	RepositoryOwner string
	RepositoryName  string
	BranchName      string
	ParentRef       string
	DryRun          bool
}

func (o *emptyCommitOptions) validate() error {
	if err := o.commitAttributeOptions.validate(); err != nil {
		return xerrors.Errorf("%w", err)
	}
	return nil
}

func (o *emptyCommitOptions) commitStrategy() commitstrategy.CommitStrategy {
	if o.ParentRef != "" {
		return commitstrategy.RebaseOn(git.RefName(o.ParentRef))
	}
	return commitstrategy.FastForward
}

func (o *emptyCommitOptions) register(f *pflag.FlagSet) {
	f.StringVarP(&o.RepositoryOwner, "owner", "u", "", "GitHub repository owner (mandatory)")
	f.StringVarP(&o.RepositoryName, "repo", "r", "", "GitHub repository name (mandatory)")
	f.StringVarP(&o.BranchName, "branch", "b", "", "Name of the branch to create or update (default: the default branch of repository)")
	f.StringVar(&o.ParentRef, "parent", "", "Create a commit from the parent branch/tag (default: fast-forward)")
	f.BoolVar(&o.DryRun, "dry-run", false, "Do not update the branch actually")
	o.commitAttributeOptions.register(f)
}
