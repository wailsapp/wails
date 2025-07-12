package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/wailsapp/wails/v3/internal/s"
)

const versionFile = "../../internal/version/version.txt"

func checkError(err error) {
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
}

func runCommand(name string, arg ...string) ([]byte, error) {
	cmd := exec.Command(name, arg...)
	return cmd.Output()
}

func hasReleaseTag() (bool, string) {
	output, err := runCommand("git", "describe", "--tags", "--exact-match", "HEAD")
	if err != nil {
		return false, ""
	}

	tag := strings.TrimSpace(string(output))
	matched, _ := regexp.MatchString(`^v3\.0\.0-alpha\.\d+(\.\d+)*$`, tag)
	return matched, tag
}

func extractChangelogSinceTag(tag string) (string, error) {
	output, err := runCommand("git", "log", "--pretty=format:- %s", tag+"..HEAD")
	if err != nil {
		return "", err
	}

	changelog := strings.TrimSpace(string(output))
	if changelog == "" {
		return "No changes since " + tag, nil
	}

	return changelog, nil
}

func extractUnreleasedChangelog() (string, error) {
	// This function assumes we're in the project root
	changelogData, err := os.ReadFile("docs/src/content/docs/changelog.mdx")
	if err != nil {
		return "", err
	}

	changelog := string(changelogData)

	// Find the [Unreleased] section
	unreleasedStart := strings.Index(changelog, "## [Unreleased]")
	if unreleasedStart == -1 {
		return "No unreleased changes found", nil
	}

	// Find the next version section
	nextVersionStart := strings.Index(changelog[unreleasedStart+len("## [Unreleased]"):], "## ")
	if nextVersionStart == -1 {
		// No next version, take everything after [Unreleased]
		content := changelog[unreleasedStart+len("## [Unreleased]"):]
		return strings.TrimSpace(content), nil
	}

	// Extract content between [Unreleased] and next version
	content := changelog[unreleasedStart+len("## [Unreleased]") : unreleasedStart+len("## [Unreleased]")+nextVersionStart]
	return strings.TrimSpace(content), nil
}

func outputReleaseMetadata(version, changelog string) error {
	// Output release metadata for GitHub Actions to consume
	fmt.Printf("RELEASE_VERSION=%s\n", version)
	fmt.Printf("RELEASE_TAG=%s\n", version)
	fmt.Printf("RELEASE_TITLE=%s\n", version)

	// Write changelog to file for GitHub Actions
	err := os.WriteFile("release-notes.txt", []byte(changelog), 0o644)
	if err != nil {
		return fmt.Errorf("failed to write release notes: %v", err)
	}

	fmt.Println("RELEASE_NOTES_FILE=release-notes.txt")
	return nil
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

//func runCommand(name string, arg ...string) {
//	cmd := exec.Command(name, arg...)
//	cmd.Stdout = os.Stdout
//	cmd.Stderr = os.Stderr
//	err := cmd.Run()
//	checkError(err)
//}

//func IsPointRelease(currentVersion string, newVersion string) bool {
//	// The first n-1 parts of the version should be the same
//	if currentVersion[:len(currentVersion)-2] != newVersion[:len(newVersion)-2] {
//		return false
//	}
//	// split on the last dot in the string
//	currentVersionSplit := strings.Split(currentVersion, ".")
//	newVersionSplit := strings.Split(newVersion, ".")
//	// if the last part of the version is the same, it's a point release
//	currentMinor := lo.Must(strconv.Atoi(currentVersionSplit[len(currentVersionSplit)-1]))
//	newMinor := lo.Must(strconv.Atoi(newVersionSplit[len(newVersionSplit)-1]))
//	return newMinor == currentMinor+1
//}

func main() {
	// Check if current commit has a release tag
	hasTag, tag := hasReleaseTag()

	var newVersion string
	var releaseChangelog string

	if hasTag {
		// Current commit has a release tag - this is a nightly release
		fmt.Printf("Found release tag: %s\n", tag)

		// Read version from version.txt
		currentVersionData, err := os.ReadFile(versionFile)
		checkError(err)
		newVersion = strings.TrimSpace(string(currentVersionData))

		// Extract changelog since the tag
		changelog, err := extractChangelogSinceTag(tag)
		checkError(err)
		releaseChangelog = changelog

		fmt.Printf("Creating GitHub release for existing tag: %s\n", newVersion)

	} else {
		// No release tag - normal release process
		fmt.Println("No release tag found - proceeding with normal release")

		if len(os.Args) > 1 {
			newVersion = os.Args[1]
			err := os.WriteFile(versionFile, []byte(newVersion), 0o755)
			checkError(err)
		} else {
			newVersion = updateVersion()
		}

		// Update ChangeLog
		s.CD("../../..")
		changelogData, err := os.ReadFile("docs/src/content/docs/changelog.mdx")
		checkError(err)
		changelog := string(changelogData)
		// Split on the line that has `## [Unreleased]`
		changelogSplit := strings.Split(changelog, "## [Unreleased]")
		// Get today's date in YYYY-MM-DD format
		today := time.Now().Format("2006-01-02")
		// Add the new version to the top of the changelog
		newChangelog := changelogSplit[0] + "## [Unreleased]\n\n## " + newVersion + " - " + today + changelogSplit[1]
		// Write the changelog back
		err = os.WriteFile("docs/src/content/docs/changelog.mdx", []byte(newChangelog), 0o755)
		checkError(err)

		// Extract unreleased changelog for GitHub release
		unreleasedChangelog, err := extractUnreleasedChangelog()
		checkError(err)
		releaseChangelog = unreleasedChangelog

		fmt.Printf("Updated version to: %s\n", newVersion)
		fmt.Println("Updated changelog")
	}

	// Output release metadata for GitHub Actions
	fmt.Printf("Preparing release metadata for version: %s\n", newVersion)
	err := outputReleaseMetadata(newVersion, releaseChangelog)
	if err != nil {
		fmt.Printf("Failed to output release metadata: %v\n", err)
		os.Exit(1)
	} else {
		fmt.Println("Release metadata prepared successfully")
	}

	// TODO: Documentation Versioning and Translations

	//if !isPointRelease {
	//	runCommand("npx", "-y", "pnpm", "install")
	//
	//	s.ECHO("Generating new Docs for version: " + newVersion)
	//
	//	runCommand("npx", "pnpm", "run", "docusaurus", "docs:version", newVersion)
	//
	//	runCommand("npx", "pnpm", "run", "write-translations")
	//
	//	// Load the version list/*
	//	versionsData, err := os.ReadFile("versions.json")
	//	checkError(err)
	//	var versions []string
	//	err = json.Unmarshal(versionsData, &versions)
	//	checkError(err)
	//	oldestVersion := versions[len(versions)-1]
	//	s.ECHO(oldestVersion)
	//	versions = versions[0 : len(versions)-1]
	//	newVersions, err := json.Marshal(&versions)
	//	checkError(err)
	//	err = os.WriteFile("versions.json", newVersions, 0o755)
	//	checkError(err)
	//
	//	s.ECHO("Removing old version: " + oldestVersion)
	//	s.CD("versioned_docs")
	//	s.RMDIR("version-" + oldestVersion)
	//	s.CD("../versioned_sidebars")
	//	s.RM("version-" + oldestVersion + "-sidebars.json")
	//	s.CD("..")
	//
	//	runCommand("npx", "pnpm", "run", "build")
	//}
}
