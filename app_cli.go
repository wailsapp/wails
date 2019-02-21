package wails

import (
	"fmt"

	"github.com/wailsapp/wails/cmd"
)

// setupCli creates a new cli handler for the application
func (app *App) setupCli() *cmd.Cli {

	// Create a new cli
	result := cmd.NewCli(app.config.Title, "Debug build")

	// Setup cli to handle loglevel and headless flags
	result.
		StringFlag("loglevel", "Sets the log level [debug|info|error|panic|fatal]. Default debug", &app.logLevel).
		// BoolFlag("headless", "Runs the app in headless mode", &app.headless).
		Action(app.start)

	// Banner
	result.PreRun(func(cli *cmd.Cli) error {
		log := cmd.NewLogger()
		log.PrintSmallBanner()
		fmt.Println()
		log.YellowUnderline(app.config.Title + " - Debug Build")
		return nil
	})

	return result
}
