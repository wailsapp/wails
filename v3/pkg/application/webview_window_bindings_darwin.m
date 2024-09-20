#include "webview_window_bindings_darwin.h"

// Create a new Window
void* windowNew(unsigned int id, int width, int height, bool fraudulentWebsiteWarningEnabled, bool frameless, bool enableDragAndDrop, struct WebviewPreferences preferences) {
	return windowOrPanelNew(true, id, width, height, fraudulentWebsiteWarningEnabled, frameless, enableDragAndDrop, preferences);
}

// Create a new Panel
void* panelNew(unsigned int id, int width, int height, bool fraudulentWebsiteWarningEnabled, bool frameless, bool enableDragAndDrop, struct WebviewPreferences preferences) {
	return windowOrPanelNew(false, id, width, height, fraudulentWebsiteWarningEnabled, frameless, enableDragAndDrop, preferences);
}

void* windowOrPanelNew(bool isWindow, unsigned int id, int width, int height, bool fraudulentWebsiteWarningEnabled, bool frameless, bool enableDragAndDrop, struct WebviewPreferences preferences) {
	NSWindowStyleMask styleMask = NSWindowStyleMaskTitled | NSWindowStyleMaskClosable | NSWindowStyleMaskMiniaturizable | NSWindowStyleMaskResizable;
	if (frameless) {
		styleMask = NSWindowStyleMaskBorderless | NSWindowStyleMaskResizable;
	}

	WebviewWindow* webviewWindow;
	if (isWindow) {
		webviewWindow = [[WebviewWindow alloc] initAsWindow:NSMakeRect(0, 0, width-1, height-1)
			styleMask:styleMask
			backing:NSBackingStoreBuffered
			defer:NO];
	} else {
		webviewWindow = [[WebviewWindow alloc] initAsPanel:NSMakeRect(0, 0, width-1, height-1)
			styleMask:styleMask
			backing:NSBackingStoreBuffered
			defer:NO];
	}
	
	NSWindow *window = webviewWindow.w;

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

	webviewWindow.webView = webView;
	return webviewWindow;
}

void printWindowStyle(void *window) {
	NSWindow* nsWindow = ((WebviewWindow *)window).w;
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
	NSWindow* nsWindow = ((WebviewWindow *)window).w;
	// Get delegate
	WebviewWindowDelegate* delegate = (WebviewWindowDelegate*)[nsWindow delegate];
	// Set height
	delegate.invisibleTitleBarHeight = height;
}

// Make NSWindow transparent
void windowSetTransparent(void* nsWindow) {
	[((WebviewWindow*)nsWindow).w setOpaque:NO];
	[((WebviewWindow*)nsWindow).w setBackgroundColor:[NSColor clearColor]];
}

void windowSetInvisibleTitleBar(void* nsWindow, unsigned int height) {
	NSWindow* window = ((WebviewWindow*)nsWindow).w;
	WebviewWindowDelegate* delegate = (WebviewWindowDelegate*)[window delegate];
	delegate.invisibleTitleBarHeight = height;
}


// Set the title of the NSWindow
void windowSetTitle(void* nsWindow, char* title) {
	NSString* nsTitle = [NSString stringWithUTF8String:title];
	[((WebviewWindow*)nsWindow).w setTitle:nsTitle];
	free(title);
}

// Set the size of the NSWindow
void windowSetSize(void* nsWindow, int width, int height) {
	// Set window size on main thread
	NSWindow* window = ((WebviewWindow*)nsWindow).w;
	NSSize contentSize = [window contentRectForFrameRect:NSMakeRect(0, 0, width, height)].size;
	[window setContentSize:contentSize];
	[window setFrame:NSMakeRect(window.frame.origin.x, window.frame.origin.y, width, height) display:YES animate:YES];
}

// Set NSWindow always on top
void windowSetAlwaysOnTop(void* nsWindow, bool alwaysOnTop) {
	// Set window always on top on main thread
	[((WebviewWindow*)nsWindow).w setLevel:alwaysOnTop ? NSFloatingWindowLevel : NSNormalWindowLevel];
}

