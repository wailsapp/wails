package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/wailsapp/wails/v2/internal/s"
)

const versionFile = "../../cmd/wails/internal/version.txt"

func checkError(err error) {
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
}

// updateVersion increments the minor segment and resets patch to 0.
// v2.11.0 -> v2.12.0, v2.12.3 -> v2.13.0
func updateVersion() string {
	currentVersionData, err := os.ReadFile(versionFile)
	checkError(err)
	currentVersion := string(currentVersionData)

	// Parse: "v2.11.0" -> ["v2", "11", "0"]
	vsplit := strings.Split(currentVersion, ".")
	if len(vsplit) < 3 {
		fmt.Printf("unexpected version format: %s\n", currentVersion)
		os.Exit(1)
	}

	minorVersion, err := strconv.Atoi(vsplit[1])
	checkError(err)
	minorVersion++
	vsplit[1] = strconv.Itoa(minorVersion)
	vsplit[2] = "0" // reset patch

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

func updateDocs(newVersion string) {
	runCommand("npx", "-y", "pnpm", "install", "--no-frozen-lockfile")

	s.ECHO("Generating new Docs for version: " + newVersion)

	runCommand("npx", "pnpm", "run", "docusaurus", "docs:version", newVersion)

	runCommand("npx", "pnpm", "run", "write-translations")

	// Load the version list
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

func main() {
	args := os.Args[1:]
	doUpdateDocs := slices.Contains(args, "--update-docs")
	// Filter out flags to get positional args
	var positional []string
	for _, arg := range args {
		if !strings.HasPrefix(arg, "--") {
			positional = append(positional, arg)
		}
	}

	var newVersion string
	if len(positional) > 0 {
		// Explicit version provided
		newVersion = positional[0]
		err := os.WriteFile(versionFile, []byte(newVersion), 0o755)
		checkError(err)
	} else {
		// No version provided: minor bump
		newVersion = updateVersion()
	}

	fmt.Printf("Releasing %s\n", newVersion)

	// Update ChangeLog
	s.CD("../../../website")

	changelogData, err := os.ReadFile("src/pages/changelog.mdx")
	checkError(err)
	changelog := string(changelogData)
	changelogSplit := strings.Split(changelog, "## [Unreleased]")
	today := time.Now().Format("2006-01-02")
	newChangelog := changelogSplit[0] + "## [Unreleased]\n\n## " + newVersion + " - " + today + changelogSplit[1]
	err = os.WriteFile("src/pages/changelog.mdx", []byte(newChangelog), 0o755)
	checkError(err)

	if doUpdateDocs {
		updateDocs(newVersion)
	}
}
