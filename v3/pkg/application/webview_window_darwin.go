//go:build darwin

package application

/*
#cgo CFLAGS: -mmacosx-version-min=10.13 -x objective-c
#cgo LDFLAGS: -framework Cocoa -framework WebKit

#include "application_darwin.h"
#include "webview_window_darwin.h"
#include <stdlib.h>
#include "Cocoa/Cocoa.h"
#import <WebKit/WebKit.h>
#import <AppKit/AppKit.h>
#import "webview_window_darwin_drag.h"

struct WebviewPreferences {
    bool *TabFocusesLinks;
    bool *TextInteractionEnabled;
    bool *FullscreenEnabled;
};

extern void registerListener(unsigned int event);

// Create a new Window
void* windowNew(unsigned int id, int width, int height, bool fraudulentWebsiteWarningEnabled, bool frameless, bool enableDragAndDrop, struct WebviewPreferences preferences) {
	NSWindowStyleMask styleMask = NSWindowStyleMaskTitled | NSWindowStyleMaskClosable | NSWindowStyleMaskMiniaturizable | NSWindowStyleMaskResizable;
	if (frameless) {
		styleMask = NSWindowStyleMaskBorderless | NSWindowStyleMaskResizable;
	}
	WebviewWindow* window = [[WebviewWindow alloc] initWithContentRect:NSMakeRect(0, 0, width-1, height-1)
		styleMask:styleMask
		backing:NSBackingStoreBuffered
		defer:NO];

	// Create delegate
	WebviewWindowDelegate* delegate = [[WebviewWindowDelegate alloc] init];
	[delegate autorelease];

	// Set delegate
	[window setDelegate:delegate];
	delegate.windowId = id;

	// Add NSView to window
	NSView* view = [[NSView alloc] initWithFrame:NSMakeRect(0, 0, width-1, height-1)];
	[view autorelease];

	[view setAutoresizingMask:NSViewWidthSizable | NSViewHeightSizable];
	if( frameless ) {
		[view setWantsLayer:YES];
		view.layer.cornerRadius = 8.0;
	}
	[window setContentView:view];

	// Embed wkwebview in window
	NSRect frame = NSMakeRect(0, 0, width, height);
	WKWebViewConfiguration* config = [[WKWebViewConfiguration alloc] init];
	[config autorelease];

	// Set preferences
    if (preferences.TabFocusesLinks != NULL) {
		config.preferences.tabFocusesLinks = *preferences.TabFocusesLinks;
	}

#if MAC_OS_X_VERSION_MAX_ALLOWED >= 110300
	if (@available(macOS 11.3, *)) {
		if (preferences.TextInteractionEnabled != NULL) {
			config.preferences.textInteractionEnabled = *preferences.TextInteractionEnabled;
		}
	}
#endif

#if MAC_OS_X_VERSION_MAX_ALLOWED >= 120300
	if (@available(macOS 12.3, *)) {
         if (preferences.FullscreenEnabled != NULL) {
             config.preferences.elementFullscreenEnabled = *preferences.FullscreenEnabled;
         }
     }
#endif

	config.suppressesIncrementalRendering = true;
    config.applicationNameForUserAgent = @"wails.io";
	[config setURLSchemeHandler:delegate forURLScheme:@"wails"];

#if MAC_OS_X_VERSION_MAX_ALLOWED >= 101500
 	if (@available(macOS 10.15, *)) {
         config.preferences.fraudulentWebsiteWarningEnabled = fraudulentWebsiteWarningEnabled;
	}
#endif

	// Setup user content controller
    WKUserContentController* userContentController = [WKUserContentController new];
	[userContentController autorelease];

    [userContentController addScriptMessageHandler:delegate name:@"external"];
    config.userContentController = userContentController;

	WKWebView* webView = [[WKWebView alloc] initWithFrame:frame configuration:config];
	[webView autorelease];

	[view addSubview:webView];

    // support webview events
    [webView setNavigationDelegate:delegate];

	// Ensure webview resizes with the window
	[webView setAutoresizingMask:NSViewWidthSizable | NSViewHeightSizable];

	if( enableDragAndDrop ) {
		WebviewDrag* dragView = [[WebviewDrag alloc] initWithFrame:NSMakeRect(0, 0, width-1, height-1)];
		[dragView autorelease];

		[view setAutoresizingMask:NSViewWidthSizable | NSViewHeightSizable];
		[view addSubview:dragView];
		dragView.windowId = id;
	}

	window.webView = webView;
	return window;
}


void printWindowStyle(void *window) {
	WebviewWindow* nsWindow = (WebviewWindow*)window;
    NSWindowStyleMask styleMask = [nsWindow styleMask];
	// Get delegate
	WebviewWindowDelegate* windowDelegate = (WebviewWindowDelegate*)[nsWindow delegate];

	printf("Window %d style mask: ", windowDelegate.windowId);

    if (styleMask & NSWindowStyleMaskTitled)
    {
        printf("NSWindowStyleMaskTitled ");
    }

    if (styleMask & NSWindowStyleMaskClosable)
    {
        printf("NSWindowStyleMaskClosable ");
    }

    if (styleMask & NSWindowStyleMaskMiniaturizable)
    {
        printf("NSWindowStyleMaskMiniaturizable ");
    }

    if (styleMask & NSWindowStyleMaskResizable)
    {
        printf("NSWindowStyleMaskResizable ");
    }

    if (styleMask & NSWindowStyleMaskFullSizeContentView)
    {
        printf("NSWindowStyleMaskFullSizeContentView ");
    }

    if (styleMask & NSWindowStyleMaskNonactivatingPanel)
    {
        printf("NSWindowStyleMaskNonactivatingPanel ");
    }

	if (styleMask & NSWindowStyleMaskFullScreen)
	{
		printf("NSWindowStyleMaskFullScreen ");
	}

	if (styleMask & NSWindowStyleMaskBorderless)
	{
		printf("MSWindowStyleMaskBorderless ");
	}

	printf("\n");
}


// setInvisibleTitleBarHeight sets the invisible title bar height
void setInvisibleTitleBarHeight(void* window, unsigned int height) {
	WebviewWindow* nsWindow = (WebviewWindow*)window;
	// Get delegate
	WebviewWindowDelegate* delegate = (WebviewWindowDelegate*)[nsWindow delegate];
	// Set height
	delegate.invisibleTitleBarHeight = height;
}

// Make NSWindow transparent
void windowSetTransparent(void* nsWindow) {
	[(WebviewWindow*)nsWindow setOpaque:NO];
	[(WebviewWindow*)nsWindow setBackgroundColor:[NSColor clearColor]];
}

void windowSetInvisibleTitleBar(void* nsWindow, unsigned int height) {
	WebviewWindow* window = (WebviewWindow*)nsWindow;
	WebviewWindowDelegate* delegate = (WebviewWindowDelegate*)[window delegate];
	delegate.invisibleTitleBarHeight = height;
}


// Set the title of the NSWindow
void windowSetTitle(void* nsWindow, char* title) {
	NSString* nsTitle = [NSString stringWithUTF8String:title];
	[(WebviewWindow*)nsWindow setTitle:nsTitle];
	free(title);
}

// Set the size of the NSWindow
void windowSetSize(void* nsWindow, int width, int height) {
	// Set window size on main thread
	WebviewWindow* window = (WebviewWindow*)nsWindow;
	NSSize contentSize = [window contentRectForFrameRect:NSMakeRect(0, 0, width, height)].size;
	[window setContentSize:contentSize];
	[window setFrame:NSMakeRect(window.frame.origin.x, window.frame.origin.y, width, height) display:YES animate:YES];
}

// Set NSWindow always on top
void windowSetAlwaysOnTop(void* nsWindow, bool alwaysOnTop) {
	// Set window always on top on main thread
	[(WebviewWindow*)nsWindow setLevel:alwaysOnTop ? NSFloatingWindowLevel : NSNormalWindowLevel];
}

void setNormalWindowLevel(void* nsWindow) { [(WebviewWindow*)nsWindow setLevel:NSNormalWindowLevel]; }
void setFloatingWindowLevel(void* nsWindow) { [(WebviewWindow*)nsWindow setLevel:NSFloatingWindowLevel];}
void setPopUpMenuWindowLevel(void* nsWindow) { [(WebviewWindow*)nsWindow setLevel:NSPopUpMenuWindowLevel]; }
void setMainMenuWindowLevel(void* nsWindow) { [(WebviewWindow*)nsWindow setLevel:NSMainMenuWindowLevel]; }
void setStatusWindowLevel(void* nsWindow) { [(WebviewWindow*)nsWindow setLevel:NSStatusWindowLevel]; }
void setModalPanelWindowLevel(void* nsWindow) { [(WebviewWindow*)nsWindow setLevel:NSModalPanelWindowLevel]; }
void setScreenSaverWindowLevel(void* nsWindow) { [(WebviewWindow*)nsWindow setLevel:NSScreenSaverWindowLevel]; }
void setTornOffMenuWindowLevel(void* nsWindow) { [(WebviewWindow*)nsWindow setLevel:NSTornOffMenuWindowLevel]; }

// Load URL in NSWindow
void navigationLoadURL(void* nsWindow, char* url) {
	// Load URL on main thread
	NSURL* nsURL = [NSURL URLWithString:[NSString stringWithUTF8String:url]];
	NSURLRequest* request = [NSURLRequest requestWithURL:nsURL];
	WebviewWindow* window = (WebviewWindow*)nsWindow;
	[window.webView loadRequest:request];
	free(url);
}

// Set NSWindow resizable
void windowSetResizable(void* nsWindow, bool resizable) {
	// Set window resizable on main thread
	WebviewWindow* window = (WebviewWindow*)nsWindow;
	if (resizable) {
		NSWindowStyleMask styleMask = [window styleMask] | NSWindowStyleMaskResizable;
		[window setStyleMask:styleMask];
	} else {
		NSWindowStyleMask styleMask = [window styleMask] & ~NSWindowStyleMaskResizable;
		[window setStyleMask:styleMask];
	}
}

// Set NSWindow min size
void windowSetMinSize(void* nsWindow, int width, int height) {
	// Set window min size on main thread
	WebviewWindow* window = (WebviewWindow*)nsWindow;
	NSSize contentSize = [window contentRectForFrameRect:NSMakeRect(0, 0, width, height)].size;
	[window setContentMinSize:contentSize];
	NSSize size = { width, height };
	[window setMinSize:size];
}

// Set NSWindow max size
void windowSetMaxSize(void* nsWindow, int width, int height) {
	// Set window max size on main thread
	NSSize size = { FLT_MAX, FLT_MAX };
	size.width = width > 0 ? width : FLT_MAX;
	size.height = height > 0 ? height : FLT_MAX;
	WebviewWindow* window = (WebviewWindow*)nsWindow;
	NSSize contentSize = [window contentRectForFrameRect:NSMakeRect(0, 0, size.width, size.height)].size;
	[window setContentMaxSize:contentSize];
	[window setMaxSize:size];
}

// windowZoomReset
void windowZoomReset(void* nsWindow) {
	WebviewWindow* window = (WebviewWindow*)nsWindow;
	[window.webView setMagnification:1.0];
}

// windowZoomSet
void windowZoomSet(void* nsWindow, double zoom) {
	WebviewWindow* window = (WebviewWindow*)nsWindow;
	// Reset zoom
	[window.webView setMagnification:zoom];
}

// windowZoomGet
float windowZoomGet(void* nsWindow) {
	WebviewWindow* window = (WebviewWindow*)nsWindow;
	// Get zoom
	return [window.webView magnification];
}

// windowZoomIn
void windowZoomIn(void* nsWindow) {
	WebviewWindow* window = (WebviewWindow*)nsWindow;
	// Zoom in
	[window.webView setMagnification:window.webView.magnification + 0.05];
}

// windowZoomOut
void windowZoomOut(void* nsWindow) {
	WebviewWindow* window = (WebviewWindow*)nsWindow;
	// Zoom out
	if( window.webView.magnification > 1.05 ) {
		[window.webView setMagnification:window.webView.magnification - 0.05];
	} else {
		[window.webView setMagnification:1.0];
	}
}

// set the window position relative to the screen
void windowSetRelativePosition(void* nsWindow, int x, int y) {
	WebviewWindow* window = (WebviewWindow*)nsWindow;
	NSScreen* screen = [window screen];
	if( screen == NULL ) {
		screen = [NSScreen mainScreen];
	}
	NSRect windowFrame = [window frame];
	NSRect screenFrame = [screen frame];
	windowFrame.origin.x = screenFrame.origin.x + (float)x;
	windowFrame.origin.y = (screenFrame.origin.y + screenFrame.size.height) - windowFrame.size.height - (float)y;

	[window setFrame:windowFrame display:TRUE animate:FALSE];
}

// Execute JS in NSWindow
void windowExecJS(void* nsWindow, const char* js) {
	WebviewWindow* window = (WebviewWindow*)nsWindow;
	[window.webView evaluateJavaScript:[NSString stringWithUTF8String:js] completionHandler:nil];
	free((void*)js);
}

// Make NSWindow backdrop translucent
void windowSetTranslucent(void* nsWindow) {
	// Get window
	WebviewWindow* window = (WebviewWindow*)nsWindow;

	id contentView = [window contentView];
	NSVisualEffectView *effectView = [NSVisualEffectView alloc];
	NSRect bounds = [contentView bounds];
	[effectView initWithFrame:bounds];
	[effectView setAutoresizingMask:NSViewWidthSizable | NSViewHeightSizable];
	[effectView setBlendingMode:NSVisualEffectBlendingModeBehindWindow];
	[effectView setState:NSVisualEffectStateActive];
	[contentView addSubview:effectView positioned:NSWindowBelow relativeTo:nil];
}

// Make webview background transparent
void webviewSetTransparent(void* nsWindow) {
	WebviewWindow* window = (WebviewWindow*)nsWindow;
	// Set webview background transparent
	[window.webView setValue:@NO forKey:@"drawsBackground"];
}

// Set webview background colour
void webviewSetBackgroundColour(void* nsWindow, int r, int g, int b, int alpha) {
	WebviewWindow* window = (WebviewWindow*)nsWindow;
	// Set webview background color
	[window.webView setValue:[NSColor colorWithRed:r/255.0 green:g/255.0 blue:b/255.0 alpha:alpha/255.0] forKey:@"backgroundColor"];
}

// Set the window background colour
void windowSetBackgroundColour(void* nsWindow, int r, int g, int b, int alpha) {
	[(WebviewWindow*)nsWindow setBackgroundColor:[NSColor colorWithRed:r/255.0 green:g/255.0 blue:b/255.0 alpha:alpha/255.0]];
}

bool windowIsMaximised(void* nsWindow) {
	return [(WebviewWindow*)nsWindow isZoomed];
}

bool windowIsFullscreen(void* nsWindow) {
	return [(WebviewWindow*)nsWindow styleMask] & NSWindowStyleMaskFullScreen;
}

bool windowIsMinimised(void* nsWindow) {
	return [(WebviewWindow*)nsWindow isMiniaturized];
}

bool windowIsFocused(void* nsWindow) {
	return [(WebviewWindow*)nsWindow isKeyWindow];
}

// Set Window fullscreen
void windowFullscreen(void* nsWindow) {
	if( windowIsFullscreen(nsWindow) ) {
		return;
	}
	dispatch_async(dispatch_get_main_queue(), ^{
		[(WebviewWindow*)nsWindow toggleFullScreen:nil];
	});}

void windowUnFullscreen(void* nsWindow) {
	if( !windowIsFullscreen(nsWindow) ) {
		return;
	}
	dispatch_async(dispatch_get_main_queue(), ^{
		[(WebviewWindow*)nsWindow toggleFullScreen:nil];
	});
}

// restore window to normal size
void windowRestore(void* nsWindow) {
	// If window is fullscreen
	if([(WebviewWindow*)nsWindow styleMask] & NSWindowStyleMaskFullScreen) {
		[(WebviewWindow*)nsWindow toggleFullScreen:nil];
	}
	// If window is maximised
	if([(WebviewWindow*)nsWindow isZoomed]) {
		[(WebviewWindow*)nsWindow zoom:nil];
	}
	// If window in minimised
	if([(WebviewWindow*)nsWindow isMiniaturized]) {
		[(WebviewWindow*)nsWindow deminiaturize:nil];
	}
}

// disable window fullscreen button
void setFullscreenButtonEnabled(void* nsWindow, bool enabled) {
	NSButton *fullscreenButton = [(WebviewWindow*)nsWindow standardWindowButton:NSWindowZoomButton];
	fullscreenButton.enabled = enabled;
}

// Set the titlebar style
void windowSetTitleBarAppearsTransparent(void* nsWindow, bool transparent) {
	if( transparent ) {
		[(WebviewWindow*)nsWindow setTitlebarAppearsTransparent:true];
	} else {
		[(WebviewWindow*)nsWindow setTitlebarAppearsTransparent:false];
	}
}

// Set window fullsize content view
void windowSetFullSizeContent(void* nsWindow, bool fullSize) {
	if( fullSize ) {
		[(WebviewWindow*)nsWindow setStyleMask:[(WebviewWindow*)nsWindow styleMask] | NSWindowStyleMaskFullSizeContentView];
	} else {
		[(WebviewWindow*)nsWindow setStyleMask:[(WebviewWindow*)nsWindow styleMask] & ~NSWindowStyleMaskFullSizeContentView];
	}
}

// Set Hide Titlebar
void windowSetHideTitleBar(void* nsWindow, bool hideTitlebar) {
	if( hideTitlebar ) {
		[(WebviewWindow*)nsWindow setStyleMask:[(WebviewWindow*)nsWindow styleMask] & ~NSWindowStyleMaskTitled];
	} else {
		[(WebviewWindow*)nsWindow setStyleMask:[(WebviewWindow*)nsWindow styleMask] | NSWindowStyleMaskTitled];
	}
}

// Set Hide Title in Titlebar
void windowSetHideTitle(void* nsWindow, bool hideTitle) {
	if( hideTitle ) {
		[(WebviewWindow*)nsWindow setTitleVisibility:NSWindowTitleHidden];
	} else {
		[(WebviewWindow*)nsWindow setTitleVisibility:NSWindowTitleVisible];
	}
}

// Set Window use toolbar
void windowSetUseToolbar(void* nsWindow, bool useToolbar) {
	WebviewWindow* window = (WebviewWindow*)nsWindow;
	if( useToolbar ) {
		NSToolbar *toolbar = [[NSToolbar alloc] initWithIdentifier:@"wails.toolbar"];
		[toolbar autorelease];
		[window setToolbar:toolbar];
	} else {
		[window setToolbar:nil];
	}
}

// Set window toolbar style
void windowSetToolbarStyle(void* nsWindow, int style) {
	WebviewWindow* window = (WebviewWindow*)nsWindow;

#if MAC_OS_X_VERSION_MAX_ALLOWED >= 110000
	if (@available(macOS 11.0, *)) {
		NSToolbar* toolbar = [window toolbar];
		if ( toolbar == nil ) {
			return;
		}
		[window setToolbarStyle:style];
	}
#endif

}
// Set Hide Toolbar Separator
void windowSetHideToolbarSeparator(void* nsWindow, bool hideSeparator) {
	NSToolbar* toolbar = [(WebviewWindow*)nsWindow toolbar];
	if( toolbar == nil ) {
		return;
	}
	[toolbar setShowsBaselineSeparator:!hideSeparator];
}

// Configure the toolbar auto-hide feature
void windowSetShowToolbarWhenFullscreen(void* window, bool setting) {
	WebviewWindow* nsWindow = (WebviewWindow*)window;
	// Get delegate
	WebviewWindowDelegate* delegate = (WebviewWindowDelegate*)[nsWindow delegate];
	// Set height
	delegate.showToolbarWhenFullscreen = setting;
}

// Set Window appearance type
void windowSetAppearanceTypeByName(void* nsWindow, const char *appearanceName) {
	// set window appearance type by name
	// Convert appearance name to NSString
	NSString* appearanceNameString = [NSString stringWithUTF8String:appearanceName];
	// Set appearance
	[(WebviewWindow*)nsWindow setAppearance:[NSAppearance appearanceNamed:appearanceNameString]];

	free((void*)appearanceName);
}

// Center window on current monitor
void windowCenter(void* nsWindow) {
	[(WebviewWindow*)nsWindow center];
}

// Get the current size of the window
void windowGetSize(void* nsWindow, int* width, int* height) {
	NSRect frame = [(WebviewWindow*)nsWindow frame];
	*width = frame.size.width;
	*height = frame.size.height;
}

// Get window width
int windowGetWidth(void* nsWindow) {
	return [(WebviewWindow*)nsWindow frame].size.width;
}

// Get window height
int windowGetHeight(void* nsWindow) {
	return [(WebviewWindow*)nsWindow frame].size.height;
}

// Get window position
void windowGetRelativePosition(void* nsWindow, int* x, int* y) {
	WebviewWindow* window = (WebviewWindow*)nsWindow;
	NSRect frame = [window frame];
	*x = frame.origin.x;

	// Translate to screen coordinates so Y=0 is the top of the screen
	NSScreen* screen = [window screen];
	if( screen == NULL ) {
		screen = [NSScreen mainScreen];
	}
	NSRect screenFrame = [screen frame];
	*y = screenFrame.size.height - frame.origin.y - frame.size.height;
}

// Get absolute window position
void windowGetPosition(void* nsWindow, int* x, int* y) {
	NSRect frame = [(WebviewWindow*)nsWindow frame];
	*x = frame.origin.x;
	*y = frame.origin.y;
}

void windowSetPosition(void* nsWindow, int x, int y) {
	NSRect frame = [(WebviewWindow*)nsWindow frame];
	frame.origin.x = x;
	frame.origin.y = y;
	[(WebviewWindow*)nsWindow setFrame:frame display:YES];
}

// Destroy window
void windowDestroy(void* nsWindow) {
	[(WebviewWindow*)nsWindow close];
}

// Remove drop shadow from window
void windowSetShadow(void* nsWindow, bool hasShadow) {
	[(WebviewWindow*)nsWindow setHasShadow:hasShadow];
}


// windowClose closes the current window
static void windowClose(void *window) {
	[(WebviewWindow*)window close];
}

// windowZoom
static void windowZoom(void *window) {
	[(WebviewWindow*)window zoom:nil];
}

// webviewRenderHTML renders the given HTML
static void windowRenderHTML(void *window, const char *html) {
	WebviewWindow* nsWindow = (WebviewWindow*)window;
	// get window delegate
	WebviewWindowDelegate* windowDelegate = (WebviewWindowDelegate*)[nsWindow delegate];
	// render html
	[nsWindow.webView loadHTMLString:[NSString stringWithUTF8String:html] baseURL:nil];
}

static void windowInjectCSS(void *window, const char *css) {
	WebviewWindow* nsWindow = (WebviewWindow*)window;
	// inject css
	[nsWindow.webView evaluateJavaScript:[NSString stringWithFormat:@"(function() { var style = document.createElement('style'); style.appendChild(document.createTextNode('%@')); document.head.appendChild(style); })();", [NSString stringWithUTF8String:css]] completionHandler:nil];
	free((void*)css);
}

static void windowMinimise(void *window) {
	[(WebviewWindow*)window miniaturize:nil];
}

// zoom maximizes the window to the screen dimensions
static void windowMaximise(void *window) {
	[(WebviewWindow*)window zoom:nil];
}

static bool isFullScreen(void *window) {
	WebviewWindow* nsWindow = (WebviewWindow*)window;
    long mask = [nsWindow styleMask];
    return (mask & NSWindowStyleMaskFullScreen) == NSWindowStyleMaskFullScreen;
}

static bool isVisible(void *window) {
	WebviewWindow* nsWindow = (WebviewWindow*)window;
    return (nsWindow.occlusionState & NSWindowOcclusionStateVisible) == NSWindowOcclusionStateVisible;
}

// windowSetFullScreen
static void windowSetFullScreen(void *window, bool fullscreen) {
	if (isFullScreen(window)) {
		return;
	}
	WebviewWindow* nsWindow = (WebviewWindow*)window;
	windowSetMaxSize(nsWindow, 0, 0);
	windowSetMinSize(nsWindow, 0, 0);
	[nsWindow toggleFullScreen:nil];
}

// windowUnminimise
static void windowUnminimise(void *window) {
	[(WebviewWindow*)window deminiaturize:nil];
}

// windowUnmaximise
static void windowUnmaximise(void *window) {
	[(WebviewWindow*)window zoom:nil];
}

static void windowDisableSizeConstraints(void *window) {
	WebviewWindow* nsWindow = (WebviewWindow*)window;
	// disable size constraints
	[nsWindow setContentMinSize:CGSizeZero];
	[nsWindow setContentMaxSize:CGSizeZero];
}

static void windowShow(void *window) {
	[(WebviewWindow*)window makeKeyAndOrderFront:nil];
}

static void windowHide(void *window) {
	[(WebviewWindow*)window orderOut:nil];
}

// setButtonState sets the state of the given button
// 0 = enabled
// 1 = disabled
// 2 = hidden
static void setButtonState(void *button, int state) {
	if (button == nil) {
		return;
	}
	NSButton *nsbutton = (NSButton*)button;
	nsbutton.hidden = state == 2;
	nsbutton.enabled = state != 1;
}

// setMinimiseButtonState sets the minimise button state
static void setMinimiseButtonState(void *window, int state) {
	WebviewWindow* nsWindow = (WebviewWindow*)window;
	NSButton *minimiseButton = [nsWindow standardWindowButton:NSWindowMiniaturizeButton];
	setButtonState(minimiseButton, state);
}

// setMaximiseButtonState sets the maximise button state
static void setMaximiseButtonState(void *window, int state) {
	WebviewWindow* nsWindow = (WebviewWindow*)window;
	NSButton *maximiseButton = [nsWindow standardWindowButton:NSWindowZoomButton];
	setButtonState(maximiseButton, state);
}

// setCloseButtonState sets the close button state
static void setCloseButtonState(void *window, int state) {
	WebviewWindow* nsWindow = (WebviewWindow*)window;
	NSButton *closeButton = [nsWindow standardWindowButton:NSWindowCloseButton];
	setButtonState(closeButton, state);
}

// windowShowMenu opens an NSMenu at the given coordinates
static void windowShowMenu(void *window, void *menu, int x, int y) {
	NSMenu* nsMenu = (NSMenu*)menu;
	WKWebView* webView = ((WebviewWindow*)window).webView;
	NSPoint point = NSMakePoint(x, y);
	[nsMenu popUpMenuPositioningItem:nil atLocation:point inView:webView];
}

// Make the given window frameless
static void windowSetFrameless(void *window, bool frameless) {
	WebviewWindow* nsWindow = (WebviewWindow*)window;
	// set the window style to be frameless
	if (frameless) {
		[nsWindow setStyleMask:([nsWindow styleMask] | NSWindowStyleMaskFullSizeContentView)];
	} else {
		[nsWindow setStyleMask:([nsWindow styleMask] & ~NSWindowStyleMaskFullSizeContentView)];
	}
}

static void startDrag(void *window) {
	WebviewWindow* nsWindow = (WebviewWindow*)window;

	// Get delegate
	WebviewWindowDelegate* windowDelegate = (WebviewWindowDelegate*)[nsWindow delegate];

	// start drag
	[windowDelegate startDrag:nsWindow];
}

// Credit: https://stackoverflow.com/q/33319295
static void windowPrint(void *window) {
#if MAC_OS_X_VERSION_MAX_ALLOWED >= 110000
	// Check if macOS 11.0 or newer
	if (@available(macOS 11.0, *)) {
		WebviewWindow* nsWindow = (WebviewWindow*)window;
		WebviewWindowDelegate* windowDelegate = (WebviewWindowDelegate*)[nsWindow delegate];
		WKWebView* webView = nsWindow.webView;

		// TODO: Think about whether to expose this as config
		NSPrintInfo *pInfo = [NSPrintInfo sharedPrintInfo];
		pInfo.horizontalPagination = NSPrintingPaginationModeAutomatic;
		pInfo.verticalPagination = NSPrintingPaginationModeAutomatic;
		pInfo.verticallyCentered = YES;
		pInfo.horizontallyCentered = YES;
		pInfo.orientation = NSPaperOrientationLandscape;
		pInfo.leftMargin = 30;
		pInfo.rightMargin = 30;
		pInfo.topMargin = 30;
		pInfo.bottomMargin = 30;

		NSPrintOperation *po = [webView printOperationWithPrintInfo:pInfo];
		po.showsPrintPanel = YES;
		po.showsProgressPanel = YES;

		// Without the next line you get an exception. Also it seems to
		// completely ignore the values in the rect. I tried changing them
		// in both x and y direction to include content scrolled off screen.
		// It had no effect whatsoever in either direction.
		po.view.frame = webView.bounds;

		// [printOperation runOperation] DOES NOT WORK WITH WKWEBVIEW, use
		[po runOperationModalForWindow:window delegate:windowDelegate didRunSelector:nil contextInfo:nil];
	}
#endif
}

void setWindowEnabled(void *window, bool enabled) {
	WebviewWindow* nsWindow = (WebviewWindow*)window;
	[nsWindow setIgnoresMouseEvents:!enabled];
}

void windowSetEnabled(void *window, bool enabled) {
	// TODO: Implement
}

void windowFocus(void *window) {
	WebviewWindow* nsWindow = (WebviewWindow*)window;
	// If the current application is not active, activate it
	if (![[NSApplication sharedApplication] isActive]) {
		[[NSApplication sharedApplication] activateIgnoringOtherApps:YES];
	}
	[nsWindow makeKeyAndOrderFront:nil];
	[nsWindow makeKeyWindow];
}

static bool isIgnoreMouseEvents(void *nsWindow) {
    NSWindow *window = (__bridge NSWindow *)nsWindow;
    return [window ignoresMouseEvents];
}

static void setIgnoreMouseEvents(void *nsWindow, bool ignore) {
    NSWindow *window = (__bridge NSWindow *)nsWindow;
    [window setIgnoresMouseEvents:ignore];
}

*/
import "C"
import (
	"sync"
	"unsafe"

	"github.com/wailsapp/wails/v3/internal/assetserver"
	"github.com/wailsapp/wails/v3/internal/runtime"

	"github.com/wailsapp/wails/v3/pkg/events"
)

