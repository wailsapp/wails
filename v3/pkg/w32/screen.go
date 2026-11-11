//go:build windows

package w32

import (
	"fmt"
	"strings"
	"sync"
	"syscall"
	"unsafe"
)

type Screen struct {
	MONITORINFOEX
	HMonitor    uintptr
	Name        string
	IsPrimary   bool
	IsCurrent   bool
	ScaleFactor float32
	Rotation    float32
}

type DISPLAY_DEVICE struct {
	cb           uint32
	DeviceName   [32]uint16
	DeviceString [128]uint16
	StateFlags   uint32
	DeviceID     [128]uint16
	DeviceKey    [128]uint16
}

func getMonitorName(deviceName string) (string, error) {
	var device DISPLAY_DEVICE
	device.cb = uint32(unsafe.Sizeof(device))
	i := uint32(0)
	for {
		res, _, _ := procEnumDisplayDevices.Call(uintptr(unsafe.Pointer(MustStringToUTF16Ptr(deviceName))), uintptr(i), uintptr(unsafe.Pointer(&device)), 0)
		if res == 0 {
			break
		}
		if device.StateFlags&0x1 != 0 {
			return syscall.UTF16ToString(device.DeviceString[:]), nil
		}
		i++
	}

	return "", fmt.Errorf("monitor name not found for device: %s", deviceName)
}

// I'm not convinced this works properly
func GetRotationForMonitor(displayName [32]uint16) (float32, error) {
	var devMode DEVMODE
	devMode.DmSize = uint16(unsafe.Sizeof(devMode))
	resp, _, _ := procEnumDisplaySettings.Call(uintptr(unsafe.Pointer(&displayName[0])), ENUM_CURRENT_SETTINGS, uintptr(unsafe.Pointer(&devMode)))
	if resp == 0 {
		return 0, fmt.Errorf("EnumDisplaySettings failed")
	}

	if (devMode.DmFields & DM_DISPLAYORIENTATION) == 0 {
		return 0, fmt.Errorf("DM_DISPLAYORIENTATION not set")
	}

	switch devMode.DmOrientation {
	case DMDO_DEFAULT:
		return 0, nil
	case DMDO_90:
		return 90, nil
	case DMDO_180:
		return 180, nil
	case DMDO_270:
		return 270, nil
	}

	return -1, nil
}

// cursorPosForScreens reports the current cursor position and whether the
// query succeeded. GetCursorPos returns FALSE (ERROR_ACCESS_DENIED) when the
// calling process is not attached to the interactive input desktop - e.g. the
// workstation is locked or on the secure/UAC desktop, or mid session switch
// (logon, RDP, fast-user-switch). Exposed as a package variable so tests can
// simulate that failure. Defaults to the real GetCursorPos syscall.
var cursorPosForScreens = func() (cursor POINT, ok bool) {
	ret, _, _ := procGetCursorPos.Call(uintptr(unsafe.Pointer(&cursor)))
	return cursor, ret != 0
}

// defaultMonitorDPI is the DPI that corresponds to a scale factor of 1.0. It is
// used for a monitor whose DPI cannot be queried.
const defaultMonitorDPI = 96

// monitorEnumeration collects the result of one EnumDisplayMonitors pass.
type monitorEnumeration struct {
	cursor   POINT
	cursorOK bool
	screens  []*Screen
	// failures describes every monitor that could not be read and was skipped.
	failures []string
}

// The EnumDisplayMonitors callback is created once. The Go runtime never frees
// the slot that syscall.NewCallback allocates and allows only a fixed number of
// them per process (about 2000), so creating a callback on every call would
// eventually be fatal: the screens are re-read on every display change. The
// state of the running pass is handed to the callback through activeEnumeration,
// which enumerationMu guards.
var (
	enumerationMu        sync.Mutex
	activeEnumeration    *monitorEnumeration
	enumMonitorsCallback = syscall.NewCallback(enumMonitorsProc)
)

func enumMonitorsProc(hMonitor uintptr, hdc uintptr, lprcMonitor *RECT, lParam uintptr) uintptr {
	activeEnumeration.add(hMonitor)
	return 1 // Always continue: one unreadable monitor must not hide the others
}

