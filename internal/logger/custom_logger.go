package logger

import (
	"fmt"
	"os"
)

// CustomLogger defines what a user can do with a logger
type CustomLogger interface {
	// Writeln writes directly to the output with no log level plus line ending
	Writeln(message string) error

	// Write writes directly to the output with no log level
	Write(message string) error

	// Trace level logging. Works like Sprintf.
	Trace(format string, args ...interface{}) error

	// Debug level logging. Works like Sprintf.
	Debug(format string, args ...interface{}) error

	// Info level logging. Works like Sprintf.
	Info(format string, args ...interface{}) error

	// Warning level logging. Works like Sprintf.
	Warning(format string, args ...interface{}) error

	// Error level logging. Works like Sprintf.
	Error(format string, args ...interface{}) error

	// Fatal level logging. Works like Sprintf.
	Fatal(format string, args ...interface{})
}

// customLogger is a utlility to log messages to a number of destinations
type customLogger struct {
	logger *Logger
	name   string
}

// New creates a new customLogger. You may pass in a number of `io.Writer`s that
// are the targets for the logs
func newcustomLogger(logger *Logger, name string) *customLogger {
	result := &customLogger{
		name:   name,
		logger: logger,
	}
	return result
}

// Writeln writes directly to the output with no log level
// Appends a carriage return to the message
func (l *customLogger) Writeln(message string) error {
	return l.logger.Writeln(message)
}

// Write writes directly to the output with no log level
func (l *customLogger) Write(message string) error {
	return l.logger.Write(message)
}

// Trace level logging. Works like Sprintf.
func (l *customLogger) Trace(format string, args ...interface{}) error {
	format = fmt.Sprintf("%s | %s", l.name, format)
	return l.logger.processLogMessage(TRACE, format, args...)
}

// Debug level logging. Works like Sprintf.
func (l *customLogger) Debug(format string, args ...interface{}) error {
	format = fmt.Sprintf("%s | %s", l.name, format)
	return l.logger.processLogMessage(DEBUG, format, args...)
}

// Info level logging. Works like Sprintf.
func (l *customLogger) Info(format string, args ...interface{}) error {
	format = fmt.Sprintf("%s | %s", l.name, format)
	return l.logger.processLogMessage(INFO, format, args...)
}

// Warning level logging. Works like Sprintf.
func (l *customLogger) Warning(format string, args ...interface{}) error {
	format = fmt.Sprintf("%s | %s", l.name, format)
	return l.logger.processLogMessage(WARNING, format, args...)
}

// Error level logging. Works like Sprintf.
func (l *customLogger) Error(format string, args ...interface{}) error {
	format = fmt.Sprintf("%s | %s", l.name, format)
	return l.logger.processLogMessage(ERROR, format, args...)

}

// Fatal level logging. Works like Sprintf.
func (l *customLogger) Fatal(format string, args ...interface{}) {
	format = fmt.Sprintf("%s | %s", l.name, format)
	l.logger.processLogMessage(FATAL, format, args...)
	os.Exit(1)
}
