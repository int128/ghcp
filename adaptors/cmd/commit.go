package cmd

import (
	"context"
	"fmt"

	"github.com/int128/ghcp/domain/git"
	"github.com/int128/ghcp/usecases/commit"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"golang.org/x/xerrors"
)

const commitCmdExample = `  To commit files to the default branch:
    ghcp commit -u OWNER -r REPO -m MESSAGE FILES...

  To commit files to the branch:
    ghcp commit -u OWNER -r REPO -b BRANCH -m MESSAGE FILES...

  If the branch does not exist, ghcp creates a branch from the default branch.
  It the branch exists, ghcp updates the branch by fast-forward.

  To commit files to a new branch from the parent branch:
    ghcp commit -u OWNER -r REPO -b BRANCH --parent PARENT -m MESSAGE FILES...

  If the branch exists, ghcp cannot update the branch by fast-forward and will fail.

  To commit files to a new branch without any parent:
    ghcp commit -u OWNER -r REPO -b BRANCH --no-parent -m MESSAGE FILES...

  If the branch exists, ghcp cannot update the branch by fast-forward and will fail.`

func (r *Runner) newCommitCmd(ctx context.Context, gOpts *globalOptions) *cobra.Command {
	var o commitOptions
	c := &cobra.Command{
		Use:     fmt.Sprintf("%s [flags] FILES...", commitCmdName),
		Short:   "Commit files to the branch",
		Long:    `This commits the files to the branch. This will create a branch if it does not exist.`,
		Example: commitCmdExample,
		Args: func(*cobra.Command, []string) error {
			if o.ParentRef != "" && o.NoParent {
				return xerrors.Errorf("do not set both --parent and --no-parent")
			}
			return nil
		},
		RunE: func(_ *cobra.Command, args []string) error {
			ir, err := r.newInternalRunner(gOpts)
			if err != nil {
				return xerrors.Errorf("error while bootstrap of the dependencies: %w", err)
			}
			in := commit.Input{
				ParentRepository: git.RepositoryID{
					Owner: o.RepositoryOwner,
					Name:  o.RepositoryName,
				},
				ParentBranch: commit.ParentBranch{
					FastForward: o.ParentRef == "" && !o.NoParent,
					NoParent:    o.NoParent,
					FromRef:     git.RefName(o.ParentRef),
				},
				TargetRepository: git.RepositoryID{
					Owner: o.RepositoryOwner,
					Name:  o.RepositoryName,
				},
				TargetBranchName: git.BranchName(o.BranchName),
				CommitMessage:    git.CommitMessage(o.CommitMessage),
				Paths:            args,
				NoFileMode:       o.NoFileMode,
				DryRun:           o.DryRun,
			}
			if err := ir.CommitUseCase.Do(ctx, in); err != nil {
				ir.Logger.Debugf("Stacktrace:\n%+v", err)
				return xerrors.Errorf("could not commit the files: %s", err)
			}
			return nil
		},
	}
	o.register(c.Flags())
	return c
}

type commitOptions struct {
	RepositoryOwner string
	RepositoryName  string
	CommitMessage   string
	BranchName      string
	ParentRef       string
	NoParent        bool
	NoFileMode      bool
	DryRun          bool
}

func (o *commitOptions) register(f *pflag.FlagSet) {
	f.StringVarP(&o.RepositoryOwner, "owner", "u", "", "GitHub repository owner (mandatory)")
	f.StringVarP(&o.RepositoryName, "repo", "r", "", "GitHub repository name (mandatory)")
	f.StringVarP(&o.CommitMessage, "message", "m", "", "Commit message (mandatory)")
	f.StringVarP(&o.BranchName, "branch", "b", "", "Name of the branch to create or update (default: the default branch of repository)")
	f.StringVar(&o.ParentRef, "parent", "", "Create a commit from the parent branch/tag (default: fast-forward)")
	f.BoolVar(&o.NoParent, "no-parent", false, "Create a commit without a parent")
	f.BoolVar(&o.NoFileMode, "no-file-mode", false, "Ignore executable bit of file and treat as 0644")
	f.BoolVar(&o.DryRun, "dry-run", false, "Upload files but do not update the branch actually")
}
