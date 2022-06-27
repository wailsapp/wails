package main

import (
	"encoding/json"
	"github.com/wailsapp/wails/v2/internal/s"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const versionFile = "../../cmd/wails/internal/version.txt"

func checkError(err error) {
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
}

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

func main() {
	newVersion := updateVersion()
	s.CD("../../../website")
	s.ECHO("Generating new Docs for version: " + newVersion)
	cmd := exec.Command("npm", "run", "docusaurus", "docs:version", newVersion)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	checkError(err)

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
	err = os.WriteFile("versions.json", newVersions, 0755)
	checkError(err)

	s.ECHO("Removing old version: " + oldestVersion)
	s.CD("versioned_docs")
	s.RMDIR("version-" + oldestVersion)
	s.CD("../versioned_sidebars")
	s.RM("version-" + oldestVersion + "-sidebars.json")

}
