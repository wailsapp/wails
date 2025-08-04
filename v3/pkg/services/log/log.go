package log

import (
	"context"
	_ "embed"
	"log/slog"
	"sync/atomic"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// A Level is the importance or severity of a log event.
// The higher the level, the more important or severe the event.
//
// Values are arbitrary, but there are four predefined ones.
type Level = int

const (
	Debug   = Level(slog.LevelDebug)
	Info    = Level(slog.LevelInfo)
	Warning = Level(slog.LevelWarn)
	Error   = Level(slog.LevelError)
)

type Config struct {
	// Logger is the logger to use. If not set, a default logger will be used.
	Logger *slog.Logger

	// LogLevel defines the log level of the logger.
	LogLevel slog.Level
}

//wails:inject export {
//wails:inject     DebugContext as Debug,
//wails:inject     InfoContext as Info,
//wails:inject     WarningContext as Warning,
//wails:inject     ErrorContext as Error,
//wails:inject };
type LogService struct {
	config atomic.Pointer[Config]
	level  slog.LevelVar
}

// New initialises a logging service with the default configuration.
func New() *LogService {
	return NewWithConfig(nil)
}

// NewWithConfig initialises a logging service with a custom configuration.
func NewWithConfig(config *Config) *LogService {
	result := &LogService{}
	result.Configure(config)
	return result
}

// ServiceName returns the name of the plugin.
// You should use the go module format e.g. github.com/myuser/myplugin
func (l *LogService) ServiceName() string {
	return "github.com/wailsapp/wails/v3/plugins/log"
}

// Configure reconfigures the logger dynamically.
// If config is nil, it falls back to the default configuration.
//
//wails:ignore
func (l *LogService) Configure(config *Config) {
	if config == nil {
		config = &Config{}
	} else {
		// Clone to prevent changes from the outside.
		clone := new(Config)
		*clone = *config
		config = clone
	}

	l.level.Set(slog.Level(config.LogLevel))

	if config.Logger == nil {
		config.Logger = application.DefaultLogger(&l.level)
	}

	l.config.Store(config)
}

// Level returns the currently configured log level,
// that is either the one configured initially
// or the last value passed to [Service.SetLogLevel].
//
// Through this method, [Service] implements the [slog.Leveler] interface.
// The intended use case is to propagate
// the service's dynamic level setting to custom loggers.
// For example:
//
//	logService := log.New()
//	customLogger := slog.New(slog.NewTextHandler(
//		customWriter,
//		&slog.HandlerOptions{
//			Level: logService,
//		},
//	))
//	logService.Configure(&log.Config{
//		Logger: customLogger
//	})
//
// By doing so, setting updates made through [Service.SetLogLevel]
// will propagate dynamically to the custom logger.
//
//wails:ignore
func (l *LogService) Level() slog.Level {
	return l.level.Level()
}

// LogLevel returns the currently configured log level,
// that is either the one configured initially
// or the last value passed to [Service.SetLogLevel].
func (l *LogService) LogLevel() Level {
	return Level(l.Level())
}

// SetLogLevel changes the current log level.
func (l *LogService) SetLogLevel(level Level) {
	l.level.Set(slog.Level(level))
}

// Log emits a log record with the current time and the given level and message.
// The Record's attributes consist of the Logger's attributes followed by
// the attributes specified by args.
//
// The attribute arguments are processed as follows:
//   - If an argument is a string and this is not the last argument,
//     the following argument is treated as the value and the two are combined
//     into an attribute.
//   - Otherwise, the argument is treated as a value with key "!BADKEY".
//
// Log feeds the binding call context into the configured logger,
// so custom handlers may access context values, e.g. the current window.
func (l *LogService) Log(ctx context.Context, level Level, message string, args ...any) {
	l.config.Load().Logger.Log(ctx, slog.Level(level), message, args...)
}

// Debug logs at level [Debug].
//
//wails:ignore
func (l *LogService) Debug(message string, args ...any) {
	l.DebugContext(context.Background(), message, args...)
}

// Info logs at level [Info].
//
//wails:ignore
func (l *LogService) Info(message string, args ...any) {
	l.InfoContext(context.Background(), message, args...)
}

// Warning logs at level [Warning].
//
//wails:ignore
func (l *LogService) Warning(message string, args ...any) {
	l.WarningContext(context.Background(), message, args...)
}

// Error logs at level [Error].
//
//wails:ignore
func (l *LogService) Error(message string, args ...any) {
	l.ErrorContext(context.Background(), message, args...)
}

// DebugContext logs at level [Debug].
//
//wails:internal
func (l *LogService) DebugContext(ctx context.Context, message string, args ...any) {
	l.Log(ctx, Debug, message, args...)
}

// InfoContext logs at level [Info].
//
//wails:internal
func (l *LogService) InfoContext(ctx context.Context, message string, args ...any) {
	l.Log(ctx, Info, message, args...)
}

// WarningContext logs at level [Warn].
//
//wails:internal
func (l *LogService) WarningContext(ctx context.Context, message string, args ...any) {
	l.Log(ctx, Warning, message, args...)
}

// ErrorContext logs at level [Error].
//
//wails:internal
func (l *LogService) ErrorContext(ctx context.Context, message string, args ...any) {
	l.Log(ctx, Error, message, args...)
}
