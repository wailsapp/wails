package main

import (
	"fmt"
	"runtime"

	"github.com/leaanthony/spinner"
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

		logger.Yellow("Checking for prerequisites...")
		// Check we have a cgo capable environment

		requiredPrograms, err := cmd.GetRequiredPrograms()
		if err != nil {
			return err
		}
		errors := false
		spinner := spinner.New()
		programHelper := cmd.NewProgramHelper()
		for _, program := range *requiredPrograms {
			spinner.Start("Looking for program '%s'", program.Name)
			bin := programHelper.FindProgram(program.Name)
			if bin == nil {
				errors = true
				spinner.Errorf("Program '%s' not found. %s", program.Name, program.Help)
			} else {
				spinner.Successf("Program '%s' found: %s", program.Name, bin.Path)
			}
		}

		// Linux has library deps
		if runtime.GOOS == "linux" {
			// Check library prerequisites
			requiredLibraries, err := cmd.GetRequiredLibraries()
			if err != nil {
				return err
			}
			distroInfo := cmd.GetLinuxDistroInfo()
			for _, library := range *requiredLibraries {
				spinner.Start()
				switch distroInfo.Distribution {
				case cmd.Ubuntu:
					installed, err := cmd.DpkgInstalled(library.Name)
					if err != nil {
						return err
					}
					if !installed {
						errors = true
						spinner.Errorf("Library '%s' not found. %s", library.Name, library.Help)
					} else {
						spinner.Successf("Library '%s' installed.", library.Name)
					}
				default:
					return fmt.Errorf("unable to check libraries on distribution '%s'. Please ensure that the '%s' equivalent is installed", distroInfo.DistributorID, library.Name)
				}
			}
		}

		logger.White("")

		if !errors {
			logger.Yellow(successMessage)
		}

		return err
	})
}
