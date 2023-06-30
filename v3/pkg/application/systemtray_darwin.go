//go:build darwin

package application

/*
#cgo CFLAGS: -mmacosx-version-min=10.13 -x objective-c
#cgo LDFLAGS: -framework Cocoa -framework WebKit

#include "Cocoa/Cocoa.h"
#include "menuitem_darwin.h"

// Create a new system tray
void* systemTrayNew() {
	NSStatusItem *statusItem = [[[NSStatusBar systemStatusBar] statusItemWithLength:NSVariableStatusItemLength] retain];
	return (void*)statusItem;
}

void systemTraySetLabel(void* nsStatusItem, char *label) {
	if( label == NULL ) {
		return;
	}
	// Set the label on the main thread
	dispatch_async(dispatch_get_main_queue(), ^{
		NSStatusItem *statusItem = (NSStatusItem *)nsStatusItem;
		statusItem.button.title = [NSString stringWithUTF8String:label];
		free(label);
	});
}

// Create an nsimage from a byte array
NSImage* imageFromBytes(const unsigned char *bytes, int length) {
	NSData *data = [NSData dataWithBytes:bytes length:length];
	NSImage *image = [[NSImage alloc] initWithData:data];
	return image;
}

// Set the icon on the system tray
void systemTraySetIcon(void* nsStatusItem, void* nsImage, int position, bool isTemplate) {
	// Set the icon on the main thread
	dispatch_async(dispatch_get_main_queue(), ^{
		NSStatusItem *statusItem = (NSStatusItem *)nsStatusItem;
		NSImage *image = (NSImage *)nsImage;

		NSStatusBar *statusBar = [NSStatusBar systemStatusBar];
		CGFloat thickness = [statusBar thickness];
		[image setSize:NSMakeSize(thickness, thickness)];
		if( isTemplate ) {
			[image setTemplate:YES];
		}
		statusItem.button.image = [image autorelease];
		statusItem.button.imagePosition = position;
	});
}

// Add menu to system tray
void systemTraySetMenu(void* nsStatusItem, void* nsMenu) {
	// Set the menu on the main thread
	dispatch_async(dispatch_get_main_queue(), ^{
		NSStatusItem *statusItem = (NSStatusItem *)nsStatusItem;
		NSMenu *menu = (NSMenu *)nsMenu;
		statusItem.menu = menu;
	});
}

// Destroy system tray
void systemTrayDestroy(void* nsStatusItem) {
	// Remove the status item from the status bar and its associated menu
	dispatch_async(dispatch_get_main_queue(), ^{
		NSStatusItem *statusItem = (NSStatusItem *)nsStatusItem;
		[[NSStatusBar systemStatusBar] removeStatusItem:statusItem];
		[statusItem release];
	});
}

void systemTrayGetBounds(void* nsStatusItem, NSRect *rect) {
	NSStatusItem *statusItem = (NSStatusItem *)nsStatusItem;
	*rect = statusItem.button.window.frame;
}

// Get the screen for the system tray
NSScreen* getScreenForSystemTray(void* nsStatusItem) {
	NSStatusItem *statusItem = (NSStatusItem *)nsStatusItem;
	NSRect frame = statusItem.button.frame;
	NSArray<NSScreen *> *screens = NSScreen.screens;
	NSScreen *associatedScreen = nil;

	for (NSScreen *screen in screens) {
		if (NSPointInRect(frame.origin, screen.frame)) {
			associatedScreen = screen;
			break;
		}
	}
	return associatedScreen;
}

*/
import "C"
import (
	"unsafe"
)

type macosSystemTray struct {
	id    uint
	label string
	icon  []byte
	menu  *Menu

	nsStatusItem   unsafe.Pointer
	nsImage        unsafe.Pointer
	nsMenu         unsafe.Pointer
	iconPosition   int
	isTemplateIcon bool
}

func (s *macosSystemTray) setIconPosition(position int) {
	s.iconPosition = position
}

func (s *macosSystemTray) setMenu(menu *Menu) {
	s.menu = menu
}

