package w32

import (
	"fmt"
	"unsafe"
)

const (
	OBJID_MENU = -3
)

var (
	menuTheme HTHEME
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
	procGetMenuItemInfo = moduser32.NewProc("GetMenuItemInfoW")
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
	MenuHoverBackground    *uint32
	MenuHoverText          *uint32
	MenuSelectedBackground *uint32
	MenuSelectedText       *uint32

	// private brushes
	titleBarBackgroundBrush     HBRUSH
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
	d.TitleBarText = createColourWithDefaultColor(d.TitleBarText, RGB(255, 255, 255))
	d.MenuSelectedText = createColourWithDefaultColor(d.MenuSelectedText, RGB(255, 255, 255))
	d.MenuSelectedBackground = createColourWithDefaultColor(d.MenuSelectedBackground, RGB(60, 60, 60))
	d.MenuHoverText = createColourWithDefaultColor(d.MenuHoverText, RGB(255, 255, 255))
	d.MenuHoverBackground = createColourWithDefaultColor(d.MenuHoverBackground, RGB(45, 45, 45))
	// Create brushes
	d.titleBarBackgroundBrush = CreateSolidBrush(*d.TitleBarBackground)
	d.menuHoverBackgroundBrush = CreateSolidBrush(*d.MenuHoverBackground)
	d.menuSelectedBackgroundBrush = CreateSolidBrush(*d.MenuSelectedBackground)
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

func TrackPopupMenu(hmenu HMENU, flags uint32, x, y int32, reserved int32, hwnd HWND, prcRect *RECT) bool {
	ret, _, _ := procTrackPopupMenu.Call(
		uintptr(hmenu),
		uintptr(flags),
		uintptr(x),
		uintptr(y),
		uintptr(reserved),
		uintptr(hwnd),
		uintptr(unsafe.Pointer(prcRect)),
	)
	return ret != 0
}

func MenuBarWndProc(hwnd HWND, msg uint32, wParam WPARAM, lParam LPARAM, theme *MenuBarTheme) (bool, LRESULT) {
	if !IsCurrentlyDarkMode() {
		return false, 0
	}
	switch msg {
	case WM_UAHDRAWMENU:
		udm := (*UAHMENU)(unsafe.Pointer(lParam))
		var rc RECT

		// get the menubar rect
		menuBarInfo, err := GetMenuBarInfo(hwnd, OBJID_MENU, 0)
		if err != nil {
			return false, 0
		}

		winRect := GetWindowRect(hwnd)

		// the rcBar is offset by the window rect
		rc = menuBarInfo.Bar
		OffsetRect(&rc, int(-winRect.Left), int(-winRect.Top))

		FillRect(udm.Hdc, &rc, theme.titleBarBackgroundBrush)

		return true, 0
	case WM_UAHDRAWMENUITEM:
		udmi := (*UAHDRAWMENUITEM)(unsafe.Pointer(lParam))

		// Create buffer for menu text
		menuString := make([]uint16, 256)

		// Setup menu item info structure
		mii := MENUITEMINFO{
			CbSize:     uint32(unsafe.Sizeof(MENUITEMINFO{})),
			FMask:      MIIM_STRING,
			DwTypeData: &menuString[0],
			Cch:        uint32(len(menuString) - 1),
		}

		GetMenuItemInfo(udmi.UM.Hmenu, uint32(udmi.UAMI.Position), true, &mii)

		if udmi.DIS.ItemState&ODS_HOTLIGHT != 0 && mii.HSubMenu != 0 {
			// If this is a menu item with a submenu, and we're hovering,
			// tell the menu to track
			TrackPopupMenu(mii.HSubMenu,
				TPM_LEFTALIGN|TPM_TOPALIGN,
				int32(udmi.DIS.RcItem.Left),
				int32(udmi.DIS.RcItem.Bottom),
				0, hwnd, nil)
		}
		dwFlags := uint32(DT_CENTER | DT_SINGLELINE | DT_VCENTER)

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
			bgBrush = theme.titleBarBackgroundBrush
			textColor = *theme.TitleBarText
		}

		// Fill background
		FillRect(udmi.UM.Hdc, &udmi.DIS.RcItem, bgBrush)

		// Draw text
		SetTextColor(udmi.UM.Hdc, textColor)
		SetBkMode(udmi.UM.Hdc, TRANSPARENT)
		DrawText(udmi.UM.Hdc, menuString, -1, &udmi.DIS.RcItem, dwFlags)

		return true, 1
	case WM_UAHMEASUREMENUITEM:
		// Cast lParam to UAHMEASUREMENUITEM pointer
		mmi := (*UAHMEASUREMENUITEM)(unsafe.Pointer(lParam))

		// Let the default window procedure handle the basic measurement
		result := DefWindowProc(hwnd, msg, wParam, lParam)

		// Modify the width to be 1/3rd wider
		mmi.Mis.ItemWidth = (mmi.Mis.ItemWidth * 4) / 3

		return true, result
	case WM_NCPAINT, WM_NCACTIVATE:
		result := DefWindowProc(hwnd, msg, wParam, lParam)
		_, err := GetMenuBarInfo(hwnd, OBJID_MENU, 0)
		if err != nil {
			return false, 0
		}

		clientRect := GetClientRect(hwnd)
		points := []POINT{
			{
				X: clientRect.Left,
				Y: clientRect.Top,
			},
			{
				X: clientRect.Right,
				Y: clientRect.Bottom,
			},
		}
		MapWindowPoints(hwnd, 0, uintptr(unsafe.Pointer(&points[0])), 2)
		clientRect.Left = points[0].X
		clientRect.Top = points[0].Y
		clientRect.Right = points[1].X
		clientRect.Bottom = points[1].Y
		winRect := GetWindowRect(hwnd)

		OffsetRect(clientRect, int(-winRect.Left), int(-winRect.Top))

		line := *clientRect
		line.Bottom = line.Top
		line.Top = line.Top - 1

		hdc := GetWindowDC(hwnd)
		FillRect(hdc, &line, theme.titleBarBackgroundBrush)
		ReleaseDC(hwnd, hdc)
		return false, result
	}
	return false, 0
}

func MapWindowPoints(hWndFrom HWND, hWndTo HWND, points uintptr, numPoints uint32) {

	// Call the MapWindowPoints function
	ret, _, _ := procMapWindowPoints.Call(
		uintptr(hWndFrom),
		uintptr(hWndTo),
		points,
		uintptr(numPoints),
	)

	// Check for errors
	if ret == 0 {
		fmt.Println("MapWindowPoints failed")
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

func GetMenuBarInfo(hwnd HWND, idObject int32, idItem uint32) (*MENUBARINFO, error) {
	var mi MENUBARINFO

	mi.CbSize = uint32(unsafe.Sizeof(mi))
	ret, _, err := procGetMenuBarInfo.Call(
		hwnd,
		uintptr(idObject),
		uintptr(idItem),
		uintptr(unsafe.Pointer(&mi)))
	if ret == 0 {
		return nil, err
	}
	return &mi, nil
}
