package logger

import (
	"os"
)

// DefaultLogger is a utlility to log messages to a number of destinations
type DefaultLogger struct{}

// NewDefaultLogger creates a new Logger.
func NewDefaultLogger() Logger {
	return &DefaultLogger{}
}

// Print works like Sprintf.
func (l *DefaultLogger) Print(message string) error {
	println(message)
	return nil
}

// Trace level logging. Works like Sprintf.
func (l *DefaultLogger) Trace(message string) error {
	println("TRACE | " + message)
	return nil
}

// Debug level logging. Works like Sprintf.
func (l *DefaultLogger) Debug(message string) error {
	println("DEBUG | " + message)
	return nil
}

// Info level logging. Works like Sprintf.
func (l *DefaultLogger) Info(message string) error {
	println("INFO  | " + message)
	return nil
}

// Warning level logging. Works like Sprintf.
func (l *DefaultLogger) Warning(message string) error {
	println("WARN  | " + message)
	return nil
}

// Error level logging. Works like Sprintf.
func (l *DefaultLogger) Error(message string) error {
	println("ERROR | " + message)
	return nil
}

// Fatal level logging. Works like Sprintf.
func (l *DefaultLogger) Fatal(message string) error {
	println("FATAL | " + message)
	os.Exit(1)
	return nil
}
