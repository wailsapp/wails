package main

import (
	"runtime"

	"github.com/wailsapp/wails/cmd"
)

func init() {

	commandDescription := `Sets up your local environment to develop Wails apps.`

	initCommand := app.Command("setup", "Setup the Wails environment").
		LongDescription(commandDescription)

	initCommand.Action(func() error {

		system := cmd.NewSystemHelper()
		err := system.Initialise()
		if err != nil {
			return err
		}

		var successMessage = `Ready for take off!
Create your first project by running 'wails init'.`
		if runtime.GOOS != "windows" {
			successMessage = "🚀 " + successMessage
		}
		logger.Yellow("Checking for prerequisites...")
		// Check we have a cgo capable environment

		successDeps, failedDeps, err := cmd.CheckBinaryPrerequisites()
		if err != nil {
			return err
		}

		for _, dep := range *successDeps {
			logger.Green("Found '%s' at '%s'", dep.Name, dep.Path)
		}

		logger.White("")

		for _, dep := range *failedDeps {
			logger.Red("PreRequisite '%s' missing. %s", dep.Name, dep.Help)
		}

		// Check non-binary prerequisites
		err = cmd.CheckNonBinaryPrerequisites()

		if err != nil {
			return err
		}

		if len(*failedDeps) == 0 {
			logger.Yellow(successMessage)
		}

		return err
	})
}
