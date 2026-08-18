package commands

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

const (
	v2ModulePath        = "github.com/wailsapp/wails/v2"
	v2OptionsPath       = "github.com/wailsapp/wails/v2/pkg/options"
	v2AssetServerPath   = "github.com/wailsapp/wails/v2/pkg/options/assetserver"
	v2RuntimePath       = "github.com/wailsapp/wails/v2/pkg/runtime"
	generatedMainName   = "main_v3.go.example"
	migrationReportName = "MIGRATION_REPORT.md"
)

// v2RGBA is an extracted v2 BackgroundColour literal.
type v2RGBA struct {
	R, G, B, A int
}

// v2WindowOptions collects the window-related fields of options.App that map
// onto application.WebviewWindowOptions.
type v2WindowOptions struct {
	Title         string
	Width         int
	Height        int
	MinWidth      int
	MinHeight     int
	MaxWidth      int
	MaxHeight     int
	DisableResize bool
	Frameless     bool
	AlwaysOnTop   bool
	Hidden        bool
	StartState    string // "Fullscreen", "Maximised" or "Minimised"
	Background    *v2RGBA
}

// v2MainInfo is everything extracted from the project's wails.Run call.
type v2MainInfo struct {
	RunFile          string
	Window           v2WindowOptions
	AssetsIdent      string
	EmbedDirective   string
	BindExprs        []string
	ConstructorStmts []string
	OnStartup        string
	OnDomReady       string
	OnShutdown       string
	OnBeforeClose    string
	SingleInstanceID string
	SingleInstanceCB bool
	Mapped           []string
	Manual           []string
}

