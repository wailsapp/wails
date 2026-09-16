package application

import (
	"errors"
	"sync"
)

// PreventSleepOptions tunes PowerManager.PreventSleep.
type PreventSleepOptions struct {
	// Display also keeps the display awake. Without it only idle system
	// sleep is prevented and the screen may still dim and lock.
	Display bool
}

// ThermalState is the system's thermal pressure level.
type ThermalState int

const (
	// ThermalStateUnknown means the platform does not report thermal state.
	ThermalStateUnknown ThermalState = iota
	// ThermalStateNominal is normal operation.
	ThermalStateNominal
	// ThermalStateFair means the system is warming; consider trimming
	// background work.
	ThermalStateFair
	// ThermalStateSerious means the system is throttling; reduce CPU, GPU
	// and I/O usage.
	ThermalStateSerious
	// ThermalStateCritical means the system is about to be throttled hard;
	// do the minimum work possible.
	ThermalStateCritical
)

// String returns the lower-case name of the thermal state.
func (t ThermalState) String() string {
	switch t {
	case ThermalStateNominal:
		return "nominal"
	case ThermalStateFair:
		return "fair"
	case ThermalStateSerious:
		return "serious"
	case ThermalStateCritical:
		return "critical"
	}
	return "unknown"
}

// MarshalText makes ThermalState serialise as its name.
func (t ThermalState) MarshalText() ([]byte, error) {
	return []byte(t.String()), nil
}

// PowerState is a snapshot of the machine's power situation.
type PowerState struct {
	// LowPowerMode reports whether the user enabled Low Power Mode
	// (macOS 12+).
	LowPowerMode bool `json:"lowPowerMode"`
	// ThermalState is the current thermal pressure.
	ThermalState ThermalState `json:"thermalState"`
	// OnBattery is true when the machine runs from its battery.
	OnBattery bool `json:"onBattery"`
	// BatteryLevel is the charge between 0 and 1, or -1 when there is no
	// battery or the level is unknown.
	BatteryLevel float64 `json:"batteryLevel"`
	// Charging is true while the battery is being charged.
	Charging bool `json:"charging"`
}

// ErrPreventSleepUnsupported is returned by PreventSleep where the platform
// has no sleep-assertion API; the returned release func is still safe to
// call.
var ErrPreventSleepUnsupported = errors.New("power: preventing sleep is not supported on this platform")

// powerImpl is the platform side of PowerManager.
type powerImpl interface {
	// beginActivity starts a sleep assertion and returns an opaque token
	// for endActivity.
	beginActivity(reason string, display bool) (token any, err error)
	endActivity(token any)
	state() PowerState
}

// PowerManager prevents sleep and reports power and thermal state.
//
// Changes are announced through events.Mac.ApplicationDidChangePowerState
// (Low Power Mode toggled) and events.Mac.ApplicationDidChangeThermalState.
//
// Platform support is documented on each method. Methods without support on
// the current platform are safe no-ops or return a documented error.
type PowerManager struct {
	app   *App
	impl  powerImpl
	lock  sync.Mutex
	holds int
}

func newPowerManager(app *App) *PowerManager {
	return &PowerManager{app: app, impl: newPowerImpl(app)}
}

// PreventSleep keeps the system (and, with opts.Display, the display) awake
// until the returned release func is called. Each call is an independent
// hold: holds from different parts of the app coexist and sleep is allowed
// again only when every hold has been released. Calling release more than
// once is harmless. reason is shown to the user in Activity Monitor and
// pmset, so make it descriptive ("Exporting video").
//
// macOS: NSProcessInfo beginActivityWithOptions: with
// NSActivityIdleSystemSleepDisabled, plus NSActivityIdleDisplaySleepDisabled
// when opts.Display is set. Visible in `pmset -g assertions`.
//
// Other platforms: ErrPreventSleepUnsupported with a no-op release.
func (pm *PowerManager) PreventSleep(reason string, opts PreventSleepOptions) (release func(), err error) {
	noop := func() {}
	if pm.impl == nil {
		return noop, ErrPreventSleepUnsupported
	}
	if reason == "" {
		reason = "Wails application activity"
	}
	token, err := pm.impl.beginActivity(reason, opts.Display)
	if err != nil {
		return noop, err
	}
	pm.lock.Lock()
	pm.holds++
	pm.lock.Unlock()

	var once sync.Once
	return func() {
		once.Do(func() {
			pm.impl.endActivity(token)
			pm.lock.Lock()
			pm.holds--
			pm.lock.Unlock()
		})
	}, nil
}

// SleepHoldCount returns the number of PreventSleep holds not yet released.
func (pm *PowerManager) SleepHoldCount() int {
	pm.lock.Lock()
	defer pm.lock.Unlock()
	return pm.holds
}

// State returns the current power and thermal state.
//
// macOS: NSProcessInfo for Low Power Mode and thermal state, IOKit power
// sources for battery information. Desktops without a battery report
// OnBattery false and BatteryLevel -1.
//
// Other platforms: zero values with BatteryLevel -1 and
// ThermalStateUnknown.
func (pm *PowerManager) State() PowerState {
	if pm.impl == nil {
		return PowerState{BatteryLevel: -1, ThermalState: ThermalStateUnknown}
	}
	return pm.impl.state()
}
