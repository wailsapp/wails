//go:build !production

package assetserver

import (
	"net"
	"net/url"
	"sync"
	"testing"
)

func TestGetStartURLDoesNotRetainDevPort(t *testing.T) {
	original := baseURL
	t.Cleanup(func() { baseURL = original })
	t.Setenv("FRONTEND_DEVSERVER_URL", "http://localhost:5173")
	wantDev := original
	wantDev.Host = net.JoinHostPort(original.Hostname(), "5173")
	if got, err := GetStartURL(""); err != nil || got != wantDev.String() {
		t.Errorf("dev start URL = %q, %v; want %q", got, err, wantDev.String())
	}
	t.Setenv("FRONTEND_DEVSERVER_URL", "")
	want := original.ResolveReference(&url.URL{Path: "/next"}).String()
	if got, err := GetStartURL("/next"); err != nil || got != want {
		t.Errorf("later start URL = %q, %v; want %q", got, err, want)
	}
	if baseURL != original {
		t.Error("GetStartURL modified the shared base URL")
	}
}

func TestGetStartURLConcurrent(t *testing.T) {
	original := baseURL
	t.Cleanup(func() { baseURL = original })
	t.Setenv("FRONTEND_DEVSERVER_URL", "http://localhost:5173")
	want := original
	want.Host = net.JoinHostPort(original.Hostname(), "5173")
	var workers sync.WaitGroup
	for range 16 {
		workers.Go(func() {
			for range 100 {
				if got, err := GetStartURL(""); err != nil || got != want.String() {
					t.Errorf("concurrent start URL = %q, %v; want %q", got, err, want.String())
				}
			}
		})
	}
	workers.Wait()
}
