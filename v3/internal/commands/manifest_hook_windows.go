//go:build windows

package commands

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

func manifestHookInterpreter(script string) (string, error) {
	name := ""
	switch strings.ToLower(filepath.Ext(script)) {
	case ".cmd", ".bat":
		name = "cmd.exe"
	case ".ps1":
		name = "powershell.exe"
	default:
		return "", fmt.Errorf("Windows hook script %q must use .cmd, .bat, or .ps1", script)
	}
	path, err := exec.LookPath(name)
	if err != nil {
		return "", fmt.Errorf("hook interpreter %s not found: %w", name, err)
	}
	return path, nil
}

func manifestHookCommand(script string) (string, []string, error) {
	interpreter, err := manifestHookInterpreter(script)
	if err != nil {
		return "", nil, err
	}
	switch strings.ToLower(filepath.Ext(script)) {
	case ".cmd", ".bat":
		return interpreter, []string{"/D", "/C", script}, nil
	case ".ps1":
		return interpreter, []string{"-NoProfile", "-NonInteractive", "-File", script}, nil
	default:
		panic("manifestHookInterpreter accepted an unsupported extension")
	}
}
