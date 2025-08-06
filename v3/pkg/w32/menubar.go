//go:build windows

package w32

import (
	"os"
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
	MENU_BARITEM   = 8 // Menu bar item part ID for theme drawing
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

// Track hover state for menubar items when maximized
var (
	currentHoverItem int  = -1
	menuIsOpen       bool = false // Track if a dropdown menu is open
)

func MenuBarWndProc(hwnd HWND, msg uint32, wParam WPARAM, lParam LPARAM, theme *MenuBarTheme) (bool, LRESULT) {
	// Only proceed if we have a theme (either for dark or light mode)
	if theme == nil {
		return false, 0
	}
	switch msg {
	case WM_UAHDRAWMENU:
		udm := (*UAHMENU)(unsafe.Pointer(lParam))

		// Check if maximized first
		isMaximized := IsZoomed(hwnd)

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

		// DEBUG: Log the coordinates
		// println("WM_UAHDRAWMENU: maximized=", isMaximized)
		// println("  menubar screen rect: L=", menuBarInfo.Bar.Left, "T=", menuBarInfo.Bar.Top,
		// 	"R=", menuBarInfo.Bar.Right, "B=", menuBarInfo.Bar.Bottom)
		// println("  window rect: L=", winRect.Left, "T=", winRect.Top,
		// 	"R=", winRect.Right, "B=", winRect.Bottom)
		// println("  converted rect: L=", rc.Left, "T=", rc.Top,
		// 	"R=", rc.Right, "B=", rc.Bottom)

		// When maximized, Windows extends the window beyond the visible area
		// We need to adjust the menubar rect to ensure it's fully visible
		if isMaximized {
			// Get the frame size - this is how much the window extends beyond visible area when maximized
			frameY := GetSystemMetrics(SM_CYSIZEFRAME)
			paddedBorder := GetSystemMetrics(SM_CXPADDEDBORDER)

			// In Windows 10/11, the actual border is frame + padding
			borderSize := frameY + paddedBorder

			// println("  Frame metrics: frameY=", frameY, "paddedBorder=", paddedBorder, "borderSize=", borderSize)

			// First, fill the area from the top of the visible area to the menubar
			topFillRect := RECT{
				Left:   rc.Left,
				Top:    int32(borderSize), // Start of visible area in window coordinates
				Right:  rc.Right,
				Bottom: rc.Top, // Up to where the menubar starts
			}
			FillRect(udm.Hdc, &topFillRect, theme.menuBarBackgroundBrush)
		}

		// Fill the entire menubar background with dark color
		FillRect(udm.Hdc, &rc, theme.menuBarBackgroundBrush)

		// Paint over the menubar border explicitly
		// The border is typically 1-2 pixels at the bottom
		borderRect := rc
		borderRect.Top = borderRect.Bottom - 1
		borderRect.Bottom = borderRect.Bottom + 2
		FillRect(udm.Hdc, &borderRect, theme.menuBarBackgroundBrush)

		// When maximized, we still need to handle the drawing ourselves
		// Some projects found that returning false here causes issues

		// When maximized, manually draw all menu items here
		if isMaximized {
			// Draw each menu item manually
			itemCount := GetMenuItemCount(menuBarInfo.Menu)
			for i := 0; i < itemCount; i++ {
				var itemRect RECT
				if GetMenuItemRect(hwnd, menuBarInfo.Menu, uint32(i), &itemRect) {
					// Convert to window coordinates
					OffsetRect(&itemRect, int(-winRect.Left), int(-winRect.Top))

					// Check if this item is hovered
					if i == currentHoverItem {
						// Fill with hover background
						FillRect(udm.Hdc, &itemRect, theme.menuHoverBackgroundBrush)
					}

					// Get menu text
					menuString := make([]uint16, 256)
					mii := MENUITEMINFO{
						CbSize:     uint32(unsafe.Sizeof(MENUITEMINFO{})),
						FMask:      MIIM_STRING,
						DwTypeData: &menuString[0],
						Cch:        uint32(len(menuString) - 1),
					}

					if GetMenuItemInfo(menuBarInfo.Menu, uint32(i), true, &mii) {
						// Draw the text
						if i == currentHoverItem {
							SetTextColor(udm.Hdc, COLORREF(*theme.MenuHoverText))
						} else {
							SetTextColor(udm.Hdc, COLORREF(*theme.TitleBarText))
						}
						SetBkMode(udm.Hdc, TRANSPARENT)
						DrawText(udm.Hdc, menuString, -1, &itemRect, DT_CENTER|DT_SINGLELINE|DT_VCENTER)
					}
				}
			}
		}

		// Return the original HDC so Windows can draw the menu text
		return true, LRESULT(udm.Hdc)
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

		// Check if we're getting menu item draw messages when maximized or fullscreen
		isMaximized := IsZoomed(hwnd)

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

		// When maximized/fullscreen, try without VCENTER to see if text appears
		if isMaximized && os.Getenv("WAILS_TEST_NO_VCENTER") == "1" {
			dwFlags = uint32(DT_CENTER | DT_SINGLELINE)
			println("  Using dwFlags without VCENTER")
		}

		// Check if this is a menubar item
		// When dwFlags has 0x0A00 (2560) it's a menubar item
		isMenuBarItem := (udmi.UM.DwFlags&0x0A00) == 0x0A00 || udmi.UM.DwFlags == 0

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
		if bgBrush != 0 {
			FillRect(udmi.UM.Hdc, &udmi.DIS.RcItem, bgBrush)
		}

		// Draw text
		SetTextColor(udmi.UM.Hdc, COLORREF(textColor))
		SetBkMode(udmi.UM.Hdc, TRANSPARENT)

		// When maximized/fullscreen and menubar item, use the same font settings as drawMenuBarText
		if isMaximized && isMenuBarItem {
			// Create a non-bold font explicitly
			menuFont := LOGFONT{
				Height:         -12, // Standard Windows menu font height (9pt)
				Weight:         400, // FW_NORMAL (not bold)
				CharSet:        1,   // DEFAULT_CHARSET
				Quality:        5,   // CLEARTYPE_QUALITY
				PitchAndFamily: 0,   // DEFAULT_PITCH
			}
			// Set font face name to "Segoe UI" (Windows default)
			fontName := []uint16{'S', 'e', 'g', 'o', 'e', ' ', 'U', 'I', 0}
			copy(menuFont.FaceName[:], fontName)

			hFont := CreateFontIndirect(&menuFont)
			if hFont != 0 {
				oldFont := SelectObject(udmi.UM.Hdc, HGDIOBJ(hFont))
				DrawText(udmi.UM.Hdc, menuString, -1, &udmi.DIS.RcItem, dwFlags)
				SelectObject(udmi.UM.Hdc, oldFont)
				DeleteObject(HGDIOBJ(hFont))
			} else {
				DrawText(udmi.UM.Hdc, menuString, -1, &udmi.DIS.RcItem, dwFlags)
			}
			return true, 4 // CDRF_SKIPDEFAULT
		} else {
			DrawText(udmi.UM.Hdc, menuString, -1, &udmi.DIS.RcItem, dwFlags)
		}

		// Return appropriate value based on whether we're in maximized/fullscreen
		// For maximized, we need to ensure Windows doesn't override our drawing
		if isMaximized {
			// Skip default processing to prevent Windows from overriding our colors
			return true, 4 // CDRF_SKIPDEFAULT
		}
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

		// CRITICAL: Force complete menubar redraw when maximized
		if msg == WM_SIZE && wParam == SIZE_MAXIMIZED {
			// Invalidate the entire menubar area to force redraw
			var mbi MENUBARINFO
			mbi.CbSize = uint32(unsafe.Sizeof(mbi))
			if GetMenuBarInfo(hwnd, OBJID_MENU, 0, &mbi) {
				InvalidateRect(hwnd, &mbi.Bar, true)
				DrawMenuBar(hwnd)
			}
		}

		return false, result
	case WM_SETFOCUS, WM_KILLFOCUS:
		// Handle focus changes (e.g., when inspector opens)
		result := DefWindowProc(hwnd, msg, wParam, lParam)

		// Repaint menubar after focus change
		paintDarkMenuBar(hwnd, theme)

		return false, result
	case WM_ERASEBKGND:
		// When maximized, draw menubar text here
		if IsZoomed(hwnd) {
			var menuBarInfo MENUBARINFO
			menuBarInfo.CbSize = uint32(unsafe.Sizeof(menuBarInfo))
			if GetMenuBarInfo(hwnd, OBJID_MENU, 0, &menuBarInfo) {
				hdc := HDC(wParam)
				drawMenuBarText(hwnd, hdc, &menuBarInfo, theme)
			}
		}
		return false, 0
	case WM_NCMOUSEMOVE, WM_MOUSEMOVE:
		// Track mouse movement for hover effects when maximized
		if IsZoomed(hwnd) {
			// Don't process hover changes while menu is open
			if menuIsOpen {
				return false, 0
			}

			var screenX, screenY int32
			if msg == WM_NCMOUSEMOVE {
				// For NC messages, lParam contains screen coordinates
				screenX = int32(LOWORD(uint32(lParam)))
				screenY = int32(HIWORD(uint32(lParam)))
			} else {
				// For regular MOUSEMOVE, convert client to screen coordinates
				clientX := int32(LOWORD(uint32(lParam)))
				clientY := int32(HIWORD(uint32(lParam)))
				sx, sy := ClientToScreen(hwnd, int(clientX), int(clientY))
				screenX = int32(sx)
				screenY = int32(sy)
			}

			// Check if we're over the menubar
			var menuBarInfo MENUBARINFO
			menuBarInfo.CbSize = uint32(unsafe.Sizeof(menuBarInfo))
			if GetMenuBarInfo(hwnd, OBJID_MENU, 0, &menuBarInfo) {
				// menuBarInfo.Bar already contains screen coordinates
				// Check if mouse is over menubar using screen coordinates
				if screenX >= menuBarInfo.Bar.Left && screenX <= menuBarInfo.Bar.Right &&
					screenY >= menuBarInfo.Bar.Top && screenY <= menuBarInfo.Bar.Bottom {

					// Always re-request mouse tracking to ensure we get leave messages
					TrackMouseEvent(&TRACKMOUSEEVENT{
						CbSize:      uint32(unsafe.Sizeof(TRACKMOUSEEVENT{})),
						DwFlags:     TME_LEAVE | TME_NONCLIENT,
						HwndTrack:   hwnd,
						DwHoverTime: 0,
					})
					// Find which menu item we're over
					itemCount := GetMenuItemCount(menuBarInfo.Menu)
					newHoverItem := -1

					for i := 0; i < itemCount; i++ {
						var itemRect RECT
						if GetMenuItemRect(hwnd, menuBarInfo.Menu, uint32(i), &itemRect) {
							// itemRect is already in screen coordinates from GetMenuItemRect
							// Check using screen coordinates
							if screenX >= itemRect.Left && screenX <= itemRect.Right &&
								screenY >= itemRect.Top && screenY <= itemRect.Bottom {
								newHoverItem = i
								break
							}
						}
					}

					// If hover item changed, update and redraw just the menubar
					if newHoverItem != currentHoverItem {
						currentHoverItem = newHoverItem
						// Get the actual menubar rect for precise invalidation
						winRect := GetWindowRect(hwnd)
						menubarRect := menuBarInfo.Bar
						// Convert to window coordinates
						menubarRect.Left -= winRect.Left
						menubarRect.Top -= winRect.Top
						menubarRect.Right -= winRect.Left
						menubarRect.Bottom -= winRect.Top
						// Invalidate only the menubar
						InvalidateRect(hwnd, &menubarRect, false)
					}
				} else {
					// Mouse left menubar
					if currentHoverItem != -1 {
						currentHoverItem = -1
						// Get the actual menubar rect
						winRect := GetWindowRect(hwnd)
						menubarRect := menuBarInfo.Bar
						// Convert to window coordinates
						menubarRect.Left -= winRect.Left
						menubarRect.Top -= winRect.Top
						menubarRect.Right -= winRect.Left
						menubarRect.Bottom -= winRect.Top
						InvalidateRect(hwnd, &menubarRect, false)
					}
				}
			}
		}
		return false, 0
	case WM_NCLBUTTONDOWN:
		// When clicking on menubar, clear hover state immediately
		if IsZoomed(hwnd) && currentHoverItem != -1 {
			// Check if click is on menubar
			var menuBarInfo MENUBARINFO
			menuBarInfo.CbSize = uint32(unsafe.Sizeof(menuBarInfo))
			if GetMenuBarInfo(hwnd, OBJID_MENU, 0, &menuBarInfo) {
				// Get click position (screen coordinates)
				clickX := int32(LOWORD(uint32(lParam)))
				clickY := int32(HIWORD(uint32(lParam)))

				if clickX >= menuBarInfo.Bar.Left && clickX <= menuBarInfo.Bar.Right &&
					clickY >= menuBarInfo.Bar.Top && clickY <= menuBarInfo.Bar.Bottom {
					// Click is on menubar - clear hover
					currentHoverItem = -1
				}
			}
		}
		return false, 0
	case WM_NCMOUSELEAVE, WM_MOUSELEAVE:
		// Clear hover state when mouse leaves (but not if menu is open)
		if IsZoomed(hwnd) && currentHoverItem != -1 && !menuIsOpen {
			currentHoverItem = -1
			// Get menubar info for precise invalidation
			var menuBarInfo MENUBARINFO
			menuBarInfo.CbSize = uint32(unsafe.Sizeof(menuBarInfo))
			if GetMenuBarInfo(hwnd, OBJID_MENU, 0, &menuBarInfo) {
				winRect := GetWindowRect(hwnd)
				menubarRect := menuBarInfo.Bar
				menubarRect.Left -= winRect.Left
				menubarRect.Top -= winRect.Top
				menubarRect.Right -= winRect.Left
				menubarRect.Bottom -= winRect.Top
				InvalidateRect(hwnd, &menubarRect, false)
			}
		}
		return false, 0
	case WM_ENTERMENULOOP:
		// Menu is being opened - clear hover state
		menuIsOpen = true
		if IsZoomed(hwnd) && currentHoverItem != -1 {
			oldHoverItem := currentHoverItem
			currentHoverItem = -1
			// Redraw the previously hovered item to remove hover effect
			var menuBarInfo MENUBARINFO
			menuBarInfo.CbSize = uint32(unsafe.Sizeof(menuBarInfo))
			if GetMenuBarInfo(hwnd, OBJID_MENU, 0, &menuBarInfo) {
				var itemRect RECT
				if GetMenuItemRect(hwnd, menuBarInfo.Menu, uint32(oldHoverItem), &itemRect) {
					winRect := GetWindowRect(hwnd)
					// Convert to window coordinates
					itemRect.Left -= winRect.Left
					itemRect.Top -= winRect.Top
					itemRect.Right -= winRect.Left
					itemRect.Bottom -= winRect.Top
					// Add some padding
					itemRect.Left -= 5
					itemRect.Right += 5
					itemRect.Top -= 5
					itemRect.Bottom += 5
					InvalidateRect(hwnd, &itemRect, false)
				}
			}
		}
		return false, 0
	case WM_EXITMENULOOP:
		// Menu has been closed
		menuIsOpen = false
		// Clear any existing hover state first
		currentHoverItem = -1
		// Force a complete menubar redraw
		if IsZoomed(hwnd) {
			var menuBarInfo MENUBARINFO
			menuBarInfo.CbSize = uint32(unsafe.Sizeof(menuBarInfo))
			if GetMenuBarInfo(hwnd, OBJID_MENU, 0, &menuBarInfo) {
				winRect := GetWindowRect(hwnd)
				menubarRect := menuBarInfo.Bar
				menubarRect.Left -= winRect.Left
				menubarRect.Top -= winRect.Top
				menubarRect.Right -= winRect.Left
				menubarRect.Bottom -= winRect.Top
				InvalidateRect(hwnd, &menubarRect, false)
			}
			// Force a timer to restart mouse tracking
			SetTimer(hwnd, 1001, 50, 0)
		}
		return false, 0
	case WM_TIMER:
		// Handle our mouse tracking restart timer
		if wParam == 1001 {
			KillTimer(hwnd, 1001)
			if IsZoomed(hwnd) {
				// Get current mouse position and simulate a mouse move
				x, y, _ := GetCursorPos()
				// Check if mouse is over the window
				winRect := GetWindowRect(hwnd)
				if x >= int(winRect.Left) && x <= int(winRect.Right) &&
					y >= int(winRect.Top) && y <= int(winRect.Bottom) {
					// Check if we're over the menubar specifically
					var menuBarInfo MENUBARINFO
					menuBarInfo.CbSize = uint32(unsafe.Sizeof(menuBarInfo))
					if GetMenuBarInfo(hwnd, OBJID_MENU, 0, &menuBarInfo) {
						if int32(x) >= menuBarInfo.Bar.Left && int32(x) <= menuBarInfo.Bar.Right &&
							int32(y) >= menuBarInfo.Bar.Top && int32(y) <= menuBarInfo.Bar.Bottom {
							// Post a non-client mouse move to restart tracking
							PostMessage(hwnd, WM_NCMOUSEMOVE, 0, uintptr(y)<<16|uintptr(x)&0xFFFF)
						} else {
							// Convert to client coordinates for regular mouse move
							clientX, clientY, _ := ScreenToClient(hwnd, x, y)
							// Post a mouse move message to restart tracking
							PostMessage(hwnd, WM_MOUSEMOVE, 0, uintptr(clientY)<<16|uintptr(clientX)&0xFFFF)
						}
					}
				}
			}
			return true, 0
		}
		return false, 0
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

	// Check if window is maximized or fullscreen
	isMaximized := IsZoomed(hwnd)
	isFullscreen := false

	// Check if window is in fullscreen by checking if it covers the monitor
	windowRect := GetWindowRect(hwnd)
	monitor := MonitorFromWindow(hwnd, MONITOR_DEFAULTTOPRIMARY)
	var monitorInfo MONITORINFO
	monitorInfo.CbSize = uint32(unsafe.Sizeof(monitorInfo))
	if GetMonitorInfo(monitor, &monitorInfo) {
		// If window matches monitor bounds, it's fullscreen
		if windowRect.Left == monitorInfo.RcMonitor.Left &&
			windowRect.Top == monitorInfo.RcMonitor.Top &&
			windowRect.Right == monitorInfo.RcMonitor.Right &&
			windowRect.Bottom == monitorInfo.RcMonitor.Bottom {
			isFullscreen = true
		}
	}

	// When maximized or fullscreen, we need to handle the special case
	if isMaximized || isFullscreen {
		// Convert menubar rect from screen to window coordinates
		menubarRect := menuBarInfo.Bar
		menubarRect.Left -= windowRect.Left
		menubarRect.Top -= windowRect.Top
		menubarRect.Right -= windowRect.Left
		menubarRect.Bottom -= windowRect.Top

		if isMaximized && !isFullscreen {
			// Get the frame size (only for maximized, not fullscreen)
			frameY := GetSystemMetrics(SM_CYSIZEFRAME)
			paddedBorder := GetSystemMetrics(SM_CXPADDEDBORDER)
			borderSize := frameY + paddedBorder

			// Fill from visible area top to menubar
			topFillRect := RECT{
				Left:   menubarRect.Left,
				Top:    int32(borderSize), // Start of visible area
				Right:  menubarRect.Right,
				Bottom: menubarRect.Top,
			}
			FillRect(hdc, &topFillRect, theme.menuBarBackgroundBrush)
		} else if isFullscreen {
			// In fullscreen, fill from the very top
			topFillRect := RECT{
				Left:   menubarRect.Left,
				Top:    0, // Start from top in fullscreen
				Right:  menubarRect.Right,
				Bottom: menubarRect.Top,
			}
			FillRect(hdc, &topFillRect, theme.menuBarBackgroundBrush)
		}

		// Fill the menubar itself
		FillRect(hdc, &menubarRect, theme.menuBarBackgroundBrush)
	} else {
		// Paint the menubar background with dark color
		FillRect(hdc, &menuBarInfo.Bar, theme.menuBarBackgroundBrush)
	}

	// Get window and client rects to find the non-client area
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

	// When maximized or fullscreen, also draw menubar text
	if isMaximized || isFullscreen {
		drawMenuBarText(hwnd, hdc, &menuBarInfo, theme)
	}
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

	// Create a non-bold font explicitly
	menuFont := LOGFONT{
		Height:         -12, // Standard Windows menu font height (9pt)
		Weight:         400, // FW_NORMAL (not bold)
		CharSet:        1,   // DEFAULT_CHARSET
		Quality:        5,   // CLEARTYPE_QUALITY
		PitchAndFamily: 0,   // DEFAULT_PITCH
	}
	// Set font face name to "Segoe UI" (Windows default)
	fontName := []uint16{'S', 'e', 'g', 'o', 'e', ' ', 'U', 'I', 0}
	copy(menuFont.FaceName[:], fontName)

	hFont := CreateFontIndirect(&menuFont)
	if hFont != 0 {
		oldFont := SelectObject(hdc, HGDIOBJ(hFont))
		defer func() {
			SelectObject(hdc, oldFont)
			DeleteObject(HGDIOBJ(hFont))
		}()
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

		// Check if this item is hovered
		if i == currentHoverItem {
			// Fill with hover background
			FillRect(hdc, &itemRect, theme.menuHoverBackgroundBrush)
		}

		// Get the menu item text
		menuString := make([]uint16, 256)
		mii := MENUITEMINFO{
			CbSize:     uint32(unsafe.Sizeof(MENUITEMINFO{})),
			FMask:      MIIM_STRING,
			DwTypeData: &menuString[0],
			Cch:        uint32(len(menuString) - 1),
		}

		if GetMenuItemInfo(hmenu, uint32(i), true, &mii) {
			// Set text color based on hover state
			if i == currentHoverItem {
				SetTextColor(hdc, COLORREF(*theme.MenuHoverText))
			} else {
				SetTextColor(hdc, COLORREF(*theme.TitleBarText))
			}
			// Draw the text
			DrawText(hdc, menuString, -1, &itemRect, DT_CENTER|DT_SINGLELINE|DT_VCENTER)
		}
	}
}
