//go:build linux
// +build linux

package packagemanager

import (
	"strings"
	"testing"
)

func TestApt_RemoveEscapeSequences(t *testing.T) {
	apt := NewApt("debian")

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "plain installed",
			input: "libgtk-3-dev/now 3.24.38 amd64 [installed]",
			want:  "libgtk-3-dev/now 3.24.38 amd64 [installed]",
		},
		{
			name:  "installed local",
			input: "libgtk-3-dev/now 3.24.38 amd64 [installed,local]",
			want:  "libgtk-3-dev/now 3.24.38 amd64 [installed,local]",
		},
		{
			name:  "green color codes",
			input: "\x1b[32mlibgtk-3-dev/now 3.24.38 amd64 [installed]\x1b[0m",
			want:  "libgtk-3-dev/now 3.24.38 amd64 [installed]",
		},
		{
			name:  "bold color codes around installed",
			input: "libgtk-3-dev/now 3.24.38 amd64 \x1b[1m[installed]\x1b[0m",
			want:  "libgtk-3-dev/now 3.24.38 amd64 [installed]",
		},
		{
			name:  "ansi escapes in status field",
			input: "\x1b[32mlibgtk-3-dev\x1b[0m/now 3.24.38 amd64 [\x1b[1minstalled,local\x1b[0m]",
			want:  "libgtk-3-dev/now 3.24.38 amd64 [installed,local]",
		},
		{
			name:  "no escapes plain text",
			input: "some plain text",
			want:  "some plain text",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := apt.removeEscapeSequences(tt.input)
			if got != tt.want {
				t.Errorf("removeEscapeSequences() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestApt_RemoveEscapeSequences_InstalledCheck(t *testing.T) {
	apt := NewApt("debian")

	tests := []struct {
		name      string
		input     string
		wantFound bool
	}{
		{
			name:      "installed with escape codes",
			input:     "\x1b[32mlibgtk-3-dev/now 3.24.38 amd64 [\x1b[1minstalled\x1b[0m]\x1b[0m",
			wantFound: true,
		},
		{
			name:      "installed local with escape codes",
			input:     "\x1b[32mlibgtk-3-dev/now 3.24.38 amd64 [\x1b[1minstalled,local\x1b[0m]\x1b[0m",
			wantFound: true,
		},
		{
			name:      "installed automatic with escape codes",
			input:     "\x1b[32mlibgtk-3-dev/now 3.24.38 amd64 [\x1b[1minstalled,automatic\x1b[0m]\x1b[0m",
			wantFound: true,
		},
		{
			name:      "not installed",
			input:     "libgtk-3-dev/stable 3.24.38 amd64",
			wantFound: false,
		},
		{
			name:      "plain installed",
			input:     "libgtk-3-dev/now 3.24.38 amd64 [installed]",
			wantFound: true,
		},
		{
			name:      "plain installed local",
			input:     "libgtk-3-dev/now 3.24.38 amd64 [installed,local]",
			wantFound: true,
		},
		{
			name:      "plain installed upgradable",
			input:     "libgtk-3-dev/now 3.24.38 amd64 [installed,upgradable to: 3.24.40]",
			wantFound: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleaned := apt.removeEscapeSequences(tt.input)
			found := strings.Contains(cleaned, "[installed")
			if found != tt.wantFound {
				t.Errorf("strings.Contains(removeEscapeSequences(%q), [installed) = %v, want %v", tt.input, found, tt.wantFound)
			}
		})
	}
}
