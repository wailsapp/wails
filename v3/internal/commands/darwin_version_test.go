package commands

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDarwinMinimumVersionInTemplates(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "wails-darwin-version-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	options := &BuildAssetsOptions{
		Dir:                filepath.Join(tempDir, "build"),
		Name:               "VersionTest",
		ProductName:        "Version Test App",
		ProductDescription: "Test",
		ProductVersion:     "1.0.0",
		ProductCompany:     "Test Co",
		ProductIdentifier:  "com.test.version",
		Silent:             true,
	}

	err = GenerateBuildAssets(options)
	if err != nil {
		t.Fatalf("GenerateBuildAssets() error = %v", err)
	}

	plistPath := filepath.Join(tempDir, "build", "darwin", "Info.plist")
	content, err := os.ReadFile(plistPath)
	if err != nil {
		t.Fatalf("Failed to read Info.plist: %v", err)
	}

	if !strings.Contains(string(content), "12.0.0") {
		t.Errorf("Expected LSMinimumSystemVersion to be 12.0.0, got:\n%s", string(content))
	}

	if strings.Contains(string(content), "10.15") {
		t.Errorf("Found outdated version 10.15 in Info.plist, should be 12.0.0")
	}
}

func TestDarwinDevPlistMinimumVersion(t *testing.T) {
	templateDir := "updatable_build_assets/darwin"
	devPlistPath := filepath.Join(templateDir, "Info.dev.plist.tmpl")

	content, err := os.ReadFile(devPlistPath)
	if err != nil {
		t.Skipf("Template file not found: %v", err)
	}

	if !strings.Contains(string(content), "12.0.0") {
		t.Errorf("Expected LSMinimumSystemVersion 12.0.0 in Info.dev.plist.tmpl, got:\n%s", string(content))
	}
}

func TestDarwinTaskfileMinimumVersion(t *testing.T) {
	taskfilePath := filepath.Join("build_assets", "darwin", "Taskfile.yml")

	content, err := os.ReadFile(taskfilePath)
	if err != nil {
		t.Skipf("Taskfile not found: %v", err)
	}

	str := string(content)
	if !strings.Contains(str, "mmacosx-version-min=12.0") {
		t.Errorf("Expected mmacosx-version-min=12.0 in Taskfile.yml")
	}
	if !strings.Contains(str, "MACOSX_DEPLOYMENT_TARGET: \"12.0\"") {
		t.Errorf("Expected MACOSX_DEPLOYMENT_TARGET: \"12.0\" in Taskfile.yml")
	}
	if strings.Contains(str, "10.15") {
		t.Errorf("Found outdated version 10.15 in Taskfile.yml, should be 12.0")
	}
}
