package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/pterm/pterm"
	"github.com/wailsapp/wails/v3/internal/flags"
)

// Entitlement represents a macOS entitlement
type Entitlement struct {
	Key         string
	Name        string
	Description string
	Category    string
}

// Common macOS entitlements organized by category
var availableEntitlements = []Entitlement{
	// Hardened Runtime - Code Execution
	{Key: "com.apple.security.cs.allow-jit", Name: "Allow JIT", Description: "Allow creating writable/executable memory using MAP_JIT (most secure option for JIT)", Category: "Code Execution"},
	{Key: "com.apple.security.cs.allow-unsigned-executable-memory", Name: "Allow Unsigned Executable Memory", Description: "Allow writable/executable memory without MAP_JIT restrictions", Category: "Code Execution"},
	{Key: "com.apple.security.cs.disable-executable-page-protection", Name: "Disable Executable Page Protection", Description: "Disable all executable memory protections (least secure)", Category: "Code Execution"},
	{Key: "com.apple.security.cs.disable-library-validation", Name: "Disable Library Validation", Description: "Allow loading unsigned or differently-signed libraries/frameworks", Category: "Code Execution"},
	{Key: "com.apple.security.cs.allow-dyld-environment-variables", Name: "Allow DYLD Environment Variables", Description: "Allow DYLD_* environment variables to modify library loading", Category: "Code Execution"},

	// Hardened Runtime - Resource Access
	{Key: "com.apple.security.device.audio-input", Name: "Audio Input (Microphone)", Description: "Access to audio input devices", Category: "Resource Access"},
	{Key: "com.apple.security.device.camera", Name: "Camera", Description: "Access to the camera", Category: "Resource Access"},
	{Key: "com.apple.security.personal-information.location", Name: "Location", Description: "Access to location services", Category: "Resource Access"},
	{Key: "com.apple.security.personal-information.addressbook", Name: "Address Book", Description: "Access to contacts", Category: "Resource Access"},
	{Key: "com.apple.security.personal-information.calendars", Name: "Calendars", Description: "Access to calendar data", Category: "Resource Access"},
	{Key: "com.apple.security.personal-information.photos-library", Name: "Photos Library", Description: "Access to the Photos library", Category: "Resource Access"},
	{Key: "com.apple.security.automation.apple-events", Name: "Apple Events", Description: "Send Apple Events to other apps (AppleScript)", Category: "Resource Access"},

	// App Sandbox - Basic
	{Key: "com.apple.security.app-sandbox", Name: "Enable App Sandbox", Description: "Enable the App Sandbox (required for Mac App Store)", Category: "App Sandbox"},

	// App Sandbox - Network
	{Key: "com.apple.security.network.client", Name: "Outgoing Network Connections", Description: "Allow outgoing network connections (client)", Category: "Network"},
	{Key: "com.apple.security.network.server", Name: "Incoming Network Connections", Description: "Allow incoming network connections (server)", Category: "Network"},

	// App Sandbox - Files
	{Key: "com.apple.security.files.user-selected.read-only", Name: "User-Selected Files (Read)", Description: "Read access to files the user selects", Category: "File Access"},
	{Key: "com.apple.security.files.user-selected.read-write", Name: "User-Selected Files (Read/Write)", Description: "Read/write access to files the user selects", Category: "File Access"},
	{Key: "com.apple.security.files.downloads.read-only", Name: "Downloads Folder (Read)", Description: "Read access to the Downloads folder", Category: "File Access"},
	{Key: "com.apple.security.files.downloads.read-write", Name: "Downloads Folder (Read/Write)", Description: "Read/write access to the Downloads folder", Category: "File Access"},

	// Development
	{Key: "com.apple.security.get-task-allow", Name: "Debugging", Description: "Allow debugging (disable for production)", Category: "Development"},
}

