// Package logrus is a wrapper of the original "github.com/sirupsen/logrus" library.
package logrus

import (
	"fmt"
	"reflect"
	"time"

	"github.com/sirupsen/logrus"
)

// LogLevel is level of log that will be printed.
// Will print level that is higher than your
// chosen one.
type LogLevel int8

// Available log level.
const (
	TraceLevel LogLevel = iota - 1
	DebugLevel
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
	PanicLevel
	Disabled
)

// Log is logging client.
type Log struct {
	log *logrus.Logger
}

// New to create new logging client.
// Color is not working in json format.
func New(level LogLevel, jsonFmt, color bool) *Log {
	l := logrus.New()
	l.SetLevel(convertLevel(level))
	l.SetFormatter(&logrus.TextFormatter{
		DisableColors:   !color,
		FullTimestamp:   true,
		TimestampFormat: colorize(time.RFC3339, 90, !color),
	})
	if jsonFmt {
		l.SetFormatter(&logrus.JSONFormatter{})
	}
	return &Log{
		log: l,
	}
}

func colorize(s interface{}, c int, disabled bool) string {
	if disabled {
		return fmt.Sprintf("%s", s)
	}
	return fmt.Sprintf("\x1b[%dm%v\x1b[0m", c, s)
}

func convertLevel(lvl LogLevel) logrus.Level {
	switch lvl {
	case TraceLevel:
		return logrus.TraceLevel
	case DebugLevel:
		return logrus.DebugLevel
	case InfoLevel:
		return logrus.InfoLevel
	case WarnLevel:
		return logrus.WarnLevel
	case ErrorLevel:
		return logrus.ErrorLevel
	case FatalLevel:
		return logrus.FatalLevel
	case PanicLevel:
		return logrus.PanicLevel
	default:
		return 0
	}
}

// Trace to print trace log.
func (l *Log) Trace(fmt string, args ...interface{}) {
	l.log.Tracef(fmt, args...)
}

// Debug to print debug log.
func (l *Log) Debug(fmt string, args ...interface{}) {
	l.log.Debugf(fmt, args...)
}

// Info to print info log.
func (l *Log) Info(fmt string, args ...interface{}) {
	l.log.Infof(fmt, args...)
}

// Warn to print warn log.
func (l *Log) Warn(fmt string, args ...interface{}) {
	l.log.Warnf(fmt, args...)
}

// Error to print error log.
func (l *Log) Error(fmt string, args ...interface{}) {
	l.log.Errorf(fmt, args...)
}

// Fatal to print fatal log.
// Will exit the program when called.
func (l *Log) Fatal(fmt string, args ...interface{}) {
	l.log.Fatalf(fmt, args...)
}

// Panic to print panic log.
// Will print panic error stack and exit
// like panic().
func (l *Log) Panic(fmt string, args ...interface{}) {
	l.log.Panicf(fmt, args...)
}

// Log to print general log.
// Key `level` can be used to differentiate
// log level.
func (l *Log) Log(fields map[string]interface{}) {
	if len(fields) == 0 {
		return
	}

	if level, ok := fields["level"]; ok {
		switch reflect.TypeOf(level).Kind() {
		case reflect.Int8:
			delete(fields, "level")
			switch LogLevel(reflect.ValueOf(level).Int()) {
			case TraceLevel:
				l.log.WithFields(fields).Trace()
			case DebugLevel:
				l.log.WithFields(fields).Debug()
			case InfoLevel:
				l.log.WithFields(fields).Info()
			case WarnLevel:
				l.log.WithFields(fields).Warn()
			case ErrorLevel:
				l.log.WithFields(fields).Error()
			case FatalLevel:
				l.log.WithFields(fields).Fatal()
			case PanicLevel:
				l.log.WithFields(fields).Panic()
			}
		}
	}
}
