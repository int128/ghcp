package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/int128/ghcp/domain/git"
	"github.com/int128/ghcp/domain/git/commitstrategy"
	"github.com/int128/ghcp/usecases/forkcommit"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"golang.org/x/xerrors"
)

func (r *Runner) newForkCommitCmd(ctx context.Context, gOpts *globalOptions) *cobra.Command {
	var o forkCommitOptions
	c := &cobra.Command{
		Use:   fmt.Sprintf("%s [flags] FILES...", forkCommitCmdName),
		Short: "Fork the repository and commit files to a branch",
		Long:  `This forks the repository and commits the files to a new branch.`,
		Args: func(*cobra.Command, []string) error {
			var errs []string
			if o.UpstreamRepositoryOwner == "" {
				errs = append(errs, "--owner is missing")
			}
			if o.UpstreamRepositoryName == "" {
				errs = append(errs, "--repo is missing")
			}
			if o.TargetBranchName == "" {
				errs = append(errs, "--branch is missing")
			}
			if len(errs) > 0 {
				return xerrors.New(strings.Join(errs, ", "))
			}
			return nil
		},
		RunE: func(_ *cobra.Command, args []string) error {
			ir, err := r.newInternalRunner(gOpts)
			if err != nil {
				return xerrors.Errorf("error while bootstrap of the dependencies: %w", err)
			}
			in := forkcommit.Input{
				ParentRepository: git.RepositoryID{
					Owner: o.UpstreamRepositoryOwner,
					Name:  o.UpstreamRepositoryName,
				},
				TargetBranchName: git.BranchName(o.TargetBranchName),
				CommitStrategy:   o.commitStrategy(),
				CommitMessage:    git.CommitMessage(o.CommitMessage),
				Paths:            args,
				NoFileMode:       o.NoFileMode,
				DryRun:           o.DryRun,
			}
			if err := ir.ForkCommitUseCase.Do(ctx, in); err != nil {
				ir.Logger.Debugf("Stacktrace:\n%+v", err)
				return xerrors.Errorf("could not commit the files: %s", err)
			}
			return nil
		},
	}
	o.register(c.Flags())
	return c
}

type forkCommitOptions struct {
	UpstreamRepositoryOwner string
	UpstreamRepositoryName  string
	UpstreamBranchName      string
	TargetBranchName        string
	CommitMessage           string
	NoFileMode              bool
	DryRun                  bool
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
	f.StringVarP(&o.CommitMessage, "message", "m", "", "Commit message (mandatory)")
	f.BoolVar(&o.NoFileMode, "no-file-mode", false, "Ignore executable bit of file and treat as 0644")
	f.BoolVar(&o.DryRun, "dry-run", false, "Upload files but do not update the branch actually")
}
