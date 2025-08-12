package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/samber/lo"

	"github.com/wailsapp/wails/v2/internal/s"
)

const versionFile = "../../cmd/wails/internal/version.txt"

func checkError(err error) {
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
}

// TODO:This can be replaced with "https://github.com/coreos/go-semver/blob/main/semver/semver.go"
func updateVersion() string {
	currentVersionData, err := os.ReadFile(versionFile)
	checkError(err)
	currentVersion := string(currentVersionData)
	vsplit := strings.Split(currentVersion, ".")
	minorVersion, err := strconv.Atoi(vsplit[len(vsplit)-1])
	checkError(err)
	minorVersion++
	vsplit[len(vsplit)-1] = strconv.Itoa(minorVersion)
	newVersion := strings.Join(vsplit, ".")
	err = os.WriteFile(versionFile, []byte(newVersion), 0o755)
	checkError(err)
	return newVersion
}

func runCommand(name string, arg ...string) {
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	checkError(err)
}

func IsPointRelease(currentVersion string, newVersion string) bool {
	// The first n-1 parts of the version should be the same
	if currentVersion[:len(currentVersion)-2] != newVersion[:len(newVersion)-2] {
		return false
	}
	// split on the last dot in the string
	currentVersionSplit := strings.Split(currentVersion, ".")
	newVersionSplit := strings.Split(newVersion, ".")
	// compare the
	// if the last part of the version is the same, it's a point release
	currentMinor := lo.Must(strconv.Atoi(currentVersionSplit[len(currentVersionSplit)-1]))
	newMinor := lo.Must(strconv.Atoi(newVersionSplit[len(newVersionSplit)-1]))
	return newMinor == currentMinor+1
}

func main() {
	var newVersion string
	var isPointRelease bool
	if len(os.Args) > 1 {
		newVersion = os.Args[1]
		currentVersion, err := os.ReadFile(versionFile)
		checkError(err)
		err = os.WriteFile(versionFile, []byte(newVersion), 0o755)
		checkError(err)
		isPointRelease = IsPointRelease(string(currentVersion), newVersion)
	} else {
		newVersion = updateVersion()
	}

	// Update ChangeLog
	s.CD("../../../website")

	// Read in `src/pages/changelog.mdx`
	changelogData, err := os.ReadFile("src/pages/changelog.mdx")
	checkError(err)
	changelog := string(changelogData)
	// Split on the line that has `## [Unreleased]`
	changelogSplit := strings.Split(changelog, "## [Unreleased]")
	// Get today's date in YYYY-MM-DD format
	today := time.Now().Format("2006-01-02")
	// Add the new version to the top of the changelog
	newChangelog := changelogSplit[0] + "## [Unreleased]\n\n## " + newVersion + " - " + today + changelogSplit[1]
	// Write the changelog back
	err = os.WriteFile("src/pages/changelog.mdx", []byte(newChangelog), 0o755)
	checkError(err)

	if !isPointRelease {
		runCommand("npx", "-y", "pnpm", "install")

		s.ECHO("Generating new Docs for version: " + newVersion)

		runCommand("npx", "pnpm", "run", "docusaurus", "docs:version", newVersion)

		runCommand("npx", "pnpm", "run", "write-translations")

		// Load the version list/*
		versionsData, err := os.ReadFile("versions.json")
		checkError(err)
		var versions []string
		err = json.Unmarshal(versionsData, &versions)
		checkError(err)
		oldestVersion := versions[len(versions)-1]
		s.ECHO(oldestVersion)
		versions = versions[0 : len(versions)-1]
		newVersions, err := json.Marshal(&versions)
		checkError(err)
		err = os.WriteFile("versions.json", newVersions, 0o755)
		checkError(err)

		s.ECHO("Removing old version: " + oldestVersion)
		s.CD("versioned_docs")
		s.RMDIR("version-" + oldestVersion)
		s.CD("../versioned_sidebars")
		s.RM("version-" + oldestVersion + "-sidebars.json")
		s.CD("..")

		runCommand("npx", "pnpm", "run", "build")
	}
}
