package log

import (
	_ "embed"
	"fmt"
	"io"
	"os"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed plugin.js
var pluginJS string

// ---------------- Plugin Setup ----------------
// This is the main plugin struct. It can be named anything you like.
// It must implement the application.Plugin interface.
// Both the Init() and Shutdown() methods are called synchronously when the app starts and stops.

type LogLevel = float64

const (
	Trace LogLevel = iota + 1
	Debug
	Info
	Warning
	Error
	Fatal
)

type Config struct {
	// Where the logs are written to. Defaults to os.Stderr
	// If you want to write to a file, use os.OpenFile()
	// e.g. os.OpenFile("mylog.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	// Closes the writer when the app shuts down
	Writer io.WriteCloser

	// The initial log level. Defaults to Debug
	Level LogLevel

	// Disables the log level prefixes
	DisablePrefix bool

	// Handles errors that occur when writing to the log
	ErrorHandler func(err error)
}

type Plugin struct {
	config *Config
	app    *application.App
	level  LogLevel
}

func NewPluginWithConfig(config *Config) *Plugin {
	if config.Level == 0 {
		config.Level = Debug
	}
	if config.Writer == nil {
		config.Writer = os.Stderr
	}
	return &Plugin{
		config: config,
		level:  config.Level,
	}
}

func NewPlugin() *Plugin {
	return NewPluginWithConfig(&Config{})
}

// Shutdown is called when the app is shutting down
// You can use this to clean up any resources you have allocated
func (p *Plugin) Shutdown() {
	p.config.Writer.Close()
}

// Name returns the name of the plugin.
// You should use the go module format e.g. github.com/myuser/myplugin
func (p *Plugin) Name() string {
	return "github.com/wailsapp/wails/v3/plugins/log"
}

func (p *Plugin) Init(app *application.App) error {
	p.app = app
	return nil
}

// CallableByJS returns a list of methods that can be called from the frontend
func (p *Plugin) CallableByJS() []string {
	return []string{
		"Trace",
		"Debug",
		"Info",
		"Warning",
		"Error",
		"Fatal",
		"SetLevel",
	}
}

func (p *Plugin) InjectJS() string {
	return pluginJS
}

// ---------------- Plugin Methods ----------------
// Plugin methods are just normal Go methods. You can add as many as you like.
// The only requirement is that they are exported (start with a capital letter).
// You can also return any type that is JSON serializable.
// See https://golang.org/pkg/encoding/json/#Marshal for more information.

func (p *Plugin) write(prefix string, level LogLevel, message string, args ...any) {
	if level >= p.level {
		if !p.config.DisablePrefix {
			message = prefix + " " + message
		}
		_, err := fmt.Fprintln(p.config.Writer, fmt.Sprintf(message, args...))
		if err != nil && p.config.ErrorHandler != nil {
			p.config.ErrorHandler(err)
		}
	}
}

func (p *Plugin) Trace(message string, args ...any) {
	p.write("[Trace]", Trace, message, args...)
}

func (p *Plugin) Debug(message string, args ...any) {
	p.write("[Debug]", Debug, message, args...)
}

func (p *Plugin) Info(message string, args ...any) {
	p.write("[Info]", Info, message, args...)
}

func (p *Plugin) Warning(message string, args ...any) {
	p.write("[Warning]", Warning, message, args...)
}

func (p *Plugin) Error(message string, args ...any) {
	p.write("[Error]", Error, message, args...)
}

func (p *Plugin) Fatal(message string, args ...any) {
	p.write("[FATAL]", Fatal, message, args...)
}

func (p *Plugin) SetLevel(newLevel LogLevel) {
	if newLevel == 0 {
		newLevel = Debug
	}
	p.level = newLevel
}
