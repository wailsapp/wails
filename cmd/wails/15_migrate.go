package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/leaanthony/spinner"
	"github.com/wailsapp/wails/cmd"
)

// Constants
var checkSpinner = spinner.NewSpinner()
var migrateProjectOptions = &cmd.ProjectOptions{}
var migrateFS = cmd.NewFSHelper()
var migrateGithub = cmd.NewGitHubHelper()
var programHelper = cmd.NewProgramHelper()
var lessThanV1 *semver.Constraints

// The user's go.mod
var goMod string
var goModFile string

// The user's main.js
var mainJSFile string
var mainJSContents string

// Frontend directory
var frontEndDir string

func init() {

	var dryrun bool
	var err error

	lessThanV1, err = semver.NewConstraint("< v1.0.0")
	if err != nil {
		log.Fatal(err)
	}

	// var forceRebuild = false
	checkSpinner.SetSpinSpeed(50)

	commandDescription := `EXPERIMENTAL - This command attempts to migrate projects to the latest Wails version.`
	updateCmd := app.Command("migrate", "Migrate projects to latest Wails release").
		LongDescription(commandDescription).
		BoolFlag("dryrun", "Only display what would be done", &dryrun)

	updateCmd.Action(func() error {

		message := "Migrate Project"
		logger.PrintSmallBanner(message)
		logger.Red("WARNING: This is an experimental command. Ensure you have backups of your project!")
		logger.Red("It currently only supports npm based projects.")
		fmt.Println()

		// Check project directory
		err := checkProjectDirectory()
		if err != nil {
			return err
		}

		// Find Wails version from go.mod
		wailsVersion, err := getWailsVersion()
		if err != nil {
			return err
		}

		// Get latest stable version
		var latestVersion *semver.Version
		latestVersion, err = getLatestWailsVersion()
		if err != nil {
			return err
		}

		var canMigrate bool
		canMigrate, err = canMigrateVersion(wailsVersion, latestVersion)
		if err != nil {
			return err
		}

		if !canMigrate {
			return nil
		}

		// Check for wailsbridge
		wailsBridge, err := checkWailsBridge()
		if err != nil {
			return err
		}

		// Is main.js using bridge.Init()
		canUpdateMainJS, err := checkMainJS()
		if err != nil {
			return err
		}

		// TODO: Check if we are using legacy js runtime

		// Operations
		logger.Yellow("Operations to perform:")

		logger.Yellowf("  - Update to Wails v%s\n", latestVersion)

		if len(wailsBridge) > 0 {
			logger.Yellow("  - Delete wailsbridge.js")
		}

		if canUpdateMainJS {
			logger.Yellow("  - Patch main.js")
		}

		logger.Yellow("  - Ensure '@wailsapp/runtime` module is installed")

		if dryrun {
			logger.White("Exiting: Dry Run")
			return nil
		}

		logger.Red("*WARNING* About to modify your project!")
		logger.Red("Type 'YES' to continue: ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		input := scanner.Text()
		if input != "YES" {
			logger.Red("ABORTED!")
			return nil
		}

		logger.Yellow("Let's do this!")

		err = updateWailsVersion(wailsVersion, latestVersion)
		if err != nil {
			return err
		}

		if len(wailsBridge) > 0 {
			err = deleteWailsBridge(wailsBridge)
			if err != nil {
				return err
			}
		}

		if canUpdateMainJS {
			err = patchMainJS()
			if err != nil {
				return err
			}
		}

		// Install runtime
		err = installWailsRuntime()
		if err != nil {
			return err
		}

		fmt.Println()
		logger.Yellow("Migration complete! Check project by running `wails build`.")
		return nil
	})
}

func checkProjectDirectory() error {
	// Get versions
	checkSpinner.Start("Check Project Directory")

	// Check we are in project directory
	err := migrateProjectOptions.LoadConfig(migrateFS.Cwd())
	if err != nil {
		checkSpinner.Error()
		return fmt.Errorf("Unable to find 'project.json'. Please check you are in a Wails project directory")
	}

	checkSpinner.Success()
	return nil
}

func getWailsVersion() (*semver.Version, error) {
	checkSpinner.Start("Get Wails Version")

	result, err := cmd.GetWailsVersion()

	if err != nil {
		checkSpinner.Error(err.Error())
		return nil, err
	}
	return result, nil

}

func canMigrateVersion(wailsVersion *semver.Version, latestVersion *semver.Version) (bool, error) {
	checkSpinner.Start("Checking ability to Migrate")

	// Check if we are at the latest version!!!!
	if wailsVersion.Equal(latestVersion) || wailsVersion.GreaterThan(latestVersion) {
		checkSpinner.Errorf("Checking ability to Migrate: No! (v%s >= v%s)", wailsVersion, latestVersion)
		return false, nil
	}

	// Check for < v1.0.0
	if lessThanV1.Check(wailsVersion) {
		checkSpinner.Successf("Checking ability to Migrate: Yes! (v%s < v1.0.0)", wailsVersion)
		return true, nil
	}
	checkSpinner.Error("Unable to migrate")
	return false, fmt.Errorf("No migration rules for version %s", wailsVersion)
}

func checkWailsBridge() (string, error) {
	checkSpinner.Start("Checking if legacy Wails Bridge present")

	// Check frontend dir is available
	if migrateProjectOptions.FrontEnd == nil ||
		len(migrateProjectOptions.FrontEnd.Dir) == 0 ||
		!migrateFS.DirExists(migrateProjectOptions.FrontEnd.Dir) {
		checkSpinner.Error("Unable to determine frontend directory")
		return "", fmt.Errorf("Unable to determine frontend directory")
	}

	frontEndDir = migrateProjectOptions.FrontEnd.Dir

	wailsBridgePath, err := filepath.Abs(filepath.Join(".", frontEndDir, "src", "wailsbridge.js"))
	if err != nil {
		checkSpinner.Error(err.Error())
		return "", err
	}

	// If it doesn't exist, return blank string
	if !migrateFS.FileExists(wailsBridgePath) {
		checkSpinner.Success("Checking if legacy Wails Bridge present: No")
		return "", nil
	}

	checkSpinner.Success("Checking if legacy Wails Bridge present: Yes")
	return wailsBridgePath, nil

}

// This function determines if the main.js file using wailsbridge can be auto-updated
func checkMainJS() (bool, error) {

	checkSpinner.Start("Checking if main.js can be migrated")
	var err error

	// Check main.js is there
	if migrateProjectOptions.FrontEnd == nil ||
		len(migrateProjectOptions.FrontEnd.Dir) == 0 ||
		!migrateFS.DirExists(migrateProjectOptions.FrontEnd.Dir) {
		checkSpinner.Error("Unable to determine frontend directory")
		return false, fmt.Errorf("Unable to determine frontend directory")
	}

	frontEndDir = migrateProjectOptions.FrontEnd.Dir

	mainJSFile, err = filepath.Abs(filepath.Join(".", frontEndDir, "src", "main.js"))
	if err != nil {
		checkSpinner.Error("Unable to find main.js")
		return false, err
	}

	mainJSContents, err = migrateFS.LoadAsString(mainJSFile)
	if err != nil {
		checkSpinner.Error("Unable to load main.js")
		return false, err
	}

	// Check we have a line like: import Bridge from "./wailsbridge";
	if strings.Index(mainJSContents, `import Bridge from "./wailsbridge";`) == -1 {
		checkSpinner.Success("Checking if main.js can be migrated: No - Cannot find `import Bridge`")
		return false, nil
	}

	// Check we have a line like: Bridge.Start(() => {
	if strings.Index(mainJSContents, `Bridge.Start(`) == -1 {
		checkSpinner.Success("Checking if main.js can be migrated: No - Cannot find `Bridge.Start`")
		return false, nil
	}
	checkSpinner.Success("Checking if main.js can be migrated: Yes")
	return true, nil
}

func getLatestWailsVersion() (*semver.Version, error) {
	checkSpinner.Start("Checking GitHub for latest Wails version")
	version, err := migrateGithub.GetLatestStableRelease()
	if err != nil {
		checkSpinner.Error("Checking GitHub for latest Wails version: Failed")
		return nil, err
	}

	checkSpinner.Successf("Checking GitHub for latest Wails version: v%s", version)
	return version.Version, nil
}

func updateWailsVersion(currentVersion, latestVersion *semver.Version) error {
	// Patch go.mod
	checkSpinner.Start("Patching go.mod")

	wailsModule := "github.com/wailsapp/wails"
	old := fmt.Sprintf("%s v%s", wailsModule, currentVersion)
	new := fmt.Sprintf("%s v%s", wailsModule, latestVersion)

	goMod = strings.Replace(goMod, old, new, -1)
	err := os.WriteFile(goModFile, []byte(goMod), 0600)
	if err != nil {
		checkSpinner.Error()
		return err
	}

	checkSpinner.Success()
	return nil
}

func deleteWailsBridge(bridgeFilename string) error {
	// Patch go.mod
	checkSpinner.Start("Delete legacy wailsbridge.js")

	err := migrateFS.RemoveFile(bridgeFilename)
	if err != nil {
		checkSpinner.Error()
		return err
	}

	checkSpinner.Success()
	return nil
}

func patchMainJS() error {
	// Patch main.js
	checkSpinner.Start("Patching main.js")

	// Patch import line
	oldImportLine := `import Bridge from "./wailsbridge";`
	newImportLine := `import * as Wails from "@wailsapp/runtime";`
	mainJSContents = strings.Replace(mainJSContents, oldImportLine, newImportLine, -1)

	// Patch Start line
	oldStartLine := `Bridge.Start`
	newStartLine := `Wails.Init`
	mainJSContents = strings.Replace(mainJSContents, oldStartLine, newStartLine, -1)

	err := os.WriteFile(mainJSFile, []byte(mainJSContents), 0600)
	if err != nil {
		checkSpinner.Error()
		return err
	}

	checkSpinner.Success()
	return nil
}

func installWailsRuntime() error {

	checkSpinner.Start("Installing @wailsapp/runtime module")

	// Change to the frontend directory
	err := os.Chdir(frontEndDir)
	if err != nil {
		checkSpinner.Error()
		return nil
	}

	// Determine package manager
	packageManager, err := migrateProjectOptions.GetNPMBinaryName()
	if err != nil {
		checkSpinner.Error()
		return nil
	}

	switch packageManager {
	case cmd.NPM:
		// npm install --save @wailsapp/runtime
		programHelper.InstallNPMPackage("@wailsapp/runtime", true)
	default:
		checkSpinner.Error()
		return fmt.Errorf("Unknown package manager")
	}

	checkSpinner.Success()
	return nil
}
