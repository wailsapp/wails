package runtime

import (
	"context"
)

// BrowserOpenURL mirrors the v2 runtime.BrowserOpenURL function. Any error is
// logged as v2 returned nothing.
// v3 equivalent: app.Browser.OpenURL.
func BrowserOpenURL(_ context.Context, url string) {
	a := app()
	if a == nil {
		return
	}
	if err := a.Browser.OpenURL(url); err != nil {
		logger().Warn("v2compat: BrowserOpenURL failed", "url", url, "error", err)
	}
}
