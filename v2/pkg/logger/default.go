package logger

import (
	"os"
)

// DefaultLogger is a utility to log messages to a number of destinations
type DefaultLogger struct{}

// NewDefaultLogger creates a new Logger.
func NewDefaultLogger() Logger {
	return &DefaultLogger{}
}

// Print works like Sprintf.
func (l *DefaultLogger) Print(message string) {
	println(message)
}

// Trace level logging. Works like Sprintf.
func (l *DefaultLogger) Trace(message string) {
	println("TRACE | " + message)
}

// Debug level logging. Works like Sprintf.
func (l *DefaultLogger) Debug(message string) {
	println("DEBUG | " + message)
}

// Info level logging. Works like Sprintf.
func (l *DefaultLogger) Info(message string) {
	println("INFO  | " + message)
}

// Warning level logging. Works like Sprintf.
func (l *DefaultLogger) Warning(message string) {
	println("WARN  | " + message)
}

// Error level logging. Works like Sprintf.
func (l *DefaultLogger) Error(message string) {
	println("ERROR | " + message)
}

// Fatal level logging. Works like Sprintf.
func (l *DefaultLogger) Fatal(message string) {
	println("FATAL | " + message)
	os.Exit(1)
}
