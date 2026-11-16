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
	envHelperReady  = "WAILS_UPDATER_HELPER_READY"  // readiness status file path
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
	readyPath := os.Getenv(envHelperReady)
	if target == "" || newPath == "" {
		_ = signalHelperStatus(readyPath, "error: missing helper target or staged artifact")
		os.Exit(2)
	}
	if _, err := os.Stat(target); err != nil {
		_ = signalHelperStatus(readyPath, "error: target is unavailable: "+err.Error())
		os.Exit(10)
	}
	if _, err := os.Stat(newPath); err != nil {
		_ = signalHelperStatus(readyPath, "error: staged artifact is unavailable: "+err.Error())
		os.Exit(11)
	}
	if err := signalHelperStatus(readyPath, "ready"); err != nil {
		os.Exit(3)
	}
	pid, _ := strconv.Atoi(os.Getenv(envHelperPID))
	logPath := os.Getenv(envHelperLog)

	code := runHelperSwap(target, newPath, pid, logPath, waitForPID, osLauncher{})
	os.Exit(code)
}

func signalHelperStatus(path, status string) error {
	if path == "" {
		return nil
	}
	return os.WriteFile(path, []byte(status), 0o600)
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

// swapRetryDelay is the pause between replace attempts. Held in a var so
// focused tests can shorten the retry window without sleeping for 10s.
var swapRetryDelay = 500 * time.Millisecond

// runHelperSwap implements the actual swap logic. It is unexported and
// dependency-injected so unit tests can drive every branch without process
// spawning. Returns the exit code the helper should use:
//
//	 0  swapped and relaunched
//	10  target missing          11  staged payload missing
//	12  backup failed (original untouched and relaunched)
//	13  swap failed; the intact original relaunched, or the backup restored
//	14  swap failed AND restore failed — the install is broken
//	15  new application would not launch, backup restored
//	16  launch failed AND restore failed — the install is broken
//	17  parent never exited; nothing was touched, original relaunched
//	18  staged payload was not installable; original untouched and relaunched
//	19  installed application failed validation, backup restored
//	20  validation failed AND restore failed — the install is broken
//
// Only 14, 16 and 20 leave the user without a working application, and each is
// a restore that itself failed after an earlier failure. Every other code has
// either changed nothing or rolled back to what was there before.
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
	//
	// If the parent never exits — most commonly because a webview dialog is
	// blocking shutdown — we must NOT proceed: on Windows the swap would
	// grind against the lock; on macOS the subsequent `open -n` would launch
	// a second instance alongside the still-running one, leaving the user
	// with two copies of the app racing on shared state.
	if parentPID > 0 {
		if err := wait(parentPID, 30*time.Second); err != nil {
			lg.logf("parent did not exit within timeout: %v — aborting swap", err)
			// The target was never touched. Leaving the user with no
			// running application is worse than briefly overlapping with a
			// parent that is still shutting down, so relaunch the intact
			// original and let the parent finish exiting.
			clearHelperEnv()
			if err := l.launch(target); err != nil {
				lg.logf("relaunch original failed: %v", err)
			} else {
				lg.logf("relaunched original after parent-timeout abort")
			}
			return 17
		}
	}

	// The parent has exited. Clear helper mode before backup or replacement
	// so both the new application and any recovered original boot normally.
	clearHelperEnv()

	// Re-validate the payload now the wait is over. The stat at the top of this
	// function ran before the parent quit, and the parent spends that interval
	// running its own shutdown — which is ample time for it to have torn the
	// staging directory down. Because a staged .app is a directory, an emptied
	// one still renames cleanly, so nothing downstream would notice. The target
	// has not been touched yet, so there is nothing to undo: relaunch what is
	// already installed and leave it alone.
	if err := validatePayloadForTarget(newPath, target); err != nil {
		lg.logf("staged update is not installable: %v — aborting swap", err)
		if err := l.launch(target); err != nil {
			lg.logf("relaunch original failed: %v", err)
		}
		return 18
	}

	backup := target + ".bak"
	_ = os.RemoveAll(backup)

	lg.logf("backing up %s → %s", target, backup)
	if err := copyAny(target, backup); err != nil {
		lg.logf("backup failed: %v", err)
		// The original is untouched; relaunch it instead of restoring a
		// potentially incomplete backup.
		if err := l.launch(target); err != nil {
			lg.logf("relaunch original failed: %v", err)
		}
		return 12
	}

	// Snapshot the original mode before the swap. The downloaded artifact
	// was created via os.Create which sets 0o666-minus-umask — on Unix
	// that drops the executable bit, so a direct rename would leave the
	// new binary non-executable and exec() would fail with EACCES.
	var origMode os.FileMode
	if origInfo, err := os.Stat(target); err == nil {
		origMode = origInfo.Mode()
	}

	// Retry the swap up to 20 times in case the OS still holds a transient
	// lock on the target. The actual replace strategy is platform-specific —
	// Unix unlinks and renames in place (open file handles keep working
	// because the inode lives on); Windows moves the target aside because
	// it can't delete a recently-running executable but it CAN rename it,
	// even with the loader still mapping its image.
	swapped := false
	for i := 0; i < 20; i++ {
		// replaceTarget re-validates the payload on every attempt, immediately
		// before it disturbs the target. A payload that has stopped being
		// installable will not recover, so that ends the loop at once.
		if err := replaceTarget(target, newPath); err != nil {
			lg.logf("replace (attempt %d): %v", i+1, err)
			if errors.Is(err, errNotInstallable) {
				break
			}
			time.Sleep(swapRetryDelay)
			continue
		}
		swapped = true
		lg.logf("swap succeeded on attempt %d", i+1)
		break
	}

	// Restore the original executable bits on the new binary.
	if swapped && origMode != 0 {
		if err := os.Chmod(target, origMode.Perm()); err != nil {
			lg.logf("chmod restored mode: %v (non-fatal)", err)
		}
	}

	if !swapped {
		// Every replace attempt left the original executable in place: the
		// Windows rename-aside is reverted on failure, and the Unix file
		// rename only takes effect on success. Calling restoreFromBackup
		// here would first os.RemoveAll(target), which fails on the very
		// lock that blocked the swap (and, on Unix .app bundles, needlessly
		// deletes an intact original). Relaunch the original instead.
		lg.logf("all swap attempts exhausted — relaunching intact original")
		if _, err := os.Stat(target); err == nil {
			if err := l.launch(target); err == nil {
				if err := os.RemoveAll(backup); err != nil {
					lg.logf("backup cleanup: %v (non-fatal)", err)
				}
				return 13
			} else {
				lg.logf("relaunch original failed: %v — restoring backup", err)
			}
		}
		// The target slot is missing (e.g. a directory swap cleared it
		// before the move failed); fall back to the backup.
		lg.logf("all swap attempts exhausted — restoring backup")
		if err := restoreFromBackup(backup, target, l); err != nil {
			lg.logf("restore failed: %v", err)
			return 14
		}
		return 13
	}

	// The payload is now installed, and the backup is still standing. This is
	// the last moment at which the old application can be recovered, so it is
	// where the new one has to prove itself — structurally, then against the
	// signature, which covers every sealed file and so catches damage the
	// layout check cannot see. The backup is the yardstick for the signature:
	// a check the original fails is not one the update is held to.
	if err := validateInstalled(target, backup, lg); err != nil {
		lg.logf("installed application failed validation: %v — restoring backup", err)
		if err := restoreFromBackup(backup, target, l); err != nil {
			lg.logf("restore failed: %v", err)
			return 20
		}
		return 19
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

// clearHelperEnv unsets every WAILS_UPDATER_HELPER_* variable in the current
// process. Called after the parent exits and before backup or replacement,
// so both the new application and any recovered original boot normally.
func clearHelperEnv() {
	for _, k := range []string{envHelperMode, envHelperTarget, envHelperNew, envHelperPID, envHelperLog, envHelperReady} {
		_ = os.Unsetenv(k)
	}
}

// errors

var (
	// ErrNotReady is returned by Restart when there is no installed update
	// staged for launch.
	ErrNotReady = errors.New("updater: nothing to restart into (call DownloadAndInstall first)")
	// ErrHelperNotReady is returned when the replacement helper does not
	// acknowledge that it reached helper mode before the application exits.
	ErrHelperNotReady = errors.New("updater: helper did not signal readiness")
	// ErrJobBreakawayDenied is returned on Windows when the helper cannot be
	// created outside a parent Job that may terminate it with the application.
	ErrJobBreakawayDenied = errors.New("updater: helper could not break away from the parent Windows Job")
)
