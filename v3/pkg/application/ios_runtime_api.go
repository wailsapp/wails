//go:build ios

package application

// Exported API for use by applications to mutate iOS WKWebView at runtime.
// These call into the internal platform-specific implementations.

func IOSSetScrollEnabled(enabled bool)              { iosSetScrollEnabled(enabled) }
func IOSSetBounceEnabled(enabled bool)              { iosSetBounceEnabled(enabled) }
func IOSSetScrollIndicatorsEnabled(enabled bool)    { iosSetScrollIndicatorsEnabled(enabled) }
func IOSSetBackForwardGesturesEnabled(enabled bool) { iosSetBackForwardGesturesEnabled(enabled) }
func IOSSetLinkPreviewEnabled(enabled bool)         { iosSetLinkPreviewEnabled(enabled) }
func IOSSetInspectableEnabled(enabled bool)         { iosSetInspectableEnabled(enabled) }
func IOSSetCustomUserAgent(ua string)               { iosSetCustomUserAgent(ua) }
