package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// AndroidOverlayGenOptions holds parameters for Android overlay generation.
type AndroidOverlayGenOptions struct {
	Out    string `description:"Path to write overlay.json" default:"build/android/overlay.json"`
	Config string `description:"Path to build/config.yml (optional)" default:"build/config.yml"`
}

// AndroidOverlayGen generates a Go build overlay JSON that injects a generated
// main_android.gen.go into the application's main package. That file's init()
// calls application.RegisterAndroidMain(main), which the c-shared library needs
// because main() is not invoked automatically in -buildmode=c-shared.
//
// This mirrors the iOS overlay (see IOSOverlayGen): the file is injected
// virtually via the overlay so the user's project root stays clean.
//
// It writes:
//   - <Out>             : overlay JSON file
//   - <dir>/gen/main_android.gen.go : the generated Go file referenced by the overlay
//
// The overlay maps <appDir>/main_android.gen.go -> <dir>/gen/main_android.gen.go
func AndroidOverlayGen(options *AndroidOverlayGenOptions) error {
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

	// Source the canonical main_android.go from the build assets embedded in
	// wails3 (NOT the wails source tree — a normal user project has no such
	// tree, so this must come from the embedded FS).
	content, err := buildAssets.ReadFile("build_assets/android/main_android.go")
	if err != nil {
		return fmt.Errorf("read embedded android main_android.go template: %w", err)
	}

	genDir := filepath.Join(targetDir, "gen")
	if err := os.MkdirAll(genDir, 0o755); err != nil {
		return err
	}
	genGo := filepath.Join(genDir, "main_android.gen.go")
	if err := os.WriteFile(genGo, content, 0o644); err != nil {
		return err
	}

	appDir, err := os.Getwd()
	if err != nil {
		return err
	}
	virtual := filepath.Join(appDir, "main_android.gen.go")

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
