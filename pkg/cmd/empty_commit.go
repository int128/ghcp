package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/int128/ghcp/pkg/git"
	"github.com/int128/ghcp/pkg/git/commitstrategy"
	"github.com/int128/ghcp/pkg/usecases/commit"
)

const emptyCommitCmdExample = `  To create an empty commit to the default branch:
    ghcp empty-commit -r OWNER/REPO -m MESSAGE

  To create an empty commit to the branch:
    ghcp empty-commit -r OWNER/REPO -b BRANCH -m MESSAGE

  If the branch does not exist, ghcp creates a branch from the default branch.
  It the branch exists, ghcp updates the branch by fast-forward.

  To create an empty commit to a new branch from the parent branch:
    ghcp empty-commit -r OWNER/REPO -b BRANCH --parent PARENT -m MESSAGE

  If the branch exists, it will fail.`

func (r *Runner) newEmptyCommitCmd(ctx context.Context, gOpts *globalOptions) *cobra.Command {
	var o emptyCommitOptions
	c := &cobra.Command{
		Use:     fmt.Sprintf("%s [flags]", emptyCommitCmdName),
		Short:   "Create an empty commit to the branch",
		Long:    `This creates an empty commit to the branch. This will create a branch if it does not exist.`,
		Example: emptyCommitCmdExample,
		Args:    cobra.NoArgs,
		RunE: func(*cobra.Command, []string) error {
			if err := o.validate(); err != nil {
				return fmt.Errorf("invalid flag: %w", err)
			}
			targetRepository, err := o.repositoryID()
			if err != nil {
				return fmt.Errorf("invalid flag: %w", err)
			}

			ir, err := r.newInternalRunner(gOpts)
			if err != nil {
				return fmt.Errorf("error while bootstrap of the dependencies: %w", err)
			}
			in := commit.Input{
				TargetRepository: targetRepository,
				TargetBranchName: git.BranchName(o.BranchName),
				ParentRepository: targetRepository,
				CommitStrategy:   o.commitStrategy(),
				CommitMessage:    git.CommitMessage(o.CommitMessage),
				Author:           o.author(),
				Committer:        o.committer(),
				DryRun:           o.DryRun,
			}
			if err := ir.CommitUseCase.Do(ctx, in); err != nil {
				ir.Logger.Debugf("Stacktrace:\n%+v", err)
				return fmt.Errorf("could not create an empty commit: %s", err)
			}
			return nil
		},
	}
	o.register(c.Flags())
	return c
}

type emptyCommitOptions struct {
	commitAttributeOptions
	repositoryOptions

	BranchName string
	ParentRef  string
	DryRun     bool
}

func (o emptyCommitOptions) validate() error {
	if err := o.commitAttributeOptions.validate(); err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}

func (o emptyCommitOptions) commitStrategy() commitstrategy.CommitStrategy {
	if o.ParentRef != "" {
		return commitstrategy.RebaseOn(git.RefName(o.ParentRef))
	}
	return commitstrategy.FastForward
}

func (o *emptyCommitOptions) register(f *pflag.FlagSet) {
	o.repositoryOptions.register(f)
	f.StringVarP(&o.BranchName, "branch", "b", "", "Name of the branch to create or update (default: the default branch of repository)")
	f.StringVar(&o.ParentRef, "parent", "", "Create a commit from the parent branch/tag (default: fast-forward)")
	f.BoolVar(&o.DryRun, "dry-run", false, "Do not update the branch actually")
	o.commitAttributeOptions.register(f)
}
