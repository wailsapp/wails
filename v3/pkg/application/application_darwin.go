//go:build darwin

package application

/*

#cgo CFLAGS: -mmacosx-version-min=10.13 -x objective-c
#cgo LDFLAGS: -framework Cocoa -mmacosx-version-min=10.13

#include "application.h"
#include "app_delegate.h"
#include "webview_window.h"
#include <stdlib.h>

extern void registerListener(unsigned int event);

#import <Cocoa/Cocoa.h>

static AppDelegate *appDelegate = nil;

static void init(void) {
    [NSApplication sharedApplication];
    appDelegate = [[AppDelegate alloc] init];
    [NSApp setDelegate:appDelegate];

	[NSEvent addLocalMonitorForEventsMatchingMask:NSEventMaskLeftMouseDown handler:^NSEvent * _Nullable(NSEvent * _Nonnull event) {
		NSWindow* eventWindow = [event window];
		if (eventWindow == nil ) {
			return event;
        }
		WebviewWindowDelegate* windowDelegate = (WebviewWindowDelegate*)[eventWindow delegate];
		if (windowDelegate == nil) {
			return event;
		}
		if ([windowDelegate respondsToSelector:@selector(handleLeftMouseDown:)]) {
			[windowDelegate handleLeftMouseDown:event];
		}
		return event;
	}];

	[NSEvent addLocalMonitorForEventsMatchingMask:NSEventMaskLeftMouseUp handler:^NSEvent * _Nullable(NSEvent * _Nonnull event) {
		NSWindow* eventWindow = [event window];
		if (eventWindow == nil ) {
			return event;
        }
		WebviewWindowDelegate* windowDelegate = (WebviewWindowDelegate*)[eventWindow delegate];
		if (windowDelegate == nil) {
			return event;
		}
		if ([windowDelegate respondsToSelector:@selector(handleLeftMouseUp:)]) {
			[windowDelegate handleLeftMouseUp:eventWindow];
		}
		return event;
	}];
}

static void setApplicationShouldTerminateAfterLastWindowClosed(bool shouldTerminate) {
	// Get the NSApp delegate
	AppDelegate *appDelegate = (AppDelegate*)[NSApp delegate];
	// Set the applicationShouldTerminateAfterLastWindowClosed boolean
	appDelegate.shouldTerminateWhenLastWindowClosed = shouldTerminate;
}

static void setActivationPolicy(int policy) {
    [NSApp setActivationPolicy:policy];
}

static void activateIgnoringOtherApps() {
	[NSApp activateIgnoringOtherApps:YES];
}

static void run(void) {
    @autoreleasepool {
        [NSApp run];
        [appDelegate release];
    }
}

// Destroy application
static void destroyApp(void) {
	[NSApp terminate:nil];
}

// Set the application menu
static void setApplicationMenu(void *menu) {
	NSMenu *nsMenu = (__bridge NSMenu *)menu;
	[NSApp setMainMenu:menu];
}

// Get the application name
static char* getAppName(void) {
	NSString *appName = [NSRunningApplication currentApplication].localizedName;
	if( appName == nil ) {
		appName = [[NSProcessInfo processInfo] processName];
	}
	return strdup([appName UTF8String]);
}

// get the current window ID
static unsigned int getCurrentWindowID(void) {
	NSWindow *window = [NSApp keyWindow];
	// Get the window delegate
	WebviewWindowDelegate *delegate = (WebviewWindowDelegate*)[window delegate];
	return delegate.windowId;
}

// Set the application icon
static void setApplicationIcon(void *icon, int length) {
    // On main thread
	dispatch_async(dispatch_get_main_queue(), ^{
		NSImage *image = [[NSImage alloc] initWithData:[NSData dataWithBytes:icon length:length]];
		[NSApp setApplicationIconImage:image];
	});
}

// Hide the application
static void hide(void) {
	[NSApp hide:nil];
}

// Show the application
static void show(void) {
	[NSApp unhide:nil];
}

*/
import "C"
import (
	"unsafe"

	"github.com/wailsapp/wails/v2/pkg/assetserver/webview"
	"github.com/wailsapp/wails/v3/pkg/events"
)

