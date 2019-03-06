package adaptors

import (
	"log"

	"github.com/int128/ghcp/adaptors/interfaces"
)

// NewLogger returns a Logger and LoggerConfig.
// By default debug logs are not shown but you can change it by LoggerConfig.
func NewLogger() (adaptors.Logger, adaptors.LoggerConfig) {
	var l Logger
	return &l, &l
}

// Logger provides logging using Go standard package.
type Logger struct {
	debug bool
}

func (*Logger) Errorf(format string, v ...interface{}) {
	log.Printf("ERROR "+format, v...)
}

func (*Logger) Warnf(format string, v ...interface{}) {
	log.Printf("WARN  "+format, v...)
}

func (*Logger) Infof(format string, v ...interface{}) {
	log.Printf("INFO  "+format, v...)
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	if l.debug {
		log.Printf("DEBUG "+format, v...)
	}
}

func (l *Logger) SetDebug(debug bool) {
	l.debug = debug
}
