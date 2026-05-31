package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/wailsapp/wails/v2/internal/project"
	"github.com/wailsapp/wails/v2/internal/staticanalysis"
	"github.com/wailsapp/wails/v2/pkg/commands/build"
)

func TestEmbedDiscoveryInTempDir(t *testing.T) {
	projectDir := t.TempDir()

	jsfrontendDir := filepath.Join(projectDir, "jsfrontend")
	jsfrontendDist := filepath.Join(jsfrontendDir, "dist")
	os.MkdirAll(jsfrontendDist, 0o755)
	os.WriteFile(filepath.Join(jsfrontendDist, "index.html"), []byte("<html><body>Hello</body></html>"), 0o644)

	goFile := filepath.Join(projectDir, "main.go")
	os.WriteFile(goFile, []byte("package main\n\nimport \"embed\"\n\n//go:embed all:frontend/dist\nvar assets embed.FS\n\nfunc main() {}\n"), 0o644)
	os.WriteFile(filepath.Join(projectDir, "go.mod"), []byte("module test\n\ngo 1.21\n"), 0o644)

	embedDetails, err := staticanalysis.GetEmbedDetails(projectDir)
	t.Logf("GetEmbedDetails returned %d results, err=%v", len(embedDetails), err)
	for i, d := range embedDetails {
		t.Logf("  [%d] BaseDir=%s EmbedPath=%s FullPath=%s All=%v", i, d.BaseDir, d.EmbedPath, d.GetFullPath(), d.All)
	}
}

func TestSyncFrontendDistToEmbedTarget(t *testing.T) {
	projectDir := t.TempDir()

	jsfrontendDir := filepath.Join(projectDir, "jsfrontend")
	jsfrontendDist := filepath.Join(jsfrontendDir, "dist")
	os.MkdirAll(jsfrontendDist, 0o755)

	indexHTML := filepath.Join(jsfrontendDist, "index.html")
	os.WriteFile(indexHTML, []byte("<html><body>Hello</body></html>"), 0o644)

	assetsDir := filepath.Join(jsfrontendDist, "assets")
	os.MkdirAll(assetsDir, 0o755)
	os.WriteFile(filepath.Join(assetsDir, "app.js"), []byte("console.log('test')"), 0o644)

	goFile := filepath.Join(projectDir, "main.go")
	os.WriteFile(goFile, []byte("package main\n\nimport \"embed\"\n\n//go:embed all:frontend/dist\nvar assets embed.FS\n\nfunc main() {}\n"), 0o644)
	os.WriteFile(filepath.Join(projectDir, "go.mod"), []byte("module test\n\ngo 1.21\n"), 0o644)

	frontendDist := filepath.Join(projectDir, "frontend", "dist")
	os.MkdirAll(frontendDist, 0o755)
	os.WriteFile(filepath.Join(frontendDist, "gitkeep"), []byte(""), 0o644)

	wailsJSON := filepath.Join(projectDir, "wails.json")
	os.WriteFile(wailsJSON, []byte(`{"projectdir": "`+filepath.ToSlash(projectDir)+`", "frontend:dir": "./jsfrontend", "name": "test"}`), 0o644)

	proj, err := project.Load(projectDir)
	if err != nil {
		t.Fatalf("failed to load project: %v", err)
	}

	options := &build.Options{
		ProjectData: proj,
	}

	err = build.CreateEmbedDirectories(projectDir, options)
	if err != nil {
		t.Fatalf("CreateEmbedDirectories failed: %v", err)
	}

	cleanup, err := build.SyncFrontendDistToEmbedTarget(options)
	if err != nil {
		t.Fatalf("SyncFrontendDistToEmbedTarget failed: %v", err)
	}
	defer cleanup()

	copiedIndex := filepath.Join(frontendDist, "index.html")
	data, err := os.ReadFile(copiedIndex)
	if err != nil {
		t.Fatalf("index.html not found in embed target: %v", err)
	}
	if string(data) != "<html><body>Hello</body></html>" {
		t.Fatalf("index.html content mismatch: %s", string(data))
	}

	copiedAppJS := filepath.Join(frontendDist, "assets", "app.js")
	data, err = os.ReadFile(copiedAppJS)
	if err != nil {
		t.Fatalf("assets/app.js not found in embed target: %v", err)
	}
	if string(data) != "console.log('test')" {
		t.Fatalf("app.js content mismatch: %s", string(data))
	}

	gitkeep := filepath.Join(frontendDist, "gitkeep")
	if _, err := os.Stat(gitkeep); !os.IsNotExist(err) {
		t.Fatal("gitkeep should have been removed when frontend/dist was replaced")
	}
}

