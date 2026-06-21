package changelog

import (
	"strings"
	"testing"
)

func TestNewParser(t *testing.T) {
	parser := NewParser()
	if parser == nil {
		t.Fatal("NewParser() returned nil")
	}
	if parser.sectionRegex == nil {
		t.Error("sectionRegex not initialized")
	}
	if parser.bulletRegex == nil {
		t.Error("bulletRegex not initialized")
	}
}

func TestParseContent_EmptyContent(t *testing.T) {
	parser := NewParser()
	reader := strings.NewReader("")

	entry, err := parser.ParseContent(reader)
	if err != nil {
		t.Fatalf("ParseContent() returned error: %v", err)
	}

	if entry.HasContent() {
		t.Error("Empty content should not have content")
	}
}

func TestParseContent_OnlyComments(t *testing.T) {
	parser := NewParser()
	content := `# Unreleased Changes

<!-- This is a comment -->
<!-- Another comment -->

## Added
<!-- Comment in section -->

## Changed
<!-- Another comment -->`

	reader := strings.NewReader(content)
	entry, err := parser.ParseContent(reader)
	if err != nil {
		t.Fatalf("ParseContent() returned error: %v", err)
	}

	if entry.HasContent() {
		t.Error("Content with only comments should not have content")
	}
}

func TestParseContent_BasicSections(t *testing.T) {
	parser := NewParser()
	content := `# Unreleased Changes

## Added
- New feature A
- New feature B

## Changed
- Changed feature C

## Fixed
- Fixed bug D
- Fixed bug E

## Deprecated
- Deprecated feature F

## Removed
- Removed feature G

## Security
- Security fix H`

	reader := strings.NewReader(content)
	entry, err := parser.ParseContent(reader)
	if err != nil {
		t.Fatalf("ParseContent() returned error: %v", err)
	}

	// Test Added section
	if len(entry.Added) != 2 {
		t.Errorf("Expected 2 Added items, got %d", len(entry.Added))
	}
	if entry.Added[0] != "New feature A" {
		t.Errorf("Expected 'New feature A', got '%s'", entry.Added[0])
	}
	if entry.Added[1] != "New feature B" {
		t.Errorf("Expected 'New feature B', got '%s'", entry.Added[1])
	}

	// Test Changed section
	if len(entry.Changed) != 1 {
		t.Errorf("Expected 1 Changed item, got %d", len(entry.Changed))
	}
	if entry.Changed[0] != "Changed feature C" {
		t.Errorf("Expected 'Changed feature C', got '%s'", entry.Changed[0])
	}

	// Test Fixed section
	if len(entry.Fixed) != 2 {
		t.Errorf("Expected 2 Fixed items, got %d", len(entry.Fixed))
	}
	if entry.Fixed[0] != "Fixed bug D" {
		t.Errorf("Expected 'Fixed bug D', got '%s'", entry.Fixed[0])
	}
	if entry.Fixed[1] != "Fixed bug E" {
		t.Errorf("Expected 'Fixed bug E', got '%s'", entry.Fixed[1])
	}

	// Test Deprecated section
	if len(entry.Deprecated) != 1 {
		t.Errorf("Expected 1 Deprecated item, got %d", len(entry.Deprecated))
	}
	if entry.Deprecated[0] != "Deprecated feature F" {
		t.Errorf("Expected 'Deprecated feature F', got '%s'", entry.Deprecated[0])
	}

	// Test Removed section
	if len(entry.Removed) != 1 {
		t.Errorf("Expected 1 Removed item, got %d", len(entry.Removed))
	}
	if entry.Removed[0] != "Removed feature G" {
		t.Errorf("Expected 'Removed feature G', got '%s'", entry.Removed[0])
	}

	// Test Security section
	if len(entry.Security) != 1 {
		t.Errorf("Expected 1 Security item, got %d", len(entry.Security))
	}
	if entry.Security[0] != "Security fix H" {
		t.Errorf("Expected 'Security fix H', got '%s'", entry.Security[0])
	}

	// Test HasContent
	if !entry.HasContent() {
		t.Error("Entry should have content")
	}
}

