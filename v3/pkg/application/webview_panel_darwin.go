//go:build darwin && !ios

package application

/*
#cgo CFLAGS: -mmacosx-version-min=10.13 -x objective-c
#cgo LDFLAGS: -framework Cocoa -framework WebKit

#include <stdlib.h>
#include "Cocoa/Cocoa.h"
#import <WebKit/WebKit.h>

// WebviewPanel delegate for handling messages
@interface WebviewPanelDelegate : NSObject <WKScriptMessageHandler, WKNavigationDelegate>
@property unsigned int panelId;
@property unsigned int windowId;
@property (assign) WKWebView* webView;
@end

@implementation WebviewPanelDelegate

- (void)userContentController:(WKUserContentController *)userContentController didReceiveScriptMessage:(WKScriptMessage *)message {
	// Handle messages from the panel's webview
	// For now, log them - in future this could route to Go
}

- (void)webView:(WKWebView *)webView didFinishNavigation:(WKNavigation *)navigation {
	// Navigation completed callback
	extern void panelNavigationCompleted(unsigned int windowId, unsigned int panelId);
	panelNavigationCompleted(self.windowId, self.panelId);
}

@end

// Create a new WebviewPanel
void* panelNew(unsigned int panelId, unsigned int windowId, void* parentWindow, int x, int y, int width, int height, bool transparent) {
	WebviewWindow* window = (WebviewWindow*)parentWindow;
	NSView* contentView = [window contentView];

	// Calculate frame (macOS uses bottom-left origin)
	NSRect contentBounds = [contentView bounds];
	NSRect frame = NSMakeRect(x, contentBounds.size.height - y - height, width, height);

	// Create WKWebView configuration
	WKWebViewConfiguration* config = [[WKWebViewConfiguration alloc] init];
	[config autorelease];

	config.suppressesIncrementalRendering = true;
	config.applicationNameForUserAgent = @"wails.io";

	// Setup user content controller
	WKUserContentController* userContentController = [WKUserContentController new];
	[userContentController autorelease];

	WebviewPanelDelegate* delegate = [[WebviewPanelDelegate alloc] init];
	[delegate autorelease];
	delegate.panelId = panelId;
	delegate.windowId = windowId;

	[userContentController addScriptMessageHandler:delegate name:@"external"];
	config.userContentController = userContentController;

	// Create the WKWebView
	WKWebView* webView = [[WKWebView alloc] initWithFrame:frame configuration:config];
	delegate.webView = webView;

	// Configure webview
	[webView setAutoresizingMask:NSViewNotSizable];

	if (transparent) {
		[webView setValue:@NO forKey:@"drawsBackground"];
	}

	// Add to parent window's content view
	[contentView addSubview:webView];

	return webView;
}

// Destroy a WebviewPanel
void panelDestroy(void* panel) {
	WKWebView* webView = (WKWebView*)panel;
	[webView removeFromSuperview];
	[webView release];
}

// Set panel bounds
void panelSetBounds(void* panel, void* parentWindow, int x, int y, int width, int height) {
	WKWebView* webView = (WKWebView*)panel;
	WebviewWindow* window = (WebviewWindow*)parentWindow;
	NSView* contentView = [window contentView];

	// Calculate frame (macOS uses bottom-left origin)
	NSRect contentBounds = [contentView bounds];
	NSRect frame = NSMakeRect(x, contentBounds.size.height - y - height, width, height);

	[webView setFrame:frame];
}

// Get panel bounds
void panelGetBounds(void* panel, void* parentWindow, int* x, int* y, int* width, int* height) {
	WKWebView* webView = (WKWebView*)panel;
	WebviewWindow* window = (WebviewWindow*)parentWindow;
	NSView* contentView = [window contentView];

	NSRect frame = [webView frame];
	NSRect contentBounds = [contentView bounds];

	*x = (int)frame.origin.x;
	*y = (int)(contentBounds.size.height - frame.origin.y - frame.size.height);
	*width = (int)frame.size.width;
	*height = (int)frame.size.height;
}

// Set panel z-index (bring to front or send to back)
// Note: This is a binary implementation - panels are either on top (zIndex > 0)
// or at the bottom (zIndex <= 0). Granular z-index ordering would require tracking
// all panels and repositioning them relative to each other using NSWindowOrderingMode.
void panelSetZIndex(void* panel, void* parentWindow, int zIndex) {
	WKWebView* webView = (WKWebView*)panel;
	WebviewWindow* window = (WebviewWindow*)parentWindow;
	NSView* contentView = [window contentView];

	if (zIndex > 0) {
		// Bring to front
		[webView removeFromSuperview];
		[contentView addSubview:webView positioned:NSWindowAbove relativeTo:nil];
	} else {
		// Send to back (but above main webview which is at index 0)
		[webView removeFromSuperview];
		[contentView addSubview:webView positioned:NSWindowBelow relativeTo:nil];
	}
}

// Navigate to URL
void panelLoadURL(void* panel, const char* url) {
	WKWebView* webView = (WKWebView*)panel;
	NSURL* nsURL = [NSURL URLWithString:[NSString stringWithUTF8String:url]];
	NSURLRequest* request = [NSURLRequest requestWithURL:nsURL];
	[webView loadRequest:request];
}

// Reload
void panelReload(void* panel) {
	WKWebView* webView = (WKWebView*)panel;
	[webView reload];
}

// Force reload (bypass cache)
void panelForceReload(void* panel) {
	WKWebView* webView = (WKWebView*)panel;
	[webView reloadFromOrigin];
}

// Show panel
void panelShow(void* panel) {
	WKWebView* webView = (WKWebView*)panel;
	[webView setHidden:NO];
}

// Hide panel
void panelHide(void* panel) {
	WKWebView* webView = (WKWebView*)panel;
	[webView setHidden:YES];
}

// Check if visible
bool panelIsVisible(void* panel) {
	WKWebView* webView = (WKWebView*)panel;
	return ![webView isHidden];
}

// Set zoom
void panelSetZoom(void* panel, double zoom) {
	WKWebView* webView = (WKWebView*)panel;
	[webView setMagnification:zoom];
}

// Get zoom
double panelGetZoom(void* panel) {
	WKWebView* webView = (WKWebView*)panel;
	return [webView magnification];
}

// Open DevTools (inspector)
void panelOpenDevTools(void* panel) {
	WKWebView* webView = (WKWebView*)panel;
	// Note: Opening inspector programmatically requires private API
	// This is a no-op for now - users can right-click -> Inspect Element if enabled
}

// Focus panel
void panelFocus(void* panel) {
	WKWebView* webView = (WKWebView*)panel;
	[[webView window] makeFirstResponder:webView];
}

// Check if focused
bool panelIsFocused(void* panel) {
	WKWebView* webView = (WKWebView*)panel;
	NSWindow* window = [webView window];
	return [window firstResponder] == webView;
}

// Set background color
void panelSetBackgroundColour(void* panel, int r, int g, int b, int a) {
	WKWebView* webView = (WKWebView*)panel;
	if (a == 0) {
		[webView setValue:@NO forKey:@"drawsBackground"];
	} else {
		[webView setValue:[NSColor colorWithRed:r/255.0 green:g/255.0 blue:b/255.0 alpha:a/255.0] forKey:@"backgroundColor"];
	}
}

*/
import "C"
import (
	"unsafe"
)

