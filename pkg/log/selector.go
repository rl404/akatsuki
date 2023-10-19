package log

import (
	"errors"

	"github.com/rl404/fairy/log"
	"github.com/rl404/fairy/log/newrelic"
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

// LogType is type for logging library.
type LogType int8

// Available types for logging.
const (
	NOP LogType = iota
	Zerolog
	Newrelic
)

// ErrInvalidLogType is error for invalid log type.
var ErrInvalidLogType = errors.New("invalid log type")

// Config is log config.
type Config struct {
	Type       LogType
	Level      LogLevel
	JsonFormat bool
	Color      bool

	// For newrelic.
	NewrelicName    string
	NewrelicLicense string
}

// New to create new log client depends on the type.
// Color will not work in json format.
func New(cfg Config) (log.Logger, error) {
	switch cfg.Type {
	case NOP:
		return nop.New(), nil
	case Zerolog:
		return zerolog.New(zerolog.LogLevel(cfg.Level), cfg.JsonFormat, cfg.Color), nil
	case Newrelic:
		return newrelic.New(newrelic.Config{
			Name:       cfg.NewrelicName,
			LicenseKey: cfg.NewrelicLicense,
			Level:      newrelic.LogLevel(cfg.Level),
		})
	default:
		return nil, ErrInvalidLogType
	}
}
