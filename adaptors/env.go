package adaptors

import (
	"os"

	"github.com/int128/ghcp/adaptors/interfaces"
)

func NewEnv() adaptors.Env {
	return &Env{}
}

// Env provides access to environment variables.
type Env struct{}

func (e *Env) Get(key string) string {
	return os.Getenv(key)
}
