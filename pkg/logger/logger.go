package logger

import (
	"log"
	"os"

	"github.com/google/wire"
)

var Set = wire.NewSet(
	wire.Value(NewFunc(New)),
)

type NewFunc func(Option) Interface

//go:generate mockgen -destination mock_logger/mock_logger.go github.com/int128/ghcp/pkg/logger Interface

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
	lgr := log.New(os.Stderr, "", log.Lmicroseconds)
	return &logger{opt: o, lgr: lgr}
}

type logger struct {
	lgr *log.Logger
	opt Option
}

func (l *logger) Errorf(format string, v ...interface{}) {
	l.lgr.Printf("ERROR "+format, v...)
}

func (l *logger) Warnf(format string, v ...interface{}) {
	l.lgr.Printf("WARN  "+format, v...)
}

func (l *logger) Infof(format string, v ...interface{}) {
	l.lgr.Printf("INFO  "+format, v...)
}

func (l *logger) Debugf(format string, v ...interface{}) {
	if l.opt.Debug {
		l.lgr.Printf("DEBUG "+format, v...)
	}
}
