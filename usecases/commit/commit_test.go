package commit

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/int128/ghcp/adaptors"
	"github.com/int128/ghcp/adaptors/mock_adaptors"
	"github.com/int128/ghcp/git"
	"github.com/int128/ghcp/usecases"
	"github.com/int128/ghcp/usecases/mock_usecases"
)

func TestCommitToBranch_Do(t *testing.T) {
	ctx := context.TODO()
	parentRepositoryID := git.RepositoryID{Owner: "upstream", Name: "repo"}
	targetRepositoryID := git.RepositoryID{Owner: "owner", Name: "repo"}
	thePathFilter := gomock.AssignableToTypeOf(&pathFilter{})
	theFiles := []adaptors.File{
		{Path: "file1"},
		{Path: "file2", Executable: true},
	}

	type testCase struct {
		dryRun               bool
		noFileMode           bool
		changedFiles         int
		branchOperationTimes int
	}

	run := func(t *testing.T, c testCase) {
		t.Run("FastForward", func(t *testing.T) {
			t.Run("CreateBranchIfItDoesNotExist", func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()
				in := usecases.CommitIn{
					ParentRepository: parentRepositoryID,
					ParentBranch:     usecases.ParentBranch{FastForward: true},
					TargetRepository: targetRepositoryID,
					TargetBranchName: "topic",
					CommitMessage:    "message",
					Paths:            []string{"path"},
					NoFileMode:       c.noFileMode,
					DryRun:           c.dryRun,
				}

				fileSystem := mock_adaptors.NewMockFileSystem(ctrl)
				fileSystem.EXPECT().FindFiles([]string{"path"}, thePathFilter).Return(theFiles, nil)

				createBlobTreeCommit := mock_usecases.NewMockCreateBlobTreeCommit(ctrl)
				createBlobTreeCommit.EXPECT().
					Do(ctx, usecases.CreateBlobTreeCommitIn{
						Files:           theFiles,
						Repository:      targetRepositoryID,
						CommitMessage:   "message",
						ParentCommitSHA: "masterCommitSHA",
						ParentTreeSHA:   "masterTreeSHA",
						NoFileMode:      c.noFileMode,
					}).
					Return(&usecases.CreateBlobTreeCommitOut{
						CommitSHA:    "commitSHA",
						ChangedFiles: c.changedFiles,
					}, nil)

				gitHub := mock_adaptors.NewMockGitHub(ctrl)
				gitHub.EXPECT().
					QueryForCommitToBranch(ctx, adaptors.QueryForCommitToBranchIn{
						ParentRepository: parentRepositoryID,
						TargetRepository: targetRepositoryID,
						TargetBranchName: "topic",
					}).
					Return(&adaptors.QueryForCommitToBranchOut{
						CurrentUserName:              "current",
						TargetRepository:             targetRepositoryID,
						TargetDefaultBranchName:      "master",
						ParentDefaultBranchCommitSHA: "masterCommitSHA",
						ParentDefaultBranchTreeSHA:   "masterTreeSHA",
					}, nil)
				gitHub.EXPECT().
					CreateBranch(ctx, git.NewBranch{
						Repository: targetRepositoryID,
						BranchName: "topic",
						CommitSHA:  "commitSHA",
					}).
					Return(nil).
					Times(c.branchOperationTimes)

				useCase := Commit{
					CreateBlobTreeCommit: createBlobTreeCommit,
					FileSystem:           fileSystem,
					Logger:               mock_adaptors.NewLogger(t),
					GitHub:               gitHub,
				}
				if err := useCase.Do(ctx, in); err != nil {
					t.Errorf("err wants nil but %+v", err)
				}
			})

			t.Run("UpdateBranch", func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()
				in := usecases.CommitIn{
					ParentRepository: parentRepositoryID,
					ParentBranch:     usecases.ParentBranch{FastForward: true},
					TargetRepository: targetRepositoryID,
					TargetBranchName: "topic",
					CommitMessage:    "message",
					Paths:            []string{"path"},
					NoFileMode:       c.noFileMode,
					DryRun:           c.dryRun,
				}

				fileSystem := mock_adaptors.NewMockFileSystem(ctrl)
				fileSystem.EXPECT().FindFiles([]string{"path"}, thePathFilter).Return(theFiles, nil)

				createBlobTreeCommit := mock_usecases.NewMockCreateBlobTreeCommit(ctrl)
				createBlobTreeCommit.EXPECT().
					Do(ctx, usecases.CreateBlobTreeCommitIn{
						Files:           theFiles,
						Repository:      targetRepositoryID,
						CommitMessage:   "message",
						ParentCommitSHA: "topicCommitSHA",
						ParentTreeSHA:   "topicTreeSHA",
						NoFileMode:      c.noFileMode,
					}).
					Return(&usecases.CreateBlobTreeCommitOut{
						CommitSHA:    "commitSHA",
						ChangedFiles: c.changedFiles,
					}, nil)

				gitHub := mock_adaptors.NewMockGitHub(ctrl)
				gitHub.EXPECT().
					QueryForCommitToBranch(ctx, adaptors.QueryForCommitToBranchIn{
						ParentRepository: parentRepositoryID,
						TargetRepository: targetRepositoryID,
						TargetBranchName: "topic",
					}).
					Return(&adaptors.QueryForCommitToBranchOut{
						CurrentUserName:              "current",
						TargetRepository:             targetRepositoryID,
						TargetDefaultBranchName:      "master",
						ParentDefaultBranchCommitSHA: "masterCommitSHA",
						ParentDefaultBranchTreeSHA:   "masterTreeSHA",
						TargetBranchCommitSHA:        "topicCommitSHA",
						TargetBranchTreeSHA:          "topicTreeSHA",
					}, nil)
				gitHub.EXPECT().
					UpdateBranch(ctx, git.NewBranch{
						Repository: targetRepositoryID,
						BranchName: "topic",
						CommitSHA:  "commitSHA",
					}, false).
					Return(nil).
					Times(c.branchOperationTimes)

				useCase := Commit{
					CreateBlobTreeCommit: createBlobTreeCommit,
					FileSystem:           fileSystem,
					Logger:               mock_adaptors.NewLogger(t),
					GitHub:               gitHub,
				}
				if err := useCase.Do(ctx, in); err != nil {
					t.Errorf("err wants nil but %+v", err)
				}
			})

			t.Run("UpdateDefaultBranch", func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()
				in := usecases.CommitIn{
					ParentRepository: parentRepositoryID,
					ParentBranch:     usecases.ParentBranch{FastForward: true},
					TargetRepository: targetRepositoryID,
					CommitMessage:    "message",
					Paths:            []string{"path"},
					NoFileMode:       c.noFileMode,
					DryRun:           c.dryRun,
				}

				fileSystem := mock_adaptors.NewMockFileSystem(ctrl)
				fileSystem.EXPECT().FindFiles([]string{"path"}, thePathFilter).Return(theFiles, nil)

				createBlobTreeCommit := mock_usecases.NewMockCreateBlobTreeCommit(ctrl)
				createBlobTreeCommit.EXPECT().
					Do(ctx, usecases.CreateBlobTreeCommitIn{
						Files:           theFiles,
						Repository:      targetRepositoryID,
						CommitMessage:   "message",
						ParentCommitSHA: "masterCommitSHA",
						ParentTreeSHA:   "masterTreeSHA",
						NoFileMode:      c.noFileMode,
					}).
					Return(&usecases.CreateBlobTreeCommitOut{
						CommitSHA:    "commitSHA",
						ChangedFiles: c.changedFiles,
					}, nil)

				gitHub := mock_adaptors.NewMockGitHub(ctrl)
				gitHub.EXPECT().
					QueryForCommitToBranch(ctx, adaptors.QueryForCommitToBranchIn{
						ParentRepository: parentRepositoryID,
						TargetRepository: targetRepositoryID,
					}).
					Return(&adaptors.QueryForCommitToBranchOut{
						CurrentUserName:              "current",
						TargetRepository:             targetRepositoryID,
						TargetDefaultBranchName:      "master",
						ParentDefaultBranchCommitSHA: "masterCommitSHA",
						ParentDefaultBranchTreeSHA:   "masterTreeSHA",
					}, nil)
				gitHub.EXPECT().
					UpdateBranch(ctx, git.NewBranch{
						Repository: targetRepositoryID,
						BranchName: "master",
						CommitSHA:  "commitSHA",
					}, false).
					Return(nil).
					Times(c.branchOperationTimes)

				useCase := Commit{
					CreateBlobTreeCommit: createBlobTreeCommit,
					FileSystem:           fileSystem,
					Logger:               mock_adaptors.NewLogger(t),
					GitHub:               gitHub,
				}
				if err := useCase.Do(ctx, in); err != nil {
					t.Errorf("err wants nil but %+v", err)
				}
			})
		})

		t.Run("NoParent", func(t *testing.T) {
			t.Run("CreateBranchIfItDoesNotExist", func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()
				in := usecases.CommitIn{
					ParentRepository: parentRepositoryID,
					ParentBranch:     usecases.ParentBranch{NoParent: true},
					TargetRepository: targetRepositoryID,
					TargetBranchName: "topic",
					CommitMessage:    "message",
					Paths:            []string{"path"},
					NoFileMode:       c.noFileMode,
					DryRun:           c.dryRun,
				}

				fileSystem := mock_adaptors.NewMockFileSystem(ctrl)
				fileSystem.EXPECT().FindFiles([]string{"path"}, thePathFilter).Return(theFiles, nil)

				createBlobTreeCommit := mock_usecases.NewMockCreateBlobTreeCommit(ctrl)
				createBlobTreeCommit.EXPECT().
					Do(ctx, usecases.CreateBlobTreeCommitIn{
						Files:         theFiles,
						Repository:    targetRepositoryID,
						CommitMessage: "message",
						NoFileMode:    c.noFileMode,
					}).
					Return(&usecases.CreateBlobTreeCommitOut{
						CommitSHA:    "commitSHA",
						ChangedFiles: c.changedFiles,
					}, nil)

				gitHub := mock_adaptors.NewMockGitHub(ctrl)
				gitHub.EXPECT().
					QueryForCommitToBranch(ctx, adaptors.QueryForCommitToBranchIn{
						ParentRepository: parentRepositoryID,
						TargetRepository: targetRepositoryID,
						TargetBranchName: "topic",
					}).
					Return(&adaptors.QueryForCommitToBranchOut{
						CurrentUserName:              "current",
						TargetRepository:             targetRepositoryID,
						TargetDefaultBranchName:      "master",
						ParentDefaultBranchCommitSHA: "masterCommitSHA",
						ParentDefaultBranchTreeSHA:   "masterTreeSHA",
					}, nil)
				gitHub.EXPECT().
					CreateBranch(ctx, git.NewBranch{
						Repository: targetRepositoryID,
						BranchName: "topic",
						CommitSHA:  "commitSHA",
					}).
					Return(nil).
					Times(c.branchOperationTimes)

				useCase := Commit{
					CreateBlobTreeCommit: createBlobTreeCommit,
					FileSystem:           fileSystem,
					Logger:               mock_adaptors.NewLogger(t),
					GitHub:               gitHub,
				}
				if err := useCase.Do(ctx, in); err != nil {
					t.Errorf("err wants nil but %+v", err)
				}
			})

			t.Run("UpdateBranch", func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()
				in := usecases.CommitIn{
					ParentRepository: parentRepositoryID,
					ParentBranch:     usecases.ParentBranch{FastForward: true},
					TargetRepository: targetRepositoryID,
					TargetBranchName: "topic",
					CommitMessage:    "message",
					Paths:            []string{"path"},
					NoFileMode:       c.noFileMode,
					DryRun:           c.dryRun,
				}

				fileSystem := mock_adaptors.NewMockFileSystem(ctrl)
				fileSystem.EXPECT().FindFiles([]string{"path"}, thePathFilter).Return(theFiles, nil)

				createBlobTreeCommit := mock_usecases.NewMockCreateBlobTreeCommit(ctrl)
				createBlobTreeCommit.EXPECT().
					Do(ctx, usecases.CreateBlobTreeCommitIn{
						Files:           theFiles,
						Repository:      targetRepositoryID,
						CommitMessage:   "message",
						ParentCommitSHA: "topicCommitSHA",
						ParentTreeSHA:   "topicTreeSHA",
						NoFileMode:      c.noFileMode,
					}).
					Return(&usecases.CreateBlobTreeCommitOut{
						CommitSHA:    "commitSHA",
						ChangedFiles: c.changedFiles,
					}, nil)

				gitHub := mock_adaptors.NewMockGitHub(ctrl)
				gitHub.EXPECT().
					QueryForCommitToBranch(ctx, adaptors.QueryForCommitToBranchIn{
						ParentRepository: parentRepositoryID,
						TargetRepository: targetRepositoryID,
						TargetBranchName: "topic",
					}).
					Return(&adaptors.QueryForCommitToBranchOut{
						CurrentUserName:              "current",
						TargetRepository:             targetRepositoryID,
						TargetDefaultBranchName:      "master",
						ParentDefaultBranchCommitSHA: "masterCommitSHA",
						ParentDefaultBranchTreeSHA:   "masterTreeSHA",
						TargetBranchCommitSHA:        "topicCommitSHA",
						TargetBranchTreeSHA:          "topicTreeSHA",
					}, nil)
				gitHub.EXPECT().
					UpdateBranch(ctx, git.NewBranch{
						Repository: targetRepositoryID,
						BranchName: "topic",
						CommitSHA:  "commitSHA",
					}, false). //TODO: force update
					Return(nil).
					Times(c.branchOperationTimes)

				useCase := Commit{
					CreateBlobTreeCommit: createBlobTreeCommit,
					FileSystem:           fileSystem,
					Logger:               mock_adaptors.NewLogger(t),
					GitHub:               gitHub,
				}
				if err := useCase.Do(ctx, in); err != nil {
					t.Errorf("err wants nil but %+v", err)
				}
			})

			t.Run("UpdateDefaultBranch", func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()
				in := usecases.CommitIn{
					ParentRepository: parentRepositoryID,
					ParentBranch:     usecases.ParentBranch{NoParent: true},
					TargetRepository: targetRepositoryID,
					TargetBranchName: "topic",
					CommitMessage:    "message",
					Paths:            []string{"path"},
					NoFileMode:       c.noFileMode,
					DryRun:           c.dryRun,
				}

				fileSystem := mock_adaptors.NewMockFileSystem(ctrl)
				fileSystem.EXPECT().FindFiles([]string{"path"}, thePathFilter).Return(theFiles, nil)

				createBlobTreeCommit := mock_usecases.NewMockCreateBlobTreeCommit(ctrl)
				createBlobTreeCommit.EXPECT().
					Do(ctx, usecases.CreateBlobTreeCommitIn{
						Files:         theFiles,
						Repository:    targetRepositoryID,
						CommitMessage: "message",
						NoFileMode:    c.noFileMode,
					}).
					Return(&usecases.CreateBlobTreeCommitOut{
						CommitSHA:    "commitSHA",
						ChangedFiles: c.changedFiles,
					}, nil)

				gitHub := mock_adaptors.NewMockGitHub(ctrl)
				gitHub.EXPECT().
					QueryForCommitToBranch(ctx, adaptors.QueryForCommitToBranchIn{
						ParentRepository: parentRepositoryID,
						TargetRepository: targetRepositoryID,
						TargetBranchName: "topic",
					}).
					Return(&adaptors.QueryForCommitToBranchOut{
						CurrentUserName:              "current",
						TargetRepository:             targetRepositoryID,
						TargetDefaultBranchName:      "master",
						ParentDefaultBranchCommitSHA: "masterCommitSHA",
						ParentDefaultBranchTreeSHA:   "masterTreeSHA",
						TargetBranchCommitSHA:        "topicCommitSHA",
						TargetBranchTreeSHA:          "topicTreeSHA",
					}, nil)
				gitHub.EXPECT().
					UpdateBranch(ctx, git.NewBranch{
						Repository: targetRepositoryID,
						BranchName: "topic",
						CommitSHA:  "commitSHA",
					}, false). //TODO: force update
					Return(nil).
					Times(c.branchOperationTimes)

				useCase := Commit{
					CreateBlobTreeCommit: createBlobTreeCommit,
					FileSystem:           fileSystem,
					Logger:               mock_adaptors.NewLogger(t),
					GitHub:               gitHub,
				}
				if err := useCase.Do(ctx, in); err != nil {
					t.Errorf("err wants nil but %+v", err)
				}
			})
		})

		t.Run("FromRef", func(t *testing.T) {
			t.Run("CreateBranchIfItDoesNotExist", func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()
				in := usecases.CommitIn{
					ParentRepository: parentRepositoryID,
					ParentBranch:     usecases.ParentBranch{FromRef: "develop"},
					TargetRepository: targetRepositoryID,
					TargetBranchName: "topic",
					CommitMessage:    "message",
					Paths:            []string{"path"},
					NoFileMode:       c.noFileMode,
					DryRun:           c.dryRun,
				}

				fileSystem := mock_adaptors.NewMockFileSystem(ctrl)
				fileSystem.EXPECT().FindFiles([]string{"path"}, thePathFilter).Return(theFiles, nil)

				createBlobTreeCommit := mock_usecases.NewMockCreateBlobTreeCommit(ctrl)
				createBlobTreeCommit.EXPECT().
					Do(ctx, usecases.CreateBlobTreeCommitIn{
						Files:           theFiles,
						Repository:      targetRepositoryID,
						CommitMessage:   "message",
						NoFileMode:      c.noFileMode,
						ParentCommitSHA: "developCommitSHA",
						ParentTreeSHA:   "developTreeSHA",
					}).
					Return(&usecases.CreateBlobTreeCommitOut{
						CommitSHA:    "commitSHA",
						ChangedFiles: c.changedFiles,
					}, nil)

				gitHub := mock_adaptors.NewMockGitHub(ctrl)
				gitHub.EXPECT().
					QueryForCommitToBranch(ctx, adaptors.QueryForCommitToBranchIn{
						ParentRepository: parentRepositoryID,
						ParentRef:        "develop",
						TargetRepository: targetRepositoryID,
						TargetBranchName: "topic",
					}).
					Return(&adaptors.QueryForCommitToBranchOut{
						CurrentUserName:              "current",
						TargetRepository:             targetRepositoryID,
						TargetDefaultBranchName:      "master",
						ParentDefaultBranchCommitSHA: "masterCommitSHA",
						ParentDefaultBranchTreeSHA:   "masterTreeSHA",
						ParentRefCommitSHA:           "developCommitSHA",
						ParentRefTreeSHA:             "developTreeSHA",
					}, nil)
				gitHub.EXPECT().
					CreateBranch(ctx, git.NewBranch{
						Repository: targetRepositoryID,
						BranchName: "topic",
						CommitSHA:  "commitSHA",
					}).
					Return(nil).
					Times(c.branchOperationTimes)

				useCase := Commit{
					CreateBlobTreeCommit: createBlobTreeCommit,
					FileSystem:           fileSystem,
					Logger:               mock_adaptors.NewLogger(t),
					GitHub:               gitHub,
				}
				if err := useCase.Do(ctx, in); err != nil {
					t.Errorf("err wants nil but %+v", err)
				}
			})

			t.Run("UpdateBranch", func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()
				in := usecases.CommitIn{
					ParentRepository: parentRepositoryID,
					ParentBranch:     usecases.ParentBranch{FromRef: "develop"},
					TargetRepository: targetRepositoryID,
					TargetBranchName: "topic",
					CommitMessage:    "message",
					Paths:            []string{"path"},
					NoFileMode:       c.noFileMode,
					DryRun:           c.dryRun,
				}

				fileSystem := mock_adaptors.NewMockFileSystem(ctrl)
				fileSystem.EXPECT().FindFiles([]string{"path"}, thePathFilter).Return(theFiles, nil)

				createBlobTreeCommit := mock_usecases.NewMockCreateBlobTreeCommit(ctrl)
				createBlobTreeCommit.EXPECT().
					Do(ctx, usecases.CreateBlobTreeCommitIn{
						Files:           theFiles,
						Repository:      targetRepositoryID,
						CommitMessage:   "message",
						ParentCommitSHA: "developCommitSHA",
						ParentTreeSHA:   "developTreeSHA",
						NoFileMode:      c.noFileMode,
					}).
					Return(&usecases.CreateBlobTreeCommitOut{
						CommitSHA:    "commitSHA",
						ChangedFiles: c.changedFiles,
					}, nil)

				gitHub := mock_adaptors.NewMockGitHub(ctrl)
				gitHub.EXPECT().
					QueryForCommitToBranch(ctx, adaptors.QueryForCommitToBranchIn{
						ParentRepository: parentRepositoryID,
						ParentRef:        "develop",
						TargetRepository: targetRepositoryID,
						TargetBranchName: "topic",
					}).
					Return(&adaptors.QueryForCommitToBranchOut{
						CurrentUserName:              "current",
						TargetRepository:             targetRepositoryID,
						TargetDefaultBranchName:      "master",
						ParentDefaultBranchCommitSHA: "masterCommitSHA",
						ParentDefaultBranchTreeSHA:   "masterTreeSHA",
						TargetBranchCommitSHA:        "topicCommitSHA",
						TargetBranchTreeSHA:          "topicTreeSHA",
						ParentRefCommitSHA:           "developCommitSHA",
						ParentRefTreeSHA:             "developTreeSHA",
					}, nil)
				gitHub.EXPECT().
					UpdateBranch(ctx, git.NewBranch{
						Repository: targetRepositoryID,
						BranchName: "topic",
						CommitSHA:  "commitSHA",
					}, false). //TODO: force update
					Return(nil).
					Times(c.branchOperationTimes)

				useCase := Commit{
					CreateBlobTreeCommit: createBlobTreeCommit,
					FileSystem:           fileSystem,
					Logger:               mock_adaptors.NewLogger(t),
					GitHub:               gitHub,
				}
				if err := useCase.Do(ctx, in); err != nil {
					t.Errorf("err wants nil but %+v", err)
				}
			})

			t.Run("UpdateDefaultBranch", func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()
				in := usecases.CommitIn{
					ParentRepository: parentRepositoryID,
					ParentBranch:     usecases.ParentBranch{FromRef: "develop"},
					TargetRepository: targetRepositoryID,
					TargetBranchName: "topic",
					CommitMessage:    "message",
					Paths:            []string{"path"},
					NoFileMode:       c.noFileMode,
					DryRun:           c.dryRun,
				}

				fileSystem := mock_adaptors.NewMockFileSystem(ctrl)
				fileSystem.EXPECT().FindFiles([]string{"path"}, thePathFilter).Return(theFiles, nil)

				createBlobTreeCommit := mock_usecases.NewMockCreateBlobTreeCommit(ctrl)
				createBlobTreeCommit.EXPECT().
					Do(ctx, usecases.CreateBlobTreeCommitIn{
						Files:           theFiles,
						Repository:      targetRepositoryID,
						CommitMessage:   "message",
						ParentCommitSHA: "developCommitSHA",
						ParentTreeSHA:   "developTreeSHA",
						NoFileMode:      c.noFileMode,
					}).
					Return(&usecases.CreateBlobTreeCommitOut{
						CommitSHA:    "commitSHA",
						ChangedFiles: c.changedFiles,
					}, nil)

				gitHub := mock_adaptors.NewMockGitHub(ctrl)
				gitHub.EXPECT().
					QueryForCommitToBranch(ctx, adaptors.QueryForCommitToBranchIn{
						ParentRepository: parentRepositoryID,
						ParentRef:        "develop",
						TargetRepository: targetRepositoryID,
						TargetBranchName: "topic",
					}).
					Return(&adaptors.QueryForCommitToBranchOut{
						CurrentUserName:              "current",
						TargetRepository:             targetRepositoryID,
						TargetDefaultBranchName:      "master",
						ParentDefaultBranchCommitSHA: "masterCommitSHA",
						ParentDefaultBranchTreeSHA:   "masterTreeSHA",
						TargetBranchCommitSHA:        "topicCommitSHA",
						TargetBranchTreeSHA:          "topicTreeSHA",
						ParentRefCommitSHA:           "developCommitSHA",
						ParentRefTreeSHA:             "developTreeSHA",
					}, nil)
				gitHub.EXPECT().
					UpdateBranch(ctx, git.NewBranch{
						Repository: targetRepositoryID,
						BranchName: "topic",
						CommitSHA:  "commitSHA",
					}, false). //TODO: force update
					Return(nil).
					Times(c.branchOperationTimes)

				useCase := Commit{
					CreateBlobTreeCommit: createBlobTreeCommit,
					FileSystem:           fileSystem,
					Logger:               mock_adaptors.NewLogger(t),
					GitHub:               gitHub,
				}
				if err := useCase.Do(ctx, in); err != nil {
					t.Errorf("err wants nil but %+v", err)
				}
			})
		})
	}

	for name, c := range map[string]testCase{
		"Success": {
			changedFiles:         1,
			branchOperationTimes: 1,
		},
		"NothingToCommit": {
			changedFiles:         0,
			branchOperationTimes: 0,
		},
		"DryRun": {
			dryRun:               true,
			branchOperationTimes: 0,
		},
		"NoFileMode": {
			noFileMode: true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			run(t, c)
		})
	}
}

func Test_pathFilter_SkipDir(t *testing.T) {
	for _, c := range []struct {
		path string
		skip bool
	}{
		{path: "."},
		{path: "foo"},
		{path: ".git", skip: true},
		{path: "foo/bar"},
		{path: "foo/.git", skip: true},
	} {
		t.Run(c.path, func(t *testing.T) {
			f := &pathFilter{Logger: mock_adaptors.NewLogger(t)}
			skip := f.SkipDir(c.path)
			if skip != c.skip {
				t.Errorf("skip wants %v but %v", c.skip, skip)
			}
		})
	}
}

func Test_pathFilter_ExcludeFile(t *testing.T) {
	f := &pathFilter{Logger: mock_adaptors.NewLogger(t)}
	exclude := f.ExcludeFile("foo")
	if exclude {
		t.Errorf("exclude wants %v but %v", false, exclude)
	}
}
