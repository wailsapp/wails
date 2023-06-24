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
