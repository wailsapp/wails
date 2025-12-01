//go:build ios

package application

// --- Haptic Feedback ---

// HapticsImpact triggers native iOS haptic feedback using the Taptic Engine.
// Valid styles: "light", "medium", "heavy", "soft", "rigid"
func (m *IOSManager) HapticsImpact(style string) {
	iosHapticsImpact(style)
}

// --- Device Information ---

// DeviceInfo returns iOS device information including model and OS version.
func (m *IOSManager) DeviceInfo() deviceInfo {
	return iosDeviceInfo()
}

// --- WebView Control ---

// SetScrollEnabled enables or disables WebView scrolling at runtime.
func (m *IOSManager) SetScrollEnabled(enabled bool) {
	iosSetScrollEnabled(enabled)
}

// SetBounceEnabled enables or disables the iOS rubber-band bounce effect.
func (m *IOSManager) SetBounceEnabled(enabled bool) {
	iosSetBounceEnabled(enabled)
}

// SetScrollIndicatorsEnabled shows or hides scroll indicators.
func (m *IOSManager) SetScrollIndicatorsEnabled(enabled bool) {
	iosSetScrollIndicatorsEnabled(enabled)
}

// SetBackForwardGesturesEnabled enables or disables swipe navigation gestures.
func (m *IOSManager) SetBackForwardGesturesEnabled(enabled bool) {
	iosSetBackForwardGesturesEnabled(enabled)
}

// SetLinkPreviewEnabled enables or disables 3D Touch / long-press link previews.
func (m *IOSManager) SetLinkPreviewEnabled(enabled bool) {
	iosSetLinkPreviewEnabled(enabled)
}

// SetInspectableEnabled enables or disables Safari Web Inspector debugging.
func (m *IOSManager) SetInspectableEnabled(enabled bool) {
	iosSetInspectableEnabled(enabled)
}

// SetCustomUserAgent sets a custom User-Agent string for the WebView.
func (m *IOSManager) SetCustomUserAgent(ua string) {
	iosSetCustomUserAgent(ua)
}

// --- Native Tabs ---

// SetNativeTabsEnabled enables or disables native tab bar.
func (m *IOSManager) SetNativeTabsEnabled(enabled bool) {
	iosSetNativeTabsEnabled(enabled)
}

// NativeTabsIsEnabled returns whether native tabs are enabled.
func (m *IOSManager) NativeTabsIsEnabled() bool {
	return iosNativeTabsIsEnabled()
}

// SelectNativeTab selects the native tab at the given index.
func (m *IOSManager) SelectNativeTab(index int) {
	iosSelectNativeTab(index)
}
