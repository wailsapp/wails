# Changelog Parser

A Go package for parsing and validating Keep a Changelog format files, with automatic detection and correction of misplaced entries.

## Overview

This package helps maintain proper changelog structure by:
- Parsing Keep a Changelog format files
- Detecting entries added to already-released versions
- Automatically moving misplaced entries to the `[Unreleased]` section
- Maintaining proper category ordering and structure

## Usage

### Basic Usage

```go
parser := changelog.NewChangelogParser()
err := parser.ParseString(changelogContent)
if err != nil {
    log.Fatal(err)
}

// Validate and fix misplaced entries
result, err := parser.ValidateAndFixMisplacedEntries()
if err != nil {
    log.Fatal(err)
}

if len(result.MisplacedEntries) > 0 {
    fmt.Printf("Found %d misplaced entries\n", len(result.MisplacedEntries))
    
    // Generate corrected changelog
    corrected := parser.GenerateChangelog()
    correctedContent := strings.Join(corrected, "\n")
    
    // Write corrected content back to file
    err = os.WriteFile("CHANGELOG.md", []byte(correctedContent), 0644)
}
```

### Integration with GitHub Actions

This package is used by the `v3-check-changelog` GitHub Actions workflow to automatically validate and fix changelog entries in pull requests.

The workflow:
1. Triggers on PRs that modify `docs/src/content/docs/changelog.mdx`
2. Detects entries added to already-released versions
3. Automatically moves them to the `[Unreleased]` section
4. Commits the fix back to the PR branch
5. Comments on the PR with the result

## Detection Heuristics

The parser uses several heuristics to detect misplaced entries:

1. **Pattern Matching**: Looks for suspicious patterns that indicate recent additions
2. **Date Analysis**: Compares entry dates with version release dates
3. **Duplicate Detection**: Identifies entries that appear in multiple versions

### Suspicious Patterns

- Recent contributor names in old releases
- PR/issue numbers that seem too high for old versions
- Specific text patterns known to be recent additions

## Keep a Changelog Format

The parser supports the standard Keep a Changelog format:

```markdown
# Changelog

## [Unreleased]

### Added
- New features go here

### Fixed
- Bug fixes go here

## v1.1.0 - 2025-01-15

### Added
- Feature that was released in v1.1.0

## v1.0.0 - 2025-01-01

### Added
- Initial release
```

### Supported Categories

- `Breaking Changes`
- `Added`
- `Changed` 
- `Deprecated`
- `Removed`
- `Fixed`
- `Security`

Categories are automatically ordered according to Keep a Changelog standards.

## API Reference

### Types

#### `ChangelogParser`
Main parser struct for handling changelog operations.

#### `ValidationResult`
Contains validation results including errors, warnings, and misplaced entries.

#### `Section`
Represents a changelog section (version or unreleased).

#### `Category` 
Represents a category within a section (Added, Fixed, etc.).

#### `Entry`
Represents a single changelog entry.

### Methods

#### `NewChangelogParser() *ChangelogParser`
Creates a new changelog parser instance.

#### `ParseString(content string) error`
Parses changelog content from a string.

#### `ParseContent(lines []string) error`
Parses changelog content from a slice of strings.

#### `ValidateAndFixMisplacedEntries() (*ValidationResult, error)`
Validates the changelog and fixes any misplaced entries.

#### `GenerateChangelog() []string`
Generates corrected changelog content as a slice of strings.

#### `GetSections() map[string]*Section`
Returns all parsed sections.

#### `GetUnreleasedSection() *Section`
Returns the unreleased section if it exists.

## Testing

Run the test suite:

```bash
go test -v
```

Run benchmarks:

```bash
go test -bench=.
```

## GitHub Actions Integration

The package integrates with the `v3-check-changelog` GitHub Actions workflow:

```yaml
name: V3 Changelog Validator
on:
  pull_request:
    branches: [ v3-alpha ]
    paths:
      - 'docs/src/content/docs/changelog.mdx'
```

The workflow automatically:
- ✅ Validates changelog structure
- 🔧 Fixes misplaced entries
- 📝 Commits corrections to the PR
- 💬 Comments on PR with results
- ❌ Fails if manual intervention is needed

## Contributing

When adding new detection heuristics or fixing bugs:

1. Add test cases in `parser_test.go`
2. Update the detection logic in `detectMisplacedEntries()`
3. Ensure all tests pass
4. Update this README if adding new features