package git

import (
	"bytes"
	"errors"
	"fmt"
	"net/url"
	"os/exec"
	"strings"
)

// ErrNotInstalled is returned when git is not found in PATH.
var ErrNotInstalled = errors.New("git is not installed; please install git from https://git-scm.com")

func isNotFound(err error) bool {
	var execErr *exec.Error
	return errors.As(err, &execErr) && errors.Is(execErr.Err, exec.ErrNotFound)
}

// redactArgs returns a copy of args with any URL credentials (user:pass@host)
// replaced by user:***@host so tokens are not leaked in error messages.
func redactArgs(args []string) []string {
	out := make([]string, len(args))
	for i, a := range args {
		if u, err := url.Parse(a); err == nil && u.User != nil {
			u.User = url.UserPassword(u.User.Username(), "***")
			a = u.String()
		}
		out[i] = a
	}
	return out
}

func run(args ...string) error {
	out, err := exec.Command("git", args...).CombinedOutput()
	if err != nil {
		if isNotFound(err) {
			return ErrNotInstalled
		}
		return fmt.Errorf("git %s: %w\n%s", strings.Join(redactArgs(args), " "), err, bytes.TrimSpace(out))
	}
	return nil
}

func output(args ...string) (string, error) {
	out, err := exec.Command("git", args...).CombinedOutput()
	if err != nil {
		if isNotFound(err) {
			return "", ErrNotInstalled
		}
		return "", fmt.Errorf("git %s: %w\n%s", strings.Join(redactArgs(args), " "), err, bytes.TrimSpace(out))
	}
	return strings.TrimSpace(string(out)), nil
}

// HeadHash returns the short (8-character) commit hash of HEAD in dir.
func HeadHash(dir string) (string, error) {
	hash, err := output("-C", dir, "rev-parse", "HEAD")
	if err != nil {
		return "", err
	}
	if len(hash) < 8 {
		return "", fmt.Errorf("git rev-parse returned unexpected output %q", hash)
	}
	return hash[:8], nil
}

// Clone clones url into dir. If tag is non-empty, checks out that tag or branch.
func Clone(url, dir, tag string) error {
	args := []string{"clone", "--quiet"}
	if tag != "" {
		args = append(args, "--branch", tag)
	}
	args = append(args, url, dir)
	return run(args...)
}

// Init initializes a new git repository at dir.
func Init(dir string) error {
	return run("-C", dir, "init", "--quiet")
}

// RemoteAdd adds a named remote to the repository at dir.
func RemoteAdd(dir, name, url string) error {
	return run("-C", dir, "remote", "add", name, url)
}

// AddAll stages all files in the repository at dir.
func AddAll(dir string) error {
	return run("-C", dir, "add", ".")
}
