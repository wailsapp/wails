package migrate

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const fixtureWailsJSON = `{
  "name": "myv2app",
  "outputfilename": "myv2app",
  "frontend:install": "pnpm install",
  "frontend:build": "pnpm run build",
  "author": { "name": "Tester" },
  "info": {
    "companyName": "Wails",
    "productName": "My V2 App",
    "productVersion": "1.2.3"
  }
}`

const fixtureGoMod = `module myv2app

go 1.23

require github.com/wailsapp/wails/v2 v2.10.2
`

const fixtureMainGo = `package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()

	err := wails.Run(&options.App{
		Title:     "My V2 App",
		Width:     1024,
		Height:    768,
		Frameless: true,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		OnShutdown:       app.shutdown,
		OnBeforeClose:    app.beforeClose,
		WindowStartState: options.Maximised,
		Bind: []interface{}{
			app,
		},
		Windows: &windows.Options{
			Theme: windows.Dark,
		},
		Mac: &mac.Options{
			TitleBar: mac.TitleBarHiddenInset(),
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
`

const fixtureAppGo = `package main

import (
	"context"
	"fmt"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) shutdown(ctx context.Context) {}

func (a *App) beforeClose(ctx context.Context) bool { return false }

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	runtime.WindowSetTitle(a.ctx, "Greeted "+name)
	return fmt.Sprintf("Hello %s!", name)
}

// Add adds two numbers
func (a *App) Add(x int, y int) int { return x + y }
`

