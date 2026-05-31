package fileexplorer

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

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
		ignoreExitCode bool = false
	)

	switch runtime.GOOS {
	case "windows":
		// NOTE: Disabling the exit code check on Windows system. Workaround for explorer.exe
		// exit code handling (https://github.com/microsoft/WSL/issues/6565)
		ignoreExitCode = true
	case "darwin", "linux":
	default:
		return errors.New("unsupported platform: " + runtime.GOOS)
	}

	explorerBin, explorerArgs, err := explorerBinArgs(path, selectFile)
	if err != nil {
		return fmt.Errorf("failed to determine the file explorer binary: %w", err)
	}

	cmd := exec.CommandContext(ctx, explorerBin, explorerArgs...)
	cmd.SysProcAttr = sysProcAttr(path, selectFile)
	cmd.Stdout = nil
	cmd.Stderr = nil

	if err := cmd.Run(); err != nil {
		if !ignoreExitCode {
			return fmt.Errorf("failed to open the file explorer: %w", err)
		}
	}
	return nil
}
