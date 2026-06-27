package macpkg

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// LoadPKGConfig loads PKG configuration from a YAML file
func LoadPKGConfig(configPath string) (*PKGConfig, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}
	
	var config PKGConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}
	
	// Expand environment variables and resolve relative paths
	if err := config.expandAndResolve(filepath.Dir(configPath)); err != nil {
		return nil, fmt.Errorf("failed to process config: %w", err)
	}
	
	return &config, nil
}

// expandAndResolve expands environment variables and resolves relative paths
func (c *PKGConfig) expandAndResolve(configDir string) error {
	// Expand environment variables
	c.AppName = os.ExpandEnv(c.AppName)
	c.AppPath = os.ExpandEnv(c.AppPath)
	c.BundleID = os.ExpandEnv(c.BundleID)
	c.Version = os.ExpandEnv(c.Version)
	c.SigningIdentity = os.ExpandEnv(c.SigningIdentity)
	c.InstallLocation = os.ExpandEnv(c.InstallLocation)
	c.Title = os.ExpandEnv(c.Title)
	c.Background = os.ExpandEnv(c.Background)
	c.WelcomeFile = os.ExpandEnv(c.WelcomeFile)
	c.ReadmeFile = os.ExpandEnv(c.ReadmeFile)
	c.LicenseFile = os.ExpandEnv(c.LicenseFile)
	c.AppleID = os.ExpandEnv(c.AppleID)
	c.AppPassword = os.ExpandEnv(c.AppPassword)
	c.TeamID = os.ExpandEnv(c.TeamID)
	c.OutputPath = os.ExpandEnv(c.OutputPath)
	
	// Resolve relative paths
	if c.AppPath != "" && !filepath.IsAbs(c.AppPath) {
		c.AppPath = filepath.Join(configDir, c.AppPath)
	}
	
	if c.Background != "" && !filepath.IsAbs(c.Background) {
		c.Background = filepath.Join(configDir, c.Background)
	}
	
	if c.WelcomeFile != "" && !filepath.IsAbs(c.WelcomeFile) {
		c.WelcomeFile = filepath.Join(configDir, c.WelcomeFile)
	}
	
	if c.ReadmeFile != "" && !filepath.IsAbs(c.ReadmeFile) {
		c.ReadmeFile = filepath.Join(configDir, c.ReadmeFile)
	}
	
	if c.LicenseFile != "" && !filepath.IsAbs(c.LicenseFile) {
		c.LicenseFile = filepath.Join(configDir, c.LicenseFile)
	}
	
	if c.OutputPath != "" && !filepath.IsAbs(c.OutputPath) {
		c.OutputPath = filepath.Join(configDir, c.OutputPath)
	}
	
	return nil
}

// GenerateTemplate creates a sample PKG configuration file
func GenerateTemplate(outputPath string) error {
	template := `# Wails v3 Mac PKG Configuration
# This file configures the creation of macOS installer packages (.pkg)

# Application information
app_name: "MyApp"                           # Name of your application
app_path: "./build/MyApp.app"              # Path to the .app bundle (relative to this config file)
bundle_id: "com.mycompany.myapp"           # Unique bundle identifier
version: "1.0.0"                           # Version number

# Code signing (required for notarization)
signing_identity: "Developer ID Installer: Your Name (TEAM_ID)"  # Signing certificate identity

# Installation
install_location: "/Applications"          # Where the app will be installed (default: /Applications)

# Distribution package appearance
title: "MyApp Installer"                   # Title shown in installer
background: ""                             # Optional background image (PNG)
welcome_file: ""                           # Optional welcome RTF file
readme_file: ""                            # Optional readme RTF file  
license_file: ""                           # Optional license RTF file

# Notarization (optional but recommended for distribution)
apple_id: "${APPLE_ID}"                    # Apple ID for notarization
app_password: "${APP_PASSWORD}"            # App-specific password
team_id: "${TEAM_ID}"                      # Developer team ID

# Output
output_path: "./dist/MyApp-installer.pkg"  # Where to save the final installer

# Environment variables can be used with ${VAR_NAME} syntax
# Relative paths are resolved relative to this config file's location
`
	
	if err := os.WriteFile(outputPath, []byte(template), 0644); err != nil {
		return fmt.Errorf("failed to write template: %w", err)
	}
	
	return nil
}

// ValidateConfig performs comprehensive validation of the PKG configuration
func ValidateConfig(config *PKGConfig) error {
	var errors []string
	
	// Required fields
	if config.AppName == "" {
		errors = append(errors, "app_name is required")
	}
	
	if config.AppPath == "" {
		errors = append(errors, "app_path is required")
	} else {
		// Verify app bundle exists
		if _, err := os.Stat(config.AppPath); os.IsNotExist(err) {
			errors = append(errors, fmt.Sprintf("app_path does not exist: %s", config.AppPath))
		} else if !strings.HasSuffix(config.AppPath, ".app") {
			errors = append(errors, "app_path must point to a .app bundle")
		}
	}
	
	if config.BundleID == "" {
		errors = append(errors, "bundle_id is required")
	} else if !isValidBundleID(config.BundleID) {
		errors = append(errors, "bundle_id must be in reverse DNS format (e.g., com.company.app)")
	}
	
	if config.Version == "" {
		errors = append(errors, "version is required")
	}
	
	if config.OutputPath == "" {
		errors = append(errors, "output_path is required")
	} else if !strings.HasSuffix(config.OutputPath, ".pkg") {
		errors = append(errors, "output_path must end with .pkg")
	}
	
	// Validate optional files exist if specified
	optionalFiles := map[string]string{
		"background":    config.Background,
		"welcome_file":  config.WelcomeFile,
		"readme_file":   config.ReadmeFile,
		"license_file":  config.LicenseFile,
	}
	
	for name, path := range optionalFiles {
		if path != "" {
			if _, err := os.Stat(path); os.IsNotExist(err) {
				errors = append(errors, fmt.Sprintf("%s does not exist: %s", name, path))
			}
		}
	}
	
	// Notarization validation
	notarizationFields := []string{config.AppleID, config.AppPassword, config.TeamID}
	notarizationProvided := 0
	for _, field := range notarizationFields {
		if field != "" {
			notarizationProvided++
		}
	}
	
	if notarizationProvided > 0 && notarizationProvided < 3 {
		errors = append(errors, "if notarization is enabled, apple_id, app_password, and team_id are all required")
	}
	
	if len(errors) > 0 {
		return fmt.Errorf("configuration validation failed:\n- %s", strings.Join(errors, "\n- "))
	}
	
	return nil
}

// isValidBundleID checks if a bundle ID follows reverse DNS format
func isValidBundleID(bundleID string) bool {
	parts := strings.Split(bundleID, ".")
	if len(parts) < 2 {
		return false
	}
	
	for _, part := range parts {
		if part == "" {
			return false
		}
		// Check for valid characters (alphanumeric and hyphens)
		for _, r := range part {
			if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || 
				 (r >= '0' && r <= '9') || r == '-') {
				return false
			}
		}
	}
	
	return true
}