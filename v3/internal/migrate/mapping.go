package migrate

import (
	"go/ast"
	"go/token"
	"sort"
	"strings"
)

// GenField is one field of a generated composite literal: Name: Expr.
type GenField struct {
	Name string
	Expr string
}

// V3Options is the mapped result of a v2 options.App literal, ready for code
// generation. All Expr values are Go source snippets valid in the migrated
// main file.
type V3Options struct {
	App      []GenField // application.Options
	AppMac   []GenField // application.Options.Mac
	AppWin   []GenField // application.Options.Windows
	AppLinux []GenField // application.Options.Linux

	Win      []GenField // application.WebviewWindowOptions
	WinMac   []GenField // WebviewWindowOptions.Mac
	WinWin   []GenField // WebviewWindowOptions.Windows
	WinLinux []GenField // WebviewWindowOptions.Linux

	SingleInstance []GenField // *application.SingleInstanceOptions, nil if unset

	// Services holds the Go expressions to wrap in application.NewService().
	Services []string

	// v2 lifecycle callbacks (original source expressions, signature
	// func(context.Context)). Empty when unset.
	OnStartup  string
	OnDomReady string
	OnShutdown string
	// OnBeforeClose (v2: return true to prevent close) is wrapped into
	// application.Options.ShouldQuit.
	OnBeforeClose string
}

// NeedsLifecycleService reports whether the generated main needs the
// v2compat lifecycle service.
func (o *V3Options) NeedsLifecycleService() bool {
	return o.OnStartup != "" || o.OnDomReady != "" || o.OnShutdown != ""
}

// mapper carries shared state while walking the options.App literal.
type mapper struct {
	proj *V2Project
	out  *V3Options
}

// MapOptions converts the parsed v2 options.App literal into V3Options,
// recording everything unmappable in the project report.
func MapOptions(proj *V2Project) *V3Options {
	m := &mapper{proj: proj, out: &V3Options{}}
	report := proj.Report

	lit := proj.Main.AppLit
	if lit == nil {
		report.Manual("options.App",
			"The value passed to wails.Run was not a literal `&options.App{...}`, so options could not be migrated automatically. A default window configuration was generated; port your options manually to application.Options / application.WebviewWindowOptions.")
		return m.out
	}

	for _, elt := range lit.Elts {
		kv, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			continue
		}
		key, ok := kv.Key.(*ast.Ident)
		if !ok {
			continue
		}
		m.mapAppField(key.Name, kv.Value)
	}

	// v2 apps always load the embedded index; v3 needs the URL explicitly.
	m.out.Win = append(m.out.Win, GenField{"URL", `"/"`})

	return m.out
}

// src returns the original source text of an expression.
func (m *mapper) src(e ast.Expr) string {
	return exprText(m.proj.Main.Fset, m.proj.Main.Source, e)
}

// referencesV2 reports whether the expression mentions any package imported
// from github.com/wailsapp/wails/v2 (such values cannot be carried verbatim).
func (m *mapper) referencesV2(e ast.Expr) bool {
	found := false
	ast.Inspect(e, func(n ast.Node) bool {
		sel, ok := n.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		ident, ok := sel.X.(*ast.Ident)
		if !ok {
			return true
		}
		if path, ok := m.proj.Main.Imports[ident.Name]; ok {
			if strings.HasPrefix(path, "github.com/wailsapp/wails/v2") {
				found = true
				return false
			}
		}
		return true
	})
	return found
}

// carry appends `name: <src>` to dst when the value is representable in the
// migrated program; otherwise it records a manual step.
func (m *mapper) carry(dst *[]GenField, v2Name, v3Name string, value ast.Expr, v3Target string) {
	if m.referencesV2(value) {
		m.proj.Report.Manual("options.App."+v2Name,
			"The value `"+m.src(value)+"` references Wails v2 packages and could not be carried over. Set `"+v3Target+"` manually.")
		return
	}
	*dst = append(*dst, GenField{v3Name, m.src(value)})
	m.proj.Report.Mapped("options.App."+v2Name, v3Target)
}

// manual is a shorthand for recording a manual step for an option.
func (m *mapper) manual(v2Name, instructions string) {
	m.proj.Report.Manual("options.App."+v2Name, instructions)
}

// note is a shorthand for adding a note to the report.
func (m *mapper) note(s string) { m.proj.Report.Note(s) }

