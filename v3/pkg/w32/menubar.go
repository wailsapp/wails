package w32

import "unsafe"

const (
	OBJID_MENU = 3
)

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

var darkModeTitleBarBrush HBRUSH

func init() {
	darkModeTitleBarBrush, _, _ = procCreateSolidBrush.Call(
		uintptr(0x8080FF),
	)
}

func CreateSolidBrush(color COLORREF) HBRUSH {
	ret, _, _ := procCreateSolidBrush.Call(
		uintptr(color),
	)
	return HBRUSH(ret)
}

func UAHDrawMenu(hwnd HWND, wParam uintptr, lParam uintptr) uintptr {
	if !IsCurrentlyDarkMode() {
		return 0
	}
	udm := (*UAHMENU)(unsafe.Pointer(lParam))
	var rc RECT

	// get the menubar rect
	{
		menuBarInfo, err := GetMenuBarInfo(hwnd, OBJID_MENU, 0)
		if err != nil {
			return 0
		}

		winRect := GetWindowRect(hwnd)

		// the rcBar is offset by the window rect
		rc = menuBarInfo.Bar
		OffsetRect(&rc, int(-winRect.Left), int(-winRect.Top))
	}

	FillRect(udm.Hdc, &rc, darkModeTitleBarBrush)

	return 1
}

func UAHDrawMenuItem(hwnd HWND, wParam uintptr, lParam uintptr) uintptr {
	pUDMI := (*UAHDRAWMENUITEM)(unsafe.Pointer(lParam))

	var pbrBackground, pbrBorder *HBRUSH
	pbrBackground = &brItemBackground
	pbrBorder = &brItemBackground

	// get the menu item string
	menuString := make([]uint16, 256)
	mii := MENUITEMINFO{
		CbSize:     uint32(unsafe.Sizeof(mii)),
		FMask:      MIIM_STRING,
		DwTypeData: uintptr(unsafe.Pointer(&menuString[0])),
		Cch:        uint32(len(menuString) - 1),
	}
	GetMenuItemInfo(pUDMI.um.hmenu, pUDMI.umi.iPosition, true, &mii)

	// get the item state for drawing
	dwFlags := DT_CENTER | DT_SINGLELINE | DT_VCENTER

	iTextStateID := 0
	iBackgroundStateID := 0
	switch {
	case pUDMI.dis.itemState&ODS_INACTIVE != 0 || pUDMI.dis.itemState&ODS_DEFAULT != 0:
		// normal display
		iTextStateID = MBI_NORMAL
		iBackgroundStateID = MBI_NORMAL
	case pUDMI.dis.itemState&ODS_HOTLIGHT != 0:
		// hot tracking
		iTextStateID = MBI_HOT
		iBackgroundStateID = MBI_HOT

		pbrBackground = &brItemBackgroundHot
		pbrBorder = &brItemBorder
	case pUDMI.dis.itemState&ODS_SELECTED != 0:
		// clicked
		iTextStateID = MBI_PUSHED
		iBackgroundStateID = MBI_PUSHED

		pbrBackground = &brItemBackgroundSelected
		pbrBorder = &brItemBorder
	case pUDMI.dis.itemState&ODS_GRAYED != 0 || pUDMI.dis.itemState&ODS_DISABLED != 0:
		// disabled / grey text
		iTextStateID = MBI_DISABLED
		iBackgroundStateID = MBI_DISABLED
	case pUDMI.dis.itemState&ODS_NOACCEL != 0:
		// hide prefix
		dwFlags |= DT_HIDEPREFIX
	}

	if g_menuTheme == 0 {
		g_menuTheme = OpenThemeData(hwnd, "Menu")
	}

	opts := DTTOPTS{
		DtSize:     uint32(unsafe.Sizeof(opts)),
		DwFlags:    DTT_TEXTCOLOR,
		CrText:     RGB(0x00, 0x00, 0x20),
		ITextState: iTextStateID,
	}
	if iTextStateID == MBI_DISABLED {
		opts.CrText = RGB(0x40, 0x40, 0x40)
	}

	FillRect(pUDMI.um.hdc, &pUDMI.dis.rcItem, *pbrBackground)
	FrameRect(pUDMI.um.hdc, &pUDMI.dis.rcItem, *pbrBorder)
	DrawThemeTextEx(g_menuTheme, pUDMI.um.hdc, MENU_BARITEM, MBI_NORMAL, uintptr(unsafe.Pointer(&menuString[0])), mii.cch, dwFlags, &pUDMI.dis.rcItem, &opts)

	return 1
}

func GetMenuBarInfo(hwnd HWND, idObject uint32, idItem uint32) (*MENUBARINFO, error) {
	var mi MENUBARINFO

	mi.CbSize = uint32(unsafe.Sizeof(&mi))
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
