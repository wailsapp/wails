package commands

import (
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/wailsapp/wails/v3/internal/templates"
	"github.com/wailsapp/wails/v3/internal/wake/manifest"
)

func writeMigrationConfigFixture(t *testing.T, root, path string, data []byte) {
	t.Helper()
	path = filepath.Join(root, filepath.FromSlash(path))
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatal(err)
	}
}

func TestLegacyMigrationDevExitPolicy(t *testing.T) {
	for _, test := range []struct {
		name, command, role, policy string
		complete                    bool
	}{
		{"stock shutdown", "wails3 task run", "primary", "shutdown", true},
		{"implicit shutdown", "wails3 task run", "primary", "", true},
		{"ignore", "wails3 task run", "primary", "ignore", false},
		{"restart", "wails3 task run", "primary", "restart", false},
		{"unknown", "wails3 task run", "primary", "surprise", false},
		{"background", "wails3 task common:dev:frontend", "background", "shutdown", false},
		{"blocking", "wails3 build DEV=true", "blocking", "shutdown", false},
		{"custom shutdown command", "wails3 task run; custom-script", "primary", "shutdown", false},
	} {
		t.Run(test.name, func(t *testing.T) {
			root := t.TempDir()
			raw := []byte("version: '3'\ndev_mode:\n  executes:\n    - cmd: " + test.command + "\n      type: " + test.role + "\n      exit_policy: " + test.policy + "\n")
			writeMigrationConfigFixture(t, root, "build/config.yml", raw)
			report := MigrationReport{Complete: true, Sources: map[string]string{}}
			doc := manifest.NewDocument(manifest.Project{Name: "app"})
			if err := applyLegacyBuildConfiguration(root, &report, &doc); err != nil {
				t.Fatal(err)
			}
			if report.Complete != test.complete {
				t.Fatalf("complete=%v diagnostics=%+v", report.Complete, report.Diagnostics)
			}
			if !test.complete && (len(report.Diagnostics) == 0 || report.Diagnostics[0].Code != "unsupported-dev-exit-policy") {
				t.Fatalf("unexpected diagnostics: %+v", report.Diagnostics)
			}
			after, err := os.ReadFile(filepath.Join(root, "build/config.yml"))
			if err != nil || string(after) != string(raw) {
				t.Fatalf("source changed: %v", err)
			}
		})
	}
}

func TestLegacyMigrationDevCommandsPreserveSuffixSemantics(t *testing.T) {
	for _, command := range []string{
		"wails3 task build EXTRA_TAGS=custom",
		"wails3 task build; custom-script",
		"wails3 task run --custom-argument",
		"wails3 task common:dev:frontend-other",
		"wails3 task common:install:frontend:deps && custom-script",
	} {
		t.Run(command, func(t *testing.T) {
			root := t.TempDir()
			writeMigrationConfigFixture(t, root, "build/config.yml", []byte("dev_mode:\n  executes:\n    - cmd: "+command+"\n      type: blocking\n"))
			report := MigrationReport{Complete: true, Sources: map[string]string{}}
			doc := manifest.NewDocument(manifest.Project{Name: "app"})
			if err := applyLegacyBuildConfiguration(root, &report, &doc); err != nil {
				t.Fatal(err)
			}
			if report.Complete || len(report.Diagnostics) != 1 || report.Diagnostics[0].Code != "unsupported-dev-command" {
				t.Fatalf("custom command was silently dropped: %+v", report)
			}
		})
	}
}

func TestLegacyMigrationUnbuiltStockFrontend(t *testing.T) {
	for _, template := range templates.GetDefaultTemplates() {
		t.Run(template.Name, func(t *testing.T) {
			root := t.TempDir()
			var configName string
			for _, name := range []string{"vite.config.js", "vite.config.ts"} {
				data, err := fs.ReadFile(template.FS, template.Name+"/frontend/"+name)
				if err != nil {
					continue
				}
				configName = name
				writeMigrationConfigFixture(t, root, "frontend/"+name, data)
			}
			if configName == "" {
				t.Skip("template has no Vite frontend")
			}
			data, err := fs.ReadFile(template.FS, template.Name+"/frontend/package.json.tmpl")
			if err != nil {
				t.Fatal(err)
			}
			writeMigrationConfigFixture(t, root, "frontend/package.json", data)
			report := MigrationReport{Complete: true, Sources: map[string]string{}}
			doc := manifest.NewDocument(manifest.Project{Name: "app"})
			if err := applyLegacyFrontendConfiguration(root, "npm", &report, &doc); err != nil {
				t.Fatal(err)
			}
			if !report.Complete {
				t.Fatalf("stock unbuilt frontend rejected: %+v", report.Diagnostics)
			}
			if doc.Frontend.OutputDirectory != "frontend/dist" || report.Sources["frontend/"+configName] == "" {
				t.Fatalf("missing output or source evidence: %+v", report)
			}
			if _, err := os.Lstat(filepath.Join(root, "frontend/dist")); !os.IsNotExist(err) {
				t.Fatalf("migration created output: %v", err)
			}
		})
	}
}

