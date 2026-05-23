package parse

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"text/template"

	"github.com/wailsapp/wails/v3/internal/wake/ast"
)

var bareVarRegex = regexp.MustCompile(`\{\{(\s*)([A-Z_][A-Z0-9_]*)(\s*)\}\}`)

func ExpandTemplates(s string, vars map[string]*ast.Var) string {
	if !strings.Contains(s, "{{") {
		return s
	}

	s = bareVarRegex.ReplaceAllString(s, "{{${1}.${2}${3}}}")

	funcMap := template.FuncMap{
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

	tmplData := make(map[string]interface{})
	for name, vr := range vars {
		val := vr.Value
		if val == "" {
			val = vr.Static
		}
		tmplData[name] = val
	}

	tmpl, err := template.New("wake").Funcs(funcMap).Parse(s)
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

	cmd := exec.Command("sh", "-c", vr.Shell)
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