type macosWebviewWindow struct {
	nsWindow unsafe.Pointer
	parent   *WebviewWindow
}

func (w *macosWebviewWindow) handleKeyEvent(acceleratorString string) {
	// Parse acceleratorString
	accelerator, err := parseAccelerator(acceleratorString)
	if err != nil {
		globalApplication.error("unable to parse accelerator: %s", err.Error())
		return
	}
	w.parent.processKeyBinding(accelerator.String())
}

func (w *macosWebviewWindow) getBorderSizes() *LRTB {
	return &LRTB{}
}

func (w *macosWebviewWindow) isFocused() bool {
	return bool(C.windowIsFocused(w.nsWindow))
}

func (w *macosWebviewWindow) setPosition(x int, y int) {
	C.windowSetPosition(w.nsWindow, C.int(x), C.int(y))
}

func (w *macosWebviewWindow) print() error {
	C.windowPrint(w.nsWindow)
	return nil
}

func (w *macosWebviewWindow) startResize(_ string) error {
	// Never called. Handled natively by the OS.
	return nil
}

func (w *macosWebviewWindow) focus() {
	// Make the window key and main
	C.windowFocus(w.nsWindow)
}

func (w *macosWebviewWindow) openContextMenu(menu *Menu, data *ContextMenuData) {
	// Create the menu
	thisMenu := newMenuImpl(menu)
	thisMenu.update()
	C.windowShowMenu(w.nsWindow, thisMenu.nsMenu, C.int(data.X), C.int(data.Y))
}