void setNormalWindowLevel(void* nsWindow) { [((WebviewWindow*)nsWindow).w setLevel:NSNormalWindowLevel]; }
void setFloatingWindowLevel(void* nsWindow) { [((WebviewWindow*)nsWindow).w setLevel:NSFloatingWindowLevel];}
void setPopUpMenuWindowLevel(void* nsWindow) { [((WebviewWindow*)nsWindow).w setLevel:NSPopUpMenuWindowLevel]; }
void setMainMenuWindowLevel(void* nsWindow) { [((WebviewWindow*)nsWindow).w setLevel:NSMainMenuWindowLevel]; }
void setStatusWindowLevel(void* nsWindow) { [((WebviewWindow*)nsWindow).w setLevel:NSStatusWindowLevel]; }
void setModalPanelWindowLevel(void* nsWindow) { [((WebviewWindow*)nsWindow).w setLevel:NSModalPanelWindowLevel]; }
void setScreenSaverWindowLevel(void* nsWindow) { [((WebviewWindow*)nsWindow).w setLevel:NSScreenSaverWindowLevel]; }
void setTornOffMenuWindowLevel(void* nsWindow) { [((WebviewWindow*)nsWindow).w setLevel:NSTornOffMenuWindowLevel]; }

// Load URL in NSWindow
void navigationLoadURL(void* nsWindow, char* url) {
	// Load URL on main thread
	NSURL* nsURL = [NSURL URLWithString:[NSString stringWithUTF8String:url]];
	NSURLRequest* request = [NSURLRequest requestWithURL:nsURL];
	WebviewWindow *webviewWindow = (WebviewWindow*)nsWindow;
	[webviewWindow.webView loadRequest:request];
	free(url);
}

// Set NSWindow resizable
void windowSetResizable(void* nsWindow, bool resizable) {
	// Set window resizable on main thread
	NSWindow* window = ((WebviewWindow*)nsWindow).w;
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
	NSWindow* window = ((WebviewWindow*)nsWindow).w;
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
	NSWindow* window = ((WebviewWindow*)nsWindow).w;
	NSSize contentSize = [window contentRectForFrameRect:NSMakeRect(0, 0, size.width, size.height)].size;
	[window setContentMaxSize:contentSize];
	[window setMaxSize:size];
}

// windowZoomReset
void windowZoomReset(void* nsWindow) {
	WebviewWindow *webviewWindow = (WebviewWindow*)nsWindow;
	[webviewWindow.webView setMagnification:1.0];
}

// windowZoomSet
void windowZoomSet(void* nsWindow, double zoom) {
	WebviewWindow *webviewWindow = (WebviewWindow*)nsWindow;
	// Reset zoom
	[webviewWindow.webView setMagnification:zoom];
}

// windowZoomGet
float windowZoomGet(void* nsWindow) {
	WebviewWindow *webviewWindow = (WebviewWindow*)nsWindow;
	// Get zoom
	return [webviewWindow.webView magnification];
}

// windowZoomIn
void windowZoomIn(void* nsWindow) {
	WebviewWindow *webviewWindow = (WebviewWindow*)nsWindow;
	// Zoom in
	[webviewWindow.webView setMagnification:webviewWindow.webView.magnification + 0.05];
}

// windowZoomOut
void windowZoomOut(void* nsWindow) {
	WebviewWindow *webviewWindow = (WebviewWindow*)nsWindow;
	// Zoom out
	if( webviewWindow.webView.magnification > 1.05 ) {
		[webviewWindow.webView setMagnification:webviewWindow.webView.magnification - 0.05];
	} else {
		[webviewWindow.webView setMagnification:1.0];
	}
}

// set the window position relative to the screen
void windowSetRelativePosition(void* nsWindow, int x, int y) {
	NSWindow* window = ((WebviewWindow*)nsWindow).w;
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
	WebviewWindow *webviewWindow = (WebviewWindow*)nsWindow;
	[webviewWindow.webView evaluateJavaScript:[NSString stringWithUTF8String:js] completionHandler:nil];
	free((void*)js);
}

// Make NSWindow backdrop translucent
void windowSetTranslucent(void* nsWindow) {
	// Get window
	NSWindow* window = ((WebviewWindow*)nsWindow).w;

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
	WebviewWindow *webviewWindow = (WebviewWindow*)nsWindow;
	// Set webview background transparent
	[webviewWindow.webView setValue:@NO forKey:@"drawsBackground"];
}

// Set webview background colour
void webviewSetBackgroundColour(void* nsWindow, int r, int g, int b, int alpha) {
	WebviewWindow *webviewWindow = (WebviewWindow*)nsWindow;
	// Set webview background color
	[webviewWindow.webView setValue:[NSColor colorWithRed:r/255.0 green:g/255.0 blue:b/255.0 alpha:alpha/255.0] forKey:@"backgroundColor"];
}

