// +build ignore

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	fmt.Println("Testing release.go --create-release-notes functionality")
	fmt.Println("=" + strings.Repeat("=", 50))

	// Test cases
	testCases := []struct {
		name     string
		content  string
		expected string
		shouldFail bool
	}{
		{
			name: "Valid changelog with content",
			content: `# Unreleased Changes

<!-- Comments -->

## Added
- Add Windows dark theme support for menus
- Add new --create-release-notes flag

## Changed
- Update Go version to 1.23
- Improve error handling

## Fixed
- Fix nightly release workflow
- Fix changelog extraction

---

### Example Entries:
Example content here`,
			expected: `## Added
- Add Windows dark theme support for menus
- Add new --create-release-notes flag

## Changed
- Update Go version to 1.23
- Improve error handling

## Fixed
- Fix nightly release workflow
- Fix changelog extraction`,
			shouldFail: false,
		},
		{
			name: "Empty changelog",
			content: `# Unreleased Changes

## Added
<!-- New features -->

## Changed
<!-- Changes -->

## Fixed
<!-- Bug fixes -->

---`,
			expected: "",
			shouldFail: true,
		},
		{
			name: "Only one section with content",
			content: `# Unreleased Changes

## Added
<!-- New features -->

## Changed
- Single change item here

## Fixed
<!-- Bug fixes -->

---`,
			expected: `## Changed
- Single change item here`,
			shouldFail: false,
		},
	}

	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "release-test-*")
	if err != nil {
		fmt.Printf("Failed to create temp dir: %v\n", err)
		os.Exit(1)
	}
	defer os.RemoveAll(tmpDir)

	// Save current directory
	originalDir, _ := os.Getwd()
	
	for i, tc := range testCases {
		fmt.Printf("\nTest %d: %s\n", i+1, tc.name)
		fmt.Println("-" + strings.Repeat("-", 40))

		// Create test changelog
		changelogPath := filepath.Join(tmpDir, "UNRELEASED_CHANGELOG.md")
		err = os.WriteFile(changelogPath, []byte(tc.content), 0644)
		if err != nil {
			fmt.Printf("❌ Failed to write test changelog: %v\n", err)
			continue
		}

		// Create release notes path
		releaseNotesPath := filepath.Join(tmpDir, fmt.Sprintf("release_notes_%d.md", i))

		// Change to temp dir (so relative paths work)
		os.Chdir(tmpDir)

		// Run the command
		cmd := exec.Command("go", "run", filepath.Join(originalDir, "release.go"), "--create-release-notes", releaseNotesPath)
		output, err := cmd.CombinedOutput()

		// Change back
		os.Chdir(originalDir)

		if tc.shouldFail {
			if err == nil {
				fmt.Printf("❌ Expected failure but command succeeded\n")
				fmt.Printf("Output: %s\n", output)
			} else {
				fmt.Printf("✅ Failed as expected: %v\n", err)
			}
		} else {
			if err != nil {
				fmt.Printf("❌ Command failed: %v\n", err)
				fmt.Printf("Output: %s\n", output)
			} else {
				fmt.Printf("✅ Command succeeded\n")
				
				// Read and verify the output
				content, err := os.ReadFile(releaseNotesPath)
				if err != nil {
					fmt.Printf("❌ Failed to read release notes: %v\n", err)
				} else {
					actualContent := strings.TrimSpace(string(content))
					expectedContent := strings.TrimSpace(tc.expected)
					
					if actualContent == expectedContent {
						fmt.Printf("✅ Content matches expected\n")
					} else {
						fmt.Printf("❌ Content mismatch\n")
						fmt.Printf("Expected:\n%s\n", expectedContent)
						fmt.Printf("Actual:\n%s\n", actualContent)
					}
				}
			}
		}

		// Clean up
		os.Remove(changelogPath)
		os.Remove(releaseNotesPath)
	}

	fmt.Println("\n" + "=" + strings.Repeat("=", 50))
	fmt.Println("Testing complete!")
}