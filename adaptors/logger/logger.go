package logger

import (
	"log"

	"github.com/google/wire"
	"github.com/int128/ghcp/adaptors"
)

var Set = wire.NewSet(
	wire.Struct(new(Logger)),
	wire.Bind(new(adaptors.Logger), new(*Logger)),
	wire.Bind(new(adaptors.LoggerConfig), new(*Logger)),
)

// Logger provides logging using Go standard package.
// By default debug logs are not shown but you can change it by LoggerConfig.
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
