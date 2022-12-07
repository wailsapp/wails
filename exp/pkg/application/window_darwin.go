//go:build darwin

package application

/*
#cgo CFLAGS: -mmacosx-version-min=10.10 -x objective-c
#cgo LDFLAGS: -framework Cocoa -framework WebKit

#include "application.h"
#include "window_delegate.h"
#include <stdlib.h>
#include "Cocoa/Cocoa.h"
#import <WebKit/WebKit.h>


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

	// Add NSView to window
	NSView* view = [[NSView alloc] initWithFrame:NSMakeRect(0, 0, width-1, height-1)];
	[view setAutoresizingMask:NSViewWidthSizable | NSViewHeightSizable];
	[window setContentView:view];

	// Embed wkwebview in window
	NSRect frame = NSMakeRect(0, 0, width, height);
	WKWebViewConfiguration* config = [[WKWebViewConfiguration alloc] init];
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
		[(NSWindow*)nsWindow setStyleMask:resizable ? NSTitledWindowMask | NSClosableWindowMask | NSMiniaturizableWindowMask | NSResizableWindowMask : NSTitledWindowMask | NSClosableWindowMask | NSMiniaturizableWindowMask];
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
		[delegate.webView setValue:@YES forKey:@"drawsTransparentBackground"];
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

func (w *macosWindow) execJS(js string) {
	C.windowExecJS(w.nsWindow, C.CString(js))
}

func (w *macosWindow) navigateToURL(url string) {
	C.navigationLoadURL(w.nsWindow, C.CString(url))
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

func (w *macosWindow) run() error {
	w.nsWindow = C.windowNew(C.int(w.options.Width), C.int(w.options.Height))
	w.setTitle(w.options.Title)
	w.setAlwaysOnTop(w.options.AlwaysOnTop)
	w.setResizable(!w.options.DisableResize)
	w.setMinSize(w.options.MinWidth, w.options.MinHeight)
	w.setMaxSize(w.options.MaxWidth, w.options.MaxHeight)
	if w.options.URL != "" {
		w.navigateToURL(w.options.URL)
	}
	if w.options.EnableDevTools {
		w.enableDevTools()
	}
	if w.options.Mac != nil {
		switch w.options.Mac.Backdrop {
		case options.MacBackdropTransparent:
			C.windowSetTransparent(w.nsWindow)
			C.webviewSetTransparent(w.nsWindow)
		case options.MacBackdropTranslucent:
			C.windowSetTranslucent(w.nsWindow)
			C.webviewSetTransparent(w.nsWindow)
		}
	}
	C.windowShow(w.nsWindow)
	return nil
}
