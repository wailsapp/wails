package commands

import (
	"fmt"
	"github.com/pterm/pterm"
	"github.com/wailsapp/wails/v3/internal/debug"
	"github.com/wailsapp/wails/v3/internal/github"
	"github.com/wailsapp/wails/v3/internal/term"
	"github.com/wailsapp/wails/v3/internal/version"
	"os"
	"os/exec"
	"path/filepath"
)

type UpdateCLIOptions struct {
	NoColour   bool   `name:"n" description:"Disable colour output"`
	PreRelease bool   `name:"pre" description:"Update to the latest pre-release (eg beta)"`
	Version    string `name:"version" description:"Update to a specific version (eg v3.0.0)"`
	Latest     bool   `name:"latest" description:"Install the latest stable release"`
}

func UpdateCLI(options *UpdateCLIOptions) error {
	if options.NoColour {
		term.DisableColor()
	}

	term.Header("Update CLI")

	// Check if this CLI has been installed from vcs
	if debug.LocalModulePath != "" && !options.Latest {
		v3Path := filepath.ToSlash(debug.LocalModulePath + "/v3")
		term.Println("This Wails CLI has been installed from source. To update to the latest stable release, run the following commands in the `" + v3Path + "` directory:")
		term.Println("   - git pull")
		term.Println("   - go install")
		term.Println("")
		term.Println("If you want to install the latest release, please run `wails update cli -latest`")
		return nil
	}

	if options.Latest {
		latestVersion, err := github.GetLatestStableRelease()
		if err != nil {
			return err
		}
		return updateToVersion(latestVersion, true, version.String())
	}

	term.Println("Checking for updates...")

	var desiredVersion *github.SemanticVersion
	var err error
	var valid bool

	if len(options.Version) > 0 {
		// Check if this is a valid version
		valid, err = github.IsValidTag(options.Version)
		if err == nil {
			if !valid {
				err = fmt.Errorf("version '%s' is invalid", options.Version)
			} else {
				desiredVersion, err = github.NewSemanticVersion(options.Version)
			}
		}
	} else {
		if options.PreRelease {
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

	currentVersion := version.String()
	pterm.Printf("    Current Version : %s\n", currentVersion)

	if len(options.Version) > 0 {
		fmt.Printf("    Desired Version : v%s\n", desiredVersion)
	} else {
		if options.PreRelease {
			fmt.Printf("  Latest Prerelease : v%s\n", desiredVersion)
		} else {
			fmt.Printf("     Latest Release : v%s\n", desiredVersion)
		}
	}

	return updateToVersion(desiredVersion, len(options.Version) > 0, currentVersion)
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
			pterm.Printf("If this is what you really want to do, use `wails update -version %s`\n", targetVersionString)
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
		return fmt.Errorf("cannot find home directory: %w", err)
	}

	cmd := exec.Command("go", "install", "github.com/wailsapp/wails/v3/cmd/wails@"+desiredVersion)
	cmd.Dir = homeDir
	sout, serr := cmd.CombinedOutput()
	if err := cmd.Run(); err != nil {
		pterm.Println("Failed.")
		pterm.Error.Println(string(sout) + "\n" + serr.Error())
		return err
	}
	pterm.Println("Done.")
	pterm.Println("\nMake sure you update your project go.mod file to use " + desiredVersion + ":")
	pterm.Println("  require github.com/wailsapp/wails/v3 " + desiredVersion)
	pterm.Println("\nTo view the release notes, please run `wails3 releasenotes`")

	return nil
}
