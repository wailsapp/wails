//go:build linux && !android && !server

package application

import (
	"sync"
	"testing"
	"time"
)

func TestNavigationPolicyPermitsOnlyTheAppsOrigins(t *testing.T) {
	var p navigationPolicy
	p.allow("http://localhost:5173/index.html")
	p.allow("https://example.org:443/")
	for raw, want := range map[string]bool{
		"wails://localhost/":            true,
		"wails://localhost/wails/x?y=1": true,
		"about:blank":                   true,
		"wails://":                      true,
		"wails:///":                     true,
		"https://example.org/x":         true,
		"https://EXAMPLE.org:443/x":     true,
		"blob:http://localhost:5173/id": true,
		"blob:https://example.com/id":   false,
		"blob:null/id":                  false,
		"http://localhost:5173/app":     true,
		"http://localhost:5174/app":     false,
		"https://example.com/":          false,
		"wails://evil/":                 false,
		"file:///etc/passwd":            false,
		"javascript:alert(1)":           false,
		"data:text/html,x":              false,
		"":                              false,
	} {
		if got := p.permits(raw); got != want {
			t.Errorf("permits(%q) = %t, want %t", raw, got, want)
		}
	}
}

func TestLinuxDecideNavigationSendsUserLinksToTheBrowser(t *testing.T) {
	var mu sync.Mutex
	var opened []string
	prev := openExternalURL
	openExternalURL = func(u string) error { mu.Lock(); opened = append(opened, u); mu.Unlock(); return nil }
	t.Cleanup(func() { openExternalURL = prev })
	calls := func() []string { mu.Lock(); defer mu.Unlock(); return append([]string(nil), opened...) }
	waitFor := func(n int) {
		t.Helper()
		for deadline := time.Now().Add(2 * time.Second); time.Now().Before(deadline); time.Sleep(5 * time.Millisecond) {
			if len(calls()) >= n {
				return
			}
		}
		t.Fatalf("%d browser calls, want %d", len(calls()), n)
	}

	if linuxDecideNavigation("https://example.com/a", true, false) {
		t.Fatal("a remote main-frame navigation was allowed")
	}
	waitFor(1)
	if got := calls(); got[0] != "https://example.com/a" {
		t.Fatalf("opened %q", got[0])
	}

	appNavigation.allow("https://trusted.example/")
	t.Cleanup(func() { appNavigation = navigationPolicy{} })
	if linuxDecideNavigation("https://trusted.example/page", true, true) {
		t.Fatal("a new window was allowed")
	}
	waitFor(2)
	if got := calls(); got[1] != "https://trusted.example/page" {
		t.Fatalf("opened %q", got[1])
	}

	if linuxDecideNavigation("wails://localhost/", true, true) {
		t.Fatal("a new window was allowed")
	}
	if !linuxDecideNavigation("wails://localhost/", false, false) {
		t.Fatal("the app's own page was refused")
	}
	if linuxDecideNavigation("https://example.com/b", false, false) {
		t.Fatal("a script-driven navigation was allowed")
	}
	time.Sleep(200 * time.Millisecond)
	if n := len(calls()); n != 2 {
		t.Fatalf("%d browser calls, want 2 (a script-driven or app-local request opened the browser)", n)
	}
}
