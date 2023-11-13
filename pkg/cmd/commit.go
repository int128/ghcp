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

const commitCmdExample = `  To commit files to the default branch:
    ghcp commit -r OWNER/REPO -m MESSAGE FILES...

  To commit files to the branch:
    ghcp commit -r OWNER/REPO -b BRANCH -m MESSAGE FILES...

  If the branch does not exist, ghcp creates a branch from the default branch.
  It the branch exists, ghcp updates the branch by fast-forward.

  To commit files to a new branch from the parent branch:
    ghcp commit -r OWNER/REPO -b BRANCH --parent PARENT -m MESSAGE FILES...

  If the branch exists, it will fail.

  To commit files to a new branch without any parent:
    ghcp commit -r OWNER/REPO -b BRANCH --no-parent -m MESSAGE FILES...

  If the branch exists, it will fail.`

func (r *Runner) newCommitCmd(ctx context.Context, gOpts *globalOptions) *cobra.Command {
	var o commitOptions
	c := &cobra.Command{
		Use:     fmt.Sprintf("%s [flags] FILES...", commitCmdName),
		Short:   "Commit files to the branch",
		Long:    `This commits the files to the branch. This will create a branch if it does not exist.`,
		Example: commitCmdExample,
		RunE: func(_ *cobra.Command, args []string) error {
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
				Paths:            args,
				DeletedPaths:     o.DeletedPaths,
				NoFileMode:       o.NoFileMode,
				DryRun:           o.DryRun,
			}
			if err := ir.CommitUseCase.Do(ctx, in); err != nil {
				ir.Logger.Debugf("Stacktrace:\n%+v", err)
				return fmt.Errorf("could not commit the files: %s", err)
			}
			return nil
		},
	}
	o.register(c.Flags())
	return c
}

type commitOptions struct {
	commitAttributeOptions
	repositoryOptions

	BranchName   string
	ParentRef    string
	NoParent     bool
	NoFileMode   bool
	DryRun       bool
	DeletedPaths []string
}

func (o commitOptions) validate() error {
	if o.ParentRef != "" && o.NoParent {
		return fmt.Errorf("do not set both --parent and --no-parent")
	}
	if err := o.commitAttributeOptions.validate(); err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}

func (o commitOptions) commitStrategy() commitstrategy.CommitStrategy {
	if o.NoParent {
		return commitstrategy.NoParent
	}
	if o.ParentRef != "" {
		return commitstrategy.RebaseOn(git.RefName(o.ParentRef))
	}
	return commitstrategy.FastForward
}

func (o *commitOptions) register(f *pflag.FlagSet) {
	o.repositoryOptions.register(f)
	f.StringVarP(&o.BranchName, "branch", "b", "", "Name of the branch to create or update (default: the default branch of repository)")
	f.StringVar(&o.ParentRef, "parent", "", "Create a commit from the parent branch/tag (default: fast-forward)")
	f.BoolVar(&o.NoParent, "no-parent", false, "Create a commit without a parent")
	f.BoolVar(&o.NoFileMode, "no-file-mode", false, "Ignore executable bit of file and treat as 0644")
	f.BoolVar(&o.DryRun, "dry-run", false, "Upload files but do not update the branch actually")
	f.StringSliceVarP(&o.DeletedPaths, "delete", "d", nil, "the list of deleted file paths")
	o.commitAttributeOptions.register(f)
}
