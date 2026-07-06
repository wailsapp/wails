package migrate

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"sort"
	"strconv"
	"strings"
)

// v2compatAlias is the local name used for the v2compat runtime import in
// generated code (avoids clashing with a stdlib runtime import).
const v2compatAlias = "v2runtime"

// GenerateMain performs textual surgery on the file containing wails.Run:
// the import block and the wails.Run statement are replaced with their v3
// equivalents, everything else is preserved byte-for-byte, and the result is
// gofmt-formatted.
func GenerateMain(proj *V2Project, opts *V3Options) ([]byte, error) {
	main := proj.Main
	src := main.Source

	appVar := pickIdent(main.File, "app", "wailsApp", "wailsV3App")

	block := buildV3Block(proj, opts, appVar)

	type edit struct {
		start, end int
		text       string
	}
	var edits []edit

	// Replace the wails.Run statement.
	stmtStart := main.Fset.Position(main.RunStmt.Pos()).Offset
	stmtEnd := main.Fset.Position(main.RunStmt.End()).Offset
	edits = append(edits, edit{stmtStart, stmtEnd, block})

	// Replace the import declaration(s).
	importText := buildImports(proj, opts)
	first := true
	for _, decl := range main.File.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok || gen.Tok != token.IMPORT {
			continue
		}
		start := main.Fset.Position(gen.Pos()).Offset
		end := main.Fset.Position(gen.End()).Offset
		if first {
			edits = append(edits, edit{start, end, importText})
			first = false
		} else {
			// Merge extra import decls away (their specs are already included
			// in importText via main.File.Imports).
			edits = append(edits, edit{start, end, ""})
		}
	}
	if first {
		return nil, fmt.Errorf("no import declaration found in %s", main.Path)
	}

	// Apply edits back-to-front so offsets stay valid.
	sort.Slice(edits, func(i, j int) bool { return edits[i].start > edits[j].start })
	out := make([]byte, len(src))
	copy(out, src)
	for _, e := range edits {
		out = append(out[:e.start], append([]byte(e.text), out[e.end:]...)...)
	}

	return pruneAndFormat(out)
}

// pickIdent returns the first candidate not used as an identifier in the file.
func pickIdent(file *ast.File, candidates ...string) string {
	used := map[string]bool{}
	ast.Inspect(file, func(n ast.Node) bool {
		if ident, ok := n.(*ast.Ident); ok {
			used[ident.Name] = true
		}
		return true
	})
	for _, c := range candidates {
		if !used[c] {
			return c
		}
	}
	return candidates[len(candidates)-1] + "2"
}

func writeFields(sb *strings.Builder, fields []GenField, indent string) {
	for _, f := range fields {
		sb.WriteString(indent + f.Name + ": " + f.Expr + ",\n")
	}
}

// buildV3Block renders the replacement for the wails.Run statement.
func buildV3Block(proj *V2Project, opts *V3Options, appVar string) string {
	var sb strings.Builder

	sb.WriteString(appVar + " := application.New(application.Options{\n")
	writeFields(&sb, opts.App, "\t")

	// Services (bound structs first, lifecycle bridge last).
	if len(opts.Services) > 0 || opts.NeedsLifecycleService() {
		sb.WriteString("\tServices: []application.Service{\n")
		for _, svc := range opts.Services {
			sb.WriteString("\t\tapplication.NewService(" + svc + "),\n")
		}
		if opts.NeedsLifecycleService() {
			args := []string{"nil", "nil", "nil"}
			for i, cb := range []string{opts.OnStartup, opts.OnDomReady, opts.OnShutdown} {
				if cb != "" {
					args[i] = cb
				}
			}
			sb.WriteString("\t\t// Bridges the v2 OnStartup/OnDomReady/OnShutdown callbacks.\n")
			sb.WriteString("\t\t" + v2compatAlias + ".NewLifecycleService(" + strings.Join(args, ", ") + "),\n")
		}
		sb.WriteString("\t},\n")
	}

	if opts.OnBeforeClose != "" {
		sb.WriteString("\t// v2 OnBeforeClose returned true to prevent closing; v3 ShouldQuit returns true to allow quitting.\n")
		sb.WriteString("\tShouldQuit: func() bool {\n\t\treturn !(" + opts.OnBeforeClose + ")(context.Background())\n\t},\n")
	}

	if len(opts.SingleInstance) > 0 {
		sb.WriteString("\tSingleInstance: &application.SingleInstanceOptions{\n")
		writeFields(&sb, opts.SingleInstance, "\t\t")
		sb.WriteString("\t},\n")
	}
	if len(opts.AppMac) > 0 {
		sb.WriteString("\tMac: application.MacOptions{\n")
		writeFields(&sb, opts.AppMac, "\t\t")
		sb.WriteString("\t},\n")
	}
	if len(opts.AppWin) > 0 {
		sb.WriteString("\tWindows: application.WindowsOptions{\n")
		writeFields(&sb, opts.AppWin, "\t\t")
		sb.WriteString("\t},\n")
	}
	if len(opts.AppLinux) > 0 {
		sb.WriteString("\tLinux: application.LinuxOptions{\n")
		writeFields(&sb, opts.AppLinux, "\t\t")
		sb.WriteString("\t},\n")
	}
	sb.WriteString("})\n\n")

	sb.WriteString(appVar + ".Window.NewWithOptions(application.WebviewWindowOptions{\n")
	writeFields(&sb, opts.Win, "\t")
	if len(opts.WinMac) > 0 {
		sb.WriteString("\tMac: application.MacWindow{\n")
		writeFields(&sb, opts.WinMac, "\t\t")
		sb.WriteString("\t},\n")
	}
	if len(opts.WinWin) > 0 {
		sb.WriteString("\tWindows: application.WindowsWindow{\n")
		writeFields(&sb, opts.WinWin, "\t\t")
		sb.WriteString("\t},\n")
	}
	if len(opts.WinLinux) > 0 {
		sb.WriteString("\tLinux: application.LinuxWindow{\n")
		writeFields(&sb, opts.WinLinux, "\t\t")
		sb.WriteString("\t},\n")
	}
	sb.WriteString("})\n\n")

	// Preserve the original error-handling shape.
	main := proj.Main
	switch {
	case main.ErrIdent != "" && main.AssignTok == token.DEFINE:
		sb.WriteString(main.ErrIdent + " := " + appVar + ".Run()")
	case main.ErrIdent != "":
		sb.WriteString(main.ErrIdent + " = " + appVar + ".Run()")
	default:
		sb.WriteString(appVar + ".Run()")
	}

	return sb.String()
}

