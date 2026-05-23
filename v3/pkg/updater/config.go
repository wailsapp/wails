package updater

import (
	"errors"
	"time"
)

// Config configures the Updater. Pass it to Updater.Init.
type Config struct {
	// CurrentVersion is the version currently running. Required.
	// Pass the same string you use to tag releases (e.g. "1.2.3" — no "v" prefix).
	CurrentVersion string

	// Providers are the update sources, tried in order. The first to return
	// a release is used; if a provider returns (nil, nil) the Updater treats
	// the application as up to date and stops walking the chain (fallback is
	// for "primary unreachable", not "providers disagree"). Required.
	Providers []Provider

	// PublicKey is the trust root used to verify release signatures. It is
	// the ONLY trust anchor for signature verification — the release source
	// cannot substitute its own key, since the whole point of pinning a key
	// here is to bind verification to a value the application developer set
	// at build time and the release feed cannot influence.
	//
	// When unset, releases that carry only a Digest still install (the digest
	// is checked against the streaming hash), but any release that carries a
	// Signature is rejected. Strongly recommended for any app whose
	// distribution channel might be compromised or proxied.
	PublicKey []byte

	// CheckInterval, when non-zero, makes the Updater poll providers on a
	// timer in the background. Pop-up-on-found behaviour is the same as a
	// manual Check finding an update.
	CheckInterval time.Duration

	// Platform / Arch / Channel override the per-platform defaults passed to
	// each Provider's Check. Leave empty to use runtime.GOOS / runtime.GOARCH
	// and the provider's default channel.
	Platform string
	Arch     string
	Channel  string

	// Window controls how the update UI is rendered. Not yet wired in v1 of
	// the package; see the upcoming Window option types.
	Window WindowOption
}

// WindowOption is the marker interface for the Window slot on Config. The set
// of concrete types is small and closed: BuiltinWindow, application.Window
// (BYO), or the sentinel WindowNone. The full window machinery is implemented
// in a follow-up commit; the interface is declared here so Config does not
// change shape between commits.
type WindowOption interface {
	isWindowOption()
}

// validate returns an error describing the first problem found in c, or nil.
func (c *Config) validate() error {
	if c == nil {
		return errors.New("updater: Config is nil")
	}
	if c.CurrentVersion == "" {
		return errors.New("updater: Config.CurrentVersion is required")
	}
	if len(c.Providers) == 0 {
		return errors.New("updater: Config.Providers must contain at least one Provider")
	}
	for i, p := range c.Providers {
		if p == nil {
			return errors.New("updater: Config.Providers contains a nil entry at index " + itoa(i))
		}
	}
	return nil
}

// itoa is a tiny dependency-free int→string for error messages, scoped to
// this package so we don't drag in strconv just for validation strings.
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}
