package cmd

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/int128/ghcp/pkg/git"
	"github.com/int128/ghcp/pkg/usecases/release"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const releaseCmdExample = `  To upload files to the release associated to tag TAG:
    ghcp release -r OWNER/REPO -t TAG FILES...

  If the release does not exist, it will create a release.
  If the tag does not exist, it will create a tag from the default branch and create a release.

  To create a tag and release on commit COMMIT_SHA and upload files to the release:
    ghcp release -r OWNER/REPO -t TAG --tagret COMMIT_SHA FILES...

  If the tag already exists, it ignores the target commit.
  If the release already exist, it only uploads the files.
`

func (r *Runner) newReleaseCmd(ctx context.Context, gOpts *globalOptions) *cobra.Command {
	var o releaseOptions
	c := &cobra.Command{
		Use:     fmt.Sprintf("%s [flags] FILES...", releaseCmdName),
		Short:   "Release files to the repository",
		Long:    `This uploads the files to the release associated to the tag. It will create a release if it does not exist.`,
		Example: releaseCmdExample,
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
			in := release.Input{
				Repository:              targetRepository,
				TagName:                 git.TagName(o.TagName),
				TargetBranchOrCommitSHA: o.TargetBranchOrCommitSHA,
				Paths:                   args,
				DryRun:                  o.DryRun,
			}
			if err := ir.ReleaseUseCase.Do(ctx, in); err != nil {
				slog.Debug("Stacktrace", "stacktrace", err)
				return fmt.Errorf("could not release the files: %s", err)
			}
			return nil
		},
	}
	o.register(c.Flags())
	return c
}

type releaseOptions struct {
	repositoryOptions

	TagName                 string
	TargetBranchOrCommitSHA string
	DryRun                  bool
}

func (o releaseOptions) validate() error {
	if o.TagName == "" {
		return errors.New("you need to set --tag")
	}
	return nil
}

func (o *releaseOptions) register(f *pflag.FlagSet) {
	o.repositoryOptions.register(f)
	f.StringVarP(&o.TagName, "tag", "t", "", "Tag name (mandatory)")
	f.StringVar(&o.TargetBranchOrCommitSHA, "target", "", "Branch name or commit SHA of a tag. Unused if the Git tag already exists (default: the default branch)")
	f.BoolVar(&o.DryRun, "dry-run", false, "Do not create a release and assets actually")
}
