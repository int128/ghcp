package adaptors

import (
	"log"

	"github.com/int128/ghcp/adaptors/interfaces"
)

// NewLogger returns a Logger.
func NewLogger() adaptors.Logger {
	return &Logger{}
}

// Logger provides logging using Go standard package.
type Logger struct{}

func (*Logger) Infof(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func (*Logger) Debugf(format string, v ...interface{}) {
	log.Printf(format, v...)
}
