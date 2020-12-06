package main

import (
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
)

// Logger struct
type Logger struct {
	runtime *wails.Runtime
}

// WailsInit is called at application startup
func (l *Logger) WailsInit(runtime *wails.Runtime) error {
	// Perform your setup here
	l.runtime = runtime
	return nil
}

// Print will log the given message
func (l *Logger) Print(message string) {
	l.runtime.Log.Print(message)
}

// Trace will log the given message
func (l *Logger) Trace(message string) {
	l.runtime.Log.Trace(message)
}

// Debug will log the given message
func (l *Logger) Debug(message string) {
	l.runtime.Log.Debug(message)
}

// Info will log the given message
func (l *Logger) Info(message string) {
	l.runtime.Log.Info(message)
}

// Warning will log the given message
func (l *Logger) Warning(message string) {
	l.runtime.Log.Warning(message)
}

// Error will log the given message
func (l *Logger) Error(message string) {
	l.runtime.Log.Error(message)
}

// Fatal will log the given message
func (l *Logger) Fatal(message string) {
	l.runtime.Log.Fatal(message)
}

// SetLogLevel will set the given loglevel
func (l *Logger) SetLogLevel(loglevel logger.LogLevel) {
	l.runtime.Log.SetLogLevel(loglevel)
}
