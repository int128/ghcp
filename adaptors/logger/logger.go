package logger

import (
	"log"

	"github.com/google/wire"
)

var Set = wire.NewSet(
	wire.Value(NewFunc(New)),
)

type NewFunc func(Option) Interface

//go:generate mockgen -destination mock_logger/mock_logger.go github.com/int128/ghcp/adaptors/logger Interface

type Interface interface {
	Errorf(format string, v ...interface{})
	Warnf(format string, v ...interface{})
	Infof(format string, v ...interface{})
	Debugf(format string, v ...interface{})
}

type Option struct {
	Debug bool
}

func New(o Option) Interface {
	return &logger{opt: o}
}

type logger struct {
	opt Option
}

func (*logger) Errorf(format string, v ...interface{}) {
	log.Printf("ERROR "+format, v...)
}

func (*logger) Warnf(format string, v ...interface{}) {
	log.Printf("WARN  "+format, v...)
}

func (*logger) Infof(format string, v ...interface{}) {
	log.Printf("INFO  "+format, v...)
}

func (l *logger) Debugf(format string, v ...interface{}) {
	if l.opt.Debug {
		log.Printf("DEBUG "+format, v...)
	}
}
