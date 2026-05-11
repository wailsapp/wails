package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/wailsapp/wails/v3/internal/debug"
	"github.com/wailsapp/wails/v3/internal/github"
	"github.com/wailsapp/wails/v3/internal/term"
	"github.com/wailsapp/wails/v3/internal/version"
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

	if debug.LocalModulePath != "" && !options.Latest {
		v3Path := filepath.ToSlash(debug.LocalModulePath + "/v3")
		term.Println("This Wails CLI was installed from source. To update, run in `" + v3Path + "`:")
		term.Println("   git pull")
		term.Println("   wails3 task install")
		term.Println("")
		term.Println("To install the latest release instead, run `wails3 update cli -latest`.")
		return nil
	}

	if options.Latest {
		latestVersion, err := github.GetLatestStableRelease()
		if err != nil {
			return err
		}
		return updateToVersion(latestVersion, true, version.String())
	}

	term.Println("Checking for updates…")

	var desiredVersion *github.SemanticVersion
	var err error
	var valid bool

	if len(options.Version) > 0 {
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
				fmt.Println()
				term.Println("No stable release found for this major version. To update to the latest pre-release, run:")
				term.Println("   wails3 update cli -pre")
				return nil
			}
		}
	}
	if err != nil {
		return err
	}

	fmt.Println()
	currentVersion := version.String()
	term.Printf("  Current Version : %s\n", currentVersion)
	if len(options.Version) > 0 {
		term.Printf("  Desired Version : v%s\n", desiredVersion)
	} else {
		if options.PreRelease {
			term.Printf("Latest Prerelease : v%s\n", desiredVersion)
		} else {
			term.Printf("   Latest Release : v%s\n", desiredVersion)
		}
	}

	return updateToVersion(desiredVersion, len(options.Version) > 0, currentVersion)
}

func updateToVersion(targetVersion *github.SemanticVersion, force bool, currentVersion string) error {
	targetVersionString := "v" + targetVersion.String()

	if targetVersionString == currentVersion {
		fmt.Println()
		term.Success("Already up to date!")
		return nil
	}

	var desiredVersion string

	if !force {
		compareVersion := currentVersion
		cur, err := github.NewSemanticVersion(compareVersion)
		if err != nil {
			return err
		}

		var success bool

		if targetVersion.IsPreRelease() && cur.IsRelease() {
			testVersion, err := github.NewSemanticVersion(compareVersion + "-0")
			if err != nil {
				return err
			}
			success, _ = targetVersion.IsGreaterThan(testVersion)
		}
		if targetVersion.IsRelease() && cur.IsPreRelease() {
			mainversion := cur.MainVersion()
			targetVersion, err = github.NewSemanticVersion(targetVersion.String())
			if err != nil {
				return err
			}
			success, _ = targetVersion.IsGreaterThanOrEqual(mainversion)
		}
		if (targetVersion.IsRelease() && cur.IsRelease()) ||
			(targetVersion.IsPreRelease() && cur.IsPreRelease()) {
			success, _ = targetVersion.IsGreaterThan(cur)
		}

		if !success {
			term.Warning("The requested version is lower than the current version.")
			term.Printf("Use `wails3 update cli -version %s` to force.\n", targetVersionString)
			return nil
		}

		desiredVersion = "v" + targetVersion.String()
	} else {
		desiredVersion = "v" + targetVersion.String()
	}

	fmt.Println()
	spinner := term.StartSpinner("Installing Wails CLI " + desiredVersion)

	homeDir, err := os.UserHomeDir()
	if err != nil {
		term.StopSpinner(spinner)
		return fmt.Errorf("cannot find home directory: %w", err)
	}

	cmd := exec.Command("go", "install", "github.com/wailsapp/wails/v3/cmd/wails@"+desiredVersion)
	cmd.Dir = homeDir
	out, cmdErr := cmd.CombinedOutput()

	if cmdErr != nil {
		spinner.Fail("Installation failed")
		term.Error(string(out))
		return cmdErr
	}

	spinner.Success("Installed Wails CLI " + desiredVersion)
	fmt.Println()
	term.Println("Update your project go.mod:")
	term.Println("  require github.com/wailsapp/wails/v3 " + desiredVersion)
	term.Println("")
	term.Println("View release notes with `wails3 releasenotes`")

	return nil
}
