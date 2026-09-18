package setupwizard

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// notarizeTimeout bounds how long the wizard waits for the terminal window to
// finish. notarytool contacts Apple after the password is entered, so the
// window can legitimately sit there for a while, but it should never wait for
// a window the user has walked away from.
const notarizeTimeout = 15 * time.Minute

// notarizePollInterval is how often the wizard checks the status file written
// by the spawned terminal.
const notarizePollInterval = 500 * time.Millisecond

// Job states reported to the frontend.
const (
	notarizeStateIdle      = "idle"
	notarizeStateRunning   = "running"
	notarizeStateSucceeded = "succeeded"
	notarizeStateFailed    = "failed"
)

// notarizeJob tracks a single `xcrun notarytool store-credentials` run that the
// wizard has handed off to a terminal window of its own. The password is typed
// into that window, so it never reaches the browser, the wizard's HTTP API or
// the process arguments of any command.
type notarizeJob struct {
	command     string // human-readable command, shown in the browser
	dir         string // temp dir holding the script and status file
	profileName string
	appleID     string
	teamID      string

	mu    sync.Mutex
	state string
	err   string
}

// statusPath is the file the spawned script writes notarytool's exit code to.
// Its appearance is what tells the wizard the terminal window is done.
func (j *notarizeJob) statusPath() string {
	return filepath.Join(j.dir, "status")
}

// scriptPath is the script handed to Terminal.app. The .command extension is
// the macOS convention for a shell script Terminal runs on open.
func (j *notarizeJob) scriptPath() string {
	return filepath.Join(j.dir, "notarize.command")
}

// finish records the job's outcome. Only the first call wins, so a cancel and a
// completing terminal window can race without clobbering each other.
func (j *notarizeJob) finish(state, errMsg string) bool {
	j.mu.Lock()
	defer j.mu.Unlock()
	if j.state != notarizeStateRunning {
		return false
	}
	j.state = state
	j.err = errMsg
	return true
}

// snapshot reads the job's current state and error message together.
func (j *notarizeJob) snapshot() (state, errMsg string) {
	j.mu.Lock()
	defer j.mu.Unlock()
	return j.state, j.err
}

// matches reports whether the job was started for exactly these credentials, so
// a repeated create can be answered with the window that is already open rather
// than quietly adopting it for different values.
func (j *notarizeJob) matches(profileName, appleID, teamID string) bool {
	return j.profileName == profileName && j.appleID == appleID && j.teamID == teamID
}

// abandon drops a job the wizard never adopted — cancelled while its window was
// still opening. The window itself is left for the user to close.
func (j *notarizeJob) abandon() {
	j.finish(notarizeStateFailed, "Cancelled")
	os.RemoveAll(j.dir)
}

// notarizeArgs builds the xcrun invocation. The app-specific password is
// deliberately absent: notarytool prompts for it securely on the terminal it is
// attached to, which keeps it out of `ps` and `/proc`.
func notarizeArgs(profileName, appleID, teamID string) []string {
	return []string{
		"xcrun", "notarytool", "store-credentials",
		profileName,
		"--apple-id", appleID,
		"--team-id", teamID,
		"--validate",
	}
}

// shellSafe matches arguments that need no quoting when shown or pasted into a
// shell.
func shellSafe(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9':
		case strings.ContainsRune("@%_+=:,./-", r):
		default:
			return false
		}
	}
	return true
}

// shellQuote renders s as a single /bin/sh word.
func shellQuote(s string) string {
	if shellSafe(s) {
		return s
	}
	return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
}

// shellJoin renders args as a copy-pasteable command line.
func shellJoin(args []string) string {
	quoted := make([]string, len(args))
	for i, arg := range args {
		quoted[i] = shellQuote(arg)
	}
	return strings.Join(quoted, " ")
}

// notarizeFieldClean rejects values carrying control characters. They are always
// quoted before they reach the shell, so this is about keeping the command the
// user is shown readable and honest rather than about injection.
func notarizeFieldClean(s string) (string, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return "", false
	}
	for _, r := range s {
		if r < 0x20 || r == 0x7f {
			return "", false
		}
	}
	return s, true
}

