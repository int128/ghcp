package env

import (
	"os"

	"github.com/google/wire"
	"github.com/int128/ghcp/adaptors"
)

var Set = wire.NewSet(
	wire.Struct(new(Env)),
	wire.Struct(new(FileSystem)),
	wire.Bind(new(adaptors.Env), new(*Env)),
	wire.Bind(new(adaptors.FileSystem), new(*FileSystem)),
)

// Env provides environment dependencies,
// such as environment variables and current directory.
type Env struct{}

func (e *Env) Getenv(key string) string {
	return os.Getenv(key)
}

func (e *Env) Chdir(dir string) error {
	return os.Chdir(dir)
}
