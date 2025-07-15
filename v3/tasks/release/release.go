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

func checkForChanges() (bool, string) {
	// Get the latest v3 alpha tag
	output, err := runCommand("git", "tag", "--list", "v3.0.0-alpha.*", "--sort=-version:refname")
	if err != nil {
		return true, "No previous tags found"
	}
	
	tags := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(tags) == 0 || tags[0] == "" {
		return true, "No previous tags found"
	}
	
	latestTag := tags[0]
	
	// Check for commits since the latest tag
	output, err = runCommand("git", "rev-list", latestTag+"..HEAD", "--count")
	if err != nil {
		return true, "Error checking commits: " + err.Error()
	}
	
	commitCount := strings.TrimSpace(string(output))
	if commitCount == "0" {
		return false, fmt.Sprintf("No changes since %s", latestTag)
	}
	
	// Get commit messages since the latest tag
	output, err = runCommand("git", "log", "--pretty=format:- %s", latestTag+"..HEAD")
	if err != nil {
		return true, "Error getting commit messages: " + err.Error()
	}
	
	return true, fmt.Sprintf("Found %s commits since %s:\n%s", commitCount, latestTag, strings.TrimSpace(string(output)))
}

func validateChangelogUpdate(version string) (bool, string, error) {
	// Check if we're in the right directory
	changelogPath := "docs/src/content/docs/changelog.mdx"
	if _, err := os.Stat(changelogPath); os.IsNotExist(err) {
		// Try from project root
		changelogPath = "../../../docs/src/content/docs/changelog.mdx"
	}
	
	changelogData, err := os.ReadFile(changelogPath)
	if err != nil {
		return false, "", fmt.Errorf("failed to read changelog: %v", err)
	}
	
	changelog := string(changelogData)
	
	// Check if the version exists in changelog
	versionHeader := "## " + version + " - "
	if !strings.Contains(changelog, versionHeader) {
		return false, "", fmt.Errorf("version %s not found in changelog", version)
	}
	
	// Extract the content for this version
	versionStart := strings.Index(changelog, versionHeader)
	if versionStart == -1 {
		return false, "", fmt.Errorf("version header not found")
	}
	
	// Find the next version section - look for next ## followed by a version pattern
	remainingContent := changelog[versionStart+len(versionHeader):]
	nextVersionStart := strings.Index(remainingContent, "\n## v")
	if nextVersionStart == -1 {
		// This is the last version, take everything until end
		content := changelog[versionStart:]
		return true, strings.TrimSpace(content), nil
	}
	
	// Extract content between this version and next version
	content := changelog[versionStart : versionStart+len(versionHeader)+nextVersionStart]
	return true, strings.TrimSpace(content), nil
}

