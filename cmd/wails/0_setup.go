package main

import (
	"fmt"
	"runtime"

	"github.com/wailsapp/wails/cmd"
)

func init() {

	commandDescription := `Sets up your local environment to develop Wails apps.`

	setupCommand := app.Command("setup", "Setup the Wails environment").
		LongDescription(commandDescription)

	app.DefaultCommand(setupCommand)

	setupCommand.Action(func() error {

		logger.PrintBanner();

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
		// Platform check
		err = platformCheck()
		if err != nil {
			return err
		}

		// Check we have a cgo capable environment
		logger.Yellow("Checking for prerequisites...")
		errors, err := checkRequiredPrograms()
		if err != nil {
			return err
		}

		// Linux has library deps
		errors, err = checkLibraries()
		if err != nil {
			return err
		}

		// Check Mewn
		err = cmd.CheckMewn()
		if err != nil {
			return err
		}

		logger.White("")

		if !errors {
			logger.Yellow(successMessage)
		}

		return err
	})
}

func platformCheck() error {
	switch runtime.GOOS {
	case "darwin":
		logger.Yellow("Detected Platform: OSX")
	case "windows":
		logger.Yellow("Detected Platform: Windows")
	case "linux":
		logger.Yellow("Detected Platform: Linux")
	default:
		return fmt.Errorf("Platform %s is currently not supported", runtime.GOOS)
	}
	return nil
}

func checkLibraries() (errors bool, err error) {
	if runtime.GOOS == "linux" {
		// Check library prerequisites
		requiredLibraries, err := cmd.GetRequiredLibraries()
		if err != nil {
			return false, err
		}
		distroInfo := cmd.GetLinuxDistroInfo()
		for _, library := range *requiredLibraries {
			switch distroInfo.Distribution {
			case cmd.Ubuntu:
				installed, err := cmd.DpkgInstalled(library.Name)
				if err != nil {
					return false, err
				}
				if !installed {
					errors = true
					logger.Red("Library '%s' not found. %s", library.Name, library.Help)
				} else {
					logger.Green("Library '%s' installed.", library.Name)
				}
			default:
				return false, fmt.Errorf("unable to check libraries on distribution '%s'. Please ensure that the '%s' equivalent is installed", distroInfo.DistributorID, library.Name)
			}
		}
	}
	return false, nil
}

func checkRequiredPrograms() (errors bool, err error) {
	requiredPrograms, err := cmd.GetRequiredPrograms()
	if err != nil {
		return true, err
	}
	errors = false
	programHelper := cmd.NewProgramHelper()
	for _, program := range *requiredPrograms {
		bin := programHelper.FindProgram(program.Name)
		if bin == nil {
			errors = true
			logger.Red("Program '%s' not found. %s", program.Name, program.Help)
		} else {
			logger.Green("Program '%s' found: %s", program.Name, bin.Path)
		}
	}
	return
}
