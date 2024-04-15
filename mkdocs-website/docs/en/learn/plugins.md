# Wails v3 Plugin Guide

Wails v3 introduces the concept of plugins. A plugin is a self-contained module that can extend the functionality of your Wails application. 
This guide will walk you through the structure and functionality of a Wails plugin.

## Plugin Structure

A Wails plugin is a standard Go module, typically consisting of the following files:

- `plugin.go`: This is the main Go file where the plugin's functionality is implemented.
- `plugin.yaml`: This is the plugin's metadata file. It contains information about the plugin such as its name, author, version, and more.
- `assets/`: This directory contains any static assets that the plugin might need.
- `README.md`: This file provides documentation for the plugin.

## Plugin Implementation

In `plugin.go`, a plugin is defined as a struct that implements the `application.Plugin` interface. This interface requires the following methods:

- `Init()`: This method is called when the plugin is initialized.
- `Shutdown()`: This method is called when the application is shutting down.
- `Name()`: This method returns the name of the plugin.
- `CallableByJS()`: This method returns a list of method names that can be called from the frontend.

In addition to these methods, you can define any number of additional methods that implement the plugin's functionality. 
These methods can be called from the frontend using the `wails.Plugin()` function.

## Plugin Metadata

The `plugin.yaml` file contains metadata about the plugin. This includes the plugin's name, description, author, version, website, repository, and license.

## Plugin Assets

Any static assets that the plugin needs can be placed in the `assets/` directory. 
These assets can be accessed by the frontend by requesting them at the plugin base path.
This path is `/wails/plugins/<plugin-name>/`.

### Example

If a plugin named `logopack` has an asset named `logo.png`, the frontend can access it at `/wails/plugins/logopack/logo.png`.

## Plugin Documentation

The `README.md` file provides documentation for the plugin. This should include instructions on how to install and use the plugin, as well as any other information that users of the plugin might find useful.

## Example

Here's the Log plugin implementation:

```go
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

```

This plugin can be added to the application like this:

```go
    Plugins: map[string]application.Plugin{
        "log": log.NewPlugin(),
    },
```
And here's how you can call a plugin method from the frontend:

```js
    wails.Plugin("log","Debug","hello world")
```

## Support

If you encounter any issues with a plugin, please raise a ticket in the plugin's repository. 

!!! note
    The Wails team does not provide support for third-party plugins.