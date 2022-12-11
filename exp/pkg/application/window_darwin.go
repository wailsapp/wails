//go:build darwin

package application

/*
#cgo CFLAGS: -mmacosx-version-min=10.13 -x objective-c
#cgo LDFLAGS: -framework Cocoa -framework WebKit

#include "application.h"
#include "window_delegate.h"
#include <stdlib.h>
#include "Cocoa/Cocoa.h"
#import <WebKit/WebKit.h>


// Create a new Window
void* windowNew(unsigned int id, int width, int height) {
	NSWindow* window = [[NSWindow alloc] initWithContentRect:NSMakeRect(0, 0, width-1, height-1)
		styleMask:NSWindowStyleMaskTitled | NSWindowStyleMaskClosable | NSWindowStyleMaskMiniaturizable | NSWindowStyleMaskResizable
		backing:NSBackingStoreBuffered
		defer:NO];

	// Create delegate
	WindowDelegate* delegate = [[WindowDelegate alloc] init];
	// Set delegate
	[window setDelegate:delegate];
	delegate.windowId = id;

	// Add NSView to window
	NSView* view = [[NSView alloc] initWithFrame:NSMakeRect(0, 0, width-1, height-1)];
	[view setAutoresizingMask:NSViewWidthSizable | NSViewHeightSizable];
	[window setContentView:view];

	// Embed wkwebview in window
	NSRect frame = NSMakeRect(0, 0, width, height);
	WKWebViewConfiguration* config = [[WKWebViewConfiguration alloc] init];
	config.suppressesIncrementalRendering = true;
    config.applicationNameForUserAgent = @"wails.io";

	// Setup user content controller
    WKUserContentController* userContentController = [WKUserContentController new];
    [userContentController addScriptMessageHandler:delegate name:@"external"];
    config.userContentController = userContentController;

	WKWebView* webView = [[WKWebView alloc] initWithFrame:frame configuration:config];
	[view addSubview:webView];
	delegate.webView = webView;


	delegate.hideOnClose = false;
	return window;
}

// Make NSWindow transparent
void windowSetTransparent(void* nsWindow) {
    // On main thread
	dispatch_async(dispatch_get_main_queue(), ^{
	NSWindow* window = (NSWindow*)nsWindow;
	[window setOpaque:NO];
	[window setBackgroundColor:[NSColor clearColor]];
	});
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
	// Show window on main thread
	dispatch_async(dispatch_get_main_queue(), ^{
		[(NSWindow*)nsWindow makeKeyAndOrderFront:nil];
	});
}

// Hide the NSWindow
void windowHide(void* nsWindow) {
	// Hide window on main thread
	dispatch_async(dispatch_get_main_queue(), ^{
		[(NSWindow*)nsWindow orderOut:nil];
	});
}

// Set NSWindow always on top
void windowSetAlwaysOnTop(void* nsWindow, bool alwaysOnTop) {
	// Set window always on top on main thread
	dispatch_async(dispatch_get_main_queue(), ^{
		[(NSWindow*)nsWindow setLevel:alwaysOnTop ? NSStatusWindowLevel : NSNormalWindowLevel];
	});
}

// Load URL in NSWindow
void navigationLoadURL(void* nsWindow, char* url) {
	// Load URL on main thread
	dispatch_async(dispatch_get_main_queue(), ^{
		NSURL* nsURL = [NSURL URLWithString:[NSString stringWithUTF8String:url]];
		NSURLRequest* request = [NSURLRequest requestWithURL:nsURL];
		[[(WindowDelegate*)[(NSWindow*)nsWindow delegate] webView] loadRequest:request];
		free(url);
	});
}

// Set NSWindow resizable
void windowSetResizable(void* nsWindow, bool resizable) {
	// Set window resizable on main thread
	dispatch_async(dispatch_get_main_queue(), ^{
		[(NSWindow*)nsWindow setStyleMask:resizable ? NSWindowStyleMaskTitled | NSWindowStyleMaskClosable | NSWindowStyleMaskMiniaturizable | NSWindowStyleMaskResizable : NSWindowStyleMaskTitled | NSWindowStyleMaskClosable | NSWindowStyleMaskMiniaturizable];
	});
}

// Set NSWindow min size
void windowSetMinSize(void* nsWindow, int width, int height) {
	// Set window min size on main thread
	dispatch_async(dispatch_get_main_queue(), ^{
		[(NSWindow*)nsWindow setContentMinSize:NSMakeSize(width, height)];
	});
}

// Set NSWindow max size
void windowSetMaxSize(void* nsWindow, int width, int height) {
	// Set window max size on main thread
	dispatch_async(dispatch_get_main_queue(), ^{
		[(NSWindow*)nsWindow setContentMaxSize:NSMakeSize(width, height)];
	});
}

// Reset NSWindow min and max size
void windowResetMinSize(void* nsWindow) {
	// Reset window min size on main thread
	dispatch_async(dispatch_get_main_queue(), ^{
		[(NSWindow*)nsWindow setContentMinSize:NSMakeSize(0, 0)];
	});
}

void windowResetMaxSize(void* nsWindow) {
	// Reset window max size on main thread
	dispatch_async(dispatch_get_main_queue(), ^{
		[(NSWindow*)nsWindow setContentMaxSize:NSMakeSize(0, 0)];
	});
}

// Enable NSWindow devtools
void windowEnableDevTools(void* nsWindow) {
	// Enable devtools on main thread
	dispatch_async(dispatch_get_main_queue(), ^{
		// Get window delegate
		WindowDelegate* delegate = (WindowDelegate*)[(NSWindow*)nsWindow delegate];
		// Enable devtools in webview
		[delegate.webView.configuration.preferences setValue:@YES forKey:@"developerExtrasEnabled"];
	});
}

// Execute JS in NSWindow
void windowExecJS(void* nsWindow, char* js) {
	// Execute JS on main thread
	dispatch_async(dispatch_get_main_queue(), ^{
		// Get window delegate
		WindowDelegate* delegate = (WindowDelegate*)[(NSWindow*)nsWindow delegate];
		// Execute JS in webview
		[delegate.webView evaluateJavaScript:[NSString stringWithUTF8String:js] completionHandler:nil];
		free(js);
	});
}

// Make NSWindow backdrop translucent
void windowSetTranslucent(void* nsWindow) {
	// Set window transparent on main thread
	dispatch_async(dispatch_get_main_queue(), ^{

		// Get window
		NSWindow* window = (NSWindow*)nsWindow;

		// Get window delegate
		WindowDelegate* delegate = (WindowDelegate*)[(NSWindow*)nsWindow delegate];

		id contentView = [window contentView];
		NSVisualEffectView *effectView = [NSVisualEffectView alloc];
		NSRect bounds = [contentView bounds];
		[effectView initWithFrame:bounds];
		[effectView setAutoresizingMask:NSViewWidthSizable | NSViewHeightSizable];
		[effectView setBlendingMode:NSVisualEffectBlendingModeBehindWindow];
		[effectView setState:NSVisualEffectStateActive];
		[contentView addSubview:effectView positioned:NSWindowBelow relativeTo:nil];
	});
}

// Make webview background transparent
void webviewSetTransparent(void* nsWindow) {
	// Set webview transparent on main thread
	dispatch_async(dispatch_get_main_queue(), ^{
		// Get window delegate
		WindowDelegate* delegate = (WindowDelegate*)[(NSWindow*)nsWindow delegate];
		// Set webview background transparent
		[delegate.webView setValue:@NO forKey:@"drawsBackground"];
	});
}

// Set webview background color
void webviewSetBackgroundColor(void* nsWindow, int r, int g, int b, int alpha) {
	// Set webview background color on main thread
	dispatch_async(dispatch_get_main_queue(), ^{
		// Get window delegate
		WindowDelegate* delegate = (WindowDelegate*)[(NSWindow*)nsWindow delegate];
		// Set webview background color
		[delegate.webView setValue:[NSColor colorWithRed:r/255.0 green:g/255.0 blue:b/255.0 alpha:alpha/255.0] forKey:@"backgroundColor"];
	});
}

// Set Window maximised
void windowSetMaximised(void* nsWindow) {
	// Set window maximized on main thread
	dispatch_async(dispatch_get_main_queue(), ^{
		[(NSWindow*)nsWindow zoom:nil];
	});
}

// Set Window fullscreen
void windowSetFullscreen(void* nsWindow) {
	// Set window fullscreen on main thread
	dispatch_async(dispatch_get_main_queue(), ^{
		[(NSWindow*)nsWindow toggleFullScreen:nil];
	});
}

// Set Window Minimised
void windowSetMinimised(void* nsWindow) {
	// Set window minimised on main thread
	dispatch_async(dispatch_get_main_queue(), ^{
		// Get screen that the window is on
		NSScreen* screen = [(NSWindow*)nsWindow screen];
		NSRect screenRect = [screen frame];
		// Set window to top left corner
		[(NSWindow*)nsWindow setFrame:NSMakeRect(0, screenRect.size.height, 0, 0) display:YES];
	});
}

// restore window to normal size
void windowRestore(void* nsWindow) {
	// Set window normal on main thread
	dispatch_async(dispatch_get_main_queue(), ^{
		// If window is fullscreen
		if([(NSWindow*)nsWindow styleMask] & NSWindowStyleMaskFullScreen) {
			[(NSWindow*)nsWindow toggleFullScreen:nil];
		}
		// If window is maximised
		if([(NSWindow*)nsWindow isZoomed]) {
			[(NSWindow*)nsWindow zoom:nil];
		}
		// If window in minimised
		if([(NSWindow*)nsWindow isMiniaturized]) {
			[(NSWindow*)nsWindow deminiaturize:nil];
		}
	});
}

bool windowIsMaximised(void* nsWindow) {
	return [(NSWindow*)nsWindow isZoomed];
}

bool windowIsFullscreen(void* nsWindow) {
	return [(NSWindow*)nsWindow styleMask] & NSWindowStyleMaskFullScreen;
}

bool windowIsMinimised(void* nsWindow) {
	return [(NSWindow*)nsWindow isMiniaturized];
}

// Set the titlebar style
void windowSetTitleBarAppearsTransparent(void* nsWindow, bool transparent) {
	// Set window titlebar style on main thread
	dispatch_async(dispatch_get_main_queue(), ^{
		if( transparent ) {
			[(NSWindow*)nsWindow setTitlebarAppearsTransparent:true];
		} else {
			[(NSWindow*)nsWindow setTitlebarAppearsTransparent:false];
		}
	});
}

// Set window fullsize content view
void windowSetFullSizeContent(void* nsWindow, bool fullSize) {
	// Set window fullsize content view on main thread
	dispatch_async(dispatch_get_main_queue(), ^{
		if( fullSize ) {
			[(NSWindow*)nsWindow setStyleMask:[(NSWindow*)nsWindow styleMask] | NSWindowStyleMaskFullSizeContentView];
		} else {
			[(NSWindow*)nsWindow setStyleMask:[(NSWindow*)nsWindow styleMask] & ~NSWindowStyleMaskFullSizeContentView];
		}
	});
}

// Set Hide Titlebar
void windowSetHideTitleBar(void* nsWindow, bool hideTitlebar) {
	// Set window titlebar hidden on main thread
	dispatch_async(dispatch_get_main_queue(), ^{
		if( hideTitlebar ) {
			[(NSWindow*)nsWindow setStyleMask:[(NSWindow*)nsWindow styleMask] & ~NSWindowStyleMaskTitled];
		} else {
			[(NSWindow*)nsWindow setStyleMask:[(NSWindow*)nsWindow styleMask] | NSWindowStyleMaskTitled];
		}
	});
}

// Set Hide Title in Titlebar
void windowSetHideTitle(void* nsWindow, bool hideTitle) {
	// Set window titlebar hidden on main thread
	dispatch_async(dispatch_get_main_queue(), ^{
		if( hideTitle ) {
			[(NSWindow*)nsWindow setTitleVisibility:NSWindowTitleHidden];
		} else {
			[(NSWindow*)nsWindow setTitleVisibility:NSWindowTitleVisible];
		}
	});
}

// Set Window use toolbar
void windowSetUseToolbar(void* nsWindow, bool useToolbar) {
	// Set window use toolbar on main thread
	dispatch_async(dispatch_get_main_queue(), ^{
		// get main window
		NSWindow* window = (NSWindow*)nsWindow;
		if( useToolbar ) {
			NSToolbar *toolbar = [[NSToolbar alloc] initWithIdentifier:@"wails.toolbar"];
			[toolbar autorelease];
			[window setToolbar:toolbar];
		} else {
			[window setToolbar:nil];
		}
	});
}

// Set Hide Toolbar Separator
void windowSetHideToolbarSeparator(void* nsWindow, bool hideSeparator) {
	// Set window hide toolbar separator on main thread
	dispatch_async(dispatch_get_main_queue(), ^{
		// get main window
		NSWindow* window = (NSWindow*)nsWindow;
		// get toolbar
		NSToolbar* toolbar = [window toolbar];
		// Return if toolbar nil
		if( toolbar == nil ) {
			return;
		}
		if( hideSeparator ) {
			[toolbar setShowsBaselineSeparator:false];
		} else {
			[toolbar setShowsBaselineSeparator:true];
		}
	});
}

// Set Window appearance type
void windowSetAppearanceTypeByName(void* nsWindow, const char *appearanceName) {
	// Set window appearance type on main thread
	dispatch_async(dispatch_get_main_queue(), ^{
		// get main window
		NSWindow* window = (NSWindow*)nsWindow;
		// set window appearance type by name
		// Convert appearance name to NSString
		NSString* appearanceNameString = [NSString stringWithUTF8String:appearanceName];
		// Set appearance
		[window setAppearance:[NSAppearance appearanceNamed:appearanceNameString]];

		free((void*)appearanceName);
	});
}

// Center window on current monitor
void windowCenter(void* nsWindow) {
	// Center window on main thread
	dispatch_async(dispatch_get_main_queue(), ^{
		// get main window
		NSWindow* window = (NSWindow*)nsWindow;
		[window center];
	});
}

// Get the current size of the window
void windowGetSize(void* nsWindow, int* width, int* height) {
	// get main window
	NSWindow* window = (NSWindow*)nsWindow;
	// get window frame
	NSRect frame = [window frame];
	// set width and height
	*width = frame.size.width;
	*height = frame.size.height;
}

// Get window width
int windowGetWidth(void* nsWindow) {
	// get main window
	NSWindow* window = (NSWindow*)nsWindow;
	// get window frame
	NSRect frame = [window frame];
	// return width
	return frame.size.width;
}

// Get window height
int windowGetHeight(void* nsWindow) {
	// get main window
	NSWindow* window = (NSWindow*)nsWindow;
	// get window frame
	NSRect frame = [window frame];
	// return height
	return frame.size.height;
}

// Get window position
void windowGetPosition(void* nsWindow, int* x, int* y) {
	// get main window
	NSWindow* window = (NSWindow*)nsWindow;
	// get window frame
	NSRect frame = [window frame];
	// set x and y
	*x = frame.origin.x;
	*y = frame.origin.y;
}


*/
import "C"
import (
	"sync"
	"unsafe"

	"github.com/wailsapp/wails/exp/pkg/options"
)