func outputReleaseMetadata(version, changelog string, hasChanges bool, changesSummary string) error {
	fmt.Println("========================================")
	fmt.Println("🧪 DRY RUN MODE - TESTING RELEASE SCRIPT")
	fmt.Println("========================================")
	
	// 1. Changes detection
	fmt.Printf("1. CHANGES DETECTED: %t\n", hasChanges)
	fmt.Printf("   SUMMARY: %s\n\n", changesSummary)
	
	// 2. Changelog validation
	fmt.Println("2. CHANGELOG VALIDATION:")
	_, extractedChangelog, err := validateChangelogUpdate(version)
	if err != nil {
		fmt.Printf("   ❌ FAILED: %v\n\n", err)
		return err
	}
	fmt.Printf("   ✅ PASSED: Version %s found in changelog\n\n", version)
	
	// 3. Release notes in memory
	fmt.Println("3. RELEASE NOTES EXTRACTED TO MEMORY:")
	fmt.Printf("   LENGTH: %d characters\n", len(extractedChangelog))
	fmt.Printf("   PREVIEW (first 200 chars): %s...\n\n", extractedChangelog[:min(200, len(extractedChangelog))])
	
	// 4. Prerelease data
	fmt.Println("4. GITHUB PRERELEASE DATA:")
	fmt.Printf("   VERSION: %s\n", version)
	fmt.Printf("   TAG: %s\n", version)
	fmt.Printf("   TITLE: Wails v3 Alpha Release - %s\n", version)
	fmt.Printf("   IS_PRERELEASE: true\n")
	fmt.Printf("   IS_LATEST: false\n")
	fmt.Printf("   DRAFT: false\n\n")
	
	// Output environment variables for GitHub Actions
	fmt.Println("5. ENVIRONMENT VARIABLES FOR GITHUB ACTIONS:")
	fmt.Printf("RELEASE_VERSION=%s\n", version)
	fmt.Printf("RELEASE_TAG=%s\n", version)
	fmt.Printf("RELEASE_TITLE=Wails v3 Alpha Release - %s\n", version)
	fmt.Printf("RELEASE_IS_PRERELEASE=true\n")
	fmt.Printf("RELEASE_IS_LATEST=false\n")
	fmt.Printf("RELEASE_DRAFT=false\n")
	fmt.Printf("HAS_CHANGES=%t\n", hasChanges)
	
	// Write changelog to file for GitHub Actions
	err = os.WriteFile("release-notes.txt", []byte(extractedChangelog), 0o644)
	if err != nil {
		return fmt.Errorf("failed to write release notes: %v", err)
	}
	
	fmt.Printf("RELEASE_NOTES_FILE=release-notes.txt\n\n")
	
	fmt.Println("========================================")
	fmt.Println("✅ DRY RUN COMPLETED SUCCESSFULLY")
	fmt.Println("========================================")
	
	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
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
	fmt.Println("🧪 STARTING DRY RUN RELEASE SCRIPT TEST")
	fmt.Println("=======================================")
	
	// Step 0: Ensure we have latest git data
	fmt.Println("STEP 0: Fetching latest git data...")
	_, gitErr := runCommand("git", "fetch", "--tags", "origin")
	if gitErr != nil {
		fmt.Printf("⚠️  Warning: Failed to fetch latest tags: %v\n", gitErr)
		fmt.Println("   Continuing with local git state...")
	} else {
		fmt.Println("   ✅ Latest tags fetched successfully")
	}
	
	// Step 1: Check for changes since last release
	fmt.Println("STEP 1: Checking for changes...")
	hasChanges, changesSummary := checkForChanges()
	
	// Step 2: Check if current commit has a release tag
	hasTag, tag := hasReleaseTag()
	
	var newVersion string
	var releaseChangelog string
	
	if hasTag {
		// Current commit has a release tag - this is a nightly release scenario
		fmt.Printf("Found release tag: %s\n", tag)
		
		// Read version from version.txt
		currentVersionData, err := os.ReadFile(versionFile)
		checkError(err)
		newVersion = strings.TrimSpace(string(currentVersionData))
		
		// Extract changelog since the tag
		changelog, err := extractChangelogSinceTag(tag)
		checkError(err)
		releaseChangelog = changelog
		
		fmt.Printf("Nightly release scenario for tag: %s\n", newVersion)
		
	} else {
		// No release tag - normal release process
		fmt.Println("STEP 2: No release tag found - proceeding with normal release process")
		
		// Don't actually update files in dry run mode - just simulate
		fmt.Println("🔄 SIMULATING VERSION UPDATE...")
		
		if len(os.Args) > 1 {
			newVersion = os.Args[1]
			fmt.Printf("   Using provided version: %s\n", newVersion)
		} else {
			// Read current version and simulate increment
			currentVersionData, err := os.ReadFile(versionFile)
			checkError(err)
			currentVersion := strings.TrimSpace(string(currentVersionData))
			vsplit := strings.Split(currentVersion, ".")
			minorVersion, err := strconv.Atoi(vsplit[len(vsplit)-1])
			checkError(err)
			minorVersion++
			vsplit[len(vsplit)-1] = strconv.Itoa(minorVersion)
			newVersion = strings.Join(vsplit, ".")
			fmt.Printf("   Current version: %s\n", currentVersion)
			fmt.Printf("   Next version would be: %s\n", newVersion)
		}

		fmt.Println("🔄 SIMULATING CHANGELOG UPDATE...")
		// Simulate changelog update by checking current structure
		s.CD("../../..")
		changelogData, err := os.ReadFile("docs/src/content/docs/changelog.mdx")
		checkError(err)
		changelog := string(changelogData)
		
		// Check if changelog structure is valid
		if !strings.Contains(changelog, "## [Unreleased]") {
			fmt.Println("   ❌ ERROR: Changelog missing [Unreleased] section")
			os.Exit(1)
		}
		
		today := time.Now().Format("2006-01-02")
		fmt.Printf("   Would add version section: ## %s - %s\n", newVersion, today)
		
		// Simulate extracting unreleased content
		unreleasedChangelog, err := extractUnreleasedChangelog()
		checkError(err)
		releaseChangelog = unreleasedChangelog
		
		fmt.Printf("   ✅ Changelog structure validated\n")
		fmt.Printf("   📝 Extracted %d characters of unreleased content\n", len(releaseChangelog))
	}
	
	// Output comprehensive test results
	fmt.Println("\nSTEP 3: Generating test results...")
	err := outputReleaseMetadata(newVersion, releaseChangelog, hasChanges, changesSummary)
	if err != nil {
		fmt.Printf("❌ Test failed: %v\n", err)
		os.Exit(1)
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
