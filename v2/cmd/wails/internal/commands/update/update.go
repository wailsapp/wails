package update

import (
	"fmt"
	"github.com/labstack/gommon/color"
	"github.com/wailsapp/wails/v2/internal/shell"
	"io"
	"log"
	"os"

	"github.com/wailsapp/wails/v2/internal/github"

	"github.com/leaanthony/clir"
	"github.com/wailsapp/wails/v2/pkg/clilogger"
)

// AddSubcommand adds the `init` command for the Wails application
func AddSubcommand(app *clir.Cli, w io.Writer, currentVersion string) error {

	command := app.NewSubCommand("update", "Update the Wails CLI")
	command.LongDescription(`This command allows you to update your version of the Wails CLI.`)

	// Setup flags
	var prereleaseRequired bool
	command.BoolFlag("pre", "Update CLI to latest Prerelease", &prereleaseRequired)

	var specificVersion string
	command.StringFlag("version", "Install a specific version (Overrides other flags) of the CLI", &specificVersion)

	command.Action(func() error {

		// Create logger
		logger := clilogger.New(w)

		// Print banner
		app.PrintBanner()
		logger.Println("Checking for updates...")

		var desiredVersion *github.SemanticVersion
		var err error
		var valid bool

		if len(specificVersion) > 0 {
			// Check if this is a valid version
			valid, err = github.IsValidTag(specificVersion)
			if err == nil {
				if !valid {
					err = fmt.Errorf("version '%s' is invalid", specificVersion)
				} else {
					desiredVersion, err = github.NewSemanticVersion(specificVersion)
				}
			}
		} else {
			if prereleaseRequired {
				desiredVersion, err = github.GetLatestPreRelease()
			} else {
				desiredVersion, err = github.GetLatestStableRelease()
				if err != nil {
					println("")
					println("No stable release found for this major version. To update to the latest pre-release (eg beta), run:")
					println("   wails update -pre")
					return nil
				}
			}
		}
		if err != nil {
			return err
		}
		fmt.Println()

		fmt.Println("    Current Version : " + currentVersion)

		if len(specificVersion) > 0 {
			fmt.Printf("    Desired Version : v%s\n", desiredVersion)
		} else {
			if prereleaseRequired {
				fmt.Printf("  Latest Prerelease : v%s\n", desiredVersion)
			} else {
				fmt.Printf("     Latest Release : v%s\n", desiredVersion)
			}
		}

		return updateToVersion(logger, desiredVersion, len(specificVersion) > 0, currentVersion)
	})

	return nil
}

func updateToVersion(logger *clilogger.CLILogger, targetVersion *github.SemanticVersion, force bool, currentVersion string) error {

	var targetVersionString = "v" + targetVersion.String()

	// Early exit
	if targetVersionString == currentVersion {
		logger.Println("\nLooks like you're up to date!")
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
			logger.Println("Error: The requested version is lower than the current version.")
			logger.Println("If this is what you really want to do, use `wails update -version %s`", targetVersionString)
			return nil
		}

		desiredVersion = "v" + targetVersion.String()

	} else {
		desiredVersion = "v" + targetVersion.String()
	}

	fmt.Println()
	logger.Print("Installing Wails CLI " + desiredVersion + "...")

	// Run command in non module directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Cannot find home directory! Please file a bug report!")
	}

	sout, serr, err := shell.RunCommand(homeDir, "go", "install", "github.com/wailsapp/wails/v2/cmd/wails@"+desiredVersion)
	if err != nil {
		logger.Println("Failed.")
		logger.Println(sout + `\n` + serr)
		return err
	}
	logger.Println("Done.")
	logger.Println(color.Green("\nMake sure you update your project go.mod file to use " + desiredVersion + ":"))
	logger.Println(color.Green("  require github.com/wailsapp/wails/v2 " + desiredVersion))
	logger.Println(color.Red("\nTo view the release notes, please run `wails show releasenotes`"))

	return nil
}
