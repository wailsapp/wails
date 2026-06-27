//go:build linux
// +build linux

package packagemanager

import (
	"strings"
	"testing"
)

func TestDnfListInstalledVersionParsing(t *testing.T) {
	tests := []struct {
		name     string
		output   string
		expected string
	}{
		{
			name:     "package with release suffix",
			output:   "webkit2gtk4.0-devel.x86_64 2.46.5-1.fc41 @updates",
			expected: "2.46.5",
		},
		{
			name:     "package simple release",
			output:   "gcc-c++.x86_64 14.3.1-1.fc41 @updates",
			expected: "14.3.1",
		},
		{
			name:     "package no dash in version",
			output:   "somepkg.x86_64 1.0 @repo",
			expected: "1.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got string
			splitoutput := strings.Split(tt.output, "\n")
			for _, line := range splitoutput {
				fields := strings.Fields(line)
				if len(fields) >= 2 {
					version := fields[1]
					if idx := strings.Index(version, "-"); idx != -1 {
						got = version[:idx]
					} else {
						got = version
					}
				}
			}
			if got != tt.expected {
				t.Errorf("expected version %q, got %q", tt.expected, got)
			}
		})
	}
}
