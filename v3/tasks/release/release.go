package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	versionFile   = "../../internal/version/version.txt"
	changelogFile = "../../../docs/src/content/docs/changelog.mdx"
)

var (
	unreleasedChangelogFile = "../../UNRELEASED_CHANGELOG.md"
)

func checkError(err error) {
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
}

// getUnreleasedChangelogTemplate returns the template content for UNRELEASED_CHANGELOG.md
func getUnreleasedChangelogTemplate() string {
	return `# Unreleased Changes

<!-- 
This file is used to collect changelog entries for the next v3-alpha release.
Add your changes under the appropriate sections below.

Guidelines:
- Follow the "Keep a Changelog" format (https://keepachangelog.com/)
- Write clear, concise descriptions of changes
- Include the impact on users when relevant
- Use present tense ("Add feature" not "Added feature")
- Reference issue/PR numbers when applicable

This file is automatically processed by the nightly release workflow.
After processing, the content will be moved to the main changelog and this file will be reset.
-->

## Added
<!-- New features, capabilities, or enhancements -->

## Changed
<!-- Changes in existing functionality -->

## Fixed
<!-- Bug fixes -->

## Deprecated
<!-- Soon-to-be removed features -->

## Removed
<!-- Features removed in this release -->

## Security
<!-- Security-related changes -->

---

### Example Entries:

**Added:**
- Add support for custom window icons in application options
- Add new ` + "`SetWindowIcon()`" + ` method to runtime API (#1234)

**Changed:**
- Update minimum Go version requirement to 1.21
- Improve error messages for invalid configuration files

**Fixed:**
- Fix memory leak in event system during window close operations (#5678)
- Fix crash when using context menus on Linux with Wayland

**Security:**
- Update dependencies to address CVE-2024-12345 in third-party library
`
}

// clearUnreleasedChangelog clears the UNRELEASED_CHANGELOG.md file and resets it with the template
func clearUnreleasedChangelog() error {
	template := getUnreleasedChangelogTemplate()

	// Write the template back to the file
	err := os.WriteFile(unreleasedChangelogFile, []byte(template), 0o644)
	if err != nil {
		return fmt.Errorf("failed to reset UNRELEASED_CHANGELOG.md: %w", err)
	}

	fmt.Printf("Successfully reset %s with template content\n", unreleasedChangelogFile)
	return nil
}

// extractChangelogContent extracts the actual changelog content from UNRELEASED_CHANGELOG.md
// It returns the content between the section headers and the example section
func extractChangelogContent() (string, error) {
	content, err := os.ReadFile(unreleasedChangelogFile)
	if err != nil {
		return "", fmt.Errorf("failed to read %s: %w", unreleasedChangelogFile, err)
	}

	contentStr := string(content)
	lines := strings.Split(contentStr, "\n")

	var result []string
	var inExampleSection bool
	var inCommentBlock bool
	var hasActualContent bool
	var currentSection string

	for i, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		// Track comment blocks (handle multi-line comments)
		if strings.Contains(line, "<!--") {
			inCommentBlock = true
			// Check if comment ends on same line
			if strings.Contains(line, "-->") {
				inCommentBlock = false
			}
			continue
		}
		if inCommentBlock {
			if strings.Contains(line, "-->") {
				inCommentBlock = false
			}
			continue
		}

		// Skip the main title
		if strings.HasPrefix(trimmedLine, "# Unreleased Changes") {
			continue
		}

		// Check if we're entering the example section
		if strings.HasPrefix(trimmedLine, "---") || strings.HasPrefix(trimmedLine, "### Example Entries") {
			inExampleSection = true
			continue
		}

		// Skip example section content
		if inExampleSection {
			continue
		}

		// Handle section headers
		if strings.HasPrefix(trimmedLine, "##") {
			currentSection = trimmedLine
			// Only include section headers that have content after them
			// We'll add it later if we find content
			continue
		}

		// Handle bullet points
		if strings.HasPrefix(trimmedLine, "-") || strings.HasPrefix(trimmedLine, "*") {
			// Check if this is actual content (not empty)
			content := strings.TrimSpace(trimmedLine[1:])
			if content != "" {
				// If this is the first content in a section, add the section header first
				if currentSection != "" {
					// Only add empty line if this isn't the first section
					if len(result) > 0 {
						result = append(result, "")
					}
					result = append(result, currentSection)
					currentSection = "" // Reset so we don't add it again
				}
				result = append(result, line)
				hasActualContent = true
			}
		} else if trimmedLine != "" && !strings.HasPrefix(trimmedLine, "<!--") {
			// Include other non-empty, non-comment lines that aren't section headers
			if !strings.HasPrefix(trimmedLine, "##") {
				// Check if next line exists and is not a comment placeholder
				if i+1 < len(lines) {
					nextLine := strings.TrimSpace(lines[i+1])
					if !strings.HasPrefix(nextLine, "<!--") {
						result = append(result, line)
					}
				}
			}
		}
	}

	if !hasActualContent {
		return "", nil
	}

	// Clean up result - remove any trailing empty lines
	for len(result) > 0 && strings.TrimSpace(result[len(result)-1]) == "" {
		result = result[:len(result)-1]
	}

	return strings.Join(result, "\n"), nil
}