func TestSyncFrontendDistCleanupRestoresGitkeep(t *testing.T) {
	projectDir := t.TempDir()

	jsfrontendDist := filepath.Join(projectDir, "jsfrontend", "dist")
	os.MkdirAll(jsfrontendDist, 0o755)
	os.WriteFile(filepath.Join(jsfrontendDist, "index.html"), []byte("<html></html>"), 0o644)

	goFile := filepath.Join(projectDir, "main.go")
	os.WriteFile(goFile, []byte("package main\n\nimport \"embed\"\n\n//go:embed all:frontend/dist\nvar assets embed.FS\n\nfunc main() {}\n"), 0o644)
	os.WriteFile(filepath.Join(projectDir, "go.mod"), []byte("module test\n\ngo 1.21\n"), 0o644)

	frontendDist := filepath.Join(projectDir, "frontend", "dist")
	os.MkdirAll(frontendDist, 0o755)
	os.WriteFile(filepath.Join(frontendDist, ".gitkeep"), []byte(""), 0o644)

	wailsJSON := filepath.Join(projectDir, "wails.json")
	os.WriteFile(wailsJSON, []byte(`{"projectdir": "`+filepath.ToSlash(projectDir)+`", "frontend:dir": "./jsfrontend", "name": "test"}`), 0o644)

	proj, err := project.Load(projectDir)
	if err != nil {
		t.Fatalf("failed to load project: %v", err)
	}

	options := &build.Options{ProjectData: proj}

	cleanup, err := build.SyncFrontendDistToEmbedTarget(options)
	if err != nil {
		t.Fatalf("SyncFrontendDistToEmbedTarget failed: %v", err)
	}

	// Verify assets were copied in.
	if _, err := os.Stat(filepath.Join(frontendDist, "index.html")); err != nil {
		t.Fatal("index.html should exist in embed target before cleanup")
	}

	// Run cleanup.
	cleanup()

	// Copied assets must be gone.
	if _, err := os.Stat(filepath.Join(frontendDist, "index.html")); !os.IsNotExist(err) {
		t.Fatal("index.html should have been removed by cleanup")
	}

	// .gitkeep must be restored.
	if _, err := os.Stat(filepath.Join(frontendDist, ".gitkeep")); err != nil {
		t.Fatal(".gitkeep should be restored by cleanup")
	}

	// No other files should exist.
	entries, _ := os.ReadDir(frontendDist)
	if len(entries) != 1 {
		t.Fatalf("frontend/dist should contain only .gitkeep after cleanup, got %d entries", len(entries))
	}
}

func TestSyncFrontendDistNoOpWhenDefault(t *testing.T) {
	projectDir := t.TempDir()

	frontendDir := filepath.Join(projectDir, "frontend")
	frontendDist := filepath.Join(frontendDir, "dist")
	os.MkdirAll(frontendDist, 0o755)
	os.WriteFile(filepath.Join(frontendDist, "index.html"), []byte("<html></html>"), 0o644)

	goFile := filepath.Join(projectDir, "main.go")
	os.WriteFile(goFile, []byte("package main\n\nimport \"embed\"\n\n//go:embed all:frontend/dist\nvar assets embed.FS\n\nfunc main() {}\n"), 0o644)
	os.WriteFile(filepath.Join(projectDir, "go.mod"), []byte("module test\n\ngo 1.21\n"), 0o644)

	wailsJSON := filepath.Join(projectDir, "wails.json")
	os.WriteFile(wailsJSON, []byte(`{"projectdir": "`+filepath.ToSlash(projectDir)+`", "name": "test"}`), 0o644)

	proj, err := project.Load(projectDir)
	if err != nil {
		t.Fatalf("failed to load project: %v", err)
	}

	options := &build.Options{
		ProjectData: proj,
	}

	cleanup, err := build.SyncFrontendDistToEmbedTarget(options)
	if err != nil {
		t.Fatalf("SyncFrontendDistToEmbedTarget should not error for default frontend dir: %v", err)
	}
	cleanup() // noop — must not panic or modify anything

	data, err := os.ReadFile(filepath.Join(frontendDist, "index.html"))
	if err != nil {
		t.Fatalf("original index.html should still exist: %v", err)
	}
	if string(data) != "<html></html>" {
		t.Fatal("original index.html should be unchanged")
	}
}

func TestSyncFrontendDistNoOpWhenNoDistInEmbed(t *testing.T) {
	projectDir := t.TempDir()

	jsfrontendDir := filepath.Join(projectDir, "jsfrontend")
	os.MkdirAll(filepath.Join(jsfrontendDir, "dist"), 0o755)

	goFile := filepath.Join(projectDir, "main.go")
	os.WriteFile(goFile, []byte("package main\n\nfunc main() {}\n"), 0o644)
	os.WriteFile(filepath.Join(projectDir, "go.mod"), []byte("module test\n\ngo 1.21\n"), 0o644)

	wailsJSON := filepath.Join(projectDir, "wails.json")
	os.WriteFile(wailsJSON, []byte(`{"projectdir": "`+filepath.ToSlash(projectDir)+`", "frontend:dir": "./jsfrontend", "name": "test"}`), 0o644)

	proj, err := project.Load(projectDir)
	if err != nil {
		t.Fatalf("failed to load project: %v", err)
	}

	options := &build.Options{
		ProjectData: proj,
	}

	cleanup, err := build.SyncFrontendDistToEmbedTarget(options)
	if err != nil {
		t.Fatalf("should not error when no embed dist dirs: %v", err)
	}
	cleanup() // noop
}