// compositeLit unwraps &T{...} / T{...} values to the composite literal, or
// nil when the value has another shape.
func compositeLit(e ast.Expr) *ast.CompositeLit {
	if unary, ok := e.(*ast.UnaryExpr); ok && unary.Op == token.AND {
		e = unary.X
	}
	lit, _ := e.(*ast.CompositeLit)
	return lit
}

// litField is a key/value pair of a composite literal, in source order.
type litField struct {
	Name  string
	Value ast.Expr
}

// orderedFields returns the key/value pairs of a composite literal in source
// order (keeps generated code and reports deterministic).
func orderedFields(lit *ast.CompositeLit) []litField {
	var fields []litField
	for _, elt := range lit.Elts {
		if kv, ok := elt.(*ast.KeyValueExpr); ok {
			if key, ok := kv.Key.(*ast.Ident); ok {
				fields = append(fields, litField{key.Name, kv.Value})
			}
		}
	}
	return fields
}

// litFields returns the key/value pairs of a composite literal as a map (for
// single-field lookups).
func litFields(lit *ast.CompositeLit) map[string]ast.Expr {
	fields := map[string]ast.Expr{}
	for _, f := range orderedFields(lit) {
		fields[f.Name] = f.Value
	}
	return fields
}

// selectorConst maps a v2 enum selector (e.g. windows.Dark) to a v3
// application-package expression via a lookup of the selector's final name.
// Returns "" if the value has another shape or the name is not in the table.
func selectorConst(value ast.Expr, table map[string]string) string {
	name := ""
	switch v := value.(type) {
	case *ast.SelectorExpr:
		name = v.Sel.Name
	case *ast.Ident:
		name = v.Name
	default:
		return ""
	}
	return table[name]
}

func isTrue(e ast.Expr) bool {
	ident, ok := e.(*ast.Ident)
	return ok && ident.Name == "true"
}

var windowStartStates = map[string]string{
	"Normal":     "application.WindowStateNormal",
	"Maximised":  "application.WindowStateMaximised",
	"Minimised":  "application.WindowStateMinimised",
	"Fullscreen": "application.WindowStateFullscreen",
}

var windowsThemes = map[string]string{
	"SystemDefault": "application.SystemDefault",
	"Dark":          "application.Dark",
	"Light":         "application.Light",
}

var windowsBackdropTypes = map[string]string{
	"Auto":    "application.Auto",
	"None":    "application.None",
	"Mica":    "application.Mica",
	"Acrylic": "application.Acrylic",
	"Tabbed":  "application.Tabbed",
}

var linuxGpuPolicies = map[string]string{
	"WebviewGpuPolicyAlways":   "application.WebviewGpuPolicyAlways",
	"WebviewGpuPolicyOnDemand": "application.WebviewGpuPolicyOnDemand",
	"WebviewGpuPolicyNever":    "application.WebviewGpuPolicyNever",
}

// mac.Appearance* names are identical in v3's MacAppearanceType consts.
var macAppearances = map[string]string{
	"DefaultAppearance":                                     "application.DefaultAppearance",
	"NSAppearanceNameAqua":                                  "application.NSAppearanceNameAqua",
	"NSAppearanceNameDarkAqua":                              "application.NSAppearanceNameDarkAqua",
	"NSAppearanceNameVibrantLight":                          "application.NSAppearanceNameVibrantLight",
	"NSAppearanceNameAccessibilityHighContrastAqua":         "application.NSAppearanceNameAccessibilityHighContrastAqua",
	"NSAppearanceNameAccessibilityHighContrastDarkAqua":     "application.NSAppearanceNameAccessibilityHighContrastDarkAqua",
	"NSAppearanceNameAccessibilityHighContrastVibrantLight": "application.NSAppearanceNameAccessibilityHighContrastVibrantLight",
	"NSAppearanceNameAccessibilityHighContrastVibrantDark":  "application.NSAppearanceNameAccessibilityHighContrastVibrantDark",
}

// v2 mac.TitleBar* preset constructors -> v3 preset variables.
var macTitleBarPresets = map[string]string{
	"TitleBarDefault":     "application.MacTitleBarDefault",
	"TitleBarHidden":      "application.MacTitleBarHidden",
	"TitleBarHiddenInset": "application.MacTitleBarHiddenInset",
}