// extractV2Main parses the Go files in the project root and extracts the
// options.App literal passed to wails.Run.
func extractV2Main(projectDir string) (*v2MainInfo, []string, error) {
	fset := token.NewFileSet()
	entries, err := os.ReadDir(projectDir)
	if err != nil {
		return nil, nil, err
	}
	var files []*ast.File
	var runtimeCalls []string
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}
		file, err := parser.ParseFile(fset, filepath.Join(projectDir, name), nil, parser.ParseComments)
		if err != nil {
			continue
		}
		files = append(files, file)
		runtimeCalls = append(runtimeCalls, scanRuntimeCalls(fset, file)...)
	}

	var lit *ast.CompositeLit
	var litFile *ast.File
	for _, file := range files {
		if found := findOptionsAppLiteral(file); found != nil {
			lit = found
			litFile = file
			break
		}
	}
	if lit == nil {
		return nil, runtimeCalls, fmt.Errorf("could not find the options.App literal passed to wails.Run")
	}

	info := &v2MainInfo{
		RunFile: filepath.Base(fset.Position(lit.Pos()).Filename),
	}
	pos := func(node ast.Node) string {
		p := fset.Position(node.Pos())
		return fmt.Sprintf("%s:%d", filepath.Base(p.Filename), p.Line)
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
		switch key.Name {
		case "Title":
			if s, ok := litString(kv.Value); ok {
				info.Window.Title = s
				info.Mapped = append(info.Mapped, "Title")
			} else {
				info.Manual = append(info.Manual, fmt.Sprintf("`Title` is not a string literal (%s); set the window title manually.", pos(kv)))
			}
		case "Width", "Height", "MinWidth", "MinHeight", "MaxWidth", "MaxHeight":
			if n, ok := litInt(kv.Value); ok {
				switch key.Name {
				case "Width":
					info.Window.Width = n
				case "Height":
					info.Window.Height = n
				case "MinWidth":
					info.Window.MinWidth = n
				case "MinHeight":
					info.Window.MinHeight = n
				case "MaxWidth":
					info.Window.MaxWidth = n
				case "MaxHeight":
					info.Window.MaxHeight = n
				}
				info.Mapped = append(info.Mapped, key.Name)
			} else {
				info.Manual = append(info.Manual, fmt.Sprintf("`%s` is not an integer literal (%s); set it manually in WebviewWindowOptions.", key.Name, pos(kv)))
			}
		case "DisableResize", "Frameless", "AlwaysOnTop", "StartHidden", "Fullscreen":
			if b, ok := litBool(kv.Value); ok {
				switch key.Name {
				case "DisableResize":
					info.Window.DisableResize = b
				case "Frameless":
					info.Window.Frameless = b
				case "AlwaysOnTop":
					info.Window.AlwaysOnTop = b
				case "StartHidden":
					info.Window.Hidden = b
				case "Fullscreen":
					if b {
						info.Window.StartState = "Fullscreen"
					}
				}
				info.Mapped = append(info.Mapped, key.Name)
			} else {
				info.Manual = append(info.Manual, fmt.Sprintf("`%s` is not a boolean literal (%s); set it manually in WebviewWindowOptions.", key.Name, pos(kv)))
			}
		case "WindowStartState":
			if sel, ok := kv.Value.(*ast.SelectorExpr); ok {
				switch sel.Sel.Name {
				case "Fullscreen", "Maximised", "Minimised":
					info.Window.StartState = sel.Sel.Name
					info.Mapped = append(info.Mapped, "WindowStartState")
				case "Normal":
					info.Mapped = append(info.Mapped, "WindowStartState")
				default:
					info.Manual = append(info.Manual, fmt.Sprintf("`WindowStartState` value not recognised (%s).", pos(kv)))
				}
			} else {
				info.Manual = append(info.Manual, fmt.Sprintf("`WindowStartState` is not a constant (%s); set StartState manually.", pos(kv)))
			}
		case "BackgroundColour", "BackgroundColor":
			if rgba, ok := extractRGBA(kv.Value); ok {
				info.Window.Background = rgba
				info.Mapped = append(info.Mapped, "BackgroundColour")
			} else {
				info.Manual = append(info.Manual, fmt.Sprintf("`BackgroundColour` could not be extracted (%s); set `BackgroundColour: application.NewRGBA(...)` manually.", pos(kv)))
			}
		case "Assets":
			if ident, ok := kv.Value.(*ast.Ident); ok {
				info.AssetsIdent = ident.Name
				info.Mapped = append(info.Mapped, "Assets")
			} else {
				info.Manual = append(info.Manual, fmt.Sprintf("`Assets` is not a plain identifier (%s); configure application.AssetOptions manually.", pos(kv)))
			}
		case "AssetServer":
			extractAssetServer(kv.Value, info, pos)
		case "AssetsHandler":
			info.Manual = append(info.Manual, fmt.Sprintf("`AssetsHandler` (%s): use `Assets: application.AssetOptions{Handler: yourHandler}` in v3.", pos(kv)))
		case "Bind":
			if bindLit, ok := kv.Value.(*ast.CompositeLit); ok {
				for _, element := range bindLit.Elts {
					info.BindExprs = append(info.BindExprs, exprString(fset, element))
				}
				info.Mapped = append(info.Mapped, "Bind")
			} else {
				info.Manual = append(info.Manual, fmt.Sprintf("`Bind` is not a slice literal (%s); register your services manually with application.NewService.", pos(kv)))
			}
		case "EnumBind":
			info.Manual = append(info.Manual, fmt.Sprintf("`EnumBind` (%s): v3 generates enum bindings automatically for types used by service methods; run `wails3 generate bindings`.", pos(kv)))
		case "OnStartup":
			info.OnStartup = exprString(fset, kv.Value)
		case "OnDomReady":
			info.OnDomReady = exprString(fset, kv.Value)
		case "OnShutdown":
			info.OnShutdown = exprString(fset, kv.Value)
		case "OnBeforeClose":
			info.OnBeforeClose = exprString(fset, kv.Value)
			info.Manual = append(info.Manual, fmt.Sprintf("`OnBeforeClose` (%s): use `ShouldQuit func() bool` in application.Options, or window close events.", pos(kv)))
		case "SingleInstanceLock":
			extractSingleInstance(kv.Value, info, pos)
		case "Menu":
			info.Manual = append(info.Manual, fmt.Sprintf("`Menu` (%s): rebuild the menu with `app.NewMenu()` / `menu.AddSubmenu(...)` and `app.Menu.SetApplicationMenu(...)`; see the v3 menus documentation.", pos(kv)))
		case "Windows", "Mac", "Linux":
			info.Manual = append(info.Manual, fmt.Sprintf("`%s` platform options (%s): v3 splits these between application.Options.%s and WebviewWindowOptions.%s; reapply them manually.", key.Name, pos(kv), key.Name, key.Name))
		case "Logger", "LogLevel", "LogLevelProduction":
			info.Manual = append(info.Manual, fmt.Sprintf("`%s` (%s): v3 uses `log/slog`; set `Logger`/`LogLevel` in application.Options.", key.Name, pos(kv)))
		case "HideWindowOnClose":
			info.Manual = append(info.Manual, fmt.Sprintf("`HideWindowOnClose` (%s): handle the window closing event and call `window.Hide()` in v3.", pos(kv)))
		case "DragAndDrop":
			info.Manual = append(info.Manual, fmt.Sprintf("`DragAndDrop` (%s): configure `EnabledFeatures`/drag-and-drop on the window in v3; see the drag-and-drop documentation.", pos(kv)))
		default:
			info.Manual = append(info.Manual, fmt.Sprintf("`%s` (%s) is not migrated automatically.", key.Name, pos(kv)))
		}
	}

	// Carry over the constructor statements for bound identifiers, e.g.
	// `app := NewApp()`.
	info.ConstructorStmts = findConstructorStmts(fset, litFile, lit, info.BindExprs)

	// Find the //go:embed directive for the assets variable.
	if info.AssetsIdent != "" {
		for _, file := range files {
			if directive := findEmbedDirective(file, info.AssetsIdent); directive != "" {
				info.EmbedDirective = directive
				break
			}
		}
	}

	sort.Strings(runtimeCalls)
	return info, runtimeCalls, nil
}