// add reads the details of a monitor and appends it to the result. A monitor
// that cannot be read is recorded in failures and skipped.
func (e *monitorEnumeration) add(hMonitor uintptr) {
	monitor := MONITORINFOEX{
		MONITORINFO: MONITORINFO{
			CbSize: uint32(unsafe.Sizeof(MONITORINFOEX{})),
		},
		SzDevice: [32]uint16{},
	}
	ret, _, callErr := procGetMonitorInfo.Call(hMonitor, uintptr(unsafe.Pointer(&monitor)))
	if ret == 0 {
		e.failures = append(e.failures, fmt.Sprintf("GetMonitorInfo failed for monitor 0x%X%s", hMonitor, lastErrorSuffix(callErr)))
		return
	}

	deviceName := syscall.UTF16ToString(monitor.SzDevice[:])
	screen := &Screen{
		MONITORINFOEX: monitor,
		HMonitor:      hMonitor,
		IsPrimary:     monitor.DwFlags == MONITORINFOF_PRIMARY,
		IsCurrent:     e.cursorOK && rectContainsPoint(monitor.RcMonitor, e.cursor),
	}

	// Get monitor name
	name, err := getMonitorName(deviceName)
	if err == nil {
		screen.Name = name
	}

	// Get DPI for monitor. The monitor can be reconfigured or removed while it
	// is enumerated, so a failure here falls back to the default DPI instead of
	// dropping the monitor. The layout is re-read on the next display change.
	var dpiX, dpiY uint
	if hr := GetDPIForMonitor(hMonitor, MDT_EFFECTIVE_DPI, &dpiX, &dpiY); hr != S_OK {
		dpiX = defaultMonitorDPI
	}
	// Convert to scale factor
	screen.ScaleFactor = float32(dpiX) / defaultMonitorDPI

	// Get rotation of monitor
	rot, err := GetRotationForMonitor(monitor.SzDevice)
	if err == nil {
		screen.Rotation = rot
	}

	e.screens = append(e.screens, screen)
}

// incompleteError reports a failure to read the primary monitor. The screen
// layout cannot be built without it, whereas a secondary monitor that failed
// can be left out. It returns nil when nothing went wrong, or when only
// non-primary monitors were skipped.
func (e *monitorEnumeration) incompleteError() error {
	if len(e.failures) == 0 {
		return nil
	}
	for _, screen := range e.screens {
		if screen.IsPrimary {
			return nil
		}
	}
	return fmt.Errorf("the primary monitor could not be read: %s", strings.Join(e.failures, "; "))
}

// lastErrorSuffix formats the last-error value returned by a Proc.Call for
// appending to an error message. Proc.Call always returns the value of
// GetLastError, so a zero value means the API did not report a reason.
func lastErrorSuffix(err error) string {
	if errno, ok := err.(syscall.Errno); !ok || errno == 0 {
		return ""
	}
	return fmt.Sprintf(": %v", err)
}

func GetAllScreens() ([]*Screen, error) {
	enumeration := &monitorEnumeration{}

	// The cursor position is used only to set the cosmetic IsCurrent flag
	// (which monitor the pointer is on). A failed query must NOT abort
	// enumeration: processAndCacheScreens() treats any error here as fatal at
	// startup (os.Exit), so a locked/secure desktop at launch would crash the
	// app before it ever shows. Degrade gracefully - no screen is marked
	// current, and the flag self-corrects on the next display-change re-cache.
	enumeration.cursor, enumeration.cursorOK = cursorPosForScreens()

	// Enumerate the monitors
	enumerationMu.Lock()
	defer enumerationMu.Unlock()
	activeEnumeration = enumeration
	defer func() { activeEnumeration = nil }()

	ret, _, callErr := procEnumDisplayMonitors.Call(0, 0, enumMonitorsCallback, 0)
	if ret == 0 {
		return nil, fmt.Errorf("EnumDisplayMonitors failed%s", lastErrorSuffix(callErr))
	}
	if err := enumeration.incompleteError(); err != nil {
		return nil, err
	}

	return enumeration.screens, nil
}

func rectContainsPoint(r RECT, p POINT) bool {
	return p.X >= r.Left && p.X < r.Right && p.Y >= r.Top && p.Y < r.Bottom
}
