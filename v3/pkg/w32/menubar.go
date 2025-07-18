package w32

import (
	"unsafe"
)

const (
	OBJID_MENU = -3
	ODT_MENU   = 1
	// Menu info flags
	MIIM_BACKGROUND      = 0x00000002
	MIIM_APPLYTOSUBMENUS = 0x80000000
)

var (
	menuTheme       HTHEME
	procSetMenuInfo = moduser32.NewProc("SetMenuInfo")
)

type DTTOPTS struct {
	DwSize              uint32
	DwFlags             uint32
	CrText              uint32
	CrBorder            uint32
	CrShadow            uint32
	ITextShadowType     int32
	PtShadowOffset      POINT
	iBorderSize         int32
	iFontPropId         int32
	IColorPropId        int32
	IStateId            int32
	FApplyOverlay       int32
	IGlowSize           int32
	PfnDrawTextCallback uintptr
	LParam              uintptr
}

const (
	MENU_POPUPITEM = 14
	DTT_TEXTCOLOR  = 1
)

// Menu item states
const (
	ODS_SELECTED    = 0x0001
	ODS_GRAYED      = 0x0002
	ODS_DISABLED    = 0x0004
	ODS_CHECKED     = 0x0008
	ODS_FOCUS       = 0x0010
	ODS_DEFAULT     = 0x0020
	ODS_HOTLIGHT    = 0x0040
	ODS_INACTIVE    = 0x0080
	ODS_NOACCEL     = 0x0100
	ODS_NOFOCUSRECT = 0x0200
)

// Menu Button Image states
const (
	MBI_NORMAL   = 1
	MBI_HOT      = 2
	MBI_PUSHED   = 3
	MBI_DISABLED = 4
)

var (
	procGetMenuItemInfo  = moduser32.NewProc("GetMenuItemInfoW")
	procGetMenuItemCount = moduser32.NewProc("GetMenuItemCount")
	procGetMenuItemRect  = moduser32.NewProc("GetMenuItemRect")
)

func GetMenuItemInfo(hmenu HMENU, item uint32, fByPosition bool, lpmii *MENUITEMINFO) bool {
	ret, _, _ := procGetMenuItemInfo.Call(
		uintptr(hmenu),
		uintptr(item),
		uintptr(boolToUint(fByPosition)),
		uintptr(unsafe.Pointer(lpmii)),
	)
	return ret != 0
}

func GetMenuItemCount(hmenu HMENU) int {
	ret, _, _ := procGetMenuItemCount.Call(uintptr(hmenu))
	return int(ret)
}

func GetMenuItemRect(hwnd HWND, hmenu HMENU, item uint32, rect *RECT) bool {
	ret, _, _ := procGetMenuItemRect.Call(
		uintptr(hwnd),
		uintptr(hmenu),
		uintptr(item),
		uintptr(unsafe.Pointer(rect)),
	)
	return ret != 0
}

// Helper function to convert bool to uint
func boolToUint(b bool) uint {
	if b {
		return 1
	}
	return 0
}

type UAHMENU struct {
	Hmenu   HMENU
	Hdc     HDC
	DwFlags uint32
}

type MENUBARINFO struct {
	CbSize     uint32
	Bar        RECT
	Menu       HMENU
	Window     HWND
	BarFocused int32
	Focused    int32
}

type DRAWITEMSTRUCT struct {
	ControlType uint32
	ControlID   uint32
	ItemID      uint32
	ItemAction  uint32
	ItemState   uint32
	HWNDItem    HWND
	HDC         HDC
	RcItem      RECT
	ItemData    uintptr
}

type UAHDRAWMENUITEM struct {
	DIS  DRAWITEMSTRUCT
	UM   UAHMENU
	UAMI UAHMENUITEM
}

type UAHMENUITEM struct {
	Position int
	Umim     UAHMENUITEMMETRICS
	Umpm     UAHMENUPOPUPMETRICS
}
type UAHMENUITEMMETRICS struct {
	data [32]byte // Total size of the union in bytes (4 DWORDs * 4 bytes each * 2 arrays)
}

func (u *UAHMENUITEMMETRICS) RgsizeBar() *[2]struct{ cx, cy uint32 } {
	return (*[2]struct{ cx, cy uint32 })(unsafe.Pointer(&u.data))
}

func (u *UAHMENUITEMMETRICS) RgsizePopup() *[4]struct{ cx, cy uint32 } {
	return (*[4]struct{ cx, cy uint32 })(unsafe.Pointer(&u.data))
}

