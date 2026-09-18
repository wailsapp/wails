package commands

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wailsapp/task/v3/taskfile/ast"
)

func TestTaskParameterPassing(t *testing.T) {
	// Skip if running in CI without proper environment
	if os.Getenv("CI") == "true" && os.Getenv("SKIP_INTEGRATION_TESTS") == "true" {
		t.Skip("Skipping integration test in CI")
	}

	// Create a temporary directory for test
	tmpDir, err := os.MkdirTemp("", "wails-task-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Create a test Taskfile
	taskfileContent := `version: '3'

tasks:
  build:
    cmds:
      - echo "PLATFORM={{.PLATFORM | default "default-platform"}}"
      - echo "CONFIG={{.CONFIG | default "default-config"}}"
    silent: true

  package:
    cmds:
      - echo "VERSION={{.VERSION | default "1.0.0"}}"
      - echo "OUTPUT={{.OUTPUT | default "output.pkg"}}"
    silent: true

  test:
    cmds:
      - echo "ENV={{.ENV | default "test"}}"
      - echo "FLAGS={{.FLAGS | default "none"}}"
    silent: true
`

	taskfilePath := filepath.Join(tmpDir, "Taskfile.yml")
	err = os.WriteFile(taskfilePath, []byte(taskfileContent), 0644)
	require.NoError(t, err)

	// Save current directory
	originalWd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(originalWd)

	// Change to test directory
	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	tests := []struct {
		name           string
		options        *RunTaskOptions
		otherArgs      []string
		expectedOutput []string
	}{
		{
			name:      "Build task with parameters",
			options:   &RunTaskOptions{Name: "build"},
			otherArgs: []string{"PLATFORM=linux", "CONFIG=production"},
			expectedOutput: []string{
				"PLATFORM=linux",
				"CONFIG=production",
			},
		},
		{
			name:      "Package task with parameters",
			options:   &RunTaskOptions{Name: "package"},
			otherArgs: []string{"VERSION=2.5.0", "OUTPUT=myapp.pkg"},
			expectedOutput: []string{
				"VERSION=2.5.0",
				"OUTPUT=myapp.pkg",
			},
		},
		{
			name:      "Task with default values",
			options:   &RunTaskOptions{Name: "build"},
			otherArgs: []string{},
			expectedOutput: []string{
				"PLATFORM=default-platform",
				"CONFIG=default-config",
			},
		},
		{
			name:      "Task with partial parameters",
			options:   &RunTaskOptions{Name: "test"},
			otherArgs: []string{"ENV=staging"},
			expectedOutput: []string{
				"ENV=staging",
				"FLAGS=none",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture output
			output := captureTaskOutput(t, tt.options, tt.otherArgs)

			// Verify expected output
			for _, expected := range tt.expectedOutput {
				assert.Contains(t, output, expected, "Output should contain: %s", expected)
			}
		})
	}
}

func TestCLIParameterFormats(t *testing.T) {
	tests := []struct {
		name         string
		otherArgs    []string
		expectError  bool
		expectedVars map[string]string
	}{
		{
			name:      "Standard KEY=VALUE format",
			otherArgs: []string{"build", "KEY1=value1", "KEY2=value2"},
			expectedVars: map[string]string{
				"KEY1": "value1",
				"KEY2": "value2",
			},
		},
		{
			name:      "Values with equals signs",
			otherArgs: []string{"build", "URL=https://example.com?key=value", "FORMULA=a=b+c"},
			expectedVars: map[string]string{
				"URL":     "https://example.com?key=value",
				"FORMULA": "a=b+c",
			},
		},
		{
			name:      "Values with spaces (quoted)",
			otherArgs: []string{"build", "MESSAGE=Hello World", "PATH=/usr/local/bin"},
			expectedVars: map[string]string{
				"MESSAGE": "Hello World",
				"PATH":    "/usr/local/bin",
			},
		},
		{
			name:      "Mixed valid and invalid arguments",
			otherArgs: []string{"build", "VALID=yes", "invalid-arg", "ANOTHER=value", "--flag"},
			expectedVars: map[string]string{
				"VALID":   "yes",
				"ANOTHER": "value",
			},
		},
		{
			name:      "Empty value",
			otherArgs: []string{"build", "EMPTY=", "KEY=value"},
			expectedVars: map[string]string{
				"EMPTY": "",
				"KEY":   "value",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			call := parseTaskCall(&RunTaskOptions{}, tt.otherArgs)

			// Verify variables
			for key, expectedValue := range tt.expectedVars {
				var actualValue string
				found := false
				if call.Vars != nil {
					call.Vars.Range(func(k string, v ast.Var) error {
						if k == key {
							actualValue = v.Value.(string)
							found = true
						}
						return nil
					})
				}
				assert.True(t, found, "Variable %s not found", key)
				assert.Equal(t, expectedValue, actualValue, "Variable %s mismatch", key)
			}
		})
	}
}

// Helper function to capture task output
func captureTaskOutput(t *testing.T, options *RunTaskOptions, otherArgs []string) string {
	// Save original stdout and stderr
	oldStdout := os.Stdout
	oldStderr := os.Stderr
	defer func() {
		os.Stdout = oldStdout
		os.Stderr = oldStderr
	}()

	// Create pipe to capture output
	r, w, err := os.Pipe()
	require.NoError(t, err)

	os.Stdout = w
	os.Stderr = w

	// Run task in a goroutine
	done := make(chan bool)
	var taskErr error
	go func() {
		// Note: This is a simplified version for testing
		// In real tests, you might want to mock the Task executor
		taskErr = RunTask(options, otherArgs)
		w.Close()
		done <- true
	}()

	// Read output
	var buf bytes.Buffer
	_, err = buf.ReadFrom(r)
	require.NoError(t, err)

	// Wait for task to complete
	<-done

	// Check for errors (might be expected in some tests)
	if taskErr != nil && !strings.Contains(taskErr.Error(), "expected") {
		t.Logf("Task error (might be expected): %v", taskErr)
	}

	return buf.String()
}

func TestBackwardCompatibility(t *testing.T) {
	// Test that the old way of calling tasks still works
	tests := []struct {
		name         string
		osArgs       []string
		expectedTask string
		expectedVars map[string]string
	}{
		{
			name:         "Legacy os.Args parsing",
			osArgs:       []string{"wails3", "task", "build", "PLATFORM=windows"},
			expectedTask: "build",
			expectedVars: map[string]string{
				"PLATFORM": "windows",
			},
		},
		{
			name:         "Legacy with flags before task",
			osArgs:       []string{"wails3", "task", "--verbose", "test", "ENV=prod"},
			expectedTask: "test",
			expectedVars: map[string]string{
				"ENV": "prod",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original os.Args
			originalArgs := os.Args
			defer func() { os.Args = originalArgs }()

			os.Args = tt.osArgs

			// Parse using the backward compatibility path
			call := parseTaskCall(&RunTaskOptions{}, []string{})

			assert.Equal(t, tt.expectedTask, call.Task)

			for key, expectedValue := range tt.expectedVars {
				var actualValue string
				if call.Vars != nil {
					call.Vars.Range(func(k string, v ast.Var) error {
						if k == key {
							actualValue = v.Value.(string)
						}
						return nil
					})
				}
				assert.Equal(t, expectedValue, actualValue, "Variable %s mismatch", key)
			}
		})
	}
}

func TestMkdirWithSpacesInPath(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("Skipping: macOS app bundle test only applies to darwin")
	}
	if os.Getenv("CI") == "true" && os.Getenv("SKIP_INTEGRATION_TESTS") == "true" {
		t.Skip("Skipping integration test in CI")
	}

	tmpDir, err := os.MkdirTemp("", "wails task test with spaces-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	taskfileContent := `version: '3'

vars:
  BIN_DIR: "` + tmpDir + `/bin"
  APP_NAME: "My App"

tasks:
  create-bundle:
    cmds:
      - mkdir -p "{{.BIN_DIR}}/{{.APP_NAME}}.app/Contents/MacOS"
      - mkdir -p "{{.BIN_DIR}}/{{.APP_NAME}}.app/Contents/Resources"
`

	taskfilePath := filepath.Join(tmpDir, "Taskfile.yml")
	err = os.WriteFile(taskfilePath, []byte(taskfileContent), 0644)
	require.NoError(t, err)

	originalWd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(originalWd)

	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	err = RunTask(&RunTaskOptions{Name: "create-bundle"}, []string{})
	require.NoError(t, err)

	appContentsDir := filepath.Join(tmpDir, "bin", "My App.app", "Contents")

	macOSDir := filepath.Join(appContentsDir, "MacOS")
	info, err := os.Stat(macOSDir)
	require.NoError(t, err, "MacOS directory should exist")
	assert.True(t, info.IsDir(), "MacOS should be a directory")

	resourcesDir := filepath.Join(appContentsDir, "Resources")
	info, err = os.Stat(resourcesDir)
	require.NoError(t, err, "Resources directory should exist")
	assert.True(t, info.IsDir(), "Resources should be a directory")
}
