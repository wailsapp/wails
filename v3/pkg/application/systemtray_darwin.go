//go:build darwin

package application

/*
#cgo CFLAGS: -mmacosx-version-min=10.13 -x objective-c
#cgo LDFLAGS: -framework Cocoa -framework WebKit

#include "Cocoa/Cocoa.h"
#include "menuitem_darwin.h"
#include "systemtray_darwin.h"

// Show the system tray icon
static void systemTrayShow(void* nsStatusItem) {
    dispatch_async(dispatch_get_main_queue(), ^{
		// Get the NSStatusItem
        NSStatusItem *statusItem = (NSStatusItem *)nsStatusItem;
        [statusItem setVisible:YES];
    });
}

// Hide the system tray icon
static void systemTrayHide(void* nsStatusItem) {
    dispatch_async(dispatch_get_main_queue(), ^{
        NSStatusItem *statusItem = (NSStatusItem *)nsStatusItem;
        [statusItem setVisible:NO];
    });
}

*/
import "C"
import (
	"errors"
	"unsafe"

	"github.com/leaanthony/go-ansi-parser"
)

type macosSystemTray struct {
	id    uint
	label string
	icon  []byte
	menu  *Menu

	nsStatusItem      unsafe.Pointer
	nsImage           unsafe.Pointer
	nsMenu            unsafe.Pointer
	iconPosition      IconPosition
	isTemplateIcon    bool
	parent            *SystemTray
	lastClickedScreen unsafe.Pointer
}

func (s *macosSystemTray) Show() {
	if s.nsStatusItem == nil {
		return
	}
	C.systemTrayShow(s.nsStatusItem)
}

func (s *macosSystemTray) Hide() {
	if s.nsStatusItem == nil {
		return
	}
	C.systemTrayHide(s.nsStatusItem)
}

func (s *macosSystemTray) openMenu() {
	if s.nsMenu == nil {
		return
	}
	C.showMenu(s.nsStatusItem, s.nsMenu)
}

type button int

const (
	leftButtonDown  button = 1
	rightButtonDown button = 3
)

// system tray map
var systemTrayMap = make(map[uint]*macosSystemTray)

//export systrayClickCallback
func systrayClickCallback(id C.long, buttonID C.int) {
	// Get the system tray
	systemTray := systemTrayMap[uint(id)]
	if systemTray == nil {
		globalApplication.error("system tray not found: %v", id)
		return
	}
	systemTray.processClick(button(buttonID))
}

func (s *macosSystemTray) setIconPosition(position IconPosition) {
	s.iconPosition = position
}

func (s *macosSystemTray) setMenu(menu *Menu) {
	s.menu = menu
}

func (s *macosSystemTray) positionWindow(window *WebviewWindow, offset int) error {
	// Get the window's native window
	impl := window.impl.(*macosWebviewWindow)

	// Position the window relative to the systray
	C.systemTrayPositionWindow(s.nsStatusItem, impl.nsWindow, C.int(offset))

	return nil
}

func (s *macosSystemTray) getScreen() (*Screen, error) {
	if s.lastClickedScreen != nil {
		// Get the screen frame
		frame := C.NSScreen_frame(s.lastClickedScreen)
		result := &Screen{
			Bounds: Rect{
				X:      int(frame.origin.x),
				Y:      int(frame.origin.y),
				Width:  int(frame.size.width),
				Height: int(frame.size.height),
			},
		}
		return result, nil
	}
	return nil, errors.New("no screen available")
}

func (s *macosSystemTray) bounds() (*Rect, error) {
	var rect C.NSRect
	var screen unsafe.Pointer
	C.systemTrayGetBounds(s.nsStatusItem, &rect, &screen)

	// Store the screen for use in positionWindow
	s.lastClickedScreen = screen

	// Return the screen-relative coordinates
	result := &Rect{
		X:      int(rect.origin.x),
		Y:      int(rect.origin.y),
		Width:  int(rect.size.width),
		Height: int(rect.size.height),
	}
	return result, nil
}