type UAHMEASUREMENUITEM struct {
	UM   UAHMENU
	UAMI UAHMENUITEM
	Mis  MEASUREITEMSTRUCT
}

type MEASUREITEMSTRUCT struct {
	CtlType    uint32
	CtlID      uint32
	ItemID     uint32
	ItemWidth  uint32
	ItemHeight uint32
	ItemData   uintptr
}

type UAHMENUPOPUPMETRICS struct {
	Rgcx             [4]uint32 // Array of 4 DWORDs
	FUpdateMaxWidths uint32    // Bit-field represented as a uint32
}

// Helper function to get the value of the fUpdateMaxWidths bit-field
func (u *UAHMENUPOPUPMETRICS) GetFUpdateMaxWidths() uint32 {
	return u.FUpdateMaxWidths & 0x3 // Mask to get the first 2 bits
}

// Helper function to set the value of the fUpdateMaxWidths bit-field
func (u *UAHMENUPOPUPMETRICS) SetFUpdateMaxWidths(value uint32) {
	u.FUpdateMaxWidths = (u.FUpdateMaxWidths &^ 0x3) | (value & 0x3) // Clear and set the first 2 bits
}

type MenuBarTheme struct {
	TitleBarBackground     *uint32
	TitleBarText           *uint32
	MenuBarBackground      *uint32 // Separate color for menubar
	MenuHoverBackground    *uint32
	MenuHoverText          *uint32
	MenuSelectedBackground *uint32
	MenuSelectedText       *uint32

	// private brushes
	titleBarBackgroundBrush     HBRUSH
	menuBarBackgroundBrush      HBRUSH // Separate brush for menubar
	menuHoverBackgroundBrush    HBRUSH
	menuSelectedBackgroundBrush HBRUSH
}

func createColourWithDefaultColor(color *uint32, def uint32) *uint32 {
	if color == nil {
		return &def
	}
	return color
}

func (d *MenuBarTheme) Init() {
	d.TitleBarBackground = createColourWithDefaultColor(d.TitleBarBackground, RGB(25, 25, 26))
	d.TitleBarText = createColourWithDefaultColor(d.TitleBarText, RGB(222, 222, 222))
	d.MenuBarBackground = createColourWithDefaultColor(d.MenuBarBackground, RGB(33, 33, 33))
	d.MenuSelectedText = createColourWithDefaultColor(d.MenuSelectedText, RGB(222, 222, 222))
	d.MenuSelectedBackground = createColourWithDefaultColor(d.MenuSelectedBackground, RGB(48, 48, 48))
	d.MenuHoverText = createColourWithDefaultColor(d.MenuHoverText, RGB(222, 222, 222))
	d.MenuHoverBackground = createColourWithDefaultColor(d.MenuHoverBackground, RGB(48, 48, 48))
	// Create brushes
	d.titleBarBackgroundBrush = CreateSolidBrush(*d.TitleBarBackground)
	d.menuBarBackgroundBrush = CreateSolidBrush(*d.MenuBarBackground)
	d.menuHoverBackgroundBrush = CreateSolidBrush(*d.MenuHoverBackground)
	d.menuSelectedBackgroundBrush = CreateSolidBrush(*d.MenuSelectedBackground)
}

// SetMenuBackground sets the menu background brush directly
func (d *MenuBarTheme) SetMenuBackground(hmenu HMENU) {
	var mi MENUINFO
	mi.CbSize = uint32(unsafe.Sizeof(mi))
	mi.FMask = MIIM_BACKGROUND | MIIM_APPLYTOSUBMENUS
	mi.HbrBack = d.menuBarBackgroundBrush // Use separate menubar brush
	SetMenuInfo(hmenu, &mi)
}

// SetMenuInfo wrapper function
func SetMenuInfo(hmenu HMENU, lpcmi *MENUINFO) bool {
	ret, _, _ := procSetMenuInfo.Call(
		uintptr(hmenu),
		uintptr(unsafe.Pointer(lpcmi)))
	return ret != 0
}

func CreateSolidBrush(color COLORREF) HBRUSH {
	ret, _, _ := procCreateSolidBrush.Call(
		uintptr(color),
	)
	return HBRUSH(ret)
}

func RGB(r, g, b byte) uint32 {
	return uint32(r) | uint32(g)<<8 | uint32(b)<<16
}

func RGBptr(r, g, b byte) *uint32 {
	result := uint32(r) | uint32(g)<<8 | uint32(b)<<16
	return &result
}