// buildImports renders the replacement import declaration: the original
// imports minus everything under github.com/wailsapp/wails/v2, plus the v3
// imports the generated code needs.
func buildImports(proj *V2Project, opts *V3Options) string {
	type imp struct{ alias, path string }
	var imports []imp
	seen := map[string]bool{}

	for _, spec := range proj.Main.File.Imports {
		path, err := strconv.Unquote(spec.Path.Value)
		if err != nil || strings.HasPrefix(path, "github.com/wailsapp/wails/v2") {
			continue
		}
		alias := ""
		if spec.Name != nil {
			alias = spec.Name.Name
		}
		imports = append(imports, imp{alias, path})
		seen[path] = true
	}

	add := func(alias, path string) {
		if !seen[path] {
			imports = append(imports, imp{alias, path})
			seen[path] = true
		}
	}
	add("", "github.com/wailsapp/wails/v3/pkg/application")
	if opts.NeedsLifecycleService() {
		add(v2compatAlias, V2CompatRuntimeImport)
	}
	if opts.OnBeforeClose != "" {
		add("", "context")
	}

	sort.Slice(imports, func(i, j int) bool { return imports[i].path < imports[j].path })

	var sb strings.Builder
	sb.WriteString("import (\n")
	for _, im := range imports {
		sb.WriteString("\t")
		if im.alias != "" {
			sb.WriteString(im.alias + " ")
		}
		sb.WriteString(strconv.Quote(im.path) + "\n")
	}
	sb.WriteString(")")
	return sb.String()
}

// pruneAndFormat removes imports that became unused after surgery (values
// that only appeared inside the replaced statement) and gofmt-formats the
// result.
func pruneAndFormat(src []byte) ([]byte, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "main.go", src, parser.ParseComments|parser.SkipObjectResolution)
	if err != nil {
		return nil, fmt.Errorf("generated main file does not parse (%w); this is a bug in wails3 migrate", err)
	}

	used := map[string]bool{}
	ast.Inspect(file, func(n ast.Node) bool {
		if sel, ok := n.(*ast.SelectorExpr); ok {
			if ident, ok := sel.X.(*ast.Ident); ok {
				used[ident.Name] = true
			}
		}
		return true
	})

	type span struct{ start, end int }
	var removals []span
	for _, decl := range file.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok || gen.Tok != token.IMPORT {
			continue
		}
		for _, spec := range gen.Specs {
			imp := spec.(*ast.ImportSpec)
			if imp.Name != nil && (imp.Name.Name == "_" || imp.Name.Name == ".") {
				continue
			}
			path, err := strconv.Unquote(imp.Path.Value)
			if err != nil {
				continue
			}
			name := ""
			if imp.Name != nil {
				name = imp.Name.Name
			} else {
				name = path[strings.LastIndex(path, "/")+1:]
			}
			if !used[name] {
				removals = append(removals, span{
					fset.Position(imp.Pos()).Offset,
					fset.Position(imp.End()).Offset,
				})
			}
		}
	}

	sort.Slice(removals, func(i, j int) bool { return removals[i].start > removals[j].start })
	for _, r := range removals {
		src = append(src[:r.start], src[r.end:]...)
	}

	formatted, err := format.Source(src)
	if err != nil {
		return nil, fmt.Errorf("could not format the generated main file (%w); this is a bug in wails3 migrate", err)
	}
	return formatted, nil
}
