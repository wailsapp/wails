package commands

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wailsapp/task/v3/taskfile/ast"
)

func TestParseTaskAndVars(t *testing.T) {
	tests := []struct {
		name           string
		options        *RunTaskOptions
		otherArgs      []string
		osArgs         []string
		expectedTask   string
		expectedVars   map[string]string
	}{
		{
			name:         "Task name in options with CLI variables",
			options:      &RunTaskOptions{Name: "build"},
			otherArgs:    []string{"PLATFORM=linux", "CONFIG=production"},
			expectedTask: "build",
			expectedVars: map[string]string{
				"PLATFORM": "linux",
				"CONFIG":   "production",
			},
		},
		{
			name:         "Task name and variables in otherArgs",
			options:      &RunTaskOptions{},
			otherArgs:    []string{"test", "ENV=staging", "DEBUG=true"},
			expectedTask: "test",
			expectedVars: map[string]string{
				"ENV":   "staging",
				"DEBUG": "true",
			},
		},
		{
			name:         "Only task name, no variables",
			options:      &RunTaskOptions{},
			otherArgs:    []string{"deploy"},
			expectedTask: "deploy",
			expectedVars: map[string]string{},
		},
		{
			name:         "Default task when no args provided",
			options:      &RunTaskOptions{},
			otherArgs:    []string{},
			osArgs:       []string{"wails3", "task"}, // Set explicit os.Args to avoid test framework interference
			expectedTask: "default",
			expectedVars: map[string]string{},
		},
		{
			name:         "Variables with equals signs in values",
			options:      &RunTaskOptions{Name: "build"},
			otherArgs:    []string{"URL=https://example.com?key=value", "CONFIG=key1=val1,key2=val2"},
			expectedTask: "build",
			expectedVars: map[string]string{
				"URL":    "https://example.com?key=value",
				"CONFIG": "key1=val1,key2=val2",
			},
		},
		{
			name:         "Skip non-variable arguments",
			options:      &RunTaskOptions{Name: "build"},
			otherArgs:    []string{"PLATFORM=linux", "some-arg", "CONFIG=debug", "--flag"},
			expectedTask: "build",
			expectedVars: map[string]string{
				"PLATFORM": "linux",
				"CONFIG":   "debug",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original os.Args
			originalArgs := os.Args
			defer func() { os.Args = originalArgs }()

			if tt.osArgs != nil {
				os.Args = tt.osArgs
			}

			// Parse the task and variables
			call := parseTaskCall(tt.options, tt.otherArgs)

			// Verify task name
			assert.Equal(t, tt.expectedTask, call.Task)

			// Verify variables
			if len(tt.expectedVars) > 0 {
				require.NotNil(t, call.Vars)
				
				// Check each expected variable
				for key, expectedValue := range tt.expectedVars {
					var actualValue string
					call.Vars.Range(func(k string, v ast.Var) error {
						if k == key {
							actualValue = v.Value.(string)
						}
						return nil
					})
					assert.Equal(t, expectedValue, actualValue, "Variable %s mismatch", key)
				}
			} else if call.Vars != nil {
				// Ensure no variables were set when none expected
				count := 0
				call.Vars.Range(func(k string, v ast.Var) error {
					count++
					return nil
				})
				assert.Equal(t, 0, count, "Expected no variables but found %d", count)
			}
		})
	}
}

// Helper function to extract the task parsing logic for testing
func parseTaskCall(options *RunTaskOptions, otherArgs []string) *ast.Call {
	var tasksAndVars []string
	
	// Check if we have a task name specified in options
	if options.Name != "" {
		// If task name is provided via options, use it and treat otherArgs as CLI variables
		tasksAndVars = append([]string{options.Name}, otherArgs...)
	} else if len(otherArgs) > 0 {
		// Use otherArgs directly if provided
		tasksAndVars = otherArgs
	} else {
		// Fall back to parsing os.Args for backward compatibility
		var index int
		var arg string
		for index, arg = range os.Args[2:] {
			if len(arg) > 0 && arg[0] != '-' {
				break
			}
		}

		for _, taskAndVar := range os.Args[index+2:] {
			if taskAndVar == "--" {
				break
			}
			tasksAndVars = append(tasksAndVars, taskAndVar)
		}
	}

	// Default task
	if len(tasksAndVars) == 0 {
		tasksAndVars = []string{"default"}
	}

	// Parse task name and CLI variables
	taskName := tasksAndVars[0]
	cliVars := tasksAndVars[1:]
	
	// Create call with CLI variables
	call := &ast.Call{
		Task: taskName,
		Vars: &ast.Vars{},
	}
	
	// Parse CLI variables (format: KEY=VALUE)
	for _, v := range cliVars {
		if idx := findEquals(v); idx != -1 {
			key := v[:idx]
			value := v[idx+1:]
			call.Vars.Set(key, ast.Var{
				Value: value,
			})
		}
	}

	return call
}

// Helper to find the first equals sign
func findEquals(s string) int {
	for i, r := range s {
		if r == '=' {
			return i
		}
	}
	return -1
}