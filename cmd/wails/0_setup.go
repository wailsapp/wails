package main

import (
	"fmt"
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

		var successMessage string

		logger.Yellow("Checking for prerequisites...")
		// Check we have a cgo capable environment
		programHelper := cmd.NewProgramHelper()
		prerequisites := make(map[string]map[string]string)
		prerequisites["darwin"] = make(map[string]string)
		prerequisites["darwin"]["clang"] = "Please install with `xcode-select --install` and try again"
		prerequisites["darwin"]["npm"] = "Please download and install npm + node from here: https://nodejs.org/en/"
		switch runtime.GOOS {
		case "darwin":
			successMessage = "🚀  Awesome! We are going to the moon! 🚀"
		default:
			return fmt.Errorf("platform '%s' is unsupported at this time", runtime.GOOS)
		}

		errors := false
		for name, help := range prerequisites[runtime.GOOS] {
			bin := programHelper.FindProgram(name)
			if bin == nil {
				errors = true
				logger.Red("Unable to find '%s' - %s", name, help)
			} else {
				logger.Green("Found program '%s' at '%s'", name, bin.Path)
			}
		}

		if errors {
			err = fmt.Errorf("There were missing dependencies")
		} else {
			logger.Yellow(successMessage)
		}

		return err
	})
}