func TestParseContent_WithExampleSection(t *testing.T) {
	parser := NewParser()
	content := `# Unreleased Changes

## Added
- Real feature A

## Changed
- Real change B

---

### Example Entries:

**Added:**
- Example feature that should be ignored
- Another example that should be ignored

**Fixed:**
- Example fix that should be ignored`

	reader := strings.NewReader(content)
	entry, err := parser.ParseContent(reader)
	if err != nil {
		t.Fatalf("ParseContent() returned error: %v", err)
	}

	// Should only have the real entries, not the examples
	if len(entry.Added) != 1 {
		t.Errorf("Expected 1 Added item, got %d", len(entry.Added))
	}
	if entry.Added[0] != "Real feature A" {
		t.Errorf("Expected 'Real feature A', got '%s'", entry.Added[0])
	}

	if len(entry.Changed) != 1 {
		t.Errorf("Expected 1 Changed item, got %d", len(entry.Changed))
	}
	if entry.Changed[0] != "Real change B" {
		t.Errorf("Expected 'Real change B', got '%s'", entry.Changed[0])
	}

	// Should not have any Fixed items from examples
	if len(entry.Fixed) != 0 {
		t.Errorf("Expected 0 Fixed items, got %d", len(entry.Fixed))
	}
}

func TestParseContent_DifferentBulletStyles(t *testing.T) {
	parser := NewParser()
	content := `# Unreleased Changes

## Added
- Feature with dash
* Feature with asterisk
  - Indented feature with dash
  * Indented feature with asterisk

## Fixed
-   Feature with extra spaces
*   Another with extra spaces`

	reader := strings.NewReader(content)
	entry, err := parser.ParseContent(reader)
	if err != nil {
		t.Fatalf("ParseContent() returned error: %v", err)
	}

	expectedAdded := []string{
		"Feature with dash",
		"Feature with asterisk",
		"Indented feature with dash",
		"Indented feature with asterisk",
	}

	if len(entry.Added) != len(expectedAdded) {
		t.Errorf("Expected %d Added items, got %d", len(expectedAdded), len(entry.Added))
	}

	for i, expected := range expectedAdded {
		if i >= len(entry.Added) || entry.Added[i] != expected {
			t.Errorf("Expected Added[%d] to be '%s', got '%s'", i, expected, entry.Added[i])
		}
	}

	expectedFixed := []string{
		"Feature with extra spaces",
		"Another with extra spaces",
	}

	if len(entry.Fixed) != len(expectedFixed) {
		t.Errorf("Expected %d Fixed items, got %d", len(expectedFixed), len(entry.Fixed))
	}

	for i, expected := range expectedFixed {
		if i >= len(entry.Fixed) || entry.Fixed[i] != expected {
			t.Errorf("Expected Fixed[%d] to be '%s', got '%s'", i, expected, entry.Fixed[i])
		}
	}
}

func TestParseContent_EmptyBulletPoints(t *testing.T) {
	parser := NewParser()
	content := `# Unreleased Changes

## Added
- Valid feature
- 
-   
- Another valid feature

## Fixed
- 
- Valid fix`

	reader := strings.NewReader(content)
	entry, err := parser.ParseContent(reader)
	if err != nil {
		t.Fatalf("ParseContent() returned error: %v", err)
	}

	// Should skip empty bullet points
	expectedAdded := []string{
		"Valid feature",
		"Another valid feature",
	}

	if len(entry.Added) != len(expectedAdded) {
		t.Errorf("Expected %d Added items, got %d", len(expectedAdded), len(entry.Added))
	}

	for i, expected := range expectedAdded {
		if i >= len(entry.Added) || entry.Added[i] != expected {
			t.Errorf("Expected Added[%d] to be '%s', got '%s'", i, expected, entry.Added[i])
		}
	}

	expectedFixed := []string{"Valid fix"}
	if len(entry.Fixed) != len(expectedFixed) {
		t.Errorf("Expected %d Fixed items, got %d", len(expectedFixed), len(entry.Fixed))
	}
	if entry.Fixed[0] != "Valid fix" {
		t.Errorf("Expected 'Valid fix', got '%s'", entry.Fixed[0])
	}
}

