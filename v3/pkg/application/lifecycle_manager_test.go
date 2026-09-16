package application

import "testing"

type fakeLifecycle struct {
	plistDefault bool
	disabled     int // net disableTermination calls
	suddenCalls  []bool
}

func (f *fakeLifecycle) suddenTerminationDefault() bool { return f.plistDefault }
func (f *fakeLifecycle) disableTermination(reason string) {
	f.disabled++
}
func (f *fakeLifecycle) enableTermination(reason string) {
	f.disabled--
	if f.disabled < 0 {
		panic("enableTermination called more often than disableTermination")
	}
}
func (f *fakeLifecycle) setSuddenTerminationEnabled(enabled bool) {
	f.suddenCalls = append(f.suddenCalls, enabled)
}

func TestHoldTerminationReferenceCounting(t *testing.T) {
	fake := &fakeLifecycle{}
	lm := &LifecycleManager{impl: fake}

	releaseA := lm.HoldTermination("saving")
	releaseB := lm.HoldTermination("uploading")
	if got := lm.TerminationHoldCount(); got != 2 {
		t.Fatalf("TerminationHoldCount = %d, want 2", got)
	}
	if fake.disabled != 2 {
		t.Fatalf("native disable count = %d, want 2", fake.disabled)
	}

	releaseA()
	releaseA() // idempotent
	if got := lm.TerminationHoldCount(); got != 1 || fake.disabled != 1 {
		t.Fatalf("after one release: holds=%d native=%d, want 1/1", got, fake.disabled)
	}

	releaseB()
	if got := lm.TerminationHoldCount(); got != 0 || fake.disabled != 0 {
		t.Fatalf("after both releases: holds=%d native=%d, want 0/0", got, fake.disabled)
	}
}

func TestSetSuddenTerminationEnabledIsBalanced(t *testing.T) {
	fake := &fakeLifecycle{plistDefault: false}
	lm := &LifecycleManager{impl: fake}

	if lm.SuddenTerminationEnabled() {
		t.Fatal("default should follow the plist default (false)")
	}
	lm.SetSuddenTerminationEnabled(false) // already the default: no native call
	if len(fake.suddenCalls) != 0 {
		t.Fatalf("native calls after redundant disable = %v, want none", fake.suddenCalls)
	}
	lm.SetSuddenTerminationEnabled(true)
	lm.SetSuddenTerminationEnabled(true) // idempotent
	lm.SetSuddenTerminationEnabled(false)
	if want := []bool{true, false}; len(fake.suddenCalls) != 2 || fake.suddenCalls[0] != want[0] || fake.suddenCalls[1] != want[1] {
		t.Fatalf("native calls = %v, want %v", fake.suddenCalls, want)
	}
	if lm.SuddenTerminationEnabled() {
		t.Fatal("SuddenTerminationEnabled should be false after the final disable")
	}
}

func TestSetSuddenTerminationEnabledHonoursPlistDefault(t *testing.T) {
	fake := &fakeLifecycle{plistDefault: true}
	lm := &LifecycleManager{impl: fake}
	if !lm.SuddenTerminationEnabled() {
		t.Fatal("default should follow the plist default (true)")
	}
	lm.SetSuddenTerminationEnabled(true) // no-op: already enabled by the plist
	if len(fake.suddenCalls) != 0 {
		t.Fatalf("native calls = %v, want none", fake.suddenCalls)
	}
	lm.SetSuddenTerminationEnabled(false)
	if len(fake.suddenCalls) != 1 || fake.suddenCalls[0] {
		t.Fatalf("native calls = %v, want [false]", fake.suddenCalls)
	}
}

func TestLifecycleManagerWithoutImpl(t *testing.T) {
	lm := &LifecycleManager{}
	release := lm.HoldTermination("x")
	release()
	release()
	lm.SetSuddenTerminationEnabled(true)
	if lm.SuddenTerminationEnabled() {
		t.Fatal("SuddenTerminationEnabled without impl should be false")
	}
	if got := lm.TerminationHoldCount(); got != 0 {
		t.Fatalf("TerminationHoldCount without impl = %d, want 0", got)
	}
}
