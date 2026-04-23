package commands

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/atterpac/refresh/engine"
	"gopkg.in/yaml.v3"
)

func TestDevConfigIgnoresTestGoFiles(t *testing.T) {
	configData, err := os.ReadFile("build_assets/config.yml")
	if err != nil {
		t.Fatalf("failed to read config template: %v", err)
	}

	type devConfig struct {
		Config engine.Config `yaml:"dev_mode"`
	}

	var cfg devConfig
	if err := yaml.Unmarshal(configData, &cfg); err != nil {
		t.Fatalf("failed to parse config: %v", err)
	}

	found := false
	for _, pattern := range cfg.Config.Ignore.File {
		if pattern == "*_test.go" {
			found = true
			break
		}
	}
	if !found {
		t.Error("config template should contain '*_test.go' in the ignore file list")
	}
}

func TestDevConfigParsesCorrectly(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yml")

	configContent := []byte(`version: '3'
dev_mode:
  root_path: .
  log_level: warn
  debounce: 1000
  ignore:
    dir:
      - .git
      - node_modules
    file:
      - .DS_Store
      - "*_test.go"
    watched_extension:
      - "*.go"
  executes:
    - cmd: echo "build"
      type: blocking
`)

	if err := os.WriteFile(configPath, configContent, 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("failed to read config: %v", err)
	}

	type devConfig struct {
		Config engine.Config `yaml:"dev_mode"`
	}

	var cfg devConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		t.Fatalf("failed to parse config: %v", err)
	}

	if len(cfg.Config.Ignore.File) < 2 {
		t.Errorf("expected at least 2 file ignore patterns, got %d", len(cfg.Config.Ignore.File))
	}

	if len(cfg.Config.ExecStruct) < 1 {
		t.Errorf("expected at least 1 execute command, got %d", len(cfg.Config.ExecStruct))
	}

	if !strings.Contains(string(data), "*_test.go") {
		t.Error("config should contain *_test.go pattern")
	}
}

func TestDefaultBuildAssetsContainTestGoIgnore(t *testing.T) {
	configData, err := os.ReadFile("build_assets/config.yml")
	if err != nil {
		t.Fatalf("failed to read build assets config: %v", err)
	}

	content := string(configData)

	if !strings.Contains(content, "*_test.go") {
		t.Error("default build assets config.yml should include *_test.go in ignore file list")
	}

	lines := strings.Split(content, "\n")
	inFileSection := false
	found := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "file:" {
			inFileSection = true
			continue
		}
		if inFileSection {
			if strings.HasPrefix(trimmed, "- ") {
				val := strings.TrimSpace(strings.TrimPrefix(trimmed, "- "))
				if val == "\"*_test.go\"" || val == "*_test.go" {
					found = true
					break
				}
			} else if trimmed != "" && !strings.HasPrefix(trimmed, "-") {
				inFileSection = false
			}
		}
	}

	if !found {
		t.Error("*_test.go should be listed as a file ignore pattern in the config")
	}
}