func (w *macosWebviewWindow) getZoom() float64 {
	return float64(C.windowZoomGet(w.nsWindow))
}

func (w *macosWebviewWindow) setZoom(zoom float64) {
	C.windowZoomSet(w.nsWindow, C.double(zoom))
}

func (w *macosWebviewWindow) setFrameless(frameless bool) {
	C.windowSetFrameless(w.nsWindow, C.bool(frameless))
	if frameless {
		C.windowSetTitleBarAppearsTransparent(w.nsWindow, C.bool(true))
		C.windowSetHideTitle(w.nsWindow, C.bool(true))
	} else {
		macOptions := w.parent.options.Mac
		appearsTransparent := macOptions.TitleBar.AppearsTransparent
		hideTitle := macOptions.TitleBar.HideTitle
		C.windowSetTitleBarAppearsTransparent(w.nsWindow, C.bool(appearsTransparent))
		C.windowSetHideTitle(w.nsWindow, C.bool(hideTitle))
	}
}

func (w *macosWebviewWindow) setHasShadow(hasShadow bool) {
	C.windowSetShadow(w.nsWindow, C.bool(hasShadow))
}

func (w *macosWebviewWindow) getScreen() (*Screen, error) {
	return getScreenForWindow(w)
}

func (w *macosWebviewWindow) show() {
	C.windowShow(w.nsWindow)
}

