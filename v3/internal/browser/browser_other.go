//go:build !darwin && !windows

package browser

import "os/exec"

var openCmd = func(target string) *exec.Cmd {
	return exec.Command("xdg-open", target)
}

func open(target string) error {
	cmd := openCmd(target)
	if err := cmd.Start(); err != nil {
		return err
	}
	go cmd.Wait() //nolint:errcheck
	return nil
}
