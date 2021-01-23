package cmd

import (
	"context"
	"fmt"

	"github.com/int128/ghcp/pkg/git"
	"github.com/int128/ghcp/pkg/usecases/release"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"golang.org/x/xerrors"
)

const releaseCmdExample = `  To upload files to the release associated to tag TAG:
    ghcp release -u OWNER -r REPO -t TAG FILES...

  If the release does not exist, it will create a release.
  If the tag does not exist, it will create a tag from the default branch and create a release.

  To create a tag and release on commit COMMIT_SHA and upload files to the release:
    ghcp release -u OWNER -r REPO -t TAG --tagret COMMIT_SHA FILES...

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
			ir, err := r.newInternalRunner(gOpts)
			if err != nil {
				return xerrors.Errorf("error while bootstrap of the dependencies: %w", err)
			}
			in := release.Input{
				Repository: git.RepositoryID{
					Owner: o.RepositoryOwner,
					Name:  o.RepositoryName,
				},
				TagName:                 git.TagName(o.TagName),
				TargetBranchOrCommitSHA: o.TargetBranchOrCommitSHA,
				Paths:                   args,
				DryRun:                  o.DryRun,
			}
			if err := ir.ReleaseUseCase.Do(ctx, in); err != nil {
				ir.Logger.Debugf("Stacktrace:\n%+v", err)
				return xerrors.Errorf("could not release the files: %s", err)
			}
			return nil
		},
	}
	o.register(c.Flags())
	return c
}

type releaseOptions struct {
	RepositoryOwner         string
	RepositoryName          string
	TagName                 string
	TargetBranchOrCommitSHA string
	DryRun                  bool
}

func (o *releaseOptions) register(f *pflag.FlagSet) {
	f.StringVarP(&o.RepositoryOwner, "owner", "u", "", "GitHub repository owner (mandatory)")
	f.StringVarP(&o.RepositoryName, "repo", "r", "", "GitHub repository name (mandatory)")
	f.StringVarP(&o.TagName, "tag", "t", "", "Tag name (mandatory)")
	f.StringVar(&o.TargetBranchOrCommitSHA, "target", "", "Branch name or commit SHA of a tag. Unused if the Git tag already exists (default: the default branch)")
	f.BoolVar(&o.DryRun, "dry-run", false, "Do not create a release and assets actually")
}
