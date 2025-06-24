//go:build linux

package fileexplorer

import (
	"bytes"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	ini "gopkg.in/ini.v1"
)

func explorerBinArgs(path string, selectFile bool) (string, []string, error) {
	// Map of field codes to their replacements
	var fieldCodes = map[string]string{
		"%d": "",
		"%D": "",
		"%n": "",
		"%N": "",
		"%v": "",
		"%m": "",
		"%f": path,
		"%F": path,
		"%u": pathToURI(path),
		"%U": pathToURI(path),
	}
	fileManagerQuery := exec.Command("xdg-mime", "query", "default", "inode/directory")
	buf := new(bytes.Buffer)
	fileManagerQuery.Stdout = buf
	fileManagerQuery.Stderr = nil

	if err := fileManagerQuery.Run(); err != nil {
		return fallbackExplorerBinArgs(path, selectFile)
	}

	desktopFile, err := findDesktopFile(strings.TrimSpace((buf.String())))
	if err != nil {
		return fallbackExplorerBinArgs(path, selectFile)
	}

	cfg, err := ini.Load(desktopFile)
	if err != nil {
		// Opting to fallback rather than fail
		return fallbackExplorerBinArgs(path, selectFile)
	}

	exec := cfg.Section("Desktop Entry").Key("Exec").String()
	for fieldCode, replacement := range fieldCodes {
		exec = strings.ReplaceAll(exec, fieldCode, replacement)
	}
	args := strings.Fields(exec)
	if !strings.Contains(strings.Join(args, " "), path) {
		args = append(args, path)
	}

	return args[0], args[1:], nil
}

func sysProcAttr(path string, selectFile bool) *syscall.SysProcAttr {
	return &syscall.SysProcAttr{}
}

func fallbackExplorerBinArgs(path string, selectFile bool) (string, []string, error) {
	// NOTE: The linux fallback explorer opening is not supporting file selection
	path = filepath.Dir(path)
	return "xdg-open", []string{path}, nil
}

func pathToURI(path string) string {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return path
	}
	return "file://" + url.PathEscape(absPath)
}

func findDesktopFile(xdgFileName string) (string, error) {
	paths := []string{
		filepath.Join(os.Getenv("XDG_DATA_HOME"), "applications"),
		filepath.Join(os.Getenv("HOME"), ".local", "share", "applications"),
		"/usr/share/applications",
	}

	for _, path := range paths {
		desktopFile := filepath.Join(path, xdgFileName)
		if _, err := os.Stat(desktopFile); err == nil {
			return desktopFile, nil
		}
	}
	err := fmt.Errorf("desktop file not found: %s", xdgFileName)
	return "", err
}
