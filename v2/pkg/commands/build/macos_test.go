package build

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewMacOSBuildConfig(t *testing.T) {
	tests := []struct {
		name          string
		deployment    string
		staged        string
		wantStaged    bool
		wantErrorText string
	}{
		{name: "default target uses public WebKit", deployment: "13.0"},
		{name: "older targets use public WebKit", deployment: "12.2"},
		{name: "invalid target", deployment: "12.x", wantErrorText: "numeric macOS version"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			config, err := newMacOSBuildConfig(test.deployment, test.staged)
			if test.wantErrorText != "" {
				if err == nil || !strings.Contains(err.Error(), test.wantErrorText) {
					t.Fatalf("newMacOSBuildConfig() error = %v, want %q", err, test.wantErrorText)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if got := config.stagedFrameworks != ""; got != test.wantStaged {
				t.Fatalf("staged framework selection = %t, want %t", got, test.wantStaged)
			}
		})
	}

	staged := makeSafariStagedFrameworks(t)
	config, err := newMacOSBuildConfig("12.3", staged)
	if err != nil {
		t.Fatal(err)
	}
	if config.stagedFrameworks != staged {
		t.Fatalf("stagedFrameworks = %q, want %q", config.stagedFrameworks, staged)
	}
	for _, want := range []string{
		"-F" + staged,
		"-framework UniformTypeIdentifiers",
		"DYLD_VERSIONED_FRAMEWORK_PATH=" + staged,
		"DYLD_VERSIONED_LIBRARY_PATH=" + staged,
	} {
		if !strings.Contains(config.externalLinkerArg, want) {
			t.Errorf("externalLinkerArg missing %q: %s", want, config.externalLinkerArg)
		}
	}
	if _, err := newMacOSBuildConfig("12.3", t.TempDir()); err == nil {
		t.Fatal("Monterey target must fail when Safari staged frameworks are incomplete")
	}
}

func makeSafariStagedFrameworks(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	for _, framework := range safariStagedFrameworks {
		binary := filepath.Join(root, framework+".framework", "Versions", "A", framework)
		if err := os.MkdirAll(filepath.Dir(binary), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(binary, nil, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	return root
}
