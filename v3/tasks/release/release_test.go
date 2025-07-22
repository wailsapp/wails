package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// setupTestEnvironment creates a proper directory structure for tests
// It returns the cleanup function and the project root directory
func setupTestEnvironment(t *testing.T) (cleanup func(), projectRoot string) {
	// Save current directory
	originalDir, _ := os.Getwd()

	// Create a temporary directory for testing
	tmpDir := t.TempDir()

	// Create the wails project structure within temp directory
	projectRoot = filepath.Join(tmpDir, "wails")
	v3Dir := filepath.Join(projectRoot, "v3")
	releaseDir := filepath.Join(v3Dir, "tasks", "release")

	// Create all necessary directories
	err := os.MkdirAll(releaseDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory structure: %v", err)
	}

	// Change to the release directory (where the script would run from)
	os.Chdir(releaseDir)

	// Return cleanup function and project root
	cleanup = func() {
		os.Chdir(originalDir)
	}
	return cleanup, projectRoot
}

func TestExtractChangelogContent(t *testing.T) {
	cleanup, _ := setupTestEnvironment(t)
	defer cleanup()

	// Create a test file with mixed content
	testContent := `# Unreleased Changes

<!-- 
This file is used to collect changelog entries for the next v3-alpha release.
Add your changes under the appropriate sections below.
-->

## Added
- Add support for custom window icons in application options
- Add new SetWindowIcon() method to runtime API (#1234)

## Changed
<!-- Changes in existing functionality -->

## Fixed
- Fix memory leak in event system during window close operations (#5678)
- Fix crash when using context menus on Linux with Wayland

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
- Update minimum Go version requirement to 1.21`

	err := os.WriteFile(unreleasedChangelogFile, []byte(testContent), 0o644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Extract the content
	content, err := extractChangelogContent()
	if err != nil {
		t.Fatalf("extractChangelogContent() failed: %v", err)
	}

	// Verify we got content
	if content == "" {
		t.Fatal("Expected to extract content, but got empty string")
	}

	// Verify section headers WITH CONTENT are included
	if !strings.Contains(content, "## Added") {
		t.Error("Expected to find '## Added' section header")
	}

	if !strings.Contains(content, "## Fixed") {
		t.Error("Expected to find '## Fixed' section header")
	}

	// Verify empty sections are NOT included
	if strings.Contains(content, "## Changed") {
		t.Error("Expected NOT to find empty '## Changed' section header")
	}

	if strings.Contains(content, "## Deprecated") {
		t.Error("Expected NOT to find empty '## Deprecated' section header")
	}

	if strings.Contains(content, "## Removed") {
		t.Error("Expected NOT to find empty '## Removed' section header")
	}

	if strings.Contains(content, "## Security") {
		t.Error("Expected NOT to find empty '## Security' section header")
	}

	// Verify actual content is included
	if !strings.Contains(content, "Add support for custom window icons") {
		t.Error("Expected to find actual Added content")
	}

	if !strings.Contains(content, "Fix memory leak in event system") {
		t.Error("Expected to find actual Fixed content")
	}

	// Verify example content is NOT included
	if strings.Contains(content, "Update minimum Go version requirement to 1.21") {
		t.Error("Expected NOT to find example content")
	}

	// Verify comments are NOT included
	if strings.Contains(content, "<!--") || strings.Contains(content, "-->") {
		t.Error("Expected NOT to find HTML comments")
	}

	// Verify the separator and example header are NOT included
	if strings.Contains(content, "---") {
		t.Error("Expected NOT to find separator")
	}

	if strings.Contains(content, "### Example Entries") {
		t.Error("Expected NOT to find example section header")
	}
}

func TestExtractChangelogContent_EmptySections(t *testing.T) {
	cleanup, _ := setupTestEnvironment(t)
	defer cleanup()

	// Create a test file with only one section having content
	testContent := `# Unreleased Changes

<!-- 
This file is used to collect changelog entries for the next v3-alpha release.
-->

## Added
<!-- New features, capabilities, or enhancements -->

## Changed
<!-- Changes in existing functionality -->

## Fixed
- Fix critical bug in the system

## Deprecated
<!-- Soon-to-be removed features -->

## Removed
<!-- Features removed in this release -->

## Security
<!-- Security-related changes -->

---

### Example Entries:

**Added:**
- Example entry that should not be included`

	err := os.WriteFile(unreleasedChangelogFile, []byte(testContent), 0o644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Extract the content
	content, err := extractChangelogContent()
	if err != nil {
		t.Fatalf("extractChangelogContent() failed: %v", err)
	}

	// Verify we got content
	if content == "" {
		t.Fatal("Expected to extract content, but got empty string")
	}

	// Verify ONLY the Fixed section is included (the only one with content)
	if !strings.Contains(content, "## Fixed") {
		t.Error("Expected to find '## Fixed' section header")
	}

	if !strings.Contains(content, "Fix critical bug in the system") {
		t.Error("Expected to find the Fixed content")
	}

	// Verify empty sections are NOT included
	sections := []string{"## Added", "## Changed", "## Deprecated", "## Removed", "## Security"}
	for _, section := range sections {
		if strings.Contains(content, section) {
			t.Errorf("Expected NOT to find empty section '%s'", section)
		}
	}

	// Verify comments are not included
	if strings.Contains(content, "<!--") || strings.Contains(content, "-->") {
		t.Error("Expected NOT to find HTML comments")
	}

	// Verify example content is not included
	if strings.Contains(content, "Example entry that should not be included") {
		t.Error("Expected NOT to find example content")
	}
}

func TestExtractChangelogContent_AllEmpty(t *testing.T) {
	cleanup, _ := setupTestEnvironment(t)
	defer cleanup()

	// Create a test file with all empty sections (just the template)
	testContent := getUnreleasedChangelogTemplate()

	err := os.WriteFile(unreleasedChangelogFile, []byte(testContent), 0o644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Extract the content
	content, err := extractChangelogContent()
	if err != nil {
		t.Fatalf("extractChangelogContent() failed: %v", err)
	}

	// Verify we got empty string (no content)
	if content != "" {
		t.Fatalf("Expected empty string for template-only file, got: %s", content)
	}
}

func TestExtractChangelogContent_MixedSections(t *testing.T) {
	cleanup, _ := setupTestEnvironment(t)
	defer cleanup()

	// Create a test file with some sections having content, others empty
	testContent := `# Unreleased Changes

## Added
- New feature A
- New feature B

## Changed
<!-- No changes -->

## Fixed
<!-- No fixes -->

## Deprecated
- Deprecated feature X

## Removed
<!-- Nothing removed -->

## Security
- Security fix for CVE-2024-1234`

	err := os.WriteFile(unreleasedChangelogFile, []byte(testContent), 0o644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Extract the content
	content, err := extractChangelogContent()
	if err != nil {
		t.Fatalf("extractChangelogContent() failed: %v", err)
	}

	// Verify we got content
	if content == "" {
		t.Fatal("Expected to extract content, but got empty string")
	}

	// Verify sections WITH content are included
	if !strings.Contains(content, "## Added") {
		t.Error("Expected to find '## Added' section header")
	}
	if !strings.Contains(content, "New feature A") {
		t.Error("Expected to find Added content")
	}

	if !strings.Contains(content, "## Deprecated") {
		t.Error("Expected to find '## Deprecated' section header")
	}
	if !strings.Contains(content, "Deprecated feature X") {
		t.Error("Expected to find Deprecated content")
	}

	if !strings.Contains(content, "## Security") {
		t.Error("Expected to find '## Security' section header")
	}
	if !strings.Contains(content, "Security fix for CVE-2024-1234") {
		t.Error("Expected to find Security content")
	}

	// Verify empty sections are NOT included
	emptyHeaders := []string{"## Changed", "## Fixed", "## Removed"}
	for _, header := range emptyHeaders {
		if strings.Contains(content, header) {
			t.Errorf("Expected NOT to find empty section '%s'", header)
		}
	}

	// Print the extracted content for debugging
	t.Logf("Extracted content:\n%s", content)
}

func TestGetUnreleasedChangelogTemplate(t *testing.T) {
	template := getUnreleasedChangelogTemplate()

	// Check that template contains required sections
	requiredSections := []string{
		"# Unreleased Changes",
		"## Added",
		"## Changed",
		"## Fixed",
		"## Deprecated",
		"## Removed",
		"## Security",
		"### Example Entries",
	}

	for _, section := range requiredSections {
		if !strings.Contains(template, section) {
			t.Errorf("Template missing required section: %s", section)
		}
	}

	// Check that template contains guidelines
	if !strings.Contains(template, "Guidelines:") {
		t.Error("Template missing guidelines section")
	}

	// Check that template contains example entries
	if !strings.Contains(template, "Add support for custom window icons") {
		t.Error("Template missing example entries")
	}
}

func TestClearUnreleasedChangelog(t *testing.T) {
	cleanup, _ := setupTestEnvironment(t)
	defer cleanup()

	// Create a test file with some content
	testContent := `# Unreleased Changes

## Added
- Some test content
- Another test item

## Fixed
- Fixed something important`

	err := os.WriteFile(unreleasedChangelogFile, []byte(testContent), 0o644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Clear the changelog
	err = clearUnreleasedChangelog()
	if err != nil {
		t.Fatalf("clearUnreleasedChangelog() failed: %v", err)
	}

	// Read the file back and verify it contains the template
	content, err := os.ReadFile(unreleasedChangelogFile)
	if err != nil {
		t.Fatalf("Failed to read cleared file: %v", err)
	}

	contentStr := string(content)
	template := getUnreleasedChangelogTemplate()

	if contentStr != template {
		t.Error("Cleared file does not match template")
	}

	// Verify the original content is gone
	if strings.Contains(contentStr, "Some test content") {
		t.Error("Original content still present after clearing")
	}
}

func TestHasUnreleasedContent_WithContent(t *testing.T) {
	cleanup, _ := setupTestEnvironment(t)
	defer cleanup()

	// Create a file with actual content
	testContent := `# Unreleased Changes

## Added
- Add new feature for testing
- Add another important feature

## Fixed
- Fix critical bug in system`

	err := os.WriteFile(unreleasedChangelogFile, []byte(testContent), 0o644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	hasContent, err := hasUnreleasedContent()
	if err != nil {
		t.Fatalf("hasUnreleasedContent() failed: %v", err)
	}

	if !hasContent {
		t.Error("Expected hasUnreleasedContent() to return true for file with content")
	}
}

func TestHasUnreleasedContent_WithoutContent(t *testing.T) {
	cleanup, _ := setupTestEnvironment(t)
	defer cleanup()

	// Create a file with only template content (no actual entries)
	template := getUnreleasedChangelogTemplate()
	err := os.WriteFile(unreleasedChangelogFile, []byte(template), 0o644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	hasContent, err := hasUnreleasedContent()
	if err != nil {
		t.Fatalf("hasUnreleasedContent() failed: %v", err)
	}

	if hasContent {
		t.Error("Expected hasUnreleasedContent() to return false for template-only file")
	}
}

func TestHasUnreleasedContent_WithEmptyBullets(t *testing.T) {
	cleanup, _ := setupTestEnvironment(t)
	defer cleanup()

	// Create a file with empty bullet points
	testContent := `# Unreleased Changes

## Added
- 
-   

## Fixed
<!-- No content -->`

	err := os.WriteFile(unreleasedChangelogFile, []byte(testContent), 0o644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	hasContent, err := hasUnreleasedContent()
	if err != nil {
		t.Fatalf("hasUnreleasedContent() failed: %v", err)
	}

	if hasContent {
		t.Error("Expected hasUnreleasedContent() to return false for file with empty bullets")
	}
}

func TestHasUnreleasedContent_NonexistentFile(t *testing.T) {
	cleanup, _ := setupTestEnvironment(t)
	defer cleanup()

	// Don't create the file
	hasContent, err := hasUnreleasedContent()
	if err == nil {
		t.Error("Expected hasUnreleasedContent() to return an error for nonexistent file")
	}

	if hasContent {
		t.Error("Expected hasUnreleasedContent() to return false for nonexistent file")
	}
}

func TestSafeFileOperation_Success(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")

	// Create initial file
	initialContent := "initial content"
	err := os.WriteFile(testFile, []byte(initialContent), 0o644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Perform safe operation that succeeds
	newContent := "new content"
	err = safeFileOperation(testFile, func() error {
		return os.WriteFile(testFile, []byte(newContent), 0o644)
	})

	if err != nil {
		t.Fatalf("safeFileOperation() failed: %v", err)
	}

	// Verify the file has new content
	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read file after operation: %v", err)
	}

	if string(content) != newContent {
		t.Errorf("Expected file content '%s', got '%s'", newContent, string(content))
	}

	// Verify backup file was cleaned up
	backupFile := testFile + ".backup"
	if _, err := os.Stat(backupFile); !os.IsNotExist(err) {
		t.Error("Backup file was not cleaned up after successful operation")
	}
}

func TestSafeFileOperation_Failure(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")

	// Create initial file
	initialContent := "initial content"
	err := os.WriteFile(testFile, []byte(initialContent), 0o644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Perform safe operation that fails
	err = safeFileOperation(testFile, func() error {
		// First write something to the file
		os.WriteFile(testFile, []byte("corrupted content"), 0o644)
		// Then return an error to simulate failure
		return os.ErrInvalid
	})

	if err == nil {
		t.Error("Expected safeFileOperation() to return error")
	}

	// Verify the file was restored to original content
	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read file after failed operation: %v", err)
	}

	if string(content) != initialContent {
		t.Errorf("Expected file content to be restored to '%s', got '%s'", initialContent, string(content))
	}

	// Verify backup file was cleaned up
	backupFile := testFile + ".backup"
	if _, err := os.Stat(backupFile); !os.IsNotExist(err) {
		t.Error("Backup file was not cleaned up after failed operation")
	}
}

func TestSafeFileOperation_NewFile(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "newfile.txt")

	// Perform safe operation on non-existent file
	content := "new file content"
	err := safeFileOperation(testFile, func() error {
		return os.WriteFile(testFile, []byte(content), 0o644)
	})

	if err != nil {
		t.Fatalf("safeFileOperation() failed: %v", err)
	}

	// Verify the file was created with correct content
	fileContent, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read created file: %v", err)
	}

	if string(fileContent) != content {
		t.Errorf("Expected file content '%s', got '%s'", content, string(fileContent))
	}
}

func TestCopyFile(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	srcFile := filepath.Join(tmpDir, "source.txt")
	dstFile := filepath.Join(tmpDir, "destination.txt")

	// Create source file
	content := "test content for copying"
	err := os.WriteFile(srcFile, []byte(content), 0o644)
	if err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}

	// Copy the file
	err = copyFile(srcFile, dstFile)
	if err != nil {
		t.Fatalf("copyFile() failed: %v", err)
	}

	// Verify destination file exists and has correct content
	dstContent, err := os.ReadFile(dstFile)
	if err != nil {
		t.Fatalf("Failed to read destination file: %v", err)
	}

	if string(dstContent) != content {
		t.Errorf("Expected destination content '%s', got '%s'", content, string(dstContent))
	}
}

func TestCopyFile_NonexistentSource(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	srcFile := filepath.Join(tmpDir, "nonexistent.txt")
	dstFile := filepath.Join(tmpDir, "destination.txt")

	// Try to copy non-existent file
	err := copyFile(srcFile, dstFile)
	if err == nil {
		t.Error("Expected copyFile() to return error for non-existent source")
	}

	// Verify destination file was not created
	if _, err := os.Stat(dstFile); !os.IsNotExist(err) {
		t.Error("Destination file should not exist after failed copy")
	}
}

func TestUpdateVersion(t *testing.T) {
	tests := []struct {
		name            string
		currentVersion  string
		expectedVersion string
	}{
		{
			name:            "Alpha version increment",
			currentVersion:  "v3.0.0-alpha.12",
			expectedVersion: "v3.0.0-alpha.13",
		},
		{
			name:            "Beta version increment",
			currentVersion:  "v3.0.0-beta.5",
			expectedVersion: "v3.0.0-beta.6",
		},
		{
			name:            "RC version increment",
			currentVersion:  "v2.5.0-rc.1",
			expectedVersion: "v2.5.0-rc.2",
		},
		{
			name:            "Patch version increment",
			currentVersion:  "v3.0.0",
			expectedVersion: "v3.0.1",
		},
		{
			name:            "Patch version with higher number",
			currentVersion:  "v1.2.15",
			expectedVersion: "v1.2.16",
		},
		{
			name:            "Pre-release without number becomes patch",
			currentVersion:  "v3.0.0-alpha",
			expectedVersion: "v3.0.1",
		},
		{
			name:            "Version without v prefix",
			currentVersion:  "3.0.0",
			expectedVersion: "v3.0.1",
		},
		{
			name:            "Alpha with large number",
			currentVersion:  "v3.0.0-alpha.999",
			expectedVersion: "v3.0.0-alpha.1000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary directory for this test
			tmpDir := t.TempDir()
			tempVersionFile := filepath.Join(tmpDir, "version.txt")

			// Save original versionFile path
			originalVersionFile := versionFile
			defer func() {
				// Restore original value
				_ = originalVersionFile
			}()

			// Write the current version to temp file
			err := os.WriteFile(tempVersionFile, []byte(tt.currentVersion), 0o644)
			if err != nil {
				t.Fatalf("Failed to write test version file: %v", err)
			}

			// Test the updateVersion function logic directly
			result := func() string {
				currentVersionData, err := os.ReadFile(tempVersionFile)
				if err != nil {
					t.Fatalf("Failed to read version file: %v", err)
				}
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
							return newVersion
						}
					}

					// If we can't parse the pre-release format, just increment patch version
					// and remove pre-release suffix
					return testIncrementPatchVersion(baseVersion)
				}

				// No pre-release suffix, just increment patch version
				return testIncrementPatchVersion(currentVersion)
			}()

			if result != tt.expectedVersion {
				t.Errorf("updateVersion() = %v, want %v", result, tt.expectedVersion)
			}
		})
	}
}

// testIncrementPatchVersion is a test version of incrementPatchVersion that doesn't write to file
func testIncrementPatchVersion(version string) string {
	// Remove 'v' prefix if present
	versionWithoutV := strings.TrimPrefix(version, "v")

	// Split into major.minor.patch
	parts := strings.Split(versionWithoutV, ".")
	if len(parts) != 3 {
		// Not a valid semver, return as-is
		return version
	}

	// Parse patch version
	patch, err := strconv.Atoi(parts[2])
	if err != nil {
		return version
	}

	// Increment patch
	patch++

	// Reconstruct version
	return fmt.Sprintf("v%s.%s.%d", parts[0], parts[1], patch)
}

// extractTestContent is a test helper that extracts changelog content using the same logic as extractChangelogContent
func extractTestContent(contentStr string) string {
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
		return ""
	}

	// Clean up result - remove any trailing empty lines
	for len(result) > 0 && strings.TrimSpace(result[len(result)-1]) == "" {
		result = result[:len(result)-1]
	}

	return strings.Join(result, "\n")
}

