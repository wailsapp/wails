//go:build linux

package fileexplorer

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"syscall"
)

// when possible; the fallback method does not support selecting a file.
func explorerBinArgs(ctx context.Context, path string, selectFile bool) (string, []string, error) {
	fileManagerQuery := exec.CommandContext(ctx, "xdg-mime", "query", "default", "inode/directory")
	buf := new(bytes.Buffer)
	fileManagerQuery.Stdout = buf
	fileManagerQuery.Stderr = nil

	if err := fileManagerQuery.Run(); err != nil {
		if err := ctx.Err(); err != nil {
			return "", nil, err
		}
		return fallbackExplorerBinArgs(path, selectFile)
	}

	desktopFilePath, err := findDesktopFile(strings.TrimSpace((buf.String())))
	if err != nil {
		return fallbackExplorerBinArgs(path, selectFile)
	}

	entry, err := ParseDesktopFile(desktopFilePath)
	if err != nil {
		// Opting to fallback rather than fail
		return fallbackExplorerBinArgs(path, selectFile)
	}

	uri := pathToURI(path)
	fieldCodes := strings.NewReplacer(
		"%d", "", "%D", "", "%n", "", "%N", "", "%v", "", "%m", "",
		"%f", path, "%F", path, "%u", uri, "%U", uri,
	)
	// Split before expanding so whitespace in a selected path stays in one argument.
	fields := strings.Fields(entry.Exec)
	args := fields[:0]
	targetExpanded := false
	for _, field := range fields {
		if expanded := fieldCodes.Replace(field); expanded != "" {
			args = append(args, expanded)
			targetExpanded = targetExpanded || strings.Contains(field, "%f") || strings.Contains(field, "%F") ||
				strings.Contains(field, "%u") || strings.Contains(field, "%U")
		}
	}
	if len(args) == 0 {
		return fallbackExplorerBinArgs(path, selectFile)
	}
	if !targetExpanded && !slices.Contains(args[1:], path) && !slices.Contains(args[1:], uri) {
		args = append(args, path)
	}

	return args[0], args[1:], nil
}

func sysProcAttr(path string, selectFile bool) *syscall.SysProcAttr {
	return &syscall.SysProcAttr{}
}

func fallbackExplorerBinArgs(path string, selectFile bool) (string, []string, error) {
	// NOTE: The linux fallback explorer opening does not support file selection

	stat, err := os.Stat(path)
	if err != nil {
		return "", []string{}, fmt.Errorf("stat path: %w", err)
	}

	// If the path is a file, we want to open the directory containing the file
	if !stat.IsDir() {
		path = filepath.Dir(path)
	}

	return "xdg-open", []string{path}, nil
}

func pathToURI(path string) string {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return path
	}
	// Use url.URL to properly construct file URIs.
	// url.PathEscape incorrectly escapes forward slashes (/ -> %2F),
	// which breaks file manager path parsing.
	u := &url.URL{
		Scheme: "file",
		Path:   absPath,
	}
	return u.String()
}

func findDesktopFile(xdgFileName string) (string, error) {
	dataHome := os.Getenv("XDG_DATA_HOME")
	if !filepath.IsAbs(dataHome) {
		dataHome = ""
		if home, err := os.UserHomeDir(); err == nil {
			dataHome = filepath.Join(home, ".local", "share")
		}
	}
	dataDirs := os.Getenv("XDG_DATA_DIRS")
	if dataDirs == "" {
		dataDirs = "/usr/local/share:/usr/share"
	}
	paths := append([]string{dataHome}, filepath.SplitList(dataDirs)...)

	for _, dataDir := range paths {
		if !filepath.IsAbs(dataDir) {
			continue
		}
		desktopFile := filepath.Join(dataDir, "applications", xdgFileName)
		if info, err := os.Stat(desktopFile); err == nil && !info.IsDir() {
			return desktopFile, nil
		}
	}
	err := fmt.Errorf("desktop file not found: %s", xdgFileName)
	return "", err
}
