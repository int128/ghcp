package mock_adaptors

import (
	"testing"

	"github.com/int128/ghcp/adaptors/interfaces"
)

func NewLogger(t *testing.T) adaptors.Logger {
	return &testLogger{t}
}

type testLogger struct {
	t *testing.T
}

func (l *testLogger) Infof(format string, v ...interface{}) {
	l.t.Logf(format, v...)
}

func (l *testLogger) Debugf(format string, v ...interface{}) {
	l.t.Logf(format, v...)
}
