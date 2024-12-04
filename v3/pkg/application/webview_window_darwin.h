#ifndef WEBVIEW_WINDOW_DARWIN
#define WEBVIEW_WINDOW_DARWIN

#import <Cocoa/Cocoa.h>
#import <WebKit/WebKit.h>
#include "webview_responder_darwin.h"

@interface WebviewWindow : NSObject

@property(assign) NSWindow *w;
@property(assign) WKWebView *webView; // We already retain WKWebView since it's part of the Window.
@property(assign) WebviewResponder *responder;

- (WebviewWindow *) initAsWindow:(NSRect)contentRect styleMask:(NSUInteger)windowStyle backing:(NSBackingStoreType)bufferingType defer:(BOOL)deferCreation;
- (WebviewWindow *) initAsPanel:(NSRect)contentRect styleMask:(NSUInteger)windowStyle backing:(NSBackingStoreType)bufferingType defer:(BOOL)deferCreation;

@end

@interface WebviewWindowDelegate : NSObject <NSWindowDelegate, WKScriptMessageHandler, WKNavigationDelegate, WKURLSchemeHandler, NSDraggingDestination>

@property unsigned int windowId;
@property (retain) NSEvent* leftMouseEvent;
@property unsigned int invisibleTitleBarHeight;
@property BOOL showToolbarWhenFullscreen;
@property NSWindowStyleMask previousStyleMask; // Used to restore the window style mask when using frameless

- (void)handleLeftMouseUp:(NSWindow *)window;
- (void)handleLeftMouseDown:(NSEvent*)event;
- (void)startDrag:(WebviewWindow*)window;

@end

#endif /* WEBVIEW_WINDOW_DARWIN */
