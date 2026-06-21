//go:build windows

package w32

import (
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
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

type WINDOWCOMPOSITIONATTRIB DWORD

type HTHEME HANDLE

const (
	WCA_UNDEFINED                     WINDOWCOMPOSITIONATTRIB = 0
	WCA_NCRENDERING_ENABLED           WINDOWCOMPOSITIONATTRIB = 1
	WCA_NCRENDERING_POLICY            WINDOWCOMPOSITIONATTRIB = 2
	WCA_TRANSITIONS_FORCEDISABLED     WINDOWCOMPOSITIONATTRIB = 3
	WCA_ALLOW_NCPAINT                 WINDOWCOMPOSITIONATTRIB = 4
	WCA_CAPTION_BUTTON_BOUNDS         WINDOWCOMPOSITIONATTRIB = 5
	WCA_NONCLIENT_RTL_LAYOUT          WINDOWCOMPOSITIONATTRIB = 6
	WCA_FORCE_ICONIC_REPRESENTATION   WINDOWCOMPOSITIONATTRIB = 7
	WCA_EXTENDED_FRAME_BOUNDS         WINDOWCOMPOSITIONATTRIB = 8
	WCA_HAS_ICONIC_BITMAP             WINDOWCOMPOSITIONATTRIB = 9
	WCA_THEME_ATTRIBUTES              WINDOWCOMPOSITIONATTRIB = 10
	WCA_NCRENDERING_EXILED            WINDOWCOMPOSITIONATTRIB = 11
	WCA_NCADORNMENTINFO               WINDOWCOMPOSITIONATTRIB = 12
	WCA_EXCLUDED_FROM_LIVEPREVIEW     WINDOWCOMPOSITIONATTRIB = 13
	WCA_VIDEO_OVERLAY_ACTIVE          WINDOWCOMPOSITIONATTRIB = 14
	WCA_FORCE_ACTIVEWINDOW_APPEARANCE WINDOWCOMPOSITIONATTRIB = 15
	WCA_DISALLOW_PEEK                 WINDOWCOMPOSITIONATTRIB = 16
	WCA_CLOAK                         WINDOWCOMPOSITIONATTRIB = 17
	WCA_CLOAKED                       WINDOWCOMPOSITIONATTRIB = 18
	WCA_ACCENT_POLICY                 WINDOWCOMPOSITIONATTRIB = 19
	WCA_FREEZE_REPRESENTATION         WINDOWCOMPOSITIONATTRIB = 20
	WCA_EVER_UNCLOAKED                WINDOWCOMPOSITIONATTRIB = 21
	WCA_VISUAL_OWNER                  WINDOWCOMPOSITIONATTRIB = 22
	WCA_HOLOGRAPHIC                   WINDOWCOMPOSITIONATTRIB = 23
	WCA_EXCLUDED_FROM_DDA             WINDOWCOMPOSITIONATTRIB = 24
	WCA_PASSIVEUPDATEMODE             WINDOWCOMPOSITIONATTRIB = 25
	WCA_USEDARKMODECOLORS             WINDOWCOMPOSITIONATTRIB = 26
	WCA_CORNER_STYLE                  WINDOWCOMPOSITIONATTRIB = 27
	WCA_PART_COLOR                    WINDOWCOMPOSITIONATTRIB = 28
	WCA_DISABLE_MOVESIZE_FEEDBACK     WINDOWCOMPOSITIONATTRIB = 29
	WCA_LAST                          WINDOWCOMPOSITIONATTRIB = 30
)

type WINDOWCOMPOSITIONATTRIBDATA struct {
	Attrib WINDOWCOMPOSITIONATTRIB
	PvData unsafe.Pointer
	CbData uintptr
}

var (
	uxtheme                         = syscall.NewLazyDLL("uxtheme.dll")
	procSetWindowTheme              = uxtheme.NewProc("SetWindowTheme")
	procOpenThemeData               = uxtheme.NewProc("OpenThemeData")
	procCloseThemeData              = uxtheme.NewProc("CloseThemeData")
	procDrawThemeBackground         = uxtheme.NewProc("DrawThemeBackground")
	procAllowDarkModeForApplication = uxtheme.NewProc("AllowDarkModeForApp")
	procDrawThemeTextEx             = uxtheme.NewProc("DrawThemeTextEx")
)

type PreferredAppMode = int32

const (
	PreferredAppModeDefault PreferredAppMode = iota
	PreferredAppModeAllowDark
	PreferredAppModeForceDark
	PreferredAppModeForceLight
	PreferredAppModeMax
)

var (
	AllowDarkModeForWindow           func(hwnd HWND, allow bool) uintptr
	SetPreferredAppMode              func(mode int32) uintptr
	FlushMenuThemes                  func()
	RefreshImmersiveColorPolicyState func()
	ShouldAppsUseDarkMode            func() bool
)

func init() {
	if IsWindowsVersionAtLeast(10, 0, 18334) {
		// AllowDarkModeForWindow is only available on Windows 10+
		localUXTheme, err := windows.LoadLibrary("uxtheme.dll")
		if err == nil {
			procAllowDarkModeForWindow, err := windows.GetProcAddressByOrdinal(localUXTheme, uintptr(133))
			if err == nil {
				AllowDarkModeForWindow = func(hwnd HWND, allow bool) uintptr {
					var allowInt int32
					if allow {
						allowInt = 1
					}
					ret, _, _ := syscall.SyscallN(procAllowDarkModeForWindow, uintptr(allowInt))
					return ret
				}
			}

			// Add ShouldAppsUseDarkMode
			procShouldAppsUseDarkMode, err := windows.GetProcAddressByOrdinal(localUXTheme, uintptr(132))
			if err == nil {
				ShouldAppsUseDarkMode = func() bool {
					ret, _, _ := syscall.SyscallN(procShouldAppsUseDarkMode)
					return ret != 0
				}
			}

			// SetPreferredAppMode is only available on Windows 10+
			procSetPreferredAppMode, err := windows.GetProcAddressByOrdinal(localUXTheme, uintptr(135))
			if err == nil {
				SetPreferredAppMode = func(mode int32) uintptr {
					ret, _, _ := syscall.SyscallN(procSetPreferredAppMode, uintptr(mode))
					return ret
				}
			}

			// Add FlushMenuThemes
			procFlushMenuThemesAddr, err := windows.GetProcAddressByOrdinal(localUXTheme, uintptr(136))
			if err == nil {
				FlushMenuThemes = func() {
					syscall.SyscallN(procFlushMenuThemesAddr)
				}
			}

			// Add RefreshImmersiveColorPolicyState
			procRefreshImmersiveColorPolicyStateAddr, err := windows.GetProcAddressByOrdinal(localUXTheme, uintptr(104))
			if err == nil {
				RefreshImmersiveColorPolicyState = func() {
					syscall.SyscallN(procRefreshImmersiveColorPolicyStateAddr)
				}
			}

			// Initialize dark mode
			if SetPreferredAppMode != nil {
				SetPreferredAppMode(PreferredAppModeAllowDark)
				if RefreshImmersiveColorPolicyState != nil {
					RefreshImmersiveColorPolicyState()
				}
			}

			windows.FreeLibrary(localUXTheme)
		}
	}
}

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

func SupportsBackdropTypes() bool {
	return IsWindowsVersionAtLeast(10, 0, 22621)
}

func SupportsImmersiveDarkMode() bool {
	return IsWindowsVersionAtLeast(10, 0, 18985)
}

func SetMenuTheme(hwnd uintptr, useDarkMode bool) {
	if !SupportsThemes() {
		return
	}

	// Check if dark mode is supported and enabled
	if useDarkMode && ShouldAppsUseDarkMode != nil && !ShouldAppsUseDarkMode() {
		useDarkMode = false
	}

	// Set the window theme
	themeName := "Explorer"
	if useDarkMode {
		themeName = "DarkMode_Explorer"
	}
	procSetWindowTheme.Call(HWND(hwnd), uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(themeName))), 0)

	// Update the theme state
	if RefreshImmersiveColorPolicyState != nil {
		RefreshImmersiveColorPolicyState()
	}

	// Flush menu themes to force a refresh
	if FlushMenuThemes != nil {
		FlushMenuThemes()
	}

	// Set dark mode for the window
	if AllowDarkModeForWindow != nil {
		AllowDarkModeForWindow(HWND(hwnd), useDarkMode)
	}

	// Force a redraw
	InvalidateRect(HWND(hwnd), nil, true)
}

