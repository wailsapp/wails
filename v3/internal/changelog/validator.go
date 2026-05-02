package changelog

import (
	"fmt"
	"regexp"
	"strings"
)

// ValidationError represents a validation error with context
type ValidationError struct {
	Field   string
	Message string
	Line    int
}

func (e ValidationError) Error() string {
	if e.Line > 0 {
		return fmt.Sprintf("validation error at line %d in %s: %s", e.Line, e.Field, e.Message)
	}
	return fmt.Sprintf("validation error in %s: %s", e.Field, e.Message)
}

// ValidationResult contains the results of validation
type ValidationResult struct {
	Valid  bool
	Errors []ValidationError
}

// Validator handles validation of changelog content and entries
type Validator struct {
	// Regex patterns for validation
	sectionHeaderRegex *regexp.Regexp
	bulletPointRegex   *regexp.Regexp
	urlRegex           *regexp.Regexp
	issueRefRegex      *regexp.Regexp
}

// NewValidator creates a new changelog validator
func NewValidator() *Validator {
	return &Validator{
		sectionHeaderRegex: regexp.MustCompile(`^##\s+(Added|Changed|Fixed|Deprecated|Removed|Security)\s*$`),
		bulletPointRegex:   regexp.MustCompile(`^[\s]*[-*]\s+(.+)$`),
		urlRegex:           regexp.MustCompile(`https?://[^\s]+`),
		issueRefRegex:      regexp.MustCompile(`#\d+`),
	}
}

// ValidateContent validates raw changelog content for proper formatting
func (v *Validator) ValidateContent(content string) ValidationResult {
	result := ValidationResult{
		Valid:  true,
		Errors: []ValidationError{},
	}

	lines := strings.Split(content, "\n")
	var currentSection string
	var hasValidSections bool
	var inExampleSection bool
	lineNum := 0

	for _, line := range lines {
		lineNum++
		trimmedLine := strings.TrimSpace(line)

		// Skip empty lines and comments
		if trimmedLine == "" || strings.HasPrefix(trimmedLine, "<!--") || strings.HasPrefix(trimmedLine, "-->") {
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

		// Check for section headers
		if strings.HasPrefix(trimmedLine, "##") {
			if matches := v.sectionHeaderRegex.FindStringSubmatch(trimmedLine); len(matches) > 1 {
				currentSection = strings.ToLower(matches[1])
				hasValidSections = true
			} else {
				result.Valid = false
				result.Errors = append(result.Errors, ValidationError{
					Field:   "section_header",
					Message: fmt.Sprintf("invalid section header format: '%s'. Expected format: '## SectionName'", trimmedLine),
					Line:    lineNum,
				})
			}
			continue
		}

		// Check bullet points
		if strings.HasPrefix(trimmedLine, "-") || strings.HasPrefix(trimmedLine, "*") {
			if currentSection == "" {
				result.Valid = false
				result.Errors = append(result.Errors, ValidationError{
					Field:   "bullet_point",
					Message: "bullet point found outside of any section",
					Line:    lineNum,
				})
				continue
			}

			// Check for empty bullet points first (just "-" or "*" with optional whitespace)
			if trimmedLine == "-" || trimmedLine == "*" || strings.TrimSpace(trimmedLine[1:]) == "" {
				result.Valid = false
				result.Errors = append(result.Errors, ValidationError{
					Field:   "bullet_point",
					Message: "empty bullet point content",
					Line:    lineNum,
				})
				continue
			}

			if matches := v.bulletPointRegex.FindStringSubmatch(trimmedLine); len(matches) > 1 {
				content := strings.TrimSpace(matches[1])
				if content == "" {
					result.Valid = false
					result.Errors = append(result.Errors, ValidationError{
						Field:   "bullet_point",
						Message: "empty bullet point content",
						Line:    lineNum,
					})
				} else {
					// Validate bullet point content
					v.validateBulletPointContent(content, lineNum, &result)
				}
			} else {
				result.Valid = false
				result.Errors = append(result.Errors, ValidationError{
					Field:   "bullet_point",
					Message: fmt.Sprintf("malformed bullet point: '%s'", trimmedLine),
					Line:    lineNum,
				})
			}
			continue
		}

		// Check for unexpected content
		if trimmedLine != "" && !strings.HasPrefix(trimmedLine, "<!--") {
			// Allow certain patterns like horizontal rules or section comments
			if !strings.HasPrefix(trimmedLine, "---") &&
				!strings.HasPrefix(trimmedLine, "###") &&
				!strings.HasPrefix(trimmedLine, "**") {
				result.Valid = false
				result.Errors = append(result.Errors, ValidationError{
					Field:   "content",
					Message: fmt.Sprintf("unexpected content outside of sections: '%s'", trimmedLine),
					Line:    lineNum,
				})
			}
		}
	}

	// Check if we have at least some valid sections
	if !hasValidSections {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Field:   "structure",
			Message: "no valid changelog sections found",
			Line:    0,
		})
	}

	return result
}

