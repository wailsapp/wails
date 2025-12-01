//go:build android

package application

import "encoding/json"

// androidMobileManager implements mobileManagerImpl for Android
type androidMobileManager struct{}

// newMobileManagerImpl creates the Android implementation
func newMobileManagerImpl() mobileManagerImpl {
	return &androidMobileManager{}
}

func (m *androidMobileManager) vibrate(durationMs int) {
	AndroidVibrate(durationMs)
}

func (m *androidMobileManager) hapticFeedback(style HapticStyle) {
	// Map HapticStyle to vibration duration
	// Android doesn't have native haptic styles like iOS, so we use vibration patterns
	var durationMs int
	switch style {
	case HapticLight:
		durationMs = 10
	case HapticMedium:
		durationMs = 25
	case HapticHeavy:
		durationMs = 50
	case HapticSoft:
		durationMs = 15
	case HapticRigid:
		durationMs = 20
	default:
		durationMs = 25
	}
	AndroidVibrate(durationMs)
}

func (m *androidMobileManager) toast(message string) {
	AndroidShowToast(message)
}

func (m *androidMobileManager) deviceInfo() DeviceInfo {
	infoJSON := AndroidGetDeviceInfo()
	if infoJSON == "" {
		return DeviceInfo{Platform: "android"}
	}

	// Parse the JSON from the native call
	var rawInfo struct {
		Platform     string `json:"platform"`
		Manufacturer string `json:"manufacturer"`
		Model        string `json:"model"`
		Brand        string `json:"brand"`
		SDKVersion   int    `json:"sdkVersion"`
		Release      string `json:"release"`
	}
	if err := json.Unmarshal([]byte(infoJSON), &rawInfo); err != nil {
		return DeviceInfo{Platform: "android"}
	}

	return DeviceInfo{
		Platform:     rawInfo.Platform,
		Model:        rawInfo.Model,
		Manufacturer: rawInfo.Manufacturer,
		Brand:        rawInfo.Brand,
		OSVersion:    rawInfo.Release,
		SDKVersion:   rawInfo.SDKVersion,
		IsVirtual:    false, // TODO: detect emulator
	}
}

func (m *androidMobileManager) screenInfo() MobileScreenInfo {
	infoJSON := AndroidGetScreenInfo()
	if infoJSON == "" {
		return MobileScreenInfo{}
	}

	var rawInfo struct {
		WidthPixels   int     `json:"widthPixels"`
		HeightPixels  int     `json:"heightPixels"`
		Density       float64 `json:"density"`
		DensityDPI    int     `json:"densityDpi"`
		ScaledDensity float64 `json:"scaledDensity"`
	}
	if err := json.Unmarshal([]byte(infoJSON), &rawInfo); err != nil {
		return MobileScreenInfo{}
	}

	return MobileScreenInfo{
		WidthPixels:   rawInfo.WidthPixels,
		HeightPixels:  rawInfo.HeightPixels,
		Density:       rawInfo.Density,
		DensityDPI:    rawInfo.DensityDPI,
		ScaledDensity: rawInfo.ScaledDensity,
	}
}

func (m *androidMobileManager) isMobile() bool {
	return true
}

// WebView control methods - Android currently doesn't expose these via JNI
// TODO: Implement these in WailsBridge.java

func (m *androidMobileManager) setScrollEnabled(enabled bool) {
	// TODO: Implement via JNI
	androidLogf("debug", "setScrollEnabled(%v) - not yet implemented on Android", enabled)
}

func (m *androidMobileManager) setBounceEnabled(enabled bool) {
	// Android calls this "overscroll"
	// TODO: Implement via JNI
	androidLogf("debug", "setBounceEnabled(%v) - not yet implemented on Android", enabled)
}

func (m *androidMobileManager) setScrollIndicatorsEnabled(enabled bool) {
	// TODO: Implement via JNI
	androidLogf("debug", "setScrollIndicatorsEnabled(%v) - not yet implemented on Android", enabled)
}

func (m *androidMobileManager) setBackForwardGesturesEnabled(enabled bool) {
	// TODO: Implement via JNI
	androidLogf("debug", "setBackForwardGesturesEnabled(%v) - not yet implemented on Android", enabled)
}

func (m *androidMobileManager) setLinkPreviewEnabled(enabled bool) {
	// Link preview is iOS-specific (3D Touch/Haptic Touch)
	// No equivalent on Android
}

func (m *androidMobileManager) setInspectableEnabled(enabled bool) {
	// Android WebView debugging is controlled via WebView.setWebContentsDebuggingEnabled()
	// This is typically set at app startup, not runtime
	androidLogf("debug", "setInspectableEnabled(%v) - set at build time on Android", enabled)
}

func (m *androidMobileManager) setCustomUserAgent(ua string) {
	// TODO: Implement via JNI - WebView.getSettings().setUserAgentString()
	androidLogf("debug", "setCustomUserAgent(%s) - not yet implemented on Android", ua)
}