func SetTheme(hwnd uintptr, useDarkMode bool) {
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
		SetMenuTheme(hwnd, useDarkMode)
	}
}

func EnableTranslucency(hwnd uintptr, backdrop uint32) {
	dwmSetWindowAttribute(hwnd, DwmwaSystemBackdropType, unsafe.Pointer(&backdrop), unsafe.Sizeof(backdrop))
}

func SetTitleBarColour(hwnd uintptr, titleBarColour uint32) {
	// Debug: Print the color value being set
	// fmt.Printf("Setting titlebar color to: 0x%08X (RGB: %d, %d, %d)\n", titleBarColour, titleBarColour&0xFF, (titleBarColour>>8)&0xFF, (titleBarColour>>16)&0xFF)
	dwmSetWindowAttribute(hwnd, DwmwaCaptionColor, unsafe.Pointer(&titleBarColour), unsafe.Sizeof(titleBarColour))
}

func SetTitleTextColour(hwnd uintptr, titleTextColour uint32) {
	dwmSetWindowAttribute(hwnd, DwmwaTextColor, unsafe.Pointer(&titleTextColour), unsafe.Sizeof(titleTextColour))
}

func SetBorderColour(hwnd uintptr, titleBorderColour uint32) {
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
	r := result.DwFlags&HCF_HIGHCONTRASTON == HCF_HIGHCONTRASTON
	return r
}