// Set the window background colour
void windowSetBackgroundColour(void* nsWindow, int r, int g, int b, int alpha) {
	[((WebviewWindow*)nsWindow).w setBackgroundColor:[NSColor colorWithRed:r/255.0 green:g/255.0 blue:b/255.0 alpha:alpha/255.0]];
}

bool windowIsMaximised(void* nsWindow) {
	return [((WebviewWindow*)nsWindow).w isZoomed];
}

bool windowIsFullscreen(void* nsWindow) {
	return [((WebviewWindow*)nsWindow).w styleMask] & NSWindowStyleMaskFullScreen;
}

bool windowIsMinimised(void* nsWindow) {
	return [((WebviewWindow*)nsWindow).w isMiniaturized];
}

bool windowIsFocused(void* nsWindow) {
	return [((WebviewWindow*)nsWindow).w isKeyWindow];
}

// Set Window fullscreen
void windowFullscreen(void* nsWindow) {
	if( windowIsFullscreen(nsWindow) ) {
		return;
	}
	dispatch_async(dispatch_get_main_queue(), ^{
		[((WebviewWindow*)nsWindow).w toggleFullScreen:nil];
	});}

void windowUnFullscreen(void* nsWindow) {
	if( !windowIsFullscreen(nsWindow) ) {
		return;
	}
	dispatch_async(dispatch_get_main_queue(), ^{
		[((WebviewWindow*)nsWindow).w toggleFullScreen:nil];
	});
}

// restore window to normal size
void windowRestore(void* nsWindow) {
	// If window is fullscreen
	if([((WebviewWindow*)nsWindow).w styleMask] & NSWindowStyleMaskFullScreen) {
		[((WebviewWindow*)nsWindow).w toggleFullScreen:nil];
	}
	// If window is maximised
	if([((WebviewWindow*)nsWindow).w isZoomed]) {
		[((WebviewWindow*)nsWindow).w zoom:nil];
	}
	// If window in minimised
	if([((WebviewWindow*)nsWindow).w isMiniaturized]) {
		[((WebviewWindow*)nsWindow).w deminiaturize:nil];
	}
}

// disable window fullscreen button
void setFullscreenButtonEnabled(void* nsWindow, bool enabled) {
	NSButton *fullscreenButton = [((WebviewWindow*)nsWindow).w standardWindowButton:NSWindowZoomButton];
	fullscreenButton.enabled = enabled;
}

// Set the titlebar style
void windowSetTitleBarAppearsTransparent(void* nsWindow, bool transparent) {
	if( transparent ) {
		[((WebviewWindow*)nsWindow).w setTitlebarAppearsTransparent:true];
	} else {
		[((WebviewWindow*)nsWindow).w setTitlebarAppearsTransparent:false];
	}
}

// Set window fullsize content view
void windowSetFullSizeContent(void* nsWindow, bool fullSize) {
	if( fullSize ) {
		[((WebviewWindow*)nsWindow).w setStyleMask:[((WebviewWindow*)nsWindow).w styleMask] | NSWindowStyleMaskFullSizeContentView];
	} else {
		[((WebviewWindow*)nsWindow).w setStyleMask:[((WebviewWindow*)nsWindow).w styleMask] & ~NSWindowStyleMaskFullSizeContentView];
	}
}

// Set Hide Titlebar
void windowSetHideTitleBar(void* nsWindow, bool hideTitlebar) {
	if( hideTitlebar ) {
		[((WebviewWindow*)nsWindow).w setStyleMask:[((WebviewWindow*)nsWindow).w styleMask] & ~NSWindowStyleMaskTitled];
	} else {
		[((WebviewWindow*)nsWindow).w setStyleMask:[((WebviewWindow*)nsWindow).w styleMask] | NSWindowStyleMaskTitled];
	}
}

// Set Hide Title in Titlebar
void windowSetHideTitle(void* nsWindow, bool hideTitle) {
	if( hideTitle ) {
		[((WebviewWindow*)nsWindow).w setTitleVisibility:NSWindowTitleHidden];
	} else {
		[((WebviewWindow*)nsWindow).w setTitleVisibility:NSWindowTitleVisible];
	}
}

// Set Window use toolbar
void windowSetUseToolbar(void* nsWindow, bool useToolbar) {
	NSWindow* window = ((WebviewWindow*)nsWindow).w;
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
	NSWindow* window = ((WebviewWindow*)nsWindow).w;

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
	NSToolbar* toolbar = [((WebviewWindow*)nsWindow).w toolbar];
	if( toolbar == nil ) {
		return;
	}
	[toolbar setShowsBaselineSeparator:!hideSeparator];
}