func (s *macosSystemTray) run() {
	globalApplication.dispatchOnMainThread(func() {
		if s.nsStatusItem != nil {
			Fatal("System tray '%d' already running", s.id)
		}
		s.nsStatusItem = unsafe.Pointer(C.systemTrayNew(C.long(s.id)))

		if s.label != "" {
			s.setLabel(s.label)
		}
		if s.icon != nil {
			s.nsImage = unsafe.Pointer(C.imageFromBytes((*C.uchar)(&s.icon[0]), C.int(len(s.icon))))
			C.systemTraySetIcon(s.nsStatusItem, s.nsImage, C.int(s.iconPosition), C.bool(s.isTemplateIcon))
		}
		if s.menu != nil {
			s.menu.Update()
			// Convert impl to macosMenu object
			s.nsMenu = (s.menu.impl).(*macosMenu).nsMenu
		}
	})
}

func (s *macosSystemTray) setIcon(icon []byte) {
	s.icon = icon
	globalApplication.dispatchOnMainThread(func() {
		s.nsImage = unsafe.Pointer(C.imageFromBytes((*C.uchar)(&icon[0]), C.int(len(icon))))
		C.systemTraySetIcon(s.nsStatusItem, s.nsImage, C.int(s.iconPosition), C.bool(s.isTemplateIcon))
	})
}

func (s *macosSystemTray) setDarkModeIcon(icon []byte) {
	s.setIcon(icon)
}

func (s *macosSystemTray) setTemplateIcon(icon []byte) {
	s.icon = icon
	s.isTemplateIcon = true
	globalApplication.dispatchOnMainThread(func() {
		s.nsImage = unsafe.Pointer(C.imageFromBytes((*C.uchar)(&icon[0]), C.int(len(icon))))
		C.systemTraySetIcon(s.nsStatusItem, s.nsImage, C.int(s.iconPosition), C.bool(s.isTemplateIcon))
	})
}

func (s *macosSystemTray) setTooltip(tooltip string) {
	// Tooltips not supported on macOS
}

func newSystemTrayImpl(s *SystemTray) systemTrayImpl {
	result := &macosSystemTray{
		parent:         s,
		id:             s.id,
		label:          s.label,
		icon:           s.icon,
		menu:           s.menu,
		iconPosition:   s.iconPosition,
		isTemplateIcon: s.isTemplateIcon,
	}
	systemTrayMap[s.id] = result
	return result
}

func extractAnsiTextParts(text *ansi.StyledText) (label *C.char, fg *C.char, bg *C.char) {
	label = C.CString(text.Label)
	if text.FgCol != nil {
		fg = C.CString(text.FgCol.Hex)
	}
	if text.BgCol != nil {
		bg = C.CString(text.BgCol.Hex)
	}
	return
}

func (s *macosSystemTray) setLabel(label string) {
	s.label = label
	if !ansi.HasEscapeCodes(label) {
		C.systemTraySetLabel(s.nsStatusItem, C.CString(label))
	} else {
		parsed, err := ansi.Parse(label)
		if err != nil {
			C.systemTraySetLabel(s.nsStatusItem, C.CString(label))
			return
		}
		if len(parsed) == 0 {
			return
		}
		label, fg, bg := extractAnsiTextParts(parsed[0])
		var attributedString = C.createAttributedString(label, fg, bg)
		if len(parsed) > 1 {
			for _, parsedPart := range parsed[1:] {
				label, fg, bg = extractAnsiTextParts(parsedPart)
				attributedString = C.appendAttributedString(attributedString, label, fg, bg)
			}
		}

		C.systemTraySetANSILabel(s.nsStatusItem, attributedString)
	}
}

func (s *macosSystemTray) destroy() {
	// Remove the status item from the status bar and its associated menu
	C.systemTrayDestroy(s.nsStatusItem)
}

func (s *macosSystemTray) processClick(b button) {
	switch b {
	case leftButtonDown:
		// Check if we have a callback
		if s.parent.clickHandler != nil {
			s.parent.clickHandler()
			return
		}
		if s.parent.attachedWindow.Window != nil {
			s.parent.defaultClickHandler()
			return
		}
		if s.menu != nil {
			C.showMenu(s.nsStatusItem, s.nsMenu)
		}
	case rightButtonDown:
		// Check if we have a callback
		if s.parent.rightClickHandler != nil {
			s.parent.rightClickHandler()
			return
		}
		if s.menu != nil {
			if s.parent.attachedWindow.Window != nil {
				s.parent.attachedWindow.Window.Hide()
			}
			C.showMenu(s.nsStatusItem, s.nsMenu)
			return
		}
		if s.parent.attachedWindow.Window != nil {
			s.parent.defaultClickHandler()
		}
	}
}
