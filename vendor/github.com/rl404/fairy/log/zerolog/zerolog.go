// Package zerolog is a wrapper of the original "github.com/rs/zerolog" library.
package zerolog

import (
	"fmt"
	"os"
	"reflect"
	"runtime"

	"github.com/mattn/go-colorable"
	"github.com/rs/zerolog"
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
	log zerolog.Logger
}

// New to create new logging client.
// Color is not working in json format.
func New(level LogLevel, jsonFmt, color bool) *Log {
	l := zerolog.
		New(os.Stderr).
		With().
		Timestamp().
		Logger().
		Level(zerolog.Level(level))

	if !jsonFmt {
		cw := zerolog.ConsoleWriter{
			Out:     os.Stderr,
			NoColor: !color,
			FormatTimestamp: func(i interface{}) string {
				return colorize(i, 90, !color)
			},
		}
		if runtime.GOOS == "windows" {
			cw.Out = colorable.NewColorableStderr()
		}
		l = l.Output(cw)
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

// Trace to print trace log.
func (l *Log) Trace(fmt string, args ...interface{}) {
	l.log.Trace().Msgf(fmt, args...)
}

// Debug to print debug log.
func (l *Log) Debug(fmt string, args ...interface{}) {
	l.log.Debug().Msgf(fmt, args...)
}

// Info to print info log.
func (l *Log) Info(fmt string, args ...interface{}) {
	l.log.Info().Msgf(fmt, args...)
}

// Warn to print warn log.
func (l *Log) Warn(fmt string, args ...interface{}) {
	l.log.Warn().Msgf(fmt, args...)
}

// Error to print error log.
func (l *Log) Error(fmt string, args ...interface{}) {
	l.log.Error().Msgf(fmt, args...)
}

// Fatal to print fatal log.
// Will exit the program when called.
func (l *Log) Fatal(fmt string, args ...interface{}) {
	l.log.Fatal().Msgf(fmt, args...)
}

// Panic to print panic log.
// Will print panic error stack and exit
// like panic().
func (l *Log) Panic(fmt string, args ...interface{}) {
	l.log.Panic().Msgf(fmt, args...)
}

// Log to print general log.
// Key `level` can be used to differentiate
// log level.
func (l *Log) Log(fields map[string]interface{}) {
	if len(fields) == 0 {
		return
	}

	ll := l.log.Log()
	if level, ok := fields["level"]; ok {
		switch reflect.TypeOf(level).Kind() {
		case reflect.Int8:
			delete(fields, "level")
			switch LogLevel(reflect.ValueOf(level).Int()) {
			case TraceLevel:
				ll = l.log.Trace()
			case DebugLevel:
				ll = l.log.Debug()
			case InfoLevel:
				ll = l.log.Info()
			case WarnLevel:
				ll = l.log.Warn()
			case ErrorLevel:
				ll = l.log.Error()
			case FatalLevel:
				ll = l.log.Fatal()
			case PanicLevel:
				ll = l.log.Panic()
			}
		}
	}

	for k, v := range fields {
		ll.Interface(k, v)
	}

	ll.Send()
}
