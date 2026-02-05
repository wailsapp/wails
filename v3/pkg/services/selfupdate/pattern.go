package selfupdate

import (
	"fmt"
	"runtime"
	"strings"
)

// DefaultAssetPattern is the default pattern for matching release assets.
const DefaultAssetPattern = "{name}_{goos}_{goarch}{ext}"

// PatternVariables holds the values used for asset pattern substitution.
type PatternVariables struct {
	// Name is the application name.
	Name string

	// Version is the release version (without 'v' prefix).
	Version string

	// GOOS is the Go operating system (darwin, linux, windows).
	GOOS string

	// GOARCH is the Go architecture (amd64, arm64, 386).
	GOARCH string

	// Variant is an optional build variant (e.g., webkit2_41).
	Variant string
}

// ResolveAssetPattern expands a pattern template into a concrete asset name.
//
// Supported template variables:
//   - {name}     - Application name
//   - {version}  - Release version (e.g., "1.2.0")
//   - {goos}     - Go OS: darwin, linux, windows
//   - {goarch}   - Go arch: amd64, arm64, 386
//   - {platform} - Combined: {goos}-{goarch} (e.g., "darwin-arm64")
//   - {variant}  - Build variant (e.g., "webkit2_41"), empty if not set
//   - {ext}      - Platform extension: .exe for windows, empty otherwise
//
// Example patterns:
//   - "{name}_{goos}_{goarch}{ext}" → "myapp_linux_amd64"
//   - "{name}-v{version}-{platform}{ext}" → "myapp-v1.2.0-darwin-arm64"
//   - "{name}_{goos}_{goarch}_{variant}{ext}" → "myapp_linux_amd64_webkit2_41"
func ResolveAssetPattern(pattern string, vars PatternVariables) string {
	if pattern == "" {
		pattern = DefaultAssetPattern
	}

	// Default to current platform if not specified
	goos := vars.GOOS
	if goos == "" {
		goos = runtime.GOOS
	}

	goarch := vars.GOARCH
	if goarch == "" {
		goarch = runtime.GOARCH
	}

	// Build the extension based on OS
	ext := platformExtension(goos)

	// Platform is a convenience combination
	platform := fmt.Sprintf("%s-%s", goos, goarch)

	// Perform substitutions
	result := pattern
	result = strings.ReplaceAll(result, "{name}", vars.Name)
	result = strings.ReplaceAll(result, "{version}", vars.Version)
	result = strings.ReplaceAll(result, "{goos}", goos)
	result = strings.ReplaceAll(result, "{goarch}", goarch)
	result = strings.ReplaceAll(result, "{platform}", platform)
	result = strings.ReplaceAll(result, "{ext}", ext)

	// Handle variant - if empty, remove any trailing underscore or dash
	if vars.Variant != "" {
		result = strings.ReplaceAll(result, "{variant}", vars.Variant)
	} else {
		// Remove {variant} and any preceding separator
		result = strings.ReplaceAll(result, "_{variant}", "")
		result = strings.ReplaceAll(result, "-{variant}", "")
		result = strings.ReplaceAll(result, "{variant}", "")
	}

	return result
}

// platformExtension returns the appropriate file extension for the given OS.
func platformExtension(goos string) string {
	switch goos {
	case "windows":
		return ".exe"
	default:
		return ""
	}
}

// MatchAssetName checks if an asset name matches the expected pattern.
// This is useful for finding the correct asset in a release with multiple assets.
//
// Parameters:
//   - assetName: The actual asset filename from the release
//   - pattern: The pattern template to match against
//   - vars: The variables to use for pattern expansion
//
// Returns true if the asset name matches the expanded pattern.
func MatchAssetName(assetName, pattern string, vars PatternVariables) bool {
	expected := ResolveAssetPattern(pattern, vars)
	return assetName == expected
}

// MatchAssetNameWithExtensions checks if an asset name matches the pattern
// with common archive extensions appended.
//
// This handles cases where releases include compressed archives:
//   - .tar.gz, .tgz
//   - .zip
//   - .tar.xz
//   - .dmg (macOS)
//   - .msi (Windows)
//
// Returns the matching extension if found, empty string otherwise.
func MatchAssetNameWithExtensions(assetName, pattern string, vars PatternVariables) string {
	baseName := ResolveAssetPattern(pattern, vars)

	// Common archive extensions to try
	extensions := []string{
		"",        // No extension (binary)
		".tar.gz", // Common Linux/macOS
		".tgz",    // Short tar.gz
		".zip",    // Common all platforms
		".tar.xz", // XZ compressed
	}

	// Platform-specific extensions
	switch vars.GOOS {
	case "darwin":
		extensions = append(extensions, ".dmg", ".app.tar.gz", ".app.zip")
	case "windows":
		extensions = append(extensions, ".msi", ".exe") // .exe might be in zip
	}

	for _, ext := range extensions {
		if assetName == baseName+ext {
			return ext
		}
	}

	return ""
}

// FindMatchingAsset searches a list of asset names for one that matches
// the given pattern and variables.
//
// Returns the matching asset name and its extension, or empty strings if not found.
func FindMatchingAsset(assetNames []string, pattern string, vars PatternVariables) (assetName, extension string) {
	for _, name := range assetNames {
		if ext := MatchAssetNameWithExtensions(name, pattern, vars); ext != "" || MatchAssetName(name, pattern, vars) {
			return name, ext
		}
	}
	return "", ""
}

// ExtractPatternVariables attempts to extract variable values from an asset name.
// This is a best-effort heuristic extraction that detects OS and architecture
// from common naming patterns in the asset filename.
//
// Returns the extracted variables and true if at least one variable was detected,
// or zero values and false if nothing could be extracted.
func ExtractPatternVariables(assetName string) (PatternVariables, bool) {
	var vars PatternVariables
	detected := false

	lower := strings.ToLower(assetName)

	// Try to detect OS from name
	switch {
	case strings.Contains(lower, "darwin") || strings.Contains(lower, "macos"):
		vars.GOOS = "darwin"
		detected = true
	case strings.Contains(lower, "linux"):
		vars.GOOS = "linux"
		detected = true
	case strings.Contains(lower, "windows") || strings.HasSuffix(lower, ".exe"):
		vars.GOOS = "windows"
		detected = true
	}

	// Try to detect arch from name
	switch {
	case strings.Contains(lower, "arm64") || strings.Contains(lower, "aarch64"):
		vars.GOARCH = "arm64"
		detected = true
	case strings.Contains(lower, "amd64") || strings.Contains(lower, "x86_64") || strings.Contains(lower, "x64"):
		vars.GOARCH = "amd64"
		detected = true
	case strings.Contains(lower, "386") || strings.Contains(lower, "i386"):
		vars.GOARCH = "386"
		detected = true
	}

	return vars, detected
}
