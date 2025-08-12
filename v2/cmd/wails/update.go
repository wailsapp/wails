package main

import (
	"fmt"
	"os"

	"github.com/labstack/gommon/color"
	"github.com/pterm/pterm"
	"github.com/wailsapp/wails/v2/cmd/wails/flags"
	"github.com/wailsapp/wails/v2/internal/colour"
	"github.com/wailsapp/wails/v2/internal/shell"

	"github.com/wailsapp/wails/v2/internal/github"
)

// AddSubcommand adds the `init` command for the Wails application
func update(f *flags.Update) error {
	if f.NoColour {
		colour.ColourEnabled = false
		pterm.DisableColor()

	}
	// Print banner
	app.PrintBanner()
	pterm.Println("Checking for updates...")

	var desiredVersion *github.SemanticVersion
	var err error
	var valid bool

	if len(f.Version) > 0 {
		// Check if this is a valid version
		valid, err = github.IsValidTag(f.Version)
		if err == nil {
			if !valid {
				err = fmt.Errorf("version '%s' is invalid", f.Version)
			} else {
				desiredVersion, err = github.NewSemanticVersion(f.Version)
			}
		}
	} else {
		if f.PreRelease {
			desiredVersion, err = github.GetLatestPreRelease()
		} else {
			desiredVersion, err = github.GetLatestStableRelease()
			if err != nil {
				pterm.Println("")
				pterm.Println("No stable release found for this major version. To update to the latest pre-release (eg beta), run:")
				pterm.Println("   wails update -pre")
				return nil
			}
		}
	}
	if err != nil {
		return err
	}
	pterm.Println()

	pterm.Println("    Current Version : " + app.Version())

	if len(f.Version) > 0 {
		fmt.Printf("    Desired Version : v%s\n", desiredVersion)
	} else {
		if f.PreRelease {
			fmt.Printf("  Latest Prerelease : v%s\n", desiredVersion)
		} else {
			fmt.Printf("     Latest Release : v%s\n", desiredVersion)
		}
	}

	return updateToVersion(desiredVersion, len(f.Version) > 0, app.Version())
}

func updateToVersion(targetVersion *github.SemanticVersion, force bool, currentVersion string) error {
	targetVersionString := "v" + targetVersion.String()

	if targetVersionString == currentVersion {
		pterm.Println("\nLooks like you're up to date!")
		return nil
	}

	var desiredVersion string

	if !force {

		compareVersion := currentVersion

		currentVersion, err := github.NewSemanticVersion(compareVersion)
		if err != nil {
			return err
		}

		var success bool

		// Release -> Pre-Release = Massage current version to prerelease format
		if targetVersion.IsPreRelease() && currentVersion.IsRelease() {
			testVersion, err := github.NewSemanticVersion(compareVersion + "-0")
			if err != nil {
				return err
			}
			success, _ = targetVersion.IsGreaterThan(testVersion)
		}
		// Pre-Release -> Release = Massage target version to prerelease format
		if targetVersion.IsRelease() && currentVersion.IsPreRelease() {
			// We are ok with greater than or equal
			mainversion := currentVersion.MainVersion()
			targetVersion, err = github.NewSemanticVersion(targetVersion.String())
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
			pterm.Println("Error: The requested version is lower than the current version.")
			pterm.Println(fmt.Sprintf("If this is what you really want to do, use `wails update -version "+"%s`", targetVersionString))

			return nil
		}

		desiredVersion = "v" + targetVersion.String()

	} else {
		desiredVersion = "v" + targetVersion.String()
	}

	pterm.Println()
	pterm.Print("Installing Wails CLI " + desiredVersion + "...")

	// Run command in non module directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fatal("Cannot find home directory! Please file a bug report!")
	}

	sout, serr, err := shell.RunCommand(homeDir, "go", "install", "github.com/wailsapp/wails/v2/cmd/wails@"+desiredVersion)
	if err != nil {
		pterm.Println("Failed.")
		pterm.Error.Println(sout + `\n` + serr)
		return err
	}
	pterm.Println("Done.")
	pterm.Println(color.Green("\nMake sure you update your project go.mod file to use " + desiredVersion + ":"))
	pterm.Println(color.Green("  require github.com/wailsapp/wails/v2 " + desiredVersion))
	pterm.Println(color.Red("\nTo view the release notes, please run `wails show releasenotes`"))

	return nil
}
