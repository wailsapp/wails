package log

import (
	"embed"
	_ "embed"
	"github.com/wailsapp/wails/v3/pkg/application"
	"io/fs"
	"log/slog"
)

//go:embed assets/*
var assets embed.FS

// ---------------- Plugin Setup ----------------
// This is the main plugin struct. It can be named anything you like.
// It must implement the application.Plugin interface.
// Both the Init() and Shutdown() methods are called synchronously when the app starts and stops.

type Config struct {
	// Logger is the logger to use. If not set, a default logger will be used.
	Logger *slog.Logger

	// LogLevel defines the log level of the logger.
	LogLevel slog.Level

	// Handles errors that occur when writing to the log
	ErrorHandler func(err error)
}

type Plugin struct {
	config *Config
	app    *application.App
	level  slog.LevelVar
}

func NewPluginWithConfig(config *Config) *Plugin {
	if config.Logger == nil {
		config.Logger = application.DefaultLogger(config.LogLevel)
	}

	result := &Plugin{
		config: config,
	}
	result.level.Set(config.LogLevel)
	return result
}

func NewPlugin() *Plugin {
	return NewPluginWithConfig(&Config{})
}

// Shutdown is called when the app is shutting down
// You can use this to clean up any resources you have allocated
func (p *Plugin) Shutdown() error { return nil }

// Name returns the name of the plugin.
// You should use the go module format e.g. github.com/myuser/myplugin
func (p *Plugin) Name() string {
	return "github.com/wailsapp/wails/v3/plugins/log"
}

func (p *Plugin) Init(api application.PluginAPI) error {
	return nil
}

// CallableByJS returns a list of methods that can be called from the frontend
func (p *Plugin) CallableByJS() []string {
	return []string{
		"Debug",
		"Info",
		"Warning",
		"Error",
		"SetLogLevel",
	}
}

func (p *Plugin) Assets() fs.FS {
	return assets
}

// ---------------- Plugin Methods ----------------
// Plugin methods are just normal Go methods. You can add as many as you like.
// The only requirement is that they are exported (start with a capital letter).
// You can also return any type that is JSON serializable.
// See https://golang.org/pkg/encoding/json/#Marshal for more information.

func (p *Plugin) Debug(message string, args ...any) {
	p.config.Logger.Debug(message, args...)
}

func (p *Plugin) Info(message string, args ...any) {
	p.config.Logger.Info(message, args...)
}

func (p *Plugin) Warning(message string, args ...any) {
	p.config.Logger.Warn(message, args...)
}

func (p *Plugin) Error(message string, args ...any) {
	p.config.Logger.Error(message, args...)
}

func (p *Plugin) SetLogLevel(level slog.Level) {
	p.level.Set(level)
}
