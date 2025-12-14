package macpkg

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// createValidTestConfig creates a test config with a real temporary app path
func createValidTestConfig(t *testing.T) *PKGConfig {
	tempDir := t.TempDir()
	appPath := filepath.Join(tempDir, "TestApp.app")
	
	// Create a mock .app directory
	if err := os.MkdirAll(appPath, 0755); err != nil {
		t.Fatalf("Failed to create mock app directory: %v", err)
	}
	
	outputPath := filepath.Join(tempDir, "output.pkg")
	
	return &PKGConfig{
		AppName:    "TestApp",
		AppPath:    appPath,
		BundleID:   "com.test.app",
		Version:    "1.0.0",
		OutputPath: outputPath,
	}
}

func TestLoadPKGConfig(t *testing.T) {
	// Create a temporary config file
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "test-pkg.yaml")
	
	configContent := `app_name: "TestApp"
app_path: "./TestApp.app"
bundle_id: "com.test.app"
version: "1.0.0"
signing_identity: "Developer ID"
title: "Test App Installer"
output_path: "./TestApp-installer.pkg"
`
	
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}
	
	// Test loading config
	config, err := LoadPKGConfig(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	
	// Verify loaded values
	if config.AppName != "TestApp" {
		t.Errorf("Expected AppName 'TestApp', got '%s'", config.AppName)
	}
	
	if config.BundleID != "com.test.app" {
		t.Errorf("Expected BundleID 'com.test.app', got '%s'", config.BundleID)
	}
	
	if config.Version != "1.0.0" {
		t.Errorf("Expected Version '1.0.0', got '%s'", config.Version)
	}
	
	// Test path resolution (should be absolute)
	expectedAppPath := filepath.Join(tempDir, "TestApp.app")
	if config.AppPath != expectedAppPath {
		t.Errorf("Expected AppPath '%s', got '%s'", expectedAppPath, config.AppPath)
	}
}

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name          string
		config        *PKGConfig
		shouldFail    bool
		expectedError string
	}{
		{
			name:       "valid minimal config",
			config:     createValidTestConfig(t),
			shouldFail: false,
		},
		{
			name: "missing app name",
			config: &PKGConfig{
				AppPath:    "/path/to/TestApp.app",
				BundleID:   "com.test.app",
				Version:    "1.0.0",
				OutputPath: "/path/to/output.pkg",
			},
			shouldFail:    true,
			expectedError: "app_name is required",
		},
		{
			name: "invalid bundle ID",
			config: &PKGConfig{
				AppName:    "TestApp",
				AppPath:    "/path/to/TestApp.app",
				BundleID:   "invalid-bundle-id",
				Version:    "1.0.0",
				OutputPath: "/path/to/output.pkg",
			},
			shouldFail:    true,
			expectedError: "bundle_id must be in reverse DNS format",
		},
		{
			name: "missing output extension",
			config: &PKGConfig{
				AppName:    "TestApp",
				AppPath:    "/path/to/TestApp.app",
				BundleID:   "com.test.app",
				Version:    "1.0.0",
				OutputPath: "/path/to/output",
			},
			shouldFail:    true,
			expectedError: "output_path must end with .pkg",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateConfig(tt.config)
			
			if tt.shouldFail {
				if err == nil {
					t.Errorf("Expected validation to fail, but it passed")
				} else if tt.expectedError != "" && !strings.Contains(err.Error(), tt.expectedError) {
					t.Errorf("Expected error to contain '%s', got '%s'", tt.expectedError, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected validation to pass, but got error: %v", err)
				}
			}
		})
	}
}

func TestIsValidBundleID(t *testing.T) {
	tests := []struct {
		bundleID string
		expected bool
	}{
		{"com.example.app", true},
		{"com.company.product.subproduct", true},
		{"org.opensource.tool", true},
		{"io.github.user.repo", true},
		{"com.example.my_app", true},        // underscores are valid
		{"com.my_company.app", true},        // underscores are valid
		{"com.example.app-name", true},      // hyphens are valid
		{"invalid", false},
		{"", false},
		{"com.", false},
		{".com.example", false},
		{"com..example", false},
		{"com.example.app!", false},
		{"com.example.app with spaces", false},
	}
	
	for _, tt := range tests {
		t.Run(tt.bundleID, func(t *testing.T) {
			result := isValidBundleID(tt.bundleID)
			if result != tt.expected {
				t.Errorf("Expected isValidBundleID(%s) = %v, got %v", tt.bundleID, tt.expected, result)
			}
		})
	}
}

func TestGenerateTemplate(t *testing.T) {
	tempDir := t.TempDir()
	templatePath := filepath.Join(tempDir, "test-template.yaml")
	
	if err := GenerateTemplate(templatePath); err != nil {
		t.Fatalf("Failed to generate template: %v", err)
	}
	
	// Verify file was created
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		t.Errorf("Template file was not created")
	}
	
	// Verify file has content
	content, err := os.ReadFile(templatePath)
	if err != nil {
		t.Fatalf("Failed to read template file: %v", err)
	}
	
	if len(content) == 0 {
		t.Errorf("Template file is empty")
	}
	
	// Verify it contains expected keys
	contentStr := string(content)
	expectedKeys := []string{"app_name", "bundle_id", "version", "output_path"}
	for _, key := range expectedKeys {
		if !strings.Contains(contentStr, key) {
			t.Errorf("Template missing expected key: %s", key)
		}
	}
}

func TestEnvironmentVariableExpansion(t *testing.T) {
	// Set test environment variable
	os.Setenv("TEST_APP_NAME", "EnvApp")
	defer os.Unsetenv("TEST_APP_NAME")
	
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "test-env.yaml")
	
	configContent := `app_name: "${TEST_APP_NAME}"
app_path: "./TestApp.app"
bundle_id: "com.test.env"
version: "1.0.0"
output_path: "./output.pkg"
`
	
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}
	
	config, err := LoadPKGConfig(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	
	if config.AppName != "EnvApp" {
		t.Errorf("Environment variable not expanded: expected 'EnvApp', got '%s'", config.AppName)
	}
}