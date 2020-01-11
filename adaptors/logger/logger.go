package logger

import (
	"log"

	"github.com/google/wire"
)

var Set = wire.NewSet(
	wire.Struct(new(Logger)),
	wire.Bind(new(Interface), new(*Logger)),
	wire.Bind(new(Config), new(*Logger)),
)

//go:generate mockgen -destination mock_logger/mock_logger.go github.com/int128/ghcp/adaptors/logger Interface,Config

type Interface interface {
	Errorf(format string, v ...interface{})
	Warnf(format string, v ...interface{})
	Infof(format string, v ...interface{})
	Debugf(format string, v ...interface{})
}

type Config interface {
	SetDebug(debug bool)
}

// Logger provides logging using Go standard package.
// By default debug logs are not shown but you can change it by Config.
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
