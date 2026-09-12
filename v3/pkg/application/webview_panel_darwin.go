//go:build darwin && !ios && !server

package application

/*
#cgo CFLAGS: -mmacosx-version-min=10.13 -x objective-c
#cgo LDFLAGS: -framework Cocoa -framework WebKit

#include <stdlib.h>
#include "Cocoa/Cocoa.h"
#import "webview_window_darwin.h"
#import <WebKit/WebKit.h>

// Create a new WebviewPanel
void* panelNew(unsigned int panelId, unsigned int windowId, void* parentWindow, int x, int y, int width, int height, bool transparent) {
	NSWindow<WailsWebviewWindow>* window = (NSWindow<WailsWebviewWindow>*)parentWindow;
	NSView* contentView = [window contentView];

	// Calculate frame (macOS uses bottom-left origin)
	NSRect contentBounds = [contentView bounds];
	NSRect frame = NSMakeRect(x, contentBounds.size.height - y - height, width, height);

	// Create WKWebView configuration
	WKWebViewConfiguration* config = [[WKWebViewConfiguration alloc] init];
	[config autorelease];

	config.suppressesIncrementalRendering = true;
	config.applicationNameForUserAgent = @"wails.io";

 // Reuse only the local asset handler; panel pages do not receive the Wails IPC bridge.
 id<WKURLSchemeHandler> handler = [window.webView.configuration urlSchemeHandlerForURLScheme:@"wails"];
 if (handler) [config setURLSchemeHandler:handler forURLScheme:@"wails"];
 WKWebView* webView = [[WKWebView alloc] initWithFrame:frame configuration:config];

 // Configure webview
	[webView setAutoresizingMask:NSViewMinYMargin];

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
	[webView stopLoading];
 [webView setNavigationDelegate:nil];
 [webView setUIDelegate:nil];
 [webView removeFromSuperview];
 [webView release];
}

// Set panel bounds
void panelSetBounds(void* panel, void* parentWindow, int x, int y, int width, int height) {
	WKWebView* webView = (WKWebView*)panel;
	NSWindow<WailsWebviewWindow>* window = (NSWindow<WailsWebviewWindow>*)parentWindow;
	NSView* contentView = [window contentView];

	// Calculate frame (macOS uses bottom-left origin)
	NSRect contentBounds = [contentView bounds];
	NSRect frame = NSMakeRect(x, contentBounds.size.height - y - height, width, height);

	[webView setFrame:frame];
}

// Get panel bounds
void panelGetBounds(void* panel, void* parentWindow, int* x, int* y, int* width, int* height) {
	WKWebView* webView = (WKWebView*)panel;
	NSWindow<WailsWebviewWindow>* window = (NSWindow<WailsWebviewWindow>*)parentWindow;
	NSView* contentView = [window contentView];

	NSRect frame = [webView frame];
	NSRect contentBounds = [contentView bounds];

	*x = (int)frame.origin.x;
	*y = (int)(contentBounds.size.height - frame.origin.y - frame.size.height);
	*width = (int)frame.size.width;
	*height = (int)frame.size.height;
}

// Place panels in ascending z-index order, above the main webview.
static void panelRaise(void* panel) {
 WKWebView* webView = (WKWebView*)panel;
 NSView* parent = webView.superview;
 [parent addSubview:webView positioned:NSWindowAbove relativeTo:nil];
}

// Navigate to URL
void panelLoadURL(void* panel, const char* url) {
	WKWebView* webView = (WKWebView*)panel;
	NSURL* nsURL = [NSURL URLWithString:[NSString stringWithUTF8String:url]];
	if (!nsURL) return;
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

// Focus panel
void panelFocus(void* panel) {
	WKWebView* webView = (WKWebView*)panel;
	[[webView window] makeFirstResponder:webView];
}

// Check if focused
bool panelIsFocused(void* panel) {
	WKWebView* webView = (WKWebView*)panel;
	NSWindow* window = [webView window];
	NSResponder* responder = [window firstResponder];
 return responder == webView || ([responder isKindOfClass:[NSView class]] && [(NSView*)responder isDescendantOf:webView]);
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

static void panelConfigure(void* panel, const char* userAgent) {
 WKWebView* view = (WKWebView*)panel;
 if (userAgent && userAgent[0]) view.customUserAgent = [NSString stringWithUTF8String:userAgent];
}
static void panelLoadInitialURL(void* panel, const char* url, const char* headersJSON) {
 WKWebView* view = (WKWebView*)panel;
 NSURL* target = [NSURL URLWithString:[NSString stringWithUTF8String:url]];
 if (!target) return;
 NSMutableURLRequest* request = [NSMutableURLRequest requestWithURL:target];
 NSData* data = [[NSString stringWithUTF8String:headersJSON] dataUsingEncoding:NSUTF8StringEncoding];
 NSDictionary* headers = [NSJSONSerialization JSONObjectWithData:data options:0 error:nil];
 for (NSString* key in headers) [request setValue:headers[key] forHTTPHeaderField:key];
 [view loadRequest:request];
}
*/
import "C"
import (
	"encoding/json"
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

func (p *darwinPanelImpl) create() {
	options := p.panel.snapshotOptions()

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

	agent := C.CString(options.UserAgent)
	C.panelConfigure(p.webview, agent)
	C.free(unsafe.Pointer(agent))
	enabled := globalApplication.isDebugMode
	if options.DevToolsEnabled != nil {
		enabled = *options.DevToolsEnabled
	}
	configureEmbeddedPanelDevTools(p.webview, enabled)
	if options.URL != "" {
		headers, _ := json.Marshal(options.Headers)
		if options.Headers == nil {
			headers = []byte("{}")
		}
		uri, rawHeaders := C.CString(options.URL), C.CString(string(headers))
		C.panelLoadInitialURL(p.webview, uri, rawHeaders)
		C.free(unsafe.Pointer(uri))
		C.free(unsafe.Pointer(rawHeaders))
	}
	if enabled && options.OpenInspectorOnStartup {
		p.openDevTools()
	}
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

func (p *darwinPanelImpl) setZIndex(_ int) {
	panels := p.panel.sortedSiblings()
	for _, panel := range panels {
		panel.destroyedLock.RLock()
		native, ok := panel.impl.(*darwinPanelImpl)
		panel.destroyedLock.RUnlock()
		if ok && native.webview != nil {
			C.panelRaise(native.webview)
		}
	}
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
	openEmbeddedPanelDevTools(p.webview)
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
