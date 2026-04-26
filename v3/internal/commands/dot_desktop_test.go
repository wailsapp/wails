package commands

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateDotDesktopPreservesName(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "wails-desktop-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name         string
		inputName    string
		expectedName string
		expectedFile string
	}{
		{
			name:         "simple name",
			inputName:    "MyApp",
			expectedName: "MyApp",
			expectedFile: "myapp.desktop",
		},
		{
			name:         "name with spaces",
			inputName:    "My App",
			expectedName: "My App",
			expectedFile: "my-app.desktop",
		},
		{
			name:         "name with mixed case",
			inputName:    "Wails",
			expectedName: "Wails",
			expectedFile: "wails.desktop",
		},
		{
			name:         "name with spaces and capitals",
			inputName:    "My Cool App",
			expectedName: "My Cool App",
			expectedFile: "my-cool-app.desktop",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outputFile := filepath.Join(tempDir, tt.expectedFile)

			options := &DotDesktopOptions{
				Name:       tt.inputName,
				Exec:       "my-app",
				OutputFile: outputFile,
			}

			err := GenerateDotDesktop(options)
			if err != nil {
				t.Fatalf("GenerateDotDesktop() error = %v", err)
			}

			data, err := os.ReadFile(outputFile)
			if err != nil {
				t.Fatalf("Failed to read output file: %v", err)
			}

			content := string(data)

			if !strings.Contains(content, "Name="+tt.expectedName) {
				t.Errorf("Expected Name=%s in desktop entry, got:\n%s", tt.expectedName, content)
			}
		})
	}
}

func TestGenerateDotDesktopFilenameNormalized(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "wails-desktop-filename-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	options := &DotDesktopOptions{
		Name: "My App",
		Exec: "my-app",
	}

	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(oldDir)

	os.Chdir(tempDir)

	err = GenerateDotDesktop(options)
	if err != nil {
		t.Fatalf("GenerateDotDesktop() error = %v", err)
	}

	expectedFile := "my-app.desktop"
	if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
		files, _ := os.ReadDir(".")
		var names []string
		for _, f := range files {
			names = append(names, f.Name())
		}
		t.Errorf("Expected file %s to exist, got files: %v", expectedFile, names)
	}
}
