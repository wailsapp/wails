package test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"howett.net/plist"
)

func TestDarwinInfoPlistMinVersion(t *testing.T) {
	templates := []string{
		"../../internal/commands/updatable_build_assets/darwin/Info.plist.tmpl",
		"../../internal/commands/updatable_build_assets/darwin/Info.dev.plist.tmpl",
	}

	expectedVersion := "12.0"

	for _, tmpl := range templates {
		t.Run(filepath.Base(tmpl), func(t *testing.T) {
			data, err := os.ReadFile(tmpl)
			if err != nil {
				t.Fatalf("Failed to read template: %v", err)
			}

			content := string(data)

			if !strings.Contains(content, expectedVersion) {
				t.Errorf("Template %s does not contain expected minimum version %s", tmpl, expectedVersion)
			}

			if strings.Contains(content, "10.15") {
				t.Errorf("Template %s still contains outdated 10.15 version", tmpl)
			}

			var plistDict map[string]any
			_, err = plist.Unmarshal(data, &plistDict)
			if err != nil {
				t.Fatalf("Failed to unmarshal plist template %s: %v", tmpl, err)
			}

			minVer, ok := plistDict["LSMinimumSystemVersion"]
			if !ok {
				t.Fatalf("Template %s is missing LSMinimumSystemVersion", tmpl)
			}
			if minVer != expectedVersion {
				t.Errorf("LSMinimumSystemVersion = %v, want %s", minVer, expectedVersion)
			}
		})
	}
}

func TestDarwinTaskfileMinVersion(t *testing.T) {
	taskfile := "../../internal/commands/build_assets/darwin/Taskfile.yml"

	data, err := os.ReadFile(taskfile)
	if err != nil {
		t.Fatalf("Failed to read Taskfile: %v", err)
	}

	content := string(data)

	if strings.Contains(content, "10.15") {
		t.Errorf("Taskfile still contains outdated 10.15 version")
	}

	if !strings.Contains(content, "12.0") {
		t.Errorf("Taskfile does not contain expected minimum version 12.0")
	}

	if !strings.Contains(content, "-mmacosx-version-min=12.0") {
		t.Errorf("Taskfile CGO flags not updated to 12.0")
	}

	if !strings.Contains(content, `MACOSX_DEPLOYMENT_TARGET: "12.0"`) {
		t.Errorf("Taskfile MACOSX_DEPLOYMENT_TARGET not updated to 12.0")
	}
}
