package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"time"

	"github.com/leaanthony/spinner"
	"github.com/wailsapp/wails/cmd"
)

func init() {

	var forceRebuild = false
	buildSpinner := spinner.NewSpinner()
	buildSpinner.SetSpinSpeed(50)

	commandDescription := `This command builds then serves your application in bridge mode. Useful for developing your app in a browser.`
	initCmd := app.Command("serve", "Run your Wails project in bridge mode.").
		LongDescription(commandDescription).
		BoolFlag("f", "Force rebuild of application components", &forceRebuild)

	initCmd.Action(func() error {
		log := cmd.NewLogger()
		message := "Building Application"
		if forceRebuild {
			message += " (force rebuild)"
		}
		log.WhiteUnderline(message)

		// Project options
		projectOptions := &cmd.ProjectOptions{}

		// Check we are in project directory
		// Check project.json loads correctly
		fs := cmd.NewFSHelper()
		err := projectOptions.LoadConfig(fs.Cwd())
		if err != nil {
			return err
		}

		// Validate config
		// Check if we have a frontend
		err = cmd.ValidateFrontendConfig(projectOptions)
		if err != nil {
			return err
		}

		// Check pre-requisites are installed

		// Program checker
		program := cmd.NewProgramHelper()

		if projectOptions.FrontEnd != nil {
			// npm
			if !program.IsInstalled("npm") {
				return fmt.Errorf("it appears npm is not installed. Please install and run again")
			}
		}

		// packr
		if !program.IsInstalled("packr") {
			buildSpinner.Start("Installing packr...")
			err := program.InstallGoPackage("github.com/gobuffalo/packr/...")
			if err != nil {
				buildSpinner.Error()
				return err
			}
			buildSpinner.Success()
		}

		// Save project directory
		projectDir := fs.Cwd()

		// Install deps
		if projectOptions.FrontEnd != nil {
			// Install frontend deps
			err = os.Chdir(projectOptions.FrontEnd.Dir)
			if err != nil {
				return err
			}

			// Check if frontend deps have been updated
			feSpinner := spinner.New("Installing frontend dependencies (This may take a while)...")
			feSpinner.SetSpinSpeed(50)
			feSpinner.Start()

			requiresNPMInstall := true

			// Read in package.json MD5
			packageJSONMD5, err := fs.FileMD5("package.json")
			if err != nil {
				return err
			}

			const md5sumFile = "package.json.md5"

			// If we aren't forcing the install and the md5sum file exists
			if !forceRebuild && fs.FileExists(md5sumFile) {
				// Yes - read contents
				savedMD5sum, err := fs.LoadAsString(md5sumFile)
				// File exists
				if err == nil {
					// Compare md5
					if savedMD5sum == packageJSONMD5 {
						// Same - no need for reinstall
						requiresNPMInstall = false
						feSpinner.Success("Skipped frontend dependencies (-f to force rebuild)")
					}
				}
			}

			// Md5 sum package.json
			// Different? Build
			if requiresNPMInstall || forceRebuild {
				// Install dependencies
				err = program.RunCommand(projectOptions.FrontEnd.Install)
				if err != nil {
					feSpinner.Error()
					return err
				}
				feSpinner.Success()

				// Update md5sum file
				ioutil.WriteFile(md5sumFile, []byte(packageJSONMD5), 0644)
			}

			bridgeFile := "wailsbridge.prod.js"

			// Copy bridge to project
			_, filename, _, _ := runtime.Caller(1)
			bridgeFileSource := filepath.Join(path.Dir(filename), "..", "assets", "default", bridgeFile)
			bridgeFileTarget := filepath.Join(projectDir, projectOptions.FrontEnd.Dir, projectOptions.FrontEnd.Bridge, "wailsbridge.js")
			err = fs.CopyFile(bridgeFileSource, bridgeFileTarget)
			if err != nil {
				return err
			}

			// Build frontend
			err = cmd.BuildFrontend(projectOptions.FrontEnd.Build)
			if err != nil {
				return err
			}
		}

		// Run packr in project directory
		err = os.Chdir(projectDir)
		if err != nil {
			return err
		}

		// Install dependencies
		err = cmd.InstallGoDependencies()
		if err != nil {
			return err
		}

		buildMode := "bridge"
		err = cmd.BuildApplication(projectOptions.BinaryName, forceRebuild, buildMode)
		if err != nil {
			return err
		}

		logger.Yellow("Awesome! Project '%s' built!", projectOptions.Name)

		go func() {
			time.Sleep(2 * time.Second)
			logger.Green(">>>>> To connect, you will need to run '" + projectOptions.FrontEnd.Serve + "' in the '" + projectOptions.FrontEnd.Dir + "' directory <<<<<")
		}()
		location, err := filepath.Abs(projectOptions.BinaryName)
		if err != nil {
			return err
		}

		logger.Yellow("Serving Application: " + location)
		cmd := exec.Command(location)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			return err
		}

		return nil
	})
}