// hasUnreleasedContent checks if UNRELEASED_CHANGELOG.md has actual content beyond the template
func hasUnreleasedContent() (bool, error) {
	content, err := extractChangelogContent()
	if err != nil {
		return false, err
	}
	return content != "", nil
}

// safeFileOperation performs a file operation with backup and rollback capability
func safeFileOperation(filePath string, operation func() error) error {
	// Create backup if file exists
	var backupPath string
	var hasBackup bool

	if _, err := os.Stat(filePath); err == nil {
		backupPath = filePath + ".backup"
		if err := copyFile(filePath, backupPath); err != nil {
			return fmt.Errorf("failed to create backup of %s: %w", filePath, err)
		}
		hasBackup = true
		defer func() {
			// Clean up backup file on success
			if hasBackup {
				_ = os.Remove(backupPath)
			}
		}()
	}

	// Perform the operation
	if err := operation(); err != nil {
		// Rollback if we have a backup
		if hasBackup {
			if rollbackErr := copyFile(backupPath, filePath); rollbackErr != nil {
				return fmt.Errorf("operation failed and rollback failed: %w (rollback error: %v)", err, rollbackErr)
			}
		}
		return err
	}

	return nil
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0o644)
}

// updateVersion increments the version number properly handling semantic versioning
// Examples:
// v3.0.0-alpha.12 -> v3.0.0-alpha.13
// v3.0.0 -> v3.0.1
// v3.0.0-beta.1 -> v3.0.0-beta.2
func updateVersion() string {
	currentVersionData, err := os.ReadFile(versionFile)
	checkError(err)
	currentVersion := strings.TrimSpace(string(currentVersionData))

	// Check if it has a pre-release suffix (e.g., -alpha.12, -beta.1)
	if strings.Contains(currentVersion, "-") {
		// Split on the dash to separate version and pre-release
		parts := strings.SplitN(currentVersion, "-", 2)
		baseVersion := parts[0]
		preRelease := parts[1]

		// Check if pre-release has a numeric suffix (e.g., alpha.12)
		lastDotIndex := strings.LastIndex(preRelease, ".")
		if lastDotIndex != -1 {
			preReleaseTag := preRelease[:lastDotIndex]
			numberStr := preRelease[lastDotIndex+1:]

			// Try to parse the number
			if number, err := strconv.Atoi(numberStr); err == nil {
				// Increment the pre-release number
				number++
				newVersion := fmt.Sprintf("%s-%s.%d", baseVersion, preReleaseTag, number)
				err = os.WriteFile(versionFile, []byte(newVersion), 0o755)
				checkError(err)
				return newVersion
			}
		}

		// If we can't parse the pre-release format, just increment patch version
		// and remove pre-release suffix
		return incrementPatchVersion(baseVersion)
	}

	// No pre-release suffix, just increment patch version
	return incrementPatchVersion(currentVersion)
}

