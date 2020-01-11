package env

import (
	"os"

	"github.com/google/wire"
)

var Set = wire.NewSet(
	wire.Struct(new(Env)),
	wire.Bind(new(Interface), new(*Env)),
)

//go:generate mockgen -destination mock_env/mock_env.go github.com/int128/ghcp/adaptors/env Interface

type Interface interface {
	Getenv(key string) string
	Chdir(dir string) error
}

// Env provides environment dependencies,
// such as environment variables and current directory.
type Env struct{}

func (e *Env) Getenv(key string) string {
	return os.Getenv(key)
}

func (e *Env) Chdir(dir string) error {
	return os.Chdir(dir)
}
