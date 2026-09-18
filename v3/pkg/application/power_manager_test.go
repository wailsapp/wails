package application

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/wailsapp/wails/v3/pkg/events"
)

type fakePower struct {
	next   int
	active map[int]string
	ended  int
}

func newFakePower() *fakePower {
	return &fakePower{active: map[int]string{}}
}

func (f *fakePower) beginActivity(reason string, display bool) (any, error) {
	f.next++
	f.active[f.next] = reason
	return f.next, nil
}

func (f *fakePower) endActivity(token any) {
	id := token.(int)
	if _, ok := f.active[id]; !ok {
		panic("endActivity for unknown token")
	}
	delete(f.active, id)
	f.ended++
}

func (f *fakePower) state() PowerState {
	return PowerState{LowPowerMode: true, ThermalState: ThermalStateFair, OnBattery: true, BatteryLevel: 0.5, Charging: false}
}

func TestPreventSleepReferenceCounting(t *testing.T) {
	fake := newFakePower()
	pm := &PowerManager{impl: fake}

	releaseA, err := pm.PreventSleep("export", PreventSleepOptions{})
	if err != nil {
		t.Fatalf("PreventSleep: %v", err)
	}
	releaseB, err := pm.PreventSleep("upload", PreventSleepOptions{Display: true})
	if err != nil {
		t.Fatalf("PreventSleep: %v", err)
	}
	if got := pm.SleepHoldCount(); got != 2 {
		t.Fatalf("SleepHoldCount = %d, want 2", got)
	}
	if len(fake.active) != 2 {
		t.Fatalf("native activities = %d, want 2", len(fake.active))
	}

	releaseA()
	releaseA() // second release is a no-op
	if got := pm.SleepHoldCount(); got != 1 {
		t.Fatalf("SleepHoldCount after one release = %d, want 1", got)
	}
	if fake.ended != 1 {
		t.Fatalf("native endActivity calls = %d, want 1", fake.ended)
	}
	if fake.active[2] != "upload" {
		t.Fatalf("remaining hold = %q, want upload", fake.active[2])
	}

	releaseB()
	if got := pm.SleepHoldCount(); got != 0 {
		t.Fatalf("SleepHoldCount after both releases = %d, want 0", got)
	}
	if len(fake.active) != 0 || fake.ended != 2 {
		t.Fatalf("native state after release: active=%d ended=%d", len(fake.active), fake.ended)
	}
}

func TestPreventSleepDefaultsReason(t *testing.T) {
	fake := newFakePower()
	pm := &PowerManager{impl: fake}
	release, err := pm.PreventSleep("", PreventSleepOptions{})
	if err != nil {
		t.Fatalf("PreventSleep: %v", err)
	}
	defer release()
	if fake.active[1] == "" {
		t.Fatal("empty reason was passed to the platform")
	}
}

func TestPowerManagerWithoutImpl(t *testing.T) {
	pm := &PowerManager{}
	release, err := pm.PreventSleep("x", PreventSleepOptions{})
	if !errors.Is(err, ErrPreventSleepUnsupported) {
		t.Fatalf("err = %v, want ErrPreventSleepUnsupported", err)
	}
	release() // must be safe
	state := pm.State()
	if state.BatteryLevel != -1 || state.ThermalState != ThermalStateUnknown {
		t.Fatalf("State without impl = %+v", state)
	}
}

func TestPowerStateDelegatesAndSerialises(t *testing.T) {
	pm := &PowerManager{impl: newFakePower()}
	state := pm.State()
	if !state.LowPowerMode || state.ThermalState != ThermalStateFair || state.BatteryLevel != 0.5 {
		t.Fatalf("State = %+v", state)
	}
	raw, err := json.Marshal(state)
	if err != nil {
		t.Fatal(err)
	}
	if string(raw) != `{"lowPowerMode":true,"thermalState":"fair","onBattery":true,"batteryLevel":0.5,"charging":false}` {
		t.Fatalf("unexpected JSON %s", raw)
	}
}

func TestThermalStateStrings(t *testing.T) {
	want := map[ThermalState]string{
		ThermalStateUnknown:  "unknown",
		ThermalStateNominal:  "nominal",
		ThermalStateFair:     "fair",
		ThermalStateSerious:  "serious",
		ThermalStateCritical: "critical",
		ThermalState(99):     "unknown",
	}
	for state, text := range want {
		if got := state.String(); got != text {
			t.Errorf("ThermalState(%d).String() = %q, want %q", int(state), got, text)
		}
	}
}

func TestPowerEventsRegistered(t *testing.T) {
	cases := map[events.ApplicationEventType]string{
		events.Mac.ApplicationDidChangePowerState:   "mac:ApplicationDidChangePowerState",
		events.Mac.ApplicationDidChangeThermalState: "mac:ApplicationDidChangeThermalState",
	}
	for id, name := range cases {
		if id == 0 {
			t.Errorf("%s has no event id", name)
		}
		if got := events.JSEvent(uint(id)); got != name {
			t.Errorf("JSEvent(%d) = %q, want %q", id, got, name)
		}
		if !events.IsKnownEvent(name) {
			t.Errorf("%s is not a known event", name)
		}
	}
}