// EntitlementsSetup runs the interactive entitlements configuration wizard
func EntitlementsSetup(options *flags.EntitlementsSetup) error {
	pterm.DefaultHeader.Println("macOS Entitlements Setup")
	fmt.Println()

	// Build all options for custom selection
	var allOptions []huh.Option[string]
	for _, e := range availableEntitlements {
		label := fmt.Sprintf("[%s] %s", e.Category, e.Name)
		allOptions = append(allOptions, huh.NewOption(label, e.Key))
	}

	// Show quick presets first
	var preset string
	presetForm := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Which entitlements profile?").
				Description("Development and production use separate files").
				Options(
					huh.NewOption("Development (entitlements.dev.plist)", "dev"),
					huh.NewOption("Production (entitlements.plist)", "prod"),
					huh.NewOption("Both (recommended)", "both"),
					huh.NewOption("App Store (entitlements.plist with sandbox)", "appstore"),
					huh.NewOption("Custom", "custom"),
				).
				Value(&preset),
		),
	)

	if err := presetForm.Run(); err != nil {
		return err
	}

	devEntitlements := []string{
		"com.apple.security.cs.allow-jit",
		"com.apple.security.cs.allow-unsigned-executable-memory",
		"com.apple.security.cs.disable-library-validation",
		"com.apple.security.get-task-allow",
		"com.apple.security.network.client",
	}

	prodEntitlements := []string{
		"com.apple.security.network.client",
	}

	appStoreEntitlements := []string{
		"com.apple.security.app-sandbox",
		"com.apple.security.network.client",
		"com.apple.security.files.user-selected.read-write",
	}

	baseDir := "build/darwin"
	if options.Output != "" {
		baseDir = filepath.Dir(options.Output)
	}

	switch preset {
	case "dev":
		return writeEntitlementsFile(filepath.Join(baseDir, "entitlements.dev.plist"), devEntitlements)

	case "prod":
		return writeEntitlementsFile(filepath.Join(baseDir, "entitlements.plist"), prodEntitlements)

	case "both":
		if err := writeEntitlementsFile(filepath.Join(baseDir, "entitlements.dev.plist"), devEntitlements); err != nil {
			return err
		}
		return writeEntitlementsFile(filepath.Join(baseDir, "entitlements.plist"), prodEntitlements)

	case "appstore":
		return writeEntitlementsFile(filepath.Join(baseDir, "entitlements.plist"), appStoreEntitlements)

	case "custom":
		// Let user choose which file and entitlements
		var targetFile string
		var selected []string

		customForm := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("Target file").
					Options(
						huh.NewOption("entitlements.plist (production)", "entitlements.plist"),
						huh.NewOption("entitlements.dev.plist (development)", "entitlements.dev.plist"),
					).
					Value(&targetFile),
			),
			huh.NewGroup(
				huh.NewMultiSelect[string]().
					Title("Select entitlements").
					Description("Use space to select, enter to confirm").
					Options(allOptions...).
					Value(&selected),
			),
		)

		if err := customForm.Run(); err != nil {
			return err
		}

		return writeEntitlementsFile(filepath.Join(baseDir, targetFile), selected)
	}

	return nil
}

func writeEntitlementsFile(path string, entitlements []string) error {
	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Generate and write the plist
	plist := generateEntitlementsPlist(entitlements)
	if err := os.WriteFile(path, []byte(plist), 0644); err != nil {
		return fmt.Errorf("failed to write entitlements file: %w", err)
	}

	pterm.Success.Printfln("Wrote %s", path)

	// Show summary
	pterm.Info.Println("Entitlements:")
	for _, key := range entitlements {
		for _, e := range availableEntitlements {
			if e.Key == key {
				fmt.Printf("  - %s\n", e.Name)
				break
			}
		}
	}
	fmt.Println()

	return nil
}

func parseExistingEntitlements(path string) (map[string]bool, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	result := make(map[string]bool)
	lines := strings.Split(string(content), "\n")

	for i, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "<key>") && strings.HasSuffix(line, "</key>") {
			key := strings.TrimPrefix(line, "<key>")
			key = strings.TrimSuffix(key, "</key>")

			// Check if next line is <true/>
			if i+1 < len(lines) {
				nextLine := strings.TrimSpace(lines[i+1])
				if nextLine == "<true/>" {
					result[key] = true
				}
			}
		}
	}

	return result, nil
}

func generateEntitlementsPlist(entitlements []string) string {
	var sb strings.Builder

	sb.WriteString(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
`)

	for _, key := range entitlements {
		sb.WriteString(fmt.Sprintf("\t<key>%s</key>\n", key))
		sb.WriteString("\t<true/>\n")
	}

	sb.WriteString(`</dict>
</plist>
`)

	return sb.String()
}