func (w *macosWebviewWindow) hide() {
	C.windowHide(w.nsWindow)
}

func (w *macosWebviewWindow) setFullscreenButtonEnabled(enabled bool) {
	C.setFullscreenButtonEnabled(w.nsWindow, C.bool(enabled))
}

func (w *macosWebviewWindow) disableSizeConstraints() {
	C.windowDisableSizeConstraints(w.nsWindow)
}

func (w *macosWebviewWindow) unfullscreen() {
	C.windowUnFullscreen(w.nsWindow)
}

func (w *macosWebviewWindow) fullscreen() {
	C.windowFullscreen(w.nsWindow)
}

func (w *macosWebviewWindow) unminimise() {
	C.windowUnminimise(w.nsWindow)
}

func (w *macosWebviewWindow) unmaximise() {
	C.windowUnmaximise(w.nsWindow)
}

func (w *macosWebviewWindow) maximise() {
	C.windowMaximise(w.nsWindow)
}

func (w *macosWebviewWindow) minimise() {
	C.windowMinimise(w.nsWindow)
}

func (w *macosWebviewWindow) on(eventID uint) {
	//C.registerListener(C.uint(eventID))
}

func (w *macosWebviewWindow) zoom() {
	C.windowZoom(w.nsWindow)
}

func (w *macosWebviewWindow) windowZoom() {
	C.windowZoom(w.nsWindow)
}

