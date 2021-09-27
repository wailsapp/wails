package main

import (
	"runtime"

	"github.com/wailsapp/wails/cmd"
)

func init() {

	commandDescription := `Sets up your local environment to develop Wails apps.`

	setupCommand := app.Command("setup", "Setup the Wails environment").
		LongDescription(commandDescription)

	app.DefaultCommand(setupCommand)

	setupCommand.Action(func() error {

		logger.PrintBanner()

		var err error

		system := cmd.NewSystemHelper()
		err = system.Initialise()
		if err != nil {
			return err
		}

		var successMessage = `Ready for take off!
Create your first project by running 'wails init'.`
		if runtime.GOOS != "windows" {
			successMessage = "🚀 " + successMessage
		}

		// Chrck for programs and libraries dependencies
		errors, err := cmd.CheckDependencies(logger)
		if err != nil {
			return err
		}

		// Check for errors
		// CheckDependencies() returns !errors
		// so to get the right message in this
		// check we have to do it in reversed
		if errors {
			logger.Yellow(successMessage)
		}

		return err
	})
}
