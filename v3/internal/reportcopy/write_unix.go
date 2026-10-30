//go:build !windows && !darwin

package reportcopy

import (
	"errors"
	"os"
	"os/exec"
)

// Write copies text using the current desktop session's clipboard utility.
func Write(text string) error {
	if os.Getenv("WAYLAND_DISPLAY") != "" {
		return writeCommand(text, "wl-copy")
	}
	if path, err := exec.LookPath("xclip"); err == nil {
		return writeCommand(text, path, "-selection", "clipboard", "-in")
	}
	if path, err := exec.LookPath("xsel"); err == nil {
		return writeCommand(text, path, "--clipboard", "--input")
	}
	return errors.New("copying a report requires xclip or xsel")
}