// findOptionsAppLiteral returns the first composite literal of type
// options.App (any alias for the v2 options import) in the file.
func findOptionsAppLiteral(file *ast.File) *ast.CompositeLit {
	optionsName := importLocalName(file, v2OptionsPath)
	if optionsName == "" {
		return nil
	}
	var result *ast.CompositeLit
	ast.Inspect(file, func(node ast.Node) bool {
		if result != nil {
			return false
		}
		lit, ok := node.(*ast.CompositeLit)
		if !ok {
			return true
		}
		sel, ok := lit.Type.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		if ident, ok := sel.X.(*ast.Ident); ok && ident.Name == optionsName && sel.Sel.Name == "App" {
			result = lit
			return false
		}
		return true
	})
	return result
}

// importLocalName returns the local identifier a file uses for the given
// import path, or "" if the file does not import it. When the import has no
// alias, the last non-version path element is assumed (the convention for
// the wails v2 packages, whose package names match their directories, with
// the root module being package "wails").
func importLocalName(file *ast.File, path string) string {
	for _, imp := range file.Imports {
		value, err := strconv.Unquote(imp.Path.Value)
		if err != nil || value != path {
			continue
		}
		if imp.Name != nil {
			if imp.Name.Name == "_" || imp.Name.Name == "." {
				return ""
			}
			return imp.Name.Name
		}
		base := filepath.Base(value)
		if base == "v2" {
			return "wails"
		}
		return base
	}
	return ""
}

func extractAssetServer(value ast.Expr, info *v2MainInfo, pos func(ast.Node) string) {
	lit, ok := compositeLit(value)
	if !ok {
		info.Manual = append(info.Manual, fmt.Sprintf("`AssetServer` could not be extracted (%s); configure application.AssetOptions manually.", pos(value)))
		return
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
		switch key.Name {
		case "Assets":
			if ident, ok := kv.Value.(*ast.Ident); ok {
				info.AssetsIdent = ident.Name
				info.Mapped = append(info.Mapped, "AssetServer.Assets")
			} else {
				info.Manual = append(info.Manual, fmt.Sprintf("`AssetServer.Assets` is not a plain identifier (%s); configure application.AssetOptions manually.", pos(kv)))
			}
		case "Handler":
			info.Manual = append(info.Manual, fmt.Sprintf("`AssetServer.Handler` (%s): set `Assets: application.AssetOptions{Handler: ...}` in v3.", pos(kv)))
		case "Middleware":
			info.Manual = append(info.Manual, fmt.Sprintf("`AssetServer.Middleware` (%s): set `Assets: application.AssetOptions{Middleware: ...}` in v3.", pos(kv)))
		default:
			info.Manual = append(info.Manual, fmt.Sprintf("`AssetServer.%s` (%s) is not migrated automatically.", key.Name, pos(kv)))
		}
	}
}

