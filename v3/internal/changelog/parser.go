package changelog

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"
)

// ChangelogParser handles parsing and manipulation of Keep a Changelog format files
type ChangelogParser struct {
	content   []string
	sections  map[string]*Section
	rawLines  []string
	headerEnd int
}

// Section represents a changelog section
type Section struct {
	Title        string
	Version      string
	Date         string
	StartLine    int
	EndLine      int
	Categories   map[string]*Category
	IsReleased   bool
	IsUnreleased bool
}

// Category represents a changelog category (Added, Changed, Fixed, etc.)
type Category struct {
	Name      string
	StartLine int
	EndLine   int
	Items     []string
}

// Entry represents a single changelog entry
type Entry struct {
	Text     string
	Category string
	Section  string
	Line     int
}

// ValidationResult contains the results of changelog validation
type ValidationResult struct {
	IsValid          bool
	Errors           []string
	Warnings         []string
	MisplacedEntries []Entry
}

// NewChangelogParser creates a new changelog parser
func NewChangelogParser() *ChangelogParser {
	return &ChangelogParser{
		sections: make(map[string]*Section),
	}
}

// ParseContent parses changelog content from a slice of strings
func (p *ChangelogParser) ParseContent(lines []string) error {
	p.rawLines = lines
	p.content = make([]string, len(lines))
	copy(p.content, lines)
	p.sections = make(map[string]*Section)

	return p.parse()
}

// ParseString parses changelog content from a string
func (p *ChangelogParser) ParseString(content string) error {
	lines := strings.Split(content, "\n")
	return p.ParseContent(lines)
}

// parse performs the actual parsing
func (p *ChangelogParser) parse() error {
	versionRegex := regexp.MustCompile(`^##\s+(.+)$`)
	categoryRegex := regexp.MustCompile(`^###?\s+(Added|Changed|Deprecated|Removed|Fixed|Security|Breaking Changes)`)
	entryRegex := regexp.MustCompile(`^[\s]*[-*]\s+(.+)$`)

	var currentSection *Section
	var currentCategory *Category

	for i, line := range p.content {
		line = strings.TrimRight(line, "\r\n")

		// Check for version/section headers
		if matches := versionRegex.FindStringSubmatch(line); matches != nil {
			// Save previous section
			if currentSection != nil {
				currentSection.EndLine = i - 1
				p.sections[currentSection.Version] = currentSection
			}

			// Create new section
			title := strings.TrimSpace(matches[1])
			section := &Section{
				Title:      title,
				StartLine:  i,
				Categories: make(map[string]*Category),
			}

			// Parse version and date
			if strings.Contains(title, "[Unreleased]") || strings.Contains(title, "Unreleased") {
				section.Version = "Unreleased"
				section.IsUnreleased = true
			} else {
				// Parse version and date from title like "v3.0.0-alpha.11 - 2025-07-12"
				parts := strings.Split(title, " - ")
				if len(parts) >= 1 {
					section.Version = strings.TrimSpace(parts[0])
				}
				if len(parts) >= 2 {
					section.Date = strings.TrimSpace(parts[1])
				}
				section.IsReleased = true
			}

			currentSection = section
			currentCategory = nil
			continue
		}

		// Check for category headers
		if matches := categoryRegex.FindStringSubmatch(line); matches != nil && currentSection != nil {
			// Save previous category
			if currentCategory != nil {
				currentCategory.EndLine = i - 1
			}

			categoryName := matches[1]
			category := &Category{
				Name:      categoryName,
				StartLine: i,
				Items:     make([]string, 0),
			}

			currentSection.Categories[categoryName] = category
			currentCategory = category
			continue
		}

		// Check for entries
		if matches := entryRegex.FindStringSubmatch(line); matches != nil && currentCategory != nil {
			entry := strings.TrimSpace(matches[1])
			currentCategory.Items = append(currentCategory.Items, entry)
			continue
		}

		// Track where the header ends (for inserting unreleased section if needed)
		if p.headerEnd == 0 && strings.Contains(line, "## ") {
			p.headerEnd = i
		}
	}

	// Save the last section
	if currentSection != nil {
		currentSection.EndLine = len(p.content) - 1
		if currentCategory != nil {
			currentCategory.EndLine = len(p.content) - 1
		}
		p.sections[currentSection.Version] = currentSection
	}

	return nil
}

