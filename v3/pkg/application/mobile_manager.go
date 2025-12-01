package application

// HapticStyle represents the intensity/type of haptic feedback
type HapticStyle int

const (
	// HapticLight provides subtle feedback for minor UI interactions
	HapticLight HapticStyle = iota
	// HapticMedium provides moderate feedback (default)
	HapticMedium
	// HapticHeavy provides strong feedback for significant actions
	HapticHeavy
	// HapticSoft provides a soft, gentle feedback
	HapticSoft
	// HapticRigid provides sharp, precise feedback
	HapticRigid
)

// DeviceInfo contains information about the mobile device
type DeviceInfo struct {
	// Platform is the operating system: "android", "ios", or "desktop"
	Platform string `json:"platform"`
	// Model is the device model name (e.g., "Pixel 7", "iPhone 15")
	Model string `json:"model"`
	// Manufacturer is the device manufacturer (e.g., "Google", "Apple")
	Manufacturer string `json:"manufacturer"`
	// Brand is the consumer-visible brand (e.g., "google", "Apple")
	Brand string `json:"brand"`
	// OSVersion is the OS version string (e.g., "14", "17.0")
	OSVersion string `json:"osVersion"`
	// SDKVersion is the SDK/API level (Android only)
	SDKVersion int `json:"sdkVersion,omitempty"`
	// IsVirtual indicates if running on an emulator/simulator
	IsVirtual bool `json:"isVirtual"`
}

// MobileScreenInfo contains display metrics for the device screen
type MobileScreenInfo struct {
	// WidthPixels is the screen width in pixels
	WidthPixels int `json:"widthPixels"`
	// HeightPixels is the screen height in pixels
	HeightPixels int `json:"heightPixels"`
	// Density is the display density scale factor
	Density float64 `json:"density"`
	// DensityDPI is the screen density in dots per inch
	DensityDPI int `json:"densityDpi"`
	// ScaledDensity is the font scaling factor
	ScaledDensity float64 `json:"scaledDensity,omitempty"`
}

// WebViewOptions contains options for controlling the mobile WebView
type WebViewOptions struct {
	// ScrollEnabled controls whether scrolling is enabled
	ScrollEnabled bool
	// BounceEnabled controls the iOS rubber-band bounce effect
	BounceEnabled bool
	// ScrollIndicatorsEnabled controls visibility of scroll indicators
	ScrollIndicatorsEnabled bool
	// BackForwardGesturesEnabled controls swipe navigation gestures
	BackForwardGesturesEnabled bool
	// LinkPreviewEnabled controls 3D Touch link previews (iOS)
	LinkPreviewEnabled bool
	// InspectableEnabled allows Safari/Chrome DevTools debugging
	InspectableEnabled bool
}

// MobileManager provides cross-platform mobile device functionality.
// On desktop platforms, most methods are no-ops that return sensible defaults.
type MobileManager struct {
	app  *App
	impl mobileManagerImpl
}

// mobileManagerImpl is the platform-specific implementation interface
type mobileManagerImpl interface {
	// Haptic feedback
	vibrate(durationMs int)
	hapticFeedback(style HapticStyle)

	// Notifications
	toast(message string)

	// Device info
	deviceInfo() DeviceInfo
	screenInfo() MobileScreenInfo
	isMobile() bool

	// WebView control
	setScrollEnabled(enabled bool)
	setBounceEnabled(enabled bool)
	setScrollIndicatorsEnabled(enabled bool)
	setBackForwardGesturesEnabled(enabled bool)
	setLinkPreviewEnabled(enabled bool)
	setInspectableEnabled(enabled bool)
	setCustomUserAgent(ua string)
}

// newMobileManager creates a new MobileManager instance
func newMobileManager(app *App) *MobileManager {
	return &MobileManager{
		app:  app,
		impl: newMobileManagerImpl(),
	}
}

// --- Haptic Feedback ---

// Vibrate triggers device vibration for the specified duration in milliseconds.
// On iOS, the duration parameter is ignored and a preset vibration is used.
// On desktop platforms, this is a no-op.
func (m *MobileManager) Vibrate(durationMs int) {
	m.impl.vibrate(durationMs)
}

// HapticFeedback triggers native haptic feedback with the specified style.
// On iOS, this uses the Taptic Engine with UIImpactFeedbackGenerator.
// On Android, this triggers a short vibration pattern mapped to the style.
// On desktop platforms, this is a no-op.
func (m *MobileManager) HapticFeedback(style HapticStyle) {
	m.impl.hapticFeedback(style)
}

// --- Notifications ---

// Toast displays a brief, non-blocking notification message.
// On Android, this shows a native Toast.
// On iOS, this could show a notification banner (implementation pending).
// On desktop platforms, this is a no-op.
func (m *MobileManager) Toast(message string) {
	m.impl.toast(message)
}

// --- Device Information ---

// DeviceInfo returns information about the device.
// On desktop platforms, returns a DeviceInfo with Platform set to "desktop".
func (m *MobileManager) DeviceInfo() DeviceInfo {
	return m.impl.deviceInfo()
}

// ScreenInfo returns display metrics for the device screen.
// On desktop platforms, returns default values.
func (m *MobileManager) ScreenInfo() MobileScreenInfo {
	return m.impl.screenInfo()
}

// IsMobile returns true if running on a mobile platform (Android or iOS).
func (m *MobileManager) IsMobile() bool {
	return m.impl.isMobile()
}

// --- WebView Control ---

// SetScrollEnabled enables or disables WebView scrolling.
// Disabling is useful for fixed-layout UIs.
func (m *MobileManager) SetScrollEnabled(enabled bool) {
	m.impl.setScrollEnabled(enabled)
}

// SetBounceEnabled enables or disables the iOS rubber-band bounce effect.
// On Android, this controls overscroll effects.
func (m *MobileManager) SetBounceEnabled(enabled bool) {
	m.impl.setBounceEnabled(enabled)
}

// SetScrollIndicatorsEnabled shows or hides scroll indicators.
func (m *MobileManager) SetScrollIndicatorsEnabled(enabled bool) {
	m.impl.setScrollIndicatorsEnabled(enabled)
}

// SetBackForwardGesturesEnabled enables or disables swipe navigation gestures.
// When enabled, users can swipe to go back/forward in navigation history.
func (m *MobileManager) SetBackForwardGesturesEnabled(enabled bool) {
	m.impl.setBackForwardGesturesEnabled(enabled)
}

// SetLinkPreviewEnabled enables or disables 3D Touch / long-press link previews (iOS).
// On Android, this has no effect.
func (m *MobileManager) SetLinkPreviewEnabled(enabled bool) {
	m.impl.setLinkPreviewEnabled(enabled)
}

// SetInspectableEnabled enables Safari Web Inspector (iOS) or Chrome DevTools (Android).
// Should typically only be enabled in debug builds.
func (m *MobileManager) SetInspectableEnabled(enabled bool) {
	m.impl.setInspectableEnabled(enabled)
}

// SetCustomUserAgent sets a custom User-Agent string for the WebView.
func (m *MobileManager) SetCustomUserAgent(ua string) {
	m.impl.setCustomUserAgent(ua)
}
