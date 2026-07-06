package runtime

import (
	"context"
	"fmt"
	"os"
	"sync"
)

// LogLevel mirrors the v2 logger.LogLevel type.
type LogLevel uint8

// Log levels, mirroring the v2 logger constants.
const (
	TRACE   LogLevel = 1
	DEBUG   LogLevel = 2
	INFO    LogLevel = 3
	WARNING LogLevel = 4
	ERROR   LogLevel = 5
)

// logLevelWarnOnce guards the one-time warning emitted by LogSetLogLevel.
var logLevelWarnOnce sync.Once

// LogPrint mirrors the v2 runtime.LogPrint function.
// v3 equivalent: app.Logger.Info.
func LogPrint(_ context.Context, message string) {
	logger().Info(message)
}

// LogTrace mirrors the v2 runtime.LogTrace function.
// v3 equivalent: app.Logger.Debug.
func LogTrace(_ context.Context, message string) {
	logger().Debug(message)
}

// LogDebug mirrors the v2 runtime.LogDebug function.
// v3 equivalent: app.Logger.Debug.
func LogDebug(_ context.Context, message string) {
	logger().Debug(message)
}

// LogInfo mirrors the v2 runtime.LogInfo function.
// v3 equivalent: app.Logger.Info.
func LogInfo(_ context.Context, message string) {
	logger().Info(message)
}

// LogWarning mirrors the v2 runtime.LogWarning function.
// v3 equivalent: app.Logger.Warn.
func LogWarning(_ context.Context, message string) {
	logger().Warn(message)
}

// LogError mirrors the v2 runtime.LogError function.
// v3 equivalent: app.Logger.Error.
func LogError(_ context.Context, message string) {
	logger().Error(message)
}

// LogFatal mirrors the v2 runtime.LogFatal function. It logs the message at
// error level and exits the process.
// v3 equivalent: app.Logger.Error followed by os.Exit(1).
func LogFatal(_ context.Context, message string) {
	logger().Error(message)
	os.Exit(1)
}

// LogPrintf mirrors the v2 runtime.LogPrintf function.
// v3 equivalent: app.Logger.Info.
func LogPrintf(_ context.Context, format string, args ...interface{}) {
	logger().Info(fmt.Sprintf(format, args...))
}

// LogTracef mirrors the v2 runtime.LogTracef function.
// v3 equivalent: app.Logger.Debug.
func LogTracef(_ context.Context, format string, args ...interface{}) {
	logger().Debug(fmt.Sprintf(format, args...))
}

// LogDebugf mirrors the v2 runtime.LogDebugf function.
// v3 equivalent: app.Logger.Debug.
func LogDebugf(_ context.Context, format string, args ...interface{}) {
	logger().Debug(fmt.Sprintf(format, args...))
}

// LogInfof mirrors the v2 runtime.LogInfof function.
// v3 equivalent: app.Logger.Info.
func LogInfof(_ context.Context, format string, args ...interface{}) {
	logger().Info(fmt.Sprintf(format, args...))
}

// LogWarningf mirrors the v2 runtime.LogWarningf function.
// v3 equivalent: app.Logger.Warn.
func LogWarningf(_ context.Context, format string, args ...interface{}) {
	logger().Warn(fmt.Sprintf(format, args...))
}

// LogErrorf mirrors the v2 runtime.LogErrorf function.
// v3 equivalent: app.Logger.Error.
func LogErrorf(_ context.Context, format string, args ...interface{}) {
	logger().Error(fmt.Sprintf(format, args...))
}

// LogFatalf mirrors the v2 runtime.LogFatalf function. It logs the message at
// error level and exits the process.
// v3 equivalent: app.Logger.Error followed by os.Exit(1).
func LogFatalf(_ context.Context, format string, args ...interface{}) {
	logger().Error(fmt.Sprintf(format, args...))
	os.Exit(1)
}

// LogSetLogLevel mirrors the v2 runtime.LogSetLogLevel function. v3 has no
// runtime log level setter: configure the slog level via
// application.Options.LogLevel instead. This is a no-op that logs a warning once.
func LogSetLogLevel(_ context.Context, level LogLevel) {
	logLevelWarnOnce.Do(func() {
		logger().Warn("v2compat: LogSetLogLevel is a no-op in v3; configure the log level via application.Options.LogLevel", "requestedLevel", level)
	})
}
