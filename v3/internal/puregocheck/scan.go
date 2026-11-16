package main

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

var (
	cIdent   = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)
	selector = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*(:([A-Za-z_][A-Za-z0-9_]*:)*)?$`)
	libName  = regexp.MustCompile(`^[^\s"]+(\.so(\.[0-9]+)*|\.dylib|\.framework/[A-Za-z]+)$`)
)

type parsed struct {
	path string
	file *ast.File
}

// buildFiles returns the non-test files of dir that `go build` selects for
// goos and tags, with cgo on or off.
func buildFiles(fset *token.FileSet, dir, goos string, tags []string, cgo bool) ([]parsed, error) {
	ctx := build.Default
	ctx.GOOS = goos
	ctx.GOARCH = "arm64"
	ctx.CgoEnabled = cgo
	ctx.BuildTags = tags
	names, err := filepath.Glob(filepath.Join(dir, "*.go"))
	if err != nil {
		return nil, err
	}
	var out []parsed
	for _, path := range names {
		name := filepath.Base(path)
		if strings.HasSuffix(name, "_test.go") {
			continue
		}
		if ok, err := ctx.MatchFile(dir, name); err != nil || !ok {
			continue
		}
		f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			return nil, err
		}
		usesC := slices.ContainsFunc(f.Imports, func(s *ast.ImportSpec) bool { return s.Path.Value == `"C"` })
		if usesC && !cgo {
			continue
		}
		out = append(out, parsed{path, f})
	}
	return out, nil
}

// cgoTags is the configuration of the cgo reference backend that a purego
// build is compared with.
func cgoTags(goos string, tags []string) []string {
	var out []string
	for _, t := range tags {
		if t != "purego" {
			out = append(out, t)
		}
	}
	return out
}

func scan(root, goos string, tags []string) (*report, error) {
	fset := token.NewFileSet()
	r := &report{GOOS: goos}
	symbols := map[string]*entry{}
	classes := map[string]*entry{}
	selectors := map[string]*entry{}
	optionalSyms := manifestOptional(root)
	var all []parsed
	var imp types.Importer
	if goos == "darwin" {
		imp = stackImporter(fset, tags)
	}

	for _, pkg := range packages {
		dir := filepath.Join(root, pkg)
		if _, err := os.Stat(dir); err != nil {
			continue
		}
		pure, err := buildFiles(fset, dir, goos, tags, false)
		if err != nil {
			return nil, err
		}
		ref, err := buildFiles(fset, dir, goos, cgoTags(goos, tags), true)
		if err != nil {
			return nil, err
		}
		r.Files += len(pure)
		for _, p := range pure {
			rel, _ := filepath.Rel(root, p.path)
			c := &collector{fset: fset, rel: rel, symbols: symbols, classes: classes, selectors: selectors, optional: optionalSyms, r: r}
			c.file(p.file)
		}
		if goos == "darwin" {
			r.PoolThread = append(r.PoolThread, findPoolThreads(fset, root, pure)...)
			var asts []*ast.File
			for _, p := range pure {
				asts = append(asts, p.file)
			}
			sr := checkStackPackage(fset, imp, root, asts)
			r.StackNarrow = append(r.StackNarrow, sr.findings...)
			r.StackChecked += sr.checked
			r.StackBundled += sr.bundled
			r.StackUnknown += sr.unknown
			r.StackForwarded += sr.forwarded
			r.StackSkipped = append(r.StackSkipped, sr.skipped...)
		}
		all = append(all, pure...)
		r.Stubs = append(r.Stubs, findStubs(fset, root, pure, ref)...)
	}
	// Guards are checked once every package's registrations are known.
	for _, p := range all {
		rel, _ := filepath.Rel(root, p.path)
		c := &collector{fset: fset, rel: rel, symbols: symbols, r: r}
		c.checkGuards(p.file)
	}
	r.Symbols = sortedEntries(symbols)
	if goos == "darwin" {
		r.Classes = sortedEntries(classes)
		r.Selectors = sortedEntries(selectors)
	}
	return r, nil
}

