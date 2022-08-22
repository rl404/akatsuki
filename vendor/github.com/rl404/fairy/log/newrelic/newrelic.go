// Package newrelic is a wrapper of the original "github.com/newrelic/go-agent/v3/newrelic".
package newrelic

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"time"

	"github.com/newrelic/go-agent/v3/newrelic"
)

// LogLevel is level of log that will be printed.
// Will print level that is higher than your
// chosen one.
type LogLevel int8

// Available log levec.
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

// Client is newrelic client.
type Client struct {
	client *newrelic.Application
	level  LogLevel
}

// Config is newrelic config.
type Config struct {
	Name       string
	LicenseKey string
	Level      LogLevel
}

// New to create new newrelic client.
func New(cfg Config) (*Client, error) {
	nr, err := newrelic.NewApplication(
		newrelic.ConfigAppName(cfg.Name),
		newrelic.ConfigLicense(cfg.LicenseKey),
		newrelic.ConfigAppLogForwardingEnabled(true),
	)
	if err != nil {
		return nil, err
	}

	// Wait connection.
	if err := nr.WaitForConnection(10 * time.Second); err != nil {
		return nil, err
	}

	return &Client{
		client: nr,
		level:  cfg.Level,
	}, nil
}

// NewFromNewrelic to create new newrelic client with
// existing newrelic app.
//
// Don't forget to use `newrelic.ConfigAppLogForwardingEnabled(true)` when
// initializing your newrelic.
func NewFromNewrelicApp(app *newrelic.Application, level LogLevel) *Client {
	return &Client{
		client: app,
		level:  level,
	}
}

// Trace to print trace log.
func (c *Client) Trace(str string, args ...interface{}) {
	if c.level <= TraceLevel {
		c.log(TraceLevel, str, args...)
	}
}

// Debug to print debug log.
func (c *Client) Debug(str string, args ...interface{}) {
	if c.level <= DebugLevel {
		c.log(DebugLevel, str, args...)
	}
}

// Info to print info log.
func (c *Client) Info(str string, args ...interface{}) {
	if c.level <= InfoLevel {
		c.log(InfoLevel, str, args...)
	}
}

// Warn to print warn log.
func (c *Client) Warn(str string, args ...interface{}) {
	if c.level <= WarnLevel {
		c.log(WarnLevel, str, args...)
	}
}

// Error to print error log.
func (c *Client) Error(str string, args ...interface{}) {
	if c.level <= ErrorLevel {
		c.log(ErrorLevel, str, args...)
	}
}

// Fatal to print fatal log.
// Will exit the program when called.
func (c *Client) Fatal(str string, args ...interface{}) {
	if c.level <= FatalLevel {
		c.log(FatalLevel, str, args...)
		os.Exit(1)
	}
}

// Panic to print panic log.
// Will print panic error stack and exit
// like panic().
func (c *Client) Panic(str string, args ...interface{}) {
	if c.level <= PanicLevel {
		c.log(PanicLevel, str, args...)
		panic(fmt.Sprintf(str, args...))
	}
}

func (c *Client) log(level LogLevel, str string, args ...interface{}) {
	c.Log(map[string]interface{}{
		"level":   level,
		"message": fmt.Sprintf(str, args...),
	})
}

// Log to print general log.
// Key `level` can be used to differentiate
// log level.
func (c *Client) Log(fields map[string]interface{}) {
	if len(fields) == 0 {
		return
	}

	// Handle level field.
	if level, ok := fields["level"]; ok {
		lvl := LogLevel(reflect.ValueOf(level).Int())
		fields["level"] = c.getLevelStr(lvl)
		if c.level > lvl {
			return
		}
	}

	data, _ := json.Marshal(fields)

	c.client.RecordLog(newrelic.LogData{
		Severity: fields["level"].(string),
		Message:  string(data),
	})
}

func (c *Client) getLevelStr(lvl LogLevel) string {
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