type macosApp struct {
	applicationMenu unsafe.Pointer
	parent          *App
}

func (m *macosApp) hide() {
	C.hide()
}

func (m *macosApp) show() {
	C.show()
}

func (m *macosApp) on(eventID uint) {
	C.registerListener(C.uint(eventID))
}

func (m *macosApp) setIcon(icon []byte) {
	C.setApplicationIcon(unsafe.Pointer(&icon[0]), C.int(len(icon)))
}

func (m *macosApp) name() string {
	appName := C.getAppName()
	defer C.free(unsafe.Pointer(appName))
	return C.GoString(appName)
}

func (m *macosApp) getCurrentWindowID() uint {
	return uint(C.getCurrentWindowID())
}

func (m *macosApp) setApplicationMenu(menu *Menu) {
	if menu == nil {
		// Create a default menu for mac
		menu = defaultApplicationMenu()
	}
	menu.Update()

	// Convert impl to macosMenu object
	m.applicationMenu = (menu.impl).(*macosMenu).nsMenu
	C.setApplicationMenu(m.applicationMenu)
}

func (m *macosApp) run() error {
	// Add a hook to the ApplicationDidFinishLaunching event
	m.parent.On(events.Mac.ApplicationDidFinishLaunching, func() {
		C.setApplicationShouldTerminateAfterLastWindowClosed(C.bool(m.parent.options.Mac.ApplicationShouldTerminateAfterLastWindowClosed))
		C.setActivationPolicy(C.int(m.parent.options.Mac.ActivationPolicy))
		C.activateIgnoringOtherApps()
	})
	// setup event listeners
	for eventID := range m.parent.applicationEventListeners {
		m.on(eventID)
	}
	C.run()
	return nil
}

func (m *macosApp) destroy() {
	C.destroyApp()
}

func newPlatformApp(app *App) *macosApp {
	C.init()
	return &macosApp{
		parent: app,
	}
}

//export processApplicationEvent
func processApplicationEvent(eventID C.uint) {
	applicationEvents <- uint(eventID)
}

//export processWindowEvent
func processWindowEvent(windowID C.uint, eventID C.uint) {
	windowEvents <- &WindowEvent{
		WindowID: uint(windowID),
		EventID:  uint(eventID),
	}
}

//export processMessage
func processMessage(windowID C.uint, message *C.char) {
	windowMessageBuffer <- &windowMessage{
		windowId: uint(windowID),
		message:  C.GoString(message),
	}
}

//export processURLRequest
func processURLRequest(windowID C.uint, wkUrlSchemeTask unsafe.Pointer) {
	webviewRequests <- &webViewAssetRequest{
		Request:    webview.NewRequest(wkUrlSchemeTask),
		windowId:   uint(windowID),
		windowName: globalApplication.getWindowForID(uint(windowID)).Name(),
	}
}

//export processDragItems
func processDragItems(windowID C.uint, arr **C.char, length C.int) {
	var filenames []string
	// Convert the C array to a Go slice
	goSlice := (*[1 << 30]*C.char)(unsafe.Pointer(arr))[:length:length]
	for _, str := range goSlice {
		filenames = append(filenames, C.GoString(str))
	}
	windowDragAndDropBuffer <- &dragAndDropMessage{
		windowId:  uint(windowID),
		filenames: filenames,
	}
}

//export processMenuItemClick
func processMenuItemClick(menuID C.uint) {
	menuItemClicked <- uint(menuID)
}

func setIcon(icon []byte) {
	if icon == nil {
		return
	}
	C.setApplicationIcon(unsafe.Pointer(&icon[0]), C.int(len(icon)))
}