func (w *macosWebviewWindow) close() {
	C.windowClose(w.nsWindow)
}

func (w *macosWebviewWindow) zoomIn() {
	C.windowZoomIn(w.nsWindow)
}

func (w *macosWebviewWindow) zoomOut() {
	C.windowZoomOut(w.nsWindow)
}

func (w *macosWebviewWindow) zoomReset() {
	C.windowZoomReset(w.nsWindow)
}

func (w *macosWebviewWindow) reload() {
	//TODO: Implement
	globalApplication.debug("reload called on WebviewWindow", "parentID", w.parent.id)
}

func (w *macosWebviewWindow) forceReload() {
	//TODO: Implement
	globalApplication.debug("force reload called on WebviewWindow", "parentID", w.parent.id)
}

func (w *macosWebviewWindow) center() {
	C.windowCenter(w.nsWindow)
}

func (w *macosWebviewWindow) isMinimised() bool {
	return w.syncMainThreadReturningBool(func() bool {
		return bool(C.windowIsMinimised(w.nsWindow))
	})
}

func (w *macosWebviewWindow) isMaximised() bool {
	return w.syncMainThreadReturningBool(func() bool {
		return bool(C.windowIsMaximised(w.nsWindow))
	})
}

func (w *macosWebviewWindow) isFullscreen() bool {
	return w.syncMainThreadReturningBool(func() bool {
		return bool(C.windowIsFullscreen(w.nsWindow))
	})
}

