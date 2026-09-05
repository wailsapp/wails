package parse

import (
	"fmt"
	"regexp"
	"runtime"
	"strings"
	"text/template"

	"github.com/wailsapp/wails/v3/internal/wake/ast"
	"github.com/wailsapp/wails/v3/internal/wake/platform"
)

var bareVarRegex = regexp.MustCompile(`\{\{(\s*)([A-Z_][A-Z0-9_]*)(\s*)\}\}`)
var defaultBareVarRegex = regexp.MustCompile(`(\|\s*default\s+)([A-Z_][A-Z0-9_]*)(\s*[\}\|])`)

func ExpandTemplates(s string, vars map[string]*ast.Var) string {
	if !strings.Contains(s, "{{") {
		return s
	}

	s = bareVarRegex.ReplaceAllString(s, "{{${1}.${2}${3}}}")
	s = defaultBareVarRegex.ReplaceAllString(s, "${1}.${2}${3}")

	funcMap := template.FuncMap{
		// Platform identifiers. Taskfile syntax uses these as both variables
		// (`{{OS}}` → "darwin") and function calls in expressions
		// (`{{if eq OS "darwin"}}...`). The bareVarRegex above rewrites the
		// variable form to `{{.OS}}`; registering them as functions covers
		// the expression form so platform-conditional task names parse.
		"OS":   func() string { return runtime.GOOS },
		"ARCH": func() string { return runtime.GOARCH },
		"default": func(def, val interface{}) interface{} {
			if val == nil || val == "" {
				return def
			}
			if s, ok := val.(string); ok && s == "" {
				return def
			}
			return val
		},
		"eq": func(a, b interface{}) bool {
			return fmt.Sprint(a) == fmt.Sprint(b)
		},
		"ne": func(a, b interface{}) bool {
			return fmt.Sprint(a) != fmt.Sprint(b)
		},
		"and": func(args ...bool) bool {
			for _, a := range args {
				if !a {
					return false
				}
			}
			return len(args) > 0
		},
		"or": func(args ...bool) bool {
			for _, a := range args {
				if a {
					return true
				}
			}
			return false
		},
		"not": func(a bool) bool {
			return !a
		},
		"trim":       strings.TrimSpace,
		"trimPrefix": func(prefix, s string) string { return strings.TrimPrefix(s, prefix) },
		"trimSuffix": func(suffix, s string) string { return strings.TrimSuffix(s, suffix) },
		"lower":      strings.ToLower,
		"upper":      strings.ToUpper,
		"replace":    func(old, new, s string) string { return strings.ReplaceAll(s, old, new) },
		"split":      func(sep, s string) []string { return strings.Split(s, sep) },
		"join":       func(sep string, parts []string) string { return strings.Join(parts, sep) },
		"hasPrefix":  func(prefix, s string) bool { return strings.HasPrefix(s, prefix) },
		"hasSuffix":  func(suffix, s string) bool { return strings.HasSuffix(s, suffix) },
		"contains":   func(substr, s string) bool { return strings.Contains(s, substr) },
	}

	// map[string]string (not interface{}) so missingkey=zero produces an
	// empty string for unknown keys instead of the nil-renders-as-"<no value>"
	// surprise. Every wake var value is a string, so the concrete type is
	// also more honest.
	tmplData := make(map[string]string)
	for name, vr := range vars {
		val := vr.Value
		if val == "" || strings.Contains(val, "{{") {
			if val == "" && !strings.Contains(vr.Static, "{{") {
				val = vr.Static
			} else if strings.Contains(val, "{{") {
				val = ""
			}
		}
		tmplData[name] = val
	}

	// missingkey=zero so an unresolved {{.X}} renders as an empty string
	// instead of the literal text "<no value>". This matches Taskfile
	// semantics, and lets the `default` function fire correctly for the
	// common `{{.X | default "y"}}` pattern even when X is undefined.
	tmpl, err := template.New("wake").Option("missingkey=zero").Funcs(funcMap).Parse(s)
	if err != nil {
		return s
	}

	var buf strings.Builder
	if err := tmpl.Execute(&buf, tmplData); err != nil {
		return s
	}

	return buf.String()
}

func ResolveVarShell(vr *ast.Var) error {
	if vr.Shell == "" || vr.Value != "" {
		return nil
	}

	cmd := platform.ShellCommand(vr.Shell)
	out, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("wake: var shell command %q failed: %w", vr.Shell, err)
	}

	vr.Value = strings.TrimSpace(string(out))
	return nil
}

func ResolveAllVarShells(vars map[string]*ast.Var) error {
	for _, vr := range vars {
		if vr.Shell != "" && vr.Value == "" {
			if err := ResolveVarShell(vr); err != nil {
				return err
			}
		}
	}
	return nil
}
