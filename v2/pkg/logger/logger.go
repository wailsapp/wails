package logger

type Logger interface {
	Print(message string, args ...interface{}) error
	Trace(message string, args ...interface{}) error
	Debug(message string, args ...interface{}) error
	Info(message string, args ...interface{}) error
	Warning(message string, args ...interface{}) error
	Error(message string, args ...interface{}) error
	Fatal(message string, args ...interface{}) error
}