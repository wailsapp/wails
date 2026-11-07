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
	for raw, want := range map[string]bool{
		"wails://localhost/":            true,
		"wails://localhost/wails/x?y=1": true,
		"about:blank":                   true,
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
	opened := make(chan string, 4)
	var once sync.Once
	prev := openExternalURL
	openExternalURL = func(u string) error { once.Do(func() { opened <- u }); return nil }
	t.Cleanup(func() { openExternalURL = prev })

	if linuxDecideNavigation("https://example.com/a", true, false) {
		t.Fatal("a remote main-frame navigation was allowed")
	}
	select {
	case u := <-opened:
		if u != "https://example.com/a" {
			t.Fatalf("opened %q", u)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("a user-gesture link did not reach the browser")
	}

	if linuxDecideNavigation("wails://localhost/", true, true) {
		t.Fatal("a new window was allowed")
	}
	if !linuxDecideNavigation("wails://localhost/", false, false) {
		t.Fatal("the app's own page was refused")
	}
	if linuxDecideNavigation("https://example.com/b", false, false) || len(opened) != 0 {
		t.Fatal("a script-driven navigation was allowed or opened the browser")
	}
}
