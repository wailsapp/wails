package application

import (
	"strings"
	"testing"
)

func TestInjectInlineEventShim_OptedIn(t *testing.T) {
	got := maybeInjectInlineEventShim("<html><body>hi</body></html>", true)
	if !strings.Contains(got, "data-wails-inline-event-shim") {
		t.Errorf("shim should have been injected; got: %s", got)
	}
	if !strings.Contains(got, "window.wails.Events") {
		t.Errorf("shim should install window.wails.Events; got: %s", got)
	}
	if !strings.Contains(got, "wails:event:emit:") {
		t.Errorf("shim should route emits through wails:event:emit:; got: %s", got)
	}
}

func TestInjectInlineEventShim_NotOptedIn(t *testing.T) {
	src := "<html><body>hi</body></html>"
	got := maybeInjectInlineEventShim(src, false)
	if got != src {
		t.Errorf("shim must NOT be injected when AllowSimpleEventEmit is false; got: %s", got)
	}
}

func TestInjectInlineEventShim_EmptyHTML(t *testing.T) {
	got := maybeInjectInlineEventShim("", true)
	if got != "" {
		t.Errorf("empty HTML should pass through unchanged; got: %s", got)
	}
}

func TestInjectInlineEventShim_PreservesDoctype(t *testing.T) {
	src := "<!doctype html>\n<html><body>hi</body></html>"
	got := maybeInjectInlineEventShim(src, true)
	if !strings.HasPrefix(strings.ToLower(got), "<!doctype html>") {
		t.Errorf("DOCTYPE must remain on line 1; got: %s", got)
	}
}

func TestInjectInlineEventShim_NoDoubleInjection(t *testing.T) {
	src := "<html><body>hi</body></html>"
	once := maybeInjectInlineEventShim(src, true)
	twice := maybeInjectInlineEventShim(once, true)
	if strings.Count(twice, "data-wails-inline-event-shim") != 1 {
		t.Errorf("re-running the helper must not inject a second copy; got %d markers",
			strings.Count(twice, "data-wails-inline-event-shim"))
	}
}
