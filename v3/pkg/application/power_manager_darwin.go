//go:build darwin && !ios && !server

package application

/*
#cgo CFLAGS: -mmacosx-version-min=10.13 -x objective-c
#cgo LDFLAGS: -framework Cocoa -framework IOKit -mmacosx-version-min=10.13

#include <stdlib.h>
#include "power_manager_darwin.h"
*/
import "C"

import (
	"errors"
	"unsafe"

	"github.com/wailsapp/wails/v3/pkg/events"
)

type macosPower struct {
	app *App
}

func newPowerImpl(app *App) powerImpl {
	return &macosPower{app: app}
}

func (p *macosPower) beginActivity(reason string, display bool) (any, error) {
	cReason := C.CString(reason)
	defer C.free(unsafe.Pointer(cReason))
	token := C.wailsPowerBeginActivity(cReason, C.bool(display))
	if token == nil {
		return nil, errors.New("power: NSProcessInfo refused the activity")
	}
	return token, nil
}

func (p *macosPower) endActivity(token any) {
	pointer, ok := token.(unsafe.Pointer)
	if !ok || pointer == nil {
		return
	}
	C.wailsPowerEndActivity(pointer)
}

func (p *macosPower) state() PowerState {
	native := C.wailsPowerState()
	return PowerState{
		LowPowerMode: bool(native.lowPowerMode),
		ThermalState: ThermalState(native.thermalState),
		OnBattery:    bool(native.onBattery),
		BatteryLevel: float64(native.batteryLevel),
		Charging:     bool(native.charging),
	}
}

//export powerNotification
func powerNotification(kind C.int) {
	switch kind {
	case 0:
		applicationEvents <- newApplicationEvent(events.Mac.ApplicationDidChangePowerState)
	case 1:
		applicationEvents <- newApplicationEvent(events.Mac.ApplicationDidChangeThermalState)
	}
}
