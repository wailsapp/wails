package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// CleanupPattern represents a cleanup rule
type CleanupPattern struct {
	Type        string // "prefix", "suffix", "exact"
	Pattern     string // The pattern to match
	TargetFiles bool   // true = target files, false = target directories
	Description string // Description for logging
}

// Patterns to clean up during test cleanup
var cleanupPatterns = []CleanupPattern{
	// Test binaries from examples
	{Type: "prefix", Pattern: "testbuild-", TargetFiles: true, Description: "test binary"},
	
	// Go test binaries
	{Type: "suffix", Pattern: ".test", TargetFiles: true, Description: "Go test binary"},
	
	// Package artifacts from packaging tests (only in internal/commands directory)
	// Note: Only clean these from the commands directory, not from test temp directories
	
	// Test template directories from template tests
	{Type: "prefix", Pattern: "test-template-", TargetFiles: false, Description: "test template directory"},
	
	// CLI test binaries (files named exactly "appimage_testfiles")
	{Type: "exact", Pattern: "appimage_testfiles", TargetFiles: true, Description: "CLI test binary"},
}

func main() {
	fmt.Println("Starting cleanup...")
	cleanedCount := 0

	// Walk through all files and directories
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Continue on errors
		}

		// Skip if we're in the .git directory
		if strings.Contains(path, ".git") {
			return nil
		}

		name := info.Name()

		// Check each cleanup pattern
		for _, pattern := range cleanupPatterns {
			shouldClean := false

			switch pattern.Type {
			case "prefix":
				shouldClean = strings.HasPrefix(name, pattern.Pattern)
			case "suffix":
				shouldClean = strings.HasSuffix(name, pattern.Pattern)
			case "exact":
				shouldClean = name == pattern.Pattern
			}

			if shouldClean {
				// Check if the pattern targets the correct type (file or directory)
				if pattern.TargetFiles && info.Mode().IsRegular() {
					// This pattern targets files and this is a file
					fmt.Printf("Removing %s: %s\n", pattern.Description, path)
					os.Remove(path)
					cleanedCount++
					break // Don't check other patterns for this file
				} else if !pattern.TargetFiles && info.IsDir() {
					// This pattern targets directories and this is a directory
					fmt.Printf("Removing %s: %s\n", pattern.Description, path)
					os.RemoveAll(path)
					cleanedCount++
					return filepath.SkipDir // Don't recurse into removed directory
				}
				// If the pattern matches but the file type doesn't match TargetFiles, continue checking other patterns
			}
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error during cleanup: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Cleanup completed. Removed %d items.\n", cleanedCount)
}