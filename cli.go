package wails

import (
	"github.com/wailsapp/wails/cmd"
)

// setupCli creates a new cli handler for the application
func (app *App) setupCli() *cmd.Cli {

	// Create a new cli
	result := cmd.NewCli(app.config.Title, "Debug build")
	result.Version(cmd.Version)

	// Setup cli to handle loglevel
	result.
		StringFlag("loglevel", "Sets the log level [debug|info|error|panic|fatal]. Default debug", &app.logLevel).
		Action(app.start)

	// Banner
	result.PreRun(func(cli *cmd.Cli) error {
		log := cmd.NewLogger()
		log.YellowUnderline(app.config.Title + " - Debug Build")
		return nil
	})

	return result
}
