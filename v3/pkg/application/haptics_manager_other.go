//go:build (!darwin || ios || server) && !android

package application

import "runtime"

func hapticsPerform(kind HapticKind) {
	if runtime.GOOS != "ios" {
		return
	}
	style := "medium"
	switch kind {
	case HapticAlignment:
		style = "light"
	case HapticLevelChange:
		style = "heavy"
	}
	iosHapticsImpact(style)
}

func hapticsSupported() bool {
	return runtime.GOOS == "ios"
}