// ValidateEntry validates a parsed changelog entry
func (v *Validator) ValidateEntry(entry *ChangelogEntry) ValidationResult {
	result := ValidationResult{
		Valid:  true,
		Errors: []ValidationError{},
	}

	// Check if entry has any content
	if !entry.HasContent() {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Field:   "content",
			Message: "changelog entry has no content",
			Line:    0,
		})
		return result
	}

	// Validate each section
	v.validateSection("Added", entry.Added, &result)
	v.validateSection("Changed", entry.Changed, &result)
	v.validateSection("Fixed", entry.Fixed, &result)
	v.validateSection("Deprecated", entry.Deprecated, &result)
	v.validateSection("Removed", entry.Removed, &result)
	v.validateSection("Security", entry.Security, &result)

	return result
}

// validateSection validates items in a specific section
func (v *Validator) validateSection(sectionName string, items []string, result *ValidationResult) {
	for i, item := range items {
		if strings.TrimSpace(item) == "" {
			result.Valid = false
			result.Errors = append(result.Errors, ValidationError{
				Field:   sectionName,
				Message: fmt.Sprintf("empty item at index %d", i),
				Line:    0,
			})
			continue
		}

		// Validate item content
		v.validateBulletPointContent(item, 0, result)
	}
}

// validateBulletPointContent validates the content of a bullet point
func (v *Validator) validateBulletPointContent(content string, lineNum int, result *ValidationResult) {
	// Check for common formatting issues
	if strings.HasSuffix(content, ".") && !v.urlRegex.MatchString(content) && !v.issueRefRegex.MatchString(content) {
		// Allow periods in URLs and issue references, but warn about other cases
		// This is a soft warning, not a hard error
	}

	// Check for very short descriptions (likely incomplete)
	if len(strings.TrimSpace(content)) < 10 {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Field:   "bullet_content",
			Message: fmt.Sprintf("bullet point content too short (less than 10 characters): '%s'", content),
			Line:    lineNum,
		})
	}

	// Check for placeholder text
	placeholders := []string{
		"TODO",
		"FIXME",
		"TBD",
		"placeholder",
		"example",
		"sample",
	}

	lowerContent := strings.ToLower(content)
	for _, placeholder := range placeholders {
		if strings.Contains(lowerContent, placeholder) {
			result.Valid = false
			result.Errors = append(result.Errors, ValidationError{
				Field:   "bullet_content",
				Message: fmt.Sprintf("bullet point contains placeholder text: '%s'", content),
				Line:    lineNum,
			})
			break
		}
	}

	// Check for proper capitalization (should start with capital letter)
	if len(content) > 0 {
		firstChar := content[0]
		if firstChar >= 'a' && firstChar <= 'z' {
			// This is a soft warning - we don't fail validation but note it
			// Could be added as a warning system in the future
		}
	}
}

// ValidateRequiredSections checks if the changelog has the minimum required sections
func (v *Validator) ValidateRequiredSections(entry *ChangelogEntry) ValidationResult {
	result := ValidationResult{
		Valid:  true,
		Errors: []ValidationError{},
	}

	// For now, we don't require specific sections to have content
	// This allows flexibility in what developers include
	// But we could add stricter requirements in the future

	// Example of stricter validation (commented out):
	/*
		if len(entry.Added) == 0 && len(entry.Changed) == 0 && len(entry.Fixed) == 0 {
			result.Valid = false
			result.Errors = append(result.Errors, ValidationError{
				Field:   "required_sections",
				Message: "at least one of Added, Changed, or Fixed sections must have content",
				Line:    0,
			})
		}
	*/

	return result
}

// GetValidationSummary returns a human-readable summary of validation results
func (result *ValidationResult) GetValidationSummary() string {
	if result.Valid {
		return "Validation passed: changelog content is properly formatted"
	}

	var summary strings.Builder
	summary.WriteString(fmt.Sprintf("Validation failed with %d error(s):\n", len(result.Errors)))

	for i, err := range result.Errors {
		summary.WriteString(fmt.Sprintf("  %d. %s\n", i+1, err.Error()))
	}

	return summary.String()
}
