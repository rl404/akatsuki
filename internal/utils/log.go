package utils

import (
	_log "github.com/rl404/fairy/log"
	"github.com/rl404/fairy/log/chain"
	"github.com/rl404/fairy/log/nop"
	"github.com/rl404/fairy/log/zerolog"
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

var l _log.Logger = zerolog.New(zerolog.LogLevel(TraceLevel), false, true)
var ls []_log.Logger = []_log.Logger{l}

// InitLog to init global logger.
func InitLog(lvl LogLevel, json, color bool) {
	l = zerolog.New(zerolog.LogLevel(lvl), json, color)
	ls = []_log.Logger{l}
}

// AddLog to add logger chain.
func AddLog(l1 _log.Logger) {
	l = chain.New(l, l1)
	ls = append(ls, l1)
}

// GetLogger to get logger.
func GetLogger(i ...int) _log.Logger {
	if len(i) > 0 {
		if len(ls) <= i[0] {
			return nop.New()
		}
		return ls[i[0]]
	}
	return l
}

// Log to log with custom field.
func Log(field map[string]interface{}) {
	GetLogger().Log(field)
}

// Info to log info.
func Info(str string, args ...interface{}) {
	GetLogger().Info(str, args...)
}

// Error to print error.
func Error(str string, args ...interface{}) {
	GetLogger().Error(str, args...)
}

// Fatal to log fatal.
func Fatal(str string, args ...interface{}) {
	GetLogger().Fatal(str, args...)
}