// incrementPatchVersion increments the patch version of a semantic version
// e.g., v3.0.0 -> v3.0.1
func incrementPatchVersion(version string) string {
	// Remove 'v' prefix if present
	versionWithoutV := strings.TrimPrefix(version, "v")

	// Split into major.minor.patch
	parts := strings.Split(versionWithoutV, ".")
	if len(parts) != 3 {
		// Not a valid semver, return as-is
		fmt.Printf("Warning: Invalid semantic version format: %s\n", version)
		return version
	}

	// Parse patch version
	patch, err := strconv.Atoi(parts[2])
	if err != nil {
		fmt.Printf("Warning: Could not parse patch version: %s\n", parts[2])
		return version
	}

	// Increment patch
	patch++

	// Reconstruct version
	newVersion := fmt.Sprintf("v%s.%s.%d", parts[0], parts[1], patch)
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

	// Check for --check-only flag
	if len(os.Args) > 1 && os.Args[1] == "--check-only" {
		// Just check if there's unreleased content and exit
		changelogContent, err := extractChangelogContent()
		if err != nil {
			fmt.Printf("Error: Failed to extract unreleased changelog content: %v\n", err)
			os.Exit(1)
		}
		if changelogContent == "" {
			fmt.Println("No unreleased changelog content found.")
			os.Exit(1)
		}
		fmt.Println("Found unreleased changelog content.")
		os.Exit(0)
	}

	// Check for --extract-changelog flag
	if len(os.Args) > 1 && os.Args[1] == "--extract-changelog" {
		// Extract and output changelog content for release notes
		changelogContent, err := extractChangelogContent()
		if err != nil {
			fmt.Printf("Error: Failed to extract unreleased changelog content: %v\n", err)
			os.Exit(1)
		}
		if changelogContent == "" {
			fmt.Println("No changelog content found.")
			os.Exit(1)
		}
		fmt.Print(changelogContent)
		os.Exit(0)
	}

	// Check for --reset-changelog flag
	if len(os.Args) > 1 && os.Args[1] == "--reset-changelog" {
		// Reset the changelog to the template
		err := clearUnreleasedChangelog()
		if err != nil {
			fmt.Printf("Error: Failed to reset changelog: %v\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	// Check for --create-release-notes flag
	if len(os.Args) > 1 && os.Args[1] == "--create-release-notes" {
		// Extract changelog content and create release_notes.md
		changelogContent, err := extractChangelogContent()
		if err != nil {
			fmt.Printf("Error: Failed to extract unreleased changelog content: %v\n", err)
			os.Exit(1)
		}
		if changelogContent == "" {
			fmt.Printf("Error: No changelog content found in UNRELEASED_CHANGELOG.md\n")
			os.Exit(1)
		}

		// Create release_notes.md file
		releaseNotesPath := "../../release_notes.md"
		if len(os.Args) > 2 {
			releaseNotesPath = os.Args[2]
		}

		err = os.WriteFile(releaseNotesPath, []byte(changelogContent), 0o644)
		if err != nil {
			fmt.Printf("Error: Failed to write release notes to %s: %v\n", releaseNotesPath, err)
			os.Exit(1)
		}

		fmt.Printf("Successfully created release notes at %s\n", releaseNotesPath)
		os.Exit(0)
	}

	// Extract changelog content first
	changelogContent, err := extractChangelogContent()
	if err != nil {
		fmt.Printf("Warning: Failed to extract unreleased changelog content: %v\n", err)
		return
	}
	if changelogContent == "" {
		fmt.Println("UNRELEASED_CHANGELOG.md is empty. Skipping changelog processing.")
		return
	}

	var newVersion string
	if len(os.Args) > 1 {
		newVersion = os.Args[1]
		//currentVersion, err := os.ReadFile(versionFile)
		//checkError(err)
		err := os.WriteFile(versionFile, []byte(newVersion), 0o755)
		checkError(err)
		//isPointRelease = IsPointRelease(string(currentVersion), newVersion)
	} else {
		newVersion = updateVersion()
	}

	// Read in the main changelog
	changelogData, err := os.ReadFile(changelogFile)
	checkError(err)
	changelog := string(changelogData)

	// Split on the line that has `## [Unreleased]`
	changelogSplit := strings.Split(changelog, "## [Unreleased]")
	if len(changelogSplit) != 2 {
		fmt.Printf("Error: Could not find '## [Unreleased]' section in changelog\n")
		os.Exit(1)
	}

	// Get today's date in YYYY-MM-DD format
	today := time.Now().Format("2006-01-02")

	// Create the new changelog with the extracted content
	newChangelog := changelogSplit[0] + "## [Unreleased]\n\n## " + newVersion + " - " + today + "\n\n" + changelogContent + changelogSplit[1]

	// Write the changelog back
	err = safeFileOperation(changelogFile, func() error {
		return os.WriteFile(changelogFile, []byte(newChangelog), 0o644)
	})
	checkError(err)

	// Clear UNRELEASED_CHANGELOG.md after successful changelog update
	fmt.Printf("Changelog updated successfully. Clearing %s...\n", unreleasedChangelogFile)
	err = safeFileOperation(unreleasedChangelogFile, func() error {
		return clearUnreleasedChangelog()
	})
	if err != nil {
		fmt.Printf("Error: Failed to clear %s: %v\n", unreleasedChangelogFile, err)
		os.Exit(1)
	}

	fmt.Printf("Release %s processed successfully!\n", newVersion)
}