// v2 mac.TitleBar field -> v3 application.MacTitleBar field.
var macTitleBarFields = map[string]string{
	"TitlebarAppearsTransparent": "AppearsTransparent",
	"HideTitleBar":               "Hide",
	"HideTitle":                  "HideTitle",
	"FullSizeContent":            "FullSizeContent",
	"UseToolbar":                 "UseToolbar",
	"HideToolbarSeparator":       "HideToolbarSeparator",
}

func (m *mapper) mapAppField(name string, value ast.Expr) {
	out := m.out
	switch name {
	case "Title":
		m.carry(&out.App, "Title", "Name", value, "application.Options.Name")
		out.Win = append(out.Win, GenField{"Title", m.src(value)})
	case "Width":
		m.carry(&out.Win, "Width", "Width", value, "WebviewWindowOptions.Width")
	case "Height":
		m.carry(&out.Win, "Height", "Height", value, "WebviewWindowOptions.Height")
	case "MinWidth":
		m.carry(&out.Win, "MinWidth", "MinWidth", value, "WebviewWindowOptions.MinWidth")
	case "MinHeight":
		m.carry(&out.Win, "MinHeight", "MinHeight", value, "WebviewWindowOptions.MinHeight")
	case "MaxWidth":
		m.carry(&out.Win, "MaxWidth", "MaxWidth", value, "WebviewWindowOptions.MaxWidth")
	case "MaxHeight":
		m.carry(&out.Win, "MaxHeight", "MaxHeight", value, "WebviewWindowOptions.MaxHeight")
	case "DisableResize":
		m.carry(&out.Win, "DisableResize", "DisableResize", value, "WebviewWindowOptions.DisableResize")
	case "Frameless":
		m.carry(&out.Win, "Frameless", "Frameless", value, "WebviewWindowOptions.Frameless")
	case "AlwaysOnTop":
		m.carry(&out.Win, "AlwaysOnTop", "AlwaysOnTop", value, "WebviewWindowOptions.AlwaysOnTop")
	case "StartHidden":
		m.carry(&out.Win, "StartHidden", "Hidden", value, "WebviewWindowOptions.Hidden")
	case "Fullscreen":
		if isTrue(value) {
			out.Win = append(out.Win, GenField{"StartState", "application.WindowStateFullscreen"})
			m.proj.Report.Mapped("options.App.Fullscreen", "WebviewWindowOptions.StartState")
		} else {
			m.manual("Fullscreen", "Non-constant Fullscreen value; set `WebviewWindowOptions.StartState` to `application.WindowStateFullscreen` as needed.")
		}
	case "WindowStartState":
		if v3 := selectorConst(value, windowStartStates); v3 != "" {
			out.Win = append(out.Win, GenField{"StartState", v3})
			m.proj.Report.Mapped("options.App.WindowStartState", "WebviewWindowOptions.StartState")
		} else {
			m.manual("WindowStartState", "Could not map the value `"+m.src(value)+"`; set `WebviewWindowOptions.StartState` manually.")
		}
	case "BackgroundColour":
		m.mapBackgroundColour(value)
	case "AssetServer":
		m.mapAssetServer(value)
	case "Assets":
		// Deprecated v2 field: Assets fs.FS.
		out.App = append(out.App, GenField{"Assets", "application.AssetOptions{\n\t\t\tHandler: application.AssetFileServerFS(" + m.src(value) + "),\n\t\t}"})
		m.proj.Report.Mapped("options.App.Assets", "application.Options.Assets")
	case "AssetsHandler":
		out.App = append(out.App, GenField{"Assets", "application.AssetOptions{\n\t\t\tHandler: " + m.src(value) + ",\n\t\t}"})
		m.proj.Report.Mapped("options.App.AssetsHandler", "application.Options.Assets.Handler")
	case "Menu":
		m.manual("Menu", "v3 menus use a different API: create the menu with `app.NewMenu()` and assign it with `app.Menu.SetApplicationMenu(menu)`. See https://v3.wails.io/learn/menus/.")
	case "Logger", "LogLevel", "LogLevelProduction":
		m.manual(name, "v3 uses the standard library `log/slog`: set `application.Options.Logger` (a *slog.Logger) and `application.Options.LogLevel`.")
	case "OnStartup":
		out.OnStartup = m.src(value)
		m.proj.Report.Mapped("options.App.OnStartup", "v2compat lifecycle service (ServiceStartup)")
	case "OnDomReady":
		out.OnDomReady = m.src(value)
		m.proj.Report.Mapped("options.App.OnDomReady", "v2compat lifecycle service (WindowRuntimeReady)")
	case "OnShutdown":
		out.OnShutdown = m.src(value)
		m.proj.Report.Mapped("options.App.OnShutdown", "v2compat lifecycle service (ServiceShutdown)")
	case "OnBeforeClose":
		out.OnBeforeClose = m.src(value)
		m.proj.Report.Mapped("options.App.OnBeforeClose", "application.Options.ShouldQuit")
	case "Bind":
		for _, bt := range m.proj.BoundTypes {
			out.Services = append(out.Services, bt.Expr)
		}
		m.proj.Report.Mapped("options.App.Bind", "application.Options.Services")
	case "EnumBind":
		m.manual("EnumBind", "v3's binding generator discovers enum types automatically from service method signatures; remove EnumBind and re-run `wails3 generate bindings`.")
	case "SingleInstanceLock":
		m.mapSingleInstance(value)
	case "Windows":
		m.mapWindowsOptions(value)
	case "Mac":
		m.mapMacOptions(value)
	case "Linux":
		m.mapLinuxOptions(value)
	case "Debug":
		if lit := compositeLit(value); lit != nil {
			if v, ok := litFields(lit)["OpenInspectorOnStartup"]; ok {
				m.carry(&out.Win, "Debug.OpenInspectorOnStartup", "OpenInspectorOnStartup", v, "WebviewWindowOptions.OpenInspectorOnStartup")
			}
		}
	case "DragAndDrop":
		m.mapDragAndDrop(value)
	case "CSSDragProperty", "CSSDragValue":
		m.note("`" + name + "` was dropped: v3 always uses the CSS property `--wails-draggable: drag`. Update your styles if you used a custom property.")
	case "EnableDefaultContextMenu":
		if isTrue(value) {
			m.note("`EnableDefaultContextMenu: true` was dropped: the v3 equivalent is `WebviewWindowOptions.DefaultContextMenuDisabled` (default enabled in dev, controllable per window).")
		}
	case "EnableFraudulentWebsiteDetection":
		m.carry(&out.WinMac, "EnableFraudulentWebsiteDetection", "EnableFraudulentWebsiteWarnings", value, "MacWindow.EnableFraudulentWebsiteWarnings")
		m.note("`EnableFraudulentWebsiteDetection` only maps to macOS in v3 (`MacWindow.EnableFraudulentWebsiteWarnings`); Windows SmartScreen is controlled by the OS.")
	case "HideWindowOnClose":
		m.manual("HideWindowOnClose", "Intercept the close in v3: `window.OnWindowEvent(events.Common.WindowClosing, func(e *application.WindowEvent) { e.Cancel(); window.Hide() })`.")
	case "ErrorFormatter":
		m.manual("ErrorFormatter", "v3 replaces ErrorFormatter with `application.Options.MarshalError func(error) []byte`.")
	case "BindingsAllowedOrigins":
		m.manual("BindingsAllowedOrigins", "Not needed in v3: binding calls are routed through the internal asset server. Remove any CORS-specific configuration.")
	case "DisablePanicRecovery":
		m.manual("DisablePanicRecovery", "v3 uses `application.Options.PanicHandler func(*application.PanicDetails)` instead.")
	case "Experimental":
		// Empty struct in v2; nothing to migrate.
	default:
		m.manual(name, "This option was not recognised by the migrator; check application.Options / application.WebviewWindowOptions for an equivalent.")
	}
}

