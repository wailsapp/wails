package buildassets

import (
	"testing"
)

func TestSafeBundleID(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"My Fab Application", "my-fab-application"},
		{"SimpleApp", "simpleapp"},
		{"App-With-Dashes", "app-with-dashes"},
		{"App_With_Underscores", "app-with-underscores"},
		{"App@#$%^&*()", "app"},
		{"123App", "123app"},
		{"  Spaces  ", "spaces"},
		{"already-lowercase", "already-lowercase"},
	}
	for _, tt := range tests {
		got := safeBundleID(tt.input)
		if got != tt.expected {
			t.Errorf("safeBundleID(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}