func MenuBarWndProc(hwnd HWND, msg uint32, wParam WPARAM, lParam LPARAM, theme *MenuBarTheme) (bool, LRESULT) {
	// Only proceed if we have a theme (either for dark or light mode)
	if theme == nil {
		return false, 0
	}
	switch msg {
	case WM_UAHDRAWMENU:
		udm := (*UAHMENU)(unsafe.Pointer(lParam))

		// get the menubar rect
		var menuBarInfo MENUBARINFO
		menuBarInfo.CbSize = uint32(unsafe.Sizeof(menuBarInfo))
		if !GetMenuBarInfo(hwnd, OBJID_MENU, 0, &menuBarInfo) {
			return false, 0
		}

		winRect := GetWindowRect(hwnd)

		// the rcBar is offset by the window rect
		rc := menuBarInfo.Bar
		OffsetRect(&rc, int(-winRect.Left), int(-winRect.Top))

		// Fill the entire menubar background with dark color
		FillRect(udm.Hdc, &rc, theme.menuBarBackgroundBrush)

		// Paint over the menubar border explicitly
		// The border is typically 1-2 pixels at the bottom
		borderRect := rc
		borderRect.Top = borderRect.Bottom - 1
		borderRect.Bottom = borderRect.Bottom + 2
		FillRect(udm.Hdc, &borderRect, theme.menuBarBackgroundBrush)

		return true, 0
	case WM_DRAWITEM:
		// Handle owner-drawn menu items
		dis := (*DRAWITEMSTRUCT)(unsafe.Pointer(lParam))

		// Check if this is a menu item
		if dis.ControlType == ODT_MENU {
			// Draw the menu item background
			var bgBrush HBRUSH
			var textColor uint32

			if dis.ItemState&ODS_SELECTED != 0 {
				// Selected state
				bgBrush = theme.menuSelectedBackgroundBrush
				textColor = *theme.MenuSelectedText
			} else {
				// Normal state
				bgBrush = theme.titleBarBackgroundBrush
				textColor = *theme.TitleBarText
			}

			// Fill background
			FillRect(dis.HDC, &dis.RcItem, bgBrush)

			// Draw text if we have item data
			if dis.ItemData != 0 {
				text := (*uint16)(unsafe.Pointer(dis.ItemData))
				if text != nil {
					// Set text color and draw
					SetTextColor(dis.HDC, COLORREF(textColor))
					SetBkMode(dis.HDC, TRANSPARENT)
					DrawText(dis.HDC, (*[256]uint16)(unsafe.Pointer(text))[:], -1, &dis.RcItem, DT_CENTER|DT_SINGLELINE|DT_VCENTER)
				}
			}

			return true, 1
		}
	case WM_UAHDRAWMENUITEM:
		udmi := (*UAHDRAWMENUITEM)(unsafe.Pointer(lParam))

		// Create buffer for menu text
		menuString := make([]uint16, 256)

		// Setup menu item info structure
		mii := MENUITEMINFO{
			CbSize:     uint32(unsafe.Sizeof(MENUITEMINFO{})),
			FMask:      MIIM_STRING | MIIM_SUBMENU,
			DwTypeData: &menuString[0],
			Cch:        uint32(len(menuString) - 1),
		}

		if !GetMenuItemInfo(udmi.UM.Hmenu, uint32(udmi.UAMI.Position), true, &mii) {
			// Failed to get menu item info, let default handler process
			return false, 0
		}

		// Remove automatic popup on hover - menus should only open on click
		// This was causing the menu to appear at wrong coordinates
		dwFlags := uint32(DT_CENTER | DT_SINGLELINE | DT_VCENTER)

		// Check if this is a menubar item (dwFlags will be 0 for menubar items)
		isMenuBarItem := udmi.UM.DwFlags == 0

		// Use different colors for menubar vs popup items
		var bgBrush HBRUSH
		var textColor uint32

		if udmi.DIS.ItemState&ODS_HOTLIGHT != 0 {
			// Hot state - use a specific color for hover
			bgBrush = theme.menuHoverBackgroundBrush
			textColor = *theme.MenuHoverText
		} else if udmi.DIS.ItemState&ODS_SELECTED != 0 {
			// Selected state
			bgBrush = theme.menuSelectedBackgroundBrush
			textColor = *theme.MenuSelectedText
		} else {
			// Normal state
			if isMenuBarItem {
				// Menubar items in normal state
				bgBrush = theme.menuBarBackgroundBrush
				textColor = *theme.TitleBarText
			} else {
				// Popup menu items in normal state - use same color as menubar
				bgBrush = theme.menuBarBackgroundBrush
				textColor = *theme.TitleBarText
			}
		}

		// Fill background
		FillRect(udmi.UM.Hdc, &udmi.DIS.RcItem, bgBrush)

		// Draw text
		// For menubar items, we need to ensure the text color is set correctly
		SetTextColor(udmi.UM.Hdc, COLORREF(textColor))
		SetBkMode(udmi.UM.Hdc, TRANSPARENT)

		// Draw text using UTF16 string directly
		// Pass -1 to let DrawText calculate the string length automatically
		DrawText(udmi.UM.Hdc, menuString, -1, &udmi.DIS.RcItem, dwFlags)

		// Return 1 to indicate we've handled the drawing
		return true, 1
	case WM_UAHMEASUREMENUITEM:
		// Let the default window procedure handle the menu item measurement
		// We're not modifying the default sizing anymore
		result := DefWindowProc(hwnd, msg, wParam, lParam)

		return true, result
	case WM_NCPAINT:
		// Paint our custom menubar first
		paintDarkMenuBar(hwnd, theme)

		// Then let Windows do its default painting
		result := DefWindowProc(hwnd, msg, wParam, lParam)

		// Paint again to ensure our painting is on top
		paintDarkMenuBar(hwnd, theme)

		return true, result
	case WM_NCACTIVATE:
		result := DefWindowProc(hwnd, msg, wParam, lParam)

		// Force paint the menubar with dark background
		paintDarkMenuBar(hwnd, theme)

		return false, result
	case WM_PAINT:
		// Let Windows paint first
		result := DefWindowProc(hwnd, msg, wParam, lParam)

		// Then paint our menubar
		paintDarkMenuBar(hwnd, theme)

		return false, result
	case WM_ACTIVATEAPP, WM_ACTIVATE:
		// Handle app activation/deactivation
		result := DefWindowProc(hwnd, msg, wParam, lParam)

		// Repaint menubar
		paintDarkMenuBar(hwnd, theme)

		return false, result
	case WM_SIZE, WM_WINDOWPOSCHANGED:
		// Handle window size changes
		result := DefWindowProc(hwnd, msg, wParam, lParam)

		// Repaint menubar after size change
		paintDarkMenuBar(hwnd, theme)

		return false, result
	case WM_SETFOCUS, WM_KILLFOCUS:
		// Handle focus changes (e.g., when inspector opens)
		result := DefWindowProc(hwnd, msg, wParam, lParam)

		// Repaint menubar after focus change
		paintDarkMenuBar(hwnd, theme)

		return false, result
	}
	return false, 0
}