func extractSingleInstance(value ast.Expr, info *v2MainInfo, pos func(ast.Node) string) {
	lit, ok := compositeLit(value)
	if !ok {
		info.Manual = append(info.Manual, fmt.Sprintf("`SingleInstanceLock` could not be extracted (%s); configure application.Options.SingleInstance manually.", pos(value)))
		return
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
		switch key.Name {
		case "UniqueId", "UniqueID":
			if s, ok := litString(kv.Value); ok {
				info.SingleInstanceID = s
				info.Mapped = append(info.Mapped, "SingleInstanceLock")
			} else {
				info.Manual = append(info.Manual, fmt.Sprintf("`SingleInstanceLock.UniqueId` is not a string literal (%s); set SingleInstanceOptions.UniqueID manually.", pos(kv)))
			}
		case "OnSecondInstanceLaunch":
			info.SingleInstanceCB = true
		}
	}
}

// findConstructorStmts returns the `x := ...` statements from the function
// enclosing the options literal that define identifiers referenced in Bind.
func findConstructorStmts(fset *token.FileSet, file *ast.File, lit *ast.CompositeLit, bindExprs []string) []string {
	wanted := map[string]bool{}
	for _, expr := range bindExprs {
		if isIdentifier(expr) {
			wanted[expr] = true
		}
	}
	if len(wanted) == 0 || file == nil {
		return nil
	}
	var enclosing *ast.FuncDecl
	for _, decl := range file.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok && fn.Body != nil {
			if fn.Body.Pos() <= lit.Pos() && lit.End() <= fn.Body.End() {
				enclosing = fn
				break
			}
		}
	}
	if enclosing == nil {
		return nil
	}
	var stmts []string
	for _, stmt := range enclosing.Body.List {
		assign, ok := stmt.(*ast.AssignStmt)
		if !ok || assign.Tok != token.DEFINE || len(assign.Lhs) != 1 {
			continue
		}
		ident, ok := assign.Lhs[0].(*ast.Ident)
		if !ok || !wanted[ident.Name] {
			continue
		}
		stmts = append(stmts, exprString(fset, stmt))
		delete(wanted, ident.Name)
	}
	return stmts
}

// findEmbedDirective returns the //go:embed comment attached to the variable
// declaration of the given identifier, if any.
func findEmbedDirective(file *ast.File, ident string) string {
	for _, decl := range file.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok || gd.Tok != token.VAR {
			continue
		}
		for _, spec := range gd.Specs {
			vs, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}
			for _, name := range vs.Names {
				if name.Name != ident {
					continue
				}
				for _, doc := range []*ast.CommentGroup{gd.Doc, vs.Doc} {
					if doc == nil {
						continue
					}
					for _, comment := range doc.List {
						if strings.HasPrefix(comment.Text, "//go:embed") {
							return comment.Text
						}
					}
				}
			}
		}
	}
	return ""
}

// scanRuntimeCalls records calls into the v2 runtime package, e.g.
// "app.go:42: runtime.EventsEmit".
func scanRuntimeCalls(fset *token.FileSet, file *ast.File) []string {
	localName := importLocalName(file, v2RuntimePath)
	if localName == "" {
		return nil
	}
	var calls []string
	ast.Inspect(file, func(node ast.Node) bool {
		call, ok := node.(*ast.CallExpr)
		if !ok {
			return true
		}
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		if ident, ok := sel.X.(*ast.Ident); ok && ident.Name == localName {
			p := fset.Position(call.Pos())
			calls = append(calls, fmt.Sprintf("%s:%d: `runtime.%s`", filepath.Base(p.Filename), p.Line, sel.Sel.Name))
		}
		return true
	})
	return calls
}