type macosWindow struct {
	id       uint
	nsWindow unsafe.Pointer
	options  *options.Window
}

func (w *macosWindow) center() {
	C.windowCenter(w.nsWindow)
}

func (w *macosWindow) isMinimised() bool {
	return w.syncMainThreadReturningBool(func() bool {
		return C.windowIsMinimised(w.nsWindow) == C.bool(true)
	})
}

func (w *macosWindow) isMaximised() bool {
	return w.syncMainThreadReturningBool(func() bool {
		return C.windowIsMaximised(w.nsWindow) == C.bool(true)
	})
}

func (w *macosWindow) isFullscreen() bool {
	return w.syncMainThreadReturningBool(func() bool {
		return C.windowIsFullscreen(w.nsWindow) == C.bool(true)
	})
}

func (w *macosWindow) syncMainThreadReturningBool(fn func() bool) bool {
	var wg sync.WaitGroup
	wg.Add(1)
	var result bool
	Dispatch(func() {
		result = fn()
		wg.Done()
	})
	wg.Wait()
	return result
}

func (w *macosWindow) restore() {
	// restore window to normal size
	C.windowRestore(w.nsWindow)
}

func (w *macosWindow) setMaximised() {
	C.windowSetMaximised(w.nsWindow)
}

func (w *macosWindow) setMinimised() {
	C.windowSetMinimised(w.nsWindow)
}

