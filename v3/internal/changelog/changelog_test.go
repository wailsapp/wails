package changelog

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewProcessor(t *testing.T) {
	processor := NewProcessor()
	if processor == nil {
		t.Fatal("NewProcessor() returned nil")
	}
	if processor.parser == nil {
		t.Error("parser not initialized")
	}
	if processor.validator == nil {
		t.Error("validator not initialized")
	}
}

func TestProcessString_ValidContent(t *testing.T) {
	processor := NewProcessor()
	content := `# Unreleased Changes

## Added
- Add support for custom window icons in application options
- Add new SetWindowIcon() method to runtime API (#1234)

## Fixed
- Fix memory leak in event system during window close operations (#5678)`

	result, err := processor.ProcessString(content)
	if err != nil {
		t.Fatalf("ProcessString() returned error: %v", err)
	}

	if !result.HasContent {
		t.Error("Result should have content")
	}

	if !result.ValidationResult.Valid {
		t.Errorf("Validation should pass, got errors: %s", result.ValidationResult.GetValidationSummary())
	}

	// Check parsed content
	if len(result.Entry.Added) != 2 {
		t.Errorf("Expected 2 Added items, got %d", len(result.Entry.Added))
	}
	if len(result.Entry.Fixed) != 1 {
		t.Errorf("Expected 1 Fixed item, got %d", len(result.Entry.Fixed))
	}
}

func TestProcessString_InvalidContent(t *testing.T) {
	processor := NewProcessor()
	content := `# Unreleased Changes

## Added
- Short
- TODO: add proper description

## InvalidSection
- This section is invalid`

	result, err := processor.ProcessString(content)
	if err != nil {
		t.Fatalf("ProcessString() returned error: %v", err)
	}

	if result.ValidationResult.Valid {
		t.Error("Validation should fail for invalid content")
	}

	// The parser will parse the Added section correctly (2 items)
	// The InvalidSection won't be parsed since it's not a valid section name
	if len(result.Entry.Added) != 2 {
		t.Errorf("Expected 2 Added items, got %d", len(result.Entry.Added))
	}

	// Should have validation errors
	if len(result.ValidationResult.Errors) == 0 {
		t.Error("Should have validation errors")
	}
}

func TestProcessString_EmptyContent(t *testing.T) {
	processor := NewProcessor()
	content := `# Unreleased Changes

## Added
<!-- No content -->

## Changed
<!-- No content -->`

	result, err := processor.ProcessString(content)
	if err != nil {
		t.Fatalf("ProcessString() returned error: %v", err)
	}

	if result.HasContent {
		t.Error("Result should not have content")
	}

	if result.ValidationResult.Valid {
		t.Error("Validation should fail for empty content")
	}
}

func TestProcessFile_ValidFile(t *testing.T) {
	processor := NewProcessor()

	// Create a temporary file
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test_changelog.md")

	content := `# Unreleased Changes

## Added
- Add support for custom window icons in application options
- Add new SetWindowIcon() method to runtime API (#1234)

## Fixed
- Fix memory leak in event system during window close operations (#5678)`

	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	result, err := processor.ProcessFile(filePath)
	if err != nil {
		t.Fatalf("ProcessFile() returned error: %v", err)
	}

	if !result.HasContent {
		t.Error("Result should have content")
	}

	if !result.ValidationResult.Valid {
		t.Errorf("Validation should pass, got errors: %s", result.ValidationResult.GetValidationSummary())
	}

	// Check parsed content
	if len(result.Entry.Added) != 2 {
		t.Errorf("Expected 2 Added items, got %d", len(result.Entry.Added))
	}
	if len(result.Entry.Fixed) != 1 {
		t.Errorf("Expected 1 Fixed item, got %d", len(result.Entry.Fixed))
	}
}

