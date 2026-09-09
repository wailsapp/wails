package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run validate-changelog.go <changelog-file> <added-lines-file> [deleted-lines-file] [base-changelog-file]")
		os.Exit(1)
	}

	changelogPath := os.Args[1]
	addedLinesPath := os.Args[2]

	// Read changelog
	content, err := readFile(changelogPath)
	if err != nil {
		fmt.Printf("ERROR: Failed to read changelog: %v\n", err)
		os.Exit(1)
	}

	// Read the lines added in this PR
	addedContent, err := readFile(addedLinesPath)
	if err != nil {
		fmt.Printf("ERROR: Failed to read PR added lines: %v\n", err)
		os.Exit(1)
	}

	addedLines := strings.Split(addedContent, "\n")
	fmt.Printf("📝 Lines added in this PR: %d\n", len(addedLines))
	var deletedLines []string
	if len(os.Args) >= 4 {
		deletedContent, err := readFile(os.Args[3])
		if err != nil {
			fmt.Printf("ERROR: Failed to read deleted changelog lines: %v\n", err)
			os.Exit(1)
		}
		deletedLines = strings.Split(deletedContent, "\n")
		fmt.Printf("📝 Lines deleted in this PR: %d\n", len(deletedLines))
	}
	var deletedEntries []changelogEntry
	if len(os.Args) >= 5 {
		baseContent, err := readFile(os.Args[4])
		if err != nil {
			fmt.Printf("ERROR: Failed to read base changelog: %v\n", err)
			os.Exit(1)
		}
		deletedEntries = deletedChangelogEntries(baseContent, content, deletedLines)
		fmt.Printf("📝 Deleted changelog entries with section metadata: %d\n", len(deletedEntries))
	}

	// Parse changelog to find where added lines ended up
	lines := strings.Split(content, "\n")

	// Find problematic entries - only check lines that were ADDED in this PR
	var issues []Issue
	currentSection := ""

	for lineNum, line := range lines {
		// Track current section
		if section := releaseSection(line); section != "" {
			currentSection = section
		}

		// Check if this line was added in this PR AND is in a released version
		if currentSection != "" && currentSection != "Unreleased" &&
			strings.HasPrefix(strings.TrimSpace(line), "- ") &&
			wasAddedInThisPR(line, addedLines) {
			if isSameSourceCorrection(line, currentSection, deletedEntries) {
				fmt.Printf("✅ CORRECTION: Same-source replacement in %s: %s\n", currentSection, strings.TrimSpace(line))
				continue
			}

			issues = append(issues, Issue{
				Line:     lineNum,
				Content:  strings.TrimSpace(line),
				Section:  currentSection,
				Category: getCurrentCategory(lines, lineNum),
			})
			fmt.Printf("🚨 MISPLACED: Line added to released version %s: %s\n", currentSection, strings.TrimSpace(line))
		}
	}

	if len(issues) == 0 {
		fmt.Println("VALIDATION_RESULT=success")
		fmt.Println("No misplaced changelog entries found ✅")
		return
	}

	// Try to fix the issues
	fmt.Printf("Found %d potentially misplaced entries:\n", len(issues))
	for _, issue := range issues {
		fmt.Printf("  - Line %d in %s: %s\n", issue.Line+1, issue.Section, issue.Content)
	}

	// Attempt automatic fix
	fixed, err := attemptFix(content, issues, changelogPath)
	if err != nil {
		fmt.Printf("VALIDATION_RESULT=error\n")
		fmt.Printf("ERROR: Failed to fix changelog: %v\n", err)
		os.Exit(1)
	}

	if fixed {
		fmt.Println("VALIDATION_RESULT=fixed")
		fmt.Println("✅ Changelog has been automatically fixed")
	} else {
		fmt.Println("VALIDATION_RESULT=cannot_fix")
		fmt.Println("❌ Cannot automatically fix changelog issues")
		os.Exit(1)
	}
}

var pullRequestReference = regexp.MustCompile(`^https://github\.com/wailsapp/wails/pull/[0-9]+$`)

func pullRequestReferenceFromLine(line string) string {
	const linkPrefix = "[PR]("
	start := strings.Index(line, linkPrefix)
	if start == -1 {
		return ""
	}
	destination := line[start+len(linkPrefix):]
	end := strings.IndexByte(destination, ')')
	if end == -1 {
		return ""
	}
	destination = destination[:end]
	if !pullRequestReference.MatchString(destination) {
		return ""
	}
	return destination
}

type changelogEntry struct {
	Line    string
	Section string
}

// deletedChangelogEntries attaches release-section provenance to deleted diff
// lines. Comparing section-specific counts between the base and current files
// prevents a line moved between releases from masquerading as a correction.
func deletedChangelogEntries(baseContent, currentContent string, deletedLines []string) []changelogEntry {
	deletedCounts := make(map[string]int)
	for _, line := range deletedLines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "- ") {
			deletedCounts[line]++
		}
	}

	currentCounts := make(map[changelogEntry]int)
	for _, entry := range changelogEntries(currentContent) {
		currentCounts[entry]++
	}

	var result []changelogEntry
	for _, entry := range changelogEntries(baseContent) {
		if currentCounts[entry] > 0 {
			currentCounts[entry]--
			continue
		}
		if deletedCounts[entry.Line] == 0 {
			continue
		}
		deletedCounts[entry.Line]--
		result = append(result, entry)
	}
	return result
}

func changelogEntries(content string) []changelogEntry {
	var entries []changelogEntry
	currentSection := ""
	for _, line := range strings.Split(content, "\n") {
		if section := releaseSection(line); section != "" {
			currentSection = section
		}
		line = strings.TrimSpace(line)
		if currentSection != "" && strings.HasPrefix(line, "- ") {
			entries = append(entries, changelogEntry{Line: line, Section: currentSection})
		}
	}
	return entries
}

