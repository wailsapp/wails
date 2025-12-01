//go:build !android && !ios

package application

import "runtime"

// desktopMobileManager implements mobileManagerImpl for desktop platforms.
// All methods are no-ops that return sensible defaults.
type desktopMobileManager struct{}

// newMobileManagerImpl creates the desktop stub implementation
func newMobileManagerImpl() mobileManagerImpl {
	return &desktopMobileManager{}
}

func (m *desktopMobileManager) vibrate(durationMs int) {
	// No-op on desktop
}

func (m *desktopMobileManager) hapticFeedback(style HapticStyle) {
	// No-op on desktop
}

func (m *desktopMobileManager) toast(message string) {
	// No-op on desktop
	// Could potentially show a desktop notification in the future
}

func (m *desktopMobileManager) deviceInfo() DeviceInfo {
	return DeviceInfo{
		Platform:  "desktop",
		OSVersion: runtime.GOOS,
		IsVirtual: false,
	}
}

func (m *desktopMobileManager) screenInfo() MobileScreenInfo {
	// Desktop screen info is handled by ScreenManager
	return MobileScreenInfo{}
}

func (m *desktopMobileManager) isMobile() bool {
	return false
}

// WebView control methods - no-ops on desktop

func (m *desktopMobileManager) setScrollEnabled(enabled bool) {
	// Desktop WebViews handle scrolling differently
}

func (m *desktopMobileManager) setBounceEnabled(enabled bool) {
	// Bounce effect is mobile-only
}

func (m *desktopMobileManager) setScrollIndicatorsEnabled(enabled bool) {
	// Desktop WebViews handle scroll indicators via native scrollbars
}

func (m *desktopMobileManager) setBackForwardGesturesEnabled(enabled bool) {
	// Navigation gestures are mobile-only
}

func (m *desktopMobileManager) setLinkPreviewEnabled(enabled bool) {
	// Link preview is mobile-only
}

func (m *desktopMobileManager) setInspectableEnabled(enabled bool) {
	// Desktop DevTools are typically always available
}

func (m *desktopMobileManager) setCustomUserAgent(ua string) {
	// Could be implemented for desktop WebViews in the future
}