// ValidateAndFixMisplacedEntries validates the changelog and detects misplaced entries
func (p *ChangelogParser) ValidateAndFixMisplacedEntries() (*ValidationResult, error) {
	result := &ValidationResult{
		IsValid:          true,
		Errors:           make([]string, 0),
		Warnings:         make([]string, 0),
		MisplacedEntries: make([]Entry, 0),
	}

	// Get all released versions and sort them
	releasedVersions := p.getReleasedVersionsSorted()

	// Check for misplaced entries in released versions
	for _, version := range releasedVersions {
		section := p.sections[version]
		if !section.IsReleased {
			continue
		}

		// Check if this version was recently modified (entries added after release)
		misplacedEntries := p.detectMisplacedEntries(section)
		if len(misplacedEntries) > 0 {
			result.MisplacedEntries = append(result.MisplacedEntries, misplacedEntries...)
			result.Warnings = append(result.Warnings,
				fmt.Sprintf("Found %d potentially misplaced entries in released version %s",
					len(misplacedEntries), version))
		}
	}

	// If we found misplaced entries, we need to fix them
	if len(result.MisplacedEntries) > 0 {
		result.IsValid = false
		err := p.moveMisplacedEntriesToUnreleased(result.MisplacedEntries)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("Failed to move misplaced entries: %v", err))
			return result, err
		}
		result.Warnings = append(result.Warnings, "Moved misplaced entries to Unreleased section")
	}

	return result, nil
}

// detectMisplacedEntries detects entries that might have been added to an already released version
func (p *ChangelogParser) detectMisplacedEntries(section *Section) []Entry {
	var misplaced []Entry

	// This is a heuristic: we detect entries that look like they were recently added
	// In a real implementation, you might compare against a previous version of the file
	// or use git history to detect when entries were added

	for categoryName, category := range section.Categories {
		for _, item := range category.Items {
			// Look for patterns that suggest recent additions
			if p.looksLikeRecentEntry(item, section.Version) {
				misplaced = append(misplaced, Entry{
					Text:     item,
					Category: categoryName,
					Section:  section.Version,
					Line:     category.StartLine, // Approximate line
				})
			}
		}
	}

	return misplaced
}

// looksLikeRecentEntry uses heuristics to determine if an entry looks recently added
func (p *ChangelogParser) looksLikeRecentEntry(entry, version string) bool {
	// Heuristic 1: Check for very recent dates in PR links or commit hashes
	currentYear := time.Now().Year()
	if strings.Contains(entry, fmt.Sprintf("%d", currentYear)) {
		// Check if the version is from a previous year or much earlier
		if p.isVersionOlderThanEntry(version, entry) {
			return true
		}
	}

	// Heuristic 2: Check for PR numbers that seem too high for an old release
	prRegex := regexp.MustCompile(`#(\d+)`)
	if matches := prRegex.FindStringSubmatch(entry); matches != nil {
		// This is a simple heuristic - you might want to make this more sophisticated
		// based on your project's PR numbering patterns
	}

	// Heuristic 3: Check for patterns indicating post-release additions
	suspiciousPatterns := []string{
		"by @leaanthony in [PR]",    // Pattern that might indicate recent manual addition
		"Add distribution-specific", // Specific to the example we saw
		"Added bindings guide",      // Another specific example
	}

	for _, pattern := range suspiciousPatterns {
		if strings.Contains(entry, pattern) {
			return true
		}
	}

	return false
}

// isVersionOlderThanEntry checks if a version seems older than an entry's content suggests
func (p *ChangelogParser) isVersionOlderThanEntry(version, entry string) bool {
	// Extract date from version if available
	if section := p.sections[version]; section != nil && section.Date != "" {
		if entryDate, err := time.Parse("2006-01-02", section.Date); err == nil {
			// If the version is more than a month old but the entry mentions current year, it's suspicious
			if time.Since(entryDate) > 30*24*time.Hour {
				currentYear := time.Now().Year()
				if strings.Contains(entry, fmt.Sprintf("%d", currentYear)) {
					return true
				}
			}
		}
	}
	return false
}

// moveMisplacedEntriesToUnreleased moves misplaced entries to the Unreleased section
func (p *ChangelogParser) moveMisplacedEntriesToUnreleased(entries []Entry) error {
	// Ensure we have an Unreleased section
	p.ensureUnreleasedSection()

	unreleasedSection := p.sections["Unreleased"]

	// Group entries by category
	entriesByCategory := make(map[string][]Entry)
	for _, entry := range entries {
		entriesByCategory[entry.Category] = append(entriesByCategory[entry.Category], entry)
	}

	// Add entries to unreleased section
	for categoryName, categoryEntries := range entriesByCategory {
		// Ensure category exists in unreleased section
		if _, exists := unreleasedSection.Categories[categoryName]; !exists {
			unreleasedSection.Categories[categoryName] = &Category{
				Name:  categoryName,
				Items: make([]string, 0),
			}
		}

		// Add entries to category
		category := unreleasedSection.Categories[categoryName]
		for _, entry := range categoryEntries {
			category.Items = append(category.Items, entry.Text)
		}
	}

	// Remove entries from their original sections
	for _, entry := range entries {
		p.removeEntryFromSection(entry)
	}

	return nil
}

