package commands

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestToolVersion(t *testing.T) {
	// Create a table of test cases
	testCases := []struct {
		name           string
		options        ToolVersionOptions
		expectedOutput string
		expectedError  bool
	}{
		{
			name: "Bump major version",
			options: ToolVersionOptions{
				Version: "1.2.3",
				Major:   true,
			},
			expectedOutput: "2.0.0",
			expectedError:  false,
		},
		{
			name: "Bump minor version",
			options: ToolVersionOptions{
				Version: "1.2.3",
				Minor:   true,
			},
			expectedOutput: "1.3.0",
			expectedError:  false,
		},
		{
			name: "Bump patch version",
			options: ToolVersionOptions{
				Version: "1.2.3",
				Patch:   true,
			},
			expectedOutput: "1.2.4",
			expectedError:  false,
		},
		{
			name: "Bump major version with v prefix",
			options: ToolVersionOptions{
				Version: "v1.2.3",
				Major:   true,
			},
			expectedOutput: "v2.0.0",
			expectedError:  false,
		},
		{
			name: "Bump minor version with v prefix",
			options: ToolVersionOptions{
				Version: "v1.2.3",
				Minor:   true,
			},
			expectedOutput: "v1.3.0",
			expectedError:  false,
		},
		{
			name: "Bump patch version with v prefix",
			options: ToolVersionOptions{
				Version: "v1.2.3",
				Patch:   true,
			},
			expectedOutput: "v1.2.4",
			expectedError:  false,
		},
		{
			name: "Bump version with prerelease",
			options: ToolVersionOptions{
				Version: "1.2.3-alpha",
				Patch:   true,
			},
			expectedOutput: "1.2.4-alpha",
			expectedError:  false,
		},
		{
			name: "Bump version with metadata",
			options: ToolVersionOptions{
				Version: "1.2.3+build123",
				Patch:   true,
			},
			expectedOutput: "1.2.4+build123",
			expectedError:  false,
		},
		{
			name: "Bump version with prerelease and metadata",
			options: ToolVersionOptions{
				Version: "1.2.3-alpha+build123",
				Patch:   true,
			},
			expectedOutput: "1.2.4-alpha+build123",
			expectedError:  false,
		},
		{
			name: "Bump version with v prefix, prerelease and metadata",
			options: ToolVersionOptions{
				Version: "v1.2.3-alpha+build123",
				Patch:   true,
			},
			expectedOutput: "v1.2.4-alpha+build123",
			expectedError:  false,
		},
		{
			name: "No version provided",
			options: ToolVersionOptions{
				Major: true,
			},
			expectedOutput: "",
			expectedError:  true,
		},
		{
			name: "No bump flag provided",
			options: ToolVersionOptions{
				Version: "1.2.3",
			},
			expectedOutput: "",
			expectedError:  true,
		},
		{
			name: "Invalid version format",
			options: ToolVersionOptions{
				Version: "invalid",
				Major:   true,
			},
			expectedOutput: "",
			expectedError:  true,
		},
		{
			name: "Bump prerelease version with numeric component",
			options: ToolVersionOptions{
				Version:    "1.2.3-alpha.5",
				Prerelease: true,
			},
			expectedOutput: "1.2.3-alpha.6",
			expectedError:  false,
		},
		{
			name: "Bump prerelease version with v prefix",
			options: ToolVersionOptions{
				Version:    "v1.2.3-alpha.5",
				Prerelease: true,
			},
			expectedOutput: "v1.2.3-alpha.6",
			expectedError:  false,
		},
		{
			name: "Bump prerelease version with metadata",
			options: ToolVersionOptions{
				Version:    "1.2.3-alpha.5+build123",
				Prerelease: true,
			},
			expectedOutput: "1.2.3-alpha.6+build123",
			expectedError:  false,
		},
		{
			name: "Bump prerelease version without numeric component",
			options: ToolVersionOptions{
				Version:    "1.2.3-alpha",
				Prerelease: true,
			},
			expectedOutput: "1.2.3-alpha",
			expectedError:  false,
		},
		{
			name: "Bump prerelease version when no prerelease part exists",
			options: ToolVersionOptions{
				Version:    "1.2.3",
				Prerelease: true,
			},
			expectedOutput: "",
			expectedError:  true,
		},
		{
			name: "Bump prerelease version for issue example v3.0.0-alpha.5",
			options: ToolVersionOptions{
				Version:    "v3.0.0-alpha.5",
				Prerelease: true,
			},
			expectedOutput: "v3.0.0-alpha.6",
			expectedError:  false,
		},
	}

	// Run each test case
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Capture stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Call the function
			err := ToolVersion(&tc.options)

			// Restore stdout
			err2 := w.Close()
			if err2 != nil {
				t.Fail()
			}
			os.Stdout = oldStdout

			// Read captured output
			var buf bytes.Buffer
			_, err2 = io.Copy(&buf, r)
			if err2 != nil {
				t.Fail()
			}

			output := strings.TrimSpace(buf.String())

			// Check error
			if tc.expectedError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tc.expectedError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}

			// Check output
			if !tc.expectedError && output != tc.expectedOutput {
				t.Errorf("Expected output '%s' but got '%s'", tc.expectedOutput, output)
			}
		})
	}
}

func TestBumpPrerelease(t *testing.T) {
	cases := []struct{ input, want string }{
		{"", ""},
		{"alpha", "alpha"},
		{"0", "1"},
		{"9", "10"},
		{"alpha.5", "alpha.6"},
		{"alpha.009", "alpha.10"},
		{"alpha.000", "alpha.1"},
		{"alpha.1299", "alpha.1300"},
		{"alpha.", "alpha."},
		{"alpha.１２", "alpha.１２"},
		{"alpha.9beta", "alpha.9beta"},
		{"alpha.-1", "alpha.-1"},
		{"alpha.9223372036854775807", "alpha.9223372036854775808"},
		{"alpha.18446744073709551616", "alpha.18446744073709551617"},
		{"999999999999999999999999999999", "1000000000000000000000000000000"},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			if got := bumpPrerelease(tc.input); got != tc.want {
				t.Errorf("bumpPrerelease(%q)=%q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

func BenchmarkBumpPrerelease(b *testing.B) {
	for _, tc := range []struct{ name, input, want string }{
		{"nightly", "nightly.123", "nightly.124"},
		{"long_prefix", strings.Repeat("preview.", 32) + "123", strings.Repeat("preview.", 32) + "124"},
		{"large_carry", "alpha." + strings.Repeat("9", 128), "alpha.1" + strings.Repeat("0", 128)},
		{"nonnumeric", "alpha.preview", "alpha.preview"},
	} {
		b.Run(tc.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				if got := bumpPrerelease(tc.input); got != tc.want {
					b.Fatalf("got %q, want %q", got, tc.want)
				}
			}
		})
	}
}