func (s *macosSystemTray) positionWindow(window *WebviewWindow) error {

	// Get the trayBounds of this system tray
	trayBounds, err := s.bounds()
	if err != nil {
		return err
	}

	// Get the current screen trayBounds
	currentScreen, err := s.getScreen()
	if err != nil {
		return err
	}

	screenBounds := currentScreen.Bounds

	// Determine which quadrant of the screen the system tray is in
	// ----------
	// | 1 | 2  |
	// ----------
	// | 3 | 4  |
	// ----------
	quadrant := 4
	if trayBounds.X < screenBounds.Width/2 {
		quadrant -= 1
	}
	if trayBounds.Y < screenBounds.Height/2 {
		quadrant -= 2
	}

	// Get the center height of the window
	windowWidthCenter := window.Width() / 2
	// Get the center height of the system tray
	systemTrayWidthCenter := trayBounds.Width / 2

	// Position the window based on the quadrant
	// It will be centered on the system tray and if it goes off-screen it will be moved back on screen
	switch quadrant {
	case 1:
		// The X will be 0 and the Y will be the system tray Y
		// Center the window on the system tray
		window.SetRelativePosition(0, trayBounds.Y)
	case 2:
		// The Y will be 0 and the X will make the center of the window line up with the center of the system tray
		windowX := trayBounds.X + systemTrayWidthCenter - windowWidthCenter
		// If the end of the window goes off-screen, move it back enough to be on screen
		if windowX+window.Width() > screenBounds.Width {
			windowX = screenBounds.Width - window.Width()
		}
		window.SetRelativePosition(windowX, 0)
	case 3:
		// The X will be 0 and the Y will be the system tray Y - the height of the window
		windowY := trayBounds.Y - window.Height()
		// If the end of the window goes off-screen, move it back enough to be on screen
		if windowY < 0 {
			windowY = 0
		}
		window.SetRelativePosition(0, windowY)
	case 4:
		// The Y will be 0 and the X will make the center of the window line up with the center of the system tray - the height of the window
		windowX := trayBounds.X + systemTrayWidthCenter - windowWidthCenter
		windowY := trayBounds.Y - window.Height()
		// If the end of the window goes off-screen, move it back enough to be on screen
		if windowX+window.Width() > screenBounds.Width {
			windowX = screenBounds.Width - window.Width()
		}
		window.SetRelativePosition(windowX, windowY)
	}
	return nil
}

func (s *macosSystemTray) getScreen() (*Screen, error) {
	cScreen := C.getScreenForSystemTray(s.nsStatusItem)
	return cScreenToScreen(cScreen), nil
}

func (s *macosSystemTray) bounds() (*Rect, error) {
	var rect C.NSRect
	rect = C.systemTrayGetBounds(s.nsStatusItem)
	return &Rect{
		X:      int(rect.origin.x),
		Y:      int(rect.origin.y),
		Width:  int(rect.size.width),
		Height: int(rect.size.height),
	}, nil
}

func (s *macosSystemTray) run() {
	globalApplication.dispatchOnMainThread(func() {
		if s.nsStatusItem != nil {
			Fatal("System tray '%d' already running", s.id)
		}
		s.nsStatusItem = unsafe.Pointer(C.systemTrayNew())
		if s.label != "" {
			C.systemTraySetLabel(s.nsStatusItem, C.CString(s.label))
		}
		if s.icon != nil {
			s.nsImage = unsafe.Pointer(C.imageFromBytes((*C.uchar)(&s.icon[0]), C.int(len(s.icon))))
			C.systemTraySetIcon(s.nsStatusItem, s.nsImage, C.int(s.iconPosition), C.bool(s.isTemplateIcon))
		}
		if s.menu != nil {
			s.menu.Update()
			// Convert impl to macosMenu object
			s.nsMenu = (s.menu.impl).(*macosMenu).nsMenu
			C.systemTraySetMenu(s.nsStatusItem, s.nsMenu)
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

func newSystemTrayImpl(s *SystemTray) systemTrayImpl {
	return &macosSystemTray{
		id:             s.id,
		label:          s.label,
		icon:           s.icon,
		menu:           s.menu,
		iconPosition:   s.iconPosition,
		isTemplateIcon: s.isTemplateIcon,
	}
}

func (s *macosSystemTray) setLabel(label string) {
	s.label = label
	C.systemTraySetLabel(s.nsStatusItem, C.CString(label))
}

func (s *macosSystemTray) destroy() {
	// Remove the status item from the status bar and its associated menu
	C.systemTrayDestroy(s.nsStatusItem)
}