// ensureUnreleasedSection ensures an Unreleased section exists
func (p *ChangelogParser) ensureUnreleasedSection() {
	if _, exists := p.sections["Unreleased"]; !exists {
		p.sections["Unreleased"] = &Section{
			Title:        "[Unreleased]",
			Version:      "Unreleased",
			IsUnreleased: true,
			Categories:   make(map[string]*Category),
		}
	}
}

// removeEntryFromSection removes an entry from its original section
func (p *ChangelogParser) removeEntryFromSection(entry Entry) {
	if section, exists := p.sections[entry.Section]; exists {
		if category, exists := section.Categories[entry.Category]; exists {
			// Remove the entry from the category
			for i, item := range category.Items {
				if item == entry.Text {
					category.Items = append(category.Items[:i], category.Items[i+1:]...)
					break
				}
			}

			// If category is now empty, remove it
			if len(category.Items) == 0 {
				delete(section.Categories, entry.Category)
			}
		}
	}
}

// getReleasedVersionsSorted returns released versions sorted by semver
func (p *ChangelogParser) getReleasedVersionsSorted() []string {
	var versions []string
	for version, section := range p.sections {
		if section.IsReleased {
			versions = append(versions, version)
		}
	}

	// Sort versions (simple string sort for now - could be improved with proper semver)
	sort.Slice(versions, func(i, j int) bool {
		return versions[i] > versions[j] // Newest first
	})

	return versions
}

// GenerateChangelog generates the corrected changelog content
func (p *ChangelogParser) GenerateChangelog() []string {
	var result []string

	// Add header content (everything before first section)
	headerAdded := false
	for _, line := range p.rawLines {
		if strings.HasPrefix(line, "## ") {
			break
		}
		result = append(result, line)
		headerAdded = true
	}

	if !headerAdded {
		// Add default header if none exists
		result = append(result,
			"# Changelog",
			"",
			"All notable changes to this project will be documented in this file.",
			"",
			"The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),",
			"and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).",
			"",
		)
	}

	// Add Unreleased section first if it exists and has content
	if unreleased, exists := p.sections["Unreleased"]; exists && p.sectionHasContent(unreleased) {
		result = append(result, p.generateSectionContent(unreleased)...)
		result = append(result, "")
	}

	// Add released versions in sorted order
	releasedVersions := p.getReleasedVersionsSorted()
	for _, version := range releasedVersions {
		section := p.sections[version]
		result = append(result, p.generateSectionContent(section)...)
		result = append(result, "")
	}

	return result
}

// sectionHasContent checks if a section has any content
func (p *ChangelogParser) sectionHasContent(section *Section) bool {
	for _, category := range section.Categories {
		if len(category.Items) > 0 {
			return true
		}
	}
	return false
}

// generateSectionContent generates the content for a single section
func (p *ChangelogParser) generateSectionContent(section *Section) []string {
	var result []string

	// Add section header
	if section.IsUnreleased {
		result = append(result, "## [Unreleased]")
	} else {
		if section.Date != "" {
			result = append(result, fmt.Sprintf("## %s - %s", section.Version, section.Date))
		} else {
			result = append(result, fmt.Sprintf("## %s", section.Version))
		}
	}

	result = append(result, "")

	// Add categories in standard order
	categoryOrder := []string{"Breaking Changes", "Added", "Changed", "Deprecated", "Removed", "Fixed", "Security"}

	for _, categoryName := range categoryOrder {
		if category, exists := section.Categories[categoryName]; exists && len(category.Items) > 0 {
			if categoryName == "Breaking Changes" {
				result = append(result, "### Breaking Changes")
			} else {
				result = append(result, fmt.Sprintf("### %s", categoryName))
			}
			result = append(result, "")

			// Add items
			for _, item := range category.Items {
				result = append(result, fmt.Sprintf("- %s", item))
			}
			result = append(result, "")
		}
	}

	return result
}

// GetSections returns all parsed sections
func (p *ChangelogParser) GetSections() map[string]*Section {
	return p.sections
}

// GetUnreleasedSection returns the unreleased section if it exists
func (p *ChangelogParser) GetUnreleasedSection() *Section {
	return p.sections["Unreleased"]
}
