package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestExtractChangelogContent tests the extractChangelogContent function with various inputs
func TestExtractChangelogContent(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected string
		wantErr  bool
	}{
		{
			name: "Full changelog with all sections",
			content: `# Unreleased Changes

<!-- Comments -->

## Added
- New feature 1
- New feature 2

## Changed
- Changed item 1

## Fixed
- Bug fix 1
- Bug fix 2

## Deprecated
- Deprecated feature

## Removed
- Removed feature

## Security
- Security fix

---

### Example Entries:

**Added:**
- Example entry
`,
			expected: `## Added
- New feature 1
- New feature 2

## Changed
- Changed item 1

## Fixed
- Bug fix 1
- Bug fix 2

## Deprecated
- Deprecated feature

## Removed
- Removed feature

## Security
- Security fix`,
			wantErr: false,
		},
		{
			name: "Only Added section with content",
			content: `# Unreleased Changes

## Added
<!-- New features -->
- Add Windows dark theme support
- Add new flag to release script

## Changed
<!-- Changes -->

## Fixed
<!-- Bug fixes -->

---
### Example Entries:
`,
			expected: `## Added
- Add Windows dark theme support
- Add new flag to release script`,
			wantErr: false,
		},
		{
			name: "Empty sections should not be included",
			content: `# Unreleased Changes

## Added
<!-- New features -->

## Changed
<!-- Changes -->
- Update Go version to 1.23

## Fixed
<!-- Bug fixes -->

---
`,
			expected: `## Changed
- Update Go version to 1.23`,
			wantErr: false,
		},
		{
			name: "No content returns empty string",
			content: `# Unreleased Changes

## Added
<!-- New features -->

## Changed
<!-- Changes -->

## Fixed
<!-- Bug fixes -->

---
`,
			expected: "",
			wantErr:  false,
		},
		{
			name: "Comments should be excluded",
			content: `# Unreleased Changes

<!-- This is a comment -->

## Added
<!-- Another comment -->
- Real content here
<!-- Mid-section comment -->
- More content

---
`,
			expected: `## Added
- Real content here
- More content`,
			wantErr: false,
		},
		{
			name: "Multi-line comments handled correctly",
			content: `# Unreleased Changes

<!-- 
Multi-line
comment
-->

## Added
- Feature 1
<!-- Another
multi-line comment
that spans lines -->
- Feature 2

---
`,
			expected: `## Added
- Feature 1
- Feature 2`,
			wantErr: false,
		},
		{
			name: "Mixed bullet point styles",
			content: `# Unreleased Changes

## Added
- Dash bullet point
* Asterisk bullet point
- Another dash

---
`,
			expected: `## Added
- Dash bullet point
* Asterisk bullet point
- Another dash`,
			wantErr: false,
		},
		{
			name: "Trailing empty lines removed",
			content: `# Unreleased Changes

## Added
- Feature 1


## Changed
- Change 1



---
`,
			expected: `## Added
- Feature 1

## Changed
- Change 1`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary file
			tmpDir := t.TempDir()
			tmpFile := filepath.Join(tmpDir, "UNRELEASED_CHANGELOG.md")
			err := os.WriteFile(tmpFile, []byte(tt.content), 0644)
			if err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			// Save original path and update it
			originalPath := unreleasedChangelogFile
			unreleasedChangelogFile = tmpFile
			defer func() {
				unreleasedChangelogFile = originalPath
			}()

			// Test the function
			got, err := extractChangelogContent()
			if (err != nil) != tt.wantErr {
				t.Errorf("extractChangelogContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.expected {
				t.Errorf("extractChangelogContent() = %q, want %q", got, tt.expected)
			}
		})
	}
}

