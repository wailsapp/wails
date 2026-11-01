package pipeline

import (
	"fmt"
	"slices"
)

// Pure-Go executables do not list GTK in their ELF imports. Use the resolved
// compiler tags, including per-executable tags, to select their runtime profile.
func appImageGTKVersion(plan Plan, configured int) (int, error) {
	if configured != 0 {
		return configured, nil
	}
	version := 0
	for _, node := range plan.Nodes {
		compile, ok := node.Spec.(CompileSpec)
		if !ok || compile.Console {
			continue
		}
		selected := 4
		if slices.Contains(compile.Tags, "gtk3") {
			selected = 3
		}
		if version != 0 && version != selected {
			return 0, fmt.Errorf("AppImage GUI executables must use the same GTK version; found GTK 3 and GTK 4 build tags")
		}
		version = selected
	}
	return version, nil
}
