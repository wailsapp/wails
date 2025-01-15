package log

import (
	"context"
	_ "embed"
	"log/slog"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type Config struct {
	// Logger is the logger to use. If not set, a default logger will be used.
	Logger *slog.Logger

	// LogLevel defines the log level of the logger.
	LogLevel slog.Level

	// Handles errors that occur when writing to the log
	ErrorHandler func(err error)
}

type LoggerService struct {
	config *Config
	app    *application.App
	level  slog.LevelVar
}

func NewLoggerService(config *Config) *LoggerService {
	if config.Logger == nil {
		config.Logger = application.DefaultLogger(config.LogLevel)
	}

	result := &LoggerService{
		config: config,
	}
	result.level.Set(config.LogLevel)
	return result
}

func New() *LoggerService {
	return NewLoggerService(&Config{})
}

// ServiceShutdown is called when the app is shutting down
// You can use this to clean up any resources you have allocated
func (l *LoggerService) ServiceShutdown() error { return nil }

// ServiceName returns the name of the plugin.
// You should use the go module format e.g. github.com/myuser/myplugin
func (l *LoggerService) ServiceName() string {
	return "github.com/wailsapp/wails/v3/plugins/log"
}

func (l *LoggerService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
	// Any initialization code here
	return nil
}

func (l *LoggerService) Debug(message string, args ...any) {
	l.config.Logger.Debug(message, args...)
}

func (l *LoggerService) Info(message string, args ...any) {
	l.config.Logger.Info(message, args...)
}

func (l *LoggerService) Warning(message string, args ...any) {
	l.config.Logger.Warn(message, args...)
}

func (l *LoggerService) Error(message string, args ...any) {
	l.config.Logger.Error(message, args...)
}

func (l *LoggerService) SetLogLevel(level slog.Level) {
	l.level.Set(level)
}
