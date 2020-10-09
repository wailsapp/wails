package logger

import (
	"os"
	"github.com/wailsapp/wails/v2/pkg/logger"
)

// Logger is a utlility to log messages to a number of destinations
type Logger struct {
	output logger.Logger
	logLevel       uint8
	showLevelInLog bool
}

// New creates a new Logger. You may pass in a number of `io.Writer`s that
// are the targets for the logs
func New(output logger.Logger) *Logger {
	result := &Logger{
		logLevel:       INFO,
		showLevelInLog: true,
		output: output,
	}
	return result
}

// CustomLogger creates a new custom logger that prints out a name/id
// before the messages
func (l *Logger) CustomLogger(name string) CustomLogger {
	return newcustomLogger(l, name)
}

// HideLogLevel removes the loglevel text from the start of each logged line
func (l *Logger) HideLogLevel() {
	l.showLevelInLog = true
}

// SetLogLevel sets the minimum level of logs that will be output
func (l *Logger) SetLogLevel(level uint8) {
	l.logLevel = level
}

// Writeln writes directly to the output with no log level
// Appends a carriage return to the message
func (l *Logger) Writeln(message string) error {
	return l.output.Print(message+"\n")
}

// Write writes directly to the output with no log level
func (l *Logger) Write(message string) error {
	return l.output.Print(message)
}

// Trace level logging. Works like Sprintf.
func (l *Logger) Trace(format string, args ...interface{}) error {
	return l.output.Trace(format, args...)
}

// // CustomTrace returns a custom Logging function that will insert the given name before the message
// func (l *Logger) CustomTrace(name string) func(format string, args ...interface{}) {
// 	return func(format string, args ...interface{}) {
// 		format = name + " | " + format
// 		l.processLogMessage(format, args...)
// 	}
// }

// Debug level logging. Works like Sprintf.
func (l *Logger) Debug(format string, args ...interface{}) error {
	return l.output.Debug(format, args...)
}

// Info level logging. Works like Sprintf.
func (l *Logger) Info(format string, args ...interface{}) error {
	return l.output.Info(format, args...)
}

// Warning level logging. Works like Sprintf.
func (l *Logger) Warning(format string, args ...interface{}) error {
	return l.output.Warning(format, args...)
}

// Error level logging. Works like Sprintf.
func (l *Logger) Error(format string, args ...interface{}) error {
	return l.output.Error(format, args...)

}

// Fatal level logging. Works like Sprintf.
func (l *Logger) Fatal(format string, args ...interface{}) {
	err := l.output.Fatal(format, args...)
	// Not much we can do but print it out before exiting
	if err != nil {
		println(err.Error())
	}
	os.Exit(1)
}