// Configure the toolbar auto-hide feature
void windowSetShowToolbarWhenFullscreen(void* window, bool setting) {
	NSWindow* nsWindow = ((WebviewWindow *)window).w;
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
	[((WebviewWindow*)nsWindow).w setAppearance:[NSAppearance appearanceNamed:appearanceNameString]];

	free((void*)appearanceName);
}

// Center window on current monitor
void windowCenter(void* nsWindow) {
	[((WebviewWindow*)nsWindow).w center];
}

// Get the current size of the window
void windowGetSize(void* nsWindow, int* width, int* height) {
	NSRect frame = [((WebviewWindow*)nsWindow).w frame];
	*width = frame.size.width;
	*height = frame.size.height;
}

// Get window width
int windowGetWidth(void* nsWindow) {
	return [((WebviewWindow*)nsWindow).w frame].size.width;
}

// Get window height
int windowGetHeight(void* nsWindow) {
	return [((WebviewWindow*)nsWindow).w frame].size.height;
}

// Get window position
void windowGetRelativePosition(void* nsWindow, int* x, int* y) {
	NSWindow* window = ((WebviewWindow*)nsWindow).w;
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
	NSRect frame = [((WebviewWindow*)nsWindow).w frame];
	*x = frame.origin.x;
	*y = frame.origin.y;
}

void windowSetPosition(void* nsWindow, int x, int y) {
	NSWindow *window = ((WebviewWindow*)nsWindow).w;
	
	NSRect frame = [window frame];
	frame.origin.x = x;
	frame.origin.y = y;
	[window setFrame:frame display:YES];
}

// Destroy window
void windowDestroy(void* nsWindow) {
	[((WebviewWindow*)nsWindow).w close];
}

// Remove drop shadow from window
void windowSetShadow(void* nsWindow, bool hasShadow) {
	[((WebviewWindow*)nsWindow).w setHasShadow:hasShadow];
}


// windowClose closes the current window
void windowClose(void *window) {
	[((WebviewWindow*)window).w close];
}

// windowZoom
void windowZoom(void *window) {
	[((WebviewWindow*)window).w zoom:nil];
}

// webviewRenderHTML renders the given HTML
void windowRenderHTML(void *window, const char *html) {
	WebviewWindow *webviewWindow = (WebviewWindow *)window;
	// get window delegate
	WebviewWindowDelegate* windowDelegate = (WebviewWindowDelegate*)[webviewWindow.w delegate];
	// render html
	[webviewWindow.webView loadHTMLString:[NSString stringWithUTF8String:html] baseURL:nil];
}

void windowInjectCSS(void *window, const char *css) {
	WebviewWindow *webviewWindow = (WebviewWindow *)window;	
	// inject css
	[webviewWindow.webView evaluateJavaScript:[NSString stringWithFormat:@"(function() { var style = document.createElement('style'); style.appendChild(document.createTextNode('%@')); document.head.appendChild(style); })();", [NSString stringWithUTF8String:css]] completionHandler:nil];
	free((void*)css);
}

void windowMinimise(void *window) {
	[((WebviewWindow*)window).w miniaturize:nil];
}

// zoom maximizes the window to the screen dimensions
void windowMaximise(void *window) {
	[((WebviewWindow*)window).w zoom:nil];
}

bool isFullScreen(void *window) {
	NSWindow* nsWindow = ((WebviewWindow *)window).w;
    long mask = [nsWindow styleMask];
    return (mask & NSWindowStyleMaskFullScreen) == NSWindowStyleMaskFullScreen;
}

bool isVisible(void *window) {
	NSWindow* nsWindow = ((WebviewWindow *)window).w;
    return (nsWindow.occlusionState & NSWindowOcclusionStateVisible) == NSWindowOcclusionStateVisible;
}

// windowSetFullScreen
void windowSetFullScreen(void *window, bool fullscreen) {
	if (isFullScreen(window)) {
		return;
	}
	NSWindow* nsWindow = ((WebviewWindow *)window).w;
	windowSetMaxSize(nsWindow, 0, 0);
	windowSetMinSize(nsWindow, 0, 0);
	[nsWindow toggleFullScreen:nil];
}

// windowUnminimise
void windowUnminimise(void *window) {
	[((WebviewWindow*)window).w deminiaturize:nil];
}

// windowUnmaximise
void windowUnmaximise(void *window) {
	[((WebviewWindow*)window).w zoom:nil];
}

