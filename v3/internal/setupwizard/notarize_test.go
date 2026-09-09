package setupwizard

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

func TestNotarizeArgsOmitPassword(t *testing.T) {
	args := notarizeArgs("wails-notary", "me@example.com", "TEAM123")
	joined := strings.Join(args, " ")

	for _, banned := range []string{"--password", "--password-stdin"} {
		if strings.Contains(joined, banned) {
			t.Fatalf("notarizeArgs must not pass %s, got %q", banned, joined)
		}
	}
	for _, want := range []string{"xcrun", "notarytool", "store-credentials", "wails-notary", "--apple-id", "me@example.com", "--team-id", "TEAM123", "--validate"} {
		if !strings.Contains(joined, want) {
			t.Fatalf("notarizeArgs missing %q, got %q", want, joined)
		}
	}
}

func TestShellJoin(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{"plain args are left alone", []string{"xcrun", "notarytool"}, "xcrun notarytool"},
		{"email and team id need no quotes", []string{"--apple-id", "me@example.com"}, "--apple-id me@example.com"},
		{"spaces are quoted", []string{"echo", "my profile"}, "echo 'my profile'"},
		{"single quotes are escaped", []string{"echo", "it's"}, `echo 'it'\''s'`},
		{"empty args are quoted", []string{"echo", ""}, "echo ''"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shellJoin(tt.args); got != tt.want {
				t.Fatalf("shellJoin(%q) = %q, want %q", tt.args, got, tt.want)
			}
		})
	}
}

func TestNotarizeFieldClean(t *testing.T) {
	tests := []struct {
		in   string
		want string
		ok   bool
	}{
		{"wails-notary", "wails-notary", true},
		{"  wails-notary  ", "wails-notary", true},
		{"", "", false},
		{"   ", "", false},
		{"two\nlines", "", false},
		{"bell\a", "", false},
	}

	for _, tt := range tests {
		got, ok := notarizeFieldClean(tt.in)
		if got != tt.want || ok != tt.ok {
			t.Fatalf("notarizeFieldClean(%q) = (%q, %v), want (%q, %v)", tt.in, got, ok, tt.want, tt.ok)
		}
	}
}

// requirePOSIXShell skips tests that run the generated script. The script is
// /bin/sh and only ever runs on macOS, so there is nothing to learn from
// wrestling it through whichever sh a Windows runner happens to have.
func requirePOSIXShell(t *testing.T) {
	t.Helper()

	if runtime.GOOS == "windows" {
		t.Skip("the notarization script is POSIX sh and only runs on macOS")
	}
	if _, err := exec.LookPath("sh"); err != nil {
		t.Skip("sh not available")
	}
}

// stubXcrun puts a fake `xcrun` on PATH that records its arguments and exits
// with the given code, so the generated script can be run for real.
func stubXcrun(t *testing.T, dir string, exitCode int) string {
	t.Helper()

	binDir := filepath.Join(dir, "bin")
	if err := os.MkdirAll(binDir, 0o700); err != nil {
		t.Fatal(err)
	}
	argsFile := filepath.Join(dir, "args")
	stub := "#!/bin/sh\nprintf '%s\\n' \"$@\" > '" + argsFile + "'\nexit " + strconv.Itoa(exitCode) + "\n"
	if err := os.WriteFile(filepath.Join(binDir, "xcrun"), []byte(stub), 0o700); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", binDir+string(os.PathListSeparator)+os.Getenv("PATH"))
	return argsFile
}

func runNotarizeScript(t *testing.T, dir string, args []string) (statusPath string) {
	t.Helper()

	requirePOSIXShell(t)

	statusPath = filepath.Join(dir, "status")
	scriptPath := filepath.Join(dir, "notarize.command")
	if err := os.WriteFile(scriptPath, []byte(notarizeScript(args, statusPath)), 0o700); err != nil {
		t.Fatal(err)
	}

	if out, err := exec.Command("sh", "-n", scriptPath).CombinedOutput(); err != nil {
		t.Fatalf("generated script is not valid sh: %v\n%s", err, out)
	}

	out, err := exec.Command("sh", scriptPath).CombinedOutput()
	if err != nil {
		t.Fatalf("script failed to run: %v\n%s", err, out)
	}
	if want := shellJoin(args); !strings.Contains(string(out), want) {
		t.Fatalf("script output does not show the command %q:\n%s", want, out)
	}
	return statusPath
}

