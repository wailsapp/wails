//go:build darwin

package application

/*

#cgo CFLAGS: -mmacosx-version-min=10.13 -x objective-c
#cgo LDFLAGS: -framework Cocoa -mmacosx-version-min=10.13

#include "application_darwin.h"
#include "application_darwin_delegate.h"
#include "webview_window_darwin.h"
#include <stdlib.h>

extern void registerListener(unsigned int event);

#import <Cocoa/Cocoa.h>
#import <Foundation/Foundation.h>

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

	NSDistributedNotificationCenter *center = [NSDistributedNotificationCenter defaultCenter];
	[center addObserver:appDelegate selector:@selector(themeChanged:) name:@"AppleInterfaceThemeChangedNotification" object:nil];

}

static bool isDarkMode(void) {
	NSUserDefaults* userDefaults = [NSUserDefaults standardUserDefaults];
	if (userDefaults == nil) {
		return false;
	}

	NSString *interfaceStyle = [userDefaults stringForKey:@"AppleInterfaceStyle"];
	if (interfaceStyle == nil) {
		return false;
	}

	return [interfaceStyle isEqualToString:@"Dark"];
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
		[NSApp abortModal];
    }
}

// destroyApp destroys the application
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

static const char* serializationNSDictionary(void *dict) {
	@autoreleasepool {
		NSDictionary *nsDict = (__bridge NSDictionary *)dict;

		if ([NSJSONSerialization isValidJSONObject:nsDict]) {
			NSError *error;
			NSData *data = [NSJSONSerialization dataWithJSONObject:nsDict options:kNilOptions error:&error];
			NSString *result = [[NSString alloc]initWithData:data encoding:NSUTF8StringEncoding];

			return strdup([result UTF8String]);
		}
	}

	return nil;
}

static void startSingleInstanceListener(const char *uniqueID) {
	// Convert to NSString
	NSString *uid = [NSString stringWithUTF8String:uniqueID];
	[[NSDistributedNotificationCenter defaultCenter] addObserver:appDelegate
          selector:@selector(handleSecondInstanceNotification:) name:uid object:nil];
}
*/
import "C"
import (
	"encoding/json"
	"unsafe"

	"github.com/wailsapp/wails/v3/internal/operatingsystem"

	"github.com/wailsapp/wails/v3/internal/assetserver/webview"
	"github.com/wailsapp/wails/v3/pkg/events"
)

type macosApp struct {
	applicationMenu unsafe.Pointer
	parent          *App
}

func (m *macosApp) isDarkMode() bool {
	return bool(C.isDarkMode())
}

func getNativeApplication() *macosApp {
	return globalApplication.impl.(*macosApp)
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
		menu = DefaultApplicationMenu()
	}
	menu.Update()

	// Convert impl to macosMenu object
	m.applicationMenu = (menu.impl).(*macosMenu).nsMenu
	C.setApplicationMenu(m.applicationMenu)
}

func (m *macosApp) run() error {
	if m.parent.options.SingleInstance != nil {
		cUniqueID := C.CString(m.parent.options.SingleInstance.UniqueID)
		defer C.free(unsafe.Pointer(cUniqueID))
		C.startSingleInstanceListener(cUniqueID)
	}
	// Add a hook to the ApplicationDidFinishLaunching event
	m.parent.OnApplicationEvent(events.Mac.ApplicationDidFinishLaunching, func(*ApplicationEvent) {
		C.setApplicationShouldTerminateAfterLastWindowClosed(C.bool(m.parent.options.Mac.ApplicationShouldTerminateAfterLastWindowClosed))
		C.setActivationPolicy(C.int(m.parent.options.Mac.ActivationPolicy))
		C.activateIgnoringOtherApps()
	})
	m.setupCommonEvents()
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

func (m *macosApp) GetFlags(options Options) map[string]any {
	if options.Flags == nil {
		options.Flags = make(map[string]any)
	}
	return options.Flags
}

func newPlatformApp(app *App) *macosApp {
	C.init()
	return &macosApp{
		parent: app,
	}
}

//export processApplicationEvent
func processApplicationEvent(eventID C.uint, data unsafe.Pointer) {
	event := newApplicationEvent(events.ApplicationEventType(eventID))

	if data != nil {
		dataCStrJSON := C.serializationNSDictionary(data)
		if dataCStrJSON != nil {
			defer C.free(unsafe.Pointer(dataCStrJSON))

			dataJSON := C.GoString(dataCStrJSON)
			var result map[string]any
			err := json.Unmarshal([]byte(dataJSON), &result)

			if err != nil {
				panic(err)
			}

			event.Context().setData(result)
		}
	}

	switch event.Id {
	case uint(events.Mac.ApplicationDidChangeTheme):
		isDark := globalApplication.IsDarkMode()
		event.Context().setIsDarkMode(isDark)
	}
	applicationEvents <- event
}

//export processWindowEvent
func processWindowEvent(windowID C.uint, eventID C.uint) {
	windowEvents <- &windowEvent{
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

//export processWindowKeyDownEvent
func processWindowKeyDownEvent(windowID C.uint, acceleratorString *C.char) {
	windowKeyEvents <- &windowKeyEvent{
		windowId:          uint(windowID),
		acceleratorString: C.GoString(acceleratorString),
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

//export shouldQuitApplication
func shouldQuitApplication() C.bool {
	// TODO: This should be configurable
	return C.bool(globalApplication.shouldQuit())
}

//export cleanup
func cleanup() {
	globalApplication.cleanup()
}

func (a *App) logPlatformInfo() {
	info, err := operatingsystem.Info()
	if err != nil {
		a.error("error getting OS info: %w", err)
		return
	}

	a.info("Platform Info:", info.AsLogSlice()...)

}

func (a *App) platformEnvironment() map[string]any {
	return map[string]any{}
}

func fatalHandler(errFunc func(error)) {
	return
}

//export HandleOpenFile
func HandleOpenFile(filePath *C.char) {
	goFilepath := C.GoString(filePath)
	// Create new application event context
	eventContext := newApplicationEventContext()
	eventContext.setOpenedWithFile(goFilepath)
	// EmitEvent application started event
	applicationEvents <- &ApplicationEvent{
		Id:  uint(events.Common.ApplicationOpenedWithFile),
		ctx: eventContext,
	}
}
