package main

import (
	"fmt"
	"log"

	"github.com/leaanthony/spinner"
	"github.com/mitchellh/go-homedir"
	"github.com/wailsapp/wails/cmd"
)

func init() {

	var prereleaseRequired bool

	// var forceRebuild = false
	checkSpinner := spinner.NewSpinner()
	checkSpinner.SetSpinSpeed(50)

	commandDescription := `This command allows you to update your version of Wails.`
	updateCmd := app.Command("update", "Check for Updates.").
		LongDescription(commandDescription).
		BoolFlag("pre", "Update to latest Prerelease", &prereleaseRequired)

	updateCmd.Action(func() error {

		message := "Checking for updates..."
		logger.PrintSmallBanner(message)
		fmt.Println()

		// Get versions
		checkSpinner.Start(message)

		github := cmd.NewGitHubHelper()
		var desiredVersion *cmd.SemanticVersion
		var err error

		if prereleaseRequired {
			desiredVersion, err = github.GetLatestPreRelease()
		} else {
			desiredVersion, err = github.GetLatestStableRelease()
		}
		if err != nil {
			checkSpinner.Error(err.Error())
			return err
		}
		checkSpinner.Success()
		fmt.Println()

		fmt.Println("  Current Version   : " + cmd.Version)
		if prereleaseRequired {
			fmt.Printf("  Latest Prerelease : v%s\n", desiredVersion)
		} else {
			fmt.Printf("  Latest Release    : v%s\n", desiredVersion)
		}

		return updateToVersion(desiredVersion)
	})
}

func updateToVersion(version *cmd.SemanticVersion) error {

	// Early exit
	if version.String() == cmd.Version {
		logger.Green("Looks like you're up to date!")
		return nil
	}

	compareVersion := cmd.Version
	if version.IsPreRelease() {
		compareVersion += "-0"
	}

	currentVersion, err := cmd.NewSemanticVersion(compareVersion)
	if err != nil {
		return err
	}

	// Compare
	success, err := version.IsGreaterThan(currentVersion)
	if !success {
		logger.Red("The requested version is lower than the current version. Aborting.")
		return nil
	}

	desiredVersion := "v" + version.String()
	fmt.Println()
	updateSpinner := spinner.NewSpinner()
	updateSpinner.SetSpinSpeed(40)
	updateSpinner.Start("Installing Wails " + desiredVersion)

	// Run command in non module directory
	homeDir, err := homedir.Dir()
	if err != nil {
		log.Fatal("Cannot find home directory! Please file a bug report!")
	}

	err = cmd.NewProgramHelper().RunCommandArray([]string{"go", "get", "github.com/wailsapp/wails/cmd/wails@" + desiredVersion}, homeDir)
	if err != nil {
		updateSpinner.Error(err.Error())
		return err
	}
	updateSpinner.Success()
	fmt.Println()
	logger.Green("Wails updated to " + desiredVersion)

	return nil
}
