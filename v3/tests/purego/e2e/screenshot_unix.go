//go:build linux || darwin

package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"runtime"
)

// screenshot writes the whole screen to file as PNG: ImageMagick's import on
// X, screencapture on macOS.
func screenshot(file string) error {
	cmd := exec.Command("import", "-window", "root", file)
	if runtime.GOOS == "darwin" {
		cmd = exec.Command("screencapture", "-x", file)
	}
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("%v: %s", err, bytes.TrimSpace(out))
	}
	return nil
}
