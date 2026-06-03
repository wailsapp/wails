package override

import (
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v3/internal/wake/ast"
	"github.com/wailsapp/wails/v3/internal/wake/parse"
)

// layers lists the local override layers in increasing precedence order: a
// committed team layer (Taskfile.override.*) first, then a personal developer
// layer (Taskfile.local.*) which wins last. Within a layer the first existing
// extension is used. Returned taskfiles are meant to be applied in order via
// resolve.MergeTaskfile so later layers take precedence.
var layers = [][]string{
	{"Taskfile.override.yml", "Taskfile.override.yaml"},
	{"Taskfile.local.yml", "Taskfile.local.yaml"},
}

// LoadLocal discovers local override taskfiles in dir and returns them parsed
// (with includes resolved), ordered from lowest to highest precedence so the
// caller can layer them over the base with resolve.MergeTaskfile. It returns an
// empty slice when no override file exists. Builtins are intentionally NOT
// populated here so only the user's own vars layer over the base.
func LoadLocal(dir string) ([]*ast.Taskfile, error) {
	var result []*ast.Taskfile
	for _, exts := range layers {
		for _, name := range exts {
			path := filepath.Join(dir, name)
			if _, err := os.Stat(path); err != nil {
				if os.IsNotExist(err) {
					continue
				}
				// Non-missing-file errors (permission, I/O) shouldn't be
				// silently swallowed — a denied override file is more
				// likely a misconfiguration than an intentional skip.
				return nil, err
			}
			ov, err := parse.Parse(path)
			if err != nil {
				return nil, err
			}
			if err := parse.ResolveIncludes(ov); err != nil {
				return nil, err
			}
			result = append(result, ov)
			break // one file per layer
		}
	}
	return result, nil
}
