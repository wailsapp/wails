//go:build ios

package application

// Exported API for use by applications to mutate iOS WKWebView at runtime.
// These call into the internal platform-specific implementations.

func (iosManager) SetScrollEnabled(enabled bool)           { iosSetScrollEnabled(enabled) }
func (iosManager) SetBounceEnabled(enabled bool)           { iosSetBounceEnabled(enabled) }
func (iosManager) SetScrollIndicatorsEnabled(enabled bool) { iosSetScrollIndicatorsEnabled(enabled) }
func (iosManager) SetBackForwardGesturesEnabled(enabled bool) {
	iosSetBackForwardGesturesEnabled(enabled)
}
func (iosManager) SetLinkPreviewEnabled(enabled bool) { iosSetLinkPreviewEnabled(enabled) }
func (iosManager) SetInspectableEnabled(enabled bool) { iosSetInspectableEnabled(enabled) }
func (iosManager) SetCustomUserAgent(ua string)       { iosSetCustomUserAgent(ua) }
