//go:build darwin && !ios && !server

package application

/*
#cgo CFLAGS: -mmacosx-version-min=10.13 -x objective-c
#cgo LDFLAGS: -framework Cocoa -framework WebKit

#include <stdlib.h>
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
	"strings"
	"unsafe"
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

// systrayPreClickCallback is called from the NSEvent local monitor BEFORE the
// button processes the mouse-down.  It returns 1 when the framework should
// show the menu via native tracking (proper highlight, no app activation),
// or 0 to let the action handler fire for custom click/window behaviour.
//
//export systrayPreClickCallback
func systrayPreClickCallback(id C.long, buttonID C.int) C.int {
	systemTray := systemTrayMap[uint(id)]
	if systemTray == nil || systemTray.nsMenu == nil {
		return 0
	}
	b := button(buttonID)
	switch b {
	case leftButtonDown:
		if systemTray.parent.clickHandler == nil &&
			systemTray.parent.attachedWindow.Window == nil {
			return 1
		}
	case rightButtonDown:
		if systemTray.parent.rightClickHandler == nil {
			// Hide the attached window before the menu appears.
			if systemTray.parent.attachedWindow.Window != nil &&
				systemTray.parent.attachedWindow.Window.IsVisible() {
				systemTray.parent.attachedWindow.Window.Hide()
			}
			return 1
		}
	}
	return 0
}

func (s *macosSystemTray) setIconPosition(position IconPosition) {
	s.iconPosition = position
}

func (s *macosSystemTray) setMenu(menu *Menu) {
	s.menu = menu
	if s.nsStatusItem != nil && menu != nil {
		menu.Update()
		s.nsMenu = (menu.impl).(*macosMenu).nsMenu
		C.systemTraySetCachedMenu(s.nsStatusItem, s.nsMenu)
	}
}

func (s *macosSystemTray) positionWindow(window Window, offset int) error {
	// Get the platform-specific window implementation
	nativeWindow := window.NativeWindow()
	if nativeWindow == nil {
		return errors.New("window native implementation unavailable")
	}

	// Position the window relative to the systray
	C.systemTrayPositionWindow(s.nsStatusItem, nativeWindow, C.int(offset))

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
			// Cache on the ObjC controller for the event monitor.
			C.systemTraySetCachedMenu(s.nsStatusItem, s.nsMenu)
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

func (s *macosSystemTray) setLabel(label string) {
	s.label = label
	if !hasANSICodes(label) {
		C.systemTraySetLabel(s.nsStatusItem, C.CString(label))
		return
	}
	parts, err := SystemTrayLabelParser(label)
	if err != nil || len(parts) == 0 {
		C.systemTraySetLabel(s.nsStatusItem, C.CString(label))
		return
	}
	cLabel, cFg, cBg := partToCStrings(parts[0])
	attr := C.createAttributedString(cLabel, cFg, cBg)
	freeCStrings(cLabel, cFg, cBg)
	for _, p := range parts[1:] {
		cLabel, cFg, cBg = partToCStrings(p)
		attr = C.appendAttributedString(attr, cLabel, cFg, cBg)
		freeCStrings(cLabel, cFg, cBg)
	}
	C.systemTraySetANSILabel(s.nsStatusItem, attr)
}

func hasANSICodes(s string) bool {
	return strings.Contains(s, "\033[")
}

func partToCStrings(p SystemTrayLabelPart) (label, fg, bg *C.char) {
	label = C.CString(p.Text)
	if p.FgColor != "" {
		fg = C.CString(p.FgColor)
	}
	if p.BgColor != "" {
		bg = C.CString(p.BgColor)
	}
	return
}

func freeCStrings(label, fg, bg *C.char) {
	C.free(unsafe.Pointer(label))
	if fg != nil {
		C.free(unsafe.Pointer(fg))
	}
	if bg != nil {
		C.free(unsafe.Pointer(bg))
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
