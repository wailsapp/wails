// +build !desktop,!hybrid,!server,!dev

package app

// This is the default application that will get run if the user compiles using `go build`.
// The reason we want to prevent that is that the `wails build` command does a lot of behind
// the scenes work such as asset compilation. If we allow `go build`, the state of these assets
// will be unknown and the application will not work as expected.

import (
	"os"

	"github.com/wailsapp/wails/v2/internal/logger"

	"github.com/wailsapp/wails/v2/pkg/options"
)

// App defines a Wails application structure
type App struct {
	Title     string
	Width     int
	Height    int
	Resizable bool

	// Indicates if the app is running in debug mode
	debug bool

	logger *logger.Logger
}

// CreateApp returns a null application
func CreateApp(_ *options.App) (*App, error) {
	return &App{}, nil
}

// Run the application
func (a *App) Run() error {
	println(`FATAL: This application was built using "go build". This is unsupported. Please compile using "wails build".`)
	os.Exit(1)
	return nil
}
