package main

import (
	wails "github.com/wailsapp/wails/v2"
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
