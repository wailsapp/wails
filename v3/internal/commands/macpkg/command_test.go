package macpkg

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestBuildMacPKG_PlatformCheck(t *testing.T) {
	if runtime.GOOS == "darwin" {
		t.Skip("Skipping platform check test on macOS")
	}
	
	options := &BuildMacPKGOptions{
		ConfigPath: "test-config.yaml",
	}
	
	err := BuildMacPKG(options)
	if err == nil {
		t.Errorf("Expected error on non-macOS platform, but got nil")
	}
	
	expectedError := "Mac PKG building is only supported on macOS"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestBuildMacPKG_GenerateTemplate(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("Skipping macOS-specific test on non-macOS platform")
	}
	
	tempDir := t.TempDir()
	templatePath := filepath.Join(tempDir, "generated-template.yaml")
	
	options := &BuildMacPKGOptions{
		GenerateTemplate: true,
		ConfigPath:       templatePath,
	}
	
	err := BuildMacPKG(options)
	if err != nil {
		t.Fatalf("Template generation failed: %v", err)
	}
	
	// Verify template was created
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		t.Errorf("Template file was not created at %s", templatePath)
	}
}

func TestBuildMacPKG_MissingConfig(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("Skipping macOS-specific test on non-macOS platform")
	}
	
	options := &BuildMacPKGOptions{
		ConfigPath: "",
	}
	
	err := BuildMacPKG(options)
	if err == nil {
		t.Errorf("Expected error for missing config, but got nil")
	}
	
	expectedError := "config file is required"
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("Expected error to contain '%s', got '%s'", expectedError, err.Error())
	}
}

func TestBuildMacPKG_ConfigNotFound(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("Skipping macOS-specific test on non-macOS platform")
	}
	
	options := &BuildMacPKGOptions{
		ConfigPath: "/nonexistent/config.yaml",
	}
	
	err := BuildMacPKG(options)
	if err == nil {
		t.Errorf("Expected error for nonexistent config, but got nil")
	}
	
	expectedError := "failed to load configuration"
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("Expected error to contain '%s', got '%s'", expectedError, err.Error())
	}
}

func TestBuildMacPKG_ValidateOnly(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("Skipping macOS-specific test on non-macOS platform")
	}
	
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "validate-test.yaml")
	
	configContent := `app_name: "TestApp"
app_path: "./TestApp.app"
bundle_id: "com.test.validate"
version: "1.0.0"
output_path: "./TestApp-installer.pkg"
`
	
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}
	
	options := &BuildMacPKGOptions{
		ConfigPath:   configPath,
		ValidateOnly: true,
	}
	
	// This should fail because the app path doesn't exist
	err := BuildMacPKG(options)
	if err == nil {
		t.Errorf("Expected validation error, but got nil")
	}
	
	expectedError := "app_path does not exist"
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("Expected error to contain '%s', got '%s'", expectedError, err.Error())
	}
}

func TestShouldNotarize(t *testing.T) {
	tests := []struct {
		name     string
		config   *PKGConfig
		expected bool
	}{
		{
			name: "all credentials provided",
			config: &PKGConfig{
				AppleID:     "test@example.com",
				AppPassword: "password",
				TeamID:      "TEAM123",
			},
			expected: true,
		},
		{
			name: "missing apple id",
			config: &PKGConfig{
				AppleID:     "",
				AppPassword: "password",
				TeamID:      "TEAM123",
			},
			expected: false,
		},
		{
			name: "missing app password",
			config: &PKGConfig{
				AppleID:     "test@example.com",
				AppPassword: "",
				TeamID:      "TEAM123",
			},
			expected: false,
		},
		{
			name: "missing team id",
			config: &PKGConfig{
				AppleID:     "test@example.com",
				AppPassword: "password",
				TeamID:      "",
			},
			expected: false,
		},
		{
			name: "no credentials",
			config: &PKGConfig{
				AppleID:     "",
				AppPassword: "",
				TeamID:      "",
			},
			expected: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := shouldNotarize(tt.config)
			if result != tt.expected {
				t.Errorf("Expected shouldNotarize() = %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestHandleTemplateGeneration(t *testing.T) {
	tempDir := t.TempDir()
	templatePath := filepath.Join(tempDir, "handle-template-test.yaml")
	
	// Test successful generation
	err := handleTemplateGeneration(templatePath)
	if err != nil {
		t.Fatalf("Template generation failed: %v", err)
	}
	
	// Verify file was created
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		t.Errorf("Template file was not created")
	}
	
	// Test file already exists
	err = handleTemplateGeneration(templatePath)
	if err == nil {
		t.Errorf("Expected error for existing file, but got nil")
	}
	
	expectedError := "file already exists"
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("Expected error to contain '%s', got '%s'", expectedError, err.Error())
	}
}

func TestHandleTemplateGeneration_DefaultPath(t *testing.T) {
	// Change to temp directory to avoid creating files in project root
	originalDir, _ := os.Getwd()
	tempDir := t.TempDir()
	os.Chdir(tempDir)
	defer os.Chdir(originalDir)
	
	// Test with empty config path (should use default)
	err := handleTemplateGeneration("")
	if err != nil {
		t.Fatalf("Template generation with default path failed: %v", err)
	}
	
	// Verify default file was created
	defaultPath := "wails-pkg.yaml"
	if _, err := os.Stat(defaultPath); os.IsNotExist(err) {
		t.Errorf("Default template file was not created")
	}
}