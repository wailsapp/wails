package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// IOSOverlayGenOptions holds parameters for overlay generation.
type IOSOverlayGenOptions struct {
	Out    string `description:"Path to write overlay.json" default:"build/ios/xcode/overlay.json"`
	Config string `description:"Path to build/config.yml (optional)" default:"build/config.yml"`
}

// IOSOverlayGen generates a Go build overlay JSON that injects a generated
// main_ios.gen.go exporting WailsIOSMain() which calls main().
//
// It writes:
// - <Out> : overlay JSON file
// - <dir>/gen/main_ios.gen.go : the generated Go file referenced by the overlay
//
// The overlay maps <appDir>/main_ios.gen.go -> <dir>/gen/main_ios.gen.go
func IOSOverlayGen(options *IOSOverlayGenOptions) error { // options currently unused beyond defaults
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

	// Source the canonical main_ios.go from the build assets embedded in wails3.
	// Reading from the wails source tree (repoRoot) only works when building
	// inside the wails repo itself; a normal user project has no such tree, so
	// this must come from the embedded FS.
	content, err := buildAssets.ReadFile("build_assets/ios/main_ios.go")
	if err != nil {
		return fmt.Errorf("read embedded ios main_ios.go template: %w", err)
	}

	genDir := filepath.Join(targetDir, "gen")
	if err := os.MkdirAll(genDir, 0o755); err != nil {
		return err
	}
	genGo := filepath.Join(genDir, "main_ios.gen.go")
	if err := os.WriteFile(genGo, content, 0o644); err != nil {
		return err
	}

	// Determine app dir (current working directory)
	appDir, err := os.Getwd()
	if err != nil {
		return err
	}
	virtual := filepath.Join(appDir, "main_ios.gen.go")

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

// IOSOverlayGenCmd is a CLI entry compatible with NewSubCommandFunction.
// Defaults:
//   config: ./build/config.yml (optional)
//   out:    ./build/ios/xcode/overlay.json
func IOSOverlayGenCmd() error {
	// Default paths relative to CWD
	out := filepath.Join("build", "ios", "xcode", "overlay.json")
	cfg := filepath.Join("build", "config.yml")
	return IOSOverlayGen(&IOSOverlayGenOptions{Out: out, Config: cfg})
}
