package migrate

import (
	"bufio"
	"fmt"
	"go/ast"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// The migrator deliberately does not rewrite runtime call sites and does not
// generate compatibility layers: half-migrated code helps nobody. Instead,
// every v2 runtime call and every wailsjs import is recorded here with its
// concrete v3 replacement, so the user gets a precise, project-specific
// checklist in MIGRATION.md and the compiler points at exactly the listed
// locations until they are ported.

// goRuntimeAdvice maps a v2 runtime function name to advice for porting the
// call site to the v3 API. `app` refers to the application instance
// (application.Get() from anywhere).
var goRuntimeAdvice = map[string]string{
	"EventsEmit":       "`app.Event.Emit(name, data...)`",
	"EventsOn":         "`app.Event.On(name, func(e *application.CustomEvent) { ... })` - the callback receives the event object; your payload is `e.Data`. Returns an unsubscribe func.",
	"EventsOnce":       "`app.Event.On(...)` and call the returned unsubscribe func inside the callback (v3 has no Once on the manager)",
	"EventsOnMultiple": "`app.Event.OnMultiple(name, callback, counter)`",
	"EventsOff":        "`app.Event.Off(name)`",
	"EventsOffAll":     "`app.Event.Reset()`",

	"Quit":        "`app.Quit()`",
	"Hide":        "`app.Hide()`",
	"Show":        "`app.Show()`",
	"Environment": "`app.Env.Info()` (returns application.EnvironmentInfo: OS, Arch, Debug)",

	"BrowserOpenURL": "`app.Browser.OpenURL(url)`",

	"ClipboardGetText": "`app.Clipboard.Text()` (returns (string, bool) instead of (string, error))",
	"ClipboardSetText": "`app.Clipboard.SetText(text)` (returns bool instead of error)",

	"ScreenGetAll": "`app.Screen.GetAll()` (the v3 Screen struct differs: Size, Bounds, PhysicalBounds, IsPrimary)",

	"LogPrint":       "`app.Logger.Info(message)`",
	"LogPrintf":      "`app.Logger.Info(fmt.Sprintf(...))`",
	"LogTrace":       "`app.Logger.Debug(message)`",
	"LogTracef":      "`app.Logger.Debug(fmt.Sprintf(...))`",
	"LogDebug":       "`app.Logger.Debug(message)`",
	"LogDebugf":      "`app.Logger.Debug(fmt.Sprintf(...))`",
	"LogInfo":        "`app.Logger.Info(message)`",
	"LogInfof":       "`app.Logger.Info(fmt.Sprintf(...))`",
	"LogWarning":     "`app.Logger.Warn(message)`",
	"LogWarningf":    "`app.Logger.Warn(fmt.Sprintf(...))`",
	"LogError":       "`app.Logger.Error(message)`",
	"LogErrorf":      "`app.Logger.Error(fmt.Sprintf(...))`",
	"LogFatal":       "`app.Logger.Error(message)` + `os.Exit(1)`",
	"LogFatalf":      "`app.Logger.Error(fmt.Sprintf(...))` + `os.Exit(1)`",
	"LogSetLogLevel": "set `application.Options.LogLevel` (log/slog level) at startup",

	"MenuSetApplicationMenu":    "rebuild the menu with `app.NewMenu()` and apply it with `app.Menu.SetApplicationMenu(menu)`",
	"MenuUpdateApplicationMenu": "`app.Menu.UpdateApplicationMenu()`",

	"OpenDirectoryDialog":     "`app.Dialog.OpenFileWithOptions(&application.OpenFileDialogOptions{CanChooseDirectories: true, CanChooseFiles: false, ...}).PromptForSingleSelection()`",
	"OpenFileDialog":          "`app.Dialog.OpenFileWithOptions(&application.OpenFileDialogOptions{...}).PromptForSingleSelection()` (field names differ slightly, e.g. DefaultDirectory -> Directory)",
	"OpenMultipleFilesDialog": "`app.Dialog.OpenFileWithOptions(&application.OpenFileDialogOptions{AllowsMultipleSelection: true, ...}).PromptForMultipleSelection()`",
	"SaveFileDialog":          "`app.Dialog.SaveFileWithOptions(&application.SaveFileDialogOptions{...}).PromptForSingleSelection()`",
	"MessageDialog":           "`app.Dialog.Info()/Question()/Warning()/Error()` with `.AddButton(label).OnClick(func(){...})` - the result arrives via button callbacks, not a return value",

	"OnFileDrop":    "`window.OnWindowEvent(events.Common.WindowFilesDropped, func(e *application.WindowEvent) { e.Context().DroppedFiles() })` - requires `WebviewWindowOptions.EnableFileDrop: true`",
	"OnFileDropOff": "call the unsubscribe func returned by `OnWindowEvent`",
}

// windowRuntimeAdvice maps v2 Window* functions to the v3 window method.
// They all operate on a window object: `app.Window.Current()` or a window you
// keep a reference to.
var windowRuntimeAdvice = map[string]string{
	"WindowSetTitle":              "`window.SetTitle(title)`",
	"WindowFullscreen":            "`window.Fullscreen()`",
	"WindowUnfullscreen":          "`window.UnFullscreen()`",
	"WindowCenter":                "`window.Center()`",
	"WindowReload":                "`window.Reload()`",
	"WindowReloadApp":             "`window.ForceReload()`",
	"WindowShow":                  "`window.Show()`",
	"WindowHide":                  "`window.Hide()`",
	"WindowSetSize":               "`window.SetSize(width, height)`",
	"WindowGetSize":               "`window.Size()`",
	"WindowSetMinSize":            "`window.SetMinSize(width, height)`",
	"WindowSetMaxSize":            "`window.SetMaxSize(width, height)`",
	"WindowSetAlwaysOnTop":        "`window.SetAlwaysOnTop(b)`",
	"WindowSetPosition":           "`window.SetRelativePosition(x, y)`",
	"WindowGetPosition":           "`window.RelativePosition()`",
	"WindowMaximise":              "`window.Maximise()`",
	"WindowToggleMaximise":        "`window.ToggleMaximise()`",
	"WindowUnmaximise":            "`window.UnMaximise()`",
	"WindowMinimise":              "`window.Minimise()`",
	"WindowUnminimise":            "`window.UnMinimise()`",
	"WindowIsFullscreen":          "`window.IsFullscreen()`",
	"WindowIsMaximised":           "`window.IsMaximised()`",
	"WindowIsMinimised":           "`window.IsMinimised()`",
	"WindowIsNormal":              "combine `!window.IsFullscreen() && !window.IsMaximised() && !window.IsMinimised()`",
	"WindowExecJS":                "`window.ExecJS(js)`",
	"WindowSetBackgroundColour":   "`window.SetBackgroundColour(application.RGBA{Red: r, Green: g, Blue: b, Alpha: a})`",
	"WindowPrint":                 "`window.Print()`",
	"WindowSetSystemDefaultTheme": "set `WebviewWindowOptions.Theme: application.SystemDefault` at window creation (v3 has no runtime theme setter)",
	"WindowSetLightTheme":         "set `WebviewWindowOptions.Theme: application.Light` at window creation (v3 has no runtime theme setter)",
	"WindowSetDarkTheme":          "set `WebviewWindowOptions.Theme: application.Dark` at window creation (v3 has no runtime theme setter)",
}

// adviseGoRuntimeCalls records every call into the v2 runtime package with
// its v3 replacement.
func adviseGoRuntimeCalls(fset *token.FileSet, files map[string]*ast.File, proj *V2Project) {
	for path, file := range files {
		localName := ""
		for name, ipath := range importMap(file) {
			if ipath == V2RuntimeImport {
				localName = name
			}
		}
		if localName == "" {
			continue
		}
		rel, err := filepath.Rel(proj.Dir, path)
		if err != nil {
			rel = path
		}
		ast.Inspect(file, func(n ast.Node) bool {
			sel, ok := n.(*ast.SelectorExpr)
			if !ok {
				return true
			}
			ident, ok := sel.X.(*ast.Ident)
			if !ok || ident.Name != localName {
				return true
			}
			name := sel.Sel.Name
			advice, ok := goRuntimeAdvice[name]
			if !ok {
				advice, ok = windowRuntimeAdvice[name]
				if ok {
					advice += " - get the window with `app.Window.Current()` or keep a reference to the one you create"
				}
			}
			if !ok {
				// Type references (runtime.OpenDialogOptions{...}) and
				// anything unknown.
				advice = "see the v3 application API and https://v3.wails.io/migration/v2-to-v3/"
			}
			pos := fset.Position(sel.Pos())
			proj.Report.CallSite(fmt.Sprintf("%s:%d", rel, pos.Line), "`runtime."+name+"`", advice)
			return true
		})
	}
}

var wailsjsImportRe = regexp.MustCompile(`(?:from\s*|require\s*\(\s*)['"]([^'"]*wailsjs/(runtime|go)/[^'"]*)['"]`)

// frontendSourceExts are the file types scanned for wailsjs imports.
var frontendSourceExts = map[string]bool{
	".js": true, ".jsx": true, ".ts": true, ".tsx": true,
	".svelte": true, ".vue": true, ".html": true, ".mjs": true, ".cjs": true,
}

// adviseFrontendImports records every wailsjs import in the frontend sources
// with its v3 replacement.
func adviseFrontendImports(proj *V2Project) error {
	frontend := proj.FrontendDir
	if _, err := os.Stat(frontend); os.IsNotExist(err) {
		return nil
	}
	return filepath.Walk(frontend, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			switch info.Name() {
			case "node_modules", "dist", "wailsjs":
				return filepath.SkipDir
			}
			return nil
		}
		if !frontendSourceExts[filepath.Ext(path)] {
			return nil
		}
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()
		rel, rerr := filepath.Rel(proj.Dir, path)
		if rerr != nil {
			rel = path
		}
		scanner := bufio.NewScanner(f)
		lineNo := 0
		for scanner.Scan() {
			lineNo++
			m := wailsjsImportRe.FindStringSubmatch(scanner.Text())
			if m == nil {
				continue
			}
			var advice string
			if m[2] == "runtime" {
				advice = "import from `@wailsio/runtime` instead: `import {Events, Window, Dialogs, ...} from '@wailsio/runtime'`. Function names change, e.g. `EventsOn(name, cb)` -> `Events.On(name, cb)` (the callback receives an event object; your payload is `event.data`), `WindowSetTitle` -> `Window.SetTitle`, `Quit` -> `Application.Quit`."
			} else {
				advice = "run `wails3 generate bindings`, then import the service from `frontend/bindings`: `import {" + importedServiceName(m[1]) + "} from './bindings/" + proj.ModulePath + "'` and call methods on it"
			}
			proj.Report.CallSite(fmt.Sprintf("%s:%d", rel, lineNo), "`"+m[1]+"`", advice)
		}
		return scanner.Err()
	})
}

// importedServiceName extracts the bound struct name from a wailsjs/go import
// path such as ../wailsjs/go/main/App.
func importedServiceName(importPath string) string {
	base := filepath.Base(importPath)
	return strings.TrimSuffix(base, filepath.Ext(base))
}