// mapBackgroundColour converts &options.RGBA{R: , G: , B: , A: } (keyed or
// positional) to application.NewRGBA(...).
func (m *mapper) mapBackgroundColour(value ast.Expr) {
	lit := compositeLit(value)
	if lit == nil {
		m.manual("BackgroundColour", "Could not parse the colour value `"+m.src(value)+"`; set `WebviewWindowOptions.BackgroundColour` with `application.NewRGBA(r, g, b, a)`.")
		return
	}
	comps := map[string]string{"R": "0", "G": "0", "B": "0", "A": "255"}
	order := []string{"R", "G", "B", "A"}
	for i, elt := range lit.Elts {
		if kv, ok := elt.(*ast.KeyValueExpr); ok {
			if key, ok := kv.Key.(*ast.Ident); ok {
				comps[key.Name] = m.src(kv.Value)
			}
			continue
		}
		if i < len(order) {
			comps[order[i]] = m.src(elt)
		}
	}
	expr := "application.NewRGBA(" + comps["R"] + ", " + comps["G"] + ", " + comps["B"] + ", " + comps["A"] + ")"
	m.out.Win = append(m.out.Win, GenField{"BackgroundColour", expr})
	m.proj.Report.Mapped("options.App.BackgroundColour", "WebviewWindowOptions.BackgroundColour")
}

