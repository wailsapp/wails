//go:build android

package application

func hapticsPerform(kind HapticKind) {
	durationMs := 20
	switch kind {
	case HapticAlignment:
		durationMs = 10
	case HapticLevelChange:
		durationMs = 40
	}
	androidHapticsVibrate(durationMs)
}

func hapticsSupported() bool {
	return true
}