// manifestOptional reads the puregogen manifests for functions marked optional.
func manifestOptional(root string) map[string]bool {
	out := map[string]bool{}
	paths, _ := filepath.Glob(filepath.Join(root, "internal/puregogen/*.json"))
	for _, p := range paths {
		data, err := os.ReadFile(p)
		if err != nil {
			continue
		}
		var m struct {
			Functions []struct {
				Name     string `json:"name"`
				Optional bool   `json:"optional"`
			} `json:"functions"`
		}
		if json.Unmarshal(data, &m) != nil {
			continue
		}
		for _, f := range m.Functions {
			if f.Optional {
				out[f.Name] = true
			}
		}
	}
	return out
}

type collector struct {
	fset      *token.FileSet
	rel       string
	symbols   map[string]*entry
	classes   map[string]*entry
	selectors map[string]*entry
	optional  map[string]bool
	r         *report
	// inside a function that defines Objective-C methods (MethodDef tables):
	// string literals in composite literals there name our own selectors.
	defining bool
	// selectors the enclosing function checks with respondsTo: a send of one
	// of these is guarded, so it is optional.
	guards map[string]bool
}

func (c *collector) where(n ast.Node) string {
	return fmt.Sprintf("%s:%d", c.rel, c.fset.Position(n.Pos()).Line)
}

func add(m map[string]*entry, name, where string, optional, ours bool) {
	e := m[name]
	if e == nil {
		e = &entry{Name: name, Optional: true, Ours: false}
		m[name] = e
	}
	// One unguarded use makes it required; one definition by the backend
	// makes it ours.
	if !optional {
		e.Optional = false
	}
	if ours {
		e.Ours = true
	}
	if len(e.Where) < 8 && !slices.Contains(e.Where, where) {
		e.Where = append(e.Where, where)
	}
}

func lit(e ast.Expr) (string, bool) {
	b, ok := e.(*ast.BasicLit)
	if !ok || b.Kind != token.STRING {
		return "", false
	}
	s, err := strconv.Unquote(b.Value)
	return s, err == nil
}

func funName(e ast.Expr) string {
	switch f := e.(type) {
	case *ast.Ident:
		return f.Name
	case *ast.SelectorExpr:
		return f.Sel.Name
	case *ast.IndexExpr:
		return funName(f.X)
	case *ast.IndexListExpr:
		return funName(f.X)
	}
	return ""
}

func mentions(n ast.Node, names ...string) bool {
	found := false
	ast.Inspect(n, func(x ast.Node) bool {
		if id, ok := x.(*ast.Ident); ok && slices.Contains(names, id.Name) {
			found = true
		}
		return !found
	})
	return found
}

func (c *collector) file(f *ast.File) {
	for _, d := range f.Decls {
		fd, ok := d.(*ast.FuncDecl)
		c.defining = ok && mentions(fd, "MethodDef", "lifecycleEvent")
		c.guards = nil
		if ok && fd.Body != nil {
			c.guards = respondedSelectors(fd.Body)
		}
		ast.Inspect(d, c.node)
	}
}

func (c *collector) node(n ast.Node) bool {
	switch x := n.(type) {
	case *ast.BasicLit:
		// Every library the backend may dlopen; resolve loads the same set.
		if s, ok := lit(x); ok && libName.MatchString(s) && !slices.Contains(c.r.libCandidates, s) {
			c.r.libCandidates = append(c.r.libCandidates, s)
		}
	case *ast.CompositeLit:
		c.composite(x)
	case *ast.CallExpr:
		c.call(x)
	}
	return true
}