func (m *mapper) mapAssetServer(value ast.Expr) {
	lit := compositeLit(value)
	if lit == nil {
		m.manual("AssetServer", "Could not parse the AssetServer options; configure `application.Options.Assets` manually.")
		return
	}
	fields := litFields(lit)
	var gen []string
	if assets, ok := fields["Assets"]; ok {
		if m.referencesV2(assets) {
			m.manual("AssetServer.Assets", "The assets value references v2 packages; set `application.Options.Assets.Handler` manually.")
		} else {
			gen = append(gen, "Handler: application.AssetFileServerFS("+m.src(assets)+"),")
			m.proj.Report.Mapped("options.App.AssetServer.Assets", "application.Options.Assets.Handler (AssetFileServerFS)")
		}
		if _, alsoHandler := fields["Handler"]; alsoHandler {
			m.manual("AssetServer.Handler", "v2 used Handler as a fallback for requests not found in Assets. In v3 there is a single Handler: chain your fallback yourself or serve it via `application.Options.Assets.Middleware`.")
		}
	} else if handler, ok := fields["Handler"]; ok {
		if m.referencesV2(handler) {
			m.manual("AssetServer.Handler", "The handler references v2 packages; set `application.Options.Assets.Handler` manually.")
		} else {
			gen = append(gen, "Handler: "+m.src(handler)+",")
			m.proj.Report.Mapped("options.App.AssetServer.Handler", "application.Options.Assets.Handler")
		}
	}
	if mw, ok := fields["Middleware"]; ok {
		if m.referencesV2(mw) {
			m.manual("AssetServer.Middleware", "The middleware references v2 packages; port it to `application.Options.Assets.Middleware` (same `func(http.Handler) http.Handler` shape).")
		} else {
			gen = append(gen, "Middleware: application.Middleware("+m.src(mw)+"),")
			m.proj.Report.Mapped("options.App.AssetServer.Middleware", "application.Options.Assets.Middleware")
		}
	}
	if len(gen) > 0 {
		m.out.App = append(m.out.App, GenField{"Assets", "application.AssetOptions{\n\t\t\t" + strings.Join(gen, "\n\t\t\t") + "\n\t\t}"})
	}
}

func (m *mapper) mapSingleInstance(value ast.Expr) {
	lit := compositeLit(value)
	if lit == nil {
		m.manual("SingleInstanceLock", "Configure `application.Options.SingleInstance` (*application.SingleInstanceOptions) manually.")
		return
	}
	fields := litFields(lit)
	if id, ok := fields["UniqueId"]; ok && !m.referencesV2(id) {
		m.out.SingleInstance = append(m.out.SingleInstance, GenField{"UniqueID", m.src(id)})
		m.proj.Report.Mapped("options.App.SingleInstanceLock.UniqueId", "application.SingleInstanceOptions.UniqueID")
	}
	if _, ok := fields["OnSecondInstanceLaunch"]; ok {
		m.manual("SingleInstanceLock.OnSecondInstanceLaunch",
			"Port the callback to `application.SingleInstanceOptions.OnSecondInstanceLaunch(data application.SecondInstanceData)`. Note the v3 struct uses `WorkingDir` (v2: `WorkingDirectory`).")
	}
}

