package logger

import (
	"os"

	"github.com/fatih/color"
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
	c := color.New(color.FgHiGreen).SprintFunc()
	println(c("DEB |"), message)
}

// Info level logging. Works like Sprintf.
func (l *DefaultLogger) Info(message string) {
	c := color.New(color.FgBlue).Add(color.Underline).SprintFunc()
	println(c("INF |"), message)
}

// Warning level logging. Works like Sprintf.
func (l *DefaultLogger) Warning(message string) {
	c := color.New(color.FgHiYellow).Add(color.Bold).SprintFunc()
	println(c("WAR |"), message)
}

// Error level logging. Works like Sprintf.
func (l *DefaultLogger) Error(message string) {
	c := color.New(color.FgRed).Add(color.Bold).SprintFunc()
	println(c("ERR |"), message)
}

// Fatal level logging. Works like Sprintf.
func (l *DefaultLogger) Fatal(message string) {
	c := color.New(color.BgRed).Add(color.Bold).SprintFunc()
	println(c("FAT |"), message)
	os.Exit(1)
}
