// Package chain is a logger chain. Useful if you want to print the log in
// local and send the log to third party at the same time.
package chain

import "github.com/rl404/fairy/log"

// Chain is logger chain.
type Chain struct {
	loggers []log.Logger
}

// New to create new logger chain.
func New(logger log.Logger, loggers ...log.Logger) log.Logger {
	return &Chain{
		loggers: append([]log.Logger{logger}, loggers...),
	}
}

// Trace to print trace log.
func (lc *Chain) Trace(format string, args ...interface{}) {
	for _, l := range lc.loggers {
		l.Trace(format, args...)
	}
}

// Debug to print debug log.
func (lc *Chain) Debug(format string, args ...interface{}) {
	for _, l := range lc.loggers {
		l.Debug(format, args...)
	}
}

// Info to print info log.
func (lc *Chain) Info(format string, args ...interface{}) {
	for _, l := range lc.loggers {
		l.Info(format, args...)
	}
}

// Warn to print warn log.
func (lc *Chain) Warn(format string, args ...interface{}) {
	for _, l := range lc.loggers {
		l.Warn(format, args...)
	}
}

// Error to print error log.
func (lc *Chain) Error(format string, args ...interface{}) {
	for _, l := range lc.loggers {
		l.Error(format, args...)
	}
}

// Fatal to print fatal log.
func (lc *Chain) Fatal(format string, args ...interface{}) {
	for _, l := range lc.loggers {
		l.Fatal(format, args...)
	}
}

// Panic to print panic log.
func (lc *Chain) Panic(format string, args ...interface{}) {
	for _, l := range lc.loggers {
		l.Panic(format, args...)
	}
}

// Log to print general log.
// Key `level` can be used to differentiate
// log level.
func (lc *Chain) Log(fields map[string]interface{}) {
	for _, l := range lc.loggers {
		// Make a copy of the fields, so when previous logger
		// modified the fields, the next logger won't be affected.
		l.Log(lc.copyMap(fields))
	}
}

func (lc *Chain) copyMap(m1 map[string]interface{}) map[string]interface{} {
	m2 := make(map[string]interface{})
	for k, v := range m1 {
		m2[k] = v
	}
	return m2
}
