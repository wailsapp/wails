package commands

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateDotDesktopPreservesAppName(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "wails-dotdesktop-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name         string
		inputName    string
		exec         string
		expectedName string
		expectedFile string
	}{
		{
			name:         "name with spaces preserved",
			inputName:    "My App",
			exec:         "my-app",
			expectedName: "My App",
			expectedFile: "my-app.desktop",
		},
		{
			name:         "name with mixed case preserved",
			inputName:    "MyAwesome App",
			exec:         "myawesome-app",
			expectedName: "MyAwesome App",
			expectedFile: "myawesome-app.desktop",
		},
		{
			name:         "simple name unchanged",
			inputName:    "MyApp",
			exec:         "myapp",
			expectedName: "MyApp",
			expectedFile: "myapp.desktop",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outputFile := filepath.Join(tempDir, tt.expectedFile)
			options := &DotDesktopOptions{
				Name:       tt.inputName,
				Exec:       tt.exec,
				OutputFile: outputFile,
			}

			err := GenerateDotDesktop(options)
			if err != nil {
				t.Fatalf("GenerateDotDesktop() error = %v", err)
			}

			content, err := os.ReadFile(outputFile)
			if err != nil {
				t.Fatalf("Failed to read output file: %v", err)
			}

			if !strings.Contains(string(content), "Name="+tt.expectedName) {
				t.Errorf("Expected Name=%s in .desktop file, got:\n%s", tt.expectedName, string(content))
			}

			if !strings.Contains(string(content), "Exec="+tt.exec) {
				t.Errorf("Expected Exec=%s in .desktop file, got:\n%s", tt.exec, string(content))
			}
		})
	}
}

func TestGenerateDotDesktopDefaultOutputFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "wails-dotdesktop-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	options := &DotDesktopOptions{
		Name: "My App",
		Exec: "my-app",
	}
	DisableFooter = true

	oldDir, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(oldDir)

	err = GenerateDotDesktop(options)
	if err != nil {
		t.Fatalf("GenerateDotDesktop() error = %v", err)
	}

	expectedFile := "my-app.desktop"
	if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
		t.Errorf("Expected output file %s to be created", expectedFile)
	}

	content, err := os.ReadFile(expectedFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	if !strings.Contains(string(content), "Name=My App") {
		t.Errorf("Expected 'Name=My App' in .desktop file (name should NOT be normalised), got:\n%s", string(content))
	}
}

func TestGenerateDotDesktopErrors(t *testing.T) {
	tests := []struct {
		name    string
		options *DotDesktopOptions
		wantErr bool
	}{
		{
			name: "missing name",
			options: &DotDesktopOptions{
				Exec: "my-app",
			},
			wantErr: true,
		},
		{
			name: "missing exec",
			options: &DotDesktopOptions{
				Name: "My App",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := GenerateDotDesktop(tt.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateDotDesktop() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
