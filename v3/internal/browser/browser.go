// Package browser provides functions to open URLs and files in the default browser.
package browser

import (
	"os/exec"
	"runtime"
)

// OpenURL opens the named URL in the default browser.
func OpenURL(url string) error {
	return open(url)
}

// OpenFile opens the named file in the default browser or file handler.
func OpenFile(path string) error {
	return open(path)
}

func open(target string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", target)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", target)
	default:
		cmd = exec.Command("xdg-open", target)
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	go cmd.Wait() //nolint:errcheck
	return nil
}
