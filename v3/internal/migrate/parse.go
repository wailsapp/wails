package migrate

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/mod/modfile"
)

// ParseV2Project reads a v2 project from dir. It is a syntax-only parse (no
// type checking), so it works without the v2 module being present in the
// module cache.
func ParseV2Project(dir string) (*V2Project, error) {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}

	cfg, err := LoadV2Config(absDir)
	if err != nil {
		return nil, err
	}

	proj := &V2Project{
		Dir:         absDir,
		Config:      cfg,
		FrontendDir: filepath.Join(absDir, cfg.FrontendDir),
		Report:      NewReport(),
	}

	// go.mod
	proj.GoModPath = filepath.Join(absDir, "go.mod")
	modData, err := os.ReadFile(proj.GoModPath)
	if err != nil {
		return nil, fmt.Errorf("could not read go.mod: %w", err)
	}
	mod, err := modfile.Parse("go.mod", modData, nil)
	if err != nil {
		return nil, fmt.Errorf("could not parse go.mod: %w", err)
	}
	if mod.Module == nil {
		return nil, fmt.Errorf("go.mod has no module directive")
	}
	proj.ModulePath = mod.Module.Mod.Path
	if !requiresWailsV2(mod) {
		return nil, fmt.Errorf("go.mod does not require github.com/wailsapp/wails/v2 - is this a Wails v2 project?")
	}

	// Parse all Go files in the module (excluding frontend, build dir and
	// hidden/vendor directories).
	fset := token.NewFileSet()
	files := map[string]*ast.File{} // abs path -> file
	err = filepath.WalkDir(absDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			name := d.Name()
			if path != absDir && (strings.HasPrefix(name, ".") || name == "vendor" || name == "node_modules") {
				return filepath.SkipDir
			}
			if path == proj.FrontendDir || path == filepath.Join(absDir, cfg.BuildDir) {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		file, perr := parser.ParseFile(fset, path, nil, parser.ParseComments|parser.SkipObjectResolution)
		if perr != nil {
			return fmt.Errorf("could not parse %s: %w", path, perr)
		}
		files[path] = file
		return nil
	})
	if err != nil {
		return nil, err
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("no Go files found in %s", absDir)
	}

	// Locate the wails.Run call.
	for path, file := range files {
		info := findRunCall(fset, path, file)
		if info == nil {
			continue
		}
		if proj.Main != nil {
			return nil, fmt.Errorf("found more than one wails.Run call (%s and %s); cannot migrate automatically", proj.Main.Path, path)
		}
		proj.Main = info
	}
	if proj.Main == nil {
		return nil, fmt.Errorf("could not find a wails.Run(&options.App{...}) call in %s", absDir)
	}
	proj.Main.Source, err = os.ReadFile(proj.Main.Path)
	if err != nil {
		return nil, err
	}

	for path := range files {
		if path != proj.Main.Path {
			proj.GoFiles = append(proj.GoFiles, path)
		}
	}

	// Resolve bound types from the Bind field, if present.
	if proj.Main.AppLit != nil {
		if bind := fieldValue(proj.Main.AppLit, "Bind"); bind != nil {
			proj.BoundTypes = resolveBoundTypes(fset, files, proj, bind)
		}
	}

	return proj, nil
}

func requiresWailsV2(mod *modfile.File) bool {
	for _, req := range mod.Require {
		if req.Mod.Path == "github.com/wailsapp/wails/v2" {
			return true
		}
	}
	return false
}

// findRunCall looks for a statement of the form
//
//	err := wails.Run(&options.App{...})
//	err = wails.Run(...)
//	wails.Run(...)
//	if err := wails.Run(...); err != nil { ... }
//
// where "wails" is the local name of the github.com/wailsapp/wails/v2 import.
func findRunCall(fset *token.FileSet, path string, file *ast.File) *MainInfo {
	imports := importMap(file)
	wailsName := ""
	for name, ipath := range imports {
		if ipath == "github.com/wailsapp/wails/v2" {
			wailsName = name
		}
	}
	if wailsName == "" {
		return nil
	}

	info := &MainInfo{Path: path, File: file, Fset: fset, Imports: imports}

	isRunCall := func(n ast.Node) *ast.CallExpr {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return nil
		}
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok || sel.Sel.Name != "Run" {
			return nil
		}
		ident, ok := sel.X.(*ast.Ident)
		if !ok || ident.Name != wailsName {
			return nil
		}
		return call
	}

	ast.Inspect(file, func(n ast.Node) bool {
		if info.RunCall != nil {
			return false
		}
		stmt, ok := n.(ast.Stmt)
		if !ok {
			return true
		}
		switch s := stmt.(type) {
		case *ast.AssignStmt:
			if len(s.Rhs) != 1 {
				return true
			}
			call := isRunCall(s.Rhs[0])
			if call == nil {
				return true
			}
			info.RunStmt = s
			info.RunCall = call
			if len(s.Lhs) == 1 {
				if ident, ok := s.Lhs[0].(*ast.Ident); ok {
					info.ErrIdent = ident.Name
					info.AssignTok = s.Tok
				}
			}
			return false
		case *ast.ExprStmt:
			call := isRunCall(s.X)
			if call == nil {
				return true
			}
			info.RunStmt = s
			info.RunCall = call
			return false
		}
		return true
	})

	if info.RunCall == nil {
		return nil
	}

	// Extract the &options.App{...} literal.
	if len(info.RunCall.Args) == 1 {
		arg := info.RunCall.Args[0]
		if unary, ok := arg.(*ast.UnaryExpr); ok && unary.Op == token.AND {
			arg = unary.X
		}
		if lit, ok := arg.(*ast.CompositeLit); ok {
			info.AppLit = lit
		}
	}
	return info
}

