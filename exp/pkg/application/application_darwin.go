//go:build darwin

package application

/*

#cgo CFLAGS:  -x objective-c
#cgo LDFLAGS: -framework Cocoa -mmacosx-version-min=10.13

#include "application.h"
#include "app_delegate.h"
#include <stdlib.h>

#import <Cocoa/Cocoa.h>

static AppDelegate *appDelegate = nil;

static void init(void) {
    [NSApplication sharedApplication];
    appDelegate = [[AppDelegate alloc] init];
    [NSApp setDelegate:appDelegate];
}

static void setActivationPolicy(int policy) {
    [appDelegate setApplicationActivationPolicy:policy];
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

*/
import "C"
import (
	"unsafe"

	"github.com/wailsapp/wails/exp/pkg/options"
)

type macosApp struct {
	options         *options.Application
	applicationMenu unsafe.Pointer
}

func (m macosApp) setApplicationMenu(menu *Menu) {
	menu.Update()
	// Convert impl to macosMenu object
	m.applicationMenu = (menu.impl).(*macosMenu).nsMenu
	C.setApplicationMenu(m.applicationMenu)
}

func (m macosApp) run() error {
	C.run()
	return nil
}

func (m macosApp) destroy() {
	C.destroyApp()
}

func newPlatformApp(appOptions *options.Application) *macosApp {
	if appOptions == nil {
		appOptions = options.ApplicationDefaults
	}
	C.init()
	C.setActivationPolicy(C.int(appOptions.Mac.ActivationPolicy))
	return &macosApp{
		options: appOptions,
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

//export processMenuItemClick
func processMenuItemClick(menuID C.uint) {
	menuItemClicked <- uint(menuID)
}
