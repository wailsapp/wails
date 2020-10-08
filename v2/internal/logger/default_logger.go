package logger

import (
	"fmt"
	"io"
	"os"
	"sync"
)

// Logger is a utlility to log messages to a number of destinations
type Logger struct {
	writers        []io.Writer
	logLevel       uint8
	showLevelInLog bool
	lock           sync.RWMutex
}

// New creates a new Logger. You may pass in a number of `io.Writer`s that
// are the targets for the logs
func New(writers ...io.Writer) *Logger {
	result := &Logger{
		logLevel:       INFO,
		showLevelInLog: true,
	}
	for _, writer := range writers {
		result.AddOutput(writer)
	}
	return result
}

// Writers gets the log writers
func (l *Logger) Writers() []io.Writer {
	return l.writers
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

// AddOutput adds the given `io.Writer` to the list of destinations
// that get logged to
func (l *Logger) AddOutput(writer io.Writer) {
	l.writers = append(l.writers, writer)
}

func (l *Logger) write(loglevel uint8, message string) error {

	// Don't print logs lower than the current log level
	if loglevel < l.logLevel {
		return nil
	}

	// Show log level text if enabled
	if l.showLevelInLog {
		message = mapLogLevel[loglevel] + message
	}

	// write out the logs
	l.lock.Lock()
	for _, writer := range l.writers {
		_, err := writer.Write([]byte(message))
		if err != nil {
			l.lock.Unlock() // Because defer is slow
			return err
		}
	}
	l.lock.Unlock()
	return nil
}

// writeln appends a newline character to the message before writing
func (l *Logger) writeln(loglevel uint8, message string) error {
	return l.write(loglevel, message+"\n")
}

// Writeln writes directly to the output with no log level
// Appends a carriage return to the message
func (l *Logger) Writeln(message string) error {
	return l.write(BYPASS, message+"\n")
}

// Write writes directly to the output with no log level
func (l *Logger) Write(message string) error {
	return l.write(BYPASS, message)
}

// processLogMessage formats the given message before writing it out
func (l *Logger) processLogMessage(loglevel uint8, format string, args ...interface{}) error {
	message := fmt.Sprintf(format, args...)
	return l.writeln(loglevel, message)
}

// Trace level logging. Works like Sprintf.
func (l *Logger) Trace(format string, args ...interface{}) error {
	return l.processLogMessage(TRACE, format, args...)
}

// CustomTrace returns a custom Logging function that will insert the given name before the message
func (l *Logger) CustomTrace(name string) func(format string, args ...interface{}) {
	return func(format string, args ...interface{}) {
		format = name + " | " + format
		l.processLogMessage(TRACE, format, args...)
	}
}

// Debug level logging. Works like Sprintf.
func (l *Logger) Debug(format string, args ...interface{}) error {
	return l.processLogMessage(DEBUG, format, args...)
}

// Info level logging. Works like Sprintf.
func (l *Logger) Info(format string, args ...interface{}) error {
	return l.processLogMessage(INFO, format, args...)
}

// Warning level logging. Works like Sprintf.
func (l *Logger) Warning(format string, args ...interface{}) error {
	return l.processLogMessage(WARNING, format, args...)
}

// Error level logging. Works like Sprintf.
func (l *Logger) Error(format string, args ...interface{}) error {
	return l.processLogMessage(ERROR, format, args...)

}

// Fatal level logging. Works like Sprintf.
func (l *Logger) Fatal(format string, args ...interface{}) {
	l.processLogMessage(FATAL, format, args...)
	os.Exit(1)
}
