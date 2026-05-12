package main

import (
	"fmt"
	"os"

	"github.com/wailsapp/wails/v2/cmd/wails/flags"
	"github.com/wailsapp/wails/v2/internal/github"
	"github.com/wailsapp/wails/v2/internal/shell"
	"github.com/wailsapp/wails/v2/internal/tui"
)

func update(f *flags.Update) error {
	if f.NoColour {
		tui.SetNoColour()
	}

	app.PrintBanner()
	fmt.Println("Checking for updates...")

	var desiredVersion *github.SemanticVersion
	var err error
	var valid bool

	if len(f.Version) > 0 {
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
				fmt.Println()
				fmt.Println("No stable release found for this major version. To update to the latest pre-release (eg beta), run:")
				fmt.Println("   wails update -pre")
				return nil
			}
		}
	}
	if err != nil {
		return err
	}
	fmt.Println()

	fmt.Println("    Current Version : " + app.Version())

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
		fmt.Println("\nLooks like you're up to date!")
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

		if targetVersion.IsPreRelease() && currentVersion.IsRelease() {
			testVersion, err := github.NewSemanticVersion(compareVersion + "-0")
			if err != nil {
				return err
			}
			success, _ = targetVersion.IsGreaterThan(testVersion)
		}
		if targetVersion.IsRelease() && currentVersion.IsPreRelease() {
			mainversion := currentVersion.MainVersion()
			targetVersion, err = github.NewSemanticVersion(targetVersion.String())
			if err != nil {
				return err
			}
			success, _ = targetVersion.IsGreaterThanOrEqual(mainversion)
		}

		if (targetVersion.IsRelease() && currentVersion.IsRelease()) ||
			(targetVersion.IsPreRelease() && currentVersion.IsPreRelease()) {
			success, _ = targetVersion.IsGreaterThan(currentVersion)
		}

		if !success {
			fmt.Println("Error: The requested version is lower than the current version.")
			fmt.Println(fmt.Sprintf("If this is what you really want to do, use `wails update -version "+"%s`", targetVersionString))
			return nil
		}

		desiredVersion = "v" + targetVersion.String()
	} else {
		desiredVersion = "v" + targetVersion.String()
	}

	fmt.Println()

	var sout, serr string
	err := tui.WithSpinner("Installing Wails CLI "+desiredVersion, func() error {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("cannot find home directory: %w", err)
		}
		sout, serr, err = shell.RunCommand(homeDir, "go", "install", "github.com/wailsapp/wails/v2/cmd/wails@"+desiredVersion)
		return err
	})
	if err != nil {
		tui.Error(sout + "\n" + serr)
		return err
	}

	fmt.Println()
	fmt.Println(tui.Green("Make sure you update your project go.mod file to use " + desiredVersion + ":"))
	fmt.Println(tui.Green("  require github.com/wailsapp/wails/v2 " + desiredVersion))
	fmt.Println(tui.Red("\nTo view the release notes, please run `wails show releasenotes`"))

	return nil
}
