package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run validate-changelog.go <changelog-file> <added-lines-file>")
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

	// Parse changelog to find where added lines ended up
	lines := strings.Split(content, "\n")

	// Find problematic entries - only check lines that were ADDED in this PR
	var issues []Issue
	currentSection := ""

	for lineNum, line := range lines {
		// Track current section
		if strings.HasPrefix(line, "## ") {
			if strings.Contains(line, "[Unreleased]") {
				currentSection = "Unreleased"
			} else if strings.Contains(line, "v3.0.0-alpha") {
				// Extract version from line like "## v3.0.0-alpha.10 - 2025-07-06"
				parts := strings.Split(strings.TrimSpace(line[3:]), " - ")
				if len(parts) >= 1 {
					currentSection = strings.TrimSpace(parts[0])
				}
			}
		}

		// Check if this line was added in this PR AND is in a released version
		if currentSection != "" && currentSection != "Unreleased" &&
			strings.HasPrefix(strings.TrimSpace(line), "- ") &&
			wasAddedInThisPR(line, addedLines) {

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
			!strings.Contains(line, "v3.0.0-alpha") {
			return strings.TrimSpace(line[3:])
		}
		if strings.HasPrefix(line, "## ") &&
			(strings.Contains(line, "[Unreleased]") || strings.Contains(line, "v3.0.0-alpha")) {
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