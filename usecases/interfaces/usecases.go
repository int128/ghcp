package usecases

import (
	"context"

	"github.com/int128/ghcp/git"
)

//go:generate mockgen -package mock_usecases -destination ../mock_usecases/mock_usecases.go github.com/int128/ghcp/usecases/interfaces CopyUseCase

type CopyUseCase interface {
	Do(ctx context.Context, in CopyUseCaseIn) error
}

type CopyUseCaseIn struct {
	Repository    git.RepositoryID
	CommitMessage git.CommitMessage
	Paths         []string
	DryRun        bool
}