func (c *collector) composite(x *ast.CompositeLit) {
	// objc.MethodDef{...}, or an element of []objc.MethodDef{{Cmd: ...}} whose
	// type is elided.
	isMethodDef := x.Type != nil && funName(x.Type) == "MethodDef"
	if x.Type == nil {
		for _, el := range x.Elts {
			if kv, ok := el.(*ast.KeyValueExpr); ok && fmt.Sprint(kv.Key) == "Cmd" {
				isMethodDef = true
			}
		}
	}
	for _, el := range x.Elts {
		// objc.MethodDef{Cmd: sel_("windowWillClose:"), ...}: a method the
		// backend implements on its own registered class.
		if kv, ok := el.(*ast.KeyValueExpr); ok && isMethodDef && fmt.Sprint(kv.Key) == "Cmd" {
			if call, ok := kv.Value.(*ast.CallExpr); ok && len(call.Args) == 1 {
				if s, ok := lit(call.Args[0]); ok {
					add(c.selectors, s, c.where(kv), false, true)
				}
			}
		}
		kv, ok := el.(*ast.KeyValueExpr)
		if ok {
			// map[string]any{"g_strfreev": &fn}: a symbol table.
			if name, ok := lit(kv.Key); ok && cIdent.MatchString(name) {
				if u, ok := kv.Value.(*ast.UnaryExpr); ok && u.Op == token.AND {
					add(c.symbols, name, c.where(kv), c.optional[name], false)
				}
			}
			el = kv.Value
		}
		if c.defining && c.r.GOOS == "darwin" {
			if s, ok := lit(el); ok && strings.Contains(s, ":") && selector.MatchString(s) {
				add(c.selectors, s, c.where(el), false, true)
			}
		}
	}
}

func (c *collector) call(x *ast.CallExpr) {
	name := funName(x.Fun)
	args := x.Args

	// C symbols: a registration call takes a function pointer (&fn) and the
	// symbol name as its last argument; Dlsym takes it second.
	if name == "Dlsym" && len(args) == 2 {
		if s, ok := lit(args[1]); ok {
			add(c.symbols, s, c.where(x), true, false)
		}
	} else if len(args) >= 2 && hasAddr(args) {
		if s, ok := lit(args[len(args)-1]); ok && cIdent.MatchString(s) {
			opt := c.optional[s] || strings.Contains(strings.ToLower(name), "optional")
			add(c.symbols, s, c.where(x), opt, false)
		} else if name == "RegisterLibFunc" {
			c.r.Dynamic = append(c.r.Dynamic, c.where(x)+" RegisterLibFunc")
		}
	}
	if c.r.GOOS != "darwin" {
		return
	}

	// Objective-C classes and selectors.
	arg := func(i int) (string, bool) {
		if i < len(args) {
			return lit(args[i])
		}
		return "", false
	}
	switch name {
	case "class", "GetClass", "getClass":
		if s, ok := arg(0); ok {
			add(c.classes, s, c.where(x), false, false)
		}
	case "classExists":
		if s, ok := arg(0); ok {
			add(c.classes, s, c.where(x), true, false)
		}
	case "registerDelegateClass", "RegisterClass":
		if s, ok := arg(0); ok {
			add(c.classes, s, c.where(x), false, true)
		}
		if s, ok := arg(1); ok {
			add(c.classes, s, c.where(x), false, false)
		}
	case "send", "sendSuper", "Send", "SendSuper", "sel_", "RegisterName", "NSSelectorFromString":
		if s, ok := arg(0); ok && selector.MatchString(s) {
			add(c.selectors, s, c.where(x), c.guards[s], false)
		}
	case "sendSuperFrom", "sendSuperFromBool":
		// sendSuperFrom(definer, self, "sel", ...): the backend's [super sel].
		if s, ok := arg(2); ok && selector.MatchString(s) {
			add(c.selectors, s, c.where(x), c.guards[s], false)
		}
	case "get":
		if s, ok := arg(1); ok && selector.MatchString(s) {
			add(c.selectors, s, c.where(x), false, false)
		}
	case "respondsTo":
		if s, ok := arg(1); ok {
			add(c.selectors, s, c.where(x), true, false)
		}
	}
}

func hasAddr(args []ast.Expr) bool {
	for _, a := range args {
		if u, ok := a.(*ast.UnaryExpr); ok && u.Op == token.AND {
			return true
		}
	}
	return false
}