func (w *macosWindow) setFullscreen() {
	C.windowSetFullscreen(w.nsWindow)
}

func (w *macosWindow) restoreWindow() {
	C.windowRestore(w.nsWindow)
}

func (w *macosWindow) execJS(js string) {
	C.windowExecJS(w.nsWindow, C.CString(js))
}

func (w *macosWindow) navigateToURL(url string) {
	C.navigationLoadURL(w.nsWindow, C.CString(url))
}

func (w *macosWindow) setAlwaysOnTop(alwaysOnTop bool) {
	C.windowSetAlwaysOnTop(w.nsWindow, C.bool(alwaysOnTop))
}

func newWindowImpl(id uint, options *options.Window) *macosWindow {
	result := &macosWindow{
		id:      id,
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

func (w *macosWindow) setMinSize(width, height int) {
	C.windowSetMinSize(w.nsWindow, C.int(width), C.int(height))
}
func (w *macosWindow) setMaxSize(width, height int) {
	C.windowSetMaxSize(w.nsWindow, C.int(width), C.int(height))
}

func (w *macosWindow) setResizable(resizable bool) {
	C.windowSetResizable(w.nsWindow, C.bool(resizable))
}
func (w *macosWindow) enableDevTools() {
	C.windowEnableDevTools(w.nsWindow)
}

func (w *macosWindow) size() (int, int) {
	var width, height C.int
	var wg sync.WaitGroup
	wg.Add(1)
	Dispatch(func() {
		C.windowGetSize(w.nsWindow, &width, &height)
		wg.Done()
	})
	wg.Wait()
	return int(width), int(height)
}

func (w *macosWindow) width() int {
	var width C.int
	var wg sync.WaitGroup
	wg.Add(1)
	Dispatch(func() {
		width = C.windowGetWidth(w.nsWindow)
		wg.Done()
	})
	wg.Wait()
	return int(width)
}
func (w *macosWindow) height() int {
	var height C.int
	var wg sync.WaitGroup
	wg.Add(1)
	Dispatch(func() {
		height = C.windowGetHeight(w.nsWindow)
		wg.Done()
	})
	wg.Wait()
	return int(height)
}

func (w *macosWindow) run() {
	Dispatch(func() {
		w.nsWindow = C.windowNew(C.uint(w.id), C.int(w.options.Width), C.int(w.options.Height))
		w.setTitle(w.options.Title)
		w.setAlwaysOnTop(w.options.AlwaysOnTop)
		w.setResizable(!w.options.DisableResize)
		w.setMinSize(w.options.MinWidth, w.options.MinHeight)
		w.setMaxSize(w.options.MaxWidth, w.options.MaxHeight)
		if w.options.EnableDevTools {
			w.enableDevTools()
		}
		w.setBackgroundColor(w.options.BackgroundColour)
		if w.options.Mac != nil {
			macOptions := w.options.Mac
			switch macOptions.Backdrop {
			case options.MacBackdropTransparent:
				C.windowSetTransparent(w.nsWindow)
				C.webviewSetTransparent(w.nsWindow)
			case options.MacBackdropTranslucent:
				C.windowSetTranslucent(w.nsWindow)
				C.webviewSetTransparent(w.nsWindow)
			}

			if macOptions.TitleBar != nil {
				titleBarOptions := macOptions.TitleBar
				C.windowSetTitleBarAppearsTransparent(w.nsWindow, C.bool(titleBarOptions.AppearsTransparent))
				C.windowSetHideTitleBar(w.nsWindow, C.bool(titleBarOptions.Hide))
				C.windowSetHideTitle(w.nsWindow, C.bool(titleBarOptions.HideTitle))
				C.windowSetFullSizeContent(w.nsWindow, C.bool(titleBarOptions.FullSizeContent))
				C.windowSetUseToolbar(w.nsWindow, C.bool(titleBarOptions.UseToolbar))
				C.windowSetHideToolbarSeparator(w.nsWindow, C.bool(titleBarOptions.HideToolbarSeparator))
			}

			if macOptions.Appearance != "" {
				C.windowSetAppearanceTypeByName(w.nsWindow, C.CString(string(macOptions.Appearance)))
			}

			switch w.options.StartState {
			case options.WindowStateMaximised:
				w.setMaximised()
			case options.WindowStateMinimised:
				w.setMinimised()
			case options.WindowStateFullscreen:
				w.setFullscreen()

			}

		}
		C.windowCenter(w.nsWindow)

		if w.options.URL != "" {
			w.navigateToURL(w.options.URL)
		}
		C.windowShow(w.nsWindow)
	})
}

func (w *macosWindow) setBackgroundColor(colour *options.RGBA) {
	if colour == nil {
		return
	}
	C.webviewSetBackgroundColor(w.nsWindow, C.int(colour.Red), C.int(colour.Green), C.int(colour.Blue), C.int(colour.Alpha))
}

func (w *macosWindow) position() (int, int) {
	var x, y C.int
	var wg sync.WaitGroup
	wg.Add(1)
	go Dispatch(func() {
		C.windowGetPosition(w.nsWindow, &x, &y)
		wg.Done()
	})
	wg.Wait()
	return int(x), int(y)
}
