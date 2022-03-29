package win32

import (
	"golang.org/x/sys/windows/registry"
	"unsafe"
)

type DWMWINDOWATTRIBUTE int32

const DwmwaUseImmersiveDarkModeBefore20h1 DWMWINDOWATTRIBUTE = 19
const DwmwaUseImmersiveDarkMode DWMWINDOWATTRIBUTE = 20
const DwmwaBorderColor DWMWINDOWATTRIBUTE = 34
const DwmwaCaptionColor DWMWINDOWATTRIBUTE = 35
const DwmwaTextColor DWMWINDOWATTRIBUTE = 36

const SPI_GETHIGHCONTRAST = 0x0042
const HCF_HIGHCONTRASTON = 0x00000001

func dwmSetWindowAttribute(hwnd uintptr, dwAttribute DWMWINDOWATTRIBUTE, pvAttribute unsafe.Pointer, cbAttribute uintptr) {
	ret, _, err := procDwmSetWindowAttribute.Call(
		hwnd,
		uintptr(dwAttribute),
		uintptr(pvAttribute),
		cbAttribute)
	if ret != 0 {
		_ = err
		// println(err.Error())
	}
}

func SupportsThemes() bool {
	// We can't support Windows versions before 17763
	return IsWindowsVersionAtLeast(10, 0, 17763)
}

func SupportsCustomThemes() bool {
	return IsWindowsVersionAtLeast(10, 0, 17763)
}

func SetTheme(hwnd uintptr, useDarkMode bool) {
	if IsWindowsVersionAtLeast(10, 0, 17763) {
		attr := DwmwaUseImmersiveDarkModeBefore20h1
		if IsWindowsVersionAtLeast(10, 0, 18985) {
			attr = DwmwaUseImmersiveDarkMode
		}
		dwmSetWindowAttribute(hwnd, attr, unsafe.Pointer(&useDarkMode), unsafe.Sizeof(&useDarkMode))
	}
}

func SetTitleBarColour(hwnd uintptr, titleBarColour int32) {
	dwmSetWindowAttribute(hwnd, DwmwaCaptionColor, unsafe.Pointer(&titleBarColour), unsafe.Sizeof(titleBarColour))
}

func SetTitleTextColour(hwnd uintptr, titleTextColour int32) {
	dwmSetWindowAttribute(hwnd, DwmwaTextColor, unsafe.Pointer(&titleTextColour), unsafe.Sizeof(titleTextColour))
}

func SetBorderColour(hwnd uintptr, titleBorderColour int32) {
	dwmSetWindowAttribute(hwnd, DwmwaBorderColor, unsafe.Pointer(&titleBorderColour), unsafe.Sizeof(titleBorderColour))
}

func IsCurrentlyDarkMode() bool {
	key, err := registry.OpenKey(registry.CURRENT_USER, `SOFTWARE\Microsoft\Windows\CurrentVersion\Themes\Personalize`, registry.QUERY_VALUE)
	if err != nil {
		return false
	}
	defer key.Close()

	AppsUseLightTheme, _, err := key.GetIntegerValue("AppsUseLightTheme")
	if err != nil {
		return false
	}

	return AppsUseLightTheme == 0
}

type highContrast struct {
	CbSize            uint32
	DwFlags           uint32
	LpszDefaultScheme *int16
}

func IsCurrentlyHighContrastMode() bool {
	var result highContrast
	result.CbSize = uint32(unsafe.Sizeof(result))
	res, _, err := procSystemParametersInfo.Call(SPI_GETHIGHCONTRAST, uintptr(result.CbSize), uintptr(unsafe.Pointer(&result)), 0)
	if res == 0 {
		_ = err
		return false
	}
	return result.DwFlags&HCF_HIGHCONTRASTON == HCF_HIGHCONTRASTON
}
