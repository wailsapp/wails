package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run validate-changelog.go <changelog-file> <added-lines-file> [deleted-lines-file] [base-changelog-file] [unreleased-file]")
		os.Exit(1)
	}

	changelogPath := os.Args[1]
	addedLinesPath := os.Args[2]
	unreleasedPath := "v3/UNRELEASED_CHANGELOG.md"
	if len(os.Args) >= 6 {
		unreleasedPath = os.Args[5]
	}

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

	// The site contains released history only. Check all newly added entries,
	// including entries accidentally placed before the first release.
	var issues []Issue
	currentSection := ""
	publishedNotes := make(map[string]string)

	for lineNum, line := range lines {
		// Track current section
		if section := releaseSection(line); section != "" {
			currentSection = section
		}

		// Historical corrections and published backfills retain their provenance checks.
		if strings.HasPrefix(strings.TrimSpace(line), "- ") &&
			wasAddedInThisPR(line, addedLines) {
			if isSameSourceCorrection(line, currentSection, deletedEntries) {
				fmt.Printf("✅ CORRECTION: Same-source replacement in %s: %s\n", currentSection, strings.TrimSpace(line))
				continue
			}

			if _, loaded := publishedNotes[currentSection]; !loaded && currentSection != "" && currentSection != "Unreleased" {
				body, err := publishedReleaseNotes(currentSection)
				if err != nil {
					fmt.Printf("⚠️ Cannot verify published notes for %s: %v\n", currentSection, err)
				}
				publishedNotes[currentSection] = body
			}
			if isPublishedBackfill(line, currentSection, publishedNotes) {
				fmt.Printf("✅ BACKFILL: Entry already published in %s: %s\n", currentSection, strings.TrimSpace(line))
				continue
			}

			issues = append(issues, Issue{
				Line:     lineNum,
				Content:  strings.TrimSpace(line),
				Section:  currentSection,
				Category: normalizeCategory(getCurrentCategory(lines, lineNum)),
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
	fixed, err := attemptFix(content, issues, changelogPath, unreleasedPath)
	if err != nil {
		fmt.Printf("VALIDATION_RESULT=error\n")
		fmt.Printf("ERROR: Failed to fix changelog: %v\n", err)
		os.Exit(1)
	}

	if fixed {
		fmt.Println("VALIDATION_RESULT=fixed")
		fmt.Println("✅ Entries moved to " + unreleasedPath)
	} else {
		fmt.Println("VALIDATION_RESULT=cannot_fix")
		fmt.Println("❌ Cannot automatically fix changelog issues")
		os.Exit(1)
	}
}

var pullRequestReference = regexp.MustCompile(`^https://github\.com/wailsapp/wails/pull/[0-9]+$`)
var shortIssueReference = regexp.MustCompile(`\(#([0-9]+)\)$`)

// Multiword prose titles may move out of code spans so they can be translated.
// Preserve every letter and space; do not strip punctuation from code or URLs.
var inlineProseCode = regexp.MustCompile("`([A-Za-z]+(?: [A-Za-z]+)+)`")

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

func changelogReferenceFromLine(line string) string {
	if strings.Contains(line, "[PR](") {
		return pullRequestReferenceFromLine(line)
	}
	if match := shortIssueReference.FindStringSubmatch(strings.TrimSpace(line)); match != nil {
		return "https://github.com/wailsapp/wails/issues/" + match[1]
	}
	return ""
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
// the same released section and either cite the same immutable Wails issue or
// pull request, or differ only in code-span versus emphasis markup around prose titles.
func isSameSourceCorrection(addedLine, addedSection string, deletedEntries []changelogEntry) bool {
	addedLine = strings.TrimSpace(addedLine)
	if !strings.HasPrefix(addedLine, "- ") {
		return false
	}
	reference := changelogReferenceFromLine(addedLine)
	for _, deletedEntry := range deletedEntries {
		if deletedEntry.Section != addedSection || deletedEntry.Line == addedLine {
			continue
		}
		if reference != "" && changelogReferenceFromLine(deletedEntry.Line) == reference ||
			inlineProseCode.ReplaceAllString(deletedEntry.Line, "*$1*") == addedLine {
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

// changelogCategories are the section headers used by UNRELEASED_CHANGELOG.md.
var changelogCategories = []string{"Added", "Changed", "Fixed", "Deprecated", "Removed", "Security"}

// normalizeCategory maps the category header an entry was found under to one
// of the UNRELEASED_CHANGELOG.md sections, defaulting to "Added".
func normalizeCategory(category string) string {
	for _, c := range changelogCategories {
		if strings.EqualFold(category, c) {
			return c
		}
	}
	return "Added"
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

// attemptFix removes the misplaced entries from the generated changelog and
// inserts them into the matching category sections of UNRELEASED_CHANGELOG.md.
func attemptFix(content string, issues []Issue, changelogPath, unreleasedPath string) (bool, error) {
	lines := strings.Split(content, "\n")

	// Remove issue lines from the changelog (highest line number first)
	linesToRemove := make([]int, 0, len(issues))
	for _, issue := range issues {
		linesToRemove = append(linesToRemove, issue.Line)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(linesToRemove)))
	for _, lineNum := range linesToRemove {
		lines = append(lines[:lineNum], lines[lineNum+1:]...)
	}

	// Insert the entries into UNRELEASED_CHANGELOG.md
	unreleasedContent, err := readFile(unreleasedPath)
	if err != nil {
		return false, fmt.Errorf("failed to read %s: %w", unreleasedPath, err)
	}
	unreleasedLines := strings.Split(unreleasedContent, "\n")

	issuesByCategory := make(map[string][]Issue)
	for _, issue := range issues {
		issuesByCategory[issue.Category] = append(issuesByCategory[issue.Category], issue)
	}

	// Iterate categories in a stable order
	for _, category := range changelogCategories {
		categoryIssues := issuesByCategory[category]
		if len(categoryIssues) == 0 {
			continue
		}
		insertPos := findCategoryInsertPos(unreleasedLines, category)
		if insertPos == -1 {
			return false, fmt.Errorf("could not find '## %s' section in %s", category, unreleasedPath)
		}
		for _, issue := range categoryIssues {
			unreleasedLines = append(unreleasedLines[:insertPos], append([]string{issue.Content}, unreleasedLines[insertPos:]...)...)
			insertPos++
		}
	}

	// Save the destination first so a failed write cannot discard an entry.
	if err := writeFile(unreleasedPath, strings.Join(unreleasedLines, "\n")); err != nil {
		return false, err
	}
	if err := writeFile(changelogPath, strings.Join(lines, "\n")); err != nil {
		if rollbackErr := writeFile(unreleasedPath, unreleasedContent); rollbackErr != nil {
			return false, fmt.Errorf("write changelog: %w; restore unreleased file: %v", err, rollbackErr)
		}
		return false, err
	}
	return true, nil
}

// findCategoryInsertPos returns the line index at which new entries should be
// inserted for the given category in UNRELEASED_CHANGELOG.md: after the
// section's placeholder comment and any existing entries, before the next
// section starts. Returns -1 if the section is missing.
func findCategoryInsertPos(lines []string, category string) int {
	sectionStart := -1
	for i, line := range lines {
		if strings.TrimSpace(line) == "## "+category {
			sectionStart = i
			break
		}
	}
	if sectionStart == -1 {
		return -1
	}

	insertPos := sectionStart + 1
	for i := sectionStart + 1; i < len(lines); i++ {
		trimmed := strings.TrimSpace(lines[i])
		if strings.HasPrefix(trimmed, "## ") || strings.HasPrefix(trimmed, "---") {
			break
		}
		if strings.HasPrefix(trimmed, "<!--") || strings.HasPrefix(trimmed, "- ") || strings.HasPrefix(trimmed, "* ") {
			insertPos = i + 1
		}
	}
	return insertPos
}

func readFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var content strings.Builder
	scanner := bufio.NewScanner(file)
	// M-Press changelog entries can contain generated links long enough to
	// exceed Scanner's default 64 KiB token limit.
	scanner.Buffer(make([]byte, 64*1024), 8*1024*1024)
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

// publishedReleaseNotes reads maintainer-published evidence from the canonical
// repository. Missing, draft, mismatched, or unavailable releases fail closed.
func publishedReleaseNotes(section string) (string, error) {
	if !regexp.MustCompile(`^v3\.[0-9]+\.[0-9]+(?:-[A-Za-z0-9.-]+)?$`).MatchString(section) {
		return "", fmt.Errorf("invalid release tag %q", section)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	output, err := exec.CommandContext(ctx, "gh", "api", "--hostname", "github.com",
		"repos/wailsapp/wails/releases/tags/"+section).Output()
	if err != nil {
		return "", fmt.Errorf("release lookup: %w", err)
	}
	var release struct {
		TagName string `json:"tag_name"`
		Body    string `json:"body"`
		Draft   bool   `json:"draft"`
	}
	if err := json.Unmarshal(output, &release); err != nil {
		return "", err
	}
	if release.Draft || release.TagName != section {
		return "", fmt.Errorf("release is draft or tag does not match")
	}
	return release.Body, nil
}

// isPublishedBackfill accepts only an exact, PR-linked bullet in the published
// notes for the same release; evidence from other releases cannot authorise it.
func isPublishedBackfill(line, section string, notes map[string]string) bool {
	line = strings.TrimSpace(line)
	if !strings.HasPrefix(line, "- ") || pullRequestReferenceFromLine(line) == "" {
		return false
	}
	for _, published := range strings.Split(notes[section], "\n") {
		if strings.TrimSpace(published) == line {
			return true
		}
	}
	return false
}
