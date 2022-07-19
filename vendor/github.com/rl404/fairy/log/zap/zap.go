// Package zap is a wrapper of the original "github.com/uber-go/zap" library.
package zap

import (
	"fmt"
	"reflect"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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
	log *zap.Logger
}

// New to create new logging client.
// Color is not working in json format.
func New(level LogLevel, jsonFmt, color bool) *Log {
	cfg := zap.Config{
		Level:             zap.NewAtomicLevelAt(convertLevel[level]),
		Development:       false,
		DisableStacktrace: true,
		DisableCaller:     true,
		Encoding:          convertFmt[jsonFmt],
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:     "message",
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "N",
			FunctionKey:    zapcore.OmitKey,
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
				enc.AppendString(colorize(t.Format(time.RFC3339), 90, !color || jsonFmt))
			},
		},
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	if jsonFmt {
		cfg.EncoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
	}

	if color && !jsonFmt {
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	l, _ := cfg.Build()

	return &Log{
		log: l,
	}
}

var convertFmt = map[bool]string{
	false: "console",
	true:  "json",
}

var convertLevel = map[LogLevel]zapcore.Level{
	DebugLevel: zapcore.DebugLevel,
	InfoLevel:  zapcore.InfoLevel,
	WarnLevel:  zapcore.WarnLevel,
	ErrorLevel: zapcore.ErrorLevel,
	FatalLevel: zapcore.FatalLevel,
	PanicLevel: zapcore.PanicLevel,
}

func colorize(s interface{}, c int, disabled bool) string {
	if disabled {
		return fmt.Sprintf("%s", s)
	}
	return fmt.Sprintf("\x1b[%dm%v\x1b[0m", c, s)
}

// Trace to print trace log.
// Do nothing because zap doesn't handle trace log.
func (l *Log) Trace(fmt string, args ...interface{}) {
}

// Debug to print debug log.
func (l *Log) Debug(fmt string, args ...interface{}) {
	l.log.Sugar().Debugf(fmt, args...)
}

// Info to print info log.
func (l *Log) Info(fmt string, args ...interface{}) {
	l.log.Sugar().Infof(fmt, args...)
}

// Warn to print warn log.
func (l *Log) Warn(fmt string, args ...interface{}) {
	l.log.Sugar().Warnf(fmt, args...)
}

// Error to print error log.
func (l *Log) Error(fmt string, args ...interface{}) {
	l.log.Sugar().Errorf(fmt, args...)
}

// Fatal to print fatal log.
// Will exit the program when called.
func (l *Log) Fatal(fmt string, args ...interface{}) {
	l.log.Sugar().Fatalf(fmt, args...)
}

// Panic to print panic log.
// Will print panic error stack and exit
// like panic().
func (l *Log) Panic(fmt string, args ...interface{}) {
	l.log.Sugar().Panicf(fmt, args...)
}

// Log to print general log.
// Key `level` can be used to differentiate
// log level.
func (l *Log) Log(fields map[string]interface{}) {
	if len(fields) == 0 {
		return
	}

	var keyVals []interface{}
	for k, v := range fields {
		keyVals = append(keyVals, k, v)
	}

	if level, ok := fields["level"]; ok {
		switch reflect.TypeOf(level).Kind() {
		case reflect.Int8:
			delete(fields, "level")
			switch LogLevel(reflect.ValueOf(level).Int()) {
			case DebugLevel:
				l.log.Sugar().Debugw("", keyVals...)
			case InfoLevel:
				l.log.Sugar().Infow("", keyVals...)
			case WarnLevel:
				l.log.Sugar().Warnw("", keyVals...)
			case ErrorLevel:
				l.log.Sugar().Errorw("", keyVals...)
			case FatalLevel:
				l.log.Sugar().Fatalw("", keyVals...)
			case PanicLevel:
				l.log.Sugar().Panicw("", keyVals...)
			}
		}
	}
}
