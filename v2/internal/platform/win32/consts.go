//go:build windows

package win32

import (
	"fmt"
	"syscall"
	"unsafe"

	"github.com/wailsapp/wails/v2/internal/system/operatingsystem"
	"golang.org/x/sys/windows"
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
	procGetKeyState                   = moduser32.NewProc("GetKeyState")
	procCreateAcceleratorTable        = moduser32.NewProc("CreateAcceleratorTableW")
	procTranslateAccelerator          = moduser32.NewProc("TranslateAcceleratorW")

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
	WM_APP                    = 32768
	WM_ACTIVATE               = 6
	WM_ACTIVATEAPP            = 28
	WM_AFXFIRST               = 864
	WM_AFXLAST                = 895
	WM_ASKCBFORMATNAME        = 780
	WM_CANCELJOURNAL          = 75
	WM_CANCELMODE             = 31
	WM_CAPTURECHANGED         = 533
	WM_CHANGECBCHAIN          = 781
	WM_CHAR                   = 258
	WM_CHARTOITEM             = 47
	WM_CHILDACTIVATE          = 34
	WM_CLEAR                  = 771
	WM_CLOSE                  = 16
	WM_COMMAND                = 273
	WM_COMMNOTIFY             = 68 /* OBSOLETE */
	WM_COMPACTING             = 65
	WM_COMPAREITEM            = 57
	WM_CONTEXTMENU            = 123
	WM_COPY                   = 769
	WM_COPYDATA               = 74
	WM_CREATE                 = 1
	WM_CTLCOLORBTN            = 309
	WM_CTLCOLORDLG            = 310
	WM_CTLCOLOREDIT           = 307
	WM_CTLCOLORLISTBOX        = 308
	WM_CTLCOLORMSGBOX         = 306
	WM_CTLCOLORSCROLLBAR      = 311
	WM_CTLCOLORSTATIC         = 312
	WM_CUT                    = 768
	WM_DEADCHAR               = 259
	WM_DELETEITEM             = 45
	WM_DESTROY                = 2
	WM_DESTROYCLIPBOARD       = 775
	WM_DEVICECHANGE           = 537
	WM_DEVMODECHANGE          = 27
	WM_DISPLAYCHANGE          = 126
	WM_DRAWCLIPBOARD          = 776
	WM_DRAWITEM               = 43
	WM_DROPFILES              = 563
	WM_ENABLE                 = 10
	WM_ENDSESSION             = 22
	WM_ENTERIDLE              = 289
	WM_ENTERMENULOOP          = 529
	WM_ENTERSIZEMOVE          = 561
	WM_ERASEBKGND             = 20
	WM_EXITMENULOOP           = 530
	WM_EXITSIZEMOVE           = 562
	WM_FONTCHANGE             = 29
	WM_GETDLGCODE             = 135
	WM_GETFONT                = 49
	WM_GETHOTKEY              = 51
	WM_GETICON                = 127
	WM_GETMINMAXINFO          = 36
	WM_GETTEXT                = 13
	WM_GETTEXTLENGTH          = 14
	WM_HANDHELDFIRST          = 856
	WM_HANDHELDLAST           = 863
	WM_HELP                   = 83
	WM_HOTKEY                 = 786
	WM_HSCROLL                = 276
	WM_HSCROLLCLIPBOARD       = 782
	WM_ICONERASEBKGND         = 39
	WM_INITDIALOG             = 272
	WM_INITMENU               = 278
	WM_INITMENUPOPUP          = 279
	WM_INPUT                  = 0x00FF
	WM_INPUTLANGCHANGE        = 81
	WM_INPUTLANGCHANGEREQUEST = 80
	WM_KEYDOWN                = 256
	WM_KEYUP                  = 257
	WM_KILLFOCUS              = 8
	WM_MDIACTIVATE            = 546
	WM_MDICASCADE             = 551
	WM_MDICREATE              = 544
	WM_MDIDESTROY             = 545
	WM_MDIGETACTIVE           = 553
	WM_MDIICONARRANGE         = 552
	WM_MDIMAXIMIZE            = 549
	WM_MDINEXT                = 548
	WM_MDIREFRESHMENU         = 564
	WM_MDIRESTORE             = 547
	WM_MDISETMENU             = 560
	WM_MDITILE                = 550
	WM_MEASUREITEM            = 44
	WM_GETOBJECT              = 0x003D
	WM_CHANGEUISTATE          = 0x0127
	WM_UPDATEUISTATE          = 0x0128
	WM_QUERYUISTATE           = 0x0129
	WM_UNINITMENUPOPUP        = 0x0125
	WM_MENURBUTTONUP          = 290
	WM_MENUCOMMAND            = 0x0126
	WM_MENUGETOBJECT          = 0x0124
	WM_MENUDRAG               = 0x0123
	WM_APPCOMMAND             = 0x0319
	WM_MENUCHAR               = 288
	WM_MENUSELECT             = 287
	WM_MOVE                   = 3
	WM_MOVING                 = 534
	WM_NCACTIVATE             = 134
	WM_NCCALCSIZE             = 131
	WM_NCCREATE               = 129
	WM_NCDESTROY              = 130
	WM_NCHITTEST              = 132
	WM_NCLBUTTONDBLCLK        = 163
	WM_NCLBUTTONDOWN          = 161
	WM_NCLBUTTONUP            = 162
	WM_NCMBUTTONDBLCLK        = 169
	WM_NCMBUTTONDOWN          = 167
	WM_NCMBUTTONUP            = 168
	WM_NCXBUTTONDOWN          = 171
	WM_NCXBUTTONUP            = 172
	WM_NCXBUTTONDBLCLK        = 173
	WM_NCMOUSEHOVER           = 0x02A0
	WM_NCMOUSELEAVE           = 0x02A2
	WM_NCMOUSEMOVE            = 160
	WM_NCPAINT                = 133
	WM_NCRBUTTONDBLCLK        = 166
	WM_NCRBUTTONDOWN          = 164
	WM_NCRBUTTONUP            = 165
	WM_NEXTDLGCTL             = 40
	WM_NEXTMENU               = 531
	WM_NOTIFY                 = 78
	WM_NOTIFYFORMAT           = 85
	WM_NULL                   = 0
	WM_PAINT                  = 15
	WM_PAINTCLIPBOARD         = 777
	WM_PAINTICON              = 38
	WM_PALETTECHANGED         = 785
	WM_PALETTEISCHANGING      = 784
	WM_PARENTNOTIFY           = 528
	WM_PASTE                  = 770
	WM_PENWINFIRST            = 896
	WM_PENWINLAST             = 911
	WM_POWER                  = 72
	WM_PRINT                  = 791
	WM_PRINTCLIENT            = 792
	WM_QUERYDRAGICON          = 55
	WM_QUERYENDSESSION        = 17
	WM_QUERYNEWPALETTE        = 783
	WM_QUERYOPEN              = 19
	WM_QUEUESYNC              = 35
	WM_QUIT                   = 18
	WM_RENDERALLFORMATS       = 774
	WM_RENDERFORMAT           = 773
	WM_SETCURSOR              = 32
	WM_SETFOCUS               = 7
	WM_SETFONT                = 48
	WM_SETHOTKEY              = 50
	WM_SETICON                = 128
	WM_SETREDRAW              = 11
	WM_SETTEXT                = 12
	WM_SETTINGCHANGE          = 26
	WM_SHOWWINDOW             = 24
	WM_SIZE                   = 5
	WM_SIZECLIPBOARD          = 779
	WM_SIZING                 = 532
	WM_SPOOLERSTATUS          = 42
	WM_STYLECHANGED           = 125
	WM_STYLECHANGING          = 124
	WM_SYSCHAR                = 262
	WM_SYSCOLORCHANGE         = 21
	WM_SYSCOMMAND             = 274
	WM_SYSDEADCHAR            = 263
	WM_SYSKEYDOWN             = 260
	WM_SYSKEYUP               = 261
	WM_TCARD                  = 82
	WM_THEMECHANGED           = 794
	WM_TIMECHANGE             = 30
	WM_TIMER                  = 275
	WM_UNDO                   = 772
	WM_USER                   = 1024
	WM_USERCHANGED            = 84
	WM_VKEYTOITEM             = 46
	WM_VSCROLL                = 277
	WM_VSCROLLCLIPBOARD       = 778
	WM_WINDOWPOSCHANGED       = 71
	WM_WINDOWPOSCHANGING      = 70
	WM_WININICHANGE           = 26
	WM_KEYFIRST               = 256
	WM_KEYLAST                = 264
	WM_SYNCPAINT              = 136
	WM_MOUSEACTIVATE          = 33
	WM_MOUSEMOVE              = 512
	WM_LBUTTONDOWN            = 513
	WM_LBUTTONUP              = 514
	WM_LBUTTONDBLCLK          = 515
	WM_RBUTTONDOWN            = 516
	WM_RBUTTONUP              = 517
	WM_RBUTTONDBLCLK          = 518
	WM_MBUTTONDOWN            = 519
	WM_MBUTTONUP              = 520
	WM_MBUTTONDBLCLK          = 521
	WM_MOUSEWHEEL             = 522
	WM_MOUSEFIRST             = 512
	WM_XBUTTONDOWN            = 523
	WM_XBUTTONUP              = 524
	WM_XBUTTONDBLCLK          = 525
	WM_MOUSELAST              = 525
	WM_MOUSEHOVER             = 0x2A1
	WM_MOUSELEAVE             = 0x2A3
	WM_CLIPBOARDUPDATE        = 0x031D

	WS_EX_APPWINDOW           = 0x00040000
	WS_OVERLAPPEDWINDOW       = 0x00000000 | 0x00C00000 | 0x00080000 | 0x00040000 | 0x00020000 | 0x00010000
	WS_EX_NOREDIRECTIONBITMAP = 0x00200000
	CW_USEDEFAULT             = ^0x7fffffff

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

	CS_VREDRAW = 0x0001
	CS_HREDRAW = 0x0002
)

