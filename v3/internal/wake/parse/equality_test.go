package parse

import (
	"testing"

	"github.com/wailsapp/wails/v3/internal/wake/ast"
)

func TestTemplateEqualityAlternatives(t *testing.T) {
	for _, tt := range []struct {
		input, value, want string
	}{
		{`{{if eq .MODE "dev" "test"}}enabled{{else}}disabled{{end}}`, "dev", "enabled"},
		{`{{if eq .MODE "dev" "test"}}enabled{{else}}disabled{{end}}`, "test", "enabled"},
		{`{{if eq .MODE "dev" "test"}}enabled{{else}}disabled{{end}}`, "release", "disabled"},
		{`{{eq .MODE "dev" "test" "release"}}`, "release", "true"},
		{`{{eq .MODE "dev"}}`, "test", "false"},
		{`{{eq .MODE true false}}`, "false", "true"},
		{`{{eq .MODE "dev"}}`, "dev", "true"},
		{`{{eq .MODE}}`, "dev", `{{eq .MODE}}`},
	} {
		t.Run(tt.input+"/"+tt.value, func(t *testing.T) {
			vars := map[string]*ast.Var{"MODE": {Value: tt.value}}
			if got := ExpandTemplates(tt.input, vars); got != tt.want {
				t.Errorf("expanded = %q, want %q", got, tt.want)
			}
		})
	}
}
