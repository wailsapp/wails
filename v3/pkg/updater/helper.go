package updater

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// Helper-mode protocol.
//
// To avoid shipping a separate binary, the Updater re-executes the running
// application with a sentinel environment variable set. The helper-mode
// process waits for the parent (the application that just initiated the
// update) to exit, swaps the on-disk binary with the downloaded artifact,
// and then relaunches the (now-replaced) application. The helper exits
// when its work is done.
//
// All communication is via environment variables so user-supplied
// command-line flags are never confused with helper plumbing.

const (
	envHelperMode   = "WAILS_UPDATER_HELPER"        // "1" to enter helper mode
	envHelperTarget = "WAILS_UPDATER_HELPER_TARGET" // path of the running app to be replaced
	envHelperNew    = "WAILS_UPDATER_HELPER_NEW"    // path of the verified new artifact
	envHelperPID    = "WAILS_UPDATER_HELPER_PID"    // parent PID to wait for
	envHelperLog    = "WAILS_UPDATER_HELPER_LOG"    // optional log file path
)

// HandleHelperMode returns immediately when the current process was not
// spawned as an updater helper. When it WAS spawned as a helper it performs
// the swap, relaunches the application, and calls os.Exit — it never returns
// in that case.
//
// The Wails application package calls this from application.New so that
// `app.Updater.Restart` works without users wiring anything by hand.
func HandleHelperMode() {
	if os.Getenv(envHelperMode) != "1" {
		return
	}
	target := os.Getenv(envHelperTarget)
	newPath := os.Getenv(envHelperNew)
	if target == "" || newPath == "" {
		os.Exit(2)
	}
	pid, _ := strconv.Atoi(os.Getenv(envHelperPID))
	logPath := os.Getenv(envHelperLog)

	code := runHelperSwap(target, newPath, pid, logPath, waitForPID, osLauncher{})
	os.Exit(code)
}

// processWaiter abstracts "wait until pid exits" so unit tests can drive the
// swap logic without spawning real processes.
type processWaiter func(pid int, timeout time.Duration) error

// launcher abstracts how the helper kicks off the replaced binary. Tests
// substitute a recorder; production uses osLauncher.
type launcher interface {
	launch(path string) error
}

type osLauncher struct{}

func (osLauncher) launch(path string) error {
	var cmd *exec.Cmd
	if runtime.GOOS == "darwin" && filepath.Ext(path) == ".app" {
		cmd = exec.Command("open", "-n", path)
	} else {
		cmd = exec.Command(path)
	}
	// Detach: we are about to exit; the new process must not depend on us.
	cmd.Stdin = nil
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Start()
}

// runHelperSwap implements the actual swap logic. It is unexported and
// dependency-injected so unit tests can drive every branch without process
// spawning. Returns the exit code the helper should use.
func runHelperSwap(target, newPath string, parentPID int, logPath string, wait processWaiter, l launcher) int {
	lg := openHelperLog(logPath)
	defer lg.Close()

	lg.logf("helper start: target=%s new=%s pid=%d", target, newPath, parentPID)

	if _, err := os.Stat(target); err != nil {
		lg.logf("stat target failed: %v", err)
		return 10
	}
	if _, err := os.Stat(newPath); err != nil {
		lg.logf("stat new failed: %v", err)
		return 11
	}

	// Wait for the parent (the running app) to exit so the file is no longer
	// locked. On Windows this is the critical step.
	if parentPID > 0 {
		if err := wait(parentPID, 30*time.Second); err != nil {
			lg.logf("parent did not exit cleanly: %v — proceeding anyway", err)
		}
	}

	backup := target + ".bak"
	_ = os.RemoveAll(backup)

	lg.logf("backing up %s → %s", target, backup)
	if err := copyAny(target, backup); err != nil {
		lg.logf("backup failed: %v", err)
		return 12
	}

	// Retry the swap up to 20 times in case the OS still holds a lock. Each
	// attempt removes the target and renames the new artifact into place.
	swapped := false
	for i := 0; i < 20; i++ {
		if err := os.RemoveAll(target); err != nil {
			lg.logf("remove old (attempt %d): %v", i+1, err)
			time.Sleep(500 * time.Millisecond)
			continue
		}
		if err := os.Rename(newPath, target); err != nil {
			lg.logf("rename new (attempt %d): %v", i+1, err)
			time.Sleep(500 * time.Millisecond)
			continue
		}
		swapped = true
		lg.logf("swap succeeded on attempt %d", i+1)
		break
	}

	if !swapped {
		lg.logf("all swap attempts exhausted — restoring backup")
		if err := restoreFromBackup(backup, target, l); err != nil {
			lg.logf("restore failed: %v", err)
			return 14
		}
		return 13
	}

	if err := l.launch(target); err != nil {
		lg.logf("launch new failed: %v — restoring backup", err)
		if err := restoreFromBackup(backup, target, l); err != nil {
			lg.logf("restore failed: %v", err)
			return 16
		}
		return 15
	}

	// Best-effort backup cleanup. The replaced app is now running.
	if err := os.RemoveAll(backup); err != nil {
		lg.logf("backup cleanup: %v (non-fatal)", err)
	}

	// Tear down the staging directory we received newPath from. The
	// download created it as `wails-update-*` under os.TempDir; after the
	// rename above the directory is empty, but absent this step it would
	// accumulate across update attempts. Guarded by the prefix so we never
	// recursively delete a caller-supplied path that happened to live in a
	// non-temp location.
	stagingDir := filepath.Dir(newPath)
	if strings.HasPrefix(filepath.Base(stagingDir), "wails-update-") {
		if err := os.RemoveAll(stagingDir); err != nil {
			lg.logf("staging cleanup: %v (non-fatal)", err)
		}
	}

	lg.logf("helper done")
	return 0
}

