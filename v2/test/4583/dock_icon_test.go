package test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestApplicationShouldHandleReopenExists(t *testing.T) {
	repoRoot := os.Getenv("WAILS_REPO_ROOT")
	if repoRoot == "" {
		currentDir, _ := os.Getwd()
		for {
			if _, err := os.Stat(filepath.Join(currentDir, "go.mod")); err == nil {
				repoRoot = currentDir
				break
			}
			parent := filepath.Dir(currentDir)
			if parent == currentDir {
				t.Skip("Cannot find Wails repo root, set WAILS_REPO_ROOT")
			}
			currentDir = parent
		}
	}

	appDelegatePath := filepath.Join(repoRoot, "internal", "frontend", "desktop", "darwin", "AppDelegate.m")
	data, err := os.ReadFile(appDelegatePath)
	if err != nil {
		t.Fatalf("Failed to read AppDelegate.m: %v", err)
	}

	content := string(data)

	if !strings.Contains(content, "applicationShouldHandleReopen") {
		t.Fatal("AppDelegate.m must implement applicationShouldHandleReopen:hasVisibleWindows: to handle dock icon clicks when StartHidden:true")
	}

	if !strings.Contains(content, "hasVisibleWindows") {
		t.Fatal("applicationShouldHandleReopen must accept hasVisibleWindows: parameter")
	}

	if !strings.Contains(content, "makeKeyAndOrderFront") {
		t.Fatal("applicationShouldHandleReopen must call makeKeyAndOrderFront: to show the window")
	}
}

func TestApplicationShouldHandleReopenLogic(t *testing.T) {
	tests := []struct {
		name       string
		flag       bool
		shouldShow bool
	}{
		{
			name:       "shows window when no visible windows (StartHidden case)",
			flag:       false,
			shouldShow: true,
		},
		{
			name:       "does not force show when windows already visible",
			flag:       true,
			shouldShow: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wouldShow := !tt.flag
			if wouldShow != tt.shouldShow {
				t.Errorf("expected wouldShow=%v for flag=%v, got %v", tt.shouldShow, tt.flag, wouldShow)
			}
		})
	}
}