func importMap(file *ast.File) map[string]string {
	m := map[string]string{}
	for _, imp := range file.Imports {
		path, err := strconv.Unquote(imp.Path.Value)
		if err != nil {
			continue
		}
		name := ""
		if imp.Name != nil {
			name = imp.Name.Name
		} else {
			// Default local name: last path element, with major-version
			// suffixes (vN) resolved to the element before them. This is a
			// heuristic (the true name is the package clause), but it is
			// correct for the packages the migrator cares about.
			name = path[strings.LastIndex(path, "/")+1:]
			if len(name) > 1 && name[0] == 'v' {
				if _, err := strconv.Atoi(name[1:]); err == nil {
					trimmed := path[:strings.LastIndex(path, "/")]
					name = trimmed[strings.LastIndex(trimmed, "/")+1:]
				}
			}
		}
		if name == "_" || name == "." {
			continue
		}
		m[name] = path
	}
	return m
}

// fieldValue returns the value of the named field in a composite literal, or
// nil if not present.
func fieldValue(lit *ast.CompositeLit, name string) ast.Expr {
	for _, elt := range lit.Elts {
		kv, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			continue
		}
		key, ok := kv.Key.(*ast.Ident)
		if !ok {
			continue
		}
		if key.Name == name {
			return kv.Value
		}
	}
	return nil
}

// exprText returns the original source text for an expression.
func exprText(fset *token.FileSet, src []byte, node ast.Node) string {
	start := fset.Position(node.Pos()).Offset
	end := fset.Position(node.End()).Offset
	if start < 0 || end > len(src) || start >= end {
		return ""
	}
	return string(src[start:end])
}

// printExpr renders an AST expression to Go source (used where original
// source bytes are not to hand).
func printExpr(fset *token.FileSet, node ast.Node) string {
	var sb strings.Builder
	_ = printer.Fprint(&sb, fset, node)
	return sb.String()
}

// resolveBoundTypes maps each element of the Bind slice literal to its struct
// type and collects the exported methods of that type from the parsed files.
func resolveBoundTypes(fset *token.FileSet, files map[string]*ast.File, proj *V2Project, bind ast.Expr) []*BoundType {
	lit, ok := bind.(*ast.CompositeLit)
	if !ok {
		proj.Report.Manual("Bind", "The Bind value is not a slice literal, so bound structs could not be discovered. The generated frontend/wailsjs shims may be incomplete; check frontend imports against the v3 bindings generated into frontend/bindings.")
		return nil
	}

	mainFile := files[proj.Main.Path]
	var result []*BoundType
	for _, elt := range lit.Elts {
		expr := printExpr(fset, elt)
		typeName := resolveElementType(fset, files, mainFile, elt)
		bt := &BoundType{Expr: expr, Name: typeName}
		if typeName == "" {
			proj.Report.Manual("Bind: "+expr,
				"Could not statically determine the struct type of this Bind entry. It is still registered as a v3 service, but no frontend/wailsjs shim was generated for it.")
			result = append(result, bt)
			continue
		}
		pkgName, pkgPath, methods := collectMethods(fset, files, proj, typeName)
		bt.PkgName = pkgName
		bt.PkgPath = pkgPath
		bt.Methods = methods
		if len(methods) == 0 {
			proj.Report.Note("Bind: " + expr + " (" + typeName + ") has no exported methods that could be discovered; no frontend/wailsjs shim was generated for it.")
		}
		result = append(result, bt)
	}
	return result
}