func (m *mapper) mapWindowsOptions(value ast.Expr) {
	lit := compositeLit(value)
	if lit == nil {
		m.manual("Windows", "Could not parse the Windows options literal; port it manually to application.WindowsOptions / application.WindowsWindow.")
		return
	}
	for _, f := range orderedFields(lit) {
		name, v := f.Name, f.Value
		prefix := "Windows." + name
		switch name {
		case "WebviewIsTransparent":
			if isTrue(v) {
				m.out.Win = append(m.out.Win, GenField{"BackgroundType", "application.BackgroundTypeTransparent"})
				m.proj.Report.Mapped("options.App."+prefix, "WebviewWindowOptions.BackgroundType")
			}
		case "WindowIsTranslucent":
			if isTrue(v) {
				m.out.Win = append(m.out.Win, GenField{"BackgroundType", "application.BackgroundTypeTranslucent"})
				m.proj.Report.Mapped("options.App."+prefix, "WebviewWindowOptions.BackgroundType")
			}
		case "DisableWindowIcon":
			m.carry(&m.out.WinWin, prefix, "DisableIcon", v, "WindowsWindow.DisableIcon")
		case "DisableFramelessWindowDecorations":
			m.carry(&m.out.WinWin, prefix, "DisableFramelessWindowDecorations", v, "WindowsWindow.DisableFramelessWindowDecorations")
		case "WebviewUserDataPath":
			m.carry(&m.out.AppWin, prefix, "WebviewUserDataPath", v, "application.WindowsOptions.WebviewUserDataPath")
		case "WebviewBrowserPath":
			m.carry(&m.out.AppWin, prefix, "WebviewBrowserPath", v, "application.WindowsOptions.WebviewBrowserPath")
		case "Theme":
			if v3 := selectorConst(v, windowsThemes); v3 != "" {
				m.out.WinWin = append(m.out.WinWin, GenField{"Theme", v3})
				m.proj.Report.Mapped("options.App."+prefix, "WindowsWindow.Theme")
			} else {
				m.manual(prefix, "Could not map the theme value `"+m.src(v)+"`; set `WindowsWindow.Theme` manually.")
			}
		case "CustomTheme":
			m.manual(prefix, "The v3 custom theme structure differs (`application.ThemeSettings` with per-mode `WindowTheme`/`MenuBarTheme` colours); port your colours manually.")
		case "BackdropType":
			if v3 := selectorConst(v, windowsBackdropTypes); v3 != "" {
				m.out.WinWin = append(m.out.WinWin, GenField{"BackdropType", v3})
				m.proj.Report.Mapped("options.App."+prefix, "WindowsWindow.BackdropType")
			} else {
				m.manual(prefix, "Could not map the backdrop value `"+m.src(v)+"`; set `WindowsWindow.BackdropType` manually.")
			}
		case "Messages":
			m.note("`Windows.Messages` was dropped: v3 handles the WebView2 bootstrap flow itself.")
		case "ResizeDebounceMS":
			m.manual(prefix, "No direct v3 equivalent. `WindowsWindow.WindowDidMoveDebounceMS` debounces move events if that is what you need.")
		case "OnSuspend":
			m.manual(prefix, "Subscribe to the v3 event instead: `app.Event.OnApplicationEvent(events.Common.SystemWillSleep, ...)`.")
		case "OnResume":
			m.manual(prefix, "Subscribe to the v3 event instead: `app.Event.OnApplicationEvent(events.Common.SystemDidWake, ...)`.")
		case "WebviewGpuIsDisabled":
			m.manual(prefix, "No direct v3 equivalent; consider `application.WindowsOptions.AdditionalBrowserArgs: []string{\"--disable-gpu\"}`.")
		case "WebviewDisableRendererCodeIntegrity":
			m.manual(prefix, "No direct v3 equivalent; consider `application.WindowsOptions.AdditionalBrowserArgs`.")
		case "EnableSwipeGestures":
			m.carry(&m.out.WinWin, prefix, "EnableSwipeGestures", v, "WindowsWindow.EnableSwipeGestures")
		case "ZoomFactor":
			m.carry(&m.out.Win, prefix, "Zoom", v, "WebviewWindowOptions.Zoom")
		case "IsZoomControlEnabled":
			m.carry(&m.out.Win, prefix, "ZoomControlEnabled", v, "WebviewWindowOptions.ZoomControlEnabled")
		case "DisablePinchZoom":
			m.manual(prefix, "No direct v3 equivalent; zoom control is `WebviewWindowOptions.ZoomControlEnabled`.")
		case "ContentProtection":
			m.carry(&m.out.Win, prefix, "ContentProtectionEnabled", v, "WebviewWindowOptions.ContentProtectionEnabled")
		case "WindowClassName":
			m.carry(&m.out.AppWin, prefix, "WndClass", v, "application.WindowsOptions.WndClass")
		case "DLLSearchPaths":
			m.manual(prefix, "No v3 equivalent; v3 loads WebView2Loader statically.")
		default:
			m.manual(prefix, "This option was not recognised by the migrator; check application.WindowsOptions / application.WindowsWindow for an equivalent.")
		}
	}
}

