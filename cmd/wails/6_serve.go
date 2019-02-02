package main

import (
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"time"

	"github.com/leaanthony/slicer"
	"github.com/leaanthony/spinner"
	"github.com/wailsapp/wails/cmd"
)

func init() {

	var forceRebuild = false
	buildSpinner := spinner.NewSpinner()
	buildSpinner.SetSpinSpeed(50)

	commandDescription := `This command builds then serves your application in bridge mode. Useful for developing your app in a browser.`
	initCmd := app.Command("serve", "Runs your Wails project in bridge mode").
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

		// Program checker
		program := cmd.NewProgramHelper()

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

		// Copy bridge to project
		var bridgeFile = "wailsbridge.js"
		_, filename, _, _ := runtime.Caller(1)
		bridgeFileSource := filepath.Join(path.Dir(filename), "..", "assets", "default", bridgeFile)
		bridgeFileTarget := filepath.Join(projectDir, projectOptions.FrontEnd.Dir, projectOptions.FrontEnd.Bridge, "wailsbridge.js")
		err = fs.CopyFile(bridgeFileSource, bridgeFileTarget)
		if err != nil {
			return err
		}

		// Run packr in project directory
		err = os.Chdir(projectDir)
		if err != nil {
			return err
		}

		// Support build tags
		buildTags := []string{}

		depSpinner := spinner.New("Installing Dependencies...")
		depSpinner.SetSpinSpeed(50)
		depSpinner.Start()
		installCommand := "go get"
		err = program.RunCommand(installCommand)
		if err != nil {
			depSpinner.Error()
			return err
		}
		depSpinner.Success()

		compileMessage := "Packing + Compiling project"

		packSpinner := spinner.New(compileMessage + "...")
		packSpinner.SetSpinSpeed(50)
		packSpinner.Start()

		buildCommand := slicer.String()
		buildCommand.AddSlice([]string{"packr", "build"})

		// Add build tags
		if len(buildTags) > 0 {
			buildCommand.Add("--tags")
			buildCommand.AddSlice(buildTags)

		}

		if projectOptions.BinaryName != "" {
			buildCommand.Add("-o")
			buildCommand.Add(projectOptions.BinaryName)
		}

		// If we are forcing a rebuild
		if forceRebuild {
			buildCommand.Add("-a")
		}

		buildCommand.AddSlice([]string{"-ldflags", "-X github.com/wailsapp/wails.BackendRenderer=headless"})
		// logger.Green("buildCommand = %+v", buildCommand)
		err = program.RunCommandArray(buildCommand.AsSlice())
		if err != nil {
			packSpinner.Error()
			return err
		}
		packSpinner.Success()

		// Run the App
		logger.Yellow("Awesome! Project '%s' built!", projectOptions.Name)
		go func() {
			time.Sleep(2 * time.Second)
			logger.Green(">>>>> To connect, you will need to run '" + projectOptions.FrontEnd.Serve + "' in the '" + projectOptions.FrontEnd.Dir + "' directory <<<<<")
		}()
		logger.Yellow("Serving Application: " + projectOptions.BinaryName)
		cmd := exec.Command(projectOptions.BinaryName)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			return err
		}

		return nil

	})
}
