//go:build windows

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
const DwmwaSystemBackdropType DWMWINDOWATTRIBUTE = 38

const SPI_GETHIGHCONTRAST = 0x0042
const HCF_HIGHCONTRASTON = 0x00000001
const WCA_ACCENT_POLICY WINDOWCOMPOSITIONATTRIB = 19

type ACCENT_STATE DWORD

const (
	ACCENT_DISABLED                   ACCENT_STATE = 0
	ACCENT_ENABLE_GRADIENT            ACCENT_STATE = 1
	ACCENT_ENABLE_TRANSPARENTGRADIENT ACCENT_STATE = 2
	ACCENT_ENABLE_BLURBEHIND          ACCENT_STATE = 3
	ACCENT_ENABLE_ACRYLICBLURBEHIND   ACCENT_STATE = 4 // RS4 1803
	ACCENT_ENABLE_HOSTBACKDROP        ACCENT_STATE = 5 // RS5 1809
	ACCENT_INVALID_STATE              ACCENT_STATE = 6
)

type ACCENT_POLICY struct {
	AccentState   ACCENT_STATE
	AccentFlags   DWORD
	GradientColor DWORD
	AnimationId   DWORD
}

type WINDOWCOMPOSITIONATTRIBDATA struct {
	Attrib WINDOWCOMPOSITIONATTRIB
	PvData unsafe.Pointer
	CbData uintptr
}

type WINDOWCOMPOSITIONATTRIB DWORD

// BackdropType defines the type of translucency we wish to use
type BackdropType int32

const (
	BackdropTypeAuto    BackdropType = 0
	BackdropTypeNone    BackdropType = 1
	BackdropTypeMica    BackdropType = 2
	BackdropTypeAcrylic BackdropType = 3
	BackdropTypeTabbed  BackdropType = 4
)

func dwmSetWindowAttribute(hwnd HWND, dwAttribute DWMWINDOWATTRIBUTE, pvAttribute unsafe.Pointer, cbAttribute uintptr) {
	ret, _, err := procDwmSetWindowAttribute.Call(
		uintptr(hwnd),
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

func SupportsBackdropTypes() bool {
	return IsWindowsVersionAtLeast(10, 0, 22621)
}

func SupportsImmersiveDarkMode() bool {
	return IsWindowsVersionAtLeast(10, 0, 18985)
}

func SetTheme(hwnd HWND, useDarkMode bool) {
	if SupportsThemes() {
		attr := DwmwaUseImmersiveDarkModeBefore20h1
		if SupportsImmersiveDarkMode() {
			attr = DwmwaUseImmersiveDarkMode
		}
		var winDark int32
		if useDarkMode {
			winDark = 1
		}
		dwmSetWindowAttribute(hwnd, attr, unsafe.Pointer(&winDark), unsafe.Sizeof(winDark))
	}
}

func EnableBlurBehind(hwnd HWND) {
	var accent = ACCENT_POLICY{
		AccentState: ACCENT_ENABLE_ACRYLICBLURBEHIND,
		AccentFlags: 0x2,
	}
	var data WINDOWCOMPOSITIONATTRIBDATA
	data.Attrib = WCA_ACCENT_POLICY
	data.PvData = unsafe.Pointer(&accent)
	data.CbData = unsafe.Sizeof(accent)

	SetWindowCompositionAttribute(hwnd, &data)
}

func SetWindowCompositionAttribute(hwnd HWND, data *WINDOWCOMPOSITIONATTRIBDATA) bool {
	if procSetWindowCompositionAttribute != nil {
		ret, _, _ := procSetWindowCompositionAttribute.Call(
			uintptr(hwnd),
			uintptr(unsafe.Pointer(data)),
		)
		return ret != 0
	}
	return false
}

func EnableTranslucency(hwnd HWND, backdrop BackdropType) {
	if SupportsBackdropTypes() {
		dwmSetWindowAttribute(hwnd, DwmwaSystemBackdropType, unsafe.Pointer(&backdrop), unsafe.Sizeof(backdrop))
	} else {
		println("Warning: Translucency type unavailable on Windows < 22621")
	}
}

func SetTitleBarColour(hwnd HWND, titleBarColour int32) {
	dwmSetWindowAttribute(hwnd, DwmwaCaptionColor, unsafe.Pointer(&titleBarColour), unsafe.Sizeof(titleBarColour))
}

func SetTitleTextColour(hwnd HWND, titleTextColour int32) {
	dwmSetWindowAttribute(hwnd, DwmwaTextColor, unsafe.Pointer(&titleTextColour), unsafe.Sizeof(titleTextColour))
}

func SetBorderColour(hwnd HWND, titleBorderColour int32) {
	dwmSetWindowAttribute(hwnd, DwmwaBorderColor, unsafe.Pointer(&titleBorderColour), unsafe.Sizeof(titleBorderColour))
}

func SetWindowTheme(hwnd HWND, appName string, subIdList string) uintptr {
	var subID uintptr
	if subIdList != "" {
		subID = MustStringToUTF16uintptr(subIdList)
	}
	ret, _, _ := procSetWindowTheme.Call(
		uintptr(hwnd),
		MustStringToUTF16uintptr(appName),
		subID,
	)

	return ret
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
	r := result.DwFlags&HCF_HIGHCONTRASTON == HCF_HIGHCONTRASTON
	return r
}