// resolveElementType attempts to resolve a Bind element expression to a
// struct type name using purely syntactic information:
//
//	&App{...}          -> App
//	NewApp()           -> return type of func NewApp
//	app                -> declaration of app in the same file (app := NewApp(),
//	                      app := &App{}, var app = ...)
func resolveElementType(fset *token.FileSet, files map[string]*ast.File, mainFile *ast.File, expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.UnaryExpr:
		if e.Op == token.AND {
			if lit, ok := e.X.(*ast.CompositeLit); ok {
				if ident, ok := lit.Type.(*ast.Ident); ok {
					return ident.Name
				}
			}
		}
	case *ast.CallExpr:
		if ident, ok := e.Fun.(*ast.Ident); ok {
			return constructorReturnType(files, ident.Name)
		}
	case *ast.Ident:
		var typeName string
		ast.Inspect(mainFile, func(n ast.Node) bool {
			if typeName != "" {
				return false
			}
			assign, ok := n.(*ast.AssignStmt)
			if !ok || len(assign.Lhs) != 1 || len(assign.Rhs) != 1 {
				return true
			}
			lhs, ok := assign.Lhs[0].(*ast.Ident)
			if !ok || lhs.Name != e.Name {
				return true
			}
			typeName = resolveElementType(fset, files, mainFile, assign.Rhs[0])
			return false
		})
		return typeName
	}
	return ""
}

// constructorReturnType finds `func Name(...) *T` in the parsed files and
// returns T.
func constructorReturnType(files map[string]*ast.File, name string) string {
	for _, file := range files {
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Recv != nil || fn.Name.Name != name {
				continue
			}
			if fn.Type.Results == nil || len(fn.Type.Results.List) == 0 {
				return ""
			}
			t := fn.Type.Results.List[0].Type
			if star, ok := t.(*ast.StarExpr); ok {
				t = star.X
			}
			if ident, ok := t.(*ast.Ident); ok {
				return ident.Name
			}
		}
	}
	return ""
}

// collectMethods gathers the exported methods declared on *typeName or
// typeName across the parsed files.
func collectMethods(fset *token.FileSet, files map[string]*ast.File, proj *V2Project, typeName string) (pkgName, pkgPath string, methods []*BoundMethod) {
	// Sort file paths for deterministic method order across map iteration.
	paths := make([]string, 0, len(files))
	for path := range files {
		paths = append(paths, path)
	}
	sort.Strings(paths)

	for _, path := range paths {
		file := files[path]
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Recv == nil || len(fn.Recv.List) != 1 {
				continue
			}
			recv := fn.Recv.List[0].Type
			if star, ok := recv.(*ast.StarExpr); ok {
				recv = star.X
			}
			ident, ok := recv.(*ast.Ident)
			if !ok || ident.Name != typeName {
				continue
			}
			if !fn.Name.IsExported() {
				continue
			}
			if pkgName == "" {
				pkgName = file.Name.Name
				pkgPath = packagePath(proj, path, pkgName)
			}
			methods = append(methods, &BoundMethod{
				Name:    fn.Name.Name,
				Params:  fieldListParams(fset, fn.Type.Params),
				Results: fieldListParams(fset, fn.Type.Results),
			})
		}
	}
	return pkgName, pkgPath, methods
}

// packagePath computes the binding FQN package path for a file: "main" for
// the main package, otherwise modulePath[/reldir].
func packagePath(proj *V2Project, filePath, pkgName string) string {
	if pkgName == "main" {
		return "main"
	}
	rel, err := filepath.Rel(proj.Dir, filepath.Dir(filePath))
	if err != nil || rel == "." {
		return proj.ModulePath
	}
	return proj.ModulePath + "/" + filepath.ToSlash(rel)
}

func fieldListParams(fset *token.FileSet, list *ast.FieldList) []Param {
	if list == nil {
		return nil
	}
	var params []Param
	for _, field := range list.List {
		goType := printExpr(fset, field.Type)
		tsType := goTypeToTS(goType)
		if len(field.Names) == 0 {
			params = append(params, Param{GoType: goType, TSType: tsType})
			continue
		}
		for _, name := range field.Names {
			params = append(params, Param{Name: name.Name, GoType: goType, TSType: tsType})
		}
	}
	return params
}

// goTypeToTS maps a printed Go type to a best-effort TypeScript type for the
// generated .d.ts shims. Anything unrecognised becomes "any"; the real v3
// bindings (frontend/bindings) carry full model types.
func goTypeToTS(goType string) string {
	goType = strings.TrimPrefix(goType, "*")
	switch goType {
	case "string":
		return "string"
	case "bool":
		return "boolean"
	case "int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64",
		"float32", "float64", "byte", "rune", "uintptr":
		return "number"
	case "interface{}", "any":
		return "any"
	case "error":
		return "void"
	}
	if strings.HasPrefix(goType, "[]") {
		return goTypeToTS(goType[2:]) + "[]"
	}
	if strings.HasPrefix(goType, "map[") {
		return "Record<string, any>"
	}
	return "any"
}
