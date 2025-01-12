//go:build windows

package w32

import (
	"fmt"
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

func GetAllScreens() ([]*Screen, error) {
	var result []*Screen
	var errMessage string

	// Get cursor position to determine the current monitor
	var cursor POINT
	ret, _, _ := procGetCursorPos.Call(uintptr(unsafe.Pointer(&cursor)))
	if ret == 0 {
		return nil, fmt.Errorf("GetCursorPos failed")
	}

	// Enumerate the monitors
	enumFunc := func(hMonitor uintptr, hdc uintptr, lprcMonitor *RECT, lParam uintptr) uintptr {
		monitor := MONITORINFOEX{
			MONITORINFO: MONITORINFO{
				CbSize: uint32(unsafe.Sizeof(MONITORINFOEX{})),
			},
			SzDevice: [32]uint16{},
		}
		ret, _, _ := procGetMonitorInfo.Call(hMonitor, uintptr(unsafe.Pointer(&monitor)))
		if ret == 0 {
			errMessage = "GetMonitorInfo failed"
			return 0 // Stop enumeration
		}

		screen := &Screen{
			MONITORINFOEX: monitor,
			HMonitor:      hMonitor,
			IsPrimary:     monitor.DwFlags == MONITORINFOF_PRIMARY,
			IsCurrent:     rectContainsPoint(monitor.RcMonitor, cursor),
		}

		// Get monitor name
		name, err := getMonitorName(syscall.UTF16ToString(monitor.SzDevice[:]))
		if err == nil {
			screen.Name = name
		}

		// Get DPI for monitor
		var dpiX, dpiY uint
		ret = GetDPIForMonitor(hMonitor, MDT_EFFECTIVE_DPI, &dpiX, &dpiY)
		if ret != S_OK {
			errMessage = "GetDpiForMonitor failed"
			return 0 // Stop enumeration
		}
		// Convert to scale factor
		screen.ScaleFactor = float32(dpiX) / 96.0

		// Get rotation of monitor
		rot, err := GetRotationForMonitor(monitor.SzDevice)
		if err == nil {
			screen.Rotation = rot
		}

		result = append(result, screen)
		return 1 // Continue enumeration
	}

	ret, _, _ = procEnumDisplayMonitors.Call(0, 0, syscall.NewCallback(enumFunc), 0)
	if ret == 0 {
		return nil, fmt.Errorf("EnumDisplayMonitors failed: %s", errMessage)
	}

	return result, nil
}

func rectContainsPoint(r RECT, p POINT) bool {
	return p.X >= r.Left && p.X < r.Right && p.Y >= r.Top && p.Y < r.Bottom
}
