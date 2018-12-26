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

		var successMessage = `Ready for take off!
Create your first project by running 'wails init'.`
		if runtime.GOOS != "windows" {
			successMessage = "🚀 " + successMessage
		}

		system := cmd.NewSystemHelper()
		err := system.Initialise()
		if err != nil {
			return err
		}

		err, success := cmd.CheckDependencies(logger)
		if err != nil {
			return err
		}
		logger.White("")

		if success {
			logger.Yellow(successMessage)
		}
		return nil
	})
}
