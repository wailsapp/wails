package cmd

import (
	"fmt"
	"path/filepath"
	"regexp"

	"github.com/Masterminds/semver"
)

func GetWailsVersion() (*semver.Version, error) {
	var FS = NewFSHelper()
	var result *semver.Version

	// Load file
	var err error
	goModFile, err := filepath.Abs(filepath.Join(".", "go.mod"))
	if err != nil {
		return nil, fmt.Errorf("Unable to load go.mod at %s", goModFile)
	}
	goMod, err := FS.LoadAsString(goModFile)
	if err != nil {
		return nil, fmt.Errorf("Unable to load go.mod")
	}

	// Find wails version
	versionRegexp := regexp.MustCompile(`.*github.com/wailsapp/wails.*(v\d+.\d+.\d+(?:-pre\d+)?)`)
	versions := versionRegexp.FindStringSubmatch(goMod)

	if len(versions) != 2 {
		return nil, fmt.Errorf("Unable to determine Wails version")
	}

	version := versions[1]
	result, err = semver.NewVersion(version)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse Wails version: %s", version)
	}
	return result, nil

}

func GetCurrentVersion() (*semver.Version, error) {
	result, err := semver.NewVersion(Version)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse Wails version: %s", Version)
	}
	return result, nil
}

func GoModOutOfSync() (bool, error) {
	gomodversion, err := GetWailsVersion()
	if err != nil {
		return true, err
	}
	currentVersion, err := GetCurrentVersion()
	if err != nil {
		return true, err
	}
	result := !currentVersion.Equal(gomodversion)
	return result, nil
}

func UpdateGoModVersion() error {
	currentVersion, err := GetCurrentVersion()
	if err != nil {
		return err
	}
	currentVersionString := currentVersion.String()

	requireLine := "-require=github.com/wailsapp/wails@v" + currentVersionString

	// Issue: go mod edit -require=github.com/wailsapp/wails@1.0.2-pre5
	helper := NewProgramHelper()
	command := []string{"go", "mod", "edit", requireLine}
	return helper.RunCommandArray(command)

}
