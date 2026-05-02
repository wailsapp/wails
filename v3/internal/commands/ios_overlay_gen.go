package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
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

	// Locate the internal template file to source content
	root, err := repoRoot()
	if err != nil {
		return err
	}
	tmplPath := filepath.Join(root, "v3", "internal", "commands", "build_assets", "ios", "main_ios.go")
	content, err := os.ReadFile(tmplPath)
	if err != nil {
		return fmt.Errorf("read template %s: %w", tmplPath, err)
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

// repoRoot attempts to find the repository root relative to this file location.
func repoRoot() (string, error) {
	// Resolve based on the location of this source at build time if possible.
	self, err := os.Getwd()
	if err != nil {
		return "", err
	}
	// Walk up until we find a directory containing v3/internal/commands
	probe := self
	for i := 0; i < 10; i++ {
		p := filepath.Join(probe, "v3", "internal", "commands")
		if st, err := os.Stat(p); err == nil && st.IsDir() {
			return probe, nil
		}
		next := filepath.Dir(probe)
		if next == probe {
			break
		}
		probe = next
	}
	return "", fs.ErrNotExist
}
