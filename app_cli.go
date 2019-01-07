package wails

import (
	"fmt"

	"github.com/wailsapp/wails/cmd"
)

func (app *App) setupCli() *cmd.Cli {

	// var apiFilename string

	result := cmd.NewCli(app.config.Title, "Debug build")

	// Gen API
	// result.Command("genapi", "Generate JS stubs for the registered Go plugins").
	// 	StringFlag("o", "Output filename", &apiFilename).
	// 	Action(func() error {
	// 		app.renderer = N
	// 	})

	result.
		StringFlag("loglevel", "Sets the log level [info|debug|error|panic|fatal]. Default debug", &app.logLevel).
		BoolFlag("headless", "Runs the app in headless mode", &app.headless).
		Action(app.start)

	// Banner
	result.PreRun(func(cli *cmd.Cli) error {
		log := cmd.NewLogger()
		log.PrintBanner()
		fmt.Println()
		result.PrintHelp()
		log.YellowUnderline(app.config.Title + " - Debug Build")
		return nil
	})

	return result
}
