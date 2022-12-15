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

// Get the application name
static char *getAppName(void) {
	NSString *appName = [NSRunningApplication currentApplication].localizedName;
	if( appName == nil ) {
		appName = [[NSProcessInfo processInfo] processName];
	}
	return strdup([appName UTF8String]);
}

// Call the copy selector on the pasteboard
static void copyToPasteboard(char *text) {
	NSPasteboard *pasteboard = [NSPasteboard generalPasteboard];
	[pasteboard clearContents];
	[pasteboard setString:[NSString stringWithUTF8String:text] forType:NSPasteboardTypeString];
}

// Call the paste selector on the pasteboard
static char *pasteFromPasteboard(void) {
	NSPasteboard *pasteboard = [NSPasteboard generalPasteboard];
	NSString *text = [pasteboard stringForType:NSPasteboardTypeString];
	if( text == nil ) {
		return NULL;
	}
	return strdup([text UTF8String]);
}

// Call paste selector to paste text
static void paste(void) {
	[NSApp sendAction:@selector(paste:) to:nil from:nil];
}

// Call copy selector to copy text
static void copy(void) {
	[NSApp sendAction:@selector(copy:) to:nil from:nil];
}

// Call cut selector to cut text
static void cut(void) {
	[NSApp sendAction:@selector(cut:) to:nil from:nil];
}

// Call selectAll selector to select all text
static void selectAll(void) {
	[NSApp sendAction:@selector(selectAll:) to:nil from:nil];
}

// Call delete selector to delete text
static void delete(void) {
	[NSApp sendAction:@selector(delete:) to:nil from:nil];
}

// Call undo selector to undo text
static void undo(void) {
	[NSApp sendAction:@selector(undo:) to:nil from:nil];
}

// Call redo selector to redo text
static void redo(void) {
	[NSApp sendAction:@selector(redo:) to:nil from:nil];
}

// Call startSpeaking selector to start speaking text
static void startSpeaking(void) {
	[NSApp sendAction:@selector(startSpeaking:) to:nil from:nil];
}

// Call stopSpeaking selector to stop speaking text
static void stopSpeaking(void) {
	[NSApp sendAction:@selector(stopSpeaking:) to:nil from:nil];
}

static void pasteAndMatchStyle(void) {
	[NSApp sendAction:@selector(pasteAndMatchStyle:) to:nil from:nil];
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

func (m *macosApp) setApplicationMenu(menu *Menu) {
	if menu == nil {
		// Create a default menu for mac
		menu = m.createDefaultApplicationMenu()
	}
	menu.Update()
	// Convert impl to macosMenu object
	m.applicationMenu = (menu.impl).(*macosMenu).nsMenu
	C.setApplicationMenu(m.applicationMenu)
}

func (m *macosApp) run() error {
	C.run()
	return nil
}

func (m *macosApp) destroy() {
	C.destroyApp()
}

func (m *macosApp) createDefaultApplicationMenu() *Menu {
	// Create a default menu for mac
	menu := NewMenu()
	cAppName := C.getAppName()
	defer C.free(unsafe.Pointer(cAppName))
	appName := C.GoString(cAppName)
	appMenu := menu.AddSubmenu(appName)
	appMenu.Add("Quit " + appName).SetAccelerator("CmdOrCtrl+q").OnClick(func(ctx *Context) {
		globalApplication.Quit()
	})
	editMenu := menu.AddSubmenu("Edit")
	editMenu.Add("Undo").SetAccelerator("CmdOrCtrl+z").OnClick(func(ctx *Context) {
		C.undo()
	})
	editMenu.Add("Redo").SetAccelerator("CmdOrCtrl+Shift+z").OnClick(func(ctx *Context) {
		C.redo()
	})
	editMenu.AddSeparator()
	editMenu.Add("Cut").SetAccelerator("CmdOrCtrl+x").OnClick(func(ctx *Context) {
		C.cut()
	})
	editMenu.Add("Copy").SetAccelerator("CmdOrCtrl+c").OnClick(func(ctx *Context) {
		C.copy()
	})
	editMenu.Add("Paste").SetAccelerator("CmdOrCtrl+v").OnClick(func(ctx *Context) {
		C.paste()
	})
	editMenu.Add("Paste and Match Style").SetAccelerator("CmdOrCtrl+OptionOrAlt+Shift+v").OnClick(func(ctx *Context) {
		C.pasteAndMatchStyle()
	})
	editMenu.Add("Delete").SetAccelerator("backspace").OnClick(func(ctx *Context) {
		C.delete()
	})
	editMenu.Add("Select All").SetAccelerator("CmdOrCtrl+a").OnClick(func(ctx *Context) {
		C.selectAll()
	})
	editMenu.AddSeparator()
	speechMenu := editMenu.AddSubmenu("Speech")
	speechMenu.Add("Start Speaking").SetAccelerator("CmdOrCtrl+OptionOrAlt+Shift+.").OnClick(func(ctx *Context) {
		C.startSpeaking()
	})
	speechMenu.Add("Stop Speaking").SetAccelerator("CmdOrCtrl+OptionOrAlt+Shift+,").OnClick(func(ctx *Context) {
		C.stopSpeaking()
	})
	return menu
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