func releaseSection(line string) string {
	if !strings.HasPrefix(line, "## ") {
		return ""
	}
	if strings.Contains(line, "[Unreleased]") {
		return "Unreleased"
	}
	if !strings.Contains(line, "v3.0.0-") {
		return ""
	}
	parts := strings.Split(strings.TrimSpace(line[3:]), " - ")
	return strings.TrimSpace(parts[0])
}

// isSameSourceCorrection distinguishes a historical correction from a new
// entry added to a released section. Both lines must be changelog bullets in
// the same released section and cite the same immutable Wails pull request.
func isSameSourceCorrection(addedLine, addedSection string, deletedEntries []changelogEntry) bool {
	addedLine = strings.TrimSpace(addedLine)
	if !strings.HasPrefix(addedLine, "- ") {
		return false
	}
	reference := pullRequestReferenceFromLine(addedLine)
	if reference == "" {
		return false
	}
	for _, deletedEntry := range deletedEntries {
		if deletedEntry.Section == addedSection &&
			deletedEntry.Line != addedLine &&
			pullRequestReferenceFromLine(deletedEntry.Line) == reference {
			return true
		}
	}
	return false
}

type Issue struct {
	Line     int
	Content  string
	Section  string
	Category string
}

func wasAddedInThisPR(line string, addedLines []string) bool {
	trimmedLine := strings.TrimSpace(line)
	for _, addedLine := range addedLines {
		trimmedAdded := strings.TrimSpace(addedLine)
		if trimmedAdded == trimmedLine {
			return true
		}
		if strings.Contains(trimmedAdded, trimmedLine) && len(trimmedAdded) > 0 {
			return true
		}
	}
	return false
}

func getCurrentCategory(lines []string, lineNum int) string {
	for i := lineNum - 1; i >= 0; i-- {
		line := strings.TrimSpace(lines[i])
		if strings.HasPrefix(line, "### ") {
			return strings.TrimSpace(line[4:])
		}
		if strings.HasPrefix(line, "## ") &&
			!strings.Contains(line, "[Unreleased]") &&
			!strings.Contains(line, "v3.0.0-") {
			return strings.TrimSpace(line[3:])
		}
		if strings.HasPrefix(line, "## ") &&
			(strings.Contains(line, "[Unreleased]") || strings.Contains(line, "v3.0.0-")) {
			break
		}
	}
	return "Added"
}

func attemptFix(content string, issues []Issue, outputPath string) (bool, error) {
	lines := strings.Split(content, "\n")

	// Find unreleased section
	unreleasedStart := -1
	unreleasedEnd := -1

	for i, line := range lines {
		if strings.Contains(line, "[Unreleased]") {
			unreleasedStart = i
			for j := i + 1; j < len(lines); j++ {
				if strings.HasPrefix(lines[j], "## ") && !strings.Contains(lines[j], "[Unreleased]") {
					unreleasedEnd = j
					break
				}
			}
			break
		}
	}

	if unreleasedStart == -1 {
		return false, fmt.Errorf("Could not find [Unreleased] section")
	}

	// Group issues by category
	issuesByCategory := make(map[string][]Issue)
	for _, issue := range issues {
		issuesByCategory[issue.Category] = append(issuesByCategory[issue.Category], issue)
	}

	// Remove issues from original locations (in reverse order)
	var linesToRemove []int
	for _, issue := range issues {
		linesToRemove = append(linesToRemove, issue.Line)
	}

	// Sort in reverse order
	for i := 0; i < len(linesToRemove); i++ {
		for j := i + 1; j < len(linesToRemove); j++ {
			if linesToRemove[i] < linesToRemove[j] {
				linesToRemove[i], linesToRemove[j] = linesToRemove[j], linesToRemove[i]
			}
		}
	}

	// Remove lines
	for _, lineNum := range linesToRemove {
		lines = append(lines[:lineNum], lines[lineNum+1:]...)
	}

	// Add entries to unreleased section
	for category, categoryIssues := range issuesByCategory {
		categoryFound := false
		insertPos := unreleasedStart + 1

		for i := unreleasedStart + 1; i < unreleasedEnd && i < len(lines); i++ {
			if strings.Contains(lines[i], "### "+category) || strings.Contains(lines[i], "## "+category) {
				categoryFound = true
				for j := i + 1; j < unreleasedEnd && j < len(lines); j++ {
					if strings.HasPrefix(lines[j], "### ") || strings.HasPrefix(lines[j], "## ") {
						insertPos = j
						break
					}
					if j == len(lines)-1 || j == unreleasedEnd-1 {
						insertPos = j + 1
						break
					}
				}
				break
			}
		}

		if !categoryFound {
			if unreleasedEnd > 0 {
				insertPos = unreleasedEnd
			} else {
				insertPos = unreleasedStart + 1
			}

			newLines := []string{
				"",
				"### " + category,
				"",
			}
			lines = append(lines[:insertPos], append(newLines, lines[insertPos:]...)...)
			insertPos += len(newLines)
			unreleasedEnd += len(newLines)
		}

		// Add entries to the category
		for _, issue := range categoryIssues {
			lines = append(lines[:insertPos], append([]string{issue.Content}, lines[insertPos:]...)...)
			insertPos++
			unreleasedEnd++
		}
	}

	// Write back to file
	newContent := strings.Join(lines, "\n")
	return true, writeFile(outputPath, newContent)
}

func readFile(path string) (string, error) {
	file, err := os.Open(path)
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

func writeFile(path, content string) error {
	dir := filepath.Dir(path)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	return os.WriteFile(path, []byte(content), 0644)
}
