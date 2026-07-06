package commands

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/wailsapp/wails/v3/internal/flags"
)

func writeV2Fixture(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	files := map[string]string{
		"wails.json": `{
  "name": "demoapp",
  "outputfilename": "demoapp",
  "frontend:install": "npm install",
  "frontend:build": "npm run build",
  "info": { "companyName": "Wails", "productName": "Demo App", "productVersion": "2.0.0" }
}`,
		"go.mod": "module demoapp\n\ngo 1.23\n\nrequire github.com/wailsapp/wails/v2 v2.10.2\n",
		"main.go": `package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app := NewApp()

	err := wails.Run(&options.App{
		Title:  "Demo App",
		Width:  800,
		Height: 600,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup: app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
`,
		"app.go": `package main

import "context"

type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) Greet(name string) string {
	return "Hello " + name
}
`,
		"frontend/package.json":  `{"name":"frontend","devDependencies":{"vite":"^3.0.0"}}`,
		"frontend/index.html":    `<html></html>`,
		"frontend/dist/.gitkeep": "",
		"build/appicon.png":      "not-really-a-png",
	}
	for name, content := range files {
		path := filepath.Join(dir, name)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

func TestMigrateEndToEnd(t *testing.T) {
	v2Dir := writeV2Fixture(t)
	outDir := filepath.Join(t.TempDir(), "v3out")

	err := Migrate(&flags.Migrate{
		V2Dir:         v2Dir,
		OutputDir:     outDir,
		Quiet:         true,
		SkipGoModTidy: true,
	})
	if err != nil {
		t.Fatalf("Migrate: %v", err)
	}

	mustExist := []string{
		"main.go",
		"app.go",
		"go.mod",
		"Taskfile.yml",
		"MIGRATION.md",
		".gitignore",
		"build/config.yml",
		"build/Taskfile.yml",
		"build/darwin/Taskfile.yml",
		"build/appicon.png",
		"frontend/package.json",
		"frontend/wailsjs/runtime/runtime.js",
		"frontend/wailsjs/go/main/App.js",
		"frontend/dist/.gitkeep",
		"v2compat/runtime/window.go",
		"v2compat/runtime/lifecycle.go",
	}
	for _, rel := range mustExist {
		if _, err := os.Stat(filepath.Join(outDir, rel)); err != nil {
			t.Errorf("expected %s to exist: %v", rel, err)
		}
	}

	mustNotExist := []string{
		"wails.json",
		"go.sum",
	}
	for _, rel := range mustNotExist {
		if _, err := os.Stat(filepath.Join(outDir, rel)); err == nil {
			t.Errorf("expected %s to be absent", rel)
		}
	}

	mainSrc, err := os.ReadFile(filepath.Join(outDir, "main.go"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(mainSrc), "application.New(application.Options{") {
		t.Errorf("main.go not migrated:\n%s", mainSrc)
	}
	if strings.Contains(string(mainSrc), "wails/v2") {
		t.Errorf("main.go still references v2:\n%s", mainSrc)
	}

	// The kept v2 appicon must win over the template icon.
	icon, err := os.ReadFile(filepath.Join(outDir, "build", "appicon.png"))
	if err != nil {
		t.Fatal(err)
	}
	if string(icon) != "not-really-a-png" {
		t.Error("v2 appicon.png was not preserved")
	}

	configYML, err := os.ReadFile(filepath.Join(outDir, "build", "config.yml"))
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{`productName: "Demo App"`, `version: "2.0.0"`, `companyName: "Wails"`, `productIdentifier: "com.wails.demoapp"`} {
		if !strings.Contains(string(configYML), want) {
			t.Errorf("config.yml missing %s:\n%s", want, configYML)
		}
	}

	goMod, err := os.ReadFile(filepath.Join(outDir, "go.mod"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(goMod), "github.com/wailsapp/wails/v3") || strings.Contains(string(goMod), "wails/v2") {
		t.Errorf("go.mod not transformed:\n%s", goMod)
	}
}

func TestMigrateRefusesNonEmptyOutput(t *testing.T) {
	v2Dir := writeV2Fixture(t)
	outDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(outDir, "existing.txt"), []byte("hi"), 0o644); err != nil {
		t.Fatal(err)
	}

	err := Migrate(&flags.Migrate{
		V2Dir:         v2Dir,
		OutputDir:     outDir,
		Quiet:         true,
		SkipGoModTidy: true,
	})
	if err == nil || !strings.Contains(err.Error(), "not empty") {
		t.Fatalf("expected non-empty error, got %v", err)
	}
}

func TestMigrateRefusesNestedOutput(t *testing.T) {
	v2Dir := writeV2Fixture(t)

	err := Migrate(&flags.Migrate{
		V2Dir:         v2Dir,
		OutputDir:     filepath.Join(v2Dir, "out"),
		Quiet:         true,
		SkipGoModTidy: true,
	})
	if err == nil || !strings.Contains(err.Error(), "must not be inside") {
		t.Fatalf("expected nested-output error, got %v", err)
	}
}

func TestMigrateRejectsNonV2Project(t *testing.T) {
	dir := t.TempDir()
	files := map[string]string{
		"wails.json": `{"name":"demo"}`,
		"go.mod":     "module demo\n\ngo 1.23\n",
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	err := Migrate(&flags.Migrate{
		V2Dir:         dir,
		OutputDir:     filepath.Join(t.TempDir(), "out"),
		Quiet:         true,
		SkipGoModTidy: true,
	})
	if err == nil || !strings.Contains(err.Error(), "wailsapp/wails/v2") {
		t.Fatalf("expected v2-detection error, got %v", err)
	}
}