func (m *mapper) mapMacOptions(value ast.Expr) {
	lit := compositeLit(value)
	if lit == nil {
		m.manual("Mac", "Could not parse the Mac options literal; port it manually to application.MacOptions / application.MacWindow.")
		return
	}
	for _, f := range orderedFields(lit) {
		name, v := f.Name, f.Value
		prefix := "Mac." + name
		switch name {
		case "TitleBar":
			m.mapMacTitleBar(v)
		case "Appearance":
			if v3 := selectorConst(v, macAppearances); v3 != "" {
				m.out.WinMac = append(m.out.WinMac, GenField{"Appearance", v3})
				m.proj.Report.Mapped("options.App."+prefix, "MacWindow.Appearance")
			} else {
				m.manual(prefix, "Could not map the appearance value `"+m.src(v)+"`; set `MacWindow.Appearance` manually.")
			}
		case "WebviewIsTransparent":
			if isTrue(v) {
				m.out.Win = append(m.out.Win, GenField{"BackgroundType", "application.BackgroundTypeTransparent"})
				m.proj.Report.Mapped("options.App."+prefix, "WebviewWindowOptions.BackgroundType")
			}
		case "WindowIsTranslucent":
			if isTrue(v) {
				m.out.WinMac = append(m.out.WinMac, GenField{"Backdrop", "application.MacBackdropTranslucent"})
				m.proj.Report.Mapped("options.App."+prefix, "MacWindow.Backdrop")
			}
		case "Preferences":
			m.manual(prefix, "Port to `MacWindow.WebviewPreferences` (application.MacWebviewPreferences) manually.")
		case "DisableZoom":
			m.manual(prefix, "No direct v3 equivalent; `WebviewWindowOptions.ZoomControlEnabled: false` disables zoom controls.")
		case "About":
			m.mapMacAbout(v)
		case "OnFileOpen":
			m.manual(prefix, "Subscribe to the v3 event instead: `app.Event.OnApplicationEvent(events.Common.ApplicationOpenedWithFile, ...)`, reading the filename from the event context.")
		case "OnUrlOpen":
			m.manual(prefix, "Subscribe to the v3 event instead: `app.Event.OnApplicationEvent(events.Common.ApplicationLaunchedWithUrl, ...)`.")
		case "DisableEscapeExitsFullscreen":
			m.carry(&m.out.WinMac, prefix, "DisableEscapeExitsFullscreen", v, "MacWindow.DisableEscapeExitsFullscreen")
		case "ContentProtection":
			m.carry(&m.out.Win, prefix, "ContentProtectionEnabled", v, "WebviewWindowOptions.ContentProtectionEnabled")
		default:
			m.manual(prefix, "This option was not recognised by the migrator; check application.MacOptions / application.MacWindow for an equivalent.")
		}
	}
}

func (m *mapper) mapMacTitleBar(value ast.Expr) {
	// Preset constructor: mac.TitleBarHiddenInset()
	if call, ok := value.(*ast.CallExpr); ok {
		if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
			if v3 := macTitleBarPresets[sel.Sel.Name]; v3 != "" {
				m.out.WinMac = append(m.out.WinMac, GenField{"TitleBar", v3})
				m.proj.Report.Mapped("options.App.Mac.TitleBar", "MacWindow.TitleBar")
				return
			}
		}
	}
	// Literal: &mac.TitleBar{...}
	if lit := compositeLit(value); lit != nil {
		var fields []string
		ok := true
		for _, f := range orderedFields(lit) {
			name, v := f.Name, f.Value
			v3Name := macTitleBarFields[name]
			if v3Name == "" || m.referencesV2(v) {
				m.manual("Mac.TitleBar."+name, "Could not map this title bar field; check application.MacTitleBar.")
				ok = false
				continue
			}
			fields = append(fields, v3Name+": "+m.src(v)+",")
		}
		if ok && len(fields) > 0 {
			sort.Strings(fields)
			m.out.WinMac = append(m.out.WinMac, GenField{"TitleBar", "application.MacTitleBar{\n\t\t\t\t" + strings.Join(fields, "\n\t\t\t\t") + "\n\t\t\t}"})
			m.proj.Report.Mapped("options.App.Mac.TitleBar", "MacWindow.TitleBar")
		}
		return
	}
	m.manual("Mac.TitleBar", "Could not map the title bar value `"+m.src(value)+"`; set `MacWindow.TitleBar` manually.")
}

