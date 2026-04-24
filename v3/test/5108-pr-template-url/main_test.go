package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPRTemplateFeedbackURL(t *testing.T) {
	repoRoot := filepath.Join("..", "..", "..")
	templatePath := filepath.Join(repoRoot, ".github", "pull_request_template.md")

	data, err := os.ReadFile(templatePath)
	if err != nil {
		t.Fatalf("Failed to read PR template at %s: %v", templatePath, err)
	}
	content := string(data)

	if strings.Contains(content, "v3alpha.wails.io/getting-started/feedback") {
		t.Error("PR template contains incorrect feedback URL with '/getting-started/' path; should be https://v3alpha.wails.io/feedback/")
	}
	if !strings.Contains(content, "https://v3alpha.wails.io/feedback/") {
		t.Error("PR template should contain correct feedback URL: https://v3alpha.wails.io/feedback/")
	}
}
