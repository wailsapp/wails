package logger

import (
	"fmt"
	"os"
)

// DefaultLogger is a utlility to log messages to a number of destinations
type DefaultLogger struct {}

// NewDefaultLogger creates a new Logger.
func NewDefaultLogger() Logger {
	return &DefaultLogger{}
}

// Print works like Sprintf.
func (l *DefaultLogger) Print(message string, args ...interface{}) error {
	fmt.Printf(message + "\n", args...)
	return nil
}
// Trace level logging. Works like Sprintf.
func (l *DefaultLogger) Trace(message string, args ...interface{}) error {
	fmt.Printf("TRACE | " + message + "\n", args...)
	return nil
}
// Debug level logging. Works like Sprintf.
func (l *DefaultLogger) Debug(message string, args ...interface{}) error {
	fmt.Printf("DEBUG | " + message + "\n", args...)
	return nil
}

// Info level logging. Works like Sprintf.
func (l *DefaultLogger) Info(message string, args ...interface{}) error {
	fmt.Printf("INFO  | " + message + "\n", args...)
	return nil
}

// Warning level logging. Works like Sprintf.
func (l *DefaultLogger) Warning(message string, args ...interface{}) error {
	fmt.Printf("WARN  | " + message + "\n", args...)
	return nil
}

// Error level logging. Works like Sprintf.
func (l *DefaultLogger) Error(message string, args ...interface{}) error {
	fmt.Printf("ERROR | " + message + "\n", args...)
	return nil
}

// Fatal level logging. Works like Sprintf.
func (l *DefaultLogger) Fatal(message string, args ...interface{}) error {
	fmt.Printf("FATAL | " + message + "\n", args...)
	os.Exit(1)
	return nil
}
