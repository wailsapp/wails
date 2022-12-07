//go:build darwin

package application

/*
#cgo CFLAGS: -mmacosx-version-min=10.10 -x objective-c
#cgo LDFLAGS: -framework Cocoa

#include "application.h"
#include "window_delegate.h"
#include <stdlib.h>
#include "Cocoa/Cocoa.h"

// Create a new Window
void* windowNew(int width, int height) {
	NSWindow* window = [[NSWindow alloc] initWithContentRect:NSMakeRect(0, 0, width-1, height-1)
		styleMask:NSTitledWindowMask | NSClosableWindowMask | NSMiniaturizableWindowMask | NSResizableWindowMask
		backing:NSBackingStoreBuffered
		defer:NO];
	// Create delegate
	WindowDelegate* delegate = [[WindowDelegate alloc] init];
	// Set delegate
	[window setDelegate:delegate];

	delegate.hideOnClose = false;
	return window;
}

// Set the title of the NSWindow
void windowSetTitle(void* nsWindow, char* title) {
	// Set window title on main thread
	dispatch_async(dispatch_get_main_queue(), ^{
		NSString* nsTitle = [NSString stringWithUTF8String:title];
		[(NSWindow*)nsWindow setTitle:nsTitle];
		free(title);
	});
}

// Set the size of the NSWindow
void windowSetSize(void* nsWindow, int width, int height) {
	// Set window size on main thread
	dispatch_async(dispatch_get_main_queue(), ^{
		[(NSWindow*)nsWindow setContentSize:NSMakeSize(width, height)];
		NSRect frame = [(NSWindow*)nsWindow frame];
		frame.size.width = width;
		frame.size.height = height;
		[(NSWindow*)nsWindow setFrame:frame display:YES];
	});

}

// Show the NSWindow
void windowShow(void* nsWindow) {
	[(NSWindow*)nsWindow makeKeyAndOrderFront:nil];
}

// Hide the NSWindow
void windowHide(void* nsWindow) {
	[(NSWindow*)nsWindow orderOut:nil];
}

// Set NSWindow always on top
void windowSetAlwaysOnTop(void* nsWindow, bool alwaysOnTop) {
	// Set window always on top on main thread
	dispatch_async(dispatch_get_main_queue(), ^{
		[(NSWindow*)nsWindow setLevel:alwaysOnTop ? NSStatusWindowLevel : NSNormalWindowLevel];
	});
}

*/
import "C"
import (
	"unsafe"

	"github.com/wailsapp/wails/exp/pkg/options"
)

type macosWindow struct {
	nsWindow unsafe.Pointer
	options  *options.Window
}

func (w *macosWindow) setAlwaysOnTop(alwaysOnTop bool) {
	C.windowSetAlwaysOnTop(w.nsWindow, C.bool(alwaysOnTop))
}

func newWindowImpl(options *options.Window) *macosWindow {
	result := &macosWindow{
		options: options,
	}
	return result
}

func (w *macosWindow) setTitle(title string) {
	cTitle := C.CString(title)
	C.windowSetTitle(w.nsWindow, cTitle)
}

func (w *macosWindow) setSize(width, height int) {
	C.windowSetSize(w.nsWindow, C.int(width), C.int(height))
}

func (w *macosWindow) Run() error {
	w.nsWindow = C.windowNew(C.int(w.options.Width), C.int(w.options.Height))
	w.setTitle(w.options.Title)
	w.setAlwaysOnTop(w.options.AlwaysOnTop)
	C.windowShow(w.nsWindow)
	return nil
}