func (w *macosWebviewWindow) isNormal() bool {
	return !w.isMinimised() && !w.isMaximised() && !w.isFullscreen()
}

func (w *macosWebviewWindow) isVisible() bool {
	return bool(C.isVisible(w.nsWindow))
}

func (w *macosWebviewWindow) syncMainThreadReturningBool(fn func() bool) bool {
	var wg sync.WaitGroup
	wg.Add(1)
	var result bool
	globalApplication.dispatchOnMainThread(func() {
		result = fn()
		wg.Done()
	})
	wg.Wait()
	return result
}

func (w *macosWebviewWindow) restore() {
	// restore window to normal size
	C.windowRestore(w.nsWindow)
}

func (w *macosWebviewWindow) restoreWindow() {
	C.windowRestore(w.nsWindow)
}

func (w *macosWebviewWindow) setEnabled(enabled bool) {
	C.windowSetEnabled(w.nsWindow, C.bool(enabled))
}

func (w *macosWebviewWindow) execJS(js string) {
	InvokeAsync(func() {
		if globalApplication.performingShutdown {
			return
		}
		if w.nsWindow == nil {
			return
		}
		C.windowExecJS(w.nsWindow, C.CString(js))
	})
}

func (w *macosWebviewWindow) setURL(uri string) {
	C.navigationLoadURL(w.nsWindow, C.CString(uri))
}

func (w *macosWebviewWindow) setAlwaysOnTop(alwaysOnTop bool) {
	C.windowSetAlwaysOnTop(w.nsWindow, C.bool(alwaysOnTop))
}

func newWindowImpl(parent *WebviewWindow) *macosWebviewWindow {
	result := &macosWebviewWindow{
		parent: parent,
	}
	result.parent.RegisterHook(events.Mac.WebViewDidFinishNavigation, func(event *WindowEvent) {
		result.execJS(runtime.Core())
	})
	return result
}

func (w *macosWebviewWindow) setTitle(title string) {
	if !w.parent.options.Frameless {
		cTitle := C.CString(title)
		C.windowSetTitle(w.nsWindow, cTitle)
	}
}

