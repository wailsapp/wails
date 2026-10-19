package git

import (
	"bytes"
	"errors"
	"fmt"
	"net/url"
	"os/exec"
	"regexp"
	"strings"
	"unicode"
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

// ErrInvalidRemoteURL is returned when a URL is not an accepted remote git URL.
var ErrInvalidRemoteURL = errors.New("invalid remote git URL")

// remoteSchemes lists the URL schemes accepted for remote clones.
var remoteSchemes = map[string]bool{
	"https": true,
	"http":  true,
	"ssh":   true,
	"git":   true,
}

// remoteProtocolArgs restricts git to the transports accepted by
// ValidateRemoteURL, so helpers such as ext:: or file:// are refused by git
// itself regardless of the user's git configuration.
var remoteProtocolArgs = []string{
	"-c", "protocol.allow=never",
	"-c", "protocol.https.allow=always",
	"-c", "protocol.http.allow=always",
	"-c", "protocol.ssh.allow=always",
	"-c", "protocol.git.allow=always",
}

var (
	// transportHelperFormat matches git's <transport>::<address> syntax.
	transportHelperFormat = regexp.MustCompile(`^[A-Za-z0-9+.-]+::`)
	// scpLikeFormat matches [user@]host:path. The path must not start with ':'.
	scpLikeFormat = regexp.MustCompile(`^(?:[^@/:]+@)?([^@/:]+):([^:].*)$`)
)

// ValidateRemoteURL reports whether url is a remote git URL that is safe to
// pass to git clone. Accepted forms are https://, http://, ssh:// and git://
// URLs and scp-like [user@]host:path addresses. Options (values starting with
// '-'), transport helpers (<transport>::<address>), file:// URLs and local
// paths are rejected.
func ValidateRemoteURL(rawURL string) error {
	// The URL is left out of the error: its userinfo may carry a token.
	invalid := func(reason string) error {
		return fmt.Errorf("%w: %s", ErrInvalidRemoteURL, reason)
	}

	if rawURL == "" {
		return invalid("empty URL")
	}
	if strings.HasPrefix(rawURL, "-") {
		return invalid("must not start with '-'")
	}
	if strings.IndexFunc(rawURL, func(r rune) bool { return unicode.IsSpace(r) || unicode.IsControl(r) }) != -1 {
		return invalid("must not contain whitespace or control characters")
	}
	if transportHelperFormat.MatchString(rawURL) {
		return invalid("transport helpers are not supported")
	}

	if strings.Contains(rawURL, "://") {
		u, err := url.Parse(rawURL)
		if err != nil {
			// url.Error repeats the input URL; report only the cause.
			var urlErr *url.Error
			if errors.As(err, &urlErr) {
				err = urlErr.Err
			}
			return invalid(err.Error())
		}
		if !remoteSchemes[strings.ToLower(u.Scheme)] {
			return invalid(fmt.Sprintf("unsupported scheme %q (use https, http, ssh or git)", u.Scheme))
		}
		if u.Hostname() == "" {
			return invalid("missing host")
		}
		if strings.HasPrefix(u.Hostname(), "-") {
			return invalid("host must not start with '-'")
		}
		return nil
	}

	match := scpLikeFormat.FindStringSubmatch(rawURL)
	if match == nil {
		return invalid("expected an https://, http://, ssh:// or git:// URL, or user@host:path")
	}
	host := match[1]
	// A single letter before ':' is a Windows drive letter, not a host.
	if len(host) < 2 {
		return invalid("not a remote host")
	}
	if strings.HasPrefix(host, "-") {
		return invalid("host must not start with '-'")
	}
	return nil
}

// cloneArgs builds the arguments for git clone, rejecting values that git
// would parse as options.
func cloneArgs(url, dir, tag string) ([]string, error) {
	for name, value := range map[string]string{"url": url, "directory": dir, "tag": tag} {
		if strings.HasPrefix(value, "-") {
			return nil, fmt.Errorf("invalid git clone %s: must not start with '-'", name)
		}
	}
	args := []string{"clone", "--quiet"}
	if tag != "" {
		args = append(args, "--branch", tag)
	}
	return append(args, "--", url, dir), nil
}

// Clone clones url into dir. If tag is non-empty, checks out that tag or branch.
func Clone(url, dir, tag string) error {
	args, err := cloneArgs(url, dir, tag)
	if err != nil {
		return err
	}
	return run(args...)
}

// CloneRemote clones the remote repository at url into dir. The URL is checked
// with ValidateRemoteURL and git is restricted to the https, http, ssh and git
// transports. If tag is non-empty, checks out that tag or branch.
func CloneRemote(url, dir, tag string) error {
	if err := ValidateRemoteURL(url); err != nil {
		return err
	}
	args, err := cloneArgs(url, dir, tag)
	if err != nil {
		return err
	}
	return run(append(append([]string{}, remoteProtocolArgs...), args...)...)
}

// Init initializes a new git repository at dir.
func Init(dir string) error {
	return run("-C", dir, "init", "--quiet")
}

// RemoteAdd adds a named remote to the repository at dir.
func RemoteAdd(dir, name, url string) error {
	if strings.HasPrefix(name, "-") || strings.HasPrefix(url, "-") {
		return fmt.Errorf("invalid git remote %q: name and URL must not start with '-'", name)
	}
	return run("-C", dir, "remote", "add", "--", name, url)
}

// AddAll stages all files in the repository at dir.
func AddAll(dir string) error {
	return run("-C", dir, "add", ".")
}
