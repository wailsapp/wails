package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"strconv"
	"strings"

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
	err = os.WriteFile(versionFile, []byte(newVersion), 0755)
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

func main() {
	var newVersion string
	if len(os.Args) > 1 {
		newVersion = os.Args[1]
		err := os.WriteFile(versionFile, []byte(newVersion), 0755)
		checkError(err)
	} else {
		newVersion = updateVersion()
	}

	s.CD("../../../website")
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
	err = os.WriteFile("versions.json", newVersions, 0755)
	checkError(err)

	s.ECHO("Removing old version: " + oldestVersion)
	s.CD("versioned_docs")
	s.RMDIR("version-" + oldestVersion)
	s.CD("../versioned_sidebars")
	s.RM("version-" + oldestVersion + "-sidebars.json")
	s.CD("..")

	runCommand("npx", "pnpm", "run", "build")
}
