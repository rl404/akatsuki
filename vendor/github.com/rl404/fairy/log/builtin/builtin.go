// Package builtin is a wrapper of the built-in log package.
package builtin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"time"
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
	log   *log.Logger
	level LogLevel
	json  bool
	color bool
}

var colorMap = map[LogLevel]int{
	TraceLevel: 94,
	DebugLevel: 93,
	InfoLevel:  32,
	WarnLevel:  33,
	ErrorLevel: 91,
	FatalLevel: 31,
	PanicLevel: 35,
}

// New to create new logging client.
// Color is not working in json format.
func New(level LogLevel, jsonFmt, color bool) *Log {
	return &Log{
		log:   log.New(os.Stderr, "", 0),
		level: level,
		json:  jsonFmt,
		color: color,
	}
}

func (l *Log) colorize(c int, s interface{}) string {
	if !l.color {
		return fmt.Sprintf("%s", s)
	}
	return fmt.Sprintf("\x1b[%dm%v\x1b[0m", c, s)
}

func (l *Log) print(lvl LogLevel, fields map[string]interface{}, str string, args ...interface{}) {
	if l.json {
		b := &bytes.Buffer{}
		b.WriteByte('{')
		b.WriteString(fmt.Sprintf(`"level":"%s"`, l.getLevelStr(lvl)))
		b.WriteString(fmt.Sprintf(`,"time":"%s"`, time.Now().Format(time.RFC3339)))
		if str != "" {
			b.WriteString(fmt.Sprintf(`,"message":"%s"`, fmt.Sprintf(str, args...)))
		}
		for k, v := range fields {
			j, _ := json.Marshal(v)
			b.WriteString(fmt.Sprintf(`,"%s":%s`, k, j))
		}
		b.WriteByte('}')
		l.log.Println(b)
	} else {
		b := &bytes.Buffer{}
		b.WriteString(l.colorize(90, time.Now().Format(time.RFC3339)))
		b.WriteByte(' ')
		l.printLevel(lvl, b)
		b.WriteByte(' ')
		b.WriteString(fmt.Sprintf(str, args...))
		for k, v := range fields {
			if k == "error" {
				b.WriteString(l.colorize(31, fmt.Sprintf(" %s=%v", k, v)))
			} else {
				b.WriteString(fmt.Sprintf(" %s=%v", l.colorize(36, k), v))
			}
		}
		l.log.Println(b)
	}
}

func (l *Log) printLevel(lvl LogLevel, b *bytes.Buffer) {
	switch lvl {
	case TraceLevel:
		b.WriteString(l.colorize(colorMap[lvl], "TRC"))
	case DebugLevel:
		b.WriteString(l.colorize(colorMap[lvl], "DBG"))
	case InfoLevel:
		b.WriteString(l.colorize(colorMap[lvl], "INF"))
	case WarnLevel:
		b.WriteString(l.colorize(colorMap[lvl], "WRN"))
	case ErrorLevel:
		b.WriteString(l.colorize(colorMap[lvl], "ERR"))
	case FatalLevel:
		b.WriteString(l.colorize(colorMap[lvl], "FTL"))
	case PanicLevel:
		b.WriteString(l.colorize(colorMap[lvl], "PNC"))
	}
}

func (l *Log) getLevelStr(lvl LogLevel) string {
	switch lvl {
	case TraceLevel:
		return "trace"
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warn"
	case ErrorLevel:
		return "error"
	case FatalLevel:
		return "fatal"
	case PanicLevel:
		return "panic"
	default:
		return "log"
	}
}

// Trace to print trace log.
func (l *Log) Trace(fmt string, args ...interface{}) {
	if l.level <= TraceLevel {
		l.print(TraceLevel, nil, fmt, args...)
	}
}

// Debug to print debug log.
func (l *Log) Debug(fmt string, args ...interface{}) {
	if l.level <= DebugLevel {
		l.print(DebugLevel, nil, fmt, args...)
	}
}

// Info to print info log.
func (l *Log) Info(fmt string, args ...interface{}) {
	if l.level <= InfoLevel {
		l.print(InfoLevel, nil, fmt, args...)
	}
}

// Warn to print warn log.
func (l *Log) Warn(fmt string, args ...interface{}) {
	if l.level <= WarnLevel {
		l.print(WarnLevel, nil, fmt, args...)
	}
}

// Error to print error log.
func (l *Log) Error(fmt string, args ...interface{}) {
	if l.level <= ErrorLevel {
		l.print(ErrorLevel, nil, fmt, args...)
	}
}

// Fatal to print fatal log.
// Will exit the program when called.
func (l *Log) Fatal(fmt string, args ...interface{}) {
	if l.level <= FatalLevel {
		l.print(FatalLevel, nil, fmt, args...)
		os.Exit(1)
	}
}

// Panic to print panic log.
// Will print panic error stack and exit
// like panic().
func (l *Log) Panic(str string, args ...interface{}) {
	if l.level <= PanicLevel {
		l.print(PanicLevel, nil, str, args...)
		panic(fmt.Sprintf(str, args...))
	}
}

// Log to print general log.
// Key `level` can be used to differentiate
// log level.
func (l *Log) Log(fields map[string]interface{}) {
	level, ok := fields["level"]
	if ok {
		switch reflect.TypeOf(level).Kind() {
		case reflect.Int8:
			delete(fields, "level")
			l.print(LogLevel(reflect.ValueOf(level).Int()), fields, "")
			return
		}
	}
	l.print(0, fields, "")
}
