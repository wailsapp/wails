package application

// HapticKind selects the haptic feedback pattern performed by
// HapticsManager.Perform.
type HapticKind int

const (
	// HapticGeneric is a general feedback pattern for events that do not fit
	// a more specific kind (NSHapticFeedbackPatternGeneric).
	HapticGeneric HapticKind = iota
	// HapticAlignment signals that an item snapped into alignment with
	// another item, for example a guide during a drag
	// (NSHapticFeedbackPatternAlignment).
	HapticAlignment
	// HapticLevelChange signals a change in a discrete level, for example a
	// slider detent or a Force click stage (NSHapticFeedbackPatternLevelChange).
	HapticLevelChange
)

// String returns the name of the haptic kind.
func (k HapticKind) String() string {
	switch k {
	case HapticGeneric:
		return "generic"
	case HapticAlignment:
		return "alignment"
	case HapticLevelChange:
		return "levelChange"
	default:
		return "unknown"
	}
}

// IsValid reports whether k is one of the defined haptic kinds.
func (k HapticKind) IsValid() bool {
	return k >= HapticGeneric && k <= HapticLevelChange
}

// HapticsManager performs system haptic feedback where the platform supports it.
//
// Platform support is documented on each method. Methods without support on
// the current platform are safe no-ops or return a documented error.
type HapticsManager struct {
	app *App
}

func newHapticsManager(app *App) *HapticsManager {
	return &HapticsManager{app: app}
}

// Perform plays the haptic pattern for kind immediately.
//
// macOS: uses NSHapticFeedbackManager, so feedback is only felt on a Force
// Touch trackpad or Magic Trackpad and only while the app is active.
// iOS: routed to UIImpactFeedbackGenerator (generic = medium, alignment =
// light, level change = heavy).
// Android: routed to the vibrator (generic = 20 ms, alignment = 10 ms,
// level change = 40 ms); requires the VIBRATE permission.
// Windows and Linux: no-op.
//
// Unknown kinds are treated as HapticGeneric.
func (h *HapticsManager) Perform(kind HapticKind) {
	if !kind.IsValid() {
		kind = HapticGeneric
	}
	hapticsPerform(kind)
}

// IsSupported reports whether Perform can produce feedback on this platform.
// It does not check for the presence of a haptic-capable device.
func (h *HapticsManager) IsSupported() bool {
	return hapticsSupported()
}