// paintDarkMenuBar paints the menubar with dark background
func paintDarkMenuBar(hwnd HWND, theme *MenuBarTheme) {
	// Get menubar info
	var menuBarInfo MENUBARINFO
	menuBarInfo.CbSize = uint32(unsafe.Sizeof(menuBarInfo))
	if !GetMenuBarInfo(hwnd, OBJID_MENU, 0, &menuBarInfo) {
		return
	}

	// Get window DC
	hdc := GetWindowDC(hwnd)
	if hdc == 0 {
		return
	}
	defer ReleaseDC(hwnd, hdc)

	// Paint the menubar background with dark color
	FillRect(hdc, &menuBarInfo.Bar, theme.menuBarBackgroundBrush)

	// Get window and client rects to find the non-client area
	windowRect := GetWindowRect(hwnd)
	clientRect := GetClientRect(hwnd)

	// Convert client rect top-left to screen coordinates
	_, screenY := ClientToScreen(hwnd, int(clientRect.Left), int(clientRect.Top))

	// Paint the entire area between menubar and client area
	// This should cover any borders
	borderRect := RECT{
		Left:   0,
		Top:    menuBarInfo.Bar.Bottom - windowRect.Top,
		Right:  windowRect.Right - windowRect.Left,
		Bottom: int32(screenY) - windowRect.Top,
	}
	FillRect(hdc, &borderRect, theme.menuBarBackgroundBrush)
}

