package changelog

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"
)

// ChangelogEntry represents a parsed changelog entry following Keep a Changelog format
type ChangelogEntry struct {
	Version    string    `json:"version"`
	Date       time.Time `json:"date"`
	Added      []string  `json:"added"`
	Changed    []string  `json:"changed"`
	Fixed      []string  `json:"fixed"`
	Deprecated []string  `json:"deprecated"`
	Removed    []string  `json:"removed"`
	Security   []string  `json:"security"`
}

// Parser handles parsing of UNRELEASED_CHANGELOG.md files
type Parser struct {
	// sectionRegex matches section headers like "## Added", "## Changed", etc.
	sectionRegex *regexp.Regexp
	// bulletRegex matches bullet points (- or *)
	bulletRegex *regexp.Regexp
}

// NewParser creates a new changelog parser
func NewParser() *Parser {
	return &Parser{
		sectionRegex: regexp.MustCompile(`^##\s+(Added|Changed|Fixed|Deprecated|Removed|Security)\s*$`),
		bulletRegex:  regexp.MustCompile(`^[\s]*[-*]\s+(.+)$`),
	}
}

// ParseContent parses changelog content from a reader and returns a ChangelogEntry
func (p *Parser) ParseContent(reader io.Reader) (*ChangelogEntry, error) {
	entry := &ChangelogEntry{
		Added:      []string{},
		Changed:    []string{},
		Fixed:      []string{},
		Deprecated: []string{},
		Removed:    []string{},
		Security:   []string{},
	}

	scanner := bufio.NewScanner(reader)
	var currentSection string
	var inExampleSection bool

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "<!--") || strings.HasPrefix(line, "-->") {
			continue
		}

		// Skip the main title
		if strings.HasPrefix(line, "# Unreleased Changes") {
			continue
		}

		// Check if we're entering the example section
		if strings.HasPrefix(line, "---") || strings.HasPrefix(line, "### Example Entries") {
			inExampleSection = true
			continue
		}

		// Skip example section content
		if inExampleSection {
			continue
		}

		// Check for section headers
		if strings.HasPrefix(line, "##") {
			if matches := p.sectionRegex.FindStringSubmatch(line); len(matches) > 1 {
				currentSection = strings.ToLower(matches[1])
			} else {
				// Invalid section header - reset current section
				currentSection = ""
			}
			continue
		}

		// Parse bullet points
		if matches := p.bulletRegex.FindStringSubmatch(line); len(matches) > 1 {
			content := strings.TrimSpace(matches[1])
			if content == "" {
				continue
			}

			switch currentSection {
			case "added":
				entry.Added = append(entry.Added, content)
			case "changed":
				entry.Changed = append(entry.Changed, content)
			case "fixed":
				entry.Fixed = append(entry.Fixed, content)
			case "deprecated":
				entry.Deprecated = append(entry.Deprecated, content)
			case "removed":
				entry.Removed = append(entry.Removed, content)
			case "security":
				entry.Security = append(entry.Security, content)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading changelog content: %w", err)
	}

	return entry, nil
}

// HasContent checks if the changelog entry contains any actual content
func (entry *ChangelogEntry) HasContent() bool {
	return len(entry.Added) > 0 ||
		len(entry.Changed) > 0 ||
		len(entry.Fixed) > 0 ||
		len(entry.Deprecated) > 0 ||
		len(entry.Removed) > 0 ||
		len(entry.Security) > 0
}

// FormatForChangelog formats the entry for insertion into the main changelog
func (entry *ChangelogEntry) FormatForChangelog() string {
	var builder strings.Builder

	if len(entry.Added) > 0 {
		builder.WriteString("### Added\n")
		for _, item := range entry.Added {
			builder.WriteString(fmt.Sprintf("- %s\n", item))
		}
		builder.WriteString("\n")
	}

	if len(entry.Changed) > 0 {
		builder.WriteString("### Changed\n")
		for _, item := range entry.Changed {
			builder.WriteString(fmt.Sprintf("- %s\n", item))
		}
		builder.WriteString("\n")
	}

	if len(entry.Fixed) > 0 {
		builder.WriteString("### Fixed\n")
		for _, item := range entry.Fixed {
			builder.WriteString(fmt.Sprintf("- %s\n", item))
		}
		builder.WriteString("\n")
	}

	if len(entry.Deprecated) > 0 {
		builder.WriteString("### Deprecated\n")
		for _, item := range entry.Deprecated {
			builder.WriteString(fmt.Sprintf("- %s\n", item))
		}
		builder.WriteString("\n")
	}

	if len(entry.Removed) > 0 {
		builder.WriteString("### Removed\n")
		for _, item := range entry.Removed {
			builder.WriteString(fmt.Sprintf("- %s\n", item))
		}
		builder.WriteString("\n")
	}

	if len(entry.Security) > 0 {
		builder.WriteString("### Security\n")
		for _, item := range entry.Security {
			builder.WriteString(fmt.Sprintf("- %s\n", item))
		}
		builder.WriteString("\n")
	}

	return strings.TrimSpace(builder.String())
}

// FormatForRelease formats the entry for GitHub release notes
func (entry *ChangelogEntry) FormatForRelease() string {
	var builder strings.Builder

	if len(entry.Added) > 0 {
		builder.WriteString("## ✨ Added\n")
		for _, item := range entry.Added {
			builder.WriteString(fmt.Sprintf("- %s\n", item))
		}
		builder.WriteString("\n")
	}

	if len(entry.Changed) > 0 {
		builder.WriteString("## 🔄 Changed\n")
		for _, item := range entry.Changed {
			builder.WriteString(fmt.Sprintf("- %s\n", item))
		}
		builder.WriteString("\n")
	}

	if len(entry.Fixed) > 0 {
		builder.WriteString("## 🐛 Fixed\n")
		for _, item := range entry.Fixed {
			builder.WriteString(fmt.Sprintf("- %s\n", item))
		}
		builder.WriteString("\n")
	}

	if len(entry.Deprecated) > 0 {
		builder.WriteString("## ⚠️ Deprecated\n")
		for _, item := range entry.Deprecated {
			builder.WriteString(fmt.Sprintf("- %s\n", item))
		}
		builder.WriteString("\n")
	}

	if len(entry.Removed) > 0 {
		builder.WriteString("## 🗑️ Removed\n")
		for _, item := range entry.Removed {
			builder.WriteString(fmt.Sprintf("- %s\n", item))
		}
		builder.WriteString("\n")
	}

	if len(entry.Security) > 0 {
		builder.WriteString("## 🔒 Security\n")
		for _, item := range entry.Security {
			builder.WriteString(fmt.Sprintf("- %s\n", item))
		}
		builder.WriteString("\n")
	}

	return strings.TrimSpace(builder.String())
}
