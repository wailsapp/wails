package main

import (
	"fmt"
	"log"
	"os"

	"github.com/leaanthony/spinner"
	"github.com/wailsapp/wails/cmd"
)

func init() {

	var prereleaseRequired bool
	var specificVersion string

	// var forceRebuild = false
	checkSpinner := spinner.NewSpinner()
	checkSpinner.SetSpinSpeed(50)

	commandDescription := `This command allows you to update your version of Wails.`
	updateCmd := app.Command("update", "Update to newer [pre]releases or specific versions").
		LongDescription(commandDescription).
		BoolFlag("pre", "Update to latest Prerelease", &prereleaseRequired).
		StringFlag("version", "Install a specific version (Overrides other flags)", &specificVersion)

	updateCmd.Action(func() error {

		message := "Checking for updates..."
		logger.PrintSmallBanner(message)
		fmt.Println()

		// Get versions
		checkSpinner.Start(message)

		github := cmd.NewGitHubHelper()
		var desiredVersion *cmd.SemanticVersion
		var err error
		var valid bool

		if len(specificVersion) > 0 {
			// Check if this is a valid version
			valid, err = github.IsValidTag(specificVersion)
			if err == nil {
				if !valid {
					err = fmt.Errorf("version '%s' is invalid", specificVersion)
				} else {
					desiredVersion, err = cmd.NewSemanticVersion(specificVersion)
				}
			}
		} else {
			if prereleaseRequired {
				desiredVersion, err = github.GetLatestPreRelease()
			} else {
				desiredVersion, err = github.GetLatestStableRelease()
			}
		}
		if err != nil {
			checkSpinner.Error(err.Error())
			return err
		}
		checkSpinner.Success()
		fmt.Println()

		fmt.Println("    Current Version : " + cmd.Version)

		if len(specificVersion) > 0 {
			fmt.Printf("    Desired Version : v%s\n", desiredVersion)
		} else {
			if prereleaseRequired {
				fmt.Printf("  Latest Prerelease : v%s\n", desiredVersion)
			} else {
				fmt.Printf("     Latest Release : v%s\n", desiredVersion)
			}
		}

		return updateToVersion(desiredVersion, len(specificVersion) > 0)
	})
}

func updateToVersion(targetVersion *cmd.SemanticVersion, force bool) error {

	var targetVersionString = "v" + targetVersion.String()

	// Early exit
	if targetVersionString == cmd.Version {
		logger.Green("Looks like you're up to date!")
		return nil
	}

	var desiredVersion string

	if !force {

		compareVersion := cmd.Version

		currentVersion, err := cmd.NewSemanticVersion(compareVersion)
		if err != nil {
			return err
		}

		var success bool

		// Release -> Pre-Release = Massage current version to prerelease format
		if targetVersion.IsPreRelease() && currentVersion.IsRelease() {
			testVersion, err := cmd.NewSemanticVersion(compareVersion + "-0")
			if err != nil {
				return err
			}
			success, _ = targetVersion.IsGreaterThan(testVersion)
		}
		// Pre-Release -> Release = Massage target version to prerelease format
		if targetVersion.IsRelease() && currentVersion.IsPreRelease() {
			// We are ok with greater than or equal
			mainversion := currentVersion.MainVersion()
			targetVersion, err = cmd.NewSemanticVersion(targetVersion.String())
			if err != nil {
				return err
			}
			success, _ = targetVersion.IsGreaterThanOrEqual(mainversion)
		}

		// Release -> Release = Standard check
		if (targetVersion.IsRelease() && currentVersion.IsRelease()) ||
			(targetVersion.IsPreRelease() && currentVersion.IsPreRelease()) {

			success, _ = targetVersion.IsGreaterThan(currentVersion)
		}

		// Compare
		if !success {
			logger.Red("The requested version is lower than the current version.")
			logger.Red("If this is what you really want to do, use `wails update -version %s`", targetVersionString)
			return nil
		}

		desiredVersion = "v" + targetVersion.String()

	} else {
		desiredVersion = "v" + targetVersion.String()
	}

	fmt.Println()
	updateSpinner := spinner.NewSpinner()
	updateSpinner.SetSpinSpeed(40)
	updateSpinner.Start("Installing Wails " + desiredVersion)

	// Run command in non module directory
	homeDir, err := os.UserHomeDir()
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
