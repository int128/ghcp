package env

import (
	"os"

	"github.com/int128/ghcp/adaptors"
)

func NewEnv() adaptors.Env {
	return &Env{}
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