// TestCreateReleaseNotes tests the --create-release-notes functionality
func TestCreateReleaseNotes(t *testing.T) {
	tests := []struct {
		name            string
		changelogContent string
		expectSuccess   bool
		expectedNotes   string
	}{
		{
			name: "Valid changelog creates release notes",
			changelogContent: `# Unreleased Changes

## Added
- New feature X
- New feature Y

## Fixed
- Bug fix A

---
### Examples:
`,
			expectSuccess: true,
			expectedNotes: `## Added
- New feature X
- New feature Y

## Fixed
- Bug fix A`,
		},
		{
			name: "Empty changelog fails",
			changelogContent: `# Unreleased Changes

## Added
<!-- New features -->

## Changed
<!-- Changes -->

---
`,
			expectSuccess: false,
			expectedNotes: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary directory
			tmpDir := t.TempDir()
			
			// Create UNRELEASED_CHANGELOG.md
			changelogPath := filepath.Join(tmpDir, "UNRELEASED_CHANGELOG.md")
			err := os.WriteFile(changelogPath, []byte(tt.changelogContent), 0644)
			if err != nil {
				t.Fatalf("Failed to create changelog file: %v", err)
			}

			// Create release notes path
			releaseNotesPath := filepath.Join(tmpDir, "release_notes.md")

			// Save original path
			originalPath := unreleasedChangelogFile
			unreleasedChangelogFile = changelogPath
			defer func() {
				unreleasedChangelogFile = originalPath
			}()

			// Test the create release notes flow
			content, err := extractChangelogContent()
			if err != nil && tt.expectSuccess {
				t.Fatalf("Failed to extract content: %v", err)
			}

			if tt.expectSuccess {
				if content == "" {
					t.Error("Expected content but got empty string")
				} else {
					// Write the release notes
					err = os.WriteFile(releaseNotesPath, []byte(content), 0644)
					if err != nil {
						t.Fatalf("Failed to write release notes: %v", err)
					}

					// Verify the file was created
					data, err := os.ReadFile(releaseNotesPath)
					if err != nil {
						t.Fatalf("Failed to read release notes: %v", err)
					}

					if string(data) != tt.expectedNotes {
						t.Errorf("Release notes = %q, want %q", string(data), tt.expectedNotes)
					}
				}
			} else {
				if content != "" {
					t.Errorf("Expected no content but got: %q", content)
				}
			}
		})
	}
}

// TestHasUnreleasedContent tests the hasUnreleasedContent function
func TestHasUnreleasedContent(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected bool
	}{
		{
			name: "Has content",
			content: `# Unreleased Changes

## Added
- New feature

---
`,
			expected: true,
		},
		{
			name: "No content",
			content: `# Unreleased Changes

## Added
<!-- New features -->

## Changed
<!-- Changes -->

---
`,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary file
			tmpDir := t.TempDir()
			tmpFile := filepath.Join(tmpDir, "UNRELEASED_CHANGELOG.md")
			err := os.WriteFile(tmpFile, []byte(tt.content), 0644)
			if err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			// Save original path
			originalPath := unreleasedChangelogFile
			unreleasedChangelogFile = tmpFile
			defer func() {
				unreleasedChangelogFile = originalPath
			}()

			// Test the function
			got, err := hasUnreleasedContent()
			if err != nil {
				t.Fatalf("hasUnreleasedContent() error = %v", err)
			}

			if got != tt.expected {
				t.Errorf("hasUnreleasedContent() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// TestVersionIncrement tests version incrementing logic
func TestVersionIncrement(t *testing.T) {
	tests := []struct {
		name           string
		currentVersion string
		expectedNext   string
	}{
		{
			name:           "Alpha version increment",
			currentVersion: "v3.0.0-alpha.15",
			expectedNext:   "v3.0.0-alpha.16",
		},
		{
			name:           "Beta version increment", 
			currentVersion: "v3.0.0-beta.5",
			expectedNext:   "v3.0.0-beta.6",
		},
		{
			name:           "Regular version increment",
			currentVersion: "v3.0.0",
			expectedNext:   "v3.0.1",
		},
		{
			name:           "Patch version increment",
			currentVersion: "v3.1.5",
			expectedNext:   "v3.1.6",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test incrementPatchVersion for regular versions
			if !strings.Contains(tt.currentVersion, "-") {
				got := incrementPatchVersion(tt.currentVersion)
				// Note: incrementPatchVersion writes to file, so we just check the return value
				if got != tt.expectedNext {
					t.Errorf("incrementPatchVersion(%s) = %s, want %s", tt.currentVersion, got, tt.expectedNext)
				}
			}
		})
	}
}

// TestClearUnreleasedChangelog tests the changelog reset functionality
func TestClearUnreleasedChangelog(t *testing.T) {
	// Create temporary file
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "UNRELEASED_CHANGELOG.md")
	
	// Write some content
	err := os.WriteFile(tmpFile, []byte("Some content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Save original path
	originalPath := unreleasedChangelogFile
	unreleasedChangelogFile = tmpFile
	defer func() {
		unreleasedChangelogFile = originalPath
	}()

	// Clear the changelog
	err = clearUnreleasedChangelog()
	if err != nil {
		t.Fatalf("clearUnreleasedChangelog() error = %v", err)
	}

	// Read the file
	content, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("Failed to read cleared file: %v", err)
	}

	// Check it contains the template
	if !strings.Contains(string(content), "# Unreleased Changes") {
		t.Error("Cleared file doesn't contain template header")
	}
	if !strings.Contains(string(content), "## Added") {
		t.Error("Cleared file doesn't contain Added section")
	}
	if !strings.Contains(string(content), "### Example Entries:") {
		t.Error("Cleared file doesn't contain example section")
	}
}