type darwinPanelImpl struct {
	panel          *WebviewPanel
	webview        unsafe.Pointer
	parentNSWindow unsafe.Pointer
}

func newPanelImpl(panel *WebviewPanel) webviewPanelImpl {
	parentWindow := panel.parent
	if parentWindow == nil || parentWindow.impl == nil {
		return nil
	}

	darwinParent, ok := parentWindow.impl.(*macosWebviewWindow)
	if !ok {
		return nil
	}

	return &darwinPanelImpl{
		panel:          panel,
		parentNSWindow: darwinParent.nsWindow,
	}
}

//export panelNavigationCompleted
func panelNavigationCompleted(windowID C.uint, panelID C.uint) {
	// Navigation completed callback - could be used for future functionality
	globalApplication.debug("panelNavigationCompleted", "windowID", uint(windowID), "panelID", uint(panelID))
}

func (p *darwinPanelImpl) create() {
	options := p.panel.options

	transparent := options.Transparent

	p.webview = C.panelNew(
		C.uint(p.panel.id),
		C.uint(p.panel.parent.id),
		p.parentNSWindow,
		C.int(options.X),
		C.int(options.Y),
		C.int(options.Width),
		C.int(options.Height),
		C.bool(transparent),
	)

	// Set background colour if not transparent
	if !transparent {
		C.panelSetBackgroundColour(
			p.webview,
			C.int(options.BackgroundColour.Red),
			C.int(options.BackgroundColour.Green),
			C.int(options.BackgroundColour.Blue),
			C.int(options.BackgroundColour.Alpha),
		)
	}

	// Set initial visibility
	if options.Visible != nil && !*options.Visible {
		C.panelHide(p.webview)
	}

	// Set zoom if specified
	if options.Zoom > 0 && options.Zoom != 1.0 {
		C.panelSetZoom(p.webview, C.double(options.Zoom))
	}

	// Navigate to initial URL
	if options.URL != "" {
		// TODO: Add support for custom headers when WKWebView supports it
		if len(options.Headers) > 0 {
			globalApplication.debug("[Panel-Darwin] Custom headers specified (not yet supported)",
				"panelID", p.panel.id,
				"headers", options.Headers)
		}

		url := C.CString(options.URL)
		C.panelLoadURL(p.webview, url)
		C.free(unsafe.Pointer(url))
	}

	// Note: markRuntimeLoaded() is called in panelNavigationCompleted callback
	// when the navigation completes
}

