package log

// Logger is logging interface.
//
// See usage example in example folder.
type Logger interface {
	Trace(format string, args ...interface{})
	Debug(format string, args ...interface{})
	Info(format string, args ...interface{})
	Warn(format string, args ...interface{})
	Error(format string, args ...interface{})
	Fatal(format string, args ...interface{})
	Panic(format string, args ...interface{})

	// General log with key value.
	Log(fields map[string]interface{})
}

type logLevel int8

// Available log level.
const (
	traceLevel logLevel = iota - 1
	debugLevel
	infoLevel
	warnLevel
	errorLevel
	fatalLevel
	panicLevel
	disabled
)