func restoreFromBackup(backup, target string, l launcher) error {
	if err := os.RemoveAll(target); err != nil {
		return fmt.Errorf("remove broken target: %w", err)
	}
	if err := os.Rename(backup, target); err != nil {
		return fmt.Errorf("restore: %w", err)
	}
	return l.launch(target)
}

// copyAny dispatches between file and directory copies. macOS .app bundles
// are directories under the hood so this naturally handles them.
func copyAny(src, dst string) error {
	info, err := os.Stat(src)
	if err != nil {
		return err
	}
	if info.IsDir() {
		return copyTree(src, dst)
	}
	return copyFile(src, dst, info.Mode())
}

func copyFile(src, dst string, mode os.FileMode) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		_ = out.Close()
		return err
	}
	if err := out.Sync(); err != nil {
		_ = out.Close()
		return err
	}
	return out.Close()
}

func copyTree(src, dst string) error {
	info, err := os.Stat(src)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dst, info.Mode()); err != nil {
		return err
	}
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	for _, e := range entries {
		srcPath := filepath.Join(src, e.Name())
		dstPath := filepath.Join(dst, e.Name())
		ei, err := e.Info()
		if err != nil {
			return err
		}
		switch {
		case ei.Mode()&os.ModeSymlink != 0:
			link, err := os.Readlink(srcPath)
			if err != nil {
				return err
			}
			if err := os.Symlink(link, dstPath); err != nil {
				return err
			}
		case ei.IsDir():
			if err := copyTree(srcPath, dstPath); err != nil {
				return err
			}
		default:
			if err := copyFile(srcPath, dstPath, ei.Mode()); err != nil {
				return err
			}
		}
	}
	return nil
}

// waitForPID polls until the named process is no longer alive or the timeout
// elapses. This is portable and avoids platform-specific process-handle APIs.
func waitForPID(pid int, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if !isAlive(pid) {
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
	return fmt.Errorf("pid %d still alive after %s", pid, timeout)
}

func isAlive(pid int) bool {
	if pid <= 0 {
		return false
	}
	return platformIsAlive(pid)
}

// helperLog is a tiny log writer that tolerates missing destinations and
// always writes to stderr too. Failures to log are never fatal.
type helperLog struct {
	w    io.Writer
	file *os.File
}

func openHelperLog(path string) *helperLog {
	if path == "" {
		path = filepath.Join(os.TempDir(), fmt.Sprintf("wails-update-%d.log", os.Getpid()))
	}
	f, err := os.Create(path)
	if err != nil {
		return &helperLog{w: os.Stderr}
	}
	return &helperLog{w: io.MultiWriter(os.Stderr, f), file: f}
}

func (h *helperLog) logf(format string, args ...any) {
	if h == nil || h.w == nil {
		return
	}
	fmt.Fprintf(h.w, "%s: %s\n", time.Now().Format(time.RFC3339), fmt.Sprintf(format, args...))
}

func (h *helperLog) Close() {
	if h != nil && h.file != nil {
		_ = h.file.Close()
	}
}

// errors

var (
	// ErrNotReady is returned by Restart when there is no installed update
	// staged for launch.
	ErrNotReady = errors.New("updater: nothing to restart into (call DownloadAndInstall first)")
)