func writeFixture(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	files := map[string]string{
		"wails.json":            fixtureWailsJSON,
		"go.mod":                fixtureGoMod,
		"main.go":               fixtureMainGo,
		"app.go":                fixtureAppGo,
		"frontend/package.json": `{"name":"frontend","devDependencies":{"vite":"^3.0.0"}}`,
		"frontend/src/main.js": `import {Greet} from '../wailsjs/go/main/App';
import {EventsOn} from '../wailsjs/runtime/runtime';
EventsOn('x', () => Greet('y'));
`,
		"frontend/index.html":                 `<html></html>`,
		"frontend/dist/.gitkeep":              "",
		"frontend/wailsjs/runtime/runtime.js": "// old",
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

func parseFixture(t *testing.T) *V2Project {
	t.Helper()
	proj, err := ParseV2Project(writeFixture(t))
	if err != nil {
		t.Fatalf("ParseV2Project: %v", err)
	}
	return proj
}

func TestLoadV2ConfigDefaults(t *testing.T) {
	dir := t.TempDir()
	err := os.WriteFile(filepath.Join(dir, "wails.json"), []byte(`{"name":"demo"}`), 0o644)
	if err != nil {
		t.Fatal(err)
	}
	cfg, err := LoadV2Config(dir)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.FrontendDir != "frontend" || cfg.BuildDir != "build" || cfg.OutputFilename != "demo" {
		t.Errorf("defaults not applied: %+v", cfg)
	}
	if cfg.Info.ProductName != "demo" || cfg.Info.ProductVersion != "1.0.0" {
		t.Errorf("info defaults not applied: %+v", cfg.Info)
	}
	if cfg.PackageManager() != "npm" {
		t.Errorf("expected npm default, got %s", cfg.PackageManager())
	}
}

func TestParseV2Project(t *testing.T) {
	proj := parseFixture(t)

	if proj.ModulePath != "myv2app" {
		t.Errorf("module path: got %q", proj.ModulePath)
	}
	if proj.Config.PackageManager() != "pnpm" {
		t.Errorf("package manager: got %q", proj.Config.PackageManager())
	}
	if proj.Main == nil || !strings.HasSuffix(proj.Main.Path, "main.go") {
		t.Fatalf("main not found: %+v", proj.Main)
	}
	if proj.Main.AppLit == nil {
		t.Fatal("options.App literal not found")
	}
	if proj.Main.ErrIdent != "err" {
		t.Errorf("err ident: got %q", proj.Main.ErrIdent)
	}

	if len(proj.BoundTypes) != 1 {
		t.Fatalf("bound types: got %d", len(proj.BoundTypes))
	}
	bt := proj.BoundTypes[0]
	if bt.Name != "App" || bt.Expr != "app" || bt.PkgName != "main" || bt.PkgPath != "main" {
		t.Errorf("bound type: %+v", bt)
	}
	if len(bt.Methods) != 2 || bt.Methods[0].Name != "Greet" || bt.Methods[1].Name != "Add" {
		t.Fatalf("methods: %+v", bt.Methods)
	}
	greet := bt.Methods[0]
	if len(greet.Params) != 1 || greet.Params[0].Name != "name" || greet.Params[0].TSType != "string" {
		t.Errorf("greet params: %+v", greet.Params)
	}
	if len(greet.Results) != 1 || greet.Results[0].TSType != "string" {
		t.Errorf("greet results: %+v", greet.Results)
	}
}

func TestMapOptions(t *testing.T) {
	proj := parseFixture(t)
	opts := MapOptions(proj)

	find := func(fields []GenField, name string) string {
		for _, f := range fields {
			if f.Name == name {
				return f.Expr
			}
		}
		return ""
	}

	if got := find(opts.App, "Name"); got != `"My V2 App"` {
		t.Errorf("Name: got %s", got)
	}
	if got := find(opts.Win, "Width"); got != "1024" {
		t.Errorf("Width: got %s", got)
	}
	if got := find(opts.Win, "Frameless"); got != "true" {
		t.Errorf("Frameless: got %s", got)
	}
	if got := find(opts.Win, "StartState"); got != "application.WindowStateMaximised" {
		t.Errorf("StartState: got %s", got)
	}
	if got := find(opts.Win, "BackgroundColour"); got != "application.NewRGBA(27, 38, 54, 1)" {
		t.Errorf("BackgroundColour: got %s", got)
	}
	if got := find(opts.Win, "URL"); got != `"/"` {
		t.Errorf("URL: got %s", got)
	}
	if got := find(opts.WinWin, "Theme"); got != "application.Dark" {
		t.Errorf("Windows theme: got %s", got)
	}
	if got := find(opts.WinMac, "TitleBar"); got != "application.MacTitleBarHiddenInset" {
		t.Errorf("Mac titlebar: got %s", got)
	}
	if !strings.Contains(find(opts.App, "Assets"), "application.AssetFileServerFS(assets)") {
		t.Errorf("Assets: got %s", find(opts.App, "Assets"))
	}
	if opts.OnStartup != "app.startup" || opts.OnShutdown != "app.shutdown" {
		t.Errorf("lifecycle: %+v", opts)
	}
	if opts.OnBeforeClose != "app.beforeClose" {
		t.Errorf("OnBeforeClose: got %s", opts.OnBeforeClose)
	}
	if len(opts.Services) != 1 || opts.Services[0] != "app" {
		t.Errorf("services: %+v", opts.Services)
	}
	if !opts.NeedsLifecycleService() {
		t.Error("expected lifecycle service")
	}
}

func TestGenerateMain(t *testing.T) {
	proj := parseFixture(t)
	opts := MapOptions(proj)
	out, err := GenerateMain(proj, opts)
	if err != nil {
		t.Fatalf("GenerateMain: %v", err)
	}
	src := string(out)

	for _, want := range []string{
		"application.New(application.Options{",
		`Name: "My V2 App"`,
		"application.NewService(app)",
		"OnShutdown: func() {",
		"(app.shutdown)(context.Background())",
		"wailsApp.Event.OnApplicationEvent(events.Common.ApplicationStarted, func(event *application.ApplicationEvent) {",
		"(app.startup)(context.Background())",
		"ShouldQuit: func() bool {",
		"!(app.beforeClose)(context.Background())",
		".Window.NewWithOptions(application.WebviewWindowOptions{",
		"//go:embed all:frontend/dist",
		"// Create an instance of the app structure",
		"err := wailsApp.Run()",
		`"github.com/wailsapp/wails/v3/pkg/events"`,
	} {
		if !strings.Contains(src, want) {
			t.Errorf("generated main.go missing %q\n---\n%s", want, src)
		}
	}
	for _, banned := range []string{
		"wailsapp/wails/v2",
		"wails.Run",
		"options.App",
	} {
		if strings.Contains(src, banned) {
			t.Errorf("generated main.go still contains %q\n---\n%s", banned, src)
		}
	}

	// Must parse as valid Go.
	fset := token.NewFileSet()
	if _, err := parser.ParseFile(fset, "main.go", out, parser.SkipObjectResolution); err != nil {
		t.Fatalf("generated main.go does not parse: %v\n---\n%s", err, src)
	}
}

func TestGenerateMainBareCall(t *testing.T) {
	dir := writeFixture(t)
	main := `package main

import (
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
)

func main() {
	wails.Run(&options.App{Title: "bare"})
}
`
	if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte(main), 0o644); err != nil {
		t.Fatal(err)
	}
	proj, err := ParseV2Project(dir)
	if err != nil {
		t.Fatal(err)
	}
	out, err := GenerateMain(proj, MapOptions(proj))
	if err != nil {
		t.Fatal(err)
	}
	src := string(out)
	if !strings.Contains(src, "app.Run()") {
		t.Errorf("expected bare Run call:\n%s", src)
	}
	if strings.Contains(src, "err :=") {
		t.Errorf("unexpected error assignment:\n%s", src)
	}
}

func TestTransformGoMod(t *testing.T) {
	proj := parseFixture(t)
	out, err := TransformGoMod(proj, "v3.0.0-alpha2.114")
	if err != nil {
		t.Fatal(err)
	}
	src := string(out)
	if strings.Contains(src, "wails/v2") {
		t.Errorf("v2 require not removed:\n%s", src)
	}
	if !strings.Contains(src, "github.com/wailsapp/wails/v3 v3.0.0-alpha2.114") {
		t.Errorf("v3 require missing:\n%s", src)
	}
	if !strings.Contains(src, "go 1.24") {
		t.Errorf("go directive not raised:\n%s", src)
	}
}

func TestAdvisor(t *testing.T) {
	proj := parseFixture(t)
	if !proj.UsesV2Runtime {
		t.Fatal("expected UsesV2Runtime")
	}
	md := proj.Report.Markdown()
	for _, want := range []string{
		"## Port these to the v3 API",
		"`app.go:", // Go call site location
		"`runtime.WindowSetTitle`",
		"window.SetTitle(title)",
		"`src/main.js` should not appear",
	} {
		if want == "`src/main.js` should not appear" {
			continue
		}
		if !strings.Contains(md, want) {
			t.Errorf("report missing %q\n---\n%s", want, md)
		}
	}
	// Frontend wailsjs imports are listed too.
	for _, want := range []string{
		"main.js:1",
		"wailsjs/go/main/App",
		"wails3 generate bindings",
		"main.js:2",
		"@wailsio/runtime",
	} {
		if !strings.Contains(md, want) {
			t.Errorf("report missing frontend advice %q\n---\n%s", want, md)
		}
	}
}

func TestMigrateFrontend(t *testing.T) {
	proj := parseFixture(t)
	outDir := t.TempDir()
	if err := MigrateFrontend(proj, outDir); err != nil {
		t.Fatal(err)
	}

	// The generated v2 wailsjs directory must not be carried over or
	// replaced by a lookalike.
	if _, err := os.Stat(filepath.Join(outDir, "frontend", "wailsjs")); err == nil {
		t.Error("frontend/wailsjs should not exist in the migrated project")
	}

	// The user's own sources are copied untouched.
	mainJS, err := os.ReadFile(filepath.Join(outDir, "frontend", "src", "main.js"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(mainJS), "../wailsjs/go/main/App") {
		t.Errorf("user source was modified:\n%s", mainJS)
	}

	pkgJSON, err := os.ReadFile(filepath.Join(outDir, "frontend", "package.json"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(pkgJSON), `"@wailsio/runtime"`) {
		t.Errorf("package.json missing runtime dep:\n%s", pkgJSON)
	}
}
