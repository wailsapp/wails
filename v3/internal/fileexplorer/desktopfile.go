package fileexplorer

import (
	"bufio"
	"io"
	"os"
	"strings"
)

// DesktopEntry represents a parsed .desktop file's [Desktop Entry] section.
// This is a minimal parser that only extracts the fields we need,
// replacing the full gopkg.in/ini.v1 dependency (~34KB + 68 transitive deps).
type DesktopEntry struct {
	Exec string
}

// ParseDesktopFile parses a .desktop file and returns the Desktop Entry section.
// It follows the Desktop Entry Specification:
// ParseDesktopFile parses the `[Desktop Entry]` section of the desktop file at path and returns a DesktopEntry.
// It returns an error if the file cannot be opened or if parsing the file fails.
func ParseDesktopFile(path string) (*DesktopEntry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ParseDesktopReader(f)
}

// ParseDesktopReader parses the [Desktop Entry] section of a .desktop file from r and extracts the Exec value.
// It ignores empty lines and lines starting with '#', treats section names as case-sensitive, and stops parsing after leaving the [Desktop Entry] section.
// The returned *DesktopEntry has Exec set to the exact value of the Exec key if present (whitespace preserved).
// An error is returned if reading from r fails.
func ParseDesktopReader(r io.Reader) (*DesktopEntry, error) {
	scanner := bufio.NewScanner(r)
	entry := &DesktopEntry{}

	inDesktopEntry := false

	for scanner.Scan() {
		line := scanner.Text()

		// Skip empty lines
		if len(line) == 0 {
			continue
		}

		// Skip comments (# at start of line)
		if line[0] == '#' {
			continue
		}

		// Handle section headers
		if line[0] == '[' {
			// Check if this is the [Desktop Entry] section
			// The spec says section names are case-sensitive
			trimmed := strings.TrimSpace(line)
			if trimmed == "[Desktop Entry]" {
				inDesktopEntry = true
			} else if inDesktopEntry {
				// We've left the [Desktop Entry] section
				// (e.g., entering [Desktop Action new-window])
				// We already have what we need, so we can stop
				break
			}
			continue
		}

		// Only process key=value pairs in [Desktop Entry] section
		if !inDesktopEntry {
			continue
		}

		// Parse key=value (spec says no spaces around =, but be lenient)
		eqIdx := strings.Index(line, "=")
		if eqIdx == -1 {
			continue
		}

		key := strings.TrimSpace(line[:eqIdx])
		value := line[eqIdx+1:] // Don't trim value - preserve intentional whitespace

		// We only need the Exec key
		// Per spec, keys are case-sensitive and Exec is always "Exec"
		if key == "Exec" {
			entry.Exec = value
			// Continue parsing in case there are multiple Exec lines (shouldn't happen but be safe)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return entry, nil
}