func TestNotarizeScriptReportsSuccess(t *testing.T) {
	dir := t.TempDir()
	argsFile := stubXcrun(t, dir, 0)
	args := notarizeArgs("wails-notary", "me@example.com", "TEAM123")

	statusPath := runNotarizeScript(t, dir, args)

	status, err := os.ReadFile(statusPath)
	if err != nil {
		t.Fatalf("no status file written: %v", err)
	}
	if string(status) != "0" {
		t.Fatalf("status = %q, want %q", status, "0")
	}

	recorded, err := os.ReadFile(argsFile)
	if err != nil {
		t.Fatalf("stub xcrun was not run: %v", err)
	}
	got := strings.Split(strings.TrimSpace(string(recorded)), "\n")
	want := args[1:] // the stub records everything after `xcrun`
	if strings.Join(got, "\x00") != strings.Join(want, "\x00") {
		t.Fatalf("xcrun called with %q, want %q", got, want)
	}
}

func TestNotarizeScriptReportsFailure(t *testing.T) {
	dir := t.TempDir()
	stubXcrun(t, dir, 1)

	statusPath := runNotarizeScript(t, dir, notarizeArgs("wails-notary", "me@example.com", "TEAM123"))

	status, err := os.ReadFile(statusPath)
	if err != nil {
		t.Fatalf("no status file written: %v", err)
	}
	if string(status) != "1" {
		t.Fatalf("status = %q, want %q", status, "1")
	}
}

func TestNotarizeScriptQuotesAwkwardValues(t *testing.T) {
	dir := t.TempDir()
	argsFile := stubXcrun(t, dir, 0)
	args := notarizeArgs("my profile; touch /tmp/pwned", "me@example.com", "it's mine")

	runNotarizeScript(t, dir, args)

	recorded, err := os.ReadFile(argsFile)
	if err != nil {
		t.Fatalf("stub xcrun was not run: %v", err)
	}
	got := strings.Split(strings.TrimSpace(string(recorded)), "\n")
	want := args[1:]
	if strings.Join(got, "\x00") != strings.Join(want, "\x00") {
		t.Fatalf("xcrun called with %q, want %q", got, want)
	}
}

func TestNotarizeJobFinishesOnce(t *testing.T) {
	job := &notarizeJob{state: notarizeStateRunning}

	if !job.finish(notarizeStateSucceeded, "") {
		t.Fatal("first finish should win")
	}
	if job.finish(notarizeStateFailed, "Cancelled") {
		t.Fatal("second finish should be ignored")
	}
	if state, errMsg := job.snapshot(); state != notarizeStateSucceeded || errMsg != "" {
		t.Fatalf("job = (%q, %q), want (%q, \"\")", state, errMsg, notarizeStateSucceeded)
	}
}

func TestNotarizeJobMatches(t *testing.T) {
	job := &notarizeJob{profileName: "wails-notary", appleID: "me@example.com", teamID: "TEAM123"}

	if !job.matches("wails-notary", "me@example.com", "TEAM123") {
		t.Fatal("the credentials the job was started for should match")
	}
	for _, tt := range []struct {
		name                         string
		profileName, appleID, teamID string
	}{
		{"different profile", "other-notary", "me@example.com", "TEAM123"},
		{"different apple id", "wails-notary", "someone@example.com", "TEAM123"},
		{"different team id", "wails-notary", "me@example.com", "OTHER99"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if job.matches(tt.profileName, tt.appleID, tt.teamID) {
				t.Fatal("a create for different credentials must not match the running job")
			}
		})
	}
}

func TestNotarizeJobAbandon(t *testing.T) {
	dir := t.TempDir()
	jobDir := filepath.Join(dir, "job")
	if err := os.MkdirAll(jobDir, 0o700); err != nil {
		t.Fatal(err)
	}
	job := &notarizeJob{dir: jobDir, state: notarizeStateRunning}

	job.abandon()

	if state, errMsg := job.snapshot(); state != notarizeStateFailed || errMsg != "Cancelled" {
		t.Fatalf("job = (%q, %q), want (%q, %q)", state, errMsg, notarizeStateFailed, "Cancelled")
	}
	if _, err := os.Stat(jobDir); !os.IsNotExist(err) {
		t.Fatalf("abandon should remove the job directory, stat err = %v", err)
	}
}
