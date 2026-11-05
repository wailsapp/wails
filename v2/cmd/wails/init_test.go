package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRemoteTemplateWarning(t *testing.T) {
	dir := t.TempDir()
	wailsJSON := `{
  "name": "test",
  "frontend:install": "npm install",
  "frontend:build": "sh ./scripts/build.sh",
  "frontend:dev:watcher": "npm run dev"
}`
	if err := os.WriteFile(filepath.Join(dir, "wails.json"), []byte(wailsJSON), 0o644); err != nil {
		t.Fatal(err)
	}

	got := remoteTemplateWarning(dir)
	for _, want := range []string{
		"remote template",
		"wails build",
		"wails dev",
		"Review wails.json",
		"run from 'frontend'",
		"frontend:install: npm install",
		"frontend:build: sh ./scripts/build.sh",
		"frontend:dev:watcher: npm run dev",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("warning missing %q:\n%s", want, got)
		}
	}
	if strings.Contains(got, "frontend:dev:build") {
		t.Errorf("warning lists unset command:\n%s", got)
	}
}

func TestRemoteTemplateWarningWithoutWailsJSON(t *testing.T) {
	got := remoteTemplateWarning(t.TempDir())
	if !strings.Contains(got, "Review wails.json") {
		t.Errorf("unexpected warning:\n%s", got)
	}
	if strings.Contains(got, "Commands in wails.json") {
		t.Errorf("warning lists commands without a wails.json:\n%s", got)
	}
}
