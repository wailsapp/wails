//go:build windows

package application_test

import (
	"path/filepath"
	"testing"

	"github.com/matryer/is"
)

func TestCleanPath(t *testing.T) {
	i := is.New(t)
	tests := []struct {
		name      string
		inputPath string
		expected  string
	}{
		{
			name:      "path with double separators",
			inputPath: `C:\\temp\\folder`,
			expected:  `C:\temp\folder`,
		},
		{
			name:      "path with forward slashes",
			inputPath: `C://temp//folder`,
			expected:  `C:\temp\folder`,
		},
		{
			name:      "path with trailing separator",
			inputPath: `C:\\temp\\folder\\`,
			expected:  `C:\temp\folder`,
		},
		{
			name:      "path with escaped tab character",
			inputPath: `C:\\Users\\test\\tab.txt`,
			expected:  `C:\Users\test\tab.txt`,
		},
		{
			name:      "newline character",
			inputPath: `C:\\Users\\test\\newline\\n.txt`,
			expected:  `C:\Users\test\newline\n.txt`,
		},
		{
			name:      "UNC path with multiple separators",
			inputPath: `\\\\\\\\host\\share\\test.txt`,
			expected:  `\\\\host\share\test.txt`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleaned := filepath.Clean(tt.inputPath)
			i.Equal(cleaned, tt.expected)
		})
	}
}
