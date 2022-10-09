package win32

import (
	"github.com/wailsapp/wails/v2/internal/system/operatingsystem"
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

var (
	modKernel32         = syscall.NewLazyDLL("kernel32.dll")
	procGetModuleHandle = modKernel32.NewProc("GetModuleHandleW")

	moduser32                         = syscall.NewLazyDLL("user32.dll")
	procRegisterClassEx               = moduser32.NewProc("RegisterClassExW")
	procLoadIcon                      = moduser32.NewProc("LoadIconW")
	procLoadCursor                    = moduser32.NewProc("LoadCursorW")
	procCreateWindowEx                = moduser32.NewProc("CreateWindowExW")
	procPostMessage                   = moduser32.NewProc("PostMessageW")
	procGetCursorPos                  = moduser32.NewProc("GetCursorPos")
	procSetForegroundWindow           = moduser32.NewProc("SetForegroundWindow")
	procCreatePopupMenu               = moduser32.NewProc("CreatePopupMenu")
	procTrackPopupMenu                = moduser32.NewProc("TrackPopupMenu")
	procDestroyMenu                   = moduser32.NewProc("DestroyMenu")
	procAppendMenuW                   = moduser32.NewProc("AppendMenuW")
	procCheckMenuItem                 = moduser32.NewProc("CheckMenuItem")
	procCheckMenuRadioItem            = moduser32.NewProc("CheckMenuRadioItem")
	procCreateIconFromResourceEx      = moduser32.NewProc("CreateIconFromResourceEx")
	procGetMessageW                   = moduser32.NewProc("GetMessageW")
	procIsDialogMessage               = moduser32.NewProc("IsDialogMessageW")
	procTranslateMessage              = moduser32.NewProc("TranslateMessage")
	procDispatchMessage               = moduser32.NewProc("DispatchMessageW")
	procPostQuitMessage               = moduser32.NewProc("PostQuitMessage")
	procSystemParametersInfo          = moduser32.NewProc("SystemParametersInfoW")
	procSetWindowCompositionAttribute = moduser32.NewProc("SetWindowCompositionAttribute")

	modshell32          = syscall.NewLazyDLL("shell32.dll")
	procShellNotifyIcon = modshell32.NewProc("Shell_NotifyIconW")

	moddwmapi                 = syscall.NewLazyDLL("dwmapi.dll")
	procDwmSetWindowAttribute = moddwmapi.NewProc("DwmSetWindowAttribute")

	moduxtheme         = syscall.NewLazyDLL("uxtheme.dll")
	procSetWindowTheme = moduxtheme.NewProc("SetWindowTheme")

	AllowDarkModeForWindow func(HWND, bool) uintptr
	SetPreferredAppMode    func(int32) uintptr
)

type PreferredAppMode = int32

const (
	PreferredAppModeDefault PreferredAppMode = iota
	PreferredAppModeAllowDark
	PreferredAppModeForceDark
	PreferredAppModeForceLight
	PreferredAppModeMax
)

/*
RtlGetNtVersionNumbers = void (LPDWORD major, LPDWORD minor, LPDWORD build) // 1809 17763
ShouldAppsUseDarkMode = bool () // ordinal 132
AllowDarkModeForWindow = bool (HWND hWnd, bool allow) // ordinal 133
AllowDarkModeForApp = bool (bool allow) // ordinal 135, removed since 18334
FlushMenuThemes = void () // ordinal 136
RefreshImmersiveColorPolicyState = void () // ordinal 104
IsDarkModeAllowedForWindow = bool (HWND hWnd) // ordinal 137
GetIsImmersiveColorUsingHighContrast = bool (IMMERSIVE_HC_CACHE_MODE mode) // ordinal 106
OpenNcThemeData = HTHEME (HWND hWnd, LPCWSTR pszClassList) // ordinal 49
// Insider 18290
ShouldSystemUseDarkMode = bool () // ordinal 138
// Insider 18334
SetPreferredAppMode = PreferredAppMode (PreferredAppMode appMode) // ordinal 135, since 18334
IsDarkModeAllowedForApp = bool () // ordinal 139
*/
func init() {
	if IsWindowsVersionAtLeast(10, 0, 18334) {

		// AllowDarkModeForWindow is only available on Windows 10+
		uxtheme, err := windows.LoadLibrary("uxtheme.dll")
		if err == nil {
			procAllowDarkModeForWindow, err := windows.GetProcAddressByOrdinal(uxtheme, uintptr(133))
			if err == nil {
				AllowDarkModeForWindow = func(hwnd HWND, allow bool) uintptr {
					var allowInt int32
					if allow {
						allowInt = 1
					}
					ret, _, _ := syscall.SyscallN(procAllowDarkModeForWindow, uintptr(hwnd), uintptr(allowInt))
					return ret
				}
			}
		}

		// SetPreferredAppMode is only available on Windows 10+
		procSetPreferredAppMode, err := windows.GetProcAddressByOrdinal(uxtheme, uintptr(135))
		if err == nil {
			SetPreferredAppMode = func(mode int32) uintptr {
				ret, _, _ := syscall.SyscallN(procSetPreferredAppMode, uintptr(mode))
				return ret
			}
			SetPreferredAppMode(PreferredAppModeAllowDark)
		}
	}

}

type HANDLE uintptr
type HINSTANCE = HANDLE
type HICON = HANDLE
type HCURSOR = HANDLE
type HBRUSH = HANDLE
type HWND = HANDLE
type HMENU = HANDLE
type DWORD = uint32
type ATOM uint16
type MenuID uint16

const (
	WM_ACTIVATE      = 0x0006
	WM_ACTIVATEAPP   = 0x001C
	WM_LBUTTONUP     = 0x0202
	WM_LBUTTONDBLCLK = 0x0203
	WM_RBUTTONUP     = 0x0205
	WM_USER          = 0x0400
	WM_TRAYICON      = WM_USER + 69
	WM_SETTINGCHANGE = 0x001A

	WS_EX_APPWINDOW           = 0x00040000
	WS_OVERLAPPEDWINDOW       = 0x00000000 | 0x00C00000 | 0x00080000 | 0x00040000 | 0x00020000 | 0x00010000
	WS_EX_NOREDIRECTIONBITMAP = 0x00200000
	CW_USEDEFAULT             = 0x80000000

	NIM_ADD        = 0x00000000
	NIM_MODIFY     = 0x00000001
	NIM_DELETE     = 0x00000002
	NIM_SETVERSION = 0x00000004

	NIF_MESSAGE = 0x00000001
	NIF_ICON    = 0x00000002
	NIF_TIP     = 0x00000004
	NIF_STATE   = 0x00000008
	NIF_INFO    = 0x00000010

	NIS_HIDDEN = 0x00000001

	NIIF_NONE               = 0x00000000
	NIIF_INFO               = 0x00000001
	NIIF_WARNING            = 0x00000002
	NIIF_ERROR              = 0x00000003
	NIIF_USER               = 0x00000004
	NIIF_NOSOUND            = 0x00000010
	NIIF_LARGE_ICON         = 0x00000020
	NIIF_RESPECT_QUIET_TIME = 0x00000080
	NIIF_ICON_MASK          = 0x0000000F

	IMAGE_BITMAP    = 0
	IMAGE_ICON      = 1
	LR_LOADFROMFILE = 0x00000010
	LR_DEFAULTSIZE  = 0x00000040

	IDC_ARROW     = 32512
	COLOR_WINDOW  = 5
	COLOR_BTNFACE = 15

	GWLP_USERDATA       = -21
	WS_CLIPSIBLINGS     = 0x04000000
	WS_EX_CONTROLPARENT = 0x00010000

	HWND_MESSAGE       = ^HWND(2)
	NOTIFYICON_VERSION = 4

	IDI_APPLICATION = 32512
	WM_APP          = 32768
	WM_COMMAND      = 273

	MenuItemMsgID       = WM_APP + 1024
	NotifyIconMessageId = WM_APP + iota

	MF_STRING       = 0x00000000
	MF_ENABLED      = 0x00000000
	MF_GRAYED       = 0x00000001
	MF_DISABLED     = 0x00000002
	MF_SEPARATOR    = 0x00000800
	MF_UNCHECKED    = 0x00000000
	MF_CHECKED      = 0x00000008
	MF_POPUP        = 0x00000010
	MF_MENUBARBREAK = 0x00000020
	MF_BYCOMMAND    = 0x00000000

	TPM_LEFTALIGN = 0x0000
	WM_NULL       = 0

	CS_VREDRAW = 0x0001
	CS_HREDRAW = 0x0002
)

var windowsVersion, _ = operatingsystem.GetWindowsVersionInfo()

func IsWindowsVersionAtLeast(major, minor, buildNumber int) bool {
	return windowsVersion.Major >= major &&
		windowsVersion.Minor >= minor &&
		windowsVersion.Build >= buildNumber
}

type WindowProc func(hwnd HWND, msg uint32, wparam, lparam uintptr) uintptr

func GetModuleHandle(value uintptr) uintptr {
	result, _, _ := procGetModuleHandle.Call(value)
	return result
}

func GetMessage(msg *MSG) uintptr {
	rt, _, _ := procGetMessageW.Call(uintptr(unsafe.Pointer(msg)), 0, 0, 0)
	return rt
}

func PostMessage(hwnd HWND, msg uint32, wParam, lParam uintptr) uintptr {
	ret, _, _ := procPostMessage.Call(
		uintptr(hwnd),
		uintptr(msg),
		wParam,
		lParam)

	return ret
}

func ShellNotifyIcon(cmd uintptr, nid *NOTIFYICONDATA) bool {
	ret, _, _ := procShellNotifyIcon.Call(cmd, uintptr(unsafe.Pointer(nid)))
	return ret == 1
}

func IsDialogMessage(hwnd HWND, msg *MSG) uintptr {
	ret, _, _ := procIsDialogMessage.Call(uintptr(hwnd), uintptr(unsafe.Pointer(msg)))
	return ret
}

func TranslateMessage(msg *MSG) uintptr {
	ret, _, _ := procTranslateMessage.Call(uintptr(unsafe.Pointer(msg)))
	return ret
}

func DispatchMessage(msg *MSG) uintptr {
	ret, _, _ := procDispatchMessage.Call(uintptr(unsafe.Pointer(msg)))
	return ret
}

func PostQuitMessage(exitCode int32) {
	procPostQuitMessage.Call(uintptr(exitCode))
}
