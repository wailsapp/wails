package changelog

import (
	"strings"
	"testing"
)

// Test changelog content with misplaced entries
const testChangelogWithMisplacedEntries = `# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## v3.0.0-alpha.11 - 2025-07-12

### Added
- Add distribution-specific build dependencies for Linux by @leaanthony in [PR](https://github.com/wailsapp/wails/pull/4345)
- Added bindings guide by @atterpac in [PR](https://github.com/wailsapp/wails/pull/4404)
- This is a legitimate entry for this version

## v3.0.0-alpha.10 - 2025-07-06

### Added
- Original feature for alpha.10
- Another legitimate feature

### Fixed
- Some bug fix
- Add distribution-specific build dependencies for Linux by @leaanthony in [PR](https://github.com/wailsapp/wails/pull/4345)

## v3.0.0-alpha.9 - 2025-01-13

### Added
- Old feature from January
- Added bindings guide by @atterpac in [PR](https://github.com/wailsapp/wails/pull/4404)
`

const testChangelogClean = `# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

## v3.0.0-alpha.11 - 2025-07-12

### Added
- This is a legitimate entry for this version

### Fixed
- Some bug fix for alpha.11

## v3.0.0-alpha.10 - 2025-07-06

### Added
- Original feature for alpha.10
- Another legitimate feature

### Fixed
- Some bug fix
`

func TestChangelogParser_ParseContent(t *testing.T) {
	parser := NewChangelogParser()
	lines := strings.Split(testChangelogClean, "\n")

	err := parser.ParseContent(lines)
	if err != nil {
		t.Fatalf("ParseContent failed: %v", err)
	}

	sections := parser.GetSections()

	// Check that we have the expected sections
	expectedSections := []string{"Unreleased", "v3.0.0-alpha.11", "v3.0.0-alpha.10"}
	for _, expected := range expectedSections {
		if _, exists := sections[expected]; !exists {
			t.Errorf("Expected section %s not found", expected)
		}
	}

	// Check unreleased section
	unreleased := sections["Unreleased"]
	if unreleased == nil {
		t.Fatal("Unreleased section not found")
	}
	if !unreleased.IsUnreleased {
		t.Error("Unreleased section should be marked as unreleased")
	}

	// Check released sections
	alpha11 := sections["v3.0.0-alpha.11"]
	if alpha11 == nil {
		t.Fatal("v3.0.0-alpha.11 section not found")
	}
	if !alpha11.IsReleased {
		t.Error("v3.0.0-alpha.11 section should be marked as released")
	}
	if alpha11.Date != "2025-07-12" {
		t.Errorf("Expected date 2025-07-12, got %s", alpha11.Date)
	}

	// Check categories
	if addedCategory, exists := alpha11.Categories["Added"]; exists {
		if len(addedCategory.Items) != 1 {
			t.Errorf("Expected 1 item in Added category, got %d", len(addedCategory.Items))
		}
	} else {
		t.Error("Added category not found in v3.0.0-alpha.11")
	}
}