void windowDisableSizeConstraints(void *window) {
	NSWindow* nsWindow = ((WebviewWindow *)window).w;
	// disable size constraints
	[nsWindow setContentMinSize:CGSizeZero];
	[nsWindow setContentMaxSize:CGSizeZero];
}

void windowShow(void *window) {
	[((WebviewWindow*)window).w makeKeyAndOrderFront:nil];
}

void windowHide(void *window) {
	[((WebviewWindow*)window).w orderOut:nil];
}

// setButtonState sets the state of the given button
// 0 = enabled
// 1 = disabled
// 2 = hidden
void setButtonState(void *button, int state) {
	if (button == nil) {
		return;
	}
	NSButton *nsbutton = (NSButton*)button;
	nsbutton.hidden = state == 2;
	nsbutton.enabled = state != 1;
}

// setMinimiseButtonState sets the minimise button state
void setMinimiseButtonState(void *window, int state) {
	NSWindow* nsWindow = ((WebviewWindow *)window).w;
	NSButton *minimiseButton = [nsWindow standardWindowButton:NSWindowMiniaturizeButton];
	setButtonState(minimiseButton, state);
}

// setMaximiseButtonState sets the maximise button state
void setMaximiseButtonState(void *window, int state) {
	NSWindow* nsWindow = ((WebviewWindow *)window).w;
	NSButton *maximiseButton = [nsWindow standardWindowButton:NSWindowZoomButton];
	setButtonState(maximiseButton, state);
}

// setCloseButtonState sets the close button state
void setCloseButtonState(void *window, int state) {
	NSWindow* nsWindow = ((WebviewWindow *)window).w;
	NSButton *closeButton = [nsWindow standardWindowButton:NSWindowCloseButton];
	setButtonState(closeButton, state);
}

// windowShowMenu opens an NSMenu at the given coordinates
void windowShowMenu(void *window, void *menu, int x, int y) {
	NSMenu* nsMenu = (NSMenu*)menu;
	WKWebView* webView = ((WebviewWindow*)window).webView;
	NSPoint point = NSMakePoint(x, y);
	[nsMenu popUpMenuPositioningItem:nil atLocation:point inView:webView];
}

// Make the given window frameless
void windowSetFrameless(void *window, bool frameless) {
	NSWindow* nsWindow = ((WebviewWindow *)window).w;
	// set the window style to be frameless
	if (frameless) {
		[nsWindow setStyleMask:([nsWindow styleMask] | NSWindowStyleMaskFullSizeContentView)];
	} else {
		[nsWindow setStyleMask:([nsWindow styleMask] & ~NSWindowStyleMaskFullSizeContentView)];
	}
}

void startDrag(void *window) {
	WebviewWindow *webviewWindow = (WebviewWindow*)window;
	// Get delegate
	WebviewWindowDelegate* windowDelegate = (WebviewWindowDelegate*)[webviewWindow.w delegate];
	// start drag
	[windowDelegate startDrag:webviewWindow];
}

// Credit: https://stackoverflow.com/q/33319295
void windowPrint(void *window) {
#if MAC_OS_X_VERSION_MAX_ALLOWED >= 110000
	// Check if macOS 11.0 or newer
	if (@available(macOS 11.0, *)) {
		WebviewWindow *webviewWindow = (WebviewWindow*)window;
		NSWindow* nsWindow = webviewWindow.w;

		WebviewWindowDelegate* windowDelegate = (WebviewWindowDelegate*)[nsWindow delegate];
		WKWebView* webView = webviewWindow.webView;

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
	NSWindow* nsWindow = ((WebviewWindow*)window).w;
	[nsWindow setIgnoresMouseEvents:!enabled];
}

void windowSetEnabled(void *window, bool enabled) {
	// TODO: Implement
}

void windowFocus(void *window) {
	NSWindow* nsWindow = ((WebviewWindow *)window).w;
	// If the current application is not active, activate it
	if (![[NSApplication sharedApplication] isActive]) {
		[[NSApplication sharedApplication] activateIgnoringOtherApps:YES];
	}
	[nsWindow makeKeyAndOrderFront:nil];
	[nsWindow makeKeyWindow];
}

bool isIgnoreMouseEvents(void *nsWindow) {
    NSWindow* window = ((WebviewWindow *)nsWindow).w;
    return [window ignoresMouseEvents];
}

void setIgnoreMouseEvents(void *nsWindow, bool ignore) {
    NSWindow* window = ((WebviewWindow *)nsWindow).w;
    [window setIgnoresMouseEvents:ignore];
}
