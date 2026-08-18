package commands

import (
	"go/format"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/wailsapp/wails/v3/internal/flags"
)

// copyFixture copies a testdata fixture project into a temp dir so the
// migrator can write into it.
func copyFixture(t *testing.T, fixture string) string {
	t.Helper()
	src := filepath.Join("testdata", "migrate", fixture)
	dest := t.TempDir()
	err := filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dest, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(target, data, 0o644)
	})
	if err != nil {
		t.Fatalf("could not copy fixture: %v", err)
	}
	return dest
}

func runMigrate(t *testing.T, dir string) {
	t.Helper()
	err := Migrate(&flags.Migrate{ProjectDir: dir, Quiet: true})
	if err != nil {
		t.Fatalf("Migrate failed: %v", err)
	}
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("could not read %s: %v", path, err)
	}
	return string(data)
}

func assertContains(t *testing.T, content, want, context string) {
	t.Helper()
	if !strings.Contains(content, want) {
		t.Errorf("%s: expected to find %q", context, want)
	}
}

func TestMigrate_RejectsNonV2Project(t *testing.T) {
	// Empty directory: no wails.json.
	err := Migrate(&flags.Migrate{ProjectDir: t.TempDir(), Quiet: true})
	if err == nil || !strings.Contains(err.Error(), "does not look like a Wails v2 project") {
		t.Fatalf("expected a polite refusal, got: %v", err)
	}

	// wails.json present but go.mod requires v3, not v2.
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "wails.json"), []byte(`{"name":"x"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	gomod := "module x\n\ngo 1.21\n\nrequire github.com/wailsapp/wails/v3 v3.0.0\n"
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte(gomod), 0o644); err != nil {
		t.Fatal(err)
	}
	err = Migrate(&flags.Migrate{ProjectDir: dir, Quiet: true})
	if err == nil || !strings.Contains(err.Error(), "does not look like a Wails v2 project") {
		t.Fatalf("expected a polite refusal, got: %v", err)
	}
}

func TestMigrate_ConfigMapping(t *testing.T) {
	dir := copyFixture(t, "basic")
	runMigrate(t, dir)

	taskfile := readFile(t, filepath.Join(dir, "Taskfile.yml"))
	assertContains(t, taskfile, `APP_NAME: "myapp-bin"`, "Taskfile.yml")
	assertContains(t, taskfile, `default "pnpm"`, "Taskfile.yml package manager")

	configYML := readFile(t, filepath.Join(dir, "build", "config.yml"))
	assertContains(t, configYML, `companyName: "Test Company"`, "config.yml")
	assertContains(t, configYML, `productName: "My Test App"`, "config.yml")
	assertContains(t, configYML, `version: "2.3.4"`, "config.yml")
	assertContains(t, configYML, `copyright: "Copyright 2024 Test Company"`, "config.yml")

	// The v3 build system must be in place.
	for _, file := range []string{
		filepath.Join("build", "Taskfile.yml"),
		filepath.Join("build", "darwin", "Taskfile.yml"),
		filepath.Join("build", "windows", "Taskfile.yml"),
		filepath.Join("build", "linux", "Taskfile.yml"),
	} {
		if _, err := os.Stat(filepath.Join(dir, file)); err != nil {
			t.Errorf("expected %s to be generated: %v", file, err)
		}
	}
}

func TestMigrate_PreservesExistingFiles(t *testing.T) {
	dir := copyFixture(t, "basic")
	runMigrate(t, dir)

	// The v2 build assets must be untouched.
	if got := readFile(t, filepath.Join(dir, "build", "appicon.png")); got != "fake v2 icon\n" {
		t.Errorf("build/appicon.png was modified: %q", got)
	}
	if got := readFile(t, filepath.Join(dir, "build", "windows", "icon.ico")); got != "fake v2 ico\n" {
		t.Errorf("build/windows/icon.ico was modified: %q", got)
	}
	// main.go must be untouched.
	mainGo := readFile(t, filepath.Join(dir, "main.go"))
	assertContains(t, mainGo, "wails.Run(&options.App{", "main.go must remain a v2 bootstrap")

	// Running the migration a second time must be safe and must not
	// overwrite the previously generated Taskfile.
	taskfilePath := filepath.Join(dir, "Taskfile.yml")
	if err := os.WriteFile(taskfilePath, []byte("# user edited\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	runMigrate(t, dir)
	if got := readFile(t, taskfilePath); got != "# user edited\n" {
		t.Errorf("second run overwrote the user's Taskfile.yml: %q", got)
	}
}

func TestMigrate_GeneratedMain(t *testing.T) {
	dir := copyFixture(t, "basic")
	runMigrate(t, dir)

	path := filepath.Join(dir, "main_v3.go.example")
	content := readFile(t, path)

	// It must parse as Go code.
	fset := token.NewFileSet()
	if _, err := parser.ParseFile(fset, path, content, parser.ParseComments); err != nil {
		t.Fatalf("generated code does not parse: %v", err)
	}
	// It must be gofmt-clean.
	formatted, err := format.Source([]byte(content))
	if err != nil {
		t.Fatalf("generated code does not format: %v", err)
	}
	if string(formatted) != content {
		t.Errorf("generated code is not gofmt-clean")
	}

	// Bind entries became services, with the constructor carried over.
	assertContains(t, content, "app := NewApp()", "generated main")
	assertContains(t, content, "application.NewService(app)", "generated main")
	assertContains(t, content, "application.NewService(&GreetService{})", "generated main")

	// Options made it across. Collapse gofmt's key alignment so the
	// assertions are independent of literal field padding.
	normalised := strings.Join(strings.Fields(content), " ")
	assertContains(t, normalised, `Name: "My Test App"`, "generated main")
	assertContains(t, normalised, `Title: "My Test App"`, "generated main")
	assertContains(t, normalised, "Width: 1024", "generated main")
	assertContains(t, normalised, "Height: 768", "generated main")
	assertContains(t, normalised, "MinWidth: 400", "generated main")
	assertContains(t, normalised, "MinHeight: 300", "generated main")
	assertContains(t, normalised, "Frameless: true", "generated main")
	assertContains(t, normalised, "application.NewRGBA(27, 38, 54, 1)", "generated main")
	assertContains(t, normalised, `UniqueID: "com.test.myapp"`, "generated main")

	// The embed directive was carried over.
	assertContains(t, content, "//go:embed all:frontend/dist", "generated main")
	assertContains(t, content, "application.AssetFileServerFS(assets)", "generated main")

	// Lifecycle hooks surface as TODOs.
	assertContains(t, content, "TODO(migrate): v2 OnStartup was `app.startup`", "generated main")
	assertContains(t, content, "TODO(migrate): v2 OnShutdown was `app.shutdown`", "generated main")
}

func TestMigrate_Report(t *testing.T) {
	dir := copyFixture(t, "basic")
	runMigrate(t, dir)

	report := readFile(t, filepath.Join(dir, "MIGRATION_REPORT.md"))

	// Unmapped config keys are called out with guidance.
	assertContains(t, report, "`wailsjsdir`", "report")
	assertContains(t, report, "`debounceMS`", "report")
	assertContains(t, report, "`frontend:dev:serverUrl`", "report")

	// v2 runtime usage is located with file:line references.
	assertContains(t, report, "`runtime.EventsEmit`", "report")
	assertContains(t, report, "`runtime.WindowSetTitle`", "report")
	assertContains(t, report, "app.go:", "report")

	// The runtime mapping table and the docs link are present.
	assertContains(t, report, "app.Event.Emit(name, data)", "report")
	assertContains(t, report, "window.SetTitle(title)", "report")
	assertContains(t, report, "https://v3.wails.io/migration/v2-to-v3/", "report")

	// Preserved files are listed.
	assertContains(t, report, "`build/appicon.png`", "report")
}

func TestDetectPackageManager(t *testing.T) {
	tests := map[string]string{
		"npm install":          "npm",
		"pnpm install":         "pnpm",
		"yarn":                 "yarn",
		"bun install":          "bun",
		"npm ci":               "npm",
		"":                     "npm",
		"./scripts/install.sh": "npm",
	}
	for command, want := range tests {
		if got := detectPackageManager(command); got != want {
			t.Errorf("detectPackageManager(%q) = %q, want %q", command, got, want)
		}
	}
}
