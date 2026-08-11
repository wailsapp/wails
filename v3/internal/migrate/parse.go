package migrate

import (
	"fmt"
	"go/ast"
	"go/build"
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
	buildContext := build.Default
	files := map[string]*ast.File{} // abs path -> file
	err = filepath.WalkDir(absDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			name := d.Name()
			if path != absDir && (strings.HasPrefix(name, ".") || name == "vendor" || name == "node_modules" || name == "testdata") {
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
		matched, merr := buildContext.MatchFile(filepath.Dir(path), d.Name())
		if merr != nil {
			return fmt.Errorf("could not evaluate build constraints for %s: %w", path, merr)
		}
		if !matched {
			return nil
		}
		file, perr := parser.ParseFile(fset, path, nil, parser.ParseComments|parser.SkipObjectResolution)
		if perr != nil {
			return fmt.Errorf("could not parse %s: %w", path, perr)
		}
		for _, imp := range file.Imports {
			if imp.Path.Value == strconv.Quote(V2RuntimeImport) {
				proj.UsesV2Runtime = true
			}
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
		candidates, ferr := findRunCalls(fset, path, file)
		if ferr != nil {
			return nil, ferr
		}
		for _, info := range candidates {
			if proj.Main != nil {
				return nil, fmt.Errorf("found more than one wails.Run call (%s and %s); cannot migrate automatically", proj.Main.Path, path)
			}
			proj.Main = info
		}
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
	sort.Strings(proj.GoFiles)

	// Resolve bound types from the Bind field, if present.
	if proj.Main.AppLit != nil {
		if bind := fieldValue(proj.Main.AppLit, "Bind"); bind != nil {
			proj.BoundTypes = resolveBoundTypes(fset, files, proj, bind)
		}
	}

	// Enumerate every v2 API call site with its v3 replacement.
	adviseGoRuntimeCalls(fset, files, proj)
	if err := adviseFrontendImports(proj); err != nil {
		return nil, err
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

// findRunCalls returns every rewriteable wails.Run call in a file. Calls
// must be direct expression statements, assignments, or the supported
// `if err := wails.Run(...); err != nil { ... }` form so GenerateMain can
// replace a complete statement without leaving invalid Go behind.
func findRunCalls(fset *token.FileSet, path string, file *ast.File) ([]*MainInfo, error) {
	imports := importMap(file)
	wailsName := ""
	for name, ipath := range imports {
		if ipath == "github.com/wailsapp/wails/v2" {
			wailsName = name
		}
	}
	if wailsName == "" {
		return nil, nil
	}

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

	isDirectAssignment := func(stmt *ast.AssignStmt, call *ast.CallExpr) bool {
		return len(stmt.Lhs) == 1 && len(stmt.Rhs) == 1 && stmt.Rhs[0] == call
	}
	isErrorCondition := func(ifStmt *ast.IfStmt, errName string) bool {
		cond, ok := ifStmt.Cond.(*ast.BinaryExpr)
		if !ok || cond.Op != token.NEQ {
			return false
		}
		matches := func(left, right ast.Expr) bool {
			ident, lok := left.(*ast.Ident)
			nilIdent, rok := right.(*ast.Ident)
			return lok && rok && ident.Name == errName && nilIdent.Name == "nil"
		}
		return matches(cond.X, cond.Y) || matches(cond.Y, cond.X)
	}

	var candidates []*MainInfo
	var stack []ast.Node
	var walkErr error
	ast.Inspect(file, func(n ast.Node) bool {
		if walkErr != nil {
			return false
		}
		if n == nil {
			stack = stack[:len(stack)-1]
			return false
		}
		call := isRunCall(n)
		if call != nil {
			var runStmt ast.Stmt
			var runIf *ast.IfStmt
			errIdent := ""
			assignTok := token.ILLEGAL

			for i := len(stack) - 1; i >= 0; i-- {
				switch parent := stack[i].(type) {
				case *ast.IfStmt:
					assign, ok := parent.Init.(*ast.AssignStmt)
					if !ok || !isDirectAssignment(assign, call) {
						continue
					}
					ident, ok := assign.Lhs[0].(*ast.Ident)
					if !ok || !isErrorCondition(parent, ident.Name) {
						walkErr = fmt.Errorf("unsupported wails.Run context at %s:%d; expected `if err := wails.Run(...); err != nil`", path, fset.Position(call.Pos()).Line)
						return false
					}
					runStmt = parent
					runIf = parent
					errIdent = ident.Name
					assignTok = assign.Tok
				case *ast.AssignStmt:
					if i > 0 {
						if parentIf, ok := stack[i-1].(*ast.IfStmt); ok && parentIf.Init == parent {
							continue
						}
					}
					if !isDirectAssignment(parent, call) {
						continue
					}
					ident, ok := parent.Lhs[0].(*ast.Ident)
					if !ok {
						walkErr = fmt.Errorf("unsupported wails.Run assignment at %s:%d", path, fset.Position(call.Pos()).Line)
						return false
					}
					runStmt = parent
					errIdent = ident.Name
					assignTok = parent.Tok
				case *ast.ExprStmt:
					if parent.X != call {
						continue
					}
					runStmt = parent
				}
				if runStmt != nil {
					break
				}
			}

			if runStmt == nil {
				walkErr = fmt.Errorf("unsupported wails.Run context at %s:%d; use a direct call, assignment, or `if err := wails.Run(...); err != nil`", path, fset.Position(call.Pos()).Line)
				return false
			}
			info := &MainInfo{
				Path: path, File: file, Fset: fset, Imports: imports,
				RunStmt: runStmt, RunIf: runIf, RunCall: call,
				ErrIdent: errIdent, AssignTok: assignTok,
			}
			if len(call.Args) == 1 {
				arg := call.Args[0]
				if unary, ok := arg.(*ast.UnaryExpr); ok && unary.Op == token.AND {
					arg = unary.X
				}
				if lit, ok := arg.(*ast.CompositeLit); ok {
					info.AppLit = lit
				}
			}
			candidates = append(candidates, info)
		}
		stack = append(stack, n)
		return true
	})
	if walkErr != nil {
		return nil, walkErr
	}
	return candidates, nil
}

// findRunCall is kept as a small compatibility wrapper for package-local
// callers that only need the first candidate.
func findRunCall(fset *token.FileSet, path string, file *ast.File) *MainInfo {
	candidates, err := findRunCalls(fset, path, file)
	if err != nil || len(candidates) == 0 {
		return nil
	}
	return candidates[0]
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
			name = defaultImportName(path)
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

	var result []*BoundType
	for _, elt := range lit.Elts {
		expr := printExpr(fset, elt)
		resolved := resolveElementInfo(fset, files, proj, proj.Main, elt)
		bt := &BoundType{Expr: expr, Name: resolved.Name, PkgName: resolved.PkgName}
		if resolved.Name == "" {
			proj.Report.Manual("Bind: "+expr,
				"Could not statically determine the struct type of this Bind entry. It is still registered as a v3 service, but no frontend/wailsjs shim was generated for it.")
			result = append(result, bt)
			continue
		}
		bt.PkgPath = resolved.PkgPath
		bt.Methods = collectMethods(fset, files, proj, resolved.PkgDir, resolved.PkgName, resolved.Name)
		if len(bt.Methods) == 0 {
			proj.Report.Note("Bind: " + expr + " (" + resolved.Name + ") has no exported methods that could be discovered; no frontend/wailsjs shim was generated for it.")
		}
		result = append(result, bt)
	}
	return result
}

type resolvedElement struct {
	Name    string
	PkgName string
	PkgDir  string
	PkgPath string
}

func resolveElementInfo(fset *token.FileSet, files map[string]*ast.File, proj *V2Project, main *MainInfo, expr ast.Expr) resolvedElement {
	return resolveElementInfoAt(fset, files, proj, main.File, main.RunCall.Pos(), expr, map[string]bool{})
}

func resolveElementInfoAt(fset *token.FileSet, files map[string]*ast.File, proj *V2Project, mainFile *ast.File, target token.Pos, expr ast.Expr, seen map[string]bool) resolvedElement {
	packageName := mainFile.Name.Name
	packageDir := filepath.Dir(proj.Main.Path)
	pkgPath := packagePath(proj, proj.Main.Path, packageName)
	fromPackage := resolvedElement{PkgName: packageName, PkgDir: packageDir, PkgPath: pkgPath}
	imports := importMap(mainFile)

	resolveType := func(name string) resolvedElement {
		return resolvedElement{Name: name, PkgName: packageName, PkgDir: packageDir, PkgPath: pkgPath}
	}
	switch e := expr.(type) {
	case *ast.UnaryExpr:
		if e.Op == token.AND {
			if lit, ok := e.X.(*ast.CompositeLit); ok {
				if ident, ok := lit.Type.(*ast.Ident); ok {
					return resolveType(ident.Name)
				}
			}
		}
	case *ast.CompositeLit:
		if ident, ok := e.Type.(*ast.Ident); ok {
			return resolveType(ident.Name)
		}
	case *ast.CallExpr:
		switch fun := e.Fun.(type) {
		case *ast.Ident:
			return constructorReturnTypeInPackage(files, packageDir, packageName, fun.Name, resolveType)
		case *ast.SelectorExpr:
			if pkgIdent, ok := fun.X.(*ast.Ident); ok {
				if importedPath, ok := imports[pkgIdent.Name]; ok {
					for path, file := range files {
						if packagePath(proj, path, file.Name.Name) != importedPath {
							continue
						}
						return constructorReturnTypeInPackage(files, filepath.Dir(path), file.Name.Name, fun.Sel.Name, func(name string) resolvedElement {
							return resolvedElement{Name: name, PkgName: file.Name.Name, PkgDir: filepath.Dir(path), PkgPath: importedPath}
						})
					}
				}
			}
		}
	case *ast.Ident:
		if seen[e.Name] {
			return resolvedElement{}
		}
		seen[e.Name] = true
		declaration := findLexicalValue(mainFile, e.Name, target)
		if declaration == nil {
			declaration = findPackageValue(files, packageDir, packageName, e.Name)
		}
		if declaration != nil {
			return resolveElementInfoAt(fset, files, proj, mainFile, target, declaration, seen)
		}
		return resolvedElement{}
	}
	return fromPackage
}

func constructorReturnTypeInPackage(files map[string]*ast.File, packageDir, packageName, name string, makeType func(string) resolvedElement) resolvedElement {
	paths := sortedFilePaths(files)
	for _, path := range paths {
		file := files[path]
		if filepath.Dir(path) != packageDir || file.Name.Name != packageName {
			continue
		}
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Recv != nil || fn.Name.Name != name || fn.Type.Results == nil || len(fn.Type.Results.List) == 0 {
				continue
			}
			if resultName := typeName(fn.Type.Results.List[0].Type); resultName != "" {
				return makeType(resultName)
			}
		}
	}
	return resolvedElement{}
}

func typeName(expr ast.Expr) string {
	if star, ok := expr.(*ast.StarExpr); ok {
		expr = star.X
	}
	switch value := expr.(type) {
	case *ast.Ident:
		return value.Name
	case *ast.SelectorExpr:
		return value.Sel.Name
	}
	return ""
}

func sortedFilePaths(files map[string]*ast.File) []string {
	paths := make([]string, 0, len(files))
	for path := range files {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	return paths
}

func findLexicalValue(file *ast.File, name string, target token.Pos) ast.Expr {
	var blocks []*ast.BlockStmt
	ast.Inspect(file, func(node ast.Node) bool {
		if block, ok := node.(*ast.BlockStmt); ok && block.Pos() <= target && target <= block.End() {
			blocks = append(blocks, block)
		}
		return true
	})
	for i := len(blocks) - 1; i >= 0; i-- {
		if value := findValueInBlock(blocks[i], name, target); value != nil {
			return value
		}
	}
	if fn := enclosingFunc(file, target); fn != nil {
		for _, field := range fn.Type.Params.List {
			for _, param := range field.Names {
				if param.Name == name {
					return nil
				}
			}
		}
	}
	return nil
}

func enclosingFunc(file *ast.File, target token.Pos) *ast.FuncDecl {
	var result *ast.FuncDecl
	ast.Inspect(file, func(node ast.Node) bool {
		fn, ok := node.(*ast.FuncDecl)
		if ok && fn.Body.Pos() <= target && target <= fn.Body.End() {
			result = fn
		}
		return true
	})
	return result
}

func findValueInBlock(block *ast.BlockStmt, name string, target token.Pos) ast.Expr {
	var result ast.Expr
	for _, stmt := range block.List {
		if stmt.Pos() >= target {
			break
		}
		switch statement := stmt.(type) {
		case *ast.DeclStmt:
			if decl, ok := statement.Decl.(*ast.GenDecl); ok && decl.Tok == token.VAR {
				for _, spec := range decl.Specs {
					valueSpec := spec.(*ast.ValueSpec)
					if value := valueForName(valueSpec, name); value != nil {
						result = value
					}
				}
			}
		case *ast.AssignStmt:
			if statement.Tok != token.DEFINE {
				continue
			}
			for i, lhs := range statement.Lhs {
				ident, ok := lhs.(*ast.Ident)
				if ok && ident.Name == name && i < len(statement.Rhs) {
					result = statement.Rhs[i]
				} else if ok && ident.Name == name && len(statement.Rhs) == 1 {
					result = statement.Rhs[0]
				}
			}
		}
	}
	return result
}

func findPackageValue(files map[string]*ast.File, packageDir, packageName, name string) ast.Expr {
	for _, path := range sortedFilePaths(files) {
		file := files[path]
		if filepath.Dir(path) != packageDir || file.Name.Name != packageName {
			continue
		}
		for _, decl := range file.Decls {
			gen, ok := decl.(*ast.GenDecl)
			if !ok || gen.Tok != token.VAR {
				continue
			}
			for _, spec := range gen.Specs {
				if value := valueForName(spec.(*ast.ValueSpec), name); value != nil {
					return value
				}
			}
		}
	}
	return nil
}

func valueForName(spec *ast.ValueSpec, name string) ast.Expr {
	for i, ident := range spec.Names {
		if ident.Name != name {
			continue
		}
		if i < len(spec.Values) {
			return spec.Values[i]
		}
		if len(spec.Values) == 1 {
			return spec.Values[0]
		}
	}
	return nil
}

// resolveElementType attempts to resolve a Bind element expression to a
// struct type name using purely syntactic information:
//
//	&App{...}          -> App
//	NewApp()           -> return type of func NewApp
//	app                -> declaration of app in the same file (app := NewApp(),
//	                      app := &App{}, var app = ...)
func resolveElementType(fset *token.FileSet, files map[string]*ast.File, mainFile *ast.File, expr ast.Expr) string {
	return resolveElementInfoAt(fset, files, &V2Project{Main: &MainInfo{Path: "main.go"}}, mainFile, expr.Pos(), expr, map[string]bool{}).Name
}

// constructorReturnType finds `func Name(...) *T` in the parsed files and
// returns T.
func constructorReturnType(files map[string]*ast.File, name string) string {
	for _, path := range sortedFilePaths(files) {
		file := files[path]
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
func collectMethods(fset *token.FileSet, files map[string]*ast.File, proj *V2Project, packageDir, packageName, typeName string) (methods []*BoundMethod) {
	for _, path := range sortedFilePaths(files) {
		file := files[path]
		if filepath.Dir(path) != packageDir || file.Name.Name != packageName {
			continue
		}
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
			methods = append(methods, &BoundMethod{
				Name:    fn.Name.Name,
				Params:  fieldListParams(fset, fn.Type.Params),
				Results: fieldListParams(fset, fn.Type.Results),
			})
		}
	}
	return methods
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