// notarizeScript builds the script the new terminal window runs. It shows the
// command before running it, so the user can see exactly what is being asked of
// their machine, then writes notarytool's exit code out as its final act — the
// wizard treats that file appearing as "the window is finished".
func notarizeScript(args []string, statusPath string) string {
	var b strings.Builder
	b.WriteString("#!/bin/sh\n")
	b.WriteString("# Created by `wails3 setup` to store your Apple notarization credentials.\n")
	b.WriteString("# It is removed automatically once the setup wizard has read the result.\n\n")
	b.WriteString("printf '\\n  Wails · macOS notarization setup\\n'\n")
	b.WriteString("printf '  ================================\\n\\n'\n")
	b.WriteString("printf '  Enter your app-specific password at the prompt below.\\n'\n")
	b.WriteString("printf '  It is not echoed as you type and is never passed on the\\n'\n")
	b.WriteString("printf '  command line, so it cannot be read out of the process list.\\n\\n'\n")
	b.WriteString("printf '  Running:\\n\\n'\n")
	fmt.Fprintf(&b, "printf '    %%s\\n\\n' %s\n\n", shellQuote(shellJoin(args)))
	b.WriteString(shellJoin(args) + "\n")
	b.WriteString("__wails_status=$?\n\n")
	b.WriteString("if [ \"$__wails_status\" -eq 0 ]; then\n")
	b.WriteString("\tprintf '\\n  Credentials stored. Head back to the setup wizard in your browser.\\n\\n'\n")
	b.WriteString("else\n")
	b.WriteString("\tprintf '\\n  notarytool exited with status %s. Read the error above, then head\\n' \"$__wails_status\"\n")
	b.WriteString("\tprintf '  back to the setup wizard in your browser to try again.\\n\\n'\n")
	b.WriteString("fi\n\n")
	b.WriteString("# Written last: the wizard is watching for this file.\n")
	fmt.Fprintf(&b, "printf '%%s' \"$__wails_status\" > %s && mv %s %s\n",
		shellQuote(statusPath+".tmp"), shellQuote(statusPath+".tmp"), shellQuote(statusPath))
	return b.String()
}

// startNotarizeJob writes the script and opens it in a new terminal window.
func startNotarizeJob(profileName, appleID, teamID string) (*notarizeJob, error) {
	dir, err := os.MkdirTemp("", "wails-notarize-")
	if err != nil {
		return nil, fmt.Errorf("could not create a working directory: %w", err)
	}

	args := notarizeArgs(profileName, appleID, teamID)
	job := &notarizeJob{
		command:     shellJoin(args),
		dir:         dir,
		profileName: profileName,
		appleID:     appleID,
		teamID:      teamID,
		state:       notarizeStateRunning,
	}

	if err := os.WriteFile(job.scriptPath(), []byte(notarizeScript(args, job.statusPath())), 0o700); err != nil {
		os.RemoveAll(dir)
		return nil, fmt.Errorf("could not write the setup script: %w", err)
	}

	// `open -a Terminal` both creates the window and brings it in front of the
	// browser, which is the point: the prompt comes to the user rather than the
	// user having to find the terminal the wizard was started from.
	if output, err := exec.Command("open", "-a", "Terminal", job.scriptPath()).CombinedOutput(); err != nil {
		os.RemoveAll(dir)
		msg := strings.TrimSpace(string(output))
		if msg == "" {
			msg = err.Error()
		}
		return nil, fmt.Errorf("could not open a terminal window: %s", msg)
	}

	return job, nil
}

// watchNotarizeJob waits for the terminal window to report an exit code, and
// records the outcome. It gives up at notarizeTimeout, or when the wizard shuts
// down.
func (w *Wizard) watchNotarizeJob(job *notarizeJob) {
	defer os.RemoveAll(job.dir)

	ticker := time.NewTicker(notarizePollInterval)
	defer ticker.Stop()
	deadline := time.After(notarizeTimeout)

	for {
		if data, err := os.ReadFile(job.statusPath()); err == nil {
			code := strings.TrimSpace(string(data))
			if code == "0" {
				w.completeNotarizeJob(job)
			} else {
				job.finish(notarizeStateFailed, "notarytool exited with status "+code+". Check the terminal window for details.")
			}
			return
		}

		select {
		case <-ticker.C:
		case <-deadline:
			job.finish(notarizeStateFailed, "Timed out waiting for the terminal window. Close it and try again.")
			return
		case <-w.shutdown:
			job.finish(notarizeStateFailed, "The setup wizard was closed before the terminal window finished.")
			return
		}
	}
}

// completeNotarizeJob persists the profile the terminal window just created.
func (w *Wizard) completeNotarizeJob(job *notarizeJob) {
	// The credentials are in the keychain either way, but a job the user has
	// cancelled is no longer the wizard's to record: a later run may already
	// have written different defaults.
	if state, _ := job.snapshot(); state != notarizeStateRunning {
		return
	}

	defs, err := LoadGlobalDefaults()
	if err != nil {
		job.finish(notarizeStateFailed, "Profile created but failed to load defaults: "+err.Error())
		return
	}
	defs.Signing.Darwin.KeychainProfile = job.profileName
	defs.Signing.Darwin.TeamID = job.teamID
	if err := SaveGlobalDefaults(defs); err != nil {
		job.finish(notarizeStateFailed, "Profile created but failed to save defaults: "+err.Error())
		return
	}
	job.finish(notarizeStateSucceeded, "")
}
