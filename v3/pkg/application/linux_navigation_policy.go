//go:build linux && !android && !server

package application

import (
	"net/url"
	"strings"
	"sync"

	"github.com/wailsapp/wails/v3/internal/assetserver"
	"github.com/wailsapp/wails/v3/internal/browser"
)

// openExternalURL hands a link to the system browser; tests replace it.
var openExternalURL = browser.OpenURL

// navigationPolicy decides which documents a Linux webview may load. The app's
// own pages talk to Go through the wails: message channel and /wails/runtime, so
// a remote document in the window must never be able to load there.
type navigationPolicy struct {
	mu      sync.Mutex
	origins map[string]struct{}
}

func originOf(raw string) (string, bool) {
	u, err := url.Parse(raw)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return "", false
	}
	return u.Scheme + "://" + strings.ToLower(u.Host), true
}

// allow trusts the origin of a URL the application itself loads into the
// window (its start URL, the dev server, a later SetURL).
func (p *navigationPolicy) allow(raw string) {
	origin, ok := originOf(raw)
	if !ok {
		return
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.origins == nil {
		p.origins = map[string]struct{}{}
	}
	p.origins[origin] = struct{}{}
}

// permits reports whether the window may navigate to raw: about:blank, the
// asset server's origin, the dev server's origin, or an origin the application
// loaded itself.
func (p *navigationPolicy) permits(raw string) bool {
	if raw == "about:blank" {
		return true
	}
	origin, ok := originOf(raw)
	if !ok {
		return false
	}
	if origin == "wails://localhost" {
		return true
	}
	if dev, ok := originOf(assetserver.GetDevServerURL()); ok && dev == origin {
		return true
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	_, trusted := p.origins[origin]
	return trusted
}

// appNavigation is shared by every window: an origin the application loads into one window is its own page.
var appNavigation navigationPolicy

// linuxDecideNavigation is what each backend's decide-policy handler asks. newWindow is true for a request to open another window (window.open, target=_blank).
func linuxDecideNavigation(raw string, userGesture, newWindow bool) bool {
	allowed := appNavigation.decideNavigation(raw, userGesture)
	// No handler creates windows, so a new-window request is never allowed to proceed.
	return allowed && !newWindow
}

// decideNavigation is the policy for one navigation or new-window request. It
// reports whether WebKit may proceed. A refused request that is an http(s) link
// opens in the system browser when a user gesture asked for it.
func (p *navigationPolicy) decideNavigation(raw string, userGesture bool) bool {
	if p.permits(raw) {
		return true
	}
	if userGesture {
		if u, err := ValidateAndSanitizeURL(raw); err == nil && (strings.HasPrefix(u, "http://") || strings.HasPrefix(u, "https://")) {
			go func() { _ = openExternalURL(u) }()
		}
	}
	return false
}
