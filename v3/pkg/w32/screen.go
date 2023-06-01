//go:build windows

package w32

import (
	"fmt"
	"syscall"
	"unsafe"
)

type Screen struct {
	MONITORINFOEX
	Name      string
	IsPrimary bool
	IsCurrent bool
	Scale     float32
	Rotation  float32
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
	var monitorList []MONITORINFOEX

	enumFunc := func(hMonitor uintptr, hdc uintptr, lprcMonitor *RECT, lParam uintptr) uintptr {
		monitor := MONITORINFOEX{
			MONITORINFO: MONITORINFO{
				CbSize: uint32(unsafe.Sizeof(MONITORINFOEX{})),
			},
			SzDevice: [32]uint16{},
		}
		ret, _, _ := procGetMonitorInfo.Call(hMonitor, uintptr(unsafe.Pointer(&monitor)))
		if ret == 0 {
			return 1 // Continue enumeration
		}

		monitorList = append(monitorList, monitor)
		return 1 // Continue enumeration
	}

	ret, _, _ := procEnumDisplayMonitors.Call(0, 0, syscall.NewCallback(enumFunc), 0)
	if ret == 0 {
		return nil, fmt.Errorf("EnumDisplayMonitors failed")
	}

	// Get the active screen
	var pt POINT
	ret, _, _ = procGetCursorPos.Call(uintptr(unsafe.Pointer(&pt)))
	if ret == 0 {
		return nil, fmt.Errorf("GetCursorPos failed")
	}

	hMonitor, _, _ := procMonitorFromPoint.Call(uintptr(unsafe.Pointer(&pt)), MONITOR_DEFAULTTONEAREST)
	if hMonitor == 0 {
		return nil, fmt.Errorf("MonitorFromPoint failed")
	}

	var monitorInfo MONITORINFO
	monitorInfo.CbSize = uint32(unsafe.Sizeof(monitorInfo))
	ret, _, _ = procGetMonitorInfo.Call(hMonitor, uintptr(unsafe.Pointer(&monitorInfo)))
	if ret == 0 {
		return nil, fmt.Errorf("GetMonitorInfo failed")
	}

	var result []*Screen

	// Iterate through the screens and set the active one
	for _, monitor := range monitorList {
		thisContainer := &Screen{
			MONITORINFOEX: monitor,
		}
		thisContainer.IsCurrent = equalRect(monitor.RcMonitor, monitorInfo.RcMonitor)
		thisContainer.IsPrimary = monitor.DwFlags == MONITORINFOF_PRIMARY
		name, err := getMonitorName(syscall.UTF16ToString(monitor.SzDevice[:]))
		if err != nil {
			name = ""
		}
		// Get DPI for monitor
		var dpiX, dpiY uint
		ret = GetDPIForMonitor(hMonitor, MDT_EFFECTIVE_DPI, &dpiX, &dpiY)
		if ret != S_OK {
			return nil, fmt.Errorf("GetDpiForMonitor failed")
		}
		// Convert to float32
		thisContainer.Scale = float32(dpiX) / 96.0

		// Get rotation of monitor
		rot, err := GetRotationForMonitor(monitor.SzDevice)
		if err != nil {
			rot = 0
		}
		thisContainer.Rotation = rot
		thisContainer.Name = name
		result = append(result, thisContainer)
	}

	return result, nil
}

func equalRect(a RECT, b RECT) bool {
	return a.Left == b.Left && a.Top == b.Top && a.Right == b.Right && a.Bottom == b.Bottom
}
