package commands

import (
	"fmt"
	"strings"

	"github.com/wailsapp/wails/v3/internal/github"
)

type ToolVersionOptions struct {
	Version    string `name:"v" description:"Current version to bump"`
	Major      bool   `name:"major" description:"Bump major version"`
	Minor      bool   `name:"minor" description:"Bump minor version"`
	Patch      bool   `name:"patch" description:"Bump patch version"`
	Prerelease bool   `name:"prerelease" description:"Bump prerelease version (e.g., alpha.5 to alpha.6)"`
}

// bumpPrerelease increments the final numeric prerelease identifier.
func bumpPrerelease(prerelease string) string {
	start := strings.LastIndexByte(prerelease, '.') + 1
	lastPart := prerelease[start:]
	if lastPart == "" || strings.Trim(lastPart, "0123456789") != "" {
		return prerelease
	}

	// Increment decimal digits directly, without imposing a machine-integer limit.
	// Strip leading zeroes to preserve the previous integer formatter's behavior.
	digits := []byte(strings.TrimLeft(lastPart, "0"))
	for i := len(digits) - 1; i >= 0; i-- {
		if digits[i] != '9' {
			digits[i]++
			return prerelease[:start] + string(digits)
		}
		digits[i] = '0'
	}
	return prerelease[:start] + "1" + string(digits)
}

// ToolVersion bumps a semantic version based on the provided flags
func ToolVersion(options *ToolVersionOptions) error {
	DisableFooter = true

	if options.Version == "" {
		return fmt.Errorf("please provide a version using the -v flag")
	}

	// Check if the version has a "v" prefix
	hasVPrefix := false
	versionStr := options.Version
	if len(versionStr) > 0 && versionStr[0] == 'v' {
		hasVPrefix = true
		versionStr = versionStr[1:]
	}

	// Parse the current version
	semver, err := github.NewSemanticVersion(versionStr)
	if err != nil {
		return fmt.Errorf("invalid version format: %v", err)
	}

	// Get the current version components
	major := semver.Version.Major
	minor := semver.Version.Minor
	patch := semver.Version.Patch
	prerelease := semver.Version.Prerelease
	metadata := semver.Version.Metadata

	// Check if at least one flag is specified
	if !options.Major && !options.Minor && !options.Patch && !options.Prerelease {
		return fmt.Errorf("please specify one of -major, -minor, -patch, or -prerelease")
	}

	// Bump the version based on the flags (major takes precedence over minor, which takes precedence over patch)
	if options.Major {
		major++
		minor = 0
		patch = 0
	} else if options.Minor {
		minor++
		patch = 0
	} else if options.Patch {
		patch++
	} else if options.Prerelease {
		// If only prerelease flag is specified, bump the prerelease version
		if prerelease == "" {
			return fmt.Errorf("cannot bump prerelease version: no prerelease part in the version")
		}
		prerelease = bumpPrerelease(prerelease)
	}

	// Format the new version
	newVersion := fmt.Sprintf("%d.%d.%d", major, minor, patch)
	if prerelease != "" {
		newVersion += "-" + prerelease
	}
	if metadata != "" {
		newVersion += "+" + metadata
	}

	// Add the "v" prefix back if it was present in the input
	if hasVPrefix {
		newVersion = "v" + newVersion
	}

	// Print the new version
	fmt.Println(newVersion)

	return nil
}
