package cmd

import (
	"context"
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/int128/ghcp/pkg/git"
	"github.com/int128/ghcp/pkg/git/commitstrategy"
	"github.com/int128/ghcp/pkg/usecases/forkcommit"
)

func (r *Runner) newForkCommitCmd(ctx context.Context, gOpts *globalOptions) *cobra.Command {
	var o forkCommitOptions
	c := &cobra.Command{
		Use:   fmt.Sprintf("%s [flags] FILES...", forkCommitCmdName),
		Short: "Fork the repository and commit files to a branch",
		Long:  `This forks the repository and commits the files to a new branch.`,
		Args: func(*cobra.Command, []string) error {
			if err := o.validate(); err != nil {
				return fmt.Errorf("invalid flag: %w", err)
			}
			return nil
		},
		RunE: func(_ *cobra.Command, args []string) error {
			ir, err := r.newInternalRunner(gOpts)
			if err != nil {
				return fmt.Errorf("error while bootstrap of the dependencies: %w", err)
			}
			in := forkcommit.Input{
				ParentRepository: git.RepositoryID{
					Owner: o.UpstreamRepositoryOwner,
					Name:  o.UpstreamRepositoryName,
				},
				TargetBranchName: git.BranchName(o.TargetBranchName),
				CommitStrategy:   o.commitStrategy(),
				CommitMessage:    git.CommitMessage(o.CommitMessage),
				Author:           o.author(),
				Committer:        o.committer(),
				Paths:            args,
				NoFileMode:       o.NoFileMode,
				DryRun:           o.DryRun,
			}
			if err := ir.ForkCommitUseCase.Do(ctx, in); err != nil {
				ir.Logger.Debugf("Stacktrace:\n%+v", err)
				return fmt.Errorf("could not commit the files: %s", err)
			}
			return nil
		},
	}
	o.register(c.Flags())
	return c
}

type forkCommitOptions struct {
	commitAttributeOptions

	UpstreamRepositoryOwner string
	UpstreamRepositoryName  string
	UpstreamBranchName      string
	TargetBranchName        string
	NoFileMode              bool
	DryRun                  bool
}

func (o *forkCommitOptions) validate() error {
	if o.UpstreamRepositoryOwner == "" {
		return errors.New("--owner is missing")
	}
	if o.UpstreamRepositoryName == "" {
		return errors.New("--repo is missing")
	}
	if o.TargetBranchName == "" {
		return errors.New("--branch is missing")
	}
	if err := o.commitAttributeOptions.validate(); err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}

func (o *forkCommitOptions) commitStrategy() commitstrategy.CommitStrategy {
	if o.UpstreamBranchName != "" {
		return commitstrategy.RebaseOn(git.RefName(o.UpstreamBranchName))
	}
	return commitstrategy.FastForward
}

func (o *forkCommitOptions) register(f *pflag.FlagSet) {
	f.StringVarP(&o.UpstreamRepositoryOwner, "owner", "u", "", "Upstream repository owner (mandatory)")
	f.StringVarP(&o.UpstreamRepositoryName, "repo", "r", "", "Upstream repository name (mandatory)")
	f.StringVar(&o.UpstreamBranchName, "parent", "", "Upstream branch name (default: the default branch of the upstream repository)")
	f.StringVarP(&o.TargetBranchName, "branch", "b", "", "Name of the branch to create (mandatory)")
	f.BoolVar(&o.NoFileMode, "no-file-mode", false, "Ignore executable bit of file and treat as 0644")
	f.BoolVar(&o.DryRun, "dry-run", false, "Upload files but do not update the branch actually")
	o.commitAttributeOptions.register(f)
}