func (w *macosWebviewWindow) flash(_ bool) {
	// Not supported on macOS
}

func (w *macosWebviewWindow) setSize(width, height int) {
	C.windowSetSize(w.nsWindow, C.int(width), C.int(height))
}

func (w *macosWebviewWindow) setMinSize(width, height int) {
	C.windowSetMinSize(w.nsWindow, C.int(width), C.int(height))
}
func (w *macosWebviewWindow) setMaxSize(width, height int) {
	C.windowSetMaxSize(w.nsWindow, C.int(width), C.int(height))
}

func (w *macosWebviewWindow) setResizable(resizable bool) {
	C.windowSetResizable(w.nsWindow, C.bool(resizable))
}

func (w *macosWebviewWindow) size() (int, int) {
	var width, height C.int
	var wg sync.WaitGroup
	wg.Add(1)
	globalApplication.dispatchOnMainThread(func() {
		C.windowGetSize(w.nsWindow, &width, &height)
		wg.Done()
	})
	wg.Wait()
	return int(width), int(height)
}

func (w *macosWebviewWindow) setRelativePosition(x, y int) {
	C.windowSetRelativePosition(w.nsWindow, C.int(x), C.int(y))
}

func (w *macosWebviewWindow) setWindowLevel(level MacWindowLevel) {
	switch level {
	case MacWindowLevelNormal:
		C.setNormalWindowLevel(w.nsWindow)
	case MacWindowLevelFloating:
		C.setFloatingWindowLevel(w.nsWindow)
	case MacWindowLevelTornOffMenu:
		C.setTornOffMenuWindowLevel(w.nsWindow)
	case MacWindowLevelModalPanel:
		C.setModalPanelWindowLevel(w.nsWindow)
	case MacWindowLevelMainMenu:
		C.setMainMenuWindowLevel(w.nsWindow)
	case MacWindowLevelStatus:
		C.setStatusWindowLevel(w.nsWindow)
	case MacWindowLevelPopUpMenu:
		C.setPopUpMenuWindowLevel(w.nsWindow)
	case MacWindowLevelScreenSaver:
		C.setScreenSaverWindowLevel(w.nsWindow)
	}
}

func (w *macosWebviewWindow) width() int {
	var width C.int
	var wg sync.WaitGroup
	wg.Add(1)
	globalApplication.dispatchOnMainThread(func() {
		width = C.windowGetWidth(w.nsWindow)
		wg.Done()
	})
	wg.Wait()
	return int(width)
}
func (w *macosWebviewWindow) height() int {
	var height C.int
	var wg sync.WaitGroup
	wg.Add(1)
	globalApplication.dispatchOnMainThread(func() {
		height = C.windowGetHeight(w.nsWindow)
		wg.Done()
	})
	wg.Wait()
	return int(height)
}

func bool2CboolPtr(value bool) *C.bool {
	v := C.bool(value)
	return &v
}

func (w *macosWebviewWindow) getWebviewPreferences() C.struct_WebviewPreferences {
	wvprefs := w.parent.options.Mac.WebviewPreferences

	var result C.struct_WebviewPreferences

	if wvprefs.TextInteractionEnabled.IsSet() {
		result.TextInteractionEnabled = bool2CboolPtr(wvprefs.TextInteractionEnabled.Get())
	}
	if wvprefs.TabFocusesLinks.IsSet() {
		result.TabFocusesLinks = bool2CboolPtr(wvprefs.TabFocusesLinks.Get())
	}
	if wvprefs.FullscreenEnabled.IsSet() {
		result.FullscreenEnabled = bool2CboolPtr(wvprefs.FullscreenEnabled.Get())
	}

	return result
}