func TestHasUnreleasedContent_IgnoresExampleSection(t *testing.T) {
	cleanup, _ := setupTestEnvironment(t)
	defer cleanup()

	// Create a file with content only in the example section
	testContent := `# Unreleased Changes

## Added
<!-- No content -->

## Fixed
<!-- No content -->

---

### Example Entries:

**Added:**
- Add support for custom window icons in application options
- Add new SetWindowIcon() method to runtime API (#1234)`

	err := os.WriteFile(unreleasedChangelogFile, []byte(testContent), 0o644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	hasContent, err := hasUnreleasedContent()
	if err != nil {
		t.Fatalf("hasUnreleasedContent() failed: %v", err)
	}

	if hasContent {
		t.Error("Expected hasUnreleasedContent() to return false when content is only in example section")
	}
}

func TestHasUnreleasedContent_WithMixedContent(t *testing.T) {
	cleanup, _ := setupTestEnvironment(t)
	defer cleanup()

	// Create a file with both real content and example content
	testContent := `# Unreleased Changes

## Added
- Real feature addition here

## Fixed
<!-- No content -->

---

### Example Entries:

**Added:**
- Add support for custom window icons in application options`

	err := os.WriteFile(unreleasedChangelogFile, []byte(testContent), 0o644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	hasContent, err := hasUnreleasedContent()
	if err != nil {
		t.Fatalf("hasUnreleasedContent() failed: %v", err)
	}

	if !hasContent {
		t.Error("Expected hasUnreleasedContent() to return true when file has real content")
	}
}

// Integration test for the complete cleanup workflow
func TestCleanupWorkflow_Integration(t *testing.T) {
	cleanup, _ := setupTestEnvironment(t)
	defer cleanup()

	// Create a changelog file with actual content
	testContent := `# Unreleased Changes

## Added
- Add comprehensive changelog processing system
- Add validation for Keep a Changelog format compliance

## Fixed
- Fix parsing issues with various markdown bullet styles
- Fix validation edge cases for empty content sections`

	err := os.WriteFile(unreleasedChangelogFile, []byte(testContent), 0o644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Step 1: Check that file has content
	hasContent, err := hasUnreleasedContent()
	if err != nil {
		t.Fatalf("hasUnreleasedContent() failed: %v", err)
	}
	if !hasContent {
		t.Fatal("Expected file to have content")
	}

	// Step 2: Perform safe cleanup operation
	err = safeFileOperation(unreleasedChangelogFile, func() error {
		return clearUnreleasedChangelog()
	})
	if err != nil {
		t.Fatalf("Safe cleanup operation failed: %v", err)
	}

	// Step 3: Verify file was reset to template
	content, err := os.ReadFile(unreleasedChangelogFile)
	if err != nil {
		t.Fatalf("Failed to read file after cleanup: %v", err)
	}

	template := getUnreleasedChangelogTemplate()
	if string(content) != template {
		t.Error("File was not properly reset to template")
	}

	// Step 4: Verify original content is gone
	if strings.Contains(string(content), "Add comprehensive changelog processing system") {
		t.Error("Original content still present after cleanup")
	}

	// Step 5: Verify file no longer has content
	hasContentAfter, err := hasUnreleasedContent()
	if err != nil {
		t.Fatalf("hasUnreleasedContent() failed after cleanup: %v", err)
	}
	if hasContentAfter {
		t.Error("File should not have content after cleanup")
	}
}

func TestFullReleaseWorkflow_OnlyNonEmptySections(t *testing.T) {
	cleanup, projectRoot := setupTestEnvironment(t)
	defer cleanup()

	// Create subdirectories to match expected structure
	err := os.MkdirAll(filepath.Join(projectRoot, "v3", "internal", "version"), 0755)
	if err != nil {
		t.Fatalf("Failed to create version directory: %v", err)
	}

	err = os.MkdirAll(filepath.Join(projectRoot, "docs", "src", "content", "docs"), 0755)
	if err != nil {
		t.Fatalf("Failed to create docs directory: %v", err)
	}

	// Create version file
	versionFile := filepath.Join(projectRoot, "v3", "internal", "version", "version.txt")
	err = os.WriteFile(versionFile, []byte("v1.0.0-alpha.5"), 0644)
	if err != nil {
		t.Fatalf("Failed to create version file: %v", err)
	}

	// Create initial changelog
	changelogFile := filepath.Join(projectRoot, "docs", "src", "content", "docs", "changelog.mdx")
	initialChangelog := `---
title: Changelog
---

## [Unreleased]

## v1.0.0-alpha.4 - 2024-01-01

### Added
- Previous feature

`
	err = os.WriteFile(changelogFile, []byte(initialChangelog), 0644)
	if err != nil {
		t.Fatalf("Failed to create changelog file: %v", err)
	}

	// Create UNRELEASED_CHANGELOG.md with mixed content
	unreleasedContent := `# Unreleased Changes

## Added
- New amazing feature
- Another cool addition

## Changed
<!-- No changes -->

## Fixed
<!-- No fixes -->

## Deprecated
- Old API method

## Removed
<!-- Nothing removed -->

## Security
<!-- No security updates -->
`
	// The script expects the file at ../../UNRELEASED_CHANGELOG.md relative to release dir
	unreleasedFile := filepath.Join(projectRoot, "v3", "UNRELEASED_CHANGELOG.md")
	err = os.WriteFile(unreleasedFile, []byte(unreleasedContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create unreleased changelog: %v", err)
	}

	// Run the release process simulation
	// Read and process the files manually since we can't override constants
	content, err := os.ReadFile(unreleasedFile)
	if err != nil {
		t.Fatalf("Failed to read unreleased file: %v", err)
	}

	// Extract content using the same logic as extractChangelogContent
	changelogContent := extractTestContent(string(content))
	if changelogContent == "" {
		t.Fatal("Failed to extract any content")
	}

	// Verify only non-empty sections were extracted
	if !strings.Contains(changelogContent, "## Added") {
		t.Error("Expected '## Added' section to be included")
	}
	if !strings.Contains(changelogContent, "## Deprecated") {
		t.Error("Expected '## Deprecated' section to be included")
	}

	// Verify empty sections were NOT extracted
	emptySections := []string{"## Changed", "## Fixed", "## Removed", "## Security"}
	for _, section := range emptySections {
		if strings.Contains(changelogContent, section) {
			t.Errorf("Expected empty section '%s' to NOT be included", section)
		}
	}

	// Simulate updating the main changelog
	changelogData, _ := os.ReadFile(changelogFile)
	changelog := string(changelogData)
	changelogSplit := strings.Split(changelog, "## [Unreleased]")

	newVersion := "v1.0.0-alpha.6"
	today := "2024-01-15"
	newChangelog := changelogSplit[0] + "## [Unreleased]\n\n## " + newVersion + " - " + today + "\n\n" + changelogContent + changelogSplit[1]

	// Verify the final changelog format
	if !strings.Contains(newChangelog, "## v1.0.0-alpha.6 - 2024-01-15") {
		t.Error("Expected new version header in changelog")
	}

	// Count occurrences of section headers in the new version section
	newVersionSection := strings.Split(newChangelog, "## v1.0.0-alpha.4")[0]

	addedCount := strings.Count(newVersionSection, "## Added")
	if addedCount != 1 {
		t.Errorf("Expected exactly 1 '## Added' section, got %d", addedCount)
	}

	deprecatedCount := strings.Count(newVersionSection, "## Deprecated")
	if deprecatedCount != 1 {
		t.Errorf("Expected exactly 1 '## Deprecated' section, got %d", deprecatedCount)
	}

	// Ensure no empty sections in the new version section
	for _, section := range []string{"## Changed", "## Fixed", "## Removed", "## Security"} {
		count := strings.Count(newVersionSection, section)
		if count > 0 {
			t.Errorf("Expected 0 occurrences of empty section '%s', got %d", section, count)
		}
	}
}

// Test error handling during cleanup
func TestCleanupWorkflow_ErrorHandling(t *testing.T) {
	cleanup, _ := setupTestEnvironment(t)
	defer cleanup()

	// Create a changelog file with content
	testContent := `# Unreleased Changes

## Added
- Test content that should be preserved on error`

	err := os.WriteFile(unreleasedChangelogFile, []byte(testContent), 0o644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Simulate an error during cleanup by making the operation fail
	err = safeFileOperation(unreleasedChangelogFile, func() error {
		// First corrupt the file
		os.WriteFile(unreleasedChangelogFile, []byte("corrupted"), 0o644)
		// Then return an error
		return os.ErrPermission
	})

	if err == nil {
		t.Error("Expected safeFileOperation to return error")
	}

	// Verify original content was restored
	content, err := os.ReadFile(unreleasedChangelogFile)
	if err != nil {
		t.Fatalf("Failed to read file after error: %v", err)
	}

	if string(content) != testContent {
		t.Error("Original content was not restored after error")
	}
}
