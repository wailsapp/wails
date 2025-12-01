//go:build ios

package application

// iosMobileManager implements mobileManagerImpl for iOS
type iosMobileManager struct{}

// newMobileManagerImpl creates the iOS implementation
func newMobileManagerImpl() mobileManagerImpl {
	return &iosMobileManager{}
}

func (m *iosMobileManager) vibrate(durationMs int) {
	// iOS ignores duration and uses a preset vibration
	// Use medium haptic as the default vibration
	iosHapticsImpact("medium")
}

func (m *iosMobileManager) hapticFeedback(style HapticStyle) {
	// Map HapticStyle to iOS haptic feedback style strings
	var styleStr string
	switch style {
	case HapticLight:
		styleStr = "light"
	case HapticMedium:
		styleStr = "medium"
	case HapticHeavy:
		styleStr = "heavy"
	case HapticSoft:
		styleStr = "soft"
	case HapticRigid:
		styleStr = "rigid"
	default:
		styleStr = "medium"
	}
	iosHapticsImpact(styleStr)
}

func (m *iosMobileManager) toast(message string) {
	// iOS doesn't have native Toast like Android
	// Could implement using UIKit banner notification in the future
	// For now, log the message
	// TODO: Implement iOS toast/banner notification
}

func (m *iosMobileManager) deviceInfo() DeviceInfo {
	// TODO: Implement iOS device info via native call
	return DeviceInfo{
		Platform:     "ios",
		Manufacturer: "Apple",
		IsVirtual:    false, // TODO: detect simulator
	}
}

func (m *iosMobileManager) screenInfo() MobileScreenInfo {
	// TODO: Implement iOS screen info via native call
	return MobileScreenInfo{}
}

func (m *iosMobileManager) isMobile() bool {
	return true
}

// WebView control methods - these call the existing iOS runtime functions

func (m *iosMobileManager) setScrollEnabled(enabled bool) {
	iosSetScrollEnabled(enabled)
}

func (m *iosMobileManager) setBounceEnabled(enabled bool) {
	iosSetBounceEnabled(enabled)
}

func (m *iosMobileManager) setScrollIndicatorsEnabled(enabled bool) {
	iosSetScrollIndicatorsEnabled(enabled)
}

func (m *iosMobileManager) setBackForwardGesturesEnabled(enabled bool) {
	iosSetBackForwardGesturesEnabled(enabled)
}

func (m *iosMobileManager) setLinkPreviewEnabled(enabled bool) {
	iosSetLinkPreviewEnabled(enabled)
}

func (m *iosMobileManager) setInspectableEnabled(enabled bool) {
	iosSetInspectableEnabled(enabled)
}

func (m *iosMobileManager) setCustomUserAgent(ua string) {
	iosSetCustomUserAgent(ua)
}
