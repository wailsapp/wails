package application

import "testing"

func TestHapticKindValues(t *testing.T) {
	if HapticGeneric != 0 || HapticAlignment != 1 || HapticLevelChange != 2 {
		t.Fatalf("unexpected HapticKind values: %d %d %d", HapticGeneric, HapticAlignment, HapticLevelChange)
	}
	tests := map[HapticKind]string{
		HapticGeneric:     "generic",
		HapticAlignment:   "alignment",
		HapticLevelChange: "levelChange",
		HapticKind(-1):    "unknown",
		HapticKind(99):    "unknown",
	}
	for kind, want := range tests {
		if got := kind.String(); got != want {
			t.Errorf("HapticKind(%d).String() = %q, want %q", int(kind), got, want)
		}
		if got := kind.IsValid(); got != (want != "unknown") {
			t.Errorf("HapticKind(%d).IsValid() = %v", int(kind), got)
		}
	}
}

func TestHapticsManagerPerformUnknownKindDoesNotPanic(t *testing.T) {
	manager := newHapticsManager(nil)
	if manager == nil {
		t.Fatal("newHapticsManager returned nil")
	}
	// Perform must degrade unknown kinds to HapticGeneric rather than panic.
	// The platform call is a no-op or asynchronous so this is safe in tests.
	if !hapticsSupported() {
		manager.Perform(HapticKind(99))
	}
}
