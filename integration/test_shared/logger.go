package test_shared

import (
	"testing"
)

type MockedLogger struct {
	silentLogger bool
	t            *testing.T
}

func NewMockedLogger(t *testing.T) *MockedLogger {
	return &MockedLogger{t: t, silentLogger: false}
}

func NewMockedSilentLogger(t *testing.T) *MockedLogger {
	return &MockedLogger{t: t, silentLogger: true}
}

func (l *MockedLogger) Error(args ...interface{}) {
	if l.silentLogger {
		return
	}
	l.t.Log(args...)
}

func (l *MockedLogger) Debug(args ...interface{}) {
	if l.silentLogger {
		return
	}
	l.t.Log(args...)
}

func (l *MockedLogger) Warn(args ...interface{}) {
	if l.silentLogger {
		return
	}
	l.t.Log(args...)
}

func (l *MockedLogger) Info(args ...interface{}) {
	if l.silentLogger {
		return
	}
	l.t.Log(args...)
}

func (l *MockedLogger) Errorf(template string, args ...interface{}) {
	if l.silentLogger {
		return
	}
	l.t.Logf(template, args...)
}

func (l *MockedLogger) Warnf(template string, args ...interface{}) {
	if l.silentLogger {
		return
	}
	l.t.Logf(template, args...)
}

func (l *MockedLogger) Infof(template string, args ...interface{}) {
	if l.silentLogger {
		return
	}
	l.t.Logf(template, args...)
}

func (l *MockedLogger) Printf(format string, v ...interface{}) {
	if l.silentLogger {
		return
	}
	l.t.Logf(format, v...)
}
