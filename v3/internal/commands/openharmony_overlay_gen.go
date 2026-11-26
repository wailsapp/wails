package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// OpenHarmonyOverlayGenOptions holds parameters for OpenHarmony overlay generation.
type OpenHarmonyOverlayGenOptions struct {
	Out    string `description:"Path to write overlay.json" default:"build/openharmony/overlay.json"`
	Config string `description:"Path to build/config.yml (optional)" default:"build/config.yml"`
}

// OpenHarmonyOverlayGen generates a Go build overlay JSON that injects a
// generated main_openharmony.gen.go into the application's main package.
//
// That file's init() calls application.RegisterOpenharmonyMain(main) and
// exports the C ABI the ArkTS shell calls (WailsBridgeOpen, WailsAsset,
// WailsRequest, ...). The c-shared library needs this because in
// -buildmode=c-shared main() is never invoked automatically.
//
// It writes:
//   - <Out>             : overlay JSON file
//   - <dir>/gen/main_openharmony.gen.go : the generated Go file referenced by the overlay
//
// The overlay maps <appDir>/main_openharmony.gen.go -> <dir>/gen/main_openharmony.gen.go
func OpenHarmonyOverlayGen(options *OpenHarmonyOverlayGenOptions) error {
	out := options.Out
	if out == "" {
		return errors.New("--out is required (path to write overlay.json)")
	}

	absOut, err := filepath.Abs(out)
	if err != nil {
		return err
	}
	targetDir := filepath.Dir(absOut)
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		return err
	}

	// Source the canonical main_openharmony.go from the build assets embedded
	// in wails3 (NOT the wails source tree — a normal user project has no such
	// tree, so this must come from the embedded FS).
	content, err := buildAssets.ReadFile("build_assets/openharmony/main_openharmony.go")
	if err != nil {
		return fmt.Errorf("read embedded openharmony main_openharmony.go template: %w", err)
	}

	genDir := filepath.Join(targetDir, "gen")
	if err := os.MkdirAll(genDir, 0o755); err != nil {
		return err
	}
	genGo := filepath.Join(genDir, "main_openharmony.gen.go")
	if err := os.WriteFile(genGo, content, 0o644); err != nil {
		return err
	}

	appDir, err := os.Getwd()
	if err != nil {
		return err
	}
	virtual := filepath.Join(appDir, "main_openharmony.gen.go")

	type overlay struct {
		Replace map[string]string `json:"Replace"`
	}
	ov := overlay{Replace: map[string]string{virtual: genGo}}
	data, err := json.MarshalIndent(ov, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(absOut, data, 0o644); err != nil {
		return err
	}
	return nil
}

// OpenHarmonyOverlayGenCmd is a CLI entry compatible with NewSubCommandFunction.
// Defaults:
//   config: ./build/config.yml (optional)
//   out:    ./build/openharmony/overlay.json
func OpenHarmonyOverlayGenCmd() error {
	out := filepath.Join("build", "openharmony", "overlay.json")
	cfg := filepath.Join("build", "config.yml")
	return OpenHarmonyOverlayGen(&OpenHarmonyOverlayGenOptions{Out: out, Config: cfg})
}
