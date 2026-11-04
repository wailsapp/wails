package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"go/ast"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Library versions. The purego backend binds symbols at runtime, so which
// functions it may call depends on the version of each .so actually loaded,
// not on headers at build time. Two rules keep that safe:
//
//  1. A function introduced after the supported floor of its library (GIR
//     "version" attribute > floors.json) must be registered as optional:
//     a required registration fails on every older system.
//  2. Every call to an optional function must be guarded in the calling
//     function: a nil check of its function variable, haveSymbol, or a
//     version check (…AtLeast / get_minor_version).

// floors maps a GIR namespace-version ("Gtk-4.0") to the oldest library
// version the backend supports ("4.6").
type floors map[string]string

func loadFloors(path string) (floors, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	f := floors{}
	return f, json.Unmarshal(data, &f)
}

type since struct {
	Namespace string // GIR namespace-version, e.g. Gtk-4.0
	Version   string // version the symbol appeared in, "" when not recorded
}

// readGIR indexes every C function in the GIR files of dir by its C name.
func readGIR(dir string) (map[string]since, error) {
	paths, err := filepath.Glob(filepath.Join(dir, "*.gir"))
	if err != nil || len(paths) == 0 {
		return nil, fmt.Errorf("no .gir files in %s", dir)
	}
	out := map[string]since{}
	for _, p := range paths {
		ns := strings.TrimSuffix(filepath.Base(p), ".gir")
		f, err := os.Open(p)
		if err != nil {
			return nil, err
		}
		dec := xml.NewDecoder(f)
		for {
			tok, err := dec.Token()
			if err != nil {
				break
			}
			se, ok := tok.(xml.StartElement)
			if !ok {
				continue
			}
			switch se.Name.Local {
			case "function", "method", "constructor", "function-macro":
			default:
				continue
			}
			var cname, version string
			for _, a := range se.Attr {
				switch {
				case a.Name.Local == "identifier" && a.Name.Space != "":
					cname = a.Value
				case a.Name.Local == "version" && a.Name.Space == "":
					version = a.Value
				}
			}
			if cname != "" {
				if _, seen := out[cname]; !seen {
					out[cname] = since{ns, version}
				}
			}
		}
		f.Close()
	}
	return out, nil
}

// newer reports whether version a is later than b ("4.14" > "4.6").
func newer(a, b string) bool {
	pa, pb := strings.Split(a, "."), strings.Split(b, ".")
	for i := 0; i < len(pa) || i < len(pb); i++ {
		var x, y int
		if i < len(pa) {
			x, _ = strconv.Atoi(pa[i])
		}
		if i < len(pb) {
			y, _ = strconv.Atoi(pb[i])
		}
		if x != y {
			return x > y
		}
	}
	return false
}

// checkVersions applies rule 1 using the GIR data and the floors.
func checkVersions(r *report, gir map[string]since, fl floors) {
	for _, e := range r.Symbols {
		s, ok := gir[e.Name]
		if !ok {
			continue
		}
		e.Since, e.Namespace = s.Version, s.Namespace
		floor, ok := fl[s.Namespace]
		if !ok || s.Version == "" || !newer(s.Version, floor) {
			continue
		}
		e.Gated = true
		if !e.Optional {
			r.Versions = append(r.Versions, fmt.Sprintf("%s appeared in %s %s, after the floor %s, but is registered as required (%s)",
				e.Name, s.Namespace, s.Version, floor, e.Where[0]))
		}
	}
}

// checkGuards applies rule 2: each function that calls an optional symbol
// must mention a guard for it.
func (c *collector) checkGuards(f *ast.File) {
	for _, d := range f.Decls {
		fd, ok := d.(*ast.FuncDecl)
		if !ok || fd.Body == nil {
			continue
		}
		ast.Inspect(fd.Body, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			name := funName(call.Fun)
			sym := strings.TrimPrefix(strings.TrimPrefix(name, "_native_"), "_Cfunc_")
			e := c.symbols[sym]
			if e == nil || !e.Optional || fd.Name.Name == sym {
				return true
			}
			if !guarded(fd.Body, sym) {
				c.r.Unguarded = append(c.r.Unguarded, fmt.Sprintf("%s calls optional %s without a guard (%s)", fd.Name.Name, sym, c.where(call)))
			}
			return true
		})
	}
}

// respondedSelectors returns the selector literals that body passes to
// respondsTo, the Objective-C guard for a selector newer than the floor.
func respondedSelectors(body *ast.BlockStmt) map[string]bool {
	out := map[string]bool{}
	ast.Inspect(body, func(n ast.Node) bool {
		if call, ok := n.(*ast.CallExpr); ok && funName(call.Fun) == "respondsTo" && len(call.Args) == 2 {
			if s, ok := lit(call.Args[1]); ok {
				out[s] = true
			}
		}
		return true
	})
	return out
}

func guarded(body *ast.BlockStmt, sym string) bool {
	found := false
	ast.Inspect(body, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.BinaryExpr:
			// _native_sym != nil / == nil
			for _, side := range []ast.Expr{x.X, x.Y} {
				if id, ok := side.(*ast.Ident); ok && strings.HasSuffix(id.Name, sym) && id.Name != sym {
					found = true
				}
			}
		case *ast.CallExpr:
			n := funName(x.Fun)
			if n == "haveSymbol" || strings.HasSuffix(n, "AtLeast") || strings.HasSuffix(n, "_get_minor_version") {
				found = true
			}
		}
		return !found
	})
	return found
}
