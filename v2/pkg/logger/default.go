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
	println("TRA | " + message)
}

// Debug level logging. Works like Sprintf.
func (l *DefaultLogger) Debug(message string) {
	println("DEB | " + message)
}

// Info level logging. Works like Sprintf.
func (l *DefaultLogger) Info(message string) {
	println("INF | " + message)
}

// Warning level logging. Works like Sprintf.
func (l *DefaultLogger) Warning(message string) {
	println("WAR | " + message)
}

// Error level logging. Works like Sprintf.
func (l *DefaultLogger) Error(message string) {
	println("ERR | " + message)
}

// Fatal level logging. Works like Sprintf.
func (l *DefaultLogger) Fatal(message string) {
	println("FAT | " + message)
	os.Exit(1)
}