func drawMenuBarText(hwnd HWND, hdc HDC, menuBarInfo *MENUBARINFO, theme *MenuBarTheme) {
	// Get the menu handle
	hmenu := menuBarInfo.Menu
	if hmenu == 0 {
		return
	}

	// Get the number of menu items
	itemCount := GetMenuItemCount(hmenu)
	if itemCount <= 0 {
		return
	}

	// Set text color and background mode
	SetTextColor(hdc, COLORREF(*theme.TitleBarText))
	SetBkMode(hdc, TRANSPARENT)

	// Get the window rect for coordinate conversion
	winRect := GetWindowRect(hwnd)

	// Iterate through each menu item
	for i := 0; i < itemCount; i++ {
		// Get the menu item rect
		var itemRect RECT
		if !GetMenuItemRect(hwnd, hmenu, uint32(i), &itemRect) {
			continue
		}

		// Convert to window coordinates
		OffsetRect(&itemRect, int(-winRect.Left), int(-winRect.Top))

		// Get the menu item text
		menuString := make([]uint16, 256)
		mii := MENUITEMINFO{
			CbSize:     uint32(unsafe.Sizeof(MENUITEMINFO{})),
			FMask:      MIIM_STRING,
			DwTypeData: &menuString[0],
			Cch:        uint32(len(menuString) - 1),
		}

		if GetMenuItemInfo(hmenu, uint32(i), true, &mii) {
			// Draw the text
			DrawText(hdc, menuString, -1, &itemRect, DT_CENTER|DT_SINGLELINE|DT_VCENTER)
		}
	}
}

func UAHDrawMenuItem(hwnd HWND, wParam uintptr, lParam uintptr) uintptr {
	return 1
	//pUDMI := (*UAHDRAWMENUITEM)(unsafe.Pointer(lParam))
	//
	//var pbrBackground, pbrBorder *HBRUSH
	//pbrBackground = &brItemBackground
	//pbrBorder = &brItemBackground
	//
	//// get the menu item string
	//menuString := make([]uint16, 256)
	//mii := MENUITEMINFO{
	//	CbSize:     uint32(unsafe.Sizeof(mii)),
	//	FMask:      MIIM_STRING,
	//	DwTypeData: uintptr(unsafe.Pointer(&menuString[0])),
	//	Cch:        uint32(len(menuString) - 1),
	//}
	//GetMenuItemInfo(pUDMI.um.hmenu, pUDMI.umi.iPosition, true, &mii)
	//
	//// get the item state for drawing
	//dwFlags := DT_CENTER | DT_SINGLELINE | DT_VCENTER
	//
	//iTextStateID := 0
	//iBackgroundStateID := 0
	//switch {
	//case pUDMI.dis.itemState&ODS_INACTIVE != 0 || pUDMI.dis.itemState&ODS_DEFAULT != 0:
	//	// normal display
	//	iTextStateID = MBI_NORMAL
	//	iBackgroundStateID = MBI_NORMAL
	//case pUDMI.dis.itemState&ODS_HOTLIGHT != 0:
	//	// hot tracking
	//	iTextStateID = MBI_HOT
	//	iBackgroundStateID = MBI_HOT
	//
	//	pbrBackground = &brItemBackgroundHot
	//	pbrBorder = &brItemBorder
	//case pUDMI.dis.itemState&ODS_SELECTED != 0:
	//	// clicked
	//	iTextStateID = MBI_PUSHED
	//	iBackgroundStateID = MBI_PUSHED
	//
	//	pbrBackground = &brItemBackgroundSelected
	//	pbrBorder = &brItemBorder
	//case pUDMI.dis.itemState&ODS_GRAYED != 0 || pUDMI.dis.itemState&ODS_DISABLED != 0:
	//	// disabled / grey text
	//	iTextStateID = MBI_DISABLED
	//	iBackgroundStateID = MBI_DISABLED
	//case pUDMI.dis.itemState&ODS_NOACCEL != 0:
	//	// hide prefix
	//	dwFlags |= DT_HIDEPREFIX
	//}
	//
	//if g_menuTheme == 0 {
	//	g_menuTheme = OpenThemeData(hwnd, "Menu")
	//}
	//
	//opts := DTTOPTS{
	//	DtSize:     uint32(unsafe.Sizeof(opts)),
	//	DwFlags:    DTT_TEXTCOLOR,
	//	CrText:     RGB(0x00, 0x00, 0x20),
	//	ITextState: iTextStateID,
	//}
	//if iTextStateID == MBI_DISABLED {
	//	opts.CrText = RGB(0x40, 0x40, 0x40)
	//}
	//
	//FillRect(pUDMI.um.hdc, &pUDMI.dis.rcItem, *pbrBackground)
	//FrameRect(pUDMI.um.hdc, &pUDMI.dis.rcItem, *pbrBorder)
	//DrawThemeTextEx(g_menuTheme, pUDMI.um.hdc, MENU_BARITEM, MBI_NORMAL, uintptr(unsafe.Pointer(&menuString[0])), mii.cch, dwFlags, &pUDMI.dis.rcItem, &opts)

	//return 1
}