func (w *macosWebviewWindow) run() {
	for eventId := range w.parent.eventListeners {
		w.on(eventId)
	}
	globalApplication.dispatchOnMainThread(func() {
		options := w.parent.options
		macOptions := options.Mac

		w.nsWindow = C.windowNew(C.uint(w.parent.id),
			C.int(options.Width),
			C.int(options.Height),
			C.bool(macOptions.EnableFraudulentWebsiteWarnings),
			C.bool(options.Frameless),
			C.bool(options.EnableDragAndDrop),
			w.getWebviewPreferences(),
		)
		w.setTitle(options.Title)
		w.setAlwaysOnTop(options.AlwaysOnTop)
		w.setResizable(!options.DisableResize)
		if options.MinWidth != 0 || options.MinHeight != 0 {
			w.setMinSize(options.MinWidth, options.MinHeight)
		}
		if options.MaxWidth != 0 || options.MaxHeight != 0 {
			w.setMaxSize(options.MaxWidth, options.MaxHeight)
		}
		//w.setZoom(options.Zoom)
		w.enableDevTools()

		w.setBackgroundColour(options.BackgroundColour)

		switch macOptions.Backdrop {
		case MacBackdropTransparent:
			C.windowSetTransparent(w.nsWindow)
			C.webviewSetTransparent(w.nsWindow)
		case MacBackdropTranslucent:
			C.windowSetTranslucent(w.nsWindow)
			C.webviewSetTransparent(w.nsWindow)
		case MacBackdropNormal:
		}

		if macOptions.WindowLevel == "" {
			macOptions.WindowLevel = MacWindowLevelNormal
		}
		w.setWindowLevel(macOptions.WindowLevel)

		// Initialise the window buttons
		w.setMinimiseButtonState(options.MinimiseButtonState)
		w.setMaximiseButtonState(options.MaximiseButtonState)
		w.setCloseButtonState(options.CloseButtonState)

		// Ignore mouse events if requested
		w.setIgnoreMouseEvents(options.IgnoreMouseEvents)

		titleBarOptions := macOptions.TitleBar
		if !w.parent.options.Frameless {
			C.windowSetTitleBarAppearsTransparent(w.nsWindow, C.bool(titleBarOptions.AppearsTransparent))
			C.windowSetHideTitleBar(w.nsWindow, C.bool(titleBarOptions.Hide))
			C.windowSetHideTitle(w.nsWindow, C.bool(titleBarOptions.HideTitle))
			C.windowSetFullSizeContent(w.nsWindow, C.bool(titleBarOptions.FullSizeContent))
			C.windowSetUseToolbar(w.nsWindow, C.bool(titleBarOptions.UseToolbar))
			C.windowSetToolbarStyle(w.nsWindow, C.int(titleBarOptions.ToolbarStyle))
			C.windowSetShowToolbarWhenFullscreen(w.nsWindow, C.bool(titleBarOptions.ShowToolbarWhenFullscreen))
			C.windowSetHideToolbarSeparator(w.nsWindow, C.bool(titleBarOptions.HideToolbarSeparator))
		}

		if macOptions.Appearance != "" {
			C.windowSetAppearanceTypeByName(w.nsWindow, C.CString(string(macOptions.Appearance)))
		}

		if macOptions.InvisibleTitleBarHeight != 0 {
			C.windowSetInvisibleTitleBar(w.nsWindow, C.uint(macOptions.InvisibleTitleBarHeight))
		}

		switch w.parent.options.StartState {
		case WindowStateMaximised:
			w.maximise()
		case WindowStateMinimised:
			w.minimise()
		case WindowStateFullscreen:
			w.fullscreen()
		case WindowStateNormal:
		}
		C.windowCenter(w.nsWindow)

		startURL, err := assetserver.GetStartURL(options.URL)
		if err != nil {
			globalApplication.fatal(err.Error())
		}

		w.setURL(startURL)

		// We need to wait for the HTML to load before we can execute the javascript
		w.parent.OnWindowEvent(events.Mac.WebViewDidFinishNavigation, func(_ *WindowEvent) {
			InvokeAsync(func() {
				if options.JS != "" {
					w.execJS(options.JS)
				}
				if options.CSS != "" {
					C.windowInjectCSS(w.nsWindow, C.CString(options.CSS))
				}
				if !options.Hidden {
					C.windowShow(w.nsWindow)
					w.setHasShadow(!options.Mac.DisableShadow)
				} else {
					// We have to wait until the window is shown before we can remove the shadow
					var cancel func()
					cancel = w.parent.OnWindowEvent(events.Mac.WindowDidBecomeKey, func(_ *WindowEvent) {
						w.setHasShadow(!options.Mac.DisableShadow)
						cancel()
					})
				}
			})
		})

		// Translate ShouldClose to common WindowClosing event
		w.parent.OnWindowEvent(events.Mac.WindowShouldClose, func(_ *WindowEvent) {
			w.parent.emit(events.Common.WindowClosing)
		})

		// Translate WindowDidResignKey to common WindowLostFocus event
		w.parent.OnWindowEvent(events.Mac.WindowDidResignKey, func(_ *WindowEvent) {
			w.parent.emit(events.Common.WindowLostFocus)
		})
		w.parent.OnWindowEvent(events.Mac.WindowDidResignMain, func(_ *WindowEvent) {
			w.parent.emit(events.Common.WindowLostFocus)
		})
		w.parent.OnWindowEvent(events.Mac.WindowDidResize, func(_ *WindowEvent) {
			w.parent.emit(events.Common.WindowDidResize)
		})

		if options.HTML != "" {
			w.setHTML(options.HTML)
		}

	})
}

func (w *macosWebviewWindow) nativeWindowHandle() uintptr {
	return uintptr(w.nsWindow)
}

func (w *macosWebviewWindow) setBackgroundColour(colour RGBA) {

	C.windowSetBackgroundColour(w.nsWindow, C.int(colour.Red), C.int(colour.Green), C.int(colour.Blue), C.int(colour.Alpha))
}

func (w *macosWebviewWindow) relativePosition() (int, int) {
	var x, y C.int
	InvokeSync(func() {
		C.windowGetRelativePosition(w.nsWindow, &x, &y)
	})

	return int(x), int(y)
}

func (w *macosWebviewWindow) position() (int, int) {
	var x, y C.int
	InvokeSync(func() {
		C.windowGetPosition(w.nsWindow, &x, &y)
	})

	return int(x), int(y)
}

func (w *macosWebviewWindow) bounds() Rect {
	// DOTO: do it in a single step + proper DPI scaling
	var x, y, width, height C.int
	InvokeSync(func() {
		C.windowGetPosition(w.nsWindow, &x, &y)
		C.windowGetSize(w.nsWindow, &width, &height)
	})

	return Rect{
		X:      int(x),
		Y:      int(y),
		Width:  int(width),
		Height: int(height),
	}
}

func (w *macosWebviewWindow) setBounds(bounds Rect) {
	// DOTO: do it in a single step + proper DPI scaling
	C.windowSetPosition(w.nsWindow, C.int(bounds.X), C.int(bounds.Y))
	C.windowSetSize(w.nsWindow, C.int(bounds.Width), C.int(bounds.Height))
}

func (w *macosWebviewWindow) physicalBounds() Rect {
	// TODO: proper DPI scaling
	return w.bounds()
}

func (w *macosWebviewWindow) setPhysicalBounds(physicalBounds Rect) {
	// TODO: proper DPI scaling
	w.setBounds(physicalBounds)
}

func (w *macosWebviewWindow) destroy() {
	w.parent.markAsDestroyed()
	C.windowDestroy(w.nsWindow)
}

func (w *macosWebviewWindow) setHTML(html string) {
	// Convert HTML to C string
	cHTML := C.CString(html)
	// Render HTML
	C.windowRenderHTML(w.nsWindow, cHTML)
}

func (w *macosWebviewWindow) startDrag() error {
	C.startDrag(w.nsWindow)
	return nil
}

func (w *macosWebviewWindow) setMinimiseButtonState(state ButtonState) {
	C.setMinimiseButtonState(w.nsWindow, C.int(state))
}

func (w *macosWebviewWindow) setMaximiseButtonState(state ButtonState) {
	C.setMaximiseButtonState(w.nsWindow, C.int(state))
}

func (w *macosWebviewWindow) setCloseButtonState(state ButtonState) {
	C.setCloseButtonState(w.nsWindow, C.int(state))
}

func (w *macosWebviewWindow) isIgnoreMouseEvents() bool {
	return bool(C.isIgnoreMouseEvents(w.nsWindow))
}

func (w *macosWebviewWindow) setIgnoreMouseEvents(ignore bool) {
	C.setIgnoreMouseEvents(w.nsWindow, C.bool(ignore))
}

func (w *macosWebviewWindow) cut() {
}

func (w *macosWebviewWindow) paste() {
}

func (w *macosWebviewWindow) copy() {
}

func (w *macosWebviewWindow) selectAll() {
}

func (w *macosWebviewWindow) undo() {
}

func (w *macosWebviewWindow) delete() {
}

func (w *macosWebviewWindow) redo() {
}
