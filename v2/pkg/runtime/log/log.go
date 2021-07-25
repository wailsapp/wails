// +build !experimental

package log

import (
	"context"
	"github.com/wailsapp/wails/v2/internal/servicebus"
	"github.com/wailsapp/wails/v2/pkg/logger"
)

type Log struct{}

// Print prints a Print level message
func (l *Log) Print(ctx context.Context, message string) {
	bus := servicebus.ExtractBus(ctx)
	bus.Publish("log:print", message)
}

// Trace prints a Trace level message
func (l *Log) Trace(ctx context.Context, message string) {
	bus := servicebus.ExtractBus(ctx)
	bus.Publish("log:trace", message)
}

// Debug prints a Debug level message
func (l *Log) Debug(ctx context.Context, message string) {
	bus := servicebus.ExtractBus(ctx)
	bus.Publish("log:debug", message)
}

// Info prints a Info level message
func (l *Log) Info(ctx context.Context, message string) {
	bus := servicebus.ExtractBus(ctx)
	bus.Publish("log:info", message)
}

// Warning prints a Warning level message
func (l *Log) Warning(ctx context.Context, message string) {
	bus := servicebus.ExtractBus(ctx)
	bus.Publish("log:warning", message)
}

// Error prints a Error level message
func (l *Log) Error(ctx context.Context, message string) {
	bus := servicebus.ExtractBus(ctx)
	bus.Publish("log:error", message)
}

// Fatal prints a Fatal level message
func (l *Log) Fatal(ctx context.Context, message string) {
	bus := servicebus.ExtractBus(ctx)
	bus.Publish("log:fatal", message)
}

// SetLogLevel sets the log level
func (l *Log) SetLogLevel(ctx context.Context, level logger.LogLevel) {
	bus := servicebus.ExtractBus(ctx)
	bus.Publish("log:setlevel", level)
}
