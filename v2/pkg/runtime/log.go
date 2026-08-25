package runtime

import (
	"context"
	"fmt"

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
	myLogger.Trace("%s", message)
}

// LogDebug prints a Debug level message
func LogDebug(ctx context.Context, message string) {
	myLogger := getLogger(ctx)
	myLogger.Debug("%s", message)
}

// LogInfo prints a Info level message
func LogInfo(ctx context.Context, message string) {
	myLogger := getLogger(ctx)
	myLogger.Info("%s", message)
}

// LogWarning prints a Warning level message
func LogWarning(ctx context.Context, message string) {
	myLogger := getLogger(ctx)
	myLogger.Warning("%s", message)
}

// LogError prints a Error level message
func LogError(ctx context.Context, message string) {
	myLogger := getLogger(ctx)
	myLogger.Error("%s", message)
}

// LogFatal prints a Fatal level message
func LogFatal(ctx context.Context, message string) {
	myLogger := getLogger(ctx)
	myLogger.Fatal("%s", message)
}

// LogPrintf prints a Print level message
func LogPrintf(ctx context.Context, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	myLogger := getLogger(ctx)
	myLogger.Print(msg)
}

// LogTracef prints a Trace level message
func LogTracef(ctx context.Context, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	myLogger := getLogger(ctx)
	myLogger.Trace("%s", msg)
}

// LogDebugf prints a Debug level message
func LogDebugf(ctx context.Context, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	myLogger := getLogger(ctx)
	myLogger.Debug("%s", msg)
}

// LogInfof prints a Info level message
func LogInfof(ctx context.Context, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	myLogger := getLogger(ctx)
	myLogger.Info("%s", msg)
}

// LogWarningf prints a Warning level message
func LogWarningf(ctx context.Context, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	myLogger := getLogger(ctx)
	myLogger.Warning("%s", msg)
}

// LogErrorf prints a Error level message
func LogErrorf(ctx context.Context, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	myLogger := getLogger(ctx)
	myLogger.Error("%s", msg)
}

// LogFatalf prints a Fatal level message
func LogFatalf(ctx context.Context, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	myLogger := getLogger(ctx)
	myLogger.Fatal("%s", msg)
}

// LogSetLogLevel sets the log level
func LogSetLogLevel(ctx context.Context, level logger.LogLevel) {
	myLogger := getLogger(ctx)
	myLogger.SetLogLevel(level)
}