func GetAccentColor() (string, error) {
	key, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\DWM`, registry.QUERY_VALUE)
	if err != nil {
		return "", err
	}
	defer key.Close()

	accentColor, _, err := key.GetIntegerValue("AccentColor")
	if err != nil {
		return "", err
	}

	// Extract RGB components from ABGR format (Alpha, Blue, Green, Red)
	red := uint8(accentColor & 0xFF)
	green := uint8((accentColor >> 8) & 0xFF)
	blue := uint8((accentColor >> 16) & 0xFF)

	return fmt.Sprintf("rgb(%d,%d,%d)", red, green, blue), nil
}

// OpenThemeData opens theme data for a window and its class
func OpenThemeData(hwnd HWND, pszClassList string) HTHEME {
	ret, _, _ := procOpenThemeData.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(pszClassList))))
	return HTHEME(ret)
}

// CloseThemeData closes theme data handle
func CloseThemeData(hTheme HTHEME) error {
	ret, _, err := procCloseThemeData.Call(uintptr(hTheme))
	if ret != 0 {
		return err
	}
	return nil
}

// DrawThemeTextEx draws theme text with extended options
func DrawThemeTextEx(hTheme HTHEME, hdc HDC, iPartId int32, iStateId int32, pszText []uint16, cchText int32, dwTextFlags uint32, pRect *RECT, pOptions *DTTOPTS) error {
	ret, _, err := procDrawThemeTextEx.Call(
		uintptr(hTheme),
		uintptr(hdc),
		uintptr(iPartId),
		uintptr(iStateId),
		uintptr(unsafe.Pointer(&pszText[0])),
		uintptr(cchText),
		uintptr(dwTextFlags),
		uintptr(unsafe.Pointer(pRect)),
		uintptr(unsafe.Pointer(pOptions)))
	if ret != 0 {
		return err
	}
	return nil
}