func TestChangelogParser_DetectMisplacedEntries(t *testing.T) {
	parser := NewChangelogParser()
	lines := strings.Split(testChangelogWithMisplacedEntries, "\n")

	err := parser.ParseContent(lines)
	if err != nil {
		t.Fatalf("ParseContent failed: %v", err)
	}

	result, err := parser.ValidateAndFixMisplacedEntries()
	if err != nil {
		t.Fatalf("ValidateAndFixMisplacedEntries failed: %v", err)
	}

	// Should detect misplaced entries
	if len(result.MisplacedEntries) == 0 {
		t.Error("Expected to find misplaced entries")
	}

	// Check that specific patterns are detected
	found := false
	for _, entry := range result.MisplacedEntries {
		if strings.Contains(entry.Text, "Add distribution-specific build dependencies") {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected to detect 'Add distribution-specific build dependencies' as misplaced")
	}

	// Should have warnings
	if len(result.Warnings) == 0 {
		t.Error("Expected warnings about misplaced entries")
	}
}

func TestChangelogParser_MoveMisplacedEntries(t *testing.T) {
	parser := NewChangelogParser()
	lines := strings.Split(testChangelogWithMisplacedEntries, "\n")

	err := parser.ParseContent(lines)
	if err != nil {
		t.Fatalf("ParseContent failed: %v", err)
	}

	// Get initial state
	sections := parser.GetSections()
	alpha11Before := sections["v3.0.0-alpha.11"]
	alpha10Before := sections["v3.0.0-alpha.10"]

	initialAlpha11AddedCount := len(alpha11Before.Categories["Added"].Items)
	initialAlpha10FixedCount := len(alpha10Before.Categories["Fixed"].Items)

	// Validate and fix
	result, err := parser.ValidateAndFixMisplacedEntries()
	if err != nil {
		t.Fatalf("ValidateAndFixMisplacedEntries failed: %v", err)
	}

	if len(result.MisplacedEntries) == 0 {
		t.Skip("No misplaced entries detected, skipping move test")
	}

	// Check that entries were moved to Unreleased
	unreleased := parser.GetUnreleasedSection()
	if unreleased == nil {
		t.Fatal("Unreleased section should exist after moving entries")
	}

	// Should have entries in Unreleased now
	hasEntries := false
	for _, category := range unreleased.Categories {
		if len(category.Items) > 0 {
			hasEntries = true
			break
		}
	}
	if !hasEntries {
		t.Error("Unreleased section should have entries after moving misplaced ones")
	}

	// Check that entries were removed from original sections
	alpha11After := sections["v3.0.0-alpha.11"]
	alpha10After := sections["v3.0.0-alpha.10"]

	if len(alpha11After.Categories["Added"].Items) >= initialAlpha11AddedCount {
		t.Error("Expected some entries to be removed from v3.0.0-alpha.11 Added section")
	}

	if len(alpha10After.Categories["Fixed"].Items) >= initialAlpha10FixedCount {
		t.Error("Expected some entries to be removed from v3.0.0-alpha.10 Fixed section")
	}
}

func TestChangelogParser_GenerateChangelog(t *testing.T) {
	parser := NewChangelogParser()
	lines := strings.Split(testChangelogClean, "\n")

	err := parser.ParseContent(lines)
	if err != nil {
		t.Fatalf("ParseContent failed: %v", err)
	}

	generated := parser.GenerateChangelog()
	generatedContent := strings.Join(generated, "\n")

	// Should contain the main sections
	// Note: Empty unreleased section won't be generated, which is correct behavior
	// if !strings.Contains(generatedContent, "## [Unreleased]") {
	//	t.Error("Generated changelog should contain Unreleased section")
	// }

	if !strings.Contains(generatedContent, "## v3.0.0-alpha.11 - 2025-07-12") {
		t.Error("Generated changelog should contain v3.0.0-alpha.11 section with date")
	}

	// Should maintain proper structure
	if !strings.Contains(generatedContent, "### Added") {
		t.Error("Generated changelog should contain Added category")
	}

	if !strings.Contains(generatedContent, "### Fixed") {
		t.Error("Generated changelog should contain Fixed category")
	}
}

func TestChangelogParser_SectionOrdering(t *testing.T) {
	parser := NewChangelogParser()
	lines := strings.Split(testChangelogClean, "\n")

	err := parser.ParseContent(lines)
	if err != nil {
		t.Fatalf("ParseContent failed: %v", err)
	}

	generated := parser.GenerateChangelog()
	generatedContent := strings.Join(generated, "\n")

	// Check that Unreleased comes before released versions
	unreleasedPos := strings.Index(generatedContent, "## [Unreleased]")
	alpha11Pos := strings.Index(generatedContent, "## v3.0.0-alpha.11")
	alpha10Pos := strings.Index(generatedContent, "## v3.0.0-alpha.10")

	if unreleasedPos > alpha11Pos && alpha11Pos > 0 {
		t.Error("Unreleased section should come before released versions")
	}

	if alpha11Pos > alpha10Pos && alpha10Pos > 0 {
		t.Error("Newer versions should come before older versions")
	}
}

func TestChangelogParser_EmptyUnreleased(t *testing.T) {
	// Test with empty unreleased section
	content := `# Changelog

## [Unreleased]

## v1.0.0 - 2025-01-01

### Added
- Initial release
`

	parser := NewChangelogParser()
	err := parser.ParseString(content)
	if err != nil {
		t.Fatalf("ParseString failed: %v", err)
	}

	generated := parser.GenerateChangelog()
	generatedContent := strings.Join(generated, "\n")

	// Empty unreleased section should not appear in output
	if strings.Contains(generatedContent, "## [Unreleased]") {
		t.Error("Empty unreleased section should not appear in generated content")
	}
}

func TestChangelogParser_CategoryOrdering(t *testing.T) {
	content := `# Changelog

## v1.0.0 - 2025-01-01

### Fixed
- Bug fix

### Security
- Security update

### Added
- New feature

### Breaking Changes
- Breaking change

### Changed
- Change

### Deprecated
- Deprecated feature

### Removed
- Removed feature
`

	parser := NewChangelogParser()
	err := parser.ParseString(content)
	if err != nil {
		t.Fatalf("ParseString failed: %v", err)
	}

	generated := parser.GenerateChangelog()
	generatedContent := strings.Join(generated, "\n")

	// Check that categories appear in the correct order
	expectedOrder := []string{
		"### Breaking Changes",
		"### Added",
		"### Changed",
		"### Deprecated",
		"### Removed",
		"### Fixed",
		"### Security",
	}

	var positions []int
	for _, category := range expectedOrder {
		pos := strings.Index(generatedContent, category)
		if pos == -1 {
			t.Errorf("Category %s not found in generated content", category)
		}
		positions = append(positions, pos)
	}

	// Check ordering
	for i := 1; i < len(positions); i++ {
		if positions[i-1] > positions[i] {
			t.Errorf("Categories are not in the correct order. %s should come before %s",
				expectedOrder[i-1], expectedOrder[i])
		}
	}
}

func TestChangelogParser_LooksLikeRecentEntry(t *testing.T) {
	parser := NewChangelogParser()

	testCases := []struct {
		entry    string
		version  string
		expected bool
		desc     string
	}{
		{
			entry:    "Add distribution-specific build dependencies for Linux by @leaanthony in [PR](https://github.com/wailsapp/wails/pull/4345)",
			version:  "v3.0.0-alpha.9",
			expected: true,
			desc:     "Should detect suspicious pattern",
		},
		{
			entry:    "Added bindings guide by @atterpac in [PR](https://github.com/wailsapp/wails/pull/4404)",
			version:  "v3.0.0-alpha.9",
			expected: true,
			desc:     "Should detect another suspicious pattern",
		},
		{
			entry:    "Regular feature addition",
			version:  "v3.0.0-alpha.11",
			expected: false,
			desc:     "Should not flag normal entries",
		},
	}

	for _, tc := range testCases {
		result := parser.looksLikeRecentEntry(tc.entry, tc.version)
		if result != tc.expected {
			t.Errorf("%s: expected %v, got %v for entry: %s", tc.desc, tc.expected, result, tc.entry)
		}
	}
}

func BenchmarkChangelogParser_ParseContent(b *testing.B) {
	lines := strings.Split(testChangelogWithMisplacedEntries, "\n")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parser := NewChangelogParser()
		parser.ParseContent(lines)
	}
}

func BenchmarkChangelogParser_ValidateAndFix(b *testing.B) {
	parser := NewChangelogParser()
	lines := strings.Split(testChangelogWithMisplacedEntries, "\n")
	parser.ParseContent(lines)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parser.ValidateAndFixMisplacedEntries()
	}
}
