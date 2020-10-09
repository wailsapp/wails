package logger

type LogLevel uint8

const (
		// TRACE level
		TRACE LogLevel = 0
	
		// DEBUG level logging
		DEBUG LogLevel = 1
	
		// INFO level logging
		INFO LogLevel = 2
	
		// WARNING level logging
		WARNING LogLevel = 4
	
		// ERROR level logging
		ERROR LogLevel = 8
	
		// FATAL level logging
		FATAL LogLevel = 16
)

type Logger interface {
	Print(message string, args ...interface{}) error
	Trace(message string, args ...interface{}) error
	Debug(message string, args ...interface{}) error
	Info(message string, args ...interface{}) error
	Warning(message string, args ...interface{}) error
	Error(message string, args ...interface{}) error
	Fatal(message string, args ...interface{}) error
}