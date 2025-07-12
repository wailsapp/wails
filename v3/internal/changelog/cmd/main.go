package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <command> [changelog-file]\n", os.Args[0])
		fmt.Printf("Commands:\n")
		fmt.Printf("  validate <file>  - Validate changelog and detect misplaced entries\n")
		fmt.Printf("  fix <file>       - Fix misplaced entries and output corrected changelog\n")
		fmt.Printf("  demo             - Run with demo data\n")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "demo":
		runDemo()
	case "validate":
		if len(os.Args) < 3 {
			fmt.Println("Error: Please provide changelog file path")
			os.Exit(1)
		}
		validateChangelog(os.Args[2])
	case "fix":
		if len(os.Args) < 3 {
			fmt.Println("Error: Please provide changelog file path")
			os.Exit(1)
		}
		fixChangelog(os.Args[2])
	default:
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}
}

func runDemo() {
	fmt.Println("=== Changelog Parser Demo ===\n")

	// Demo data with misplaced entries
	demoContent := `# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## v3.0.0-alpha.11 - 2025-07-12

### Added
- Add distribution-specific build dependencies for Linux by @leaanthony in [PR](https://github.com/wailsapp/wails/pull/4345)
- Added bindings guide by @atterpac in [PR](https://github.com/wailsapp/wails/pull/4404)
- Legitimate feature for alpha.11

### Fixed
- Bug fix for alpha.11

## v3.0.0-alpha.10 - 2025-07-06

### Added
- Original feature for alpha.10
- Another legitimate feature
- Add distribution-specific build dependencies for Linux by @leaanthony in [PR](https://github.com/wailsapp/wails/pull/4345)

### Fixed
- Some bug fix
- Added bindings guide by @atterpac in [PR](https://github.com/wailsapp/wails/pull/4404)

## v3.0.0-alpha.9 - 2025-01-13

### Added
- Old feature from January
`

	fmt.Println("Demo changelog content:")
	fmt.Println(strings.Repeat("-", 50))
	fmt.Println(demoContent)
	fmt.Println(strings.Repeat("-", 50))

	parser := changelog.NewChangelogParser()
	err := parser.ParseString(demoContent)
	if err != nil {
		fmt.Printf("Error parsing changelog: %v\n", err)
		return
	}

	fmt.Printf("\n=== Validation Results ===\n")

	result, err := parser.ValidateAndFixMisplacedEntries()
	if err != nil {
		fmt.Printf("Error validating: %v\n", err)
		return
	}

	fmt.Printf("Valid: %v\n", result.IsValid)
	fmt.Printf("Errors: %d\n", len(result.Errors))
	fmt.Printf("Warnings: %d\n", len(result.Warnings))
	fmt.Printf("Misplaced entries found: %d\n", len(result.MisplacedEntries))

	if len(result.Warnings) > 0 {
		fmt.Printf("\nWarnings:\n")
		for _, warning := range result.Warnings {
			fmt.Printf("  - %s\n", warning)
		}
	}

	if len(result.MisplacedEntries) > 0 {
		fmt.Printf("\nMisplaced entries detected:\n")
		for _, entry := range result.MisplacedEntries {
			fmt.Printf("  - [%s] %s (from section: %s)\n",
				entry.Category, entry.Text, entry.Section)
		}
	}

	fmt.Printf("\n=== Corrected Changelog ===\n")
	corrected := parser.GenerateChangelog()
	for _, line := range corrected {
		fmt.Println(line)
	}
}

func validateChangelog(filePath string) {
	content, err := readFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	parser := changelog.NewChangelogParser()
	err = parser.ParseString(content)
	if err != nil {
		fmt.Printf("Error parsing changelog: %v\n", err)
		os.Exit(1)
	}

	result, err := parser.ValidateAndFixMisplacedEntries()
	if err != nil {
		fmt.Printf("Error validating: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Validation Results for %s:\n", filePath)
	fmt.Printf("=========================\n")
	fmt.Printf("Valid: %v\n", result.IsValid)
	fmt.Printf("Errors: %d\n", len(result.Errors))
	fmt.Printf("Warnings: %d\n", len(result.Warnings))
	fmt.Printf("Misplaced entries: %d\n", len(result.MisplacedEntries))

	if len(result.Errors) > 0 {
		fmt.Printf("\nErrors:\n")
		for _, error := range result.Errors {
			fmt.Printf("  - %s\n", error)
		}
	}

	if len(result.Warnings) > 0 {
		fmt.Printf("\nWarnings:\n")
		for _, warning := range result.Warnings {
			fmt.Printf("  - %s\n", warning)
		}
	}

	if len(result.MisplacedEntries) > 0 {
		fmt.Printf("\nMisplaced entries:\n")
		for _, entry := range result.MisplacedEntries {
			fmt.Printf("  - [%s] %s\n    From: %s\n",
				entry.Category, entry.Text, entry.Section)
		}

		fmt.Printf("\nSuggestion: Run 'fix %s' to automatically correct these issues.\n", filePath)
	} else {
		fmt.Printf("\n✅ No misplaced entries detected!\n")
	}
}

func fixChangelog(filePath string) {
	content, err := readFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	parser := changelog.NewChangelogParser()
	err = parser.ParseString(content)
	if err != nil {
		fmt.Printf("Error parsing changelog: %v\n", err)
		os.Exit(1)
	}

	result, err := parser.ValidateAndFixMisplacedEntries()
	if err != nil {
		fmt.Printf("Error validating: %v\n", err)
		os.Exit(1)
	}

	if len(result.MisplacedEntries) == 0 {
		fmt.Printf("No misplaced entries found in %s\n", filePath)
		return
	}

	fmt.Printf("Fixed %d misplaced entries in %s\n", len(result.MisplacedEntries), filePath)

	// Generate corrected content
	corrected := parser.GenerateChangelog()
	correctedContent := strings.Join(corrected, "\n")

	// Create backup
	backupPath := filePath + ".backup"
	err = writeFile(backupPath, content)
	if err != nil {
		fmt.Printf("Warning: Could not create backup: %v\n", err)
	} else {
		fmt.Printf("Backup created: %s\n", backupPath)
	}

	// Write corrected content
	err = writeFile(filePath, correctedContent)
	if err != nil {
		fmt.Printf("Error writing corrected file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Changelog fixed successfully!\n")

	// Show summary
	fmt.Printf("\nSummary of changes:\n")
	for _, entry := range result.MisplacedEntries {
		fmt.Printf("  - Moved '%s' from %s to Unreleased\n",
			truncateString(entry.Text, 60), entry.Section)
	}
}

func readFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var content strings.Builder
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		content.WriteString(scanner.Text())
		content.WriteString("\n")
	}

	return content.String(), scanner.Err()
}

func writeFile(filePath, content string) error {
	dir := filepath.Dir(filePath)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, []byte(content), 0644)
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
