package usecases

import "context"

//go:generate mockgen -package mock_usecases -destination ../mock_usecases/mock_usecases.go github.com/int128/ghcp/usecases/interfaces Push

type Push interface {
	Do(ctx context.Context, in PushIn) error
}

type PushIn struct {
	RepositoryOwner string
	RepositoryName  string
	Paths           []string
}
