#!/bin/bash

# Test script to simulate v3-check-changelog workflow locally
# This simulates exactly what the workflow would do for PR #4392

set -e

echo "🧪 Testing v3-check-changelog workflow locally..."
echo "================================================"

# Create temp directory for test
TEST_DIR="/tmp/wails-changelog-test"
rm -rf "$TEST_DIR"
mkdir -p "$TEST_DIR"

echo "📁 Test directory: $TEST_DIR"

# Copy current changelog to test location
cp "docs/src/content/docs/changelog.mdx" "$TEST_DIR/changelog.mdx"

# Simulate the PR diff for #4392 (the problematic entries)
cat > "$TEST_DIR/pr_added_lines.txt" << 'EOF'
- Add distribution-specific build dependencies for Linux by @leaanthony in [PR](https://github.com/wailsapp/wails/pull/4345)
- Added bindings guide by @atterpac in [PR](https://github.com/wailsapp/wails/pull/4404)
EOF

echo "📝 Simulated PR diff (lines that would be added):"
cat "$TEST_DIR/pr_added_lines.txt"
echo ""

# Create the validation script (same as in workflow)
cat > "$TEST_DIR/validate_and_fix.go" << 'EOF'
package main

import (
  "bufio"
  "fmt"
  "os"
  "path/filepath"
  "strings"
)

// Simplified validator for GitHub Actions
func main() {
  changelogPath := "changelog.mdx"
  addedLinesPath := "pr_added_lines.txt"
  
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
        Line:        lineNum,
        Content:     strings.TrimSpace(line),
        Section:     currentSection,
        Category:    getCurrentCategory(lines, lineNum),
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
  fixed, err := attemptFix(content, issues)
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
  // Check if this exact line (trimmed) was added in this PR
  trimmedLine := strings.TrimSpace(line)
  
  for _, addedLine := range addedLines {
    trimmedAdded := strings.TrimSpace(addedLine)
    if trimmedAdded == trimmedLine {
      return true
    }
    
    // Also check if the content matches (handles whitespace differences)
    if strings.Contains(trimmedAdded, trimmedLine) && len(trimmedAdded) > 0 {
      return true
    }
  }
  
  return false
}

func getCurrentCategory(lines []string, lineNum int) string {
  // Look backwards to find the current category
  for i := lineNum - 1; i >= 0; i-- {
    line := strings.TrimSpace(lines[i])
    if strings.HasPrefix(line, "### ") {
      return strings.TrimSpace(line[4:])
    }
    if strings.HasPrefix(line, "## ") && 
       !strings.Contains(line, "[Unreleased]") && 
       !strings.Contains(line, "v3.0.0-alpha") {
      // This is a malformed category like "## Added" - should be "### Added"
      // But we'll handle it for backward compatibility
      return strings.TrimSpace(line[3:])
    }
    if strings.HasPrefix(line, "## ") && 
       (strings.Contains(line, "[Unreleased]") || strings.Contains(line, "v3.0.0-alpha")) {
      // This is a version section header, stop looking
      break
    }
  }
  return "Added" // Default fallback for new entries
}

func attemptFix(content string, issues []Issue) (bool, error) {
  lines := strings.Split(content, "\n")
  
  // Find unreleased section
  unreleasedStart := -1
  unreleasedEnd := -1
  
  for i, line := range lines {
    if strings.Contains(line, "[Unreleased]") {
      unreleasedStart = i
      // Find where unreleased section ends (next ## section)
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
  
  // Remove issues from original locations (in reverse order to maintain line numbers)
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
  // Find where to insert (after existing categories or create new ones)
  for category, categoryIssues := range issuesByCategory {
    // Look for existing category in unreleased section
    categoryFound := false
    insertPos := unreleasedStart + 1
    
    for i := unreleasedStart + 1; i < unreleasedEnd && i < len(lines); i++ {
      if strings.Contains(lines[i], "### " + category) || strings.Contains(lines[i], "## " + category) {
        categoryFound = true
        // Find the end of this category to append entries
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
      // Add new category at the end of unreleased section
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
      
      // Update unreleasedEnd since we added lines
      unreleasedEnd += len(newLines)
    }
    
    // Add entries to the category
    for _, issue := range categoryIssues {
      lines = append(lines[:insertPos], append([]string{issue.Content}, lines[insertPos:]...)...)
      insertPos++
      unreleasedEnd++ // Update end position
    }
  }
  
  // Write back to file
  newContent := strings.Join(lines, "\n")
  return true, writeFile("changelog_fixed.mdx", newContent)
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
EOF

echo "🔄 Running validation script..."
cd "$TEST_DIR"

# Initialize go module for the test
go mod init test 2>/dev/null || true

# Run the validation
OUTPUT=$(go run validate_and_fix.go 2>&1)
echo "$OUTPUT"

# Check if a fixed file was created
if [ -f "changelog_fixed.mdx" ]; then
  echo ""
  echo "📄 Fixed changelog was created. Showing differences:"
  echo "=================================================="
  
  # Show the before/after diff
  echo "🔍 Showing changes made:"
  diff -u changelog.mdx changelog_fixed.mdx || true
  
  echo ""
  echo "📋 Summary: The workflow would automatically fix the changelog by moving misplaced entries to [Unreleased]"
else
  echo "📋 Summary: No fixes were needed or could not fix automatically"
fi

# Cleanup
cd - > /dev/null
rm -rf "$TEST_DIR"

echo ""
echo "✅ Local test completed! This simulates exactly what the GitHub workflow would do."