func WMMessageToString(msg uintptr) string {
	// Convert windows message to string
	switch msg {
	case WM_APP:
		return "WM_APP"
	case WM_ACTIVATE:
		return "WM_ACTIVATE"
	case WM_ACTIVATEAPP:
		return "WM_ACTIVATEAPP"
	case WM_AFXFIRST:
		return "WM_AFXFIRST"
	case WM_AFXLAST:
		return "WM_AFXLAST"
	case WM_ASKCBFORMATNAME:
		return "WM_ASKCBFORMATNAME"
	case WM_CANCELJOURNAL:
		return "WM_CANCELJOURNAL"
	case WM_CANCELMODE:
		return "WM_CANCELMODE"
	case WM_CAPTURECHANGED:
		return "WM_CAPTURECHANGED"
	case WM_CHANGECBCHAIN:
		return "WM_CHANGECBCHAIN"
	case WM_CHAR:
		return "WM_CHAR"
	case WM_CHARTOITEM:
		return "WM_CHARTOITEM"
	case WM_CHILDACTIVATE:
		return "WM_CHILDACTIVATE"
	case WM_CLEAR:
		return "WM_CLEAR"
	case WM_CLOSE:
		return "WM_CLOSE"
	case WM_COMMAND:
		return "WM_COMMAND"
	case WM_COMMNOTIFY /* OBSOLETE */ :
		return "WM_COMMNOTIFY"
	case WM_COMPACTING:
		return "WM_COMPACTING"
	case WM_COMPAREITEM:
		return "WM_COMPAREITEM"
	case WM_CONTEXTMENU:
		return "WM_CONTEXTMENU"
	case WM_COPY:
		return "WM_COPY"
	case WM_COPYDATA:
		return "WM_COPYDATA"
	case WM_CREATE:
		return "WM_CREATE"
	case WM_CTLCOLORBTN:
		return "WM_CTLCOLORBTN"
	case WM_CTLCOLORDLG:
		return "WM_CTLCOLORDLG"
	case WM_CTLCOLOREDIT:
		return "WM_CTLCOLOREDIT"
	case WM_CTLCOLORLISTBOX:
		return "WM_CTLCOLORLISTBOX"
	case WM_CTLCOLORMSGBOX:
		return "WM_CTLCOLORMSGBOX"
	case WM_CTLCOLORSCROLLBAR:
		return "WM_CTLCOLORSCROLLBAR"
	case WM_CTLCOLORSTATIC:
		return "WM_CTLCOLORSTATIC"
	case WM_CUT:
		return "WM_CUT"
	case WM_DEADCHAR:
		return "WM_DEADCHAR"
	case WM_DELETEITEM:
		return "WM_DELETEITEM"
	case WM_DESTROY:
		return "WM_DESTROY"
	case WM_DESTROYCLIPBOARD:
		return "WM_DESTROYCLIPBOARD"
	case WM_DEVICECHANGE:
		return "WM_DEVICECHANGE"
	case WM_DEVMODECHANGE:
		return "WM_DEVMODECHANGE"
	case WM_DISPLAYCHANGE:
		return "WM_DISPLAYCHANGE"
	case WM_DRAWCLIPBOARD:
		return "WM_DRAWCLIPBOARD"
	case WM_DRAWITEM:
		return "WM_DRAWITEM"
	case WM_DROPFILES:
		return "WM_DROPFILES"
	case WM_ENABLE:
		return "WM_ENABLE"
	case WM_ENDSESSION:
		return "WM_ENDSESSION"
	case WM_ENTERIDLE:
		return "WM_ENTERIDLE"
	case WM_ENTERMENULOOP:
		return "WM_ENTERMENULOOP"
	case WM_ENTERSIZEMOVE:
		return "WM_ENTERSIZEMOVE"
	case WM_ERASEBKGND:
		return "WM_ERASEBKGND"
	case WM_EXITMENULOOP:
		return "WM_EXITMENULOOP"
	case WM_EXITSIZEMOVE:
		return "WM_EXITSIZEMOVE"
	case WM_FONTCHANGE:
		return "WM_FONTCHANGE"
	case WM_GETDLGCODE:
		return "WM_GETDLGCODE"
	case WM_GETFONT:
		return "WM_GETFONT"
	case WM_GETHOTKEY:
		return "WM_GETHOTKEY"
	case WM_GETICON:
		return "WM_GETICON"
	case WM_GETMINMAXINFO:
		return "WM_GETMINMAXINFO"
	case WM_GETTEXT:
		return "WM_GETTEXT"
	case WM_GETTEXTLENGTH:
		return "WM_GETTEXTLENGTH"
	case WM_HANDHELDFIRST:
		return "WM_HANDHELDFIRST"
	case WM_HANDHELDLAST:
		return "WM_HANDHELDLAST"
	case WM_HELP:
		return "WM_HELP"
	case WM_HOTKEY:
		return "WM_HOTKEY"
	case WM_HSCROLL:
		return "WM_HSCROLL"
	case WM_HSCROLLCLIPBOARD:
		return "WM_HSCROLLCLIPBOARD"
	case WM_ICONERASEBKGND:
		return "WM_ICONERASEBKGND"
	case WM_INITDIALOG:
		return "WM_INITDIALOG"
	case WM_INITMENU:
		return "WM_INITMENU"
	case WM_INITMENUPOPUP:
		return "WM_INITMENUPOPUP"
	case WM_INPUT:
		return "WM_INPUT"
	case WM_INPUTLANGCHANGE:
		return "WM_INPUTLANGCHANGE"
	case WM_INPUTLANGCHANGEREQUEST:
		return "WM_INPUTLANGCHANGEREQUEST"
	case WM_KEYDOWN:
		return "WM_KEYDOWN"
	case WM_KEYUP:
		return "WM_KEYUP"
	case WM_KILLFOCUS:
		return "WM_KILLFOCUS"
	case WM_MDIACTIVATE:
		return "WM_MDIACTIVATE"
	case WM_MDICASCADE:
		return "WM_MDICASCADE"
	case WM_MDICREATE:
		return "WM_MDICREATE"
	case WM_MDIDESTROY:
		return "WM_MDIDESTROY"
	case WM_MDIGETACTIVE:
		return "WM_MDIGETACTIVE"
	case WM_MDIICONARRANGE:
		return "WM_MDIICONARRANGE"
	case WM_MDIMAXIMIZE:
		return "WM_MDIMAXIMIZE"
	case WM_MDINEXT:
		return "WM_MDINEXT"
	case WM_MDIREFRESHMENU:
		return "WM_MDIREFRESHMENU"
	case WM_MDIRESTORE:
		return "WM_MDIRESTORE"
	case WM_MDISETMENU:
		return "WM_MDISETMENU"
	case WM_MDITILE:
		return "WM_MDITILE"
	case WM_MEASUREITEM:
		return "WM_MEASUREITEM"
	case WM_GETOBJECT:
		return "WM_GETOBJECT"
	case WM_CHANGEUISTATE:
		return "WM_CHANGEUISTATE"
	case WM_UPDATEUISTATE:
		return "WM_UPDATEUISTATE"
	case WM_QUERYUISTATE:
		return "WM_QUERYUISTATE"
	case WM_UNINITMENUPOPUP:
		return "WM_UNINITMENUPOPUP"
	case WM_MENURBUTTONUP:
		return "WM_MENURBUTTONUP"
	case WM_MENUCOMMAND:
		return "WM_MENUCOMMAND"
	case WM_MENUGETOBJECT:
		return "WM_MENUGETOBJECT"
	case WM_MENUDRAG:
		return "WM_MENUDRAG"
	case WM_APPCOMMAND:
		return "WM_APPCOMMAND"
	case WM_MENUCHAR:
		return "WM_MENUCHAR"
	case WM_MENUSELECT:
		return "WM_MENUSELECT"
	case WM_MOVE:
		return "WM_MOVE"
	case WM_MOVING:
		return "WM_MOVING"
	case WM_NCACTIVATE:
		return "WM_NCACTIVATE"
	case WM_NCCALCSIZE:
		return "WM_NCCALCSIZE"
	case WM_NCCREATE:
		return "WM_NCCREATE"
	case WM_NCDESTROY:
		return "WM_NCDESTROY"
	case WM_NCHITTEST:
		return "WM_NCHITTEST"
	case WM_NCLBUTTONDBLCLK:
		return "WM_NCLBUTTONDBLCLK"
	case WM_NCLBUTTONDOWN:
		return "WM_NCLBUTTONDOWN"
	case WM_NCLBUTTONUP:
		return "WM_NCLBUTTONUP"
	case WM_NCMBUTTONDBLCLK:
		return "WM_NCMBUTTONDBLCLK"
	case WM_NCMBUTTONDOWN:
		return "WM_NCMBUTTONDOWN"
	case WM_NCMBUTTONUP:
		return "WM_NCMBUTTONUP"
	case WM_NCXBUTTONDOWN:
		return "WM_NCXBUTTONDOWN"
	case WM_NCXBUTTONUP:
		return "WM_NCXBUTTONUP"
	case WM_NCXBUTTONDBLCLK:
		return "WM_NCXBUTTONDBLCLK"
	case WM_NCMOUSEHOVER:
		return "WM_NCMOUSEHOVER"
	case WM_NCMOUSELEAVE:
		return "WM_NCMOUSELEAVE"
	case WM_NCMOUSEMOVE:
		return "WM_NCMOUSEMOVE"
	case WM_NCPAINT:
		return "WM_NCPAINT"
	case WM_NCRBUTTONDBLCLK:
		return "WM_NCRBUTTONDBLCLK"
	case WM_NCRBUTTONDOWN:
		return "WM_NCRBUTTONDOWN"
	case WM_NCRBUTTONUP:
		return "WM_NCRBUTTONUP"
	case WM_NEXTDLGCTL:
		return "WM_NEXTDLGCTL"
	case WM_NEXTMENU:
		return "WM_NEXTMENU"
	case WM_NOTIFY:
		return "WM_NOTIFY"
	case WM_NOTIFYFORMAT:
		return "WM_NOTIFYFORMAT"
	case WM_NULL:
		return "WM_NULL"
	case WM_PAINT:
		return "WM_PAINT"
	case WM_PAINTCLIPBOARD:
		return "WM_PAINTCLIPBOARD"
	case WM_PAINTICON:
		return "WM_PAINTICON"
	case WM_PALETTECHANGED:
		return "WM_PALETTECHANGED"
	case WM_PALETTEISCHANGING:
		return "WM_PALETTEISCHANGING"
	case WM_PARENTNOTIFY:
		return "WM_PARENTNOTIFY"
	case WM_PASTE:
		return "WM_PASTE"
	case WM_PENWINFIRST:
		return "WM_PENWINFIRST"
	case WM_PENWINLAST:
		return "WM_PENWINLAST"
	case WM_POWER:
		return "WM_POWER"
	case WM_PRINT:
		return "WM_PRINT"
	case WM_PRINTCLIENT:
		return "WM_PRINTCLIENT"
	case WM_QUERYDRAGICON:
		return "WM_QUERYDRAGICON"
	case WM_QUERYENDSESSION:
		return "WM_QUERYENDSESSION"
	case WM_QUERYNEWPALETTE:
		return "WM_QUERYNEWPALETTE"
	case WM_QUERYOPEN:
		return "WM_QUERYOPEN"
	case WM_QUEUESYNC:
		return "WM_QUEUESYNC"
	case WM_QUIT:
		return "WM_QUIT"
	case WM_RENDERALLFORMATS:
		return "WM_RENDERALLFORMATS"
	case WM_RENDERFORMAT:
		return "WM_RENDERFORMAT"
	case WM_SETCURSOR:
		return "WM_SETCURSOR"
	case WM_SETFOCUS:
		return "WM_SETFOCUS"
	case WM_SETFONT:
		return "WM_SETFONT"
	case WM_SETHOTKEY:
		return "WM_SETHOTKEY"
	case WM_SETICON:
		return "WM_SETICON"
	case WM_SETREDRAW:
		return "WM_SETREDRAW"
	case WM_SETTEXT:
		return "WM_SETTEXT"
	case WM_SETTINGCHANGE:
		return "WM_SETTINGCHANGE"
	case WM_SHOWWINDOW:
		return "WM_SHOWWINDOW"
	case WM_SIZE:
		return "WM_SIZE"
	case WM_SIZECLIPBOARD:
		return "WM_SIZECLIPBOARD"
	case WM_SIZING:
		return "WM_SIZING"
	case WM_SPOOLERSTATUS:
		return "WM_SPOOLERSTATUS"
	case WM_STYLECHANGED:
		return "WM_STYLECHANGED"
	case WM_STYLECHANGING:
		return "WM_STYLECHANGING"
	case WM_SYSCHAR:
		return "WM_SYSCHAR"
	case WM_SYSCOLORCHANGE:
		return "WM_SYSCOLORCHANGE"
	case WM_SYSCOMMAND:
		return "WM_SYSCOMMAND"
	case WM_SYSDEADCHAR:
		return "WM_SYSDEADCHAR"
	case WM_SYSKEYDOWN:
		return "WM_SYSKEYDOWN"
	case WM_SYSKEYUP:
		return "WM_SYSKEYUP"
	case WM_TCARD:
		return "WM_TCARD"
	case WM_THEMECHANGED:
		return "WM_THEMECHANGED"
	case WM_TIMECHANGE:
		return "WM_TIMECHANGE"
	case WM_TIMER:
		return "WM_TIMER"
	case WM_UNDO:
		return "WM_UNDO"
	case WM_USER:
		return "WM_USER"
	case WM_USERCHANGED:
		return "WM_USERCHANGED"
	case WM_VKEYTOITEM:
		return "WM_VKEYTOITEM"
	case WM_VSCROLL:
		return "WM_VSCROLL"
	case WM_VSCROLLCLIPBOARD:
		return "WM_VSCROLLCLIPBOARD"
	case WM_WINDOWPOSCHANGED:
		return "WM_WINDOWPOSCHANGED"
	case WM_WINDOWPOSCHANGING:
		return "WM_WINDOWPOSCHANGING"
	case WM_KEYLAST:
		return "WM_KEYLAST"
	case WM_SYNCPAINT:
		return "WM_SYNCPAINT"
	case WM_MOUSEACTIVATE:
		return "WM_MOUSEACTIVATE"
	case WM_MOUSEMOVE:
		return "WM_MOUSEMOVE"
	case WM_LBUTTONDOWN:
		return "WM_LBUTTONDOWN"
	case WM_LBUTTONUP:
		return "WM_LBUTTONUP"
	case WM_LBUTTONDBLCLK:
		return "WM_LBUTTONDBLCLK"
	case WM_RBUTTONDOWN:
		return "WM_RBUTTONDOWN"
	case WM_RBUTTONUP:
		return "WM_RBUTTONUP"
	case WM_RBUTTONDBLCLK:
		return "WM_RBUTTONDBLCLK"
	case WM_MBUTTONDOWN:
		return "WM_MBUTTONDOWN"
	case WM_MBUTTONUP:
		return "WM_MBUTTONUP"
	case WM_MBUTTONDBLCLK:
		return "WM_MBUTTONDBLCLK"
	case WM_MOUSEWHEEL:
		return "WM_MOUSEWHEEL"
	case WM_XBUTTONDOWN:
		return "WM_XBUTTONDOWN"
	case WM_XBUTTONUP:
		return "WM_XBUTTONUP"
	case WM_MOUSELAST:
		return "WM_MOUSELAST"
	case WM_MOUSEHOVER:
		return "WM_MOUSEHOVER"
	case WM_MOUSELEAVE:
		return "WM_MOUSELEAVE"
	case WM_CLIPBOARDUPDATE:
		return "WM_CLIPBOARDUPDATE"
	default:
		return fmt.Sprintf("0x%08x", msg)
	}
}

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

func LoHiWords(input uint32) (uint16, uint16) {
	return uint16(input & 0xffff), uint16(input >> 16 & 0xffff)
}
