package changelog

import (
	"fmt"
)

func ExampleChangelogParser_basic() {
	content := `# Changelog

## [Unreleased]

## v1.1.0 - 2025-01-15

### Added
- New feature A
- Accidentally added to wrong version by @dev in [PR](https://github.com/example/repo/pull/123)

### Fixed
- Bug fix B

## v1.0.0 - 2025-01-01

### Added  
- Initial release
- New feature A
`

	parser := NewChangelogParser()
	err := parser.ParseString(content)
	if err != nil {
		fmt.Printf("Error parsing: %v\n", err)
		return
	}

	// Validate and fix misplaced entries
	result, err := parser.ValidateAndFixMisplacedEntries()
	if err != nil {
		fmt.Printf("Error validating: %v\n", err)
		return
	}

	fmt.Printf("Validation result:\n")
	fmt.Printf("- Valid: %v\n", result.IsValid)
	fmt.Printf("- Misplaced entries found: %d\n", len(result.MisplacedEntries))
	fmt.Printf("- Warnings: %d\n", len(result.Warnings))

	if len(result.MisplacedEntries) > 0 {
		fmt.Printf("\nMisplaced entries:\n")
		for _, entry := range result.MisplacedEntries {
			fmt.Printf("- %s (from %s -> Unreleased)\n", entry.Text, entry.Section)
		}
	}

	// Generate corrected changelog
	corrected := parser.GenerateChangelog()
	fmt.Printf("\nCorrected changelog preview:\n")
	for i, line := range corrected {
		if i > 15 { // Only show first few lines
			fmt.Printf("... (truncated)\n")
			break
		}
		fmt.Printf("%s\n", line)
	}
}

func ExampleChangelogParser_detectPatterns() {
	parser := NewChangelogParser()

	// Test the heuristic detection
	testEntries := []string{
		"Add distribution-specific build dependencies for Linux by @leaanthony in [PR](https://github.com/wailsapp/wails/pull/4345)",
		"Added bindings guide by @atterpac in [PR](https://github.com/wailsapp/wails/pull/4404)",
		"Regular bug fix",
		"Normal feature addition",
	}

	fmt.Printf("Pattern detection test:\n")
	for _, entry := range testEntries {
		suspicious := parser.looksLikeRecentEntry(entry, "v3.0.0-alpha.9")
		status := "NORMAL"
		if suspicious {
			status = "SUSPICIOUS"
		}
		fmt.Printf("- [%s] %s\n", status, entry)
	}
}

func ExampleChangelogParser_sectionInfo() {
	content := `# Changelog

## [Unreleased]

### Added
- Future feature

## v2.0.0 - 2025-02-01

### Breaking Changes
- Major API change

### Added
- New feature

## v1.0.0 - 2025-01-01

### Added
- Initial release
`

	parser := NewChangelogParser()
	parser.ParseString(content)

	fmt.Printf("Changelog sections:\n")
	sections := parser.GetSections()

	// Sort sections for consistent output
	for version, section := range sections {
		fmt.Printf("\n%s:\n", version)
		fmt.Printf("  - Title: %s\n", section.Title)
		fmt.Printf("  - Date: %s\n", section.Date)
		fmt.Printf("  - Released: %v\n", section.IsReleased)
		fmt.Printf("  - Categories: %d\n", len(section.Categories))

		for categoryName, category := range section.Categories {
			fmt.Printf("    - %s: %d items\n", categoryName, len(category.Items))
		}
	}
}