func (p *darwinPanelImpl) destroy() {
	if p.webview != nil {
		C.panelDestroy(p.webview)
		p.webview = nil
	}
}

func (p *darwinPanelImpl) setBounds(bounds Rect) {
	if p.webview == nil {
		return
	}
	C.panelSetBounds(
		p.webview,
		p.parentNSWindow,
		C.int(bounds.X),
		C.int(bounds.Y),
		C.int(bounds.Width),
		C.int(bounds.Height),
	)
}

func (p *darwinPanelImpl) bounds() Rect {
	if p.webview == nil {
		return Rect{}
	}
	var x, y, width, height C.int
	C.panelGetBounds(p.webview, p.parentNSWindow, &x, &y, &width, &height)
	return Rect{
		X:      int(x),
		Y:      int(y),
		Width:  int(width),
		Height: int(height),
	}
}

func (p *darwinPanelImpl) setZIndex(zIndex int) {
	if p.webview == nil {
		return
	}
	C.panelSetZIndex(p.webview, p.parentNSWindow, C.int(zIndex))
}

func (p *darwinPanelImpl) setURL(url string) {
	if p.webview == nil {
		return
	}
	urlStr := C.CString(url)
	defer C.free(unsafe.Pointer(urlStr))
	C.panelLoadURL(p.webview, urlStr)
}

func (p *darwinPanelImpl) reload() {
	if p.webview == nil {
		return
	}
	C.panelReload(p.webview)
}

func (p *darwinPanelImpl) forceReload() {
	if p.webview == nil {
		return
	}
	C.panelForceReload(p.webview)
}

func (p *darwinPanelImpl) show() {
	if p.webview == nil {
		return
	}
	C.panelShow(p.webview)
}

func (p *darwinPanelImpl) hide() {
	if p.webview == nil {
		return
	}
	C.panelHide(p.webview)
}

func (p *darwinPanelImpl) isVisible() bool {
	if p.webview == nil {
		return false
	}
	return bool(C.panelIsVisible(p.webview))
}

func (p *darwinPanelImpl) setZoom(zoom float64) {
	if p.webview == nil {
		return
	}
	C.panelSetZoom(p.webview, C.double(zoom))
}

func (p *darwinPanelImpl) getZoom() float64 {
	if p.webview == nil {
		return 1.0
	}
	return float64(C.panelGetZoom(p.webview))
}

func (p *darwinPanelImpl) openDevTools() {
	if p.webview == nil {
		return
	}
	C.panelOpenDevTools(p.webview)
}

func (p *darwinPanelImpl) focus() {
	if p.webview == nil {
		return
	}
	C.panelFocus(p.webview)
}

func (p *darwinPanelImpl) isFocused() bool {
	if p.webview == nil {
		return false
	}
	return bool(C.panelIsFocused(p.webview))
}
