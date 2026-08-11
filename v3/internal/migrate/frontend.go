package migrate

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// MigrateFrontend copies the v2 frontend into the output project and adds the
// @wailsio/runtime dependency to package.json.
//
// The generated wailsjs directory is NOT copied and NOT replaced with a
// compatibility layer: it is v2 build output that cannot work against a v3
// backend, and shipping a lookalike would encourage code to stay on the v2
// API. Every import of it is listed in the migration report with its v3
// replacement instead.
func MigrateFrontend(proj *V2Project, outDir string) error {
	srcFrontend := proj.FrontendDir
	dstFrontend := filepath.Join(outDir, "frontend")

	if _, err := os.Stat(srcFrontend); os.IsNotExist(err) {
		proj.Report.Manual("frontend", "No frontend directory was found; create one or point the v3 Taskfile at your assets.")
		return nil
	}

	hadWailsJS := false
	if _, err := os.Stat(filepath.Join(srcFrontend, "wailsjs")); err == nil {
		hadWailsJS = true
	}

	err := copyTree(srcFrontend, dstFrontend, func(rel string, isDir bool) bool {
		switch {
		case rel == "node_modules" || strings.HasPrefix(rel, "node_modules"+string(filepath.Separator)):
			return false
		case rel == "wailsjs" || strings.HasPrefix(rel, "wailsjs"+string(filepath.Separator)):
			return false
		}
		return true
	})
	if err != nil {
		return err
	}

	// Ensure frontend/dist exists so `go build` can embed it before the first
	// frontend build.
	distDir := filepath.Join(dstFrontend, "dist")
	if err := os.MkdirAll(distDir, 0o755); err != nil {
		return err
	}
	gitkeep := filepath.Join(distDir, ".gitkeep")
	if _, err := os.Stat(gitkeep); os.IsNotExist(err) {
		if err := os.WriteFile(gitkeep, nil, 0o644); err != nil {
			return err
		}
	}

	if hadWailsJS {
		proj.Report.Note("The generated `frontend/wailsjs` directory was not copied: it is v2 build output and cannot work with v3. Imports of it are listed in *Port these to the v3 API*; run `wails3 generate bindings` to generate the v3 bindings into `frontend/bindings`.")
	}

	if err := addRuntimeDependency(proj, dstFrontend); err != nil {
		return err
	}

	return nil
}

// addRuntimeDependency inserts "@wailsio/runtime" into the frontend
// package.json dependencies as structured JSON so empty objects and files
// remain valid JSON.
func addRuntimeDependency(proj *V2Project, frontendDir string) error {
	path := filepath.Join(frontendDir, "package.json")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			proj.Report.Manual("frontend/package.json",
				"No package.json found. The v3 frontend runtime is the `@wailsio/runtime` npm package, which requires a bundler; plain-JS frontends must load the runtime differently (see the v3 docs).")
			return nil
		}
		return err
	}
	var packageJSON map[string]any
	if err := json.Unmarshal(data, &packageJSON); err != nil {
		return fmt.Errorf("could not parse frontend/package.json: %w", err)
	}

	deps := map[string]any{}
	if raw, ok := packageJSON["dependencies"]; ok {
		var valid bool
		deps, valid = raw.(map[string]any)
		if !valid {
			proj.Report.Manual("frontend/package.json", "The dependencies field is not a JSON object; add `@wailsio/runtime` manually and run your package manager's install.")
			return nil
		}
	}
	if _, ok := deps["@wailsio/runtime"]; ok {
		return nil
	}
	deps["@wailsio/runtime"] = "latest"
	packageJSON["dependencies"] = deps
	updated, err := json.MarshalIndent(packageJSON, "", "  ")
	if err != nil {
		return fmt.Errorf("could not encode frontend/package.json: %w", err)
	}
	updated = append(updated, '\n')

	if err := os.WriteFile(path, updated, 0o644); err != nil {
		return err
	}
	proj.Report.Note("Added `@wailsio/runtime` to frontend/package.json; run your package manager's install (or let the Taskfile do it on first build).")
	return nil
}