func TestHasContent(t *testing.T) {
	tests := []struct {
		name     string
		entry    ChangelogEntry
		expected bool
	}{
		{
			name:     "Empty entry",
			entry:    ChangelogEntry{},
			expected: false,
		},
		{
			name: "Entry with Added items",
			entry: ChangelogEntry{
				Added: []string{"Feature A"},
			},
			expected: true,
		},
		{
			name: "Entry with Changed items",
			entry: ChangelogEntry{
				Changed: []string{"Change A"},
			},
			expected: true,
		},
		{
			name: "Entry with Fixed items",
			entry: ChangelogEntry{
				Fixed: []string{"Fix A"},
			},
			expected: true,
		},
		{
			name: "Entry with Deprecated items",
			entry: ChangelogEntry{
				Deprecated: []string{"Deprecated A"},
			},
			expected: true,
		},
		{
			name: "Entry with Removed items",
			entry: ChangelogEntry{
				Removed: []string{"Removed A"},
			},
			expected: true,
		},
		{
			name: "Entry with Security items",
			entry: ChangelogEntry{
				Security: []string{"Security A"},
			},
			expected: true,
		},
		{
			name: "Entry with multiple sections",
			entry: ChangelogEntry{
				Added:   []string{"Feature A"},
				Fixed:   []string{"Fix A"},
				Changed: []string{"Change A"},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.entry.HasContent(); got != tt.expected {
				t.Errorf("HasContent() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestFormatForChangelog(t *testing.T) {
	entry := ChangelogEntry{
		Added:      []string{"New feature A", "New feature B"},
		Changed:    []string{"Changed feature C"},
		Fixed:      []string{"Fixed bug D"},
		Deprecated: []string{"Deprecated feature E"},
		Removed:    []string{"Removed feature F"},
		Security:   []string{"Security fix G"},
	}

	result := entry.FormatForChangelog()

	expected := `### Added
- New feature A
- New feature B

### Changed
- Changed feature C

### Fixed
- Fixed bug D

### Deprecated
- Deprecated feature E

### Removed
- Removed feature F

### Security
- Security fix G`

	if result != expected {
		t.Errorf("FormatForChangelog() mismatch.\nExpected:\n%s\n\nGot:\n%s", expected, result)
	}
}

func TestFormatForChangelog_PartialSections(t *testing.T) {
	entry := ChangelogEntry{
		Added: []string{"New feature A"},
		Fixed: []string{"Fixed bug B"},
		// Other sections empty
	}

	result := entry.FormatForChangelog()

	expected := `### Added
- New feature A

### Fixed
- Fixed bug B`

	if result != expected {
		t.Errorf("FormatForChangelog() mismatch.\nExpected:\n%s\n\nGot:\n%s", expected, result)
	}
}

func TestFormatForRelease(t *testing.T) {
	entry := ChangelogEntry{
		Added:      []string{"New feature A", "New feature B"},
		Changed:    []string{"Changed feature C"},
		Fixed:      []string{"Fixed bug D"},
		Deprecated: []string{"Deprecated feature E"},
		Removed:    []string{"Removed feature F"},
		Security:   []string{"Security fix G"},
	}

	result := entry.FormatForRelease()

	expected := `## ✨ Added
- New feature A
- New feature B

## 🔄 Changed
- Changed feature C

## 🐛 Fixed
- Fixed bug D

## ⚠️ Deprecated
- Deprecated feature E

## 🗑️ Removed
- Removed feature F

## 🔒 Security
- Security fix G`

	if result != expected {
		t.Errorf("FormatForRelease() mismatch.\nExpected:\n%s\n\nGot:\n%s", expected, result)
	}
}

func TestFormatForRelease_EmptyEntry(t *testing.T) {
	entry := ChangelogEntry{}

	result := entry.FormatForRelease()

	if result != "" {
		t.Errorf("FormatForRelease() for empty entry should return empty string, got: %s", result)
	}
}
