package doctor

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
)

// DiagnosticTest represents a single diagnostic test to be run
type DiagnosticTest struct {
	Name    string
	Run     func() (bool, string) // Returns success and error message if failed
	HelpURL string
}

// DiagnosticResult represents the result of a diagnostic test
type DiagnosticResult struct {
	TestName string
	ErrorMsg string
	HelpURL  string
}

// platformDiagnostics maps platform names to their diagnostic tests
var platformDiagnostics = map[string][]DiagnosticTest{
	// Tests that run on all platforms
	"all": {
		{
			Name: "Check Go installation",
			Run: func() (bool, string) {
				// This is just an example test for all platforms
				if runtime.Version() == "" {
					return false, "Go installation not found"
				}
				return true, ""
			},
			HelpURL: "/getting-started/installation/",
		},
	},
	// Platform specific tests
	"darwin": {
		{
			Name: "Check for .syso file",
			Run: func() (bool, string) {
				// Check for .syso files in current directory
				matches, err := filepath.Glob("*.syso")
				if err != nil {
					return false, "Error checking for .syso files"
				}
				if len(matches) > 0 {
					return false, fmt.Sprintf("Found .syso file(s): %v. These may cause issues when building on macOS", strings.Join(matches, ", "))
				}
				return true, ""
			},
			HelpURL: "/troubleshooting/mac-syso",
		},
	},
}

// RunDiagnostics executes all diagnostic tests for the current platform
func RunDiagnostics() []DiagnosticResult {
	var results []DiagnosticResult

	// Run tests that apply to all platforms
	if tests, exists := platformDiagnostics["all"]; exists {
		for _, test := range tests {
			success, errMsg := test.Run()
			if !success {
				results = append(results, DiagnosticResult{
					TestName: test.Name,
					ErrorMsg: errMsg,
					HelpURL:  test.HelpURL,
				})
			}
		}
	}

	// Run platform-specific tests
	if tests, exists := platformDiagnostics[runtime.GOOS]; exists {
		for _, test := range tests {
			success, errMsg := test.Run()
			if !success {
				results = append(results, DiagnosticResult{
					TestName: test.Name,
					ErrorMsg: errMsg,
					HelpURL:  test.HelpURL,
				})
			}
		}
	}

	return results
}