func (m *mapper) mapMacAbout(value ast.Expr) {
	lit := compositeLit(value)
	if lit == nil {
		m.manual("Mac.About", "v3 renders the about box from application.Options Name/Description/Icon; set those manually.")
		return
	}
	fields := litFields(lit)
	if msg, ok := fields["Message"]; ok && !m.referencesV2(msg) {
		m.out.App = append(m.out.App, GenField{"Description", m.src(msg)})
		m.proj.Report.Mapped("options.App.Mac.About.Message", "application.Options.Description")
	}
	if icon, ok := fields["Icon"]; ok && !m.referencesV2(icon) {
		m.out.App = append(m.out.App, GenField{"Icon", m.src(icon)})
		m.proj.Report.Mapped("options.App.Mac.About.Icon", "application.Options.Icon")
	}
	if _, ok := fields["Title"]; ok {
		m.note("`Mac.About.Title` was dropped: the v3 about box title comes from `application.Options.Name`.")
	}
}

func (m *mapper) mapLinuxOptions(value ast.Expr) {
	lit := compositeLit(value)
	if lit == nil {
		m.manual("Linux", "Could not parse the Linux options literal; port it manually to application.LinuxOptions / application.LinuxWindow.")
		return
	}
	for _, f := range orderedFields(lit) {
		name, v := f.Name, f.Value
		prefix := "Linux." + name
		switch name {
		case "Icon":
			m.carry(&m.out.WinLinux, prefix, "Icon", v, "LinuxWindow.Icon")
		case "WindowIsTranslucent":
			m.carry(&m.out.WinLinux, prefix, "WindowIsTranslucent", v, "LinuxWindow.WindowIsTranslucent")
		case "WebviewGpuPolicy":
			if v3 := selectorConst(v, linuxGpuPolicies); v3 != "" {
				m.out.WinLinux = append(m.out.WinLinux, GenField{"WebviewGpuPolicy", v3})
				m.proj.Report.Mapped("options.App."+prefix, "LinuxWindow.WebviewGpuPolicy")
			} else {
				m.manual(prefix, "Could not map the GPU policy value `"+m.src(v)+"`; set `LinuxWindow.WebviewGpuPolicy` manually.")
			}
		case "ProgramName":
			m.carry(&m.out.AppLinux, prefix, "ProgramName", v, "application.LinuxOptions.ProgramName")
		case "Messages":
			m.note("`Linux.Messages` was dropped: v3 reports missing webkit dependencies itself.")
		default:
			m.manual(prefix, "This option was not recognised by the migrator; check application.LinuxOptions / application.LinuxWindow for an equivalent.")
		}
	}
}

func (m *mapper) mapDragAndDrop(value ast.Expr) {
	lit := compositeLit(value)
	if lit == nil {
		m.manual("DragAndDrop", "Set `WebviewWindowOptions.EnableFileDrop` manually.")
		return
	}
	for _, f := range orderedFields(lit) {
		name, v := f.Name, f.Value
		prefix := "DragAndDrop." + name
		switch name {
		case "EnableFileDrop":
			m.carry(&m.out.Win, prefix, "EnableFileDrop", v, "WebviewWindowOptions.EnableFileDrop")
		case "DisableWebViewDrop":
			m.manual(prefix, "No direct v3 equivalent; see `WebviewWindowOptions.EnableFileDrop` and the WindowFilesDropped event.")
		case "CSSDropProperty", "CSSDropValue":
			m.note("`" + prefix + "` was dropped: v3 controls drop targets with the `--wails-drop-target: drop` CSS property.")
		default:
			m.manual(prefix, "This option was not recognised by the migrator.")
		}
	}
}
