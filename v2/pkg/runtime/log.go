package runtime

import (
	"github.com/wailsapp/wails/v2/internal/servicebus"
	"github.com/wailsapp/wails/v2/pkg/logger"
)

// Log defines all Log related operations
type Log interface {
	Print(message string)
	Trace(message string)
	Debug(message string)
	Info(message string)
	Warning(message string)
	Error(message string)
	Fatal(message string)
	SetLogLevel(level logger.LogLevel)
}

type log struct {
	bus *servicebus.ServiceBus
}

// newLog creates a new Log struct
func newLog(bus *servicebus.ServiceBus) Log {
	return &log{
		bus: bus,
	}
}

// Print prints a Print level message
func (r *log) Print(message string) {
	r.bus.Publish("log:print", message)
}

// Trace prints a Trace level message
func (r *log) Trace(message string) {
	r.bus.Publish("log:trace", message)
}

// Debug prints a Debug level message
func (r *log) Debug(message string) {
	r.bus.Publish("log:debug", message)
}

// Info prints a Info level message
func (r *log) Info(message string) {
	r.bus.Publish("log:info", message)
}

// Warning prints a Warning level message
func (r *log) Warning(message string) {
	r.bus.Publish("log:warning", message)
}

// Error prints a Error level message
func (r *log) Error(message string) {
	r.bus.Publish("log:error", message)
}

// Fatal prints a Fatal level message
func (r *log) Fatal(message string) {
	r.bus.Publish("log:fatal", message)
}

// Sets the log level
func (r *log) SetLogLevel(level logger.LogLevel) {
	r.bus.Publish("log:setlevel", level)
}
