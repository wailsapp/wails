package runtime

import (
	"context"
	"github.com/wailsapp/wails/v2/pkg/logger"
)

// LogPrint prints a Print level message
func LogPrint(ctx context.Context, message string) {
	myLogger := getLogger(ctx)
	myLogger.Print(message)
}

// LogTrace prints a Trace level message
func LogTrace(ctx context.Context, message string) {
	myLogger := getLogger(ctx)
	myLogger.Trace(message)
}

// LogDebug prints a Debug level message
func LogDebug(ctx context.Context, message string) {
	myLogger := getLogger(ctx)
	myLogger.Debug(message)
}

// LogInfo prints a Info level message
func LogInfo(ctx context.Context, message string) {
	myLogger := getLogger(ctx)
	myLogger.Info(message)
}

// LogWarning prints a Warning level message
func LogWarning(ctx context.Context, message string) {
	myLogger := getLogger(ctx)
	myLogger.Warning(message)
}

// LogError prints a Error level message
func LogError(ctx context.Context, message string) {
	myLogger := getLogger(ctx)
	myLogger.Error(message)
}

// LogFatal prints a Fatal level message
func LogFatal(ctx context.Context, message string) {
	myLogger := getLogger(ctx)
	myLogger.Fatal(message)
}

// LogSetLogLevel sets the log level
func LogSetLogLevel(ctx context.Context, level logger.LogLevel) {
	myLogger := getLogger(ctx)
	myLogger.SetLogLevel(level)
}