func TestLegacyMigrationMissingOutputDoesNotGuessCustomConfiguration(t *testing.T) {
	var stockConfig, stockPackage []byte
	for _, template := range templates.GetDefaultTemplates() {
		if template.Name == "vanilla-js" {
			stockConfig, _ = fs.ReadFile(template.FS, "vanilla-js/frontend/vite.config.js")
			stockPackage, _ = fs.ReadFile(template.FS, "vanilla-js/frontend/package.json.tmpl")
		}
	}
	if len(stockConfig) == 0 || len(stockPackage) == 0 {
		t.Fatal("missing stock fixture")
	}
	for _, name := range []string{"custom output", "custom script", "prebuild", "postbuild", "competing config", "symlink config", "no config", "output file", "dangling output"} {
		t.Run(name, func(t *testing.T) {
			root := t.TempDir()
			config := append([]byte(nil), stockConfig...)
			var metadata map[string]any
			if err := json.Unmarshal(stockPackage, &metadata); err != nil {
				t.Fatal(err)
			}
			scripts := metadata["scripts"].(map[string]any)
			switch name {
			case "custom output":
				config = []byte(strings.Replace(string(config), "defineConfig({", "defineConfig({ build: { outDir: 'release' },", 1))
			case "custom script":
				scripts["build"] = "vite build --outDir release"
			case "prebuild", "postbuild":
				scripts[name] = "node custom.js"
			}
			data, err := json.Marshal(metadata)
			if err != nil {
				t.Fatal(err)
			}
			writeMigrationConfigFixture(t, root, "frontend/package.json", data)
			if name != "no config" {
				writeMigrationConfigFixture(t, root, "frontend/vite.config.js", config)
			}
			if name == "competing config" {
				writeMigrationConfigFixture(t, root, "frontend/vite.config.ts", stockConfig)
			}
			if name == "symlink config" {
				path := filepath.Join(root, "frontend/vite.config.js")
				if err := os.Rename(path, filepath.Join(root, "frontend/linked.js")); err != nil {
					t.Fatal(err)
				}
				if err := os.Symlink("linked.js", path); err != nil {
					t.Skipf("symlinks unavailable: %v", err)
				}
			}
			if name == "output file" {
				writeMigrationConfigFixture(t, root, "frontend/dist", []byte("not a directory"))
			}
			if name == "dangling output" {
				if err := os.Symlink("missing", filepath.Join(root, "frontend/dist")); err != nil {
					t.Skipf("symlinks unavailable: %v", err)
				}
			}
			report := MigrationReport{Complete: true, Sources: map[string]string{}}
			doc := manifest.NewDocument(manifest.Project{Name: "app"})
			if err := applyLegacyFrontendConfiguration(root, "npm", &report, &doc); err != nil {
				t.Fatal(err)
			}
			if report.Complete {
				t.Fatal("custom output accepted without evidence")
			}
			found := false
			for _, diagnostic := range report.Diagnostics {
				found = found || diagnostic.Code == "frontend-output"
			}
			if !found {
				t.Fatalf("missing output diagnostic: %+v", report.Diagnostics)
			}
		})
	}
}

func TestLegacyMigrationDevRoles(t *testing.T) {
	for _, test := range []struct {
		command string
		roles   []string
	}{
		{"go mod tidy", []string{"blocking"}},
		{"wails3 build DEV=true", []string{"blocking"}},
		{"wails3 task build", []string{"blocking"}},
		{"wails3 task common:install:frontend:deps", []string{"once", "blocking"}},
		{"wails3 task common:dev:frontend", []string{"background"}},
		{"wails3 task run", []string{"primary"}},
	} {
		for _, role := range []string{"", "blocking", "once", "background", "primary", "custom"} {
			t.Run(test.command+"/"+role, func(t *testing.T) {
				root := t.TempDir()
				raw := []byte("dev_mode:\n  executes:\n    - cmd: " + test.command + "\n      type: " + role + "\n")
				writeMigrationConfigFixture(t, root, "build/config.yml", raw)
				report := MigrationReport{Complete: true, Sources: map[string]string{}}
				doc := manifest.NewDocument(manifest.Project{Name: "app"})
				if err := applyLegacyBuildConfiguration(root, &report, &doc); err != nil {
					t.Fatal(err)
				}
				supported := false
				for _, expected := range test.roles {
					supported = supported || role == expected
				}
				if report.Complete != supported {
					t.Fatalf("command %q role %q: complete=%v diagnostics=%+v", test.command, role, report.Complete, report.Diagnostics)
				}
				if !supported && (len(report.Diagnostics) != 1 || report.Diagnostics[0].Code != "unsupported-dev-role") {
					t.Fatalf("unexpected diagnostics: %+v", report.Diagnostics)
				}
				after, err := os.ReadFile(filepath.Join(root, "build/config.yml"))
				if err != nil || string(after) != string(raw) {
					t.Fatalf("source changed: %v", err)
				}
			})
		}
	}
}
