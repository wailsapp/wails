package cmd

import (
	"testing"
)

func TestSemanticVersion_IsPreRelease(t *testing.T) {
	tests := []struct {
		name    string
		version string
		want    bool
	}{
		{"v1.6.7-pre0", "v1.6.7-pre0", true},
		{"v2.6.7+pre0", "v2.6.7+pre0", false},
		{"v2.6.7", "v2.6.7", false},
		{"v2.0.0+alpha.1", "v2.0.0+alpha.1", false},
		{"v2.0.0-alpha.1", "v2.0.0-alpha.1", false},
		{"v1.6.7", "v1.6.7", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			semanticversion, err := NewSemanticVersion(tt.version)
			if err != nil {
				t.Errorf("Invalid semantic version: %s", semanticversion)
				return
			}
			s := &SemanticVersion{
				Version: semanticversion.Version,
			}
			if got := s.IsPreRelease(); got != tt.want {
				t.Errorf("IsPreRelease() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSemanticVersion_IsRelease(t *testing.T) {
	tests := []struct {
		name    string
		version string
		want    bool
	}{
		{"v1.6.7", "v1.6.7", true},
		{"v2.6.7-pre0", "v2.6.7-pre0", false},
		{"v2.6.7", "v2.6.7", false},
		{"v2.6.7+release", "v2.6.7+release", false},
		{"v2.0.0-alpha.1", "v2.0.0-alpha.1", false},
		{"v1.6.7-pre0", "v1.6.7-pre0", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			semanticversion, err := NewSemanticVersion(tt.version)
			if err != nil {
				t.Errorf("Invalid semantic version: %s", semanticversion)
				return
			}
			s := &SemanticVersion{
				Version: semanticversion.Version,
			}
			if got := s.IsRelease(); got != tt.want {
				t.Errorf("IsRelease() = %v, want %v", got, tt.want)
			}
		})
	}
}
