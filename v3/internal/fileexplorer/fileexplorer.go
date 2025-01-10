package fileexplorer

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	ini "gopkg.in/ini.v1"
)

type explorerBinArgs func(path string, selectFile bool) (string, []string, error)

func OpenFileManager(path string, selectFile bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	path = os.ExpandEnv(path)
	path = filepath.Clean(path)
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("failed to resolve the absolute path: %w", err)
	}
	path = absPath
	if pathInfo, err := os.Stat(path); err != nil {
		return fmt.Errorf("failed to access the specified path: %w", err)
	} else {
		selectFile = selectFile && !pathInfo.IsDir()
	}

	var (
		explorerBinArgs explorerBinArgs
		ignoreExitCode  bool = false
	)

	switch runtime.GOOS {
	case "windows":
		explorerBinArgs = windowsExplorerBinArgs
		// NOTE: Disabling the exit code check on Windows system. Workaround for explorer.exe
		// exit code handling (https://github.com/microsoft/WSL/issues/6565)
		ignoreExitCode = true
	case "darwin":
		explorerBinArgs = darwinExplorerBinArgs
	case "linux":
		explorerBinArgs = linuxExplorerBinArgs
	default:
		return errors.New("unsupported platform: " + runtime.GOOS)
	}

	explorerBin, explorerArgs, err := explorerBinArgs(path, selectFile)
	if err != nil {
		return fmt.Errorf("failed to determine the file explorer binary: %w", err)
	}

	cmd := exec.CommandContext(ctx, explorerBin, explorerArgs...)
	cmd.Stdout = nil
	cmd.Stderr = nil

	if err := cmd.Run(); err != nil {
		if !ignoreExitCode {
			return fmt.Errorf("failed to open the file explorer: %w", err)
		}
	}
	return nil
}

var windowsExplorerBinArgs explorerBinArgs = func(path string, selectFile bool) (string, []string, error) {
	args := []string{}
	if selectFile {
		args = append(args, fmt.Sprintf("/select,\"%s\"", path))
	} else {
		args = append(args, path)
	}
	return "explorer.exe", args, nil
}

var darwinExplorerBinArgs explorerBinArgs = func(path string, selectFile bool) (string, []string, error) {
	args := []string{}
	if selectFile {
		args = append(args, "-R")
	}

	args = append(args, path)
	return "open", args, nil
}

var linuxExplorerBinArgs explorerBinArgs = func(path string, selectFile bool) (string, []string, error) {
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
		return linuxFallbackExplorerBinArgs(path, selectFile)
	}

	desktopFile, err := findDesktopFile(strings.TrimSpace((buf.String())))
	if err != nil {
		return linuxFallbackExplorerBinArgs(path, selectFile)
	}

	cfg, err := ini.Load(desktopFile)
	if err != nil {
		// Opting to fallback rather than fail
		return linuxFallbackExplorerBinArgs(path, selectFile)
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

var linuxFallbackExplorerBinArgs explorerBinArgs = func(path string, selectFile bool) (string, []string, error) {
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
