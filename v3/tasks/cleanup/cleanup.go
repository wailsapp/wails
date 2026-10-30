package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// CleanupPattern represents a cleanup rule
type CleanupPattern struct {
	Type         string // "prefix", "suffix", "exact"
	Pattern      string // The pattern to match
	TargetFiles  bool   // true = target files, false = target directories
	Description  string // Description for logging
	CommandsOnly bool   // Restrict package artifacts to the commands test directory
}

// Patterns to clean up during test cleanup
var cleanupPatterns = []CleanupPattern{
	// Test binaries from examples
	{Type: "prefix", Pattern: "testbuild-", TargetFiles: true, Description: "test binary"},

	// Go test binaries
	{Type: "suffix", Pattern: ".test", TargetFiles: true, Description: "Go test binary"},

	// Package artifacts from packaging tests (only in internal/commands directory)
	// Note: Only clean these from the commands directory, not from test temp directories
	{Type: "exact", Pattern: "myapp.ARCHLINUX", TargetFiles: true, CommandsOnly: true, Description: "Linux ARCHLINUX package"},
	{Type: "exact", Pattern: "myapp.DEB", TargetFiles: true, CommandsOnly: true, Description: "Linux DEB package"},
	{Type: "exact", Pattern: "myapp.RPM", TargetFiles: true, CommandsOnly: true, Description: "Linux RPM package"},

	// Test template directories from template tests
	{Type: "prefix", Pattern: "test-template-", TargetFiles: false, Description: "test template directory"},

	// CLI test binaries (files named exactly "appimage_testfiles")
	{Type: "exact", Pattern: "appimage_testfiles", TargetFiles: true, Description: "CLI test binary"},
}

func main() {
	fmt.Println("Starting cleanup...")
	count, err := cleanup(".", os.Stdout)
	if err != nil {
		fmt.Printf("Error during cleanup: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Cleanup completed. Removed %d items.\n", count)
}

func cleanup(root string, output io.Writer) (int, error) {
	count := 0
	err := filepath.WalkDir(root, func(path string, info os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && info.Name() == ".git" {
			return filepath.SkipDir
		}
		for _, pattern := range cleanupPatterns {
			matches, err := pattern.matches(root, path, info)
			if err != nil {
				return err
			}
			if !matches {
				continue
			}

			remove := os.Remove
			if info.IsDir() {
				remove = os.RemoveAll
			}
			if err := remove(path); err != nil {
				return fmt.Errorf("remove %s: %w", path, err)
			}
			fmt.Fprintf(output, "Removing %s: %s\n", pattern.Description, path)
			count++
			if info.IsDir() {
				return filepath.SkipDir
			}
			break
		}
		return nil
	})
	return count, err
}

func (pattern CleanupPattern) matches(root, path string, entry os.DirEntry) (bool, error) {
	if !pattern.matchesName(entry.Name()) {
		return false, nil
	}
	if pattern.TargetFiles {
		if !entry.Type().IsRegular() {
			return false, nil
		}
	} else if !entry.IsDir() {
		return false, nil
	}
	if !pattern.CommandsOnly {
		return true, nil
	}
	relative, err := filepath.Rel(root, path)
	if err != nil {
		return false, err
	}
	directory := filepath.Dir(relative)
	return directory == filepath.Join("internal", "commands") || directory == filepath.Join("v3", "internal", "commands"), nil
}

func (pattern CleanupPattern) matchesName(name string) bool {
	switch pattern.Type {
	case "prefix":
		return strings.HasPrefix(name, pattern.Pattern)
	case "suffix":
		return strings.HasSuffix(name, pattern.Pattern)
	case "exact":
		return name == pattern.Pattern
	}
	return false
}
