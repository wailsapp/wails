//go:build windows && !server

package application

import "testing"

func TestPanelInitialRequestURLCanonicalization(t *testing.T) {
	for _, url := range []string{"https://example.com", "https://example.com/", "https://example.com/#section"} {
		if got := panelRequestURL(url); got != "https://example.com/" {
			t.Errorf("initial headers would miss %q: %q", url, got)
		}
	}
	if panelRequestURL("https://example.com/?page=2") == panelRequestURL("https://example.com/?page=1") {
		t.Fatal("headers could be applied to a different request")
	}
}
