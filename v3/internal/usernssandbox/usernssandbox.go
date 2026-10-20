// Package usernssandbox detects when WebKitGTK's process sandbox cannot start
// because the kernel or AppArmor forbids the unprivileged user namespaces that
// bubblewrap (bwrap) and xdg-dbus-proxy need.
//
// On Ubuntu 23.10+/24.04 (and derivatives) AppArmor restricts unprivileged
// user namespaces, so WebKitGTK aborts with SIGTRAP and a message such as
// "bwrap: setting up uid map: Permission denied" the first time a web view is
// created. See https://github.com/wailsapp/wails/issues/6172.
//
// The package is pure Go and has no platform build tags so the logic can be
// unit tested anywhere. It never changes any setting: it only reports.
package usernssandbox

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

// DisableEnvVar is WebKitGTK's own switch that turns its sandbox off. It is
// dangerous: it removes a security boundary, so Wails never sets it.
const DisableEnvVar = "WEBKIT_DISABLE_SANDBOX_THIS_IS_DANGEROUS"

const (
	procRestrictUserns = "/proc/sys/kernel/apparmor_restrict_unprivileged_userns"
	procUsernsClone    = "/proc/sys/kernel/unprivileged_userns_clone"
	procMaxUserNS      = "/proc/sys/user/max_user_namespaces"
)

// Env abstracts the system so detection can be tested.
type Env struct {
	// ReadFile reads a file (os.ReadFile when nil).
	ReadFile func(path string) ([]byte, error)
	// Getenv reads an environment variable (os.Getenv when nil).
	Getenv func(key string) string
	// LookPath finds an executable (exec.LookPath when nil).
	LookPath func(file string) (string, error)
	// Run runs a command and returns its combined output (a real command with a timeout when nil).
	Run func(name string, args ...string) ([]byte, error)
}

func (e Env) readFile(p string) ([]byte, error) {
	if e.ReadFile != nil {
		return e.ReadFile(p)
	}
	return os.ReadFile(p)
}

func (e Env) getenv(k string) string {
	if e.Getenv != nil {
		return e.Getenv(k)
	}
	return os.Getenv(k)
}

func (e Env) lookPath(f string) (string, error) {
	if e.LookPath != nil {
		return e.LookPath(f)
	}
	return exec.LookPath(f)
}

func (e Env) run(name string, args ...string) ([]byte, error) {
	if e.Run != nil {
		return e.Run(name, args...)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return exec.CommandContext(ctx, name, args...).CombinedOutput()
}

// Status is the outcome of Check.
type Status int

const (
	// OK means no problem was found, or the sandbox is not in use.
	OK Status = iota
	// Blocked means bwrap was run and could not create a user namespace.
	Blocked
	// Unknown means restrictions are visible but could not be confirmed
	// (bwrap is not installed, so it cannot be probed).
	Unknown
)

// Result describes what Check found.
type Result struct {
	Status Status
	// Reasons lists the restrictions that were observed.
	Reasons []string
	// ProbeOutput is bwrap's output when the probe failed.
	ProbeOutput string
	// SandboxDisabled is true when DisableEnvVar is set, so WebKit will not use bwrap at all.
	SandboxDisabled bool
}

func readInt(e Env, path string) (int, bool) {
	b, err := e.readFile(path)
	if err != nil {
		return 0, false
	}
	var n int
	if _, err := fmt.Sscanf(strings.TrimSpace(string(b)), "%d", &n); err != nil {
		return 0, false
	}
	return n, true
}

// restrictions returns the sysctl-visible reasons user namespaces may be unavailable.
func restrictions(e Env) []string {
	var r []string
	if v, ok := readInt(e, procRestrictUserns); ok && v == 1 {
		r = append(r, "kernel.apparmor_restrict_unprivileged_userns = 1")
	}
	if v, ok := readInt(e, procUsernsClone); ok && v == 0 {
		r = append(r, "kernel.unprivileged_userns_clone = 0")
	}
	if v, ok := readInt(e, procMaxUserNS); ok && v == 0 {
		r = append(r, "user.max_user_namespaces = 0")
	}
	return r
}

// Check reports whether WebKitGTK's sandbox is likely to fail. It is cheap when
// no restriction is visible: bwrap is only run when a sysctl suggests trouble,
// unless probe is true (used by "wails3 doctor"), which always runs it.
func Check(e Env, probe bool) Result {
	res := Result{}
	if e.getenv(DisableEnvVar) != "" {
		res.SandboxDisabled = true
		return res
	}
	res.Reasons = restrictions(e)
	if len(res.Reasons) == 0 && !probe {
		return res
	}
	if _, err := e.lookPath("bwrap"); err != nil {
		if len(res.Reasons) > 0 {
			res.Status = Unknown
		}
		return res
	}
	out, err := e.run("bwrap", "--unshare-user", "--ro-bind", "/", "/", "true")
	if err != nil {
		var exitErr *exec.ExitError
		// Only a bwrap failure counts; a timeout/start error proves nothing.
		if errors.As(err, &exitErr) || e.Run != nil {
			res.Status = Blocked
			res.ProbeOutput = strings.TrimSpace(string(out))
		}
	}
	return res
}

// Advice is actionable text for a Blocked or Unknown result.
func Advice() string {
	return `WebKitGTK's sandbox (bubblewrap) cannot create a user namespace, so the web view will crash on start (SIGTRAP, "bwrap: setting up uid map: Permission denied").
This is typically AppArmor restricting unprivileged user namespaces (Ubuntu 23.10+/24.04). To fix, either:
  - allow your app (or bwrap) via an AppArmor profile with "userns," permission, e.g. /etc/apparmor.d/<app>:
        abi <abi/4.0>,
        include <tunables/global>
        profile myapp /path/to/myapp flags=(unconfined) {
          userns,
          include if exists <local/myapp>
        }
    then run: sudo apparmor_parser -r /etc/apparmor.d/<app>
  - or, system-wide (weakens a host protection): sudo sysctl -w kernel.apparmor_restrict_unprivileged_userns=0
As a last resort for local development only, set ` + DisableEnvVar + `=1 to run WebKit WITHOUT its sandbox. Wails never sets this for you.
See https://github.com/wailsapp/wails/issues/6172`
}

// Message is a one-paragraph summary for a result, or "" when there is nothing to say.
func (r Result) Message() string {
	switch r.Status {
	case Blocked:
		s := "bwrap could not create a user namespace"
		if r.ProbeOutput != "" {
			s += " (" + r.ProbeOutput + ")"
		}
		if len(r.Reasons) > 0 {
			s += "; observed: " + strings.Join(r.Reasons, ", ")
		}
		return s
	case Unknown:
		return "user namespaces look restricted (" + strings.Join(r.Reasons, ", ") + ") and bwrap is not installed to confirm"
	}
	return ""
}