func compositeLit(expr ast.Expr) (*ast.CompositeLit, bool) {
	if unary, ok := expr.(*ast.UnaryExpr); ok && unary.Op == token.AND {
		expr = unary.X
	}
	lit, ok := expr.(*ast.CompositeLit)
	return lit, ok
}

func extractRGBA(expr ast.Expr) (*v2RGBA, bool) {
	// options.NewRGBA(r, g, b, a) / options.NewRGB(r, g, b)
	if call, ok := expr.(*ast.CallExpr); ok {
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return nil, false
		}
		values := make([]int, 0, 4)
		for _, arg := range call.Args {
			n, ok := litInt(arg)
			if !ok {
				return nil, false
			}
			values = append(values, n)
		}
		switch {
		case sel.Sel.Name == "NewRGBA" && len(values) == 4:
			return &v2RGBA{R: values[0], G: values[1], B: values[2], A: values[3]}, true
		case sel.Sel.Name == "NewRGB" && len(values) == 3:
			return &v2RGBA{R: values[0], G: values[1], B: values[2], A: 255}, true
		}
		return nil, false
	}
	// &options.RGBA{R: r, G: g, B: b, A: a} (keyed or positional)
	lit, ok := compositeLit(expr)
	if !ok {
		return nil, false
	}
	result := &v2RGBA{A: 255}
	positional := make([]int, 0, 4)
	for _, elt := range lit.Elts {
		if kv, ok := elt.(*ast.KeyValueExpr); ok {
			key, ok := kv.Key.(*ast.Ident)
			if !ok {
				return nil, false
			}
			n, ok := litInt(kv.Value)
			if !ok {
				return nil, false
			}
			switch key.Name {
			case "R":
				result.R = n
			case "G":
				result.G = n
			case "B":
				result.B = n
			case "A":
				result.A = n
			}
			continue
		}
		n, ok := litInt(elt)
		if !ok {
			return nil, false
		}
		positional = append(positional, n)
	}
	if len(positional) >= 3 {
		result.R, result.G, result.B = positional[0], positional[1], positional[2]
		if len(positional) == 4 {
			result.A = positional[3]
		}
	}
	return result, true
}

func litString(expr ast.Expr) (string, bool) {
	lit, ok := expr.(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return "", false
	}
	s, err := strconv.Unquote(lit.Value)
	if err != nil {
		return "", false
	}
	return s, true
}

func litInt(expr ast.Expr) (int, bool) {
	lit, ok := expr.(*ast.BasicLit)
	if !ok || lit.Kind != token.INT {
		return 0, false
	}
	n, err := strconv.Atoi(lit.Value)
	if err != nil {
		return 0, false
	}
	return n, true
}

func litBool(expr ast.Expr) (bool, bool) {
	ident, ok := expr.(*ast.Ident)
	if !ok {
		return false, false
	}
	switch ident.Name {
	case "true":
		return true, true
	case "false":
		return false, true
	}
	return false, false
}

func exprString(fset *token.FileSet, node ast.Node) string {
	var buf bytes.Buffer
	if err := printer.Fprint(&buf, fset, node); err != nil {
		return ""
	}
	return buf.String()
}

