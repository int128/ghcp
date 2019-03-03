// Package usecases provides use cases of this application.
package usecases

import (
	"context"
	"log"

	"github.com/int128/ghcp/adaptors/interfaces"
	"github.com/int128/ghcp/usecases/interfaces"
	"github.com/pkg/errors"
	"go.uber.org/dig"
)

func NewPush(i Push) usecases.Push {
	return &i
}

// Push performs commit and push files to the repository.
type Push struct {
	dig.In
	GitHub adaptors.GitHub
}

func (u *Push) Do(ctx context.Context, in usecases.PushIn) error {
	repo, err := u.GitHub.GetRepository(ctx, adaptors.GetRepositoryIn{
		Owner: in.RepositoryOwner,
		Name:  in.RepositoryName,
	})
	if err != nil {
		return errors.Wrapf(err, "error while getting the repository")
	}
	log.Printf("Logging in as %s", repo.CurrentUserName)

	//TODO: commit and push
	log.Printf("repo=%+v", repo)

	return nil
}
