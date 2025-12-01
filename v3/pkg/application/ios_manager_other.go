//go:build !ios

package application

// --- Haptic Feedback ---

// HapticsImpact is a no-op on non-iOS platforms.
func (m *IOSManager) HapticsImpact(style string) {}

// --- Device Information ---

// iosDeviceInfoStub is the device info type for non-iOS platforms
type iosDeviceInfoStub struct {
	Model         string `json:"model"`
	SystemName    string `json:"systemName"`
	SystemVersion string `json:"systemVersion"`
	IsSimulator   bool   `json:"isSimulator"`
}

// DeviceInfo returns empty device info on non-iOS platforms.
func (m *IOSManager) DeviceInfo() iosDeviceInfoStub {
	return iosDeviceInfoStub{}
}

// --- WebView Control ---

// SetScrollEnabled is a no-op on non-iOS platforms.
func (m *IOSManager) SetScrollEnabled(enabled bool) {}

// SetBounceEnabled is a no-op on non-iOS platforms.
func (m *IOSManager) SetBounceEnabled(enabled bool) {}

// SetScrollIndicatorsEnabled is a no-op on non-iOS platforms.
func (m *IOSManager) SetScrollIndicatorsEnabled(enabled bool) {}

// SetBackForwardGesturesEnabled is a no-op on non-iOS platforms.
func (m *IOSManager) SetBackForwardGesturesEnabled(enabled bool) {}

// SetLinkPreviewEnabled is a no-op on non-iOS platforms.
func (m *IOSManager) SetLinkPreviewEnabled(enabled bool) {}

// SetInspectableEnabled is a no-op on non-iOS platforms.
func (m *IOSManager) SetInspectableEnabled(enabled bool) {}

// SetCustomUserAgent is a no-op on non-iOS platforms.
func (m *IOSManager) SetCustomUserAgent(ua string) {}

// --- Native Tabs ---

// SetNativeTabsEnabled is a no-op on non-iOS platforms.
func (m *IOSManager) SetNativeTabsEnabled(enabled bool) {}

// NativeTabsIsEnabled returns false on non-iOS platforms.
func (m *IOSManager) NativeTabsIsEnabled() bool { return false }

// SelectNativeTab is a no-op on non-iOS platforms.
func (m *IOSManager) SelectNativeTab(index int) {}
