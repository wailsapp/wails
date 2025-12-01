//go:build !ios

package application

func iosHapticsImpact(style string) {
	// no-op on non-iOS
}

type deviceInfo struct {
	Model         string `json:"model"`
	SystemName    string `json:"systemName"`
	SystemVersion string `json:"systemVersion"`
	IsSimulator   bool   `json:"isSimulator"`
}

func iosDeviceInfo() deviceInfo {
	return deviceInfo{}
}

// Live mutation stubs
func iosSetScrollEnabled(enabled bool)              {}
func iosSetBounceEnabled(enabled bool)              {}
func iosSetScrollIndicatorsEnabled(enabled bool)    {}
func iosSetBackForwardGesturesEnabled(enabled bool) {}
func iosSetLinkPreviewEnabled(enabled bool)         {}
func iosSetInspectableEnabled(enabled bool)         {}
func iosSetCustomUserAgent(ua string)               {}

// Native tabs stubs
func iosSetNativeTabsEnabled(enabled bool) {}
func iosNativeTabsIsEnabled() bool         { return false }
func iosSelectNativeTab(index int)         {}
