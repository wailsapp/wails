package setupwizard

import (
	"reflect"
	"testing"
)

func TestIsCommandAllowed(t *testing.T) {
	tests := []struct {
		name    string
		parts   []string
		allowed bool
	}{
		// Valid package manager commands
		{"apt install", []string{"apt", "install", "pkg"}, true},
		{"apt-get install", []string{"apt-get", "install", "pkg"}, true},
		{"brew install", []string{"brew", "install", "pkg"}, true},
		{"pacman -S", []string{"pacman", "-S", "pkg"}, true},
		{"dnf install", []string{"dnf", "install", "pkg"}, true},

		// Valid sudo commands
		{"sudo apt install", []string{"sudo", "apt", "install", "pkg"}, true},
		{"sudo apt-get install", []string{"sudo", "apt-get", "install", "pkg"}, true},
		{"sudo pacman -S", []string{"sudo", "pacman", "-S", "pkg"}, true},
		{"pkexec apt install", []string{"pkexec", "apt", "install", "pkg"}, true},
		{"doas apt install", []string{"doas", "apt", "install", "pkg"}, true},

		// CRITICAL: Bypass attempts that MUST be blocked
		{"sudo -u apt bash (bypass)", []string{"sudo", "-u", "apt", "bash", "-c", "malicious"}, false},
		{"sudo -E bash", []string{"sudo", "-E", "bash"}, false},
		{"sudo --user=root bash", []string{"sudo", "--user=root", "bash"}, false},
		{"doas -u apt bash", []string{"doas", "-u", "apt", "bash"}, false},
		{"pkexec --user apt bash", []string{"pkexec", "--user", "apt", "bash"}, false},

		// Invalid commands
		{"bash", []string{"bash", "-c", "malicious"}, false},
		{"rm -rf", []string{"rm", "-rf", "/"}, false},
		{"curl", []string{"curl", "http://evil.com"}, false},
		{"wget", []string{"wget", "http://evil.com"}, false},
		{"empty", []string{}, false},
		{"sudo only", []string{"sudo"}, false},

		// Nested sudo attempts
		{"sudo sudo apt", []string{"sudo", "sudo", "apt"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isCommandAllowed(tt.parts)
			if got != tt.allowed {
				t.Errorf("isCommandAllowed(%v) = %v, want %v", tt.parts, got, tt.allowed)
			}
		})
	}
}

func TestGetSafeCommand(t *testing.T) {
	tests := []struct {
		name            string
		parts           []string
		wantSafeCmd     string
		wantElevatedCmd string
		wantArgs        []string
		wantOk          bool
	}{
		// Direct package manager commands
		{
			name:            "apt install",
			parts:           []string{"apt", "install", "pkg"},
			wantSafeCmd:     "apt",
			wantElevatedCmd: "",
			wantArgs:        []string{"install", "pkg"},
			wantOk:          true,
		},
		{
			name:            "brew install",
			parts:           []string{"brew", "install", "pkg"},
			wantSafeCmd:     "brew",
			wantElevatedCmd: "",
			wantArgs:        []string{"install", "pkg"},
			wantOk:          true,
		},
		// Sudo commands - verify elevated command comes from whitelist
		{
			name:            "sudo apt install",
			parts:           []string{"sudo", "apt", "install", "pkg"},
			wantSafeCmd:     "sudo",
			wantElevatedCmd: "apt",
			wantArgs:        []string{"install", "pkg"},
			wantOk:          true,
		},
		{
			name:            "pkexec pacman -S",
			parts:           []string{"pkexec", "pacman", "-S", "pkg"},
			wantSafeCmd:     "pkexec",
			wantElevatedCmd: "pacman",
			wantArgs:        []string{"-S", "pkg"},
			wantOk:          true,
		},
		// Bypass attempts
		{
			name:            "sudo -u bypass",
			parts:           []string{"sudo", "-u", "apt", "bash"},
			wantSafeCmd:     "",
			wantElevatedCmd: "",
			wantArgs:        nil,
			wantOk:          false,
		},
		// Invalid commands
		{
			name:            "bash command",
			parts:           []string{"bash", "-c", "evil"},
			wantSafeCmd:     "",
			wantElevatedCmd: "",
			wantArgs:        nil,
			wantOk:          false,
		},
		{
			name:            "empty",
			parts:           []string{},
			wantSafeCmd:     "",
			wantElevatedCmd: "",
			wantArgs:        nil,
			wantOk:          false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			safeCmd, elevatedCmd, args, ok := getSafeCommand(tt.parts)
			if ok != tt.wantOk {
				t.Errorf("getSafeCommand() ok = %v, want %v", ok, tt.wantOk)
			}
			if safeCmd != tt.wantSafeCmd {
				t.Errorf("getSafeCommand() safeCmd = %v, want %v", safeCmd, tt.wantSafeCmd)
			}
			if elevatedCmd != tt.wantElevatedCmd {
				t.Errorf("getSafeCommand() elevatedCmd = %v, want %v", elevatedCmd, tt.wantElevatedCmd)
			}
			if !reflect.DeepEqual(args, tt.wantArgs) {
				t.Errorf("getSafeCommand() args = %v, want %v", args, tt.wantArgs)
			}
		})
	}
}