func isIdentifier(s string) bool {
	if s == "" {
		return false
	}
	for i, r := range s {
		if r == '_' || (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (i > 0 && r >= '0' && r <= '9') {
			continue
		}
		return false
	}
	return true
}

// generateMainV3 renders the v3 bootstrap file from the extracted v2
// information. The output is gofmt-formatted and parsed as a self-check.
func generateMainV3(info *v2MainInfo, cfg *v2ProjectConfig) ([]byte, error) {
	var b strings.Builder

	b.WriteString("// Generated by `wails3 migrate` from " + info.RunFile + ". Review before use.\n")
	b.WriteString("//\n")
	b.WriteString("// This file is the v3 equivalent of your v2 application bootstrap. It is\n")
	b.WriteString("// not part of the build (the .example extension keeps it out). To adopt it:\n")
	b.WriteString("//\n")
	b.WriteString("//  1. Migrate your bound structs to v3 services (see MIGRATION_REPORT.md).\n")
	b.WriteString("//  2. Replace the v2 bootstrap in main.go with this file's contents and\n")
	b.WriteString("//     remove the wails/v2 imports.\n")
	b.WriteString("//  3. Run: go get github.com/wailsapp/wails/v3@latest && go mod tidy\n")
	b.WriteString("package main\n\n")

	b.WriteString("import (\n")
	if info.AssetsIdent != "" {
		b.WriteString("\t\"embed\"\n")
	}
	b.WriteString("\t\"log\"\n\n")
	b.WriteString("\t\"github.com/wailsapp/wails/v3/pkg/application\"\n")
	b.WriteString(")\n\n")

	if info.AssetsIdent != "" {
		if info.EmbedDirective != "" {
			b.WriteString(info.EmbedDirective + "\n")
		} else {
			b.WriteString("// TODO(migrate): copy the //go:embed directive for `" + info.AssetsIdent + "` from your v2 code here.\n")
		}
		b.WriteString("var " + info.AssetsIdent + " embed.FS\n\n")
	}

	b.WriteString("func main() {\n")
	for _, stmt := range info.ConstructorStmts {
		b.WriteString("\t" + stmt + "\n")
	}
	if len(info.ConstructorStmts) > 0 {
		b.WriteString("\n")
	}

	appName := info.Window.Title
	if appName == "" {
		appName = cfg.Name
	}

	b.WriteString("\twailsApp := application.New(application.Options{\n")
	b.WriteString(fmt.Sprintf("\t\tName: %q,\n", appName))

	// Bind -> Services.
	b.WriteString("\t\tServices: []application.Service{\n")
	if len(info.BindExprs) == 0 {
		b.WriteString("\t\t\t// TODO(migrate): register your services here, e.g.\n")
		b.WriteString("\t\t\t//   application.NewService(&MyService{}),\n")
	}
	for _, expr := range info.BindExprs {
		b.WriteString("\t\t\tapplication.NewService(" + expr + "),\n")
	}
	b.WriteString("\t\t},\n")

	// Assets.
	if info.AssetsIdent != "" {
		b.WriteString("\t\tAssets: application.AssetOptions{\n")
		b.WriteString("\t\t\tHandler: application.AssetFileServerFS(" + info.AssetsIdent + "),\n")
		b.WriteString("\t\t},\n")
	} else {
		b.WriteString("\t\t// TODO(migrate): could not locate the embedded assets in your v2 code.\n")
		b.WriteString("\t\t// Configure the asset server like a fresh v3 project:\n")
		b.WriteString("\t\t//   Assets: application.AssetOptions{Handler: application.AssetFileServerFS(assets)},\n")
	}

	if info.SingleInstanceID != "" {
		b.WriteString("\t\tSingleInstance: &application.SingleInstanceOptions{\n")
		b.WriteString(fmt.Sprintf("\t\t\tUniqueID: %q,\n", info.SingleInstanceID))
		if info.SingleInstanceCB {
			b.WriteString("\t\t\t// TODO(migrate): port your v2 OnSecondInstanceLaunch callback:\n")
			b.WriteString("\t\t\t//   OnSecondInstanceLaunch: func(data application.SecondInstanceData) { ... },\n")
		}
		b.WriteString("\t\t},\n")
	}

	if info.OnShutdown != "" {
		b.WriteString("\t\t// TODO(migrate): v2 OnShutdown was `" + info.OnShutdown + "`.\n")
		b.WriteString("\t\t// v3 equivalent: OnShutdown: func() { ... } (no context parameter).\n")
	}
	b.WriteString("\t})\n\n")

	if info.OnStartup != "" {
		b.WriteString("\t// TODO(migrate): v2 OnStartup was `" + info.OnStartup + "`.\n")
		b.WriteString("\t// Move that logic into a service ServiceStartup(ctx context.Context,\n")
		b.WriteString("\t// options application.ServiceOptions) method, or subscribe to the\n")
		b.WriteString("\t// application start event:\n")
		b.WriteString("\t//   wailsApp.Event.OnApplicationEvent(events.Common.ApplicationStarted, func(*application.ApplicationEvent) { ... })\n")
	}
	if info.OnDomReady != "" {
		b.WriteString("\t// TODO(migrate): v2 OnDomReady was `" + info.OnDomReady + "`.\n")
		b.WriteString("\t// Subscribe to the window runtime-ready event in v3:\n")
		b.WriteString("\t//   window.OnWindowEvent(events.Common.WindowRuntimeReady, func(*application.WindowEvent) { ... })\n")
	}

	b.WriteString("\twailsApp.Window.NewWithOptions(application.WebviewWindowOptions{\n")
	b.WriteString(fmt.Sprintf("\t\tTitle: %q,\n", info.Window.Title))
	width, height := info.Window.Width, info.Window.Height
	if width == 0 {
		width = 1024 // the v2 default
	}
	if height == 0 {
		height = 768 // the v2 default
	}
	b.WriteString(fmt.Sprintf("\t\tWidth:  %d,\n", width))
	b.WriteString(fmt.Sprintf("\t\tHeight: %d,\n", height))
	b.WriteString("\t\tURL:    \"/\",\n")
	if info.Window.MinWidth != 0 {
		b.WriteString(fmt.Sprintf("\t\tMinWidth: %d,\n", info.Window.MinWidth))
	}
	if info.Window.MinHeight != 0 {
		b.WriteString(fmt.Sprintf("\t\tMinHeight: %d,\n", info.Window.MinHeight))
	}
	if info.Window.MaxWidth != 0 {
		b.WriteString(fmt.Sprintf("\t\tMaxWidth: %d,\n", info.Window.MaxWidth))
	}
	if info.Window.MaxHeight != 0 {
		b.WriteString(fmt.Sprintf("\t\tMaxHeight: %d,\n", info.Window.MaxHeight))
	}
	if info.Window.DisableResize {
		b.WriteString("\t\tDisableResize: true,\n")
	}
	if info.Window.Frameless {
		b.WriteString("\t\tFrameless: true,\n")
	}
	if info.Window.AlwaysOnTop {
		b.WriteString("\t\tAlwaysOnTop: true,\n")
	}
	if info.Window.Hidden {
		b.WriteString("\t\tHidden: true,\n")
	}
	if info.Window.StartState != "" {
		b.WriteString("\t\tStartState: application.WindowState" + info.Window.StartState + ",\n")
	}
	if bg := info.Window.Background; bg != nil {
		b.WriteString(fmt.Sprintf("\t\tBackgroundColour: application.NewRGBA(%d, %d, %d, %d),\n", bg.R, bg.G, bg.B, bg.A))
	}
	b.WriteString("\t})\n\n")

	b.WriteString("\tif err := wailsApp.Run(); err != nil {\n")
	b.WriteString("\t\tlog.Fatal(err)\n")
	b.WriteString("\t}\n")
	b.WriteString("}\n")

	formatted, err := format.Source([]byte(b.String()))
	if err != nil {
		return nil, fmt.Errorf("generated code failed to format: %w", err)
	}
	// Self-check: the generated file must parse.
	fset := token.NewFileSet()
	if _, err := parser.ParseFile(fset, generatedMainName, formatted, parser.ParseComments); err != nil {
		return nil, fmt.Errorf("generated code failed to parse: %w", err)
	}
	return formatted, nil
}

// migrateMain extracts the v2 wails.Run options and writes the v3 bootstrap
// example file. Extraction problems degrade to report entries.
func (m *migrator) migrateMain() {
	info, runtimeCalls, err := extractV2Main(m.projectDir)
	m.runtimeCalls = runtimeCalls
	if err != nil {
		m.mainManual = append(m.mainManual,
			fmt.Sprintf("%v. Follow the main.go section of the migration guide to write the v3 bootstrap by hand.", err))
		return
	}
	m.main = info
	m.mainManual = append(m.mainManual, info.Manual...)

	generated, err := generateMainV3(info, m.cfg)
	if err != nil {
		m.mainManual = append(m.mainManual, fmt.Sprintf("could not generate %s: %v", generatedMainName, err))
		return
	}
	if err := os.WriteFile(filepath.Join(m.projectDir, generatedMainName), generated, 0o644); err != nil {
		m.mainManual = append(m.mainManual, fmt.Sprintf("could not write %s: %v", generatedMainName, err))
		return
	}
	m.created = append(m.created, generatedMainName)
}
