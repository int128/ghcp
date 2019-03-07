package adaptors

import (
	"os"

	"github.com/int128/ghcp/adaptors/interfaces"
)

func NewEnv() adaptors.Env {
	return &Env{}
}

// Env provides environment dependencies,
// such as environment variables.
type Env struct{}

func (e *Env) Getenv(key string) string {
	return os.Getenv(key)
}
