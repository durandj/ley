package configuration

import (
	"fmt"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// LoggingConfiguration holds any logging specific service configuration.
type LoggingConfiguration struct {
	Level LogLevel `default:"warn"`
}

// LogLevel tells what level of log messages we should keep while throwing
// away the rest.
type LogLevel zapcore.Level

// Decode converts a string value that is pulled from an environment
// variable into the given type.
func (logLevel *LogLevel) Decode(value string) error {
	level, err := zapcore.ParseLevel(strings.ToLower(value))
	if err != nil {
		return fmt.Errorf("Invalid value for LogLevel '%s'", value)
	}

	*logLevel = LogLevel(level)

	return nil
}

// AsAtomicLevel converts the log level value into an atomic level
// that is recognized by Zap.
func (logLevel *LogLevel) AsAtomicLevel() zap.AtomicLevel {
	return zap.NewAtomicLevelAt(zapcore.Level(*logLevel))
}

const (
	// LogLevelDefault gives the default log level if nothing is specified.
	LogLevelDefault LogLevel = LogLevel(zapcore.WarnLevel)

	// LogLevelDebug sets the log level to debug level and higher.
	//
	// Debug can be used anywhere where something useful would be logged
	// but maybe isn't something that you would actually want showing
	// up in production. For example maybe there's some setup work that
	// is useful when figuring out why something is broken during development.
	//
	// This level is also useful what you want to log would be incredibly
	// noisey in production.
	LogLevelDebug LogLevel = LogLevel(zapcore.DebugLevel)

	// LogLevelInfo sets the log level to info level and higher.
	//
	// This level is useful when what you're logging could maybe be
	// useful in some situations in production but isn't incredibly
	// noisey.
	//
	// This is a good default value to use but in general it wouldn't
	// be shown when running in production.
	LogLevelInfo LogLevel = LogLevel(zapcore.InfoLevel)

	// LogLevelWarn sets the log level to warn level and higher.
	//
	// This level is useful when you're logging an anomalous condition
	// that was recoverable but maybe is worth reporting. For example,
	// maybe we received a garbage message on SQS that we're throwing out.
	// We recovered but that's maybe something an operator would care about.
	//
	// This is the level that production should ideally be at.
	LogLevelWarn LogLevel = LogLevel(zapcore.WarnLevel)

	// LogLevelError sets the log level to error level and higher.
	//
	// This level is used whenever a big error occured that is still
	// recoverable but is definitely a problem. For example, if an
	// HTTP request handler failed to handle a request. The service
	// recovered enough to keep working but it resulted in something
	// the consumer noticed.
	LogLevelError LogLevel = LogLevel(zapcore.ErrorLevel)

	// LogLevelFatal sets the log level to fatal level and higher.
	//
	// This is reserved for a condition so bad the service can't
	// recover from it and must terminate.
	LogLevelFatal LogLevel = LogLevel(zapcore.FatalLevel)
)