func TestProcessFile_NonexistentFile(t *testing.T) {
	processor := NewProcessor()

	result, err := processor.ProcessFile("/nonexistent/file.md")
	if err == nil {
		t.Error("ProcessFile() should return error for nonexistent file")
	}
	if result != nil {
		t.Error("ProcessFile() should return nil result for nonexistent file")
	}
}

func TestValidateString_ValidContent(t *testing.T) {
	processor := NewProcessor()
	content := `# Unreleased Changes

## Added
- Add support for custom window icons in application options

## Fixed
- Fix memory leak in event system during window close operations`

	result := processor.ValidateString(content)
	if !result.Valid {
		t.Errorf("ValidateString() should be valid, got errors: %s", result.GetValidationSummary())
	}
}

func TestValidateString_InvalidContent(t *testing.T) {
	processor := NewProcessor()
	content := `# Unreleased Changes

## InvalidSection
- This section is invalid

- Bullet point outside section`

	result := processor.ValidateString(content)
	if result.Valid {
		t.Error("ValidateString() should be invalid")
	}

	if len(result.Errors) == 0 {
		t.Error("ValidateString() should return validation errors")
	}
}

func TestValidateFile_ValidFile(t *testing.T) {
	processor := NewProcessor()

	// Create a temporary file
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test_changelog.md")

	content := `# Unreleased Changes

## Added
- Add support for custom window icons in application options

## Fixed
- Fix memory leak in event system during window close operations`

	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	result, err := processor.ValidateFile(filePath)
	if err != nil {
		t.Fatalf("ValidateFile() returned error: %v", err)
	}

	if !result.Valid {
		t.Errorf("ValidateFile() should be valid, got errors: %s", result.GetValidationSummary())
	}
}

func TestValidateFile_NonexistentFile(t *testing.T) {
	processor := NewProcessor()

	result, err := processor.ValidateFile("/nonexistent/file.md")
	if err == nil {
		t.Error("ValidateFile() should return error for nonexistent file")
	}
	if result.Valid {
		t.Error("ValidateFile() should return invalid result for nonexistent file")
	}
}

func TestProcessorResult_Integration(t *testing.T) {
	processor := NewProcessor()
	content := `# Unreleased Changes

## Added
- Add comprehensive changelog processing system
- Add validation for Keep a Changelog format compliance

## Changed
- Update changelog workflow to use automated processing

## Fixed
- Fix parsing issues with various markdown bullet styles
- Fix validation edge cases for empty content sections`

	result, err := processor.ProcessString(content)
	if err != nil {
		t.Fatalf("ProcessString() returned error: %v", err)
	}

	// Test that we can format the result for different outputs
	changelogFormat := result.Entry.FormatForChangelog()
	if !strings.Contains(changelogFormat, "### Added") {
		t.Error("FormatForChangelog() should contain Added section")
	}
	if !strings.Contains(changelogFormat, "### Changed") {
		t.Error("FormatForChangelog() should contain Changed section")
	}
	if !strings.Contains(changelogFormat, "### Fixed") {
		t.Error("FormatForChangelog() should contain Fixed section")
	}

	releaseFormat := result.Entry.FormatForRelease()
	if !strings.Contains(releaseFormat, "## ✨ Added") {
		t.Error("FormatForRelease() should contain Added section with emoji")
	}
	if !strings.Contains(releaseFormat, "## 🔄 Changed") {
		t.Error("FormatForRelease() should contain Changed section with emoji")
	}
	if !strings.Contains(releaseFormat, "## 🐛 Fixed") {
		t.Error("FormatForRelease() should contain Fixed section with emoji")
	}

	// Test validation summary
	if !result.ValidationResult.Valid {
		t.Errorf("Validation should pass, got: %s", result.ValidationResult.GetValidationSummary())
	}

	summary := result.ValidationResult.GetValidationSummary()
	if !strings.Contains(summary, "Validation passed") {
		t.Errorf("Validation summary should indicate success, got: %s", summary)
	}
}
