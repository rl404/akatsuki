// Package nolog is just a mock of logging which do nothing.
package nolog

// Log is mock of logging.
type Log struct{}

// New to create new mock logging client.
func New() *Log {
	return &Log{}
}

// Trace to do nothing.
func (l *Log) Trace(fmt string, args ...interface{}) {}

// Debug to do nothing.
func (l *Log) Debug(fmt string, args ...interface{}) {}

// Info to do nothing.
func (l *Log) Info(fmt string, args ...interface{}) {}

// Warn to do nothing.
func (l *Log) Warn(fmt string, args ...interface{}) {}

// Error to do nothing.
func (l *Log) Error(fmt string, args ...interface{}) {}

// Fatal to do nothing.
func (l *Log) Fatal(fmt string, args ...interface{}) {}

// Panic to do nothing.
func (l *Log) Panic(fmt string, args ...interface{}) {}

// Log to do nothing.
func (l *Log) Log(fields map[string]interface{}) {}
