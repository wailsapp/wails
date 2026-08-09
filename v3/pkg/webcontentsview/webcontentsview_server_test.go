//go:build server && !android && !ios

package webcontentsview

import "testing"

func TestServerWebContentsViewIsExplicitlyNonNative(t *testing.T) {
	view := NewWebContentsView(WebContentsViewOptions{URL: "https://example.com"})
	if _, ok := view.impl.(*serverWebContentsView); !ok {
		t.Fatalf("implementation = %T, want *serverWebContentsView", view.impl)
	}
	if got := view.GetURL(); got != "https://example.com" {
		t.Fatalf("GetURL() = %q, want requested URL", got)
	}
}
