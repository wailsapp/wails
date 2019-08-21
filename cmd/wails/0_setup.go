package main

import (
	"fmt"
	"log"
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
		// Platform check
		err = platformCheck()
		if err != nil {
			return err
		}

		// Check we have a cgo capable environment
		logger.Yellow("Checking for prerequisites...")
		var requiredProgramErrors bool
		requiredProgramErrors, err = checkRequiredPrograms()
		if err != nil {
			return err
		}

		// Linux has library deps
		var libraryErrors bool
		libraryErrors, err = checkLibraries()
		if err != nil {
			return err
		}

		// Check Mewn
		err = cmd.CheckMewn()
		if err != nil {
			return err
		}

		logger.White("")

		// Check for errors
		var errors = libraryErrors || requiredProgramErrors
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
			case cmd.Ubuntu, cmd.Debian, cmd.Zorin, cmd.Parrot, cmd.Linuxmint, cmd.Elementary:
				installed, err := cmd.DpkgInstalled(library.Name)
				if err != nil {
					return false, err
				}
				if !installed {
					errors = true
					logger.Error("Library '%s' not found. %s", library.Name, library.Help)
				} else {
					logger.Green("Library '%s' installed.", library.Name)
				}
			case cmd.Arch:
				installed, err := cmd.PacmanInstalled(library.Name)
				if err != nil {
					return false, err
				}
				if !installed {
					errors = true
					logger.Error("Library '%s' not found. %s", library.Name, library.Help)
				} else {
					logger.Green("Library '%s' installed.", library.Name)
				}
			case cmd.CentOS, cmd.Fedora:
				installed, err := cmd.RpmInstalled(library.Name)
				if err != nil {
					return false, err
				}
				if !installed {
					errors = true
					logger.Error("Library '%s' not found. %s", library.Name, library.Help)
				} else {
					logger.Green("Library '%s' installed.", library.Name)
				}
			case cmd.Gentoo:
				installed, err := cmd.EqueryInstalled(library.Name)
				if err != nil {
					return false, err
				}
				if !installed {
					errors = true
					logger.Error("Library '%s' not found. %s", library.Name, library.Help)
				} else {
					logger.Green("Library '%s' installed.", library.Name)
				}
			default:
				return false, cmd.RequestSupportForDistribution(distroInfo, library.Name)
			}
		}
	}
	return false, nil
}

func checkRequiredPrograms() (errors bool, err error) {
	// TODO: remove
	log.Println("pre-required programs")
	requiredPrograms, err := cmd.GetRequiredPrograms()
	if err != nil {
		return false, err
	}
	// TODO: remove
	log.Println("post-required programs")
	errors = false
	programHelper := cmd.NewProgramHelper()
	// TODO: remove
	log.Println("post-programhelper programs")
	for _, program := range *requiredPrograms {
		bin := programHelper.FindProgram(program.Name)
		if bin == nil {
			log.Println("aaa")
			errors = true
			logger.Red("Program '%s' not found. %s", program.Name, program.Help)
		} else {
			log.Println("bbb")
			logger.Green("Program '%s' found: %s", program.Name, bin.Path)
		}
	}
	// TODO: remove
	log.Println("end-functio programs")
	return